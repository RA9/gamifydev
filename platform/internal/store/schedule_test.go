package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/RA9/gamifydev/platform/internal/schedule"
)

// mkCourse creates a course with n lessons and optionally a required checkpoint,
// wired into a path at the given position.
func mkCourse(t *testing.T, st *Store, pathID int64, slug string, sort, lessons int, checkpoint bool) int64 {
	t.Helper()
	var cid int64
	if err := st.db.QueryRow(
		`INSERT INTO courses (slug, title, published) VALUES (?, ?, 1) RETURNING id`, slug, slug).
		Scan(&cid); err != nil {
		t.Fatalf("course: %v", err)
	}
	mustExec(t, st, `INSERT INTO path_courses (path_id, course_id, sort) VALUES (?, ?, ?)`, pathID, cid, sort)
	for i := 0; i < lessons; i++ {
		mustExec(t, st,
			`INSERT INTO lessons (course_id, slug, title, sort) VALUES (?, ?, ?, ?)`,
			cid, slug+"-l"+string(rune('a'+i)), "Lesson", i)
	}
	if checkpoint {
		mustExec(t, st,
			`INSERT INTO assignments (course_id, slug, title, required, published) VALUES (?, ?, ?, 1, 1)`,
			cid, slug+"-check", "Checkpoint")
	}
	return cid
}

func scheduledCohort(t *testing.T, st *Store) (cohortID, userID, pathID int64) {
	t.Helper()
	ctx := context.Background()
	pathID = mkPath(t, st, "frontend")
	mkCourse(t, st, pathID, "c1", 0, 3, true)
	mkCourse(t, st, pathID, "c2", 1, 2, false)

	// Start on a known Monday so weekday arithmetic is predictable.
	cohortID, err := st.CreateCohort(ctx, Cohort{
		PathID: pathID, Name: "C", TZBand: "europe_africa", StartsOn: "2026-09-07",
	})
	if err != nil {
		t.Fatalf("cohort: %v", err)
	}
	userID = enroll(t, st, "sched@x.com", pathID, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, cohortID, pathID, "europe_africa", 5); err != nil {
		t.Fatalf("fill: %v", err)
	}
	return
}

func TestMaterializeScheduleIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)

	n, err := st.MaterializeSchedule(ctx, cid)
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	// 3 lessons + 1 checkpoint + 2 lessons = 6 items.
	if n != 6 {
		t.Fatalf("wrote %d items, want 6", n)
	}

	// Re-running writes nothing — the job runs every 10 minutes.
	again, err := st.MaterializeSchedule(ctx, cid)
	if err != nil {
		t.Fatalf("re-materialize: %v", err)
	}
	if again != 0 {
		t.Fatalf("re-materialize wrote %d items, want 0", again)
	}
	items, _ := st.CohortSchedule(ctx, cid, uid)
	if len(items) != 6 {
		t.Fatalf("schedule has %d items after re-run, want 6", len(items))
	}

	// And it drops out of the unscheduled list.
	pending, _ := st.UnscheduledCohorts(ctx)
	for _, id := range pending {
		if id == cid {
			t.Fatal("a scheduled cohort is still listed as unscheduled")
		}
	}
}

func TestRequiredCourseProblemsAreScheduledAndCompleted(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cohortID, userID, _ := scheduledCohort(t, st)
	var courseID, problemID int64
	if err := st.db.QueryRow(`SELECT id FROM courses WHERE slug = 'c1'`).Scan(&courseID); err != nil {
		t.Fatalf("course: %v", err)
	}
	if err := st.db.QueryRow(`
		INSERT INTO problems (slug, title, language, workload_minutes, mode, published)
		VALUES ('guided-c', 'Guided C', 'c', 70, 'practical', 1) RETURNING id`).Scan(&problemID); err != nil {
		t.Fatalf("problem: %v", err)
	}
	mustExec(t, st, `INSERT INTO course_problems (course_id, problem_id, sort, required) VALUES (?, ?, 0, 1)`, courseID, problemID)
	if n, err := st.MaterializeSchedule(ctx, cohortID); err != nil || n != 7 {
		t.Fatalf("materialize with problem = %d, %v; want 7, nil", n, err)
	}
	items, err := st.CohortSchedule(ctx, cohortID, userID)
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	var problemItem *ScheduleItem
	for i := range items {
		if items[i].Kind == schedule.KindProblem {
			problemItem = &items[i]
			break
		}
	}
	if problemItem == nil || !problemItem.ProblemID.Valid || problemItem.ProblemID.Int64 != problemID || problemItem.URL() != "/problems/guided-c" || problemItem.Done {
		t.Fatalf("scheduled problem = %+v", problemItem)
	}
	mustExec(t, st, `INSERT INTO problem_submissions
		(problem_id, user_id, language, code, verdict) VALUES (?, ?, 'c', 'ok', 'accepted')`, problemID, userID)
	items, err = st.CohortSchedule(ctx, cohortID, userID)
	if err != nil {
		t.Fatalf("completed schedule: %v", err)
	}
	for i := range items {
		if items[i].Kind == schedule.KindProblem && !items[i].Done {
			t.Fatalf("accepted problem is not complete: %+v", items[i])
		}
	}
}

func TestScheduleSkipsWeekendsEndToEnd(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	items, _ := st.CohortSchedule(ctx, cid, uid)
	for _, it := range items {
		d, err := time.Parse("2006-01-02", it.DueOn)
		if err != nil {
			t.Fatalf("bad due date %q", it.DueOn)
		}
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			t.Fatalf("item due on a %s (%s)", d.Weekday(), it.DueOn)
		}
	}
}

func TestExemptionsAreAPerLearnerOverlay(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, pathID := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}

	// A second learner joins the same cohort with c1 in the immutable exemption
	// snapshot copied onto their enrollment.
	other := enroll(t, st, "exempt@x.com", pathID, "europe_africa")
	mustExec(t, st, `UPDATE enrollments SET placement_exemptions = '["c1"]' WHERE user_id = ?`, other)
	if _, err := st.TakeFromQueue(ctx, cid, pathID, "europe_africa", 1); err != nil {
		t.Fatalf("fill exempt learner: %v", err)
	}

	// A newer result with a different decision must not rewrite an active cohort.
	mustExec(t, st, `INSERT INTO assessment_attempts (id, user_id, assessment_id, expires_at, submitted_at)
		VALUES (900, ?, 1, datetime('now','+1 hour'), datetime('now'))`, other)
	mustExec(t, st, `INSERT INTO placement_results (attempt_id, user_id, foundations_required, passed, exemptions)
		VALUES (900, ?, 0, 1, '[]')`, other)

	mine, _ := st.CohortSchedule(ctx, cid, uid)
	theirs, _ := st.CohortSchedule(ctx, cid, other)

	// The shared plan is identical — one cohort, one calendar. Only the overlay
	// differs, which is what keeps the group on the same day.
	if len(mine) != len(theirs) {
		t.Fatalf("plans differ in length: %d vs %d", len(mine), len(theirs))
	}
	minExempt, theirExempt := 0, 0
	for i := range mine {
		if mine[i].ID != theirs[i].ID {
			t.Fatalf("item %d differs between learners", i)
		}
		if mine[i].Exempt {
			minExempt++
		}
		if theirs[i].Exempt {
			theirExempt++
		}
	}
	if minExempt != 0 {
		t.Fatalf("unexempted learner has %d exempt items", minExempt)
	}
	// c1 contributes 3 lessons + 1 checkpoint.
	if theirExempt != 4 {
		t.Fatalf("exempt learner has %d exempt items, want 4", theirExempt)
	}
}

func TestExemptWorkIsNeverOverdue(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pathID := mkPath(t, st, "frontend")
	mkCourse(t, st, pathID, "c1", 0, 2, false)
	// A cohort that started well in the past, so everything is overdue.
	cid, _ := st.CreateCohort(ctx, Cohort{
		PathID: pathID, Name: "old", TZBand: "europe_africa", StartsOn: "2020-01-06",
	})
	uid := enroll(t, st, "late@x.com", pathID, "europe_africa")
	mustExec(t, st, `UPDATE enrollments SET placement_exemptions = '["c1"]' WHERE user_id = ?`, uid)
	if _, err := st.TakeFromQueue(ctx, cid, pathID, "europe_africa", 1); err != nil {
		t.Fatalf("fill: %v", err)
	}
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}

	items, err := st.CohortSchedule(ctx, cid, uid)
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if len(items) != 2 || !items[0].Exempt || !items[1].Exempt {
		t.Fatalf("enrollment snapshot was not applied to schedule: %+v", items)
	}
	overdue, err := st.Overdue(ctx, cid, uid, 10)
	if err != nil || len(overdue) != 0 {
		t.Fatalf("exempt learner overdue = %d, %v; want 0, nil", len(overdue), err)
	}
	progress, err := st.Progress(ctx, cid, uid)
	if err != nil || progress.Exempt != 2 || progress.Owed() != 0 || progress.Overdue != 0 {
		t.Fatalf("exempt learner progress = %+v, %v", progress, err)
	}

	// Even a newer placement row that removes the exemption cannot mutate the
	// active schedule, overdue queue, or progress denominator.
	mustExec(t, st, `INSERT INTO assessment_attempts (id, user_id, assessment_id, expires_at, submitted_at)
		VALUES (901, ?, 1, datetime('now'), datetime('now'))`, uid)
	mustExec(t, st, `INSERT INTO placement_results (attempt_id, user_id, foundations_required, passed, exemptions)
		VALUES (901, ?, 0, 1, '[]')`, uid)
	items, err = st.CohortSchedule(ctx, cid, uid)
	if err != nil || len(items) != 2 || !items[0].Exempt || !items[1].Exempt {
		t.Fatalf("newer result changed schedule: %+v, %v", items, err)
	}
	overdue, err = st.Overdue(ctx, cid, uid, 10)
	if err != nil || len(overdue) != 0 {
		t.Fatalf("newer result changed overdue work: %d, %v", len(overdue), err)
	}
	progress, err = st.Progress(ctx, cid, uid)
	if err != nil || progress.Exempt != 2 || progress.Owed() != 0 || progress.Overdue != 0 {
		t.Fatalf("newer result changed progress: %+v, %v", progress, err)
	}
}

func TestCompletionClearsDueAndCountsProgress(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}

	items, _ := st.CohortSchedule(ctx, cid, uid)
	var firstLesson int64
	for _, it := range items {
		if it.LessonID.Valid {
			firstLesson = it.LessonID.Int64
			break
		}
	}
	before, _ := st.Progress(ctx, cid, uid)
	if before.Done != 0 {
		t.Fatalf("progress starts at %d done, want 0", before.Done)
	}
	if before.Owed() != 6 {
		t.Fatalf("owed = %d, want 6", before.Owed())
	}

	if err := st.MarkLessonComplete(ctx, uid, firstLesson); err != nil {
		t.Fatalf("complete: %v", err)
	}
	after, _ := st.Progress(ctx, cid, uid)
	if after.Done != 1 {
		t.Fatalf("progress = %d done, want 1", after.Done)
	}
	if after.Pct() == before.Pct() {
		t.Fatal("percentage did not move after completing a lesson")
	}
	// Idempotent.
	if err := st.MarkLessonComplete(ctx, uid, firstLesson); err != nil {
		t.Fatalf("re-complete: %v", err)
	}
	again, _ := st.Progress(ctx, cid, uid)
	if again.Done != 1 {
		t.Fatalf("re-completing counted twice: %d", again.Done)
	}
}

func TestStepLessonAutoCompletes(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "stepper@x.com")
	mustExec(t, st, `INSERT INTO courses (id, slug, title) VALUES (7, 'c', 'C')`)
	mustExec(t, st, `INSERT INTO lessons (id, course_id, slug, title) VALUES (7, 7, 'l', 'L')`)
	mustExec(t, st, `INSERT INTO lesson_steps (id, lesson_id, sort) VALUES (70, 7, 0), (71, 7, 1)`)

	// Half done — the lesson must not be marked complete yet.
	if err := st.MarkStepComplete(ctx, uid, 70); err != nil {
		t.Fatalf("step: %v", err)
	}
	done, err := st.MarkLessonCompleteIfStepsDone(ctx, uid, 7)
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if done {
		t.Fatal("lesson completed with a step still outstanding")
	}

	// Last step finishes the lesson without a second confirmation.
	if err := st.MarkStepComplete(ctx, uid, 71); err != nil {
		t.Fatalf("step: %v", err)
	}
	done, err = st.MarkLessonCompleteIfStepsDone(ctx, uid, 7)
	if err != nil || !done {
		t.Fatalf("MarkLessonCompleteIfStepsDone = %v, %v; want true, nil", done, err)
	}
	if ok, _ := st.LessonComplete(ctx, uid, 7); !ok {
		t.Fatal("lesson not recorded complete")
	}

	// A lesson with no steps is never auto-completed — it needs the explicit mark.
	mustExec(t, st, `INSERT INTO lessons (id, course_id, slug, title) VALUES (8, 7, 'l2', 'L2')`)
	if done, _ := st.MarkLessonCompleteIfStepsDone(ctx, uid, 8); done {
		t.Fatal("a stepless lesson was auto-completed")
	}
}

func TestScheduledPracticalWorkHonorsReleaseAndPrerequisites(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	mustExec(t, st, `UPDATE lessons SET mode = 'practical' WHERE course_id = (SELECT id FROM courses WHERE slug = 'c1')`)
	mustExec(t, st, `UPDATE cohort_schedule SET due_on = date('now', '+1 day') WHERE cohort_id = ?`, cid)

	var first, second int64
	if err := st.db.QueryRow(`
		SELECT MIN(lesson_id), MAX(lesson_id) FROM (
		  SELECT lesson_id FROM cohort_schedule
		  WHERE cohort_id = ? AND lesson_id IS NOT NULL ORDER BY day_index, sort LIMIT 2
		)`, cid).Scan(&first, &second); err != nil {
		t.Fatalf("scheduled lessons: %v", err)
	}
	if err := st.RequireLessonAvailable(ctx, uid, first); !errors.Is(err, ErrWorkLocked) {
		t.Fatalf("future practical access err = %v, want ErrWorkLocked", err)
	}

	mustExec(t, st, `UPDATE cohort_schedule SET due_on = date('now') WHERE cohort_id = ?`, cid)
	if err := st.RequireLessonAvailable(ctx, uid, first); err != nil {
		t.Fatalf("first released practical item: %v", err)
	}
	if err := st.RequireLessonAvailable(ctx, uid, second); !errors.Is(err, ErrWorkLocked) {
		t.Fatalf("second practical access err = %v, want prerequisite lock", err)
	}
	if err := st.MarkLessonComplete(ctx, uid, first); err != nil {
		t.Fatalf("complete prerequisite: %v", err)
	}
	if err := st.RequireLessonAvailable(ctx, uid, second); err != nil {
		t.Fatalf("second practical item stayed locked: %v", err)
	}
}

func TestCheckpointPassRuleIsSharedByGateAndSchedule(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	var assignmentID, courseID int64
	if err := st.db.QueryRow(`SELECT assignment_id, course_id FROM cohort_schedule WHERE cohort_id = ? AND assignment_id IS NOT NULL`, cid).
		Scan(&assignmentID, &courseID); err != nil {
		t.Fatalf("checkpoint: %v", err)
	}
	mustExec(t, st, `UPDATE assignments SET pass_points = 80 WHERE id = ?`, assignmentID)
	mustExec(t, st, `INSERT INTO submissions (assignment_id, user_id, status, score) VALUES (?, ?, 'graded', 79)`, assignmentID, uid)

	items, err := st.CohortSchedule(ctx, cid, uid)
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	for _, item := range items {
		if item.AssignmentID.Valid && item.Done {
			t.Fatal("below-threshold checkpoint counted complete in schedule")
		}
	}
	if gate, err := st.CourseGateFor(ctx, uid, courseID); err != nil || gate.RequiredPassed != 0 {
		t.Fatalf("gate after failed checkpoint = %+v, %v", gate, err)
	}

	mustExec(t, st, `INSERT INTO submissions (assignment_id, user_id, status, score) VALUES (?, ?, 'graded', 80)`, assignmentID, uid)
	items, err = st.CohortSchedule(ctx, cid, uid)
	if err != nil {
		t.Fatalf("schedule after pass: %v", err)
	}
	passed := false
	for _, item := range items {
		if item.AssignmentID.Valid {
			passed = item.Done
		}
	}
	gate, err := st.CourseGateFor(ctx, uid, courseID)
	if err != nil || !passed || gate.RequiredPassed != 1 {
		t.Fatalf("shared pass state: schedule=%t gate=%+v err=%v", passed, gate, err)
	}

	// Retrying after a pass cannot revoke earned completion.
	mustExec(t, st, `INSERT INTO submissions (assignment_id, user_id, status, score) VALUES (?, ?, 'graded', 0)`, assignmentID, uid)
	gate, _ = st.CourseGateFor(ctx, uid, courseID)
	if gate.RequiredPassed != 1 {
		t.Fatal("later failed retry revoked a passed checkpoint")
	}
}

func TestAdvanceCompletionsFinishesPathAndCohortAfterScheduleEnd(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	items, err := st.CohortSchedule(ctx, cid, uid)
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	for _, item := range items {
		switch {
		case item.LessonID.Valid:
			mustExec(t, st, `INSERT INTO lesson_progress (user_id, lesson_id) VALUES (?, ?)`, uid, item.LessonID.Int64)
		case item.AssignmentID.Valid:
			mustExec(t, st, `INSERT INTO submissions (assignment_id, user_id, status, score) VALUES (?, ?, 'graded', 100)`, item.AssignmentID.Int64, uid)
		}
	}

	mustExec(t, st, `UPDATE cohort_schedule SET due_on = date('now', '+1 day') WHERE cohort_id = ?`, cid)
	if result, err := st.AdvanceCompletions(ctx); err != nil || result.Enrollments != 0 || result.Cohorts != 0 {
		t.Fatalf("early completion = %+v, %v", result, err)
	}
	mustExec(t, st, `UPDATE cohort_schedule SET due_on = date('now', '-1 day') WHERE cohort_id = ?`, cid)
	result, err := st.AdvanceCompletions(ctx)
	if err != nil || result.Enrollments != 1 || result.Cohorts != 1 {
		t.Fatalf("completion = %+v, %v", result, err)
	}
	if live, err := st.LiveEnrollment(ctx, uid); err != nil || live != nil {
		t.Fatalf("completed enrollment still live: %+v, %v", live, err)
	}
	var enrollmentState, cohortState string
	var completedAt sql.NullString
	if err := st.db.QueryRow(`SELECT state FROM enrollments WHERE user_id = ? ORDER BY id DESC LIMIT 1`, uid).Scan(&enrollmentState); err != nil {
		t.Fatalf("enrollment state: %v", err)
	}
	if err := st.db.QueryRow(`SELECT state, completed_at FROM cohorts WHERE id = ?`, cid).Scan(&cohortState, &completedAt); err != nil {
		t.Fatalf("cohort state: %v", err)
	}
	if enrollmentState != EnrollCompleted || cohortState != CohortCompleted || !completedAt.Valid {
		t.Fatalf("terminal states = enrollment:%s cohort:%s completed:%+v", enrollmentState, cohortState, completedAt)
	}
	if again, err := st.AdvanceCompletions(ctx); err != nil || again.Enrollments != 0 || again.Cohorts != 0 {
		t.Fatalf("completion was not idempotent: %+v, %v", again, err)
	}
	if _, err := st.EnsureEnrollment(ctx, uid); !errors.Is(err, ErrReapplicationCooldown) {
		t.Fatalf("completed learner reapplication err = %v, want cooldown", err)
	}
}

func TestDueOnReturnsOnlyThatDay(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid, _ := scheduledCohort(t, st)
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	// Day 0 is the Monday the cohort starts: two lessons at the daily budget.
	due, err := st.DueOn(ctx, cid, uid, "2026-09-07")
	if err != nil {
		t.Fatalf("due: %v", err)
	}
	if len(due) != 2 {
		t.Fatalf("day 0 has %d items, want 2 (the daily budget)", len(due))
	}
	for _, it := range due {
		if it.DueOn != "2026-09-07" {
			t.Fatalf("DueOn returned an item for %s", it.DueOn)
		}
	}
	// A weekend is empty.
	if sat, _ := st.DueOn(ctx, cid, uid, "2026-09-12"); len(sat) != 0 {
		t.Fatalf("Saturday has %d items, want 0", len(sat))
	}
}
