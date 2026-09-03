package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RA9/gamifydev/platform/internal/attendance"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// Attendance jobs: resolving each scheduled day, and evaluating the rolling
// window against it.
//
// `enforce` is the switch the PRD insists on (§07). When false — the default —
// the evaluator records what it *would* have done and changes nothing. Run it
// that way for at least one full cohort, look at the shadow sanctions, then set
// the thresholds. Enforcing an uncalibrated rule removes your first cohort and
// destroys the evidence that the number was wrong.
func registerAttendanceJobs(r *Runner, st *store.Store, enforce bool) {
	// attendance:resolve records present/excused/absent for each scheduled day
	// that has passed. Re-scans a trailing window so a late-arriving signal
	// (a standup posted at 23:58, a rollup that ran after midnight) still counts.
	r.Register(Job{
		Name:    "attendance:resolve",
		Every:   time.Hour,
		Timeout: 2 * time.Minute,
		Run: func(ctx context.Context) (string, error) {
			cohorts, err := st.OpenCohorts(ctx)
			if err != nil {
				return "", err
			}
			since := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
			marks := 0
			for _, c := range cohorts {
				n, err := st.ResolveAttendance(ctx, c.ID, since)
				if err != nil {
					return "", err
				}
				marks += n
			}
			if marks == 0 {
				return "", nil
			}
			return fmt.Sprintf("resolved %d learner-day(s)", marks), nil
		},
	})

	// attendance:evaluate applies the rolling window and records sanctions.
	r.Register(Job{
		Name:    "attendance:evaluate",
		Every:   6 * time.Hour,
		Timeout: 2 * time.Minute,
		Run: func(ctx context.Context) (string, error) {
			learners, err := st.LearnersToEvaluate(ctx)
			if err != nil {
				return "", err
			}
			recorded, applied := 0, 0
			for _, l := range learners {
				days, err := st.AttendanceDays(ctx, l.UserID)
				if err != nil {
					return "", err
				}
				d := attendance.Evaluate(days)
				if !d.Sanction {
					continue
				}
				id, isNew, err := st.RecordSanction(ctx, store.Sanction{
					UserID:     l.UserID,
					CohortID:   sql.NullInt64{Int64: l.CohortID, Valid: true},
					Kind:       d.Kind,
					Reason:     d.Reason,
					WindowFrom: d.Run.From,
					WindowTo:   d.Run.To,
				}, !enforce)
				if err != nil {
					return "", err
				}
				if !isNew {
					continue // already recorded for this run of absences
				}
				recorded++
				if enforce {
					if err := st.ApplySanction(ctx, id); err != nil {
						return "", err
					}
					applied++
				}
			}
			if recorded == 0 {
				return "", nil
			}
			if enforce {
				return fmt.Sprintf("applied %d sanction(s)", applied), nil
			}
			// Say plainly that nothing happened to anyone — an operator reading
			// the job log should never have to infer that.
			return fmt.Sprintf("SHADOW: would have sanctioned %d learner(s); none applied", recorded), nil
		},
	})

	if enforce {
		log.Printf("attendance: ENFORCEMENT ENABLED — sanctions will be applied")
	} else {
		log.Printf("attendance: shadow mode (set ENFORCE_ATTENDANCE=1 to apply sanctions)")
	}
}
