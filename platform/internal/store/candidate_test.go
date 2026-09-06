package store

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestNormalizeEmailCollapsesTheEasyDisguises(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"Alice@Example.com", "alice@example.com"},
		{"  alice@example.com  ", "alice@example.com"},
		{"alice+placement@example.com", "alice@example.com"},
		{"alice+one+two@example.com", "alice@example.com"},
		// A dot is not a tag. Gmail treats first.last and firstlast as one
		// account; most providers do not, and merging them would fuse two real
		// people into one set of attempts.
		{"first.last@example.com", "first.last@example.com"},
		// A leading + is the whole local part, not a tag on an empty one.
		{"+weird@example.com", "+weird@example.com"},
		{"not-an-email", "not-an-email"},
	} {
		if got := NormalizeEmail(tc.in); got != tc.want {
			t.Errorf("NormalizeEmail(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestACandidateIsOnePersonHoweverTheyTypeTheirAddress(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	first, err := st.UpsertCandidate(ctx, "Ada Lovelace", "Ada@Example.com")
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	again, err := st.UpsertCandidate(ctx, "Ada L", "ada+retry@example.com")
	if err != nil {
		t.Fatalf("upsert again: %v", err)
	}
	if first.ID != again.ID {
		t.Errorf("the same person got two identities (%d and %d)", first.ID, again.ID)
	}
	if again.Name != "Ada L" {
		t.Errorf("name = %q, want the newer spelling", again.Name)
	}
	if again.RawEmail != "ada+retry@example.com" {
		t.Errorf("raw email = %q, want the address as typed", again.RawEmail)
	}
}

func TestUnusableDetailsAreRefused(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	for _, tc := range []struct{ name, email string }{
		{"", "a@b.com"},
		{"A", "a@b.com"}, // a single letter is not a name
		{"Ada", "not-an-email"},
		{"Ada", "a@b"},       // no dot in the domain
		{"Ada", "a b@c.com"}, // whitespace
		{"Ada", ""},
	} {
		if _, err := st.UpsertCandidate(ctx, tc.name, tc.email); !errors.Is(err, ErrInvalidCandidate) {
			t.Errorf("UpsertCandidate(%q, %q) = %v, want ErrInvalidCandidate", tc.name, tc.email, err)
		}
	}
}

// The reason candidates exist. Sitting the test, failing, then opening a
// private window used to hand you a clean slate; the cooldown was fitted to the
// one thing the candidate fully controls.
func TestAFreshBrowserDoesNotResetTheCooldown(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	cand, err := st.UpsertCandidate(ctx, "Ada Lovelace", "ada@example.com")
	if err != nil {
		t.Fatalf("candidate: %v", err)
	}

	// First sitting, in one browser, failed.
	browserA := mkGuest(t, st, freshToken(t))
	att, err := st.StartPlacementFor(ctx, browserA, cand.ID, a)
	if err != nil {
		t.Fatalf("first sitting: %v", err)
	}
	if _, err := st.ScoreAttemptFor(ctx, browserA, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("score: %v", err)
	}
	if _, err := st.SavePlacementResultFor(ctx, browserA, PlacementResult{AttemptID: att.ID, Passed: false}); err != nil {
		t.Fatalf("save result: %v", err)
	}

	// A completely different browser — new cookie, same person.
	browserB := mkGuest(t, st, freshToken(t))
	_, err = st.StartPlacementFor(ctx, browserB, cand.ID, a)
	if !errors.Is(err, ErrPlacementRetryTooSoon) {
		t.Errorf("second browser = %v, want ErrPlacementRetryTooSoon — the cooldown reset", err)
	}
	// And without the identity, it would have reset: this is what the gate buys.
	if _, err := st.StartPlacementFor(ctx, mkGuest(t, st, freshToken(t)), 0, a); err != nil {
		t.Fatalf("an unidentified sitting should still be allowed: %v", err)
	}
}

// A pass follows the person too, so a second identity-less browser cannot be
// used to sit it again for a better score.
func TestAPassFollowsTheCandidateAcrossBrowsers(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	cand, _ := st.UpsertCandidate(ctx, "Grace Hopper", "grace@example.com")

	browserA := mkGuest(t, st, freshToken(t))
	att, err := st.StartPlacementFor(ctx, browserA, cand.ID, a)
	if err != nil {
		t.Fatalf("sitting: %v", err)
	}
	if _, err := st.ScoreAttemptFor(ctx, browserA, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("score: %v", err)
	}
	if _, err := st.SavePlacementResultFor(ctx, browserA, PlacementResult{AttemptID: att.ID, Passed: true}); err != nil {
		t.Fatalf("save result: %v", err)
	}

	browserB := mkGuest(t, st, freshToken(t))
	if _, err := st.StartPlacementFor(ctx, browserB, cand.ID, a); !errors.Is(err, ErrPlacementAlreadyPassed) {
		t.Errorf("second browser after passing = %v, want ErrPlacementAlreadyPassed", err)
	}
}

// The rolling limit counts the person's sittings, not the browser's.
func TestTheAttemptLimitCountsThePersonNotTheBrowser(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	cand, _ := st.UpsertCandidate(ctx, "Alan Turing", "alan@example.com")

	// Three sittings, each abandoned, each in its own fresh browser.
	for i := 0; i < PlacementMaxAttempts; i++ {
		by := mkGuest(t, st, freshToken(t))
		att, err := st.StartPlacementFor(ctx, by, cand.ID, a)
		if err != nil {
			t.Fatalf("sitting %d: %v", i, err)
		}
		if _, err := st.db.ExecContext(ctx,
			`UPDATE assessment_attempts SET abandoned = 1, submitted_at = datetime('now') WHERE id = ?`,
			att.ID); err != nil {
			t.Fatalf("abandon: %v", err)
		}
	}
	if _, err := st.StartPlacementFor(ctx, mkGuest(t, st, freshToken(t)), cand.ID, a); !errors.Is(err, ErrPlacementAttemptLimit) {
		t.Errorf("a fourth browser = %v, want ErrPlacementAttemptLimit", err)
	}
}

// A sitting the proctor closed is a spent attempt for the person, not for the
// cookie — otherwise being caught cheating costs nothing but a new tab.
func TestBeingCaughtCheatingFollowsTheCandidate(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	cand, _ := st.UpsertCandidate(ctx, "Someone Else", "else@example.com")

	browserA := mkGuest(t, st, freshToken(t))
	att, err := st.StartPlacementFor(ctx, browserA, cand.ID, a)
	if err != nil {
		t.Fatalf("sitting: %v", err)
	}
	for i := 0; i < ProctorStrikeLimit; i++ {
		if _, err := st.RecordViolation(ctx, browserA, att.ID, "copy"); err != nil {
			t.Fatalf("violation: %v", err)
		}
	}
	if _, err := st.SavePlacementResultFor(ctx, browserA, PlacementResult{AttemptID: att.ID, Passed: false}); err != nil {
		t.Fatalf("save result: %v", err)
	}

	browserB := mkGuest(t, st, freshToken(t))
	if _, err := st.StartPlacementFor(ctx, browserB, cand.ID, a); !errors.Is(err, ErrPlacementRetryTooSoon) {
		t.Errorf("a fresh browser after being caught = %v, want the cooldown to hold", err)
	}
}

// Once an identity holds an account, an enrollment locks placement from any
// browser — signing out must not reopen the decision a schedule was built from.
func TestAnEnrollmentLocksPlacementForTheIdentityNotJustTheAccount(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	cand, _ := st.UpsertCandidate(ctx, "Enrolled Person", "enrolled@example.com")
	uid := mkUser(t, st, "enrolled@example.com")
	if err := st.LinkCandidateToUser(ctx, cand.ID, uid); err != nil {
		t.Fatalf("link: %v", err)
	}
	if _, err := st.db.ExecContext(ctx,
		`INSERT INTO enrollments (user_id, state) VALUES (?, 'active')`, uid); err != nil {
		t.Fatalf("enroll: %v", err)
	}

	// Signed out, in a browser that knows nothing about the account.
	guest := mkGuest(t, st, freshToken(t))
	if _, err := st.StartPlacementFor(ctx, guest, cand.ID, a); !errors.Is(err, ErrPlacementEnrollmentLocked) {
		t.Errorf("signed-out sitting = %v, want ErrPlacementEnrollmentLocked", err)
	}
}

// The link is one-way. Letting a candidate be repointed at a second account is
// how a spent identity would be laundered into a clean one.
func TestACandidateCannotBeMovedToAnotherAccount(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	cand, _ := st.UpsertCandidate(ctx, "Once Only", "once@example.com")
	first := mkUser(t, st, "first@example.com")
	second := mkUser(t, st, "second@example.com")

	if err := st.LinkCandidateToUser(ctx, cand.ID, first); err != nil {
		t.Fatalf("first link: %v", err)
	}
	if err := st.LinkCandidateToUser(ctx, cand.ID, second); err != nil {
		t.Fatalf("second link: %v", err)
	}
	got, err := st.CandidateByEmail(ctx, "once@example.com")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if got.UserID != first {
		t.Errorf("candidate now belongs to %d, want %d — the link was moved", got.UserID, first)
	}
}

// freshToken is a unique guest-cookie value, so each "browser" in these tests
// really is a separate one. Uniqueness is all that matters here, not secrecy.
var browserSeq atomic.Int64

func freshToken(t *testing.T) string {
	t.Helper()
	return "browser-" + strconv.FormatInt(browserSeq.Add(1), 10)
}
