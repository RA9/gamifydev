package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// sitting starts a real placement paper as a guest and returns the cookie that
// identifies them.
//
// Each gets its own email, because the gate is now keyed to the candidate: two
// tests sharing an address would share a cooldown and the second would be
// refused for reasons that have nothing to do with what it is testing.
func (ts *testServer) sitting(t *testing.T) *http.Cookie {
	t.Helper()
	w := ts.do(t, http.MethodPost, "/placement/start", url.Values{
		"name":  {"Test Candidate"},
		"email": {strings.ToLower(strings.NewReplacer("/", "-", " ", "-").Replace(t.Name())) + "@example.com"},
	})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("start placement: status %d, want 303", w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == "gd_guest" && c.Value != "" {
			return c
		}
	}
	t.Fatal("starting a paper did not identify the candidate")
	return nil
}

func (ts *testServer) reportViolation(t *testing.T, kind string, cookie *http.Cookie) (int, map[string]any) {
	t.Helper()
	var w = ts.do(t, http.MethodPost, "/placement/violation", url.Values{"kind": {kind}}, cookie)
	out := map[string]any{}
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode violation response: %v (%s)", err, w.Body.String())
		}
	}
	return w.Code, out
}

// The rule end to end, through the real endpoint: warned once, closed twice.
func TestTheSecondViolationClosesTheSittingOverHTTP(t *testing.T) {
	ts := newTestServer(t)
	cookie := ts.sitting(t)

	code, first := ts.reportViolation(t, "copy", cookie)
	if code != http.StatusOK {
		t.Fatalf("first report: status %d", code)
	}
	if first["strikes"].(float64) != 1 || first["voided"].(bool) {
		t.Fatalf("first report = %v, want one strike and no void", first)
	}

	code, second := ts.reportViolation(t, "hidden", cookie)
	if code != http.StatusOK {
		t.Fatalf("second report: status %d", code)
	}
	if !second["voided"].(bool) {
		t.Fatalf("second report = %v, want the sitting voided", second)
	}
}

// A void has to leave a failed result behind, or the candidate lands on a page
// with nothing to read and the retry rules never see the failure.
func TestAVoidedSittingLeavesAFailedResult(t *testing.T) {
	ts := newTestServer(t)
	cookie := ts.sitting(t)
	for _, kind := range []string{"copy", "capture"} {
		if code, _ := ts.reportViolation(t, kind, cookie); code != http.StatusOK {
			t.Fatalf("report %s: status %d", kind, code)
		}
	}

	w := ts.do(t, http.MethodGet, "/placement/result", nil, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("result: status %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "flagged for cheating") {
		t.Error("the result does not tell the candidate why their sitting was closed")
	}

	by := store.Solver{}
	if id, err := ts.st.GuestByToken(t.Context(), cookie.Value); err == nil {
		by = store.Solver{GuestID: id}
	}
	res, err := ts.st.LatestPlacementFor(t.Context(), by)
	if err != nil || res == nil {
		t.Fatalf("no placement result was recorded: %v", err)
	}
	if res.Passed {
		t.Error("a voided sitting was recorded as a pass")
	}
}

// Handing the paper in after it was closed must not score it.
func TestSubmittingAVoidedPaperIsRefused(t *testing.T) {
	ts := newTestServer(t)
	cookie := ts.sitting(t)
	for _, kind := range []string{"copy", "blur"} {
		ts.reportViolation(t, kind, cookie)
	}
	w := ts.do(t, http.MethodPost, "/placement/submit", url.Values{}, cookie)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement/result" {
		t.Errorf("submitting a voided paper = %d -> %q, want a redirect to the result",
			w.Code, w.Header().Get("Location"))
	}
}

// Nobody without a live sitting can spend strikes — least of all on somebody
// else's paper. This is the endpoint's whole attack surface.
func TestOnlyACandidateMidSittingCanReportAViolation(t *testing.T) {
	ts := newTestServer(t)
	sitter := ts.sitting(t)

	if code, _ := ts.reportViolation(t, "copy", &http.Cookie{Name: "gd_guest", Value: "not-a-session"}); code != http.StatusForbidden {
		t.Errorf("a stranger reporting = %d, want 403", code)
	}
	w := ts.do(t, http.MethodPost, "/placement/violation", url.Values{"kind": {"copy"}})
	if w.Code != http.StatusForbidden {
		t.Errorf("a visitor with no sitting = %d, want 403", w.Code)
	}
	if code, _ := ts.reportViolation(t, "sneezed", sitter); code != http.StatusBadRequest {
		t.Errorf("an invented violation kind = %d, want 400", code)
	}
	// And the invented kind cost the candidate nothing.
	if code, out := ts.reportViolation(t, "copy", sitter); code != http.StatusOK || out["strikes"].(float64) != 1 {
		t.Errorf("after a rejected kind the first real strike = %v, want 1", out)
	}
}

// The candidate has to be told the rules before the clock starts, and the paper
// has to carry the proctor. A rule nobody was shown is a trap.
func TestTheRulesAreStatedBeforeTheClockStarts(t *testing.T) {
	ts := newTestServer(t)
	intro := ts.do(t, http.MethodGet, "/placement", nil).Body.String()
	for _, want := range []string{"supervised", "warning", "closes"} {
		if !strings.Contains(strings.ToLower(intro), want) {
			t.Errorf("the introduction does not mention %q", want)
		}
	}

	cookie := ts.sitting(t)
	paper := ts.do(t, http.MethodGet, "/placement", nil, cookie).Body.String()
	if !strings.Contains(paper, "data-proctored") {
		t.Error("the paper is not marked as proctored, so the invigilator never starts")
	}
	if !strings.Contains(paper, "proctor.js") {
		t.Error("the paper does not load the proctor")
	}
}

// The bypass this whole identity layer exists to close: sit the test, get
// caught, then open a private window. Two different cookie jars, one email.
func TestAnotherBrowserDoesNotGetAFreshPlacement(t *testing.T) {
	ts := newTestServer(t)
	details := url.Values{"name": {"Ada Lovelace"}, "email": {"ada@example.com"}}

	// First browser: start, get caught twice, sitting closed and failed.
	w := ts.do(t, http.MethodPost, "/placement/start", details)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("first start = %d, want 303", w.Code)
	}
	var first *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "gd_guest" {
			first = c
		}
	}
	for _, kind := range []string{"copy", "hidden"} {
		ts.reportViolation(t, kind, first)
	}

	// Second browser: no cookie at all, same person.
	w = ts.do(t, http.MethodPost, "/placement/start", details)
	if w.Code == http.StatusSeeOther {
		t.Fatal("a fresh browser was handed a new sitting — the gate is still cookie-deep")
	}
	if !strings.Contains(w.Body.String(), "seven-day wait") {
		t.Errorf("the refusal does not explain itself:\n%s", w.Body.String()[:min(600, w.Body.Len())])
	}

	// And plus-addressing is not a second identity either.
	tagged := url.Values{"name": {"Ada Lovelace"}, "email": {"Ada+again@Example.com"}}
	if w := ts.do(t, http.MethodPost, "/placement/start", tagged); w.Code == http.StatusSeeOther {
		t.Error("a +tag on the same address opened a fresh sitting")
	}
}

// Details are required, and the page says why rather than silently doing
// nothing.
func TestAPaperCannotStartWithoutIdentifyingTheCandidate(t *testing.T) {
	ts := newTestServer(t)
	for _, bad := range []url.Values{
		{},
		{"name": {"Ada"}},
		{"email": {"ada@example.com"}},
		{"name": {"A"}, "email": {"ada@example.com"}},
		{"name": {"Ada Lovelace"}, "email": {"not-an-email"}},
	} {
		w := ts.do(t, http.MethodPost, "/placement/start", bad)
		if w.Code == http.StatusSeeOther {
			t.Errorf("%v started a sitting", bad)
			continue
		}
		if !strings.Contains(w.Body.String(), "full name and a real email") {
			t.Errorf("%v was refused without saying why", bad)
		}
	}
}

// A signed-in learner is already identified, so they are never asked again.
func TestASignedInLearnerIsNotAskedWhoTheyAre(t *testing.T) {
	ts := newTestServer(t)
	admin := ts.register(t, "admin@example.com")
	body := ts.do(t, http.MethodGet, "/placement", nil, admin).Body.String()
	if strings.Contains(body, `name="email"`) {
		t.Error("a signed-in user is being asked for details the account already has")
	}
	if w := ts.do(t, http.MethodPost, "/placement/start", url.Values{}, admin); w.Code != http.StatusSeeOther {
		t.Errorf("a signed-in start = %d, want 303", w.Code)
	}
}
