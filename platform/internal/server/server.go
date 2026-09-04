// Package server wires routes, middleware, and handlers together.
package server

import (
	"net/http"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/email"
	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	st     *store.Store
	rdb    *redis.Client // optional; nil when REDIS_URL is unset
	mail   *email.Mailer // optional; falls back to showing links when unconfigured
	hub    *hub          // realtime WebSocket fan-out
	rnd    *renderer
	exec   runner.Executor // server-side code execution; disabled by default
	jobs   *jobs.Runner    // scheduled work; nil when JOBS=off
	runlim *runLimiter     // throttles /api/run
	// authlim throttles failed sign-ins and reset requests. Separate from
	// runlim because its keys are attacker-supplied rather than user ids.
	authlim *attemptLimiter
	// Attendance enforcement; false means sanctions are recorded but never
	// applied. Mirrors what was passed to jobs.Register.
	enforceAttendance bool
	secure            bool // set Secure cookies (true in production/HTTPS)
}

func New(st *store.Store, rc *redis.Client, mail *email.Mailer, exec runner.Executor, jr *jobs.Runner, enforceAttendance, secure bool) (*Server, error) {
	rnd, err := newRenderer()
	if err != nil {
		return nil, err
	}
	if mail == nil {
		mail = email.New()
	}
	if exec == nil {
		exec, _ = runner.New(runner.Config{}) // disabled
	}
	return &Server{
		st: st, rdb: rc, mail: mail, hub: newHub(), rnd: rnd,
		exec: exec, jobs: jr, runlim: newRunLimiter(4, 1500),
		// Five misses inside fifteen minutes, then a fifteen-minute cool-off:
		// generous for a person who forgot which password they used, ruinous
		// for anything working through a list.
		authlim:           newAttemptLimiter(5, 15*time.Minute, 15*time.Minute),
		enforceAttendance: enforceAttendance, secure: secure,
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", StaticHandler())
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public
	mux.HandleFunc("GET /{$}", s.handleHome)
	mux.HandleFunc("GET /about", s.handleAbout)
	mux.HandleFunc("GET /paths", s.handlePaths)
	mux.HandleFunc("GET /paths/{slug}", s.handlePath)
	mux.HandleFunc("GET /courses", s.handleCourses)
	mux.HandleFunc("GET /courses/{course}", s.handleCourse)
	mux.HandleFunc("GET /courses/{course}/{lesson}", s.handleLesson)
	mux.HandleFunc("GET /forum", s.handleForum)
	mux.HandleFunc("GET /forum/{id}", s.handleForumThread)
	mux.HandleFunc("GET /competitions", s.handleCompetitions)
	mux.HandleFunc("GET /competitions/{slug}", s.handleCompetition)
	mux.HandleFunc("GET /blog", s.handleBlog)
	mux.HandleFunc("GET /blog/{slug}", s.handleBlogPost)
	mux.HandleFunc("GET /login", s.handleLoginForm)
	mux.HandleFunc("POST /login", s.handleLogin)
	mux.HandleFunc("GET /register", s.handleRegisterForm)
	mux.HandleFunc("POST /register", s.handleRegister)
	mux.HandleFunc("GET /forgot", s.handleForgotForm)
	mux.HandleFunc("POST /forgot", s.handleForgot)
	mux.HandleFunc("GET /reset/{token}", s.handleResetForm)
	mux.HandleFunc("POST /reset/{token}", s.handleReset)
	mux.HandleFunc("POST /logout", s.handleLogout)
	// Invitation acceptance (public — the invitee has no account yet).
	mux.HandleFunc("GET /invite/{token}", s.handleInvite)
	mux.HandleFunc("POST /invite/{token}", s.handleInviteSubmit)

	// Learner (auth required)
	in := s.requireAuth
	mux.Handle("POST /steps/{id}/complete", in(http.HandlerFunc(s.handleStepComplete)))
	mux.Handle("POST /steps/{id}/run", in(http.HandlerFunc(s.handleStepRun)))
	// Placement diagnostic and the enrollment gate it feeds.
	mux.Handle("GET /placement", in(http.HandlerFunc(s.handlePlacement)))
	mux.Handle("POST /placement/start", in(http.HandlerFunc(s.handlePlacementStart)))
	mux.Handle("POST /placement/submit", in(http.HandlerFunc(s.handlePlacementSubmit)))
	mux.Handle("GET /placement/result", in(http.HandlerFunc(s.handlePlacementResult)))
	mux.Handle("GET /placement/review", in(http.HandlerFunc(s.handlePlacementReview)))
	mux.Handle("POST /enroll", in(http.HandlerFunc(s.handleEnroll)))
	// Cohort space and the daily standup.
	mux.Handle("GET /cohort", in(http.HandlerFunc(s.handleCohort)))
	mux.Handle("POST /cohort/standup", in(http.HandlerFunc(s.handleStandupPost)))
	mux.Handle("GET /schedule", in(http.HandlerFunc(s.handleSchedule)))
	mux.Handle("GET /attendance", in(http.HandlerFunc(s.handleAttendance)))
	mux.Handle("POST /attendance/absence", in(http.HandlerFunc(s.handleAbsenceFile)))
	mux.Handle("POST /attendance/appeal", in(http.HandlerFunc(s.handleAppealFile)))
	mux.Handle("POST /lessons/{id}/complete", in(http.HandlerFunc(s.handleLessonComplete)))
	mux.Handle("GET /dashboard", in(http.HandlerFunc(s.handleDashboard)))
	mux.Handle("GET /dashboard/live", in(http.HandlerFunc(s.handleDashboardLive)))
	mux.Handle("GET /ws", in(http.HandlerFunc(s.handleWS)))
	mux.Handle("GET /playground", in(http.HandlerFunc(s.handlePlayground)))
	mux.Handle("POST /api/run", in(http.HandlerFunc(s.handleRunCode)))
	mux.Handle("GET /assignments", in(http.HandlerFunc(s.handleAssignments)))
	mux.Handle("GET /assignments/{slug}", in(http.HandlerFunc(s.handleAssignment)))
	mux.Handle("POST /assignments/{slug}/submit", in(http.HandlerFunc(s.handleSubmitAssignment)))
	mux.Handle("GET /forum/new", in(http.HandlerFunc(s.handleForumNewForm)))
	mux.Handle("POST /forum", in(http.HandlerFunc(s.handleForumCreate)))
	mux.Handle("POST /forum/{id}/reply", in(http.HandlerFunc(s.handleForumReply)))
	mux.Handle("POST /competitions/{slug}/submit", in(http.HandlerFunc(s.handleCompetitionSubmit)))

	// Admin
	admin := s.requireRole("admin")
	mux.Handle("GET /admin", admin(http.HandlerFunc(s.handleAdminHome)))
	mux.Handle("GET /admin/cohorts", admin(http.HandlerFunc(s.handleAdminCohorts)))
	mux.Handle("GET /admin/attendance", admin(http.HandlerFunc(s.handleAdminAttendance)))
	mux.Handle("POST /admin/appeals/{id}", admin(http.HandlerFunc(s.handleAdminAppealDecide)))
	mux.Handle("GET /admin/jobs", admin(http.HandlerFunc(s.handleAdminJobs)))
	mux.Handle("POST /admin/jobs/{name}/run", admin(http.HandlerFunc(s.handleAdminJobRun)))
	mux.Handle("GET /admin/users", admin(http.HandlerFunc(s.handleAdminUsers)))
	mux.Handle("GET /admin/users/new", admin(http.HandlerFunc(s.handleAdminInviteForm)))
	mux.Handle("POST /admin/users/invite", admin(http.HandlerFunc(s.handleAdminInvite)))
	mux.Handle("POST /admin/users/{id}/role", admin(http.HandlerFunc(s.handleAdminSetRole)))
	mux.Handle("POST /admin/invitations/{id}/revoke", admin(http.HandlerFunc(s.handleAdminRevokeInvite)))
	mux.Handle("GET /admin/paths", admin(http.HandlerFunc(s.handleAdminPaths)))
	mux.Handle("GET /admin/paths/new", admin(http.HandlerFunc(s.handleAdminPathNew)))
	mux.Handle("POST /admin/paths", admin(http.HandlerFunc(s.handleAdminPathCreate)))
	mux.Handle("GET /admin/paths/{id}", admin(http.HandlerFunc(s.handleAdminPathEdit)))
	mux.Handle("POST /admin/paths/{id}", admin(http.HandlerFunc(s.handleAdminPathUpdate)))
	mux.Handle("POST /admin/paths/{id}/delete", admin(http.HandlerFunc(s.handleAdminPathDelete)))
	mux.Handle("GET /admin/courses", admin(http.HandlerFunc(s.handleAdminCourses)))
	mux.Handle("GET /admin/courses/new", admin(http.HandlerFunc(s.handleAdminCourseNew)))
	mux.Handle("POST /admin/courses", admin(http.HandlerFunc(s.handleAdminCourseCreate)))
	mux.Handle("GET /admin/courses/{id}", admin(http.HandlerFunc(s.handleAdminCourseManage)))
	mux.Handle("POST /admin/courses/{id}", admin(http.HandlerFunc(s.handleAdminCourseUpdate)))
	mux.Handle("GET /admin/courses/{id}/lessons/new", admin(http.HandlerFunc(s.handleAdminLessonNew)))
	mux.Handle("POST /admin/courses/{id}/lessons", admin(http.HandlerFunc(s.handleAdminLessonCreate)))
	mux.Handle("GET /admin/courses/{id}/assignments/new", admin(http.HandlerFunc(s.handleAdminCourseAssignmentNew)))
	mux.Handle("POST /admin/courses/{id}/assignments", admin(http.HandlerFunc(s.handleAdminCourseAssignmentCreate)))
	mux.Handle("GET /admin/lessons/{id}", admin(http.HandlerFunc(s.handleAdminLessonEdit)))
	mux.Handle("POST /admin/lessons/{id}", admin(http.HandlerFunc(s.handleAdminLessonUpdate)))
	mux.Handle("POST /admin/lessons/{id}/delete", admin(http.HandlerFunc(s.handleAdminLessonDelete)))
	mux.Handle("POST /admin/lessons/preview", admin(http.HandlerFunc(s.handleAdminLessonPreview)))
	mux.Handle("GET /admin/lessons/{id}/steps/new", admin(http.HandlerFunc(s.handleAdminStepNew)))
	mux.Handle("POST /admin/lessons/{id}/steps", admin(http.HandlerFunc(s.handleAdminStepCreate)))
	mux.Handle("GET /admin/steps/{id}", admin(http.HandlerFunc(s.handleAdminStepEdit)))
	mux.Handle("POST /admin/steps/{id}", admin(http.HandlerFunc(s.handleAdminStepUpdate)))
	mux.Handle("POST /admin/steps/{id}/delete", admin(http.HandlerFunc(s.handleAdminStepDelete)))
	mux.Handle("POST /admin/steps/{id}/{dir}", admin(http.HandlerFunc(s.handleAdminStepMove)))
	mux.Handle("GET /admin/assignments", admin(http.HandlerFunc(s.handleAdminAssignments)))
	mux.Handle("GET /admin/assignments/new", admin(http.HandlerFunc(s.handleAdminAssignmentNew)))
	mux.Handle("POST /admin/assignments", admin(http.HandlerFunc(s.handleAdminAssignmentCreate)))
	mux.Handle("GET /admin/assignments/{id}", admin(http.HandlerFunc(s.handleAdminAssignmentEdit)))
	mux.Handle("POST /admin/assignments/{id}", admin(http.HandlerFunc(s.handleAdminAssignmentUpdate)))
	mux.Handle("POST /admin/assignments/{id}/delete", admin(http.HandlerFunc(s.handleAdminAssignmentDelete)))
	mux.Handle("GET /admin/competitions", admin(http.HandlerFunc(s.handleAdminCompetitions)))
	mux.Handle("GET /admin/competitions/new", admin(http.HandlerFunc(s.handleAdminCompetitionNew)))
	mux.Handle("POST /admin/competitions", admin(http.HandlerFunc(s.handleAdminCompetitionCreate)))
	mux.Handle("GET /admin/competitions/{slug}", admin(http.HandlerFunc(s.handleAdminCompetitionManage)))
	mux.Handle("POST /admin/competitions/{slug}/entries/{id}/score", admin(http.HandlerFunc(s.handleAdminScoreEntry)))
	mux.Handle("GET /admin/forum", admin(http.HandlerFunc(s.handleAdminForum)))
	mux.Handle("POST /admin/forum/{id}/{action}", admin(http.HandlerFunc(s.handleForumModerate)))
	mux.Handle("GET /admin/blog", admin(http.HandlerFunc(s.handleAdminBlog)))
	mux.Handle("GET /admin/blog/new", admin(http.HandlerFunc(s.handleAdminBlogNew)))
	mux.Handle("POST /admin/blog/preview", admin(http.HandlerFunc(s.handleAdminBlogPreview)))
	mux.Handle("POST /admin/blog", admin(http.HandlerFunc(s.handleAdminBlogCreate)))
	mux.Handle("GET /admin/blog/{id}", admin(http.HandlerFunc(s.handleAdminBlogEdit)))
	mux.Handle("POST /admin/blog/{id}", admin(http.HandlerFunc(s.handleAdminBlogUpdate)))
	mux.Handle("POST /admin/blog/{id}/delete", admin(http.HandlerFunc(s.handleAdminBlogDelete)))

	// Grading — admins and graders
	staff := s.requireRole("admin", "grader")
	mux.Handle("GET /admin/grading", staff(http.HandlerFunc(s.handleAdminGrading)))
	mux.Handle("GET /admin/grading/{id}", staff(http.HandlerFunc(s.handleAdminGradeForm)))
	mux.Handle("POST /admin/grading/{id}", staff(http.HandlerFunc(s.handleAdminGrade)))

	// loadUser runs on every request so templates know who's logged in.
	return s.loadUser(mux)
}

// loadUser attaches the current user (if any) to the request context.
func (s *Server) loadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(auth.SessionCookie); err == nil && c.Value != "" {
			if u, err := s.st.UserBySession(r.Context(), c.Value); err == nil {
				r = r.WithContext(auth.WithUser(r.Context(), u))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.CurrentUser(r.Context()) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := auth.CurrentUser(r.Context())
			if u == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			if !allowed[u.Role] {
				http.Error(w, "403 — you don't have access to this area.", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
