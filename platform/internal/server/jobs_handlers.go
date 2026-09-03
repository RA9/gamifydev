package server

import (
	"context"
	"net/http"
	"time"
)

// jobView is one row of the admin jobs table: the job's schedule joined with
// what its history says actually happened.
type jobView struct {
	Name     string
	Every    string
	LastRun  string
	LastOK   bool
	HasRun   bool
	Runs24h  int
	Fails24h int
}

// handleAdminJobs renders the scheduled-work overview: what is registered, when
// each job last ran, and whether it succeeded. Without this the schedule is
// invisible until something has already gone wrong.
func (s *Server) handleAdminJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	health, _ := s.st.JobHealthSummary(ctx)
	byName := map[string]int{}
	for i, h := range health {
		byName[h.Name] = i
	}

	var views []jobView
	if s.jobs != nil {
		for _, j := range s.jobs.Jobs() {
			v := jobView{Name: j.Name, Every: humanEvery(j.Every)}
			if i, ok := byName[j.Name]; ok {
				h := health[i]
				v.HasRun = h.LastRun.Valid
				v.LastRun = h.LastRun.String
				v.LastOK = h.LastOK.Valid && h.LastOK.Bool
				v.Runs24h, v.Fails24h = h.Runs24h, h.Fails24h
			}
			views = append(views, v)
		}
	}

	runs, _ := s.st.RecentJobRuns(ctx, 40)
	s.render(w, r, "admin_jobs.html", ViewData{
		Title: "Scheduled work",
		Data: map[string]any{
			"jobs":     views,
			"runs":     runs,
			"disabled": s.jobs == nil,
		},
	})
}

// handleAdminJobRun triggers a job immediately. It still takes the job's lease,
// so a manual run cannot collide with a scheduled one.
func (s *Server) handleAdminJobRun(w http.ResponseWriter, r *http.Request) {
	if s.jobs == nil {
		http.Error(w, "scheduled work is disabled", http.StatusServiceUnavailable)
		return
	}
	name := r.PathValue("name")
	// Run detached from the request: the job outlives the response, and a
	// client disconnect must not cancel it mid-write.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		_ = s.jobs.RunNow(ctx, name)
	}()
	http.Redirect(w, r, "/admin/jobs", http.StatusSeeOther)
}

// humanEvery renders a cadence the way an operator would say it.
func humanEvery(d time.Duration) string {
	switch {
	case d >= 24*time.Hour:
		if n := int(d / (24 * time.Hour)); n == 1 {
			return "daily"
		} else {
			return "every " + itoa(n) + " days"
		}
	case d >= time.Hour:
		if n := int(d / time.Hour); n == 1 {
			return "hourly"
		} else {
			return "every " + itoa(n) + "h"
		}
	case d >= time.Minute:
		return "every " + itoa(int(d/time.Minute)) + "m"
	default:
		return d.String()
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
