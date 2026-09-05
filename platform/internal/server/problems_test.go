package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// guest creates an anonymous session directly and returns the cookie that
// stands for it, so these tests can set up work that was done before signing in
// without going through a judge the test server has no sandbox for.
func (ts *testServer) guest(t *testing.T) (*http.Cookie, store.Solver) {
	t.Helper()
	token, err := auth.NewToken()
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	id, err := ts.st.CreateGuestSession(t.Context(), token)
	if err != nil {
		t.Fatalf("create guest session: %v", err)
	}
	return &http.Cookie{Name: auth.GuestCookie, Value: token}, store.Solver{GuestID: id}
}

// solve records an accepted submission for a solver against a seeded problem.
func (ts *testServer) solve(t *testing.T, slug string, by store.Solver) int64 {
	t.Helper()
	p, err := ts.st.GetProblemBySlug(t.Context(), slug)
	if err != nil {
		t.Fatalf("get problem %s: %v", slug, err)
	}
	id, err := ts.st.CreateProblemSubmission(t.Context(), p.ID, by, "python", "print('hi')")
	if err != nil {
		t.Fatalf("create submission: %v", err)
	}
	if err := ts.st.RecordVerdict(t.Context(), id, store.Verdict{
		Name: store.VerdictAccepted, Passed: 3, Total: 3,
	}); err != nil {
		t.Fatalf("record verdict: %v", err)
	}
	return id
}

func TestTheProblemBankIsReadableWithoutAnAccount(t *testing.T) {
	ts := newTestServer(t)
	for _, path := range []string{"/problems", "/problems/affordable-pairs"} {
		w := ts.do(t, http.MethodGet, path, nil)
		if w.Code != http.StatusOK {
			t.Errorf("%s: status %d, want 200 — the bank exists to be tried without signing up", path, w.Code)
		}
	}
}

// Reading the bank must not mint anything. /problems is a public page, and a
// database row per crawler hit is a table that grows with traffic rather than
// with people.
func TestBrowsingTheBankDoesNotCreateAGuestSession(t *testing.T) {
	ts := newTestServer(t)
	w := ts.do(t, http.MethodGet, "/problems", nil)
	for _, c := range w.Result().Cookies() {
		if c.Name == auth.GuestCookie && c.Value != "" {
			t.Fatalf("a plain page view set a guest cookie")
		}
	}
}

// A submission holds someone's code. Knowing its id is not permission to read
// it — this is what lets the route be public without leaking anyone's work.
func TestASubmissionIsNotReadableByItsID(t *testing.T) {
	ts := newTestServer(t)
	mine, myID := ts.guest(t)
	theirs, _ := ts.guest(t)
	subID := ts.solve(t, "affordable-pairs", myID)
	path := "/problems/result/" + strconv.FormatInt(subID, 10)

	if w := ts.do(t, http.MethodGet, path, nil, mine); w.Code != http.StatusOK {
		t.Errorf("the solver who wrote it got %d, want 200", w.Code)
	}
	if w := ts.do(t, http.MethodGet, path, nil, theirs); w.Code != http.StatusNotFound {
		t.Errorf("another guest got %d, want 404", w.Code)
	}
	if w := ts.do(t, http.MethodGet, path, nil); w.Code != http.StatusNotFound {
		t.Errorf("a visitor with no cookie at all got %d, want 404", w.Code)
	}
	learner := ts.register(t, "nosy@example.com")
	if w := ts.do(t, http.MethodGet, path, nil, learner); w.Code != http.StatusNotFound {
		t.Errorf("a signed-in stranger got %d, want 404", w.Code)
	}
}

// The promise the bank makes to a guest: solve things now, keep them if you
// sign up. Losing that work at the moment someone commits to an account would
// be the worst possible time to lose it.
func TestSigningUpKeepsWhatYouSolvedAsAGuest(t *testing.T) {
	ts := newTestServer(t)
	ts.register(t, "admin@example.com")
	cookie, by := ts.guest(t)
	ts.solve(t, "affordable-pairs", by)
	ts.solve(t, "double-tiles", by)
	ts.qualifyGuest(t, by, true)

	w := ts.do(t, http.MethodPost, "/register", url.Values{
		"name": {"Guest Turned Learner"}, "email": {"keeper@example.com"},
		"password": {"averysafepassword"},
	}, cookie)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("register: status %d, want 303", w.Code)
	}

	u, err := ts.st.GetUserByEmail(t.Context(), "keeper@example.com")
	if err != nil {
		t.Fatalf("get new user: %v", err)
	}
	n, err := ts.st.SolvedCount(t.Context(), store.Solver{UserID: u.ID})
	if err != nil {
		t.Fatalf("solved count: %v", err)
	}
	if n != 2 {
		t.Errorf("the new account has %d solved problems, want 2", n)
	}
	// And the guest cookie is spent, so it can't hand the same work to a second
	// account later.
	cleared := false
	for _, c := range w.Result().Cookies() {
		if c.Name == auth.GuestCookie && c.MaxAge < 0 {
			cleared = true
		}
	}
	if !cleared {
		t.Errorf("the guest cookie was not cleared after its work was claimed")
	}
	if _, err := ts.st.GuestByToken(t.Context(), cookie.Value); err == nil {
		t.Errorf("a claimed guest session still resolves, so the cookie could be replayed")
	}
}

// Someone who practised as a guest and then signed into an account they already
// had should keep their work just as much as someone who registered.
func TestSigningInAlsoClaimsGuestWork(t *testing.T) {
	ts := newTestServer(t)
	ts.register(t, "admin@example.com")
	ts.register(t, "returning@example.com")

	cookie, by := ts.guest(t)
	ts.solve(t, "affordable-pairs", by)

	w := ts.do(t, http.MethodPost, "/login", url.Values{
		"email": {"returning@example.com"}, "password": {"averysafepassword"},
	}, cookie)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("login: status %d, want 303", w.Code)
	}
	u, err := ts.st.GetUserByEmail(t.Context(), "returning@example.com")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if n, _ := ts.st.SolvedCount(t.Context(), store.Solver{UserID: u.ID}); n != 1 {
		t.Errorf("the returning account has %d solved problems, want 1", n)
	}
}

// The list marks what you've solved, and the sign-up nudge appears only once
// there is something to lose by ignoring it.
func TestTheListShowsAGuestWhatTheyWouldKeep(t *testing.T) {
	ts := newTestServer(t)
	cookie, by := ts.guest(t)

	// Asserted through the sign-up link and the progress bar's own value rather
	// than through wording, so a copy edit doesn't fail the test but dropping
	// the invitation — or miscounting — still does.
	before := ts.do(t, http.MethodGet, "/problems", nil, cookie).Body.String()
	if strings.Contains(before, `href="/register"`) {
		t.Errorf("a guest who has solved nothing is being asked to save nothing")
	}

	ts.solve(t, "affordable-pairs", by)
	after := ts.do(t, http.MethodGet, "/problems", nil, cookie).Body.String()
	if !strings.Contains(after, `href="/register"`) {
		t.Errorf("a guest with work to lose is not being offered a way to keep it")
	}
	if !strings.Contains(after, `aria-valuenow="1"`) {
		t.Errorf("the solved count does not reflect the guest's accepted submission")
	}
}

// The picker only offers languages the problem has starters for, but a POST body
// can say anything. A submission judged against tests written for another
// language would fail for reasons that have nothing to do with the learner.
func TestASubmissionCannotNameALanguageTheProblemDoesNotAccept(t *testing.T) {
	ts := newTestServer(t)
	p, err := ts.st.GetProblemBySlug(t.Context(), "affordable-pairs")
	if err != nil {
		t.Fatalf("get problem: %v", err)
	}
	if _, err := ts.acceptedLanguage(t.Context(), p.ID, "c"); err == nil {
		t.Errorf("a language with no starter was accepted")
	}
	if got, err := ts.acceptedLanguage(t.Context(), p.ID, "python"); err != nil || got != "python" {
		t.Errorf("acceptedLanguage(python) = %q, %v; want python, nil", got, err)
	}
}
