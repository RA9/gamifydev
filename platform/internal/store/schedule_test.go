package store

import (
	"context"
	"testing"
	"time"
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
