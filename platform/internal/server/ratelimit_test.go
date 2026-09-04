package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiterAllowsUntilTheLimitThenRefuses(t *testing.T) {
	l := newAttemptLimiter(3, time.Minute, time.Minute)
	for i := 0; i < 2; i++ {
		l.fail("k")
		if w := l.retryAfter("k"); w != 0 {
			t.Fatalf("refused after %d failure(s), want %d allowed", i+1, 3)
		}
	}
	l.fail("k") // the third trips it
	if l.retryAfter("k") == 0 {
		t.Fatal("still allowed after reaching the limit")
	}
}

func TestSuccessForgetsEarlierMisses(t *testing.T) {
	// Someone who mistypes twice in the morning must not be two misses closer
	// to a lockout in the afternoon.
	l := newAttemptLimiter(3, time.Minute, time.Minute)
	l.fail("k")
	l.fail("k")
	l.succeed("k")
	l.fail("k")
	l.fail("k")
	if w := l.retryAfter("k"); w != 0 {
		t.Fatalf("locked out after 2 failures following a success (wait %v)", w)
	}
}

func TestFailuresOutsideTheWindowAreForgotten(t *testing.T) {
	l := newAttemptLimiter(3, 30*time.Millisecond, time.Minute)
	l.fail("k")
	l.fail("k")
	time.Sleep(50 * time.Millisecond)
	l.fail("k")
	l.fail("k")
	if w := l.retryAfter("k"); w != 0 {
		t.Fatalf("stale failures counted toward the limit (wait %v)", w)
	}
}

func TestLockoutExpires(t *testing.T) {
	l := newAttemptLimiter(2, time.Minute, 40*time.Millisecond)
	l.fail("k")
	l.fail("k")
	if l.retryAfter("k") == 0 {
		t.Fatal("not locked after reaching the limit")
	}
	time.Sleep(60 * time.Millisecond)
	if w := l.retryAfter("k"); w != 0 {
		t.Fatalf("still locked past the cool-off (wait %v)", w)
	}
}

func TestKeysAreIndependent(t *testing.T) {
	// One account being attacked must not lock everyone else out.
	l := newAttemptLimiter(2, time.Minute, time.Minute)
	l.fail("victim")
	l.fail("victim")
	if l.retryAfter("victim") == 0 {
		t.Fatal("victim not locked")
	}
	if w := l.retryAfter("bystander"); w != 0 {
		t.Fatalf("an unrelated key was locked too (wait %v)", w)
	}
}

func TestTrackedKeysStayBounded(t *testing.T) {
	// The keys are whatever a stranger types, so the map must not be a way to
	// exhaust memory.
	l := newAttemptLimiter(5, time.Minute, time.Minute)
	l.cap = 500
	for i := 0; i < 5000; i++ {
		l.fail(string(rune(i%1114111)) + "-key")
	}
	l.mu.Lock()
	n := len(l.entries)
	l.mu.Unlock()
	if n > l.cap {
		t.Fatalf("tracking %d keys, cap is %d", n, l.cap)
	}
}

func TestClientIPIgnoresForwardedHeaderWhenNotBehindAProxy(t *testing.T) {
	// Trusting X-Forwarded-For unconditionally would let anyone mint a fresh
	// identity per request and walk straight past the per-address limit.
	s := &Server{secure: false}
	r := httptest.NewRequest(http.MethodPost, "/login", nil)
	r.RemoteAddr = "10.1.2.3:5555"
	r.Header.Set("X-Forwarded-For", "1.2.3.4")
	if got := s.clientIP(r); got != "10.1.2.3" {
		t.Fatalf("clientIP = %q, want the real peer 10.1.2.3", got)
	}
}

func TestClientIPUsesForwardedHeaderBehindAProxy(t *testing.T) {
	s := &Server{secure: true}
	r := httptest.NewRequest(http.MethodPost, "/login", nil)
	r.RemoteAddr = "10.1.2.3:5555"
	r.Header.Set("X-Forwarded-For", "1.2.3.4, 10.0.0.1")
	if got := s.clientIP(r); got != "1.2.3.4" {
		t.Fatalf("clientIP = %q, want the originating client 1.2.3.4", got)
	}
}

func TestHumanWaitReadsLikeAPerson(t *testing.T) {
	for _, tc := range []struct {
		in   time.Duration
		want string
	}{
		{500 * time.Millisecond, "1 second"},
		{30 * time.Second, "31 seconds"},
		{90 * time.Second, "2 minutes"},
	} {
		if got := humanWait(tc.in); got != tc.want {
			t.Errorf("humanWait(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
