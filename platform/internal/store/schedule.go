package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RA9/gamifydev/platform/internal/schedule"
)

// The cohort schedule: materializing a path into dated work, and answering
// "what is due today" for one learner.

// --- Lesson completion -------------------------------------------------------

// MarkLessonComplete records that a user finished a lesson (idempotent).
func (s *Store) MarkLessonComplete(ctx context.Context, userID, lessonID int64) error {
	if err := s.RequireLessonAvailable(ctx, userID, lessonID); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO lesson_progress (user_id, lesson_id) VALUES (?, ?)
		 ON CONFLICT(user_id, lesson_id) DO NOTHING`, userID, lessonID)
	return err
}

// LessonComplete reports whether a user has finished a lesson.
func (s *Store) LessonComplete(ctx context.Context, userID, lessonID int64) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM lesson_progress WHERE user_id = ? AND lesson_id = ?`,
		userID, lessonID).Scan(&n)
	return n > 0, err
}

// MarkLessonCompleteIfStepsDone completes a step-based lesson once every one of
// its steps is done. Called after a step completes so a learner never has to
// separately confirm a lesson they've finished by working through it.
//
// A lesson with no steps is left alone — that one is completed explicitly.
func (s *Store) MarkLessonCompleteIfStepsDone(ctx context.Context, userID, lessonID int64) (bool, error) {
	var total, done int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COUNT(p.step_id)
		FROM lesson_steps s
		LEFT JOIN step_progress p ON p.step_id = s.id AND p.user_id = ?
		WHERE s.lesson_id = ?`, userID, lessonID).Scan(&total, &done); err != nil {
		return false, err
	}
	if total == 0 || done < total {
		return false, nil
	}
	return true, s.MarkLessonComplete(ctx, userID, lessonID)
}

// --- Materializing -----------------------------------------------------------

// pathCoursesForPlan loads a path's courses with their lessons and checkpoint,
// in path order, shaped for the planner.
func (s *Store) pathCoursesForPlan(ctx context.Context, pathID int64) ([]schedule.Course, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.title, c.slug,
		       (SELECT a.id FROM assignments a
		         WHERE a.course_id = c.id AND a.required = 1 AND a.published = 1
		         ORDER BY a.sort LIMIT 1)
		FROM path_courses pc JOIN courses c ON c.id = pc.course_id
		WHERE pc.path_id = ? ORDER BY pc.sort`, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []schedule.Course
	for rows.Next() {
		var c schedule.Course
		var checkpoint sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Title, &c.Slug, &checkpoint); err != nil {
			return nil, err
		}
		if checkpoint.Valid {
			c.CheckpointID = checkpoint.Int64
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		lrows, err := s.db.QueryContext(ctx,
			`SELECT id, title FROM lessons WHERE course_id = ? ORDER BY sort`, out[i].ID)
		if err != nil {
			return nil, err
		}
		for lrows.Next() {
			var l schedule.Lesson
			if err := lrows.Scan(&l.ID, &l.Title); err != nil {
				lrows.Close()
				return nil, err
			}
			out[i].Lessons = append(out[i].Lessons, l)
		}
		lrows.Close()
		if err := lrows.Err(); err != nil {
			return nil, err
		}
		prows, err := s.db.QueryContext(ctx, `
			SELECT p.id, p.title
			FROM course_problems cp JOIN problems p ON p.id = cp.problem_id
			WHERE cp.course_id = ? AND cp.required = 1 AND p.published = 1
			ORDER BY cp.sort, p.id`, out[i].ID)
		if err != nil {
			return nil, err
		}
		for prows.Next() {
			var problem schedule.Problem
			if err := prows.Scan(&problem.ID, &problem.Title); err != nil {
				prows.Close()
				return nil, err
			}
			out[i].Problems = append(out[i].Problems, problem)
		}
		prows.Close()
		if err := prows.Err(); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MaterializeSchedule builds a cohort's plan if it has none. Returns how many
// items were written; zero means it was already scheduled.
//
// Idempotent twice over: the cohort's scheduled_at short-circuits the common
// case, and the unique indexes make a concurrent second run write nothing.
func (s *Store) MaterializeSchedule(ctx context.Context, cohortID int64) (int, error) {
	var pathID int64
	var startsOn string
	var scheduled sql.NullString
	if err := s.db.QueryRowContext(ctx,
		`SELECT path_id, starts_on, scheduled_at FROM cohorts WHERE id = ?`, cohortID).
		Scan(&pathID, &startsOn, &scheduled); err != nil {
		return 0, err
	}
	if scheduled.Valid {
		return 0, nil
	}
	courses, err := s.pathCoursesForPlan(ctx, pathID)
	if err != nil {
		return 0, err
	}
	start, err := time.Parse("2006-01-02", startsOn)
	if err != nil {
		start = time.Now().UTC()
	}
	items := schedule.Plan(courses, start)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck
	for _, it := range items {
		var lesson, assignment, problem any
		if it.LessonID != 0 {
			lesson = it.LessonID
		}
		if it.AssignmentID != 0 {
			assignment = it.AssignmentID
		}
		if it.ProblemID != 0 {
			problem = it.ProblemID
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO cohort_schedule
				(cohort_id, day_index, sprint, due_on, kind, course_id,
				 lesson_id, assignment_id, problem_id, sort)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			cohortID, it.DayIndex, it.Sprint, it.DueOn.Format("2006-01-02"), it.Kind,
			it.CourseID, lesson, assignment, problem, it.Sort); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE cohorts SET scheduled_at = datetime('now') WHERE id = ?`, cohortID); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(items), nil
}

// UnscheduledCohorts lists active cohorts whose plan has not been built.
func (s *Store) UnscheduledCohorts(ctx context.Context) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id FROM cohorts WHERE state = 'active' AND scheduled_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// --- Reading the schedule ----------------------------------------------------

// ScheduleItem is one dated piece of work, resolved for a specific learner.
type ScheduleItem struct {
	ID             int64
	DayIndex       int
	Sprint         int
	DueOn          string
	Kind           string
	CourseID       int64
	Course         string
	CourseSlug     string
	LessonID       sql.NullInt64
	Lesson         string
	LessonSlug     string
	AssignmentID   sql.NullInt64
	Assignment     string
	AssignmentSlug string
	ProblemID      sql.NullInt64
	Problem        string
	ProblemSlug    string
	Mode           string

	// Per-learner state.
	Done    bool
	Exempt  bool
	Overdue bool
}

// URL is where the learner goes to do this item.
func (i ScheduleItem) URL() string {
	switch i.Kind {
	case schedule.KindCheckpoint:
		return "/assignments/" + i.AssignmentSlug
	case schedule.KindProblem:
		return "/problems/" + i.ProblemSlug
	default:
		return "/courses/" + i.CourseSlug + "/" + i.LessonSlug
	}
}

// IsCheckpoint reports whether this item is a gating assignment.
func (i ScheduleItem) IsCheckpoint() bool { return i.Kind == schedule.KindCheckpoint }

// IsProblem reports whether this item is integrated guided practice.
func (i ScheduleItem) IsProblem() bool { return i.Kind == schedule.KindProblem }

// Practical reports whether pacing and prerequisite gates apply to this item.
func (i ScheduleItem) Practical() bool { return i.Mode == "practical" || i.Mode == "fun" }

const scheduleSelect = `
	SELECT s.id, s.day_index, s.sprint, s.due_on, s.kind, s.course_id,
	       COALESCE(c.title,''), COALESCE(c.slug,''),
	       s.lesson_id, COALESCE(l.title,''), COALESCE(l.slug,''),
	       s.assignment_id, COALESCE(a.title,''), COALESCE(a.slug,''),
	       s.problem_id, COALESCE(pr.title,''), COALESCE(pr.slug,''),
	       COALESCE(l.mode, a.mode, pr.mode, 'theoretical'),
	       CASE WHEN s.lesson_id IS NOT NULL AND lp.lesson_id IS NOT NULL THEN 1
	            WHEN s.assignment_id IS NOT NULL AND EXISTS (
	              SELECT 1 FROM passed_checkpoints pc
	              WHERE pc.assignment_id = s.assignment_id AND pc.user_id = ?
	            ) THEN 1
	            WHEN s.problem_id IS NOT NULL AND EXISTS (
	              SELECT 1 FROM problem_submissions ps
	              WHERE ps.problem_id = s.problem_id AND ps.user_id = ?
	                AND ps.verdict = 'accepted'
	            ) THEN 1
	            ELSE 0 END AS done,
	       CASE WHEN date(s.due_on) < date('now') THEN 1 ELSE 0 END AS overdue
	FROM cohort_schedule s
	LEFT JOIN courses c ON c.id = s.course_id
	LEFT JOIN lessons l ON l.id = s.lesson_id
	LEFT JOIN assignments a ON a.id = s.assignment_id
	LEFT JOIN problems pr ON pr.id = s.problem_id
	LEFT JOIN lesson_progress lp ON lp.lesson_id = s.lesson_id AND lp.user_id = ?
`

func (s *Store) scanSchedule(rows *sql.Rows, exempt map[string]bool) ([]ScheduleItem, error) {
	defer rows.Close()
	var out []ScheduleItem
	for rows.Next() {
		var it ScheduleItem
		var done, overdue int
		if err := rows.Scan(&it.ID, &it.DayIndex, &it.Sprint, &it.DueOn, &it.Kind, &it.CourseID,
			&it.Course, &it.CourseSlug, &it.LessonID, &it.Lesson, &it.LessonSlug,
			&it.AssignmentID, &it.Assignment, &it.AssignmentSlug,
			&it.ProblemID, &it.Problem, &it.ProblemSlug, &it.Mode, &done, &overdue); err != nil {
			return nil, err
		}
		it.Done, it.Overdue = done == 1, overdue == 1
		// Exemptions are a per-learner overlay on the shared cohort plan: the
		// item stays on the calendar, it just isn't owed by this learner.
		it.Exempt = exempt[it.CourseSlug]
		if it.Exempt {
			it.Overdue = false
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// exemptCourses returns the immutable exemption snapshot attached to the
// learner's enrollment in this cohort. It intentionally does not read the latest
// placement result: a later diagnostic must never rewrite an existing schedule.
func (s *Store) exemptCourses(ctx context.Context, cohortID, userID int64) (map[string]bool, error) {
	out := map[string]bool{}
	var raw string
	err := s.db.QueryRowContext(ctx, `
		SELECT placement_exemptions FROM enrollments
		WHERE user_id = ? AND cohort_id = ?
		  AND state IN ('active','paused','completed')
		ORDER BY id DESC LIMIT 1`, userID, cohortID).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return out, nil // legacy/no enrollment: nothing exempt
	}
	if err != nil {
		return nil, err
	}
	var slugs []string
	if err := json.Unmarshal([]byte(raw), &slugs); err != nil {
		return nil, fmt.Errorf("decode enrollment exemptions: %w", err)
	}
	for _, slug := range slugs {
		out[slug] = true
	}
	return out, nil
}

// DueOn returns a learner's scheduled work for one date.
func (s *Store) DueOn(ctx context.Context, cohortID, userID int64, day string) ([]ScheduleItem, error) {
	exempt, err := s.exemptCourses(ctx, cohortID, userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		scheduleSelect+` WHERE s.cohort_id = ? AND date(s.due_on) = date(?) ORDER BY s.sort`,
		userID, userID, userID, cohortID, day)
	if err != nil {
		return nil, err
	}
	return s.scanSchedule(rows, exempt)
}

// Overdue returns unfinished work whose due date has passed. Exempted items are
// filtered out entirely — a learner is never chased for a course they skipped.
func (s *Store) Overdue(ctx context.Context, cohortID, userID int64, limit int) ([]ScheduleItem, error) {
	exempt, err := s.exemptCourses(ctx, cohortID, userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		scheduleSelect+` WHERE s.cohort_id = ? AND date(s.due_on) < date('now')
		 ORDER BY s.due_on, s.sort LIMIT ?`,
		userID, userID, userID, cohortID, limit*4)
	if err != nil {
		return nil, err
	}
	all, err := s.scanSchedule(rows, exempt)
	if err != nil {
		return nil, err
	}
	var out []ScheduleItem
	for _, it := range all {
		if it.Done || it.Exempt {
			continue
		}
		out = append(out, it)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// CohortSchedule returns the whole plan for a cohort, resolved for a learner.
func (s *Store) CohortSchedule(ctx context.Context, cohortID, userID int64) ([]ScheduleItem, error) {
	exempt, err := s.exemptCourses(ctx, cohortID, userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		scheduleSelect+` WHERE s.cohort_id = ? ORDER BY s.day_index, s.sort`,
		userID, userID, userID, cohortID)
	if err != nil {
		return nil, err
	}
	return s.scanSchedule(rows, exempt)
}

// ScheduleProgress summarises a learner's standing against the plan.
type ScheduleProgress struct {
	Total, Done, Exempt, Overdue int
	DayIndex, TotalDays          int
	Sprint, TotalSprints         int
	StartsOn, EndsOn             string
}

// Pct is completion across everything actually owed by this learner.
func (p ScheduleProgress) Pct() int {
	owed := p.Total - p.Exempt
	if owed <= 0 {
		return 0
	}
	return p.Done * 100 / owed
}

// Owed is how many items this learner is actually responsible for.
func (p ScheduleProgress) Owed() int { return p.Total - p.Exempt }

// Progress computes a learner's standing against their cohort's plan.
func (s *Store) Progress(ctx context.Context, cohortID, userID int64) (ScheduleProgress, error) {
	var p ScheduleProgress
	items, err := s.CohortSchedule(ctx, cohortID, userID)
	if err != nil {
		return p, err
	}
	for _, it := range items {
		p.Total++
		switch {
		case it.Exempt:
			p.Exempt++
		case it.Done:
			p.Done++
		case it.Overdue:
			p.Overdue++
		}
		if it.DayIndex+1 > p.TotalDays {
			p.TotalDays = it.DayIndex + 1
		}
		if it.Sprint > p.TotalSprints {
			p.TotalSprints = it.Sprint
		}
	}
	if len(items) > 0 {
		p.StartsOn = items[0].DueOn
		p.EndsOn = items[len(items)-1].DueOn
	}
	// Which day the cohort is on today, clamped to the plan.
	_ = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(day_index), 0) FROM cohort_schedule
		WHERE cohort_id = ? AND date(due_on) <= date('now')`, cohortID).Scan(&p.DayIndex)
	p.Sprint = p.DayIndex/schedule.DaysPerSprint + 1
	return p, nil
}

// CompletionResult reports terminal transitions made by AdvanceCompletions.
type CompletionResult struct {
	Enrollments int
	Cohorts     int
}

// AdvanceCompletions transactionally completes active path enrollments once the
// cohort has reached its final scheduled date and every non-exempt item is done.
// A cohort completes after its last live learner finishes; dropped learners do
// not block the group. Re-running is an idempotent no-op.
func (s *Store) AdvanceCompletions(ctx context.Context) (CompletionResult, error) {
	var result CompletionResult
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback() //nolint:errcheck

	rows, err := tx.QueryContext(ctx, `
		SELECT e.id, e.user_id, e.cohort_id, e.placement_exemptions
		FROM enrollments e
		WHERE e.state = 'active' AND e.cohort_id IS NOT NULL
		  AND EXISTS (SELECT 1 FROM cohort_schedule cs WHERE cs.cohort_id = e.cohort_id)
		  AND date((SELECT MAX(cs.due_on) FROM cohort_schedule cs WHERE cs.cohort_id = e.cohort_id)) <= date('now')`)
	if err != nil {
		return result, err
	}
	type candidate struct {
		enrollmentID int64
		userID       int64
		cohortID     int64
		exemptions   string
	}
	var candidates []candidate
	for rows.Next() {
		var c candidate
		if err := rows.Scan(&c.enrollmentID, &c.userID, &c.cohortID, &c.exemptions); err != nil {
			rows.Close()
			return result, err
		}
		candidates = append(candidates, c)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return result, err
	}

	for _, c := range candidates {
		var incomplete int
		if err := tx.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM cohort_schedule cs
			JOIN courses course ON course.id = cs.course_id
			WHERE cs.cohort_id = ?
			  AND NOT EXISTS (
			    SELECT 1 FROM json_each(?) ex
			    WHERE CAST(ex.value AS TEXT) = course.slug
			  )
			  AND (
			    (cs.lesson_id IS NOT NULL AND NOT EXISTS (
			      SELECT 1 FROM lesson_progress lp
			      WHERE lp.user_id = ? AND lp.lesson_id = cs.lesson_id
			    ))
			    OR (cs.assignment_id IS NOT NULL AND NOT EXISTS (
			      SELECT 1 FROM passed_checkpoints pc
			      WHERE pc.user_id = ? AND pc.assignment_id = cs.assignment_id
			    ))
			    OR (cs.problem_id IS NOT NULL AND NOT EXISTS (
			      SELECT 1 FROM problem_submissions ps
			      WHERE ps.user_id = ? AND ps.problem_id = cs.problem_id AND ps.verdict = 'accepted'
			    ))
			  )`, c.cohortID, c.exemptions, c.userID, c.userID, c.userID).Scan(&incomplete); err != nil {
			return result, err
		}
		if incomplete != 0 {
			continue
		}
		res, err := tx.ExecContext(ctx, `
			UPDATE enrollments
			SET state = 'completed', reason = 'path requirements completed',
			    ended_at = datetime('now'), updated_at = datetime('now')
			WHERE id = ? AND state = 'active'`, c.enrollmentID)
		if err != nil {
			return result, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return result, err
		}
		if n == 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE cohort_members SET left_at = datetime('now')
			WHERE cohort_id = ? AND user_id = ? AND left_at IS NULL`, c.cohortID, c.userID); err != nil {
			return result, err
		}
		result.Enrollments++
	}

	cohortRows, err := tx.QueryContext(ctx, `
		SELECT c.id
		FROM cohorts c
		WHERE c.state = 'active'
		  AND EXISTS (SELECT 1 FROM cohort_schedule cs WHERE cs.cohort_id = c.id)
		  AND date((SELECT MAX(cs.due_on) FROM cohort_schedule cs WHERE cs.cohort_id = c.id)) <= date('now')
		  AND EXISTS (SELECT 1 FROM enrollments e WHERE e.cohort_id = c.id AND e.state = 'completed')
		  AND NOT EXISTS (
		    SELECT 1 FROM enrollments e
		    WHERE e.cohort_id = c.id AND e.state IN ('placed','active','paused')
		  )`)
	if err != nil {
		return result, err
	}
	var completableCohorts []int64
	for cohortRows.Next() {
		var cohortID int64
		if err := cohortRows.Scan(&cohortID); err != nil {
			cohortRows.Close()
			return result, err
		}
		completableCohorts = append(completableCohorts, cohortID)
	}
	cohortRows.Close()
	if err := cohortRows.Err(); err != nil {
		return result, err
	}

	for _, cohortID := range completableCohorts {
		res, err := tx.ExecContext(ctx, `
			UPDATE cohorts SET state = 'completed', completed_at = datetime('now')
			WHERE id = ? AND state = 'active'`, cohortID)
		if err != nil {
			return result, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return result, err
		}
		if n > 0 {
			if _, err := tx.ExecContext(ctx, `
				UPDATE cohort_members SET left_at = datetime('now')
				WHERE cohort_id = ? AND left_at IS NULL`, cohortID); err != nil {
				return result, err
			}
			result.Cohorts++
		}
	}

	if err := tx.Commit(); err != nil {
		return CompletionResult{}, err
	}
	return result, nil
}
