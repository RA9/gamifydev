package server

import (
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/runner"
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
	// Every problem in the bank accepts C, Go and Python, so the language that
	// must be refused is one nobody wrote a starter for.
	for _, lang := range []string{"shell", "ruby", "", "javascript"} {
		if _, err := ts.acceptedLanguage(t.Context(), p.ID, lang); err == nil {
			t.Errorf("%q has no starter but was accepted", lang)
		}
	}
	for _, lang := range []string{"c", "go", "python"} {
		if got, err := ts.acceptedLanguage(t.Context(), p.ID, lang); err != nil || got != lang {
			t.Errorf("acceptedLanguage(%q) = %q, %v; want it accepted", lang, got, err)
		}
	}
}

// Run and Submit are different actions, and the difference is the point: Run
// checks the cases the learner can read and costs them nothing; Submit judges
// the hidden tests too and goes on the record. If Run started recording, every
// experiment would show up as a failed attempt.
func TestRunningTheExampleCasesRecordsNothing(t *testing.T) {
	ts := newTestServerWith(t, pythonExecutor(t))
	cookie, by := ts.guest(t)
	form := url.Values{"language": {"python"}, "code": {"def total(nums):\n    return 0\n"}}

	w := ts.do(t, http.MethodPost, "/problems/double-tiles/run", form, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("run: status %d", w.Code)
	}
	p, err := ts.st.GetProblemBySlug(t.Context(), "double-tiles")
	if err != nil {
		t.Fatalf("get problem: %v", err)
	}
	subs, err := ts.st.ProblemSubmissionsFor(t.Context(), p.ID, by, 10)
	if err != nil {
		t.Fatalf("submissions: %v", err)
	}
	if len(subs) != 0 {
		t.Errorf("Run recorded %d submission(s); it should leave no trace", len(subs))
	}
}

// Whatever Run reports, it must never be the hidden tests — those are the
// actual bar, and showing them would hand over the answer.
func TestRunOnlyReportsTheVisibleCases(t *testing.T) {
	ts := newTestServerWith(t, pythonExecutor(t))
	cookie, _ := ts.guest(t)
	p, _ := ts.st.GetProblemBySlug(t.Context(), "double-tiles")
	all, _ := ts.st.ProblemTests(t.Context(), p.ID)
	visible := 0
	var hiddenLabels []string
	for _, c := range all {
		if c.Hidden {
			hiddenLabels = append(hiddenLabels, c.Label)
		} else {
			visible++
		}
	}
	if visible == 0 || len(hiddenLabels) == 0 {
		t.Skip("this problem has no mix of visible and hidden tests to distinguish")
	}

	body := ts.do(t, http.MethodPost, "/problems/double-tiles/run", url.Values{
		"language": {"python"}, "code": {"def score(values, gold):\n    return 0\n"},
	}, cookie).Body.String()
	if got := strings.Count(body, `class="case`); got < visible {
		t.Errorf("run listed %d case rows, want at least the %d visible ones", got, visible)
	}
	for _, label := range hiddenLabels {
		if strings.Contains(body, label) {
			t.Errorf("run leaked the hidden test %q", label)
		}
	}
}

// The left pane's tab comes from the query string, and an unknown one must land
// somewhere real rather than on an empty pane.
func TestTheWorkspaceTabsAreBothReachable(t *testing.T) {
	ts := newTestServer(t)
	for _, tc := range []struct{ query, want string }{
		{"", `Example cases`},
		{"?tab=problem", `Example cases`},
		{"?tab=submissions", `Your submissions`},
		{"?tab=nonsense", `Example cases`},
	} {
		w := ts.do(t, http.MethodGet, "/problems/double-tiles"+tc.query, nil)
		if w.Code != http.StatusOK {
			t.Errorf("%q: status %d", tc.query, w.Code)
			continue
		}
		if !strings.Contains(w.Body.String(), tc.want) {
			t.Errorf("%q does not show %q", tc.query, tc.want)
		}
	}
}

// pythonExecutor is a real local sandbox for the tests that need code to run.
func pythonExecutor(t *testing.T) runner.Executor {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3 on PATH")
	}
	e, err := runner.New(runner.Config{Mode: "local", PythonPath: "python3", AllowUnsafe: true})
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	return e
}
