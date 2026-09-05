package server

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/placement"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func (ts *testServer) qualifyGuest(t *testing.T, by store.Solver, passed bool) {
	t.Helper()
	assessment, err := ts.st.GetAssessment(t.Context(), "placement")
	if err != nil {
		t.Fatalf("get placement: %v", err)
	}
	attempt, err := ts.st.StartAttemptFor(t.Context(), by, assessment)
	if err != nil {
		t.Fatalf("start placement: %v", err)
	}
	if _, err := ts.st.ScoreAttemptFor(t.Context(), by, attempt.ID, map[int64]int{}); err != nil {
		t.Fatalf("score placement: %v", err)
	}
	result := store.PlacementResult{AttemptID: attempt.ID, By: by, Passed: passed}
	if passed {
		path, err := ts.st.GetPathBySlug(t.Context(), placement.PathFoundations)
		if err != nil {
			t.Fatalf("get foundations path: %v", err)
		}
		result.FoundationsRequired = true
		result.RecommendedPathID.Valid = true
		result.RecommendedPathID.Int64 = path.ID
	}
	if _, err := ts.st.SavePlacementResultFor(t.Context(), by, result); err != nil {
		t.Fatalf("save placement result: %v", err)
	}
}

func TestPlacementIsPublicButBrowsingDoesNotCreateAGuest(t *testing.T) {
	ts := newTestServer(t)
	w := ts.do(t, http.MethodGet, "/placement", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /placement = %d, want 200", w.Code)
	}
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "gd_guest" && cookie.Value != "" {
			t.Fatal("browsing placement created a guest identity")
		}
	}
}

func TestRegistrationRequiresPassingGuestPlacement(t *testing.T) {
	ts := newTestServer(t)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		var form url.Values
		if method == http.MethodPost {
			form = url.Values{"name": {"Bypass"}, "email": {"bypass@example.com"}, "password": {"averysafepassword"}}
		}
		w := ts.do(t, method, "/register", form)
		if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement" {
			t.Errorf("%s /register = %d -> %q, want 303 -> /placement", method, w.Code, w.Header().Get("Location"))
		}
	}
	if _, err := ts.st.GetUserByEmail(t.Context(), "bypass@example.com"); err == nil {
		t.Fatal("direct registration bypass created a user")
	}

	cookie, by := ts.guest(t)
	ts.qualifyGuest(t, by, false)
	w := ts.do(t, http.MethodPost, "/register", url.Values{
		"name": {"Failed Guest"}, "email": {"failed@example.com"}, "password": {"averysafepassword"},
	}, cookie)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement" {
		t.Fatalf("failed guest registration = %d -> %q", w.Code, w.Header().Get("Location"))
	}
}

func (ts *testServer) registerQualifiedGuest(t *testing.T, email string) (*http.Cookie, *store.User, *store.PlacementResult) {
	t.Helper()
	guestCookie, by := ts.guest(t)
	ts.qualifyGuest(t, by, true)
	w := ts.do(t, http.MethodPost, "/register", url.Values{
		"name": {"Qualified Learner"}, "email": {email}, "password": {"averysafepassword"},
	}, guestCookie)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement/result" {
		t.Fatalf("register = %d -> %q, want 303 -> /placement/result", w.Code, w.Header().Get("Location"))
	}
	var session *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == auth.SessionCookie && cookie.Value != "" {
			session = cookie
			break
		}
	}
	if session == nil {
		t.Fatal("qualified registration did not create a session")
	}
	u, err := ts.st.GetUserByEmail(t.Context(), email)
	if err != nil {
		t.Fatalf("get learner: %v", err)
	}
	result, err := ts.st.LatestPlacement(t.Context(), u.ID)
	if err != nil || result == nil {
		t.Fatalf("get claimed result: %+v, %v", result, err)
	}
	return session, u, result
}

func TestPassingGuestCreatesLearnerAndClaimsResult(t *testing.T) {
	ts := newTestServer(t)
	cookie, by := ts.guest(t)
	ts.qualifyGuest(t, by, true)

	w := ts.do(t, http.MethodPost, "/register", url.Values{
		"name": {"Qualified Learner"}, "email": {"qualified@example.com"}, "password": {"averysafepassword"},
	}, cookie)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement/result" {
		t.Fatalf("register = %d -> %q, want 303 -> /placement/result", w.Code, w.Header().Get("Location"))
	}
	u, err := ts.st.GetUserByEmail(t.Context(), "qualified@example.com")
	if err != nil {
		t.Fatalf("get learner: %v", err)
	}
	if u.Role != "learner" {
		t.Fatalf("first public account role = %q, want learner", u.Role)
	}
	result, err := ts.st.LatestPlacement(t.Context(), u.ID)
	if err != nil || result == nil || !result.Passed {
		t.Fatalf("claimed result = %+v, %v", result, err)
	}
	if result.GuestID != 0 || result.UserID != u.ID {
		t.Fatalf("result owner = user:%d guest:%d, want user:%d", result.UserID, result.GuestID, u.ID)
	}
	enrollment, err := ts.st.LiveEnrollment(t.Context(), u.ID)
	if err != nil {
		t.Fatalf("get enrollment: %v", err)
	}
	if enrollment == nil || enrollment.State != store.EnrollUnplaced || enrollment.PathID.Valid || enrollment.PlacementResultID.Valid {
		t.Fatalf("new learner enrollment = %+v, want unplaced and unbound", enrollment)
	}
}

func TestGeneratedEnrollmentBindsExactResultAndIgnoresPostedPath(t *testing.T) {
	ts := newTestServer(t)
	session, user, result := ts.registerQualifiedGuest(t, "enroll@example.com")

	form := url.Values{
		"placement_result": {strconv.FormatInt(result.ID, 10)},
		"band":             {"americas"},
		"path":             {"forged-path-that-must-be-ignored"},
	}
	w := ts.do(t, http.MethodPost, "/enroll", form, session)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/dashboard" {
		t.Fatalf("enroll = %d -> %q, want 303 -> /dashboard", w.Code, w.Header().Get("Location"))
	}
	enrollment, err := ts.st.LiveEnrollment(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("get enrollment: %v", err)
	}
	if enrollment == nil || enrollment.State != store.EnrollPlaced {
		t.Fatalf("enrollment = %+v, want placed", enrollment)
	}
	if !enrollment.PathID.Valid || enrollment.PathID != result.RecommendedPathID {
		t.Fatalf("bound path = %+v, want generated %+v", enrollment.PathID, result.RecommendedPathID)
	}
	if !enrollment.PlacementResultID.Valid || enrollment.PlacementResultID.Int64 != result.ID {
		t.Fatalf("bound result = %+v, want %d", enrollment.PlacementResultID, result.ID)
	}
	if enrollment.TZBand != "americas" {
		t.Fatalf("timezone band = %q, want americas", enrollment.TZBand)
	}
	if len(enrollment.Exemptions) != len(result.Exemptions) {
		t.Fatalf("exemption snapshot = %v, want %v", enrollment.Exemptions, result.Exemptions)
	}

	// Browser retries of the same POST are idempotent.
	w = ts.do(t, http.MethodPost, "/enroll", form, session)
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/dashboard" {
		t.Fatalf("repeat enroll = %d -> %q", w.Code, w.Header().Get("Location"))
	}
	again, err := ts.st.LiveEnrollment(t.Context(), user.ID)
	if err != nil || again.ID != enrollment.ID || again.PlacementResultID.Int64 != result.ID {
		t.Fatalf("repeat enroll changed binding: %+v, %v", again, err)
	}
}

func TestGeneratedEnrollmentRejectsInvalidBandAndForeignResult(t *testing.T) {
	ts := newTestServer(t)
	session, user, result := ts.registerQualifiedGuest(t, "protected@example.com")

	w := ts.do(t, http.MethodPost, "/enroll", url.Values{
		"placement_result": {strconv.FormatInt(result.ID, 10)},
		"band":             {"not-a-band"},
	}, session)
	if w.Code != http.StatusOK {
		t.Fatalf("invalid band response = %d, want rendered result page", w.Code)
	}
	enrollment, err := ts.st.LiveEnrollment(t.Context(), user.ID)
	if err != nil || enrollment == nil || enrollment.State != store.EnrollUnplaced || enrollment.PathID.Valid || enrollment.PlacementResultID.Valid {
		t.Fatalf("invalid band mutated enrollment: %+v, %v", enrollment, err)
	}

	_, _, foreign := ts.registerQualifiedGuest(t, "other@example.com")
	w = ts.do(t, http.MethodPost, "/enroll", url.Values{
		"placement_result": {strconv.FormatInt(foreign.ID, 10)},
		"band":             {"europe_africa"},
	}, session)
	if w.Code != http.StatusOK {
		t.Fatalf("foreign result response = %d, want rendered result page", w.Code)
	}
	enrollment, err = ts.st.LiveEnrollment(t.Context(), user.ID)
	if err != nil || enrollment.State != store.EnrollUnplaced || enrollment.PathID.Valid || enrollment.PlacementResultID.Valid {
		t.Fatalf("foreign result mutated enrollment: %+v, %v", enrollment, err)
	}
}

func TestPlacementSubmissionAppliesAcceptedScoreBands(t *testing.T) {
	for _, tc := range []struct {
		name           string
		correctPerCore int
		passed         bool
		foundations    bool
	}{
		{name: "below 50 fails", correctPerCore: 1},
		{name: "50 passes to foundations", correctPerCore: 2, passed: true, foundations: true},
		{name: "100 passes to specialization", correctPerCore: 4, passed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestServer(t)
			assessment, err := ts.st.GetAssessment(t.Context(), "placement")
			if err != nil {
				t.Fatalf("assessment: %v", err)
			}
			var bank []store.Item
			for _, topic := range placement.AllTopics {
				for i := 0; i < 4; i++ {
					bank = append(bank, store.Item{Topic: topic, Prompt: topic, Options: []string{"wrong", "right"}, Answer: 1})
				}
			}
			if err := ts.st.ReplaceItems(t.Context(), assessment.ID, bank); err != nil {
				t.Fatalf("replace bank: %v", err)
			}
			cookie, by := ts.guest(t)
			if w := ts.do(t, http.MethodPost, "/placement/start", url.Values{}, cookie); w.Code != http.StatusSeeOther {
				t.Fatalf("start = %d", w.Code)
			}
			attempt, _ := ts.st.LiveAttemptFor(t.Context(), by)
			items, _ := ts.st.AttemptItemsFor(t.Context(), by, attempt.ID)
			form := url.Values{}
			seen := map[string]int{}
			for _, item := range items {
				choice := "0"
				if !placement.IsCore(item.Topic) || seen[item.Topic] < tc.correctPerCore {
					choice = "1"
				}
				seen[item.Topic]++
				form.Set("q"+strconv.FormatInt(item.ID, 10), choice)
			}
			w := ts.do(t, http.MethodPost, "/placement/submit", form, cookie)
			if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/placement/result" {
				t.Fatalf("submit = %d -> %q", w.Code, w.Header().Get("Location"))
			}
			result, err := ts.st.LatestPlacementFor(t.Context(), by)
			if err != nil || result == nil {
				t.Fatalf("result = %+v, %v", result, err)
			}
			if result.Passed != tc.passed || result.FoundationsRequired != tc.foundations {
				t.Fatalf("result passed=%t foundations=%t, want %t/%t", result.Passed, result.FoundationsRequired, tc.passed, tc.foundations)
			}
		})
	}
}
