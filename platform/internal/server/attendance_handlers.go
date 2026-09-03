package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/attendance"
	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// Attendance: a learner's own record, the "I'll be out" flow, and appeals.

// handleAttendance shows a learner their attendance, their remaining excused
// days, and any sanction against them.
//
// Framed as a record, not a warning. The page exists so nothing about
// enforcement is a surprise — a learner should always be able to see exactly
// where they stand before anything happens.
func (s *Server) handleAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	days, err := s.st.AttendanceDays(ctx, u.ID)
	if err != nil {
		http.Error(w, "could not load your attendance", http.StatusInternalServerError)
		return
	}
	notices, _ := s.st.AbsenceNotices(ctx, u.ID, 10)
	sanction, _ := s.st.LatestSanction(ctx, u.ID)
	c, _ := s.st.CohortForUser(ctx, u.ID)

	// Most recent first reads better; the policy works on chronological order,
	// so reverse a copy rather than the slice the policy saw.
	shown := make([]attendance.Day, len(days))
	for i, d := range days {
		shown[len(days)-1-i] = d
	}
	if len(shown) > 30 {
		shown = shown[:30]
	}

	s.render(w, r, "attendance.html", ViewData{Title: "Your attendance", Data: map[string]any{
		"bodyClass":  "cohort-dark",
		"days":       shown,
		"rate":       attendance.Rate(days),
		"budgetLeft": attendance.BudgetLeft(days),
		"budget":     attendance.ExcusedBudget,
		"currentRun": attendance.CurrentRun(days),
		"limit":      attendance.UnexcusedLimit,
		"notices":    notices,
		"sanction":   sanction,
		"cohort":     c,
		"today":      time.Now().UTC().Format("2006-01-02"),
	}})
}

// handleAbsenceFile records an "I'll be out" notice.
//
// This is the whole mechanism the policy turns on: filing converts silence into
// an excused absence. It is deliberately trivial to do and asks for no proof.
func (s *Server) handleAbsenceFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	from := strings.TrimSpace(r.FormValue("from"))
	to := strings.TrimSpace(r.FormValue("to"))
	reason := strings.TrimSpace(r.FormValue("reason"))
	if from == "" {
		http.Redirect(w, r, "/attendance", http.StatusSeeOther)
		return
	}
	if to == "" {
		to = from
	}
	// A backwards range is a slip, not an attack — normalise rather than reject.
	if to < from {
		from, to = to, from
	}
	var cohortID int64
	if c, _ := s.st.CohortForUser(ctx, u.ID); c != nil {
		cohortID = c.ID
	}
	if err := s.st.FileAbsence(ctx, u.ID, cohortID, from, to, reason); err != nil {
		http.Error(w, "could not file your notice", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/attendance", http.StatusSeeOther)
}

// handleAppealFile opens an appeal against a sanction.
func (s *Server) handleAppealFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	id, _ := strconv.ParseInt(r.FormValue("sanction"), 10, 64)
	body := strings.TrimSpace(r.FormValue("body"))
	if id == 0 || body == "" {
		http.Redirect(w, r, "/attendance", http.StatusSeeOther)
		return
	}
	// FileAppeal verifies the sanction belongs to this learner.
	if err := s.st.FileAppeal(ctx, id, u.ID, body); err != nil {
		s.notFound(w, r)
		return
	}
	http.Redirect(w, r, "/attendance", http.StatusSeeOther)
}

// --- Admin -------------------------------------------------------------------

// handleAdminAttendance is the calibration view: what the policy has resolved,
// and — crucially — what it *would* have done if enforcement were on.
//
// This page is the reason phase 5 ships inert. It is where the thresholds get
// chosen, from evidence rather than intuition.
func (s *Server) handleAdminAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, _ := s.st.AttendanceStats(ctx)
	shadow, _ := s.st.ListSanctions(ctx, true, 50)
	real, _ := s.st.ListSanctions(ctx, false, 50)
	appeals, _ := s.st.OpenAppeals(ctx)

	s.render(w, r, "admin_attendance.html", ViewData{Title: "Attendance", Data: map[string]any{
		"stats":    stats,
		"shadow":   shadow,
		"real":     real,
		"appeals":  appeals,
		"enforced": s.enforceAttendance,
		"limit":    attendance.UnexcusedLimit,
		"budget":   attendance.ExcusedBudget,
	}})
}

// handleAdminAppealDecide records a human decision on an appeal.
func (s *Server) handleAdminAppealDecide(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	grant := r.FormValue("decision") == "grant"
	note := strings.TrimSpace(r.FormValue("note"))

	if err := s.st.DecideAppeal(ctx, id, u.ID, grant, note); err != nil {
		http.Error(w, "could not record the decision", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/attendance", http.StatusSeeOther)
}

// attendanceSummary is the compact standing shown on the dashboard.
type attendanceSummary struct {
	Rate       int
	RunLength  int
	Limit      int
	BudgetLeft int
	AtRisk     bool
	Sanction   *store.Sanction
}

// attendanceFor builds a learner's standing for the dashboard.
func (s *Server) attendanceFor(ctx context.Context, userID int64) attendanceSummary {
	days, _ := s.st.AttendanceDays(ctx, userID)
	run := attendance.CurrentRun(days)
	sanction, _ := s.st.LatestSanction(ctx, userID)
	return attendanceSummary{
		Rate:       attendance.Rate(days),
		RunLength:  run.Length,
		Limit:      attendance.UnexcusedLimit,
		BudgetLeft: attendance.BudgetLeft(days),
		// Warn one day before the rule would bite, so the first a learner hears
		// of it is never the sanction itself.
		AtRisk:   run.Length > 0 && run.Length >= attendance.UnexcusedLimit-1,
		Sanction: sanction,
	}
}
