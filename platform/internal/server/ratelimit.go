package server

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// attemptLimiter throttles repeated failures against a key — a login email, a
// client address — and is the guard on the credential-facing endpoints.
//
// It counts only failures. A learner who signs in correctly twenty times in a
// morning is never slowed down; someone guessing is, after a handful of misses.
// Success clears the key, so one fat-fingered password costs nothing later.
//
// Unlike runLimiter, whose map is keyed by user id and therefore bounded by the
// number of real accounts, these keys come from whatever a stranger types. The
// map is pruned on write and hard-capped, because an unbounded map fed by an
// anonymous endpoint is itself the denial of service.
type attemptLimiter struct {
	mu      sync.Mutex
	entries map[string]*attemptRecord
	max     int           // failures tolerated inside window
	window  time.Duration // how long failures are remembered
	lockout time.Duration // how long a tripped key stays refused
	cap     int           // hard bound on tracked keys
}

type attemptRecord struct {
	n     int
	first time.Time
	until time.Time // non-zero once locked
}

func newAttemptLimiter(max int, window, lockout time.Duration) *attemptLimiter {
	return &attemptLimiter{
		entries: map[string]*attemptRecord{},
		max:     max,
		window:  window,
		lockout: lockout,
		cap:     20000,
	}
}

// retryAfter reports how long a key must wait, or 0 when it may proceed.
func (l *attemptLimiter) retryAfter(key string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	rec, ok := l.entries[key]
	if !ok {
		return 0
	}
	now := time.Now()
	if now.Before(rec.until) {
		return rec.until.Sub(now)
	}
	if now.Sub(rec.first) > l.window {
		delete(l.entries, key)
	}
	return 0
}

// fail records one failure, locking the key once it passes the limit.
func (l *attemptLimiter) fail(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.pruneLocked(now)

	rec, ok := l.entries[key]
	if !ok || now.Sub(rec.first) > l.window {
		l.entries[key] = &attemptRecord{n: 1, first: now}
		return
	}
	rec.n++
	if rec.n >= l.max {
		// Restart the count with the lock so a key that keeps hammering keeps
		// tripping, rather than trickling one attempt through per window.
		rec.until = now.Add(l.lockout)
		rec.n = 0
		rec.first = now
	}
}

// succeed forgets a key, so a correct password immediately restores a learner
// who mistyped their way close to the limit.
func (l *attemptLimiter) succeed(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}

// pruneLocked drops entries that are past both their window and their lock. It
// runs on write, which is the only path that grows the map; when even that
// leaves it over cap, the map is dropped wholesale rather than allowed to grow
// — forgetting every counter is far better than exhausting memory, and a
// flood large enough to reach here is already visible in the logs.
func (l *attemptLimiter) pruneLocked(now time.Time) {
	if len(l.entries) < l.cap {
		if len(l.entries)%64 != 0 { // amortise: sweep occasionally, not every call
			return
		}
		for k, rec := range l.entries {
			if now.After(rec.until) && now.Sub(rec.first) > l.window {
				delete(l.entries, k)
			}
		}
		return
	}
	l.entries = map[string]*attemptRecord{}
}

// clientIP is the caller's address for rate-limiting purposes.
//
// X-Forwarded-For is honoured only when the app knows it sits behind a proxy
// (the same signal that turns on Secure cookies), because otherwise anyone can
// set the header and mint a fresh identity per request. Even behind a proxy the
// header is a hint rather than proof, which is why the login limiter also keys
// on the email address: spoofing XFF spreads an attack across addresses, but
// every guess against one account still lands on that account's counter.
func (s *Server) clientIP(r *http.Request) string {
	if s.secure {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if first, _, ok := strings.Cut(xff, ","); ok {
				return strings.TrimSpace(first)
			}
			return strings.TrimSpace(xff)
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
