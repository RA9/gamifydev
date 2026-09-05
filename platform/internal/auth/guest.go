package auth

import (
	"net/http"
	"time"
)

const (
	// GuestCookie identifies an anonymous visitor taking placement or practising
	// problems. It is separate from the account session and acts as a bearer key
	// for that guest-owned work, including a passing placement result.
	GuestCookie = "gd_guest"
	// GuestTTL allows someone to return to unfinished guest work. Passing results
	// remain single-use because claiming the account spends the server-side
	// guest session even when a browser still holds the old cookie.
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
