package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/cohort"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// The cohort space: who you're learning with, and today's standup.

// handleCohort renders the learner's cohort — members, today's standup, and the
// recent history.
//
// Three states it must handle honestly:
//   - not enrolled yet          → point at the diagnostic
//   - enrolled, awaiting a group → say so, with what we're waiting for
//   - in a cohort                → the real page
func (s *Server) handleCohort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	c, err := s.st.CohortForUser(ctx, u.ID)
	if err != nil {
		http.Error(w, "could not load your cohort", http.StatusInternalServerError)
		return
	}
	if c == nil {
		enr, _ := s.st.LiveEnrollment(ctx, u.ID)
		if enr == nil {
			enr, _ = s.st.LatestEnrollment(ctx, u.ID)
		}
		placed, _ := s.st.LatestPlacement(ctx, u.ID)
		s.render(w, r, "cohort_waiting.html", ViewData{Title: "Your cohort", Data: map[string]any{
			"bodyClass":  "cohort-dark",
			"enrollment": enr,
			"placed":     placed != nil,
			"completed":  enr != nil && enr.State == store.EnrollCompleted,
			"cooldown":   enr != nil && store.EnrollmentTerminal(enr.State),
			"target":     cohort.TargetSize,
		}})
		return
	}

	band := cohort.LookupBand(c.TZBand)
	day := cohort.LocalDay(time.Now().UTC(), band)
	dayStr := day.Format("2006-01-02")

	members, _ := s.st.CohortMembers(ctx, c.ID)
	today, _ := s.st.TodayStandup(ctx, c.ID, dayStr)

	var entries []store.StandupEntry
	var posted bool
	if today != nil {
		entries, _ = s.st.StandupEntries(ctx, today.ID)
		posted, _ = s.st.HasPosted(ctx, today.ID, u.ID)
		// Flag who has already posted so the roster shows the day at a glance —
		// this is the social pressure the model actually runs on.
		byUser := map[int64]bool{}
		for _, e := range entries {
			byUser[e.UserID] = true
		}
		for i := range members {
			members[i].PostedToday = byUser[members[i].UserID]
		}
	}

	recent, _ := s.st.RecentStandups(ctx, c.ID, 7)
	// The current day is rendered in its own section; don't repeat it below.
	var history []store.Standup
	for _, st := range recent {
		if st.Day != dayStr {
			history = append(history, st)
		}
	}

	// Blockers are what a mentor acts on, so surface them separately.
	var blockers []store.StandupEntry
	for _, e := range entries {
		if strings.TrimSpace(e.Blockers) != "" {
			blockers = append(blockers, e)
		}
	}

	s.render(w, r, "cohort.html", ViewData{Title: c.Name, Data: map[string]any{
		"bodyClass": "cohort-dark",
		"cohort":    c,
		"band":      band,
		"members":   members,
		"standup":   today,
		"entries":   entries,
		"blockers":  blockers,
		"posted":    posted,
		"history":   history,
		"day":       dayStr,
	}})
}

// handleStandupPost records a member's standup entry.
func (s *Server) handleStandupPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	c, err := s.st.CohortForUser(ctx, u.ID)
	if err != nil || c == nil {
		http.Redirect(w, r, "/cohort", http.StatusSeeOther)
		return
	}
	band := cohort.LookupBand(c.TZBand)
	day := cohort.LocalDay(time.Now().UTC(), band).Format("2006-01-02")
	today, err := s.st.TodayStandup(ctx, c.ID, day)
	if err != nil || today == nil {
		http.Redirect(w, r, "/cohort", http.StatusSeeOther)
		return
	}

	yesterday := strings.TrimSpace(r.FormValue("yesterday"))
	todayText := strings.TrimSpace(r.FormValue("today"))
	blockers := strings.TrimSpace(r.FormValue("blockers"))
	if todayText == "" {
		// "What are you doing today" is the one field that carries the standup;
		// an empty post is not a check-in.
		http.Redirect(w, r, "/cohort?empty=1", http.StatusSeeOther)
		return
	}

	err = s.st.PostStandup(ctx, today.ID, u.ID, yesterday, todayText, blockers)
	if errors.Is(err, store.ErrStandupClosed) {
		http.Redirect(w, r, "/cohort?closed=1", http.StatusSeeOther)
		return
	}
	if err != nil {
		http.Error(w, "could not post your standup", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cohort", http.StatusSeeOther)
}

// --- Admin -------------------------------------------------------------------

// handleAdminCohorts shows cohort health: how many are running, how many
// learners are queued, and which groups are thin enough to need repacking.
func (s *Server) handleAdminCohorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, _ := s.st.CohortOverview(ctx)
	list, _ := s.st.ListCohorts(ctx, 50)
	queue, _ := s.st.PlacementQueue(ctx)

	type queueRow struct {
		Path, Band string
		Waiting    int
		WaitedFor  string
		Ready      bool
	}
	var rows []queueRow
	for _, g := range queue {
		rows = append(rows, queueRow{
			Path: g.PathTitle, Band: cohort.LookupBand(g.TZBand).Label,
			Waiting:   g.Waiting,
			WaitedFor: humanEvery(g.OldestAge.Truncate(time.Minute)),
			Ready:     cohort.ShouldForm(g.Waiting, g.OldestAge),
		})
	}

	s.render(w, r, "admin_cohorts.html", ViewData{Title: "Cohorts", Data: map[string]any{
		"stats":     stats,
		"cohorts":   list,
		"queue":     rows,
		"target":    cohort.TargetSize,
		"minViable": cohort.MinViable,
	}})
}

// handleSchedule renders the cohort's full plan for this learner: what's due
// today, what slipped, and the sprint-by-sprint map of the path ahead.
func (s *Server) handleSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	c, err := s.st.CohortForUser(ctx, u.ID)
	if err != nil || c == nil {
		http.Redirect(w, r, "/cohort", http.StatusSeeOther)
		return
	}
	today := time.Now().UTC().Format("2006-01-02")
	due, _ := s.st.DueOn(ctx, c.ID, u.ID, today)
	overdue, _ := s.st.Overdue(ctx, c.ID, u.ID, 20)
	all, _ := s.st.CohortSchedule(ctx, c.ID, u.ID)
	prog, _ := s.st.Progress(ctx, c.ID, u.ID)

	// Group the plan by sprint so the page reads as a map rather than a list.
	type sprintGroup struct {
		Sprint int
		Items  []store.ScheduleItem
	}
	var sprints []sprintGroup
	for _, it := range all {
		if len(sprints) == 0 || sprints[len(sprints)-1].Sprint != it.Sprint {
			sprints = append(sprints, sprintGroup{Sprint: it.Sprint})
		}
		g := &sprints[len(sprints)-1]
		g.Items = append(g.Items, it)
	}

	s.render(w, r, "schedule.html", ViewData{Title: "Your schedule", Data: map[string]any{
		"bodyClass": "cohort-dark",
		"cohort":    c,
		"today":     due,
		"overdue":   overdue,
		"sprints":   sprints,
		"progress":  prog,
		"todayDate": today,
	}})
}
