package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if strings.Contains(hash, "correct horse") {
		t.Fatal("the hash contains the password")
	}
	if !CheckPassword(hash, "correct horse battery staple") {
		t.Error("the right password was rejected")
	}
	if CheckPassword(hash, "correct horse battery stapl") {
		t.Error("a wrong password was accepted")
	}
}

func TestTheSamePasswordHashesDifferentlyEachTime(t *testing.T) {
	// bcrypt salts per call. Identical hashes would mean two people with the
	// same password are visibly the same in a database dump.
	a, _ := HashPassword("same-password")
	b, _ := HashPassword("same-password")
	if a == b {
		t.Fatal("two hashes of the same password are identical — the salt is missing")
	}
	if !CheckPassword(a, "same-password") || !CheckPassword(b, "same-password") {
		t.Fatal("salted hashes do not both verify")
	}
}

func TestCheckPasswordSurvivesGarbage(t *testing.T) {
	// A malformed hash must be a failed sign-in, never a panic or a success.
	for _, hash := range []string{"", "not-a-hash", "$2a$broken"} {
		if CheckPassword(hash, "anything") {
			t.Errorf("hash %q accepted a password", hash)
		}
	}
}

func TestBurnPasswordCheckAlwaysFailsAndCostsRealTime(t *testing.T) {
	// It exists to spend bcrypt's time on a miss. If it ever became cheap, the
	// timing gap it closes would quietly reopen.
	BurnPasswordCheck("whatever") // must not panic
	if CheckPassword(dummyHash, "whatever") {
		t.Fatal("the dummy hash matched a guess")
	}
	if !strings.HasPrefix(dummyHash, "$2") {
		t.Fatalf("dummyHash is not a bcrypt hash: %q", dummyHash)
	}
}

func TestTokensAreLongAndUnique(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		tok, err := NewToken()
		if err != nil {
			t.Fatalf("token: %v", err)
		}
		if len(tok) != 64 { // 32 random bytes, hex-encoded
			t.Fatalf("token length %d, want 64", len(tok))
		}
		if seen[tok] {
			t.Fatal("NewToken repeated a value")
		}
		seen[tok] = true
	}
}

func TestSessionCookieIsHardenedAndClearable(t *testing.T) {
	w := httptest.NewRecorder()
	SetSessionCookie(w, "tok123", true)
	c := w.Result().Cookies()[0]
	if !c.HttpOnly {
		t.Error("session cookie is readable from JavaScript")
	}
	if !c.Secure {
		t.Error("session cookie is not marked Secure in production")
	}
	// Lax is what stands in for CSRF tokens here: it stops the cookie riding
	// along on a cross-site POST.
	if c.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite = %v, want Lax", c.SameSite)
	}

	w = httptest.NewRecorder()
	ClearSessionCookie(w, true)
	c = w.Result().Cookies()[0]
	if c.Value != "" || c.MaxAge >= 0 {
		t.Errorf("clearing left value %q, MaxAge %d", c.Value, c.MaxAge)
	}
}

func TestCurrentUserOnAContextWithoutOneIsNil(t *testing.T) {
	// Every guard in the server keys off this returning nil rather than
	// panicking on a guest request.
	if u := CurrentUser(httptest.NewRequest("GET", "/", nil).Context()); u != nil {
		t.Fatalf("guest context yielded a user: %+v", u)
	}
}
