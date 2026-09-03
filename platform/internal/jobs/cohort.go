package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/cohort"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// Cohort lifecycle jobs: forming groups, repacking ones that decayed, and
// running the daily standup window.

// registerCohortJobs adds phase 3's scheduled work.
func registerCohortJobs(r *Runner, st *store.Store) {
	// cohort:form groups queued learners by path and timezone band.
	//
	// A group qualifies either by reaching the target size or by having someone
	// who has waited long enough that starting matters more than group size —
	// which is also what makes a cohort of one work on day zero.
	r.Register(Job{
		Name:    "cohort:form",
		Every:   15 * time.Minute,
		Timeout: time.Minute,
		Run: func(ctx context.Context) (string, error) {
			groups, err := st.PlacementQueue(ctx)
			if err != nil {
				return "", err
			}
			botID, err := st.EnsureBotMentor(ctx)
			if err != nil {
				return "", err
			}
			var formed []string
			for _, g := range groups {
				if !cohort.ShouldForm(g.Waiting, g.OldestAge) {
					continue
				}
				band := cohort.LookupBand(g.TZBand)
				start := time.Now().UTC()
				id, err := st.CreateCohort(ctx, store.Cohort{
					PathID:   g.PathID,
					Name:     cohort.Name(g.PathTitle, band, start),
					TZBand:   band.Key,
					StartsOn: start.Format("2006-01-02"),
				})
				if err != nil {
					return "", err
				}
				// Every cohort gets a mentor, including a cohort of one.
				if err := st.AddMember(ctx, id, botID, store.RoleMentor); err != nil {
					return "", err
				}
				users, err := st.TakeFromQueue(ctx, id, g.PathID, g.TZBand, cohort.TakeSize(g.Waiting))
				if err != nil {
					return "", err
				}
				formed = append(formed, fmt.Sprintf("%s/%s×%d", g.PathTitle, band.Key, len(users)))
			}
			if len(formed) == 0 {
				return "", nil
			}
			return "formed " + strings.Join(formed, ", "), nil
		},
	})

	// cohort:repack merges cohorts that decayed below the viable size.
	//
	// This is the counterweight to an aggressive drop policy: without it, a
	// cohort bleeds down to one person and the social pressure that makes the
	// whole model work disappears.
	r.Register(Job{
		Name:    "cohort:repack",
		Every:   6 * time.Hour,
		Timeout: time.Minute,
		Run: func(ctx context.Context) (string, error) {
			thin, err := st.ThinCohorts(ctx, cohort.MinViable)
			if err != nil {
				return "", err
			}
			merged := 0
			for _, t := range thin {
				// An empty cohort has nobody to move; leave it for the operator
				// rather than silently archiving a group that may still fill.
				if t.Members == 0 {
					continue
				}
				target, err := st.MergeTarget(ctx, t.ID, t.PathID, t.TZBand, t.Members, cohort.MaxSize)
				if err != nil {
					return "", err
				}
				if target == 0 {
					continue // nowhere to merge into yet; try again next run
				}
				n, err := st.MergeCohorts(ctx, t.ID, target)
				if err != nil {
					return "", err
				}
				merged += n
			}
			if merged == 0 {
				return "", nil
			}
			return fmt.Sprintf("repacked %d learner(s)", merged), nil
		},
	})

	// standup:open creates each active cohort's standup for its local day, with
	// the day's rotating prompt.
	r.Register(Job{
		Name:    "standup:open",
		Every:   30 * time.Minute,
		Timeout: time.Minute,
		Run: func(ctx context.Context) (string, error) {
			cohorts, err := st.OpenCohorts(ctx)
			if err != nil {
				return "", err
			}
			opened := 0
			now := time.Now().UTC()
			for _, c := range cohorts {
				band := cohort.LookupBand(c.TZBand)
				day := cohort.LocalDay(now, band)
				opens, closes := cohort.Window(day, band)
				_, created, err := st.EnsureStandup(ctx, c.ID,
					day.Format("2006-01-02"),
					opens.Format("2006-01-02 15:04:05"),
					closes.Format("2006-01-02 15:04:05"),
					cohort.Prompt(day))
				if err != nil {
					return "", err
				}
				if created {
					opened++
				}
			}
			if opened == 0 {
				return "", nil
			}
			return fmt.Sprintf("opened %d standup(s)", opened), nil
		},
	})

	// schedule:materialize expands a cohort's path into dated work.
	//
	// Runs shortly after formation rather than inside it: a cohort exists the
	// moment it is formed, and a plan that fails to build must not roll back the
	// grouping. Cheap to re-run — a scheduled cohort short-circuits.
	r.Register(Job{
		Name:    "schedule:materialize",
		Every:   10 * time.Minute,
		Timeout: 2 * time.Minute,
		Run: func(ctx context.Context) (string, error) {
			ids, err := st.UnscheduledCohorts(ctx)
			if err != nil {
				return "", err
			}
			total, built := 0, 0
			for _, id := range ids {
				n, err := st.MaterializeSchedule(ctx, id)
				if err != nil {
					return "", err
				}
				if n > 0 {
					built++
					total += n
				}
			}
			if built == 0 {
				return "", nil
			}
			return fmt.Sprintf("scheduled %d cohort(s), %d item(s)", built, total), nil
		},
	})

	// standup:close shuts windows whose time has passed.
	//
	// Phase 5 will hang attendance resolution off this; for now closing the
	// window is all it does, and a post after close is recorded as late rather
	// than rejected outright.
	r.Register(Job{
		Name:    "standup:close",
		Every:   30 * time.Minute,
		Timeout: time.Minute,
		Run: func(ctx context.Context) (string, error) {
			n, err := st.CloseExpiredStandups(ctx)
			if err != nil {
				return "", err
			}
			if n == 0 {
				return "", nil
			}
			return fmt.Sprintf("closed %d standup(s)", n), nil
		},
	})
}
