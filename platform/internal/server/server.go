// Package server wires routes, middleware, and handlers together.
package server

import (
	"net/http"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

type Server struct {
	st     *store.Store
	rnd    *renderer
	secure bool // set Secure cookies (true in production/HTTPS)
}

func New(st *store.Store, secure bool) (*Server, error) {
	rnd, err := newRenderer()
	if err != nil {
		return nil, err
	}
	return &Server{st: st, rnd: rnd, secure: secure}, nil
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
	mux.HandleFunc("GET /login", s.handleLoginForm)
	mux.HandleFunc("POST /login", s.handleLogin)
	mux.HandleFunc("GET /register", s.handleRegisterForm)
	mux.HandleFunc("POST /register", s.handleRegister)
	mux.HandleFunc("POST /logout", s.handleLogout)

	// Learner (auth required)
	mux.Handle("GET /dashboard", s.requireAuth(http.HandlerFunc(s.handleDashboard)))

	// Admin
	admin := s.requireRole("admin")
	mux.Handle("GET /admin", admin(http.HandlerFunc(s.handleAdminHome)))
	mux.Handle("GET /admin/users", admin(http.HandlerFunc(s.handleAdminUsers)))
	mux.Handle("POST /admin/users/{id}/role", admin(http.HandlerFunc(s.handleAdminSetRole)))
	mux.Handle("GET /admin/courses", admin(http.HandlerFunc(s.handleAdminCourses)))
	mux.Handle("GET /admin/grading", admin(http.HandlerFunc(s.handleAdminSoon("Grading queue"))))
	mux.Handle("GET /admin/competitions", admin(http.HandlerFunc(s.handleAdminSoon("Competitions"))))
	mux.Handle("GET /admin/forum", admin(http.HandlerFunc(s.handleAdminSoon("Forum"))))

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
