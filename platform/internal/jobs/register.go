package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// Register wires the platform's scheduled jobs onto a runner.
//
// Phase 1 registers only the jobs whose backing tables exist today. The cohort,
// standup, attendance and sanction jobs named in the PRD are added by their own
// phases rather than being stubbed here — a registered job that does nothing is
// worse than an absent one, because it reports healthy.
func Register(r *Runner, st *store.Store) {
	// activity:rollup materializes activity_days from step completions and
	// submissions. It re-scans a trailing window rather than only "since last
	// run" so a missed run, a clock skew, or a late-arriving row still lands;
	// the insert is a no-op for days already recorded.
	r.Register(Job{
		Name:    "activity:rollup",
		Every:   15 * time.Minute,
		Timeout: time.Minute,
		Run: func(ctx context.Context) (string, error) {
			since := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
			n, err := st.RollupActivity(ctx, since)
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", nil // quiet: nothing new to say
			}
			return fmt.Sprintf("recorded %d active day(s)", n), nil
		},
	})

	// sanction:expire restores accounts whose suspension has elapsed. The read
	// path in AccountState already treats an expired suspension as active, so
	// this is bookkeeping rather than enforcement — it keeps the table honest
	// and makes the admin view show the real state.
	r.Register(Job{
		Name:    "sanction:expire",
		Every:   time.Hour,
		Timeout: 30 * time.Second,
		Run: func(ctx context.Context) (string, error) {
			n, err := st.ExpireSanctions(ctx)
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", nil
			}
			return fmt.Sprintf("restored %d account(s)", n), nil
		},
	})

	// placement:expire closes out sittings that ran past their deadline without
	// a submission, so "you have a diagnostic in progress" cannot stick forever.
	// The submit path already refuses an expired attempt; this is what lets the
	// learner start a fresh one.
	r.Register(Job{
		Name:    "placement:expire",
		Every:   10 * time.Minute,
		Timeout: 30 * time.Second,
		Run: func(ctx context.Context) (string, error) {
			n, err := st.ExpireAbandonedAttempts(ctx)
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", nil
			}
			return fmt.Sprintf("closed %d expired sitting(s)", n), nil
		},
	})

	// jobs:prune keeps run history bounded.
	r.Register(Job{
		Name:    "jobs:prune",
		Every:   24 * time.Hour,
		Timeout: 30 * time.Second,
		Run: func(ctx context.Context) (string, error) {
			n, err := st.PruneJobRuns(ctx, 30)
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", nil
			}
			return fmt.Sprintf("pruned %d old run(s)", n), nil
		},
	})
}
