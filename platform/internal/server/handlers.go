package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "home.html", ViewData{Title: "GamifyDev — Learn to code by building"})
}

// --- Auth -------------------------------------------------------------------

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	if auth.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.render(w, r, "login.html", ViewData{Title: "Sign in"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	u, err := s.st.GetUserByEmail(r.Context(), email)
	if err != nil || !auth.CheckPassword(u.PasswordHash, password) {
		w.WriteHeader(http.StatusUnauthorized)
		s.render(w, r, "login.html", ViewData{Title: "Sign in", Flash: "Wrong email or password.", Data: map[string]any{"email": email}})
		return
	}
	s.startSession(w, r, u)
}

func (s *Server) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	if auth.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.render(w, r, "register.html", ViewData{Title: "Create your account"})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	fail := func(msg string) {
		w.WriteHeader(http.StatusBadRequest)
		s.render(w, r, "register.html", ViewData{Title: "Create your account", Flash: msg, Data: map[string]any{"name": name, "email": email}})
	}
	if len(name) < 2 || !strings.Contains(email, "@") || len(password) < 8 {
		fail("Enter a name, a valid email, and a password of at least 8 characters.")
		return
	}
	if _, err := s.st.GetUserByEmail(r.Context(), email); err == nil {
		fail("An account with that email already exists.")
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// First-ever user becomes the admin so the portal is reachable.
	role := "learner"
	if n, _ := s.st.CountUsers(r.Context()); n == 0 {
		role = "admin"
	}
	u, err := s.st.CreateUser(r.Context(), email, hash, name, role)
	if err != nil {
		fail("Could not create the account. Try a different email.")
		return
	}
	s.startSession(w, r, u)
}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, u *store.User) {
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
	dest := "/dashboard"
	if u.IsAdmin() {
		dest = "/admin"
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
	s.render(w, r, "dashboard.html", ViewData{Title: "Your dashboard"})
}

// --- Admin ------------------------------------------------------------------

func (s *Server) handleAdminHome(w http.ResponseWriter, r *http.Request) {
	n, _ := s.st.CountUsers(r.Context())
	s.render(w, r, "admin_home.html", ViewData{Title: "Admin", Data: map[string]any{"userCount": n}})
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.st.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "could not load users", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "admin_users.html", ViewData{Title: "User management", Data: map[string]any{"users": users}})
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
