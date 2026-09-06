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
func (ts *testServer) sitting(t *testing.T) *http.Cookie {
	t.Helper()
	w := ts.do(t, http.MethodPost, "/placement/start", url.Values{})
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
