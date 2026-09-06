package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/seed"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// These are smoke tests: they drive the real mux, the real templates and a real
// migrated database with the real seed content, and assert that pages come back
// and that the doors that are supposed to be locked are locked.
//
// They are deliberately shallow. Their job is to catch a route that stopped
// rendering, a template that stopped compiling, and — most of all — a handler
// that was registered without the middleware that guards it, which is the class
// of mistake nothing else here would notice.

type testServer struct {
	*Server
	h  http.Handler
	st *store.Store
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	return newTestServerWith(t, nil)
}

// newTestServerWith is newTestServer with an executor, for the handful of tests
// that need code to actually run rather than to be told execution is off.
func newTestServerWith(t *testing.T, exec runner.Executor) *testServer {
	t.Helper()
	ctx := context.Background()

	st, err := store.Open(filepath.Join(t.TempDir(), "smoke.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// Real content, so the public pages render against something rather than
	// against empty tables — an empty database hides exactly the template bugs
	// worth catching.
	if _, err := seed.Run(ctx, st); err != nil {
		t.Fatalf("seed: %v", err)
	}

	srv, err := New(st, nil, nil, exec, jobs.New(st, "test"), false, false)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return &testServer{Server: srv, h: srv.Routes(), st: st}
}

// do issues a request carrying the given cookies and returns the recorder.
// Redirects are never followed: which redirect a handler chose is the thing
// under test.
func (ts *testServer) do(t *testing.T, method, path string, form url.Values, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if form != nil {
		r = httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	ts.h.ServeHTTP(w, r)
	return w
}

// register provisions an authenticated fixture directly. Public registration
// is deliberately tested separately because it now requires a passing guest
// placement; unrelated route tests should not repeat that whole admissions flow.
func (ts *testServer) register(t *testing.T, email string) *http.Cookie {
	t.Helper()
	role := "learner"
	if email == "admin@example.com" {
		role = "admin"
	}
	hash, err := auth.HashPassword("averysafepassword")
	if err != nil {
		t.Fatalf("hash fixture password: %v", err)
	}
	u, err := ts.st.CreateUser(t.Context(), email, hash, "Test Person", role)
	if err != nil {
		t.Fatalf("create fixture user %s: %v", email, err)
	}
	token, err := auth.NewToken()
	if err != nil {
		t.Fatalf("create fixture session: %v", err)
	}
	if err := ts.st.CreateSession(t.Context(), token, u.ID, time.Now().Add(auth.SessionTTL)); err != nil {
		t.Fatalf("create fixture session: %v", err)
	}
	return &http.Cookie{Name: auth.SessionCookie, Value: token}
}

func TestRestrictedAccountCannotUseParticipationRoutes(t *testing.T) {
	ts := newTestServer(t)
	session := ts.register(t, "restricted@example.com")
	u, err := ts.st.GetUserByEmail(t.Context(), "restricted@example.com")
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	if err := ts.st.SetAccountState(t.Context(), u.ID, store.AccountSuspended, "review", nil); err != nil {
		t.Fatalf("suspend: %v", err)
	}

	for _, path := range []string{"/api/run", "/cohort/standup", "/attendance/absence", "/forum"} {
		w := ts.do(t, http.MethodPost, path, url.Values{}, session)
		if w.Code != http.StatusForbidden {
			t.Fatalf("POST %s = %d, want 403", path, w.Code)
		}
	}
	// Record and appeal pages remain readable while participation is restricted.
	if w := ts.do(t, http.MethodGet, "/attendance", nil, session); w.Code != http.StatusOK {
		t.Fatalf("GET /attendance = %d, want 200", w.Code)
	}
}

// --- public pages -----------------------------------------------------------

func TestPublicPagesRender(t *testing.T) {
	ts := newTestServer(t)
	for _, tc := range []struct{ path, wantText string }{
		{"/", "GamifyDev"},
		{"/about", ""},
		{"/paths", "Computer Science Foundations"},
		{"/paths/cs-foundations", "Data Structures"},
		{"/courses", ""},
		{"/courses/c", ""},
		{"/login", "Forgot it?"},
		{"/placement", `action="/placement/start"`},
		{"/forgot", "Reset your password"},
		{"/forum", ""},
		{"/competitions", ""},
		{"/blog", ""},
		{"/healthz", "ok"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			w := ts.do(t, http.MethodGet, tc.path, nil)
			if w.Code != http.StatusOK {
				t.Fatalf("status %d, want 200", w.Code)
			}
			if tc.wantText != "" && !strings.Contains(w.Body.String(), tc.wantText) {
				t.Errorf("body does not mention %q", tc.wantText)
			}
		})
	}
}

func TestUnknownPagesAreNotFoundRatherThanAnError(t *testing.T) {
	ts := newTestServer(t)
	for _, path := range []string{"/courses/nope", "/courses/c/nope", "/paths/nope", "/blog/nope"} {
		w := ts.do(t, http.MethodGet, path, nil)
		if w.Code != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", path, w.Code)
		}
	}
}

// --- the auth boundary ------------------------------------------------------

// routePattern pulls the registered method+path out of the routing table.
//
// Reading the source rather than listing routes by hand is the point: a route
// added next month is covered by these tests the day it is added, which is the
// only way this stays a guard rather than a snapshot of what was true once.
var routePattern = regexp.MustCompile(`mux\.(?:Handle|HandleFunc)\("(GET|POST) ([^"]+)"`)

type route struct{ method, path string }

func registeredRoutes(t *testing.T) []route {
	t.Helper()
	src, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatalf("read routing table: %v", err)
	}
	var routes []route
	for _, m := range routePattern.FindAllStringSubmatch(string(src), -1) {
		path := m[2]
		if strings.Contains(path, "{$}") {
			path = "/"
		}
		// Give every wildcard a concrete value; these requests are meant to be
		// turned away by middleware long before the value is looked at.
		path = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(path, "1")
		routes = append(routes, route{m[1], path})
	}
	if len(routes) < 50 {
		t.Fatalf("found only %d routes — the pattern stopped matching the routing table", len(routes))
	}
	return routes
}

// publicPaths are reachable without signing in, by design.
func isPublic(p string) bool {
	for _, prefix := range []string{
		"/static/", "/healthz", "/about", "/paths", "/courses", "/forum",
		"/competitions", "/blog", "/login", "/register", "/placement", "/forgot", "/reset/",
		"/invite/", "/logout",
		// The practice bank is open to anyone by design — solving problems
		// without an account is the point of it. What it does not expose is
		// other people's work: /problems/result/{id} answers only the solver
		// who wrote the submission, which TestASubmissionIsNotReadableByItsID
		// is what actually holds that line.
		"/problems",
	} {
		if p == prefix || strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return p == "/"
}

func TestEveryPrivateRouteTurnsAwayAGuest(t *testing.T) {
	ts := newTestServer(t)
	for _, rt := range registeredRoutes(t) {
		if isPublic(rt.path) {
			continue
		}
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			var form url.Values
			if rt.method == http.MethodPost {
				form = url.Values{}
			}
			w := ts.do(t, rt.method, rt.path, form)
			if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/login" {
				t.Errorf("guest got %d -> %q, want 303 -> /login. A handler registered "+
					"without requireAuth/requireRole is reachable by anyone",
					w.Code, w.Header().Get("Location"))
			}
		})
	}
}

func TestEveryAdminRouteTurnsAwayALearner(t *testing.T) {
	ts := newTestServer(t)
	ts.register(t, "admin@example.com")
	learner := ts.register(t, "learner@example.com")

	for _, rt := range registeredRoutes(t) {
		if !strings.HasPrefix(rt.path, "/admin") {
			continue
		}
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			var form url.Values
			if rt.method == http.MethodPost {
				form = url.Values{}
			}
			w := ts.do(t, rt.method, rt.path, form, learner)
			if w.Code != http.StatusForbidden {
				t.Errorf("learner got %d, want 403 — this admin route is reachable "+
					"by any signed-in user", w.Code)
			}
		})
	}
}

func TestGradingIsOpenToGradersButTheRestOfAdminIsNot(t *testing.T) {
	ts := newTestServer(t)
	admin := ts.register(t, "admin@example.com")
	grader := ts.register(t, "grader@example.com")

	u, err := ts.st.GetUserByEmail(context.Background(), "grader@example.com")
	if err != nil {
		t.Fatalf("lookup grader: %v", err)
	}
	if err := ts.st.SetUserRole(context.Background(), u.ID, "grader"); err != nil {
		t.Fatalf("set role: %v", err)
	}

	if w := ts.do(t, http.MethodGet, "/admin/grading", nil, grader); w.Code != http.StatusOK {
		t.Errorf("grader on /admin/grading: %d, want 200", w.Code)
	}
	if w := ts.do(t, http.MethodGet, "/admin/users", nil, grader); w.Code != http.StatusForbidden {
		t.Errorf("grader on /admin/users: %d, want 403", w.Code)
	}
	if w := ts.do(t, http.MethodGet, "/admin/users", nil, admin); w.Code != http.StatusOK {
		t.Errorf("admin on /admin/users: %d, want 200", w.Code)
	}
}

// --- sessions ---------------------------------------------------------------

func TestSessionLifecycle(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "someone@example.com")

	if w := ts.do(t, http.MethodGet, "/dashboard", nil, c); w.Code != http.StatusOK {
		t.Fatalf("dashboard with a session: %d, want 200", w.Code)
	}
	if w := ts.do(t, http.MethodPost, "/logout", url.Values{}, c); w.Code != http.StatusSeeOther {
		t.Fatalf("logout: %d, want 303", w.Code)
	}
	// The cookie value is now a dead token; presenting it must not sign anyone in.
	w := ts.do(t, http.MethodGet, "/dashboard", nil, c)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("dashboard after logout: %d, want 303 — the session outlived the sign-out", w.Code)
	}
}

func TestAForgedSessionCookieIsIgnored(t *testing.T) {
	ts := newTestServer(t)
	forged := &http.Cookie{Name: "gd_session", Value: strings.Repeat("a", 64)}
	if w := ts.do(t, http.MethodGet, "/dashboard", nil, forged); w.Code != http.StatusSeeOther {
		t.Fatalf("a made-up session token got %d, want 303 to /login", w.Code)
	}
}

func TestSignedInUsersAreSentOnFromTheAuthPages(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "someone@example.com")
	for _, path := range []string{"/login", "/register", "/forgot"} {
		w := ts.do(t, http.MethodGet, path, nil, c)
		if w.Code != http.StatusSeeOther {
			t.Errorf("%s while signed in: %d, want a redirect away", path, w.Code)
		}
	}
}

// --- learner pages ----------------------------------------------------------

func TestLearnerPagesRender(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "learner@example.com")
	for _, path := range []string{
		"/dashboard", "/assignments", "/assignments/c-max-three",
		"/playground", "/placement", "/attendance", "/cohort",
	} {
		t.Run(path, func(t *testing.T) {
			w := ts.do(t, http.MethodGet, path, nil, c)
			if w.Code != http.StatusOK {
				t.Fatalf("status %d, want 200", w.Code)
			}
		})
	}
}

func TestScheduleSendsALearnerWithNoCohortToTheCohortPage(t *testing.T) {
	// A schedule is a cohort's shared plan, so there is nothing to show before
	// the learner is in one. The redirect is the answer, not a 200 with an
	// empty page.
	ts := newTestServer(t)
	c := ts.register(t, "learner@example.com")
	w := ts.do(t, http.MethodGet, "/schedule", nil, c)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/cohort" {
		t.Fatalf("got %d -> %q, want 303 -> /cohort", w.Code, w.Header().Get("Location"))
	}
}

func TestAdminPagesRender(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "admin@example.com")
	for _, path := range []string{
		"/admin", "/admin/users", "/admin/courses", "/admin/assignments",
		"/admin/paths", "/admin/grading", "/admin/cohorts", "/admin/attendance",
		"/admin/blog", "/admin/competitions", "/admin/forum", "/admin/jobs",
	} {
		t.Run(path, func(t *testing.T) {
			w := ts.do(t, http.MethodGet, path, nil, c)
			if w.Code != http.StatusOK {
				t.Fatalf("status %d, want 200", w.Code)
			}
		})
	}
}

// --- the gate ---------------------------------------------------------------

func TestThePathPageLocksCoursesBehindAnUnmetCheckpoint(t *testing.T) {
	// The gate is the platform's main mechanic. If the page stopped rendering
	// lock state, every learner would see the whole curriculum as open.
	ts := newTestServer(t)
	c := ts.register(t, "learner@example.com")
	body := ts.do(t, http.MethodGet, "/paths/cs-foundations", nil, c).Body.String()

	if !strings.Contains(body, "Locked") {
		t.Fatal("nothing is locked for a learner who has passed no checkpoint")
	}
	// The first course is the way in and must never be locked.
	first := strings.Index(body, "Program close to the metal")
	firstLock := strings.Index(body, "Locked")
	if first == -1 {
		t.Fatal("the first course is missing from the page")
	}
	if firstLock < first {
		t.Error("the first course in the path is locked, so there is no way to start")
	}
}

func TestSubmittingAnAssignmentIsRecorded(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "learner@example.com")

	w := ts.do(t, http.MethodPost, "/assignments/py-sum-list/submit",
		url.Values{"code": {"def sum_list(nums):\n    return sum(nums)\n"}, "note": {"done"}}, c)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("submit: %d, want 303", w.Code)
	}
	body := ts.do(t, http.MethodGet, "/assignments/py-sum-list", nil, c).Body.String()
	if !strings.Contains(body, "Your latest submission") {
		t.Error("the submission does not show on the assignment page")
	}
}

func TestAnEmptySubmissionIsRejected(t *testing.T) {
	ts := newTestServer(t)
	c := ts.register(t, "learner@example.com")
	w := ts.do(t, http.MethodPost, "/assignments/py-sum-list/submit", url.Values{"code": {"   "}}, c)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("status %d, want 303", w.Code)
	}
	if strings.Contains(w.Header().Get("Location"), "submitted=1") {
		t.Error("an empty submission was accepted")
	}
}
