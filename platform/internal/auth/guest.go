package auth

import (
	"net/http"
	"time"
)

const (
	// GuestCookie identifies an anonymous visitor working through the practice
	// problems. Separate from the session cookie on purpose: it grants nothing
	// and outlives sign-in, so signing out must not throw away work the visitor
	// did before they had an account.
	GuestCookie = "gd_guest"
	// GuestTTL is long because the promise the problem bank makes is that you
	// can try it, close the tab, and come back to your solved set. A cookie
	// that expired in a week would quietly break that for exactly the casual
	// visitor it exists to serve.
	GuestTTL = 180 * 24 * time.Hour
)

// SetGuestCookie writes the anonymous-visitor cookie (HttpOnly, SameSite=Lax).
func SetGuestCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     GuestCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(GuestTTL),
	})
}

// ClearGuestCookie spends the cookie once its work has been handed to an
// account. The token is a bearer credential for that work, so leaving it in the
// browser after the handover would be leaving a key to a door that moved.
func ClearGuestCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     GuestCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
