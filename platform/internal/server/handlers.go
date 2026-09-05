package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/attendance"
	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/cohort"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	// The curriculum figures come from the database so the page cannot drift
	// from the thing it describes.
	stats, err := s.st.PublicCurriculumStats(r.Context())
	if err != nil {
		log.Printf("home: curriculum stats: %v", err)
	}
	s.render(w, r, "home.html", ViewData{
		Title: "GamifyDev — a free CS program that expects you to show up",
		Data: map[string]any{
			"stats": stats,
			// Stated on the page, so they are read from the policy packages
			// rather than retyped into the copy where they would rot.
			"cohortMax":     cohort.MaxSize,
			"missedLimit":   attendance.UnexcusedLimit,
			"excusedBudget": attendance.ExcusedBudget,
			"reapplyDays":   int(attendance.ReapplyAfter.Hours() / 24),
		},
	})
}

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "about.html", ViewData{Title: "About us"})
}

// --- Auth -------------------------------------------------------------------

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	if auth.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	vd := ViewData{Title: "Sign in"}
	if r.URL.Query().Get("reset") == "1" {
		vd.Flash = "Password changed. Sign in with your new one."
	}
	s.render(w, r, "login.html", vd)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	// Two keys, because the two attacks are different shapes. Keying on the
	// address stops one account being ground down from many hosts; keying on
	// the address of the caller stops one host spraying a password across many
	// accounts. Either one tripping is enough to refuse.
	ipKey := "login-ip:" + s.clientIP(r)
	userKey := "login-user:" + email
	if wait := longer(s.authlim.retryAfter(ipKey), s.authlim.retryAfter(userKey)); wait > 0 {
		w.WriteHeader(http.StatusTooManyRequests)
		s.render(w, r, "login.html", ViewData{Title: "Sign in",
			Flash: "Too many sign-in attempts. Try again in " + humanWait(wait) + ".",
			Data:  map[string]any{"email": email}})
		return
	}

	u, err := s.st.GetUserByEmail(r.Context(), email)
	if err != nil {
		// Spend the same time bcrypt would have spent on a real hash. Returning
		// early here is measurably faster than a wrong password, which turns
		// the login form into a way to test whether an address is registered.
		auth.BurnPasswordCheck(password)
	}
	if err != nil || !auth.CheckPassword(u.PasswordHash, password) {
		s.authlim.fail(ipKey)
		s.authlim.fail(userKey)
		w.WriteHeader(http.StatusUnauthorized)
		s.render(w, r, "login.html", ViewData{Title: "Sign in", Flash: "Wrong email or password.", Data: map[string]any{"email": email}})
		return
	}
	// A correct password forgets the misses that came before it, so someone who
	// mistyped twice this morning isn't closer to a lockout this afternoon.
	s.authlim.succeed(ipKey)
	s.authlim.succeed(userKey)
	s.startSession(w, r, u)
}

func (s *Server) passingGuestPlacement(r *http.Request) (int64, *store.PlacementResult, error) {
	cookie, err := r.Cookie(auth.GuestCookie)
	if err != nil || cookie.Value == "" {
		return 0, nil, store.ErrPlacementNotPassed
	}
	guestID, err := s.st.GuestByToken(r.Context(), cookie.Value)
	if err != nil {
		return 0, nil, store.ErrPlacementNotPassed
	}
	result, err := s.st.LatestPlacementFor(r.Context(), store.Solver{GuestID: guestID})
	if err != nil {
		return 0, nil, err
	}
	if result == nil || !result.Passed {
		return 0, nil, store.ErrPlacementNotPassed
	}
	return guestID, result, nil
}

func (s *Server) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	if auth.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	_, result, err := s.passingGuestPlacement(r)
	if err != nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	var path *store.Path
	if result.RecommendedPathID.Valid {
		path, _ = s.st.GetPathByID(r.Context(), result.RecommendedPathID.Int64)
	}
	s.render(w, r, "register.html", ViewData{Title: "Create your account", Data: map[string]any{"path": path}})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	guestID, result, err := s.passingGuestPlacement(r)
	if err != nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")
	var path *store.Path
	if result.RecommendedPathID.Valid {
		path, _ = s.st.GetPathByID(r.Context(), result.RecommendedPathID.Int64)
	}

	fail := func(msg string) {
		w.WriteHeader(http.StatusBadRequest)
		s.render(w, r, "register.html", ViewData{Title: "Create your account", Flash: msg, Data: map[string]any{
			"name": name, "email": email, "path": path,
		}})
	}
	if len(name) < 2 || !strings.Contains(email, "@") || len(password) < 8 {
		fail("Enter a name, a valid email, and a password of at least 8 characters.")
		return
	}
	if _, err := s.st.GetUserByEmail(r.Context(), email); err == nil {
		fail("An account with that email already exists. Sign in instead.")
		return
	} else if !errors.Is(err, store.ErrNotFound) {
		log.Printf("register: lookup %q failed: %v", email, err)
		fail("Could not create the account right now. Please try again.")
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// Public signup can only create a learner, and the account/result/work transfer
	// commits as one transaction. There is no first-user admin promotion.
	u, err := s.st.CreateLearnerFromGuestPlacement(r.Context(), guestID, email, hash, name)
	if err != nil {
		log.Printf("register: claim placement for %q failed: %v", email, err)
		if errors.Is(err, store.ErrPlacementNotPassed) || errors.Is(err, store.ErrGuestAlreadyClaimed) {
			http.Redirect(w, r, "/placement", http.StatusSeeOther)
			return
		}
		fail("Could not create the account. Try a different email.")
		return
	}
	auth.ClearGuestCookie(w, s.secure)
	s.startSessionAt(w, r, u, "/placement/result")
}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, u *store.User) {
	// Existing accounts may still claim work completed as a guest. New learner
	// registration has already claimed it atomically before calling startSessionAt.
	s.claimGuestWork(w, r, u.ID)
	s.startSessionAt(w, r, u, "")
}

func (s *Server) startSessionAt(w http.ResponseWriter, r *http.Request, u *store.User, dest string) {
	token, err := auth.NewToken()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if err := s.st.CreateSession(r.Context(), token, u.ID, time.Now().Add(auth.SessionTTL)); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	auth.SetSessionCookie(w, token, s.secure)
	if dest == "" {
		dest = "/dashboard"
		if u.IsAdmin() {
			dest = "/admin"
		} else if u.CanGrade() {
			dest = "/admin/grading"
		}
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(auth.SessionCookie); err == nil {
		_ = s.st.DeleteSession(r.Context(), c.Value)
	}
	auth.ClearSessionCookie(w, s.secure)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Learner ----------------------------------------------------------------

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := s.dashboardData(r.Context())
	if err != nil {
		http.Error(w, "could not load your dashboard", http.StatusInternalServerError)
		return
	}
	data["bodyClass"] = "dash-dark"
	s.render(w, r, "dashboard.html", ViewData{Title: "Your dashboard", Data: data})
}

// handleDashboardLive renders just the realtime region of the dashboard. htmx
// fetches it on a "refresh" signal pushed over the WebSocket (e.g. after a
// submission is graded), so the cards/charts/table update without a reload.
func (s *Server) handleDashboardLive(w http.ResponseWriter, r *http.Request) {
	data, err := s.dashboardData(r.Context())
	if err != nil {
		http.Error(w, "could not load your dashboard", http.StatusInternalServerError)
		return
	}
	vd := ViewData{User: auth.CurrentUser(r.Context()), Path: r.URL.Path, Data: data}
	s.renderPartial(w, r, "dashboard.html", "dashLive", vd)
}

// dashboardData assembles every figure the learner dashboard shows.
func (s *Server) dashboardData(ctx context.Context) (map[string]any, error) {
	u := auth.CurrentUser(ctx)
	if u == nil {
		return nil, errors.New("no user")
	}
	st, err := s.st.LearnerStats(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	series, err := s.st.LearnerSubmissionSeries(ctx, u.ID, 14)
	if err != nil {
		return nil, err
	}
	recent, err := s.st.RecentSubmissionsByUser(ctx, u.ID, 6)
	if err != nil {
		return nil, err
	}
	paths, _ := s.st.ListPaths(ctx, false)
	if len(paths) > 4 {
		paths = paths[:4]
	}
	continueLesson, _ := s.st.ContinueLearning(ctx, u.ID)

	// Placement state drives the dashboard's top card: an unplaced learner is
	// pointed at the diagnostic rather than at content they cannot yet enroll in.
	enrollment, _ := s.st.EnsureEnrollment(ctx, u.ID)
	placed, _ := s.st.LatestPlacement(ctx, u.ID)
	var enrolledPath *store.Path
	if enrollment != nil && enrollment.PathID.Valid {
		enrolledPath, _ = s.st.GetPathByID(ctx, enrollment.PathID.Int64)
	}

	// Cohort state: whether today's standup is still waiting on this learner is
	// the single most actionable thing the dashboard can say.
	myCohort, _ := s.st.CohortForUser(ctx, u.ID)
	standupDone, standupOpen := false, false
	var dueToday, overdue []store.ScheduleItem
	var progress store.ScheduleProgress
	var attend attendanceSummary
	if myCohort != nil {
		// Phase 4: the dashboard leads with today's work rather than "pick up
		// where you left off" — that is what ending self-pacing looks like.
		today := time.Now().UTC().Format("2006-01-02")
		dueToday, _ = s.st.DueOn(ctx, myCohort.ID, u.ID, today)
		overdue, _ = s.st.Overdue(ctx, myCohort.ID, u.ID, 5)
		progress, _ = s.st.Progress(ctx, myCohort.ID, u.ID)
		attend = s.attendanceFor(ctx, u.ID)
		band := cohort.LookupBand(myCohort.TZBand)
		day := cohort.LocalDay(time.Now().UTC(), band).Format("2006-01-02")
		if su, _ := s.st.TodayStandup(ctx, myCohort.ID, day); su != nil {
			standupOpen = su.IsOpen
			standupDone, _ = s.st.HasPosted(ctx, su.ID, u.ID)
		}
	}

	periodTotal := 0
	for _, d := range series {
		periodTotal += d.Count
	}

	// Submission-status donut.
	pct := func(n int) int {
		if st.TotalSubmissions == 0 {
			return 0
		}
		return int(float64(n) / float64(st.TotalSubmissions) * 100)
	}
	segments := []DonutSegment{
		{Label: "Graded", Count: st.Graded, Class: "donut-learner", Pct: pct(st.Graded)},
		{Label: "Awaiting review", Count: st.Pending, Class: "donut-grader", Pct: pct(st.Pending)},
		{Label: "Needs rework", Count: st.Returned, Class: "donut-admin", Pct: pct(st.Returned)},
	}

	return map[string]any{
		"stats":        st,
		"chart":        areaChartSVG(series),
		"periodTotal":  periodTotal,
		"donut":        donutSVG(st.TotalSubmissions, "submissions", segments),
		"segments":     segments,
		"recent":       recent,
		"paths":        paths,
		"continue":     continueLesson,
		"enrollment":   enrollment,
		"placed":       placed != nil,
		"enrolledPath": enrolledPath,
		"cohort":       myCohort,
		"standupDone":  standupDone,
		"standupOpen":  standupOpen,
		"dueToday":     dueToday,
		"overdue":      overdue,
		"schedule":     progress,
		"attendance":   attend,
	}, nil
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, _ := s.st.AdminStats(ctx)
	series, _ := s.st.SignupSeries(ctx, 14)
	recent, _ := s.st.RecentUsers(ctx, 6)
	courses, _ := s.st.ListCourses(ctx, true)

	// Period total for the chart card header.
	periodTotal := 0
	for _, d := range series {
		periodTotal += d.Count
	}

	// Role donut segments (skip-zero handled by the renderer).
	pct := func(n int) int {
		if stats.TotalUsers == 0 {
			return 0
		}
		return int(float64(n)/float64(stats.TotalUsers)*100 + 0.5)
	}
	segments := []DonutSegment{
		{Label: "Learners", Count: stats.Learners, Class: "donut-learner", Pct: pct(stats.Learners)},
		{Label: "Graders", Count: stats.Graders, Class: "donut-grader", Pct: pct(stats.Graders)},
		{Label: "Admins", Count: stats.Admins, Class: "donut-admin", Pct: pct(stats.Admins)},
	}

	s.render(w, r, "admin_home.html", ViewData{
		Title: "Overview",
		Data: map[string]any{
			"stats":       stats,
			"chart":       areaChartSVG(series),
			"donut":       donutSVG(stats.TotalUsers, "members", segments),
			"segments":    segments,
			"recent":      recent,
			"courses":     courses,
			"periodTotal": periodTotal,
			"today":       time.Now().Format("Monday, Jan 2"),
		},
	})
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.st.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "could not load users", http.StatusInternalServerError)
		return
	}
	pending, _ := s.st.ListPendingInvitations(r.Context())
	s.render(w, r, "admin_users.html", ViewData{Title: "User management",
		Data: map[string]any{"users": users, "pending": pending}})
}

func (s *Server) handleAdminSetRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	role := r.FormValue("role")
	if role != "admin" && role != "grader" && role != "learner" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}
	if err := s.st.SetUserRole(r.Context(), id, role); err != nil {
		http.Error(w, "could not update role", http.StatusInternalServerError)
		return
	}
	// htmx swaps just the updated row; fall back to a full redirect otherwise.
	if r.Header.Get("HX-Request") == "true" {
		u, err := s.st.GetUserByID(r.Context(), id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		s.renderPartial(w, r, "admin_users.html", "userRow", u)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
