package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/attendance"
)

// attendanceFixture builds a cohort with a materialized schedule that started in
// the past, so there are scheduled days to resolve.
func attendanceFixture(t *testing.T, st *Store) (cohortID, userID int64) {
	t.Helper()
	ctx := context.Background()
	pathID := mkPath(t, st, "frontend")
	mkCourse(t, st, pathID, "c1", 0, 10, false)
	// A Monday well in the past, so every scheduled day has already elapsed.
	cid, err := st.CreateCohort(ctx, Cohort{
		PathID: pathID, Name: "C", TZBand: "europe_africa", StartsOn: "2020-01-06",
	})
	if err != nil {
		t.Fatalf("cohort: %v", err)
	}
	uid := enroll(t, st, "attend@x.com", pathID, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, cid, pathID, "europe_africa", 1); err != nil {
		t.Fatalf("fill: %v", err)
	}
	if _, err := st.MaterializeSchedule(ctx, cid); err != nil {
		t.Fatalf("materialize: %v", err)
	}
	return cid, uid
}

func TestResolveAttendanceOnlyCoversScheduledDays(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)

	n, err := st.ResolveAttendance(ctx, cid, "2020-01-01")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	// 10 lessons at 2/day = 5 scheduled days.
	if n != 5 {
		t.Fatalf("resolved %d learner-days, want 5", n)
	}
	days, _ := st.AttendanceDays(ctx, uid)
	if len(days) != 5 {
		t.Fatalf("got %d marks, want 5", len(days))
	}
	// Weekends must never appear — counting them would punish the rest the
	// schedule deliberately builds in.
	for _, d := range days {
		var isScheduled int
		_ = st.db.QueryRow(
			`SELECT COUNT(*) FROM cohort_schedule WHERE cohort_id = ? AND due_on = ?`,
			cid, d.Day).Scan(&isScheduled)
		if isScheduled == 0 {
			t.Fatalf("resolved %s, which is not a scheduled day", d.Day)
		}
	}
	// Silence everywhere: all absent.
	for _, d := range days {
		if d.State != attendance.Absent {
			t.Fatalf("day %s = %q with no activity, want absent", d.Day, d.State)
		}
	}
}

func TestResolveIsIdempotentAndReflectsLaterSignals(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)

	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	days, _ := st.AttendanceDays(ctx, uid)
	first := days[0].Day

	// A signal arrives late (the rollup ran after midnight). Re-resolving must
	// upgrade the day rather than leaving a stale absence.
	if err := st.MarkActive(ctx, uid, first, "standup"); err != nil {
		t.Fatalf("mark active: %v", err)
	}
	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("re-resolve: %v", err)
	}
	days, _ = st.AttendanceDays(ctx, uid)
	if len(days) != 5 {
		t.Fatalf("re-resolving duplicated marks: %d rows", len(days))
	}
	if days[0].State != attendance.Present {
		t.Fatalf("day %s = %q after activity arrived, want present", days[0].Day, days[0].State)
	}
}

func TestAbsenceNoticeExcusesTheDay(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	days, _ := st.AttendanceDays(ctx, uid)
	from, to := days[0].Day, days[2].Day

	// This is the mechanism the whole policy turns on: saying "I'll be out"
	// converts silence into an excused absence.
	if err := st.FileAbsence(ctx, uid, cid, from, to, "exams"); err != nil {
		t.Fatalf("file: %v", err)
	}
	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("re-resolve: %v", err)
	}
	days, _ = st.AttendanceDays(ctx, uid)
	for i := 0; i < 3; i++ {
		if days[i].State != attendance.Excused {
			t.Fatalf("day %s = %q, want excused (covered by notice)", days[i].Day, days[i].State)
		}
	}
	if days[3].State != attendance.Absent {
		t.Fatalf("day outside the notice = %q, want absent", days[3].State)
	}

	// And the run is broken, so no sanction.
	if d := attendance.Evaluate(days); d.Sanction {
		t.Fatal("sanctioned despite an absence notice covering the run")
	}
}

func TestShadowSanctionChangesNothing(t *testing.T) {
	// The single most important guarantee in this phase: a shadow sanction is a
	// record, not an action.
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("resolve: %v", err)
	}

	id, isNew, err := st.RecordSanction(ctx, Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, Reason: "test", WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}, true /* shadow */)
	if err != nil || !isNew {
		t.Fatalf("record = %v, %v", isNew, err)
	}

	// Still in the cohort, still enrolled, still active.
	if c, _ := st.CohortForUser(ctx, uid); c == nil {
		t.Fatal("a shadow sanction removed the learner from their cohort")
	}
	e, _ := st.LiveEnrollment(ctx, uid)
	if e == nil || e.State != EnrollActive {
		t.Fatalf("enrollment = %+v after a shadow sanction, want active", e)
	}
	var applied sql.NullString
	_ = st.db.QueryRow(`SELECT applied_at FROM sanctions WHERE id = ?`, id).Scan(&applied)
	if applied.Valid {
		t.Fatal("a shadow sanction recorded an applied_at timestamp")
	}

	// ApplySanction must refuse to act on a shadow record.
	if err := st.ApplySanction(ctx, id); err == nil {
		t.Fatal("ApplySanction acted on a shadow sanction")
	}
	if c, _ := st.CohortForUser(ctx, uid); c == nil {
		t.Fatal("learner removed after ApplySanction on a shadow record")
	}
}

func TestSanctionIsRecordedOncePerWindow(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	sc := Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}
	if _, isNew, _ := st.RecordSanction(ctx, sc, true); !isNew {
		t.Fatal("first record was not new")
	}
	// The evaluator runs every six hours; the same run of silence must not pile
	// up duplicate records.
	if _, isNew, _ := st.RecordSanction(ctx, sc, true); isNew {
		t.Fatal("the same absence window produced a second sanction")
	}
	var n int
	_ = st.db.QueryRow(`SELECT COUNT(*) FROM sanctions WHERE user_id = ?`, uid).Scan(&n)
	if n != 1 {
		t.Fatalf("%d sanctions recorded, want 1", n)
	}
}

func TestRealSanctionDropsButDoesNotBan(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)

	id, _, err := st.RecordSanction(ctx, Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, Reason: "silence", WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}, false /* apply for real */)
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := st.ApplySanction(ctx, id); err != nil {
		t.Fatalf("apply: %v", err)
	}

	// Out of the cohort and enrollment ended...
	if c, _ := st.CohortForUser(ctx, uid); c != nil {
		t.Fatal("learner still in the cohort after an applied drop")
	}
	if e, _ := st.LiveEnrollment(ctx, uid); e != nil {
		t.Fatalf("live enrollment survived a drop: %+v", e)
	}
	// ...but the account is untouched. Attendance drops; it never bans.
	state, err := st.AccountState(ctx, uid)
	if err != nil {
		t.Fatalf("account state: %v", err)
	}
	if state != AccountActive {
		t.Fatalf("account state = %q after an attendance drop, want active", state)
	}
	// And re-application is possible — that's what "drop, not ban" means.
	if _, err := st.EnsureEnrollment(ctx, uid); err != nil {
		t.Fatalf("dropped learner cannot re-enroll: %v", err)
	}
}

func TestAppealGrantedRestoresTheLearner(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	admin := mkUser(t, st, "admin@x.com")

	id, _, _ := st.RecordSanction(ctx, Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}, false)
	if err := st.ApplySanction(ctx, id); err != nil {
		t.Fatalf("apply: %v", err)
	}

	if err := st.FileAppeal(ctx, id, uid, "I was in hospital"); err != nil {
		t.Fatalf("appeal: %v", err)
	}
	open, _ := st.OpenAppeals(ctx)
	if len(open) != 1 {
		t.Fatalf("%d open appeals, want 1", len(open))
	}

	if err := st.DecideAppeal(ctx, open[0].ID, admin, true, "confirmed"); err != nil {
		t.Fatalf("decide: %v", err)
	}
	// Back in the cohort and active again.
	if c, _ := st.CohortForUser(ctx, uid); c == nil || c.ID != cid {
		t.Fatal("granting an appeal did not restore cohort membership")
	}
	e, _ := st.LiveEnrollment(ctx, uid)
	if e == nil || e.State != EnrollActive {
		t.Fatalf("enrollment = %+v after a granted appeal, want active", e)
	}
	if left, _ := st.OpenAppeals(ctx); len(left) != 0 {
		t.Fatal("appeal still open after a decision")
	}
}

func TestAppealUpheldLeavesSanctionInPlace(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	admin := mkUser(t, st, "admin2@x.com")

	id, _, _ := st.RecordSanction(ctx, Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}, false)
	_ = st.ApplySanction(ctx, id)
	_ = st.FileAppeal(ctx, id, uid, "please")
	open, _ := st.OpenAppeals(ctx)

	if err := st.DecideAppeal(ctx, open[0].ID, admin, false, "record stands"); err != nil {
		t.Fatalf("decide: %v", err)
	}
	if c, _ := st.CohortForUser(ctx, uid); c != nil {
		t.Fatal("upholding an appeal restored the learner anyway")
	}
	sc, _ := st.LatestSanction(ctx, uid)
	if sc == nil || sc.AppealState != "upheld" {
		t.Fatalf("appeal state = %+v, want upheld", sc)
	}
}

func TestCannotAppealSomeoneElsesSanction(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)
	stranger := mkUser(t, st, "stranger@x.com")

	id, _, _ := st.RecordSanction(ctx, Sanction{
		UserID: uid, CohortID: sql.NullInt64{Int64: cid, Valid: true},
		Kind: attendance.KindDrop, WindowFrom: "2020-01-06", WindowTo: "2020-01-08",
	}, true)
	if err := st.FileAppeal(ctx, id, stranger, "not mine"); err == nil {
		t.Fatal("a stranger filed an appeal against someone else's sanction")
	}
}

func TestLessonCompletionCountsAsAttendance(t *testing.T) {
	// Doing the work is showing up, even on a day you skip standup.
	ctx := context.Background()
	st := newTestStore(t)
	cid, uid := attendanceFixture(t, st)

	var lessonID int64
	var day string
	if err := st.db.QueryRow(
		`SELECT lesson_id, due_on FROM cohort_schedule WHERE cohort_id = ? AND lesson_id IS NOT NULL
		 ORDER BY day_index LIMIT 1`, cid).Scan(&lessonID, &day); err != nil {
		t.Fatalf("pick lesson: %v", err)
	}
	if err := st.MarkLessonComplete(ctx, uid, lessonID); err != nil {
		t.Fatalf("complete: %v", err)
	}
	// Backdate it onto the scheduled day, then roll up.
	mustExec(t, st, `UPDATE lesson_progress SET completed_at = ? || ' 12:00:00' WHERE lesson_id = ?`, day, lessonID)
	if _, err := st.RollupActivity(ctx, "2019-01-01"); err != nil {
		t.Fatalf("rollup: %v", err)
	}
	if _, err := st.ResolveAttendance(ctx, cid, "2020-01-01"); err != nil {
		t.Fatalf("resolve: %v", err)
	}

	days, _ := st.AttendanceDays(ctx, uid)
	var found bool
	for _, d := range days {
		if d.Day == day {
			found = true
			if d.State != attendance.Present {
				t.Fatalf("day %s = %q after completing a lesson, want present", d.Day, d.State)
			}
		}
	}
	if !found {
		t.Fatalf("scheduled day %s was not resolved", day)
	}
}
