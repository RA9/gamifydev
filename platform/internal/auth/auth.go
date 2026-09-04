// Package auth handles password hashing, session token generation, and the
// request-context plumbing that carries the current user through handlers.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/RA9/gamifydev/platform/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionCookie = "gd_session"
	SessionTTL    = 30 * 24 * time.Hour
)

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// dummyHash is a bcrypt hash of a value nobody knows, generated once at start.
// It exists only to be compared against.
var dummyHash = func() string {
	secret, err := NewToken()
	if err != nil {
		secret = "gamifydev-fallback-comparison-secret"
	}
	h, _ := HashPassword(secret)
	return h
}()

// BurnPasswordCheck spends the time a real password comparison would, and
// always fails.
//
// Sign-in for an address with no account otherwise returns in microseconds
// while a wrong password takes the ~60ms bcrypt costs. That gap is measurable
// over the network, and it turns the login form into a way to ask which email
// addresses are registered here.
func BurnPasswordCheck(plain string) {
	_ = CheckPassword(dummyHash, plain)
}

// NewToken returns a cryptographically random session token.
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SetSessionCookie writes the session cookie (HttpOnly, SameSite=Lax).
func SetSessionCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(SessionTTL),
	})
}

func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// --- request context --------------------------------------------------------

type ctxKey struct{}

func WithUser(ctx context.Context, u *store.User) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

// CurrentUser returns the authenticated user, or nil for a guest.
func CurrentUser(ctx context.Context) *store.User {
	u, _ := ctx.Value(ctxKey{}).(*store.User)
	return u
}
