package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// resetTTL is deliberately short. The link is a bearer key to the account, and
// it usually travels through an inbox the learner may not control forever.
const resetTTL = time.Hour

// minPasswordLen matches the rule the registration form enforces. Reset is a
// second door into the same account and must not be the weaker one.
const minPasswordLen = 8

func (s *Server) handleForgotForm(w http.ResponseWriter, r *http.Request) {
	if auth.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.render(w, r, "forgot.html", ViewData{Title: "Reset your password"})
}

// handleForgot issues a reset link.
//
// It answers identically whether or not the address has an account. A form that
// says "no such user" is a free membership oracle, and this one is public.
func (s *Server) handleForgot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addr := strings.ToLower(strings.TrimSpace(r.FormValue("email")))

	// Sending mail on demand from an anonymous endpoint is abusable in its own
	// right — as a way to flood one person's inbox, and as a way to burn the
	// sending reputation of the domain.
	ipKey := "forgot-ip:" + s.clientIP(r)
	addrKey := "forgot-addr:" + addr
	if wait := longer(s.authlim.retryAfter(ipKey), s.authlim.retryAfter(addrKey)); wait > 0 {
		w.WriteHeader(http.StatusTooManyRequests)
		s.render(w, r, "forgot.html", ViewData{Title: "Reset your password",
			Flash: "Too many requests. Try again in " + humanWait(wait) + "."})
		return
	}
	s.authlim.fail(ipKey)
	s.authlim.fail(addrKey)

	// Rendered on every path below, whatever actually happened.
	sent := func() {
		s.render(w, r, "forgot.html", ViewData{Title: "Reset your password",
			Data: map[string]any{"sent": true, "email": addr}})
	}
	if !strings.Contains(addr, "@") {
		sent()
		return
	}
	u, err := s.st.GetUserByEmail(ctx, addr)
	if err != nil {
		sent() // no account: say the same thing, do nothing
		return
	}
	token, err := auth.NewToken()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if err := s.st.CreatePasswordReset(ctx, u.ID, token, time.Now().Add(resetTTL)); err != nil {
		log.Printf("forgot: create reset for user %d failed: %v", u.ID, err)
		sent()
		return
	}

	link := s.baseURL(r) + "/reset/" + token
	if s.mail.Configured() {
		if err := s.mail.Send(addr, "Reset your GamifyDev password", resetEmailBody(u.Name, link)); err != nil {
			log.Printf("forgot: send to %q failed: %v", addr, err)
		}
	} else {
		// No SMTP: the link goes to the server log, where an operator can find
		// it — never to the page. Showing it in the browser would turn "type
		// any address" into "take over that account".
		log.Printf("forgot: SMTP not configured; reset link for %q: %s", addr, link)
	}
	sent()
}

func (s *Server) handleResetForm(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if _, err := s.st.UserByResetToken(r.Context(), token); err != nil {
		s.render(w, r, "reset_invalid.html", ViewData{Title: "Link expired"})
		return
	}
	s.render(w, r, "reset.html", ViewData{Title: "Choose a new password",
		Data: map[string]any{"token": token}})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	fail := func(msg string) {
		w.WriteHeader(http.StatusBadRequest)
		s.render(w, r, "reset.html", ViewData{Title: "Choose a new password", Flash: msg,
			Data: map[string]any{"token": token}})
	}
	if len(password) < minPasswordLen {
		fail(fmt.Sprintf("Use at least %d characters.", minPasswordLen))
		return
	}
	if password != confirm {
		fail("Those two passwords don't match.")
		return
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	userID, err := s.st.ResetPassword(r.Context(), token, hash)
	if errors.Is(err, store.ErrResetSpent) {
		s.render(w, r, "reset_invalid.html", ViewData{Title: "Link expired"})
		return
	} else if err != nil {
		log.Printf("reset: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Every session was just deleted, this browser's included if it had one.
	// Clear the cookie too so the next page doesn't present a dead token.
	auth.ClearSessionCookie(w, s.secure)
	log.Printf("reset: password changed for user %d; all sessions revoked", userID)
	http.Redirect(w, r, "/login?reset=1", http.StatusSeeOther)
}

func resetEmailBody(name, link string) string {
	var b strings.Builder
	if name != "" {
		fmt.Fprintf(&b, "Hi %s,\n\n", name)
	} else {
		b.WriteString("Hi,\n\n")
	}
	b.WriteString("Someone asked to reset the password on your GamifyDev account.\n\n")
	b.WriteString(link)
	b.WriteString("\n\nThe link works once and expires in an hour.\n\n")
	b.WriteString("If this wasn't you, ignore this email — nothing has changed, and\n")
	b.WriteString("your current password still works.\n")
	return b.String()
}

// longer returns whichever wait is greater, so a caller checking several
// limiter keys is held by the strictest of them.
func longer(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// humanWait renders a delay the way a person would say it.
func humanWait(d time.Duration) string {
	if d < time.Minute {
		secs := int(d.Seconds()) + 1
		return fmt.Sprintf("%d second%s", secs, plural(secs))
	}
	mins := int(d.Minutes()) + 1
	return fmt.Sprintf("%d minute%s", mins, plural(mins))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
