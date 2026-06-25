// Package server wires routes, middleware, and handlers together.
package server

import (
	"net/http"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/email"
	"github.com/RA9/gamifydev/platform/internal/store"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	st     *store.Store
	rdb    *redis.Client // optional; nil when REDIS_URL is unset
	mail   *email.Mailer // optional; falls back to showing links when unconfigured
	rnd    *renderer
	secure bool // set Secure cookies (true in production/HTTPS)
}

func New(st *store.Store, rc *redis.Client, mail *email.Mailer, secure bool) (*Server, error) {
	rnd, err := newRenderer()
	if err != nil {
		return nil, err
	}
	if mail == nil {
		mail = email.New()
	}
	return &Server{st: st, rdb: rc, mail: mail, rnd: rnd, secure: secure}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", StaticHandler())
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public
	mux.HandleFunc("GET /{$}", s.handleHome)
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
	mux.HandleFunc("POST /logout", s.handleLogout)
	// Invitation acceptance (public — the invitee has no account yet).
	mux.HandleFunc("GET /invite/{token}", s.handleInvite)
	mux.HandleFunc("POST /invite/{token}", s.handleInviteSubmit)

	// Learner (auth required)
	in := s.requireAuth
	mux.Handle("GET /dashboard", in(http.HandlerFunc(s.handleDashboard)))
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
	mux.Handle("GET /admin/users", admin(http.HandlerFunc(s.handleAdminUsers)))
	mux.Handle("GET /admin/users/new", admin(http.HandlerFunc(s.handleAdminInviteForm)))
	mux.Handle("POST /admin/users/invite", admin(http.HandlerFunc(s.handleAdminInvite)))
	mux.Handle("POST /admin/users/{id}/role", admin(http.HandlerFunc(s.handleAdminSetRole)))
	mux.Handle("POST /admin/invitations/{id}/revoke", admin(http.HandlerFunc(s.handleAdminRevokeInvite)))
	mux.Handle("GET /admin/courses", admin(http.HandlerFunc(s.handleAdminCourses)))
	mux.Handle("GET /admin/courses/new", admin(http.HandlerFunc(s.handleAdminCourseNew)))
	mux.Handle("POST /admin/courses", admin(http.HandlerFunc(s.handleAdminCourseCreate)))
	mux.Handle("GET /admin/courses/{id}", admin(http.HandlerFunc(s.handleAdminCourseManage)))
	mux.Handle("POST /admin/courses/{id}", admin(http.HandlerFunc(s.handleAdminCourseUpdate)))
	mux.Handle("GET /admin/courses/{id}/lessons/new", admin(http.HandlerFunc(s.handleAdminLessonNew)))
	mux.Handle("POST /admin/courses/{id}/lessons", admin(http.HandlerFunc(s.handleAdminLessonCreate)))
	mux.Handle("GET /admin/lessons/{id}", admin(http.HandlerFunc(s.handleAdminLessonEdit)))
	mux.Handle("POST /admin/lessons/{id}", admin(http.HandlerFunc(s.handleAdminLessonUpdate)))
	mux.Handle("POST /admin/lessons/{id}/delete", admin(http.HandlerFunc(s.handleAdminLessonDelete)))
	mux.Handle("GET /admin/assignments", admin(http.HandlerFunc(s.handleAdminAssignments)))
	mux.Handle("GET /admin/assignments/new", admin(http.HandlerFunc(s.handleAdminAssignmentNew)))
	mux.Handle("POST /admin/assignments", admin(http.HandlerFunc(s.handleAdminAssignmentCreate)))
	mux.Handle("GET /admin/assignments/{id}", admin(http.HandlerFunc(s.handleAdminAssignmentEdit)))
	mux.Handle("POST /admin/assignments/{id}", admin(http.HandlerFunc(s.handleAdminAssignmentUpdate)))
	mux.Handle("GET /admin/competitions", admin(http.HandlerFunc(s.handleAdminCompetitions)))
	mux.Handle("GET /admin/competitions/new", admin(http.HandlerFunc(s.handleAdminCompetitionNew)))
	mux.Handle("POST /admin/competitions", admin(http.HandlerFunc(s.handleAdminCompetitionCreate)))
	mux.Handle("GET /admin/competitions/{slug}", admin(http.HandlerFunc(s.handleAdminCompetitionManage)))
	mux.Handle("POST /admin/competitions/{slug}/entries/{id}/score", admin(http.HandlerFunc(s.handleAdminScoreEntry)))
	mux.Handle("GET /admin/forum", admin(http.HandlerFunc(s.handleAdminForum)))
	mux.Handle("POST /admin/forum/{id}/{action}", admin(http.HandlerFunc(s.handleForumModerate)))
	mux.Handle("GET /admin/blog", admin(http.HandlerFunc(s.handleAdminBlog)))
	mux.Handle("GET /admin/blog/new", admin(http.HandlerFunc(s.handleAdminBlogNew)))
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
