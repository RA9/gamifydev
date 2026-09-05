package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

const inviteTTL = 7 * 24 * time.Hour

func validInviteRole(role string) bool {
	// Learners enter through placement. Invitations are reserved for staff so an
	// admin link cannot become a second, untested learner-admission path.
	return role == "admin" || role == "grader"
}

// baseURL builds the public origin (scheme://host) for the current request,
// honoring the proxy headers Railway sets in front of the app.
func (s *Server) baseURL(r *http.Request) string {
	scheme := "http"
	if s.secure || r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := r.Host
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		host = h
	}
	return scheme + "://" + host
}

// --- Admin: invite a member -------------------------------------------------

func (s *Server) handleAdminInviteForm(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin_user_invite.html", ViewData{Title: "Invite a member", Data: map[string]any{"role": "grader"}})
}

func (s *Server) handleAdminInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	emailAddr := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	name := strings.TrimSpace(r.FormValue("name"))
	role := r.FormValue("role")
	if !validInviteRole(role) {
		role = "grader"
	}

	fail := func(msg string) {
		w.WriteHeader(http.StatusBadRequest)
		s.render(w, r, "admin_user_invite.html", ViewData{Title: "Invite a member", Flash: msg,
			Data: map[string]any{"email": emailAddr, "name": name, "role": role}})
	}
	if !strings.Contains(emailAddr, "@") {
		fail("Enter a valid email address.")
		return
	}
	if _, err := s.st.GetUserByEmail(ctx, emailAddr); err == nil {
		fail("A user with that email already exists.")
		return
	}

	token, err := auth.NewToken()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	var invitedBy int64
	if u := auth.CurrentUser(ctx); u != nil {
		invitedBy = u.ID
	}
	inv, err := s.st.CreateInvitation(ctx, emailAddr, name, role, token, invitedBy, time.Now().Add(inviteTTL))
	if err != nil {
		log.Printf("invite: create %q failed: %v", emailAddr, err)
		fail("Could not create the invitation. Please try again.")
		return
	}

	link := s.baseURL(r) + "/invite/" + token
	sent := false
	if s.mail.Configured() {
		if err := s.mail.Send(emailAddr, "You're invited to GamifyDev", inviteEmailBody(name, role, link)); err != nil {
			log.Printf("invite: send to %q failed: %v", emailAddr, err)
		} else {
			sent = true
		}
	}

	// Re-render the user list with the result (the link is shown when it wasn't emailed).
	users, _ := s.st.ListUsers(ctx)
	pending, _ := s.st.ListPendingInvitations(ctx)
	s.render(w, r, "admin_users.html", ViewData{Title: "User management", Data: map[string]any{
		"users":          users,
		"pending":        pending,
		"invited":        inv,
		"inviteLink":     link,
		"sent":           sent,
		"mailConfigured": s.mail.Configured(),
	}})
}

func (s *Server) handleAdminRevokeInvite(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := s.st.DeleteInvitation(r.Context(), id); err != nil {
		http.Error(w, "could not revoke invitation", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// --- Public: accept an invitation -------------------------------------------

func (s *Server) handleInvite(w http.ResponseWriter, r *http.Request) {
	inv, ok := s.validInvitation(r)
	if !ok {
		s.render(w, r, "invite_invalid.html", ViewData{Title: "Invitation"})
		return
	}
	if inv.Role == "learner" {
		s.render(w, r, "invite_invalid.html", ViewData{Title: "Placement required",
			Flash: "Learner invitations have been replaced by the placement process. Pass placement to create your account."})
		return
	}
	s.render(w, r, "invite.html", ViewData{Title: "Set your password",
		Data: map[string]any{"inv": inv, "token": inv.Token}})
}

func (s *Server) handleInviteSubmit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	inv, ok := s.validInvitation(r)
	if !ok {
		s.render(w, r, "invite_invalid.html", ViewData{Title: "Invitation"})
		return
	}
	if inv.Role == "learner" {
		s.render(w, r, "invite_invalid.html", ViewData{Title: "Placement required",
			Flash: "Learner invitations have been replaced by the placement process. Pass placement to create your account."})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = inv.Name
	}
	password := r.FormValue("password")
	if len(name) < 2 || len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		s.render(w, r, "invite.html", ViewData{Title: "Set your password",
			Flash: "Enter your name and a password of at least 8 characters.",
			Data:  map[string]any{"inv": inv, "token": inv.Token}})
		return
	}

	// Guard against the email having been claimed since the invite was sent.
	if _, err := s.st.GetUserByEmail(ctx, inv.Email); err == nil {
		s.render(w, r, "invite_invalid.html", ViewData{Title: "Invitation",
			Flash: "An account with this email already exists — please sign in."})
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	u, err := s.st.CreateUser(ctx, inv.Email, hash, name, inv.Role)
	if err != nil {
		log.Printf("invite: accept create %q failed: %v", inv.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		s.render(w, r, "invite.html", ViewData{Title: "Set your password",
			Flash: "Could not create your account. Please try again.",
			Data:  map[string]any{"inv": inv, "token": inv.Token}})
		return
	}
	_ = s.st.AcceptInvitation(ctx, inv.Token)
	s.startSession(w, r, u) // logs the new user in and redirects by role
}

// validInvitation loads the invitation for the request token and reports whether
// it is usable (exists, not accepted, not expired).
func (s *Server) validInvitation(r *http.Request) (*store.Invitation, bool) {
	inv, err := s.st.GetInvitationByToken(r.Context(), r.PathValue("token"))
	if err != nil || inv.Accepted() || inv.Expired() {
		return nil, false
	}
	return inv, true
}

func inviteEmailBody(name, role, link string) string {
	greeting := "Hello"
	if name != "" {
		greeting = "Hello " + name
	}
	return fmt.Sprintf(`%s,

You've been invited to join GamifyDev as a %s.

Set your password and activate your account here:

%s

This link expires in 7 days. If you weren't expecting this, you can ignore this email.

— The GamifyDev team
`, greeting, role, link)
}
