package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/cohort"
	"github.com/RA9/gamifydev/platform/internal/placement"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// The placement diagnostic is the admissions gate. Guests may sit it before
// creating an account; only a passing, unclaimed result can create a learner.
// A pass from 50–69 routes to CS Foundations, while 70+ routes to a
// specialization. Placement never grants course credit.

func (s *Server) handlePlacement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by := s.solver(r)

	if by != (store.Solver{}) {
		live, err := s.st.LiveAttemptFor(ctx, by)
		if err != nil {
			http.Error(w, "could not load your placement test", http.StatusInternalServerError)
			return
		}
		if live != nil {
			items, err := s.st.AttemptItemsFor(ctx, by, live.ID)
			if err != nil {
				http.Error(w, "could not load your questions", http.StatusInternalServerError)
				return
			}
			s.render(w, r, "placement.html", ViewData{Title: "Placement test", Data: map[string]any{
				"bodyClass": "placement-dark",
				"attempt":   live,
				"items":     items,
				"expiresAt": live.ExpiresAt,
				"total":     len(items),
			}})
			return
		}
	}

	s.renderPlacementIntro(w, r, by, "")
}

func (s *Server) renderPlacementIntro(w http.ResponseWriter, r *http.Request, by store.Solver, flash string) {
	ctx := r.Context()
	a, err := s.st.GetAssessment(ctx, "placement")
	if err != nil {
		http.Error(w, "the placement test is not available", http.StatusServiceUnavailable)
		return
	}

	var prev *store.PlacementResult
	retry := store.PlacementRetryStatus{Allowed: true}
	if by != (store.Solver{}) {
		prev, _ = s.st.LatestPlacementFor(ctx, by)
		if status, err := s.st.PlacementRetryFor(ctx, by); err == nil {
			retry = status
		}
	}

	s.render(w, r, "placement_intro.html", ViewData{Title: "Placement test", Flash: flash, Data: map[string]any{
		"bodyClass":            "placement-dark",
		"assessment":           a,
		"minutes":              a.TimeLimitS / 60,
		"done":                 prev != nil,
		"result":               prev,
		"retry":                retry,
		"topics":               placement.AllTopics,
		"passThreshold":        placement.PassThreshold,
		"foundationsThreshold": placement.FoundationsThreshold,
	}})
}

// handlePlacementStart creates a guest identity only when a visitor actually
// starts a paper. Merely browsing the introduction creates no database state.
func (s *Server) handlePlacementStart(w http.ResponseWriter, r *http.Request) {
	by, err := s.solverForWrite(w, r)
	if err != nil {
		http.Error(w, "could not start the placement test", http.StatusInternalServerError)
		return
	}
	if live, _ := s.st.LiveAttemptFor(r.Context(), by); live != nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}

	a, err := s.st.GetAssessment(r.Context(), "placement")
	if err != nil {
		http.Error(w, "the placement test is not available", http.StatusServiceUnavailable)
		return
	}
	if _, err := s.st.StartAttemptFor(r.Context(), by, a); err != nil {
		switch {
		case errors.Is(err, store.ErrPlacementAlreadyPassed), errors.Is(err, store.ErrPlacementEnrollmentLocked):
			http.Redirect(w, r, "/placement/result", http.StatusSeeOther)
		case errors.Is(err, store.ErrPlacementRetryTooSoon):
			s.renderPlacementIntro(w, r, by, "You need to wait seven days after an unsuccessful attempt before trying again.")
		case errors.Is(err, store.ErrPlacementAttemptLimit):
			s.renderPlacementIntro(w, r, by, "You have used three attempts in the last 30 days. Your next attempt opens when the rolling limit resets.")
		default:
			http.Error(w, "could not start the placement test", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/placement", http.StatusSeeOther)
}

func (s *Server) handlePlacementSubmit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by := s.solver(r)
	if by == (store.Solver{}) {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	live, err := s.st.LiveAttemptFor(ctx, by)
	if err != nil || live == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not read your answers", http.StatusBadRequest)
		return
	}

	responses := map[int64]int{}
	for key, vals := range r.Form {
		if !strings.HasPrefix(key, "q") || len(vals) == 0 {
			continue
		}
		itemID, err := strconv.ParseInt(key[1:], 10, 64)
		if err != nil {
			continue
		}
		choice, err := strconv.Atoi(vals[0])
		if err == nil {
			responses[itemID] = choice
		}
	}

	attempt, err := s.st.ScoreAttemptFor(ctx, by, live.ID, responses)
	if errors.Is(err, store.ErrAttemptExpired) {
		s.render(w, r, "placement_expired.html", ViewData{Title: "Time's up", Data: map[string]any{"bodyClass": "placement-dark"}})
		return
	}
	if err != nil {
		http.Error(w, "could not score your placement test", http.StatusInternalServerError)
		return
	}

	decision := placement.Decide(attempt.TopicScores)
	result := store.PlacementResult{
		AttemptID:           attempt.ID,
		By:                  by,
		Passed:              decision.Passed,
		FoundationsRequired: decision.FoundationsRequired,
		Exemptions:          decision.Exemptions,
	}
	if decision.Passed {
		p, err := s.st.GetPathBySlug(ctx, decision.RecommendedPath)
		if err != nil || !p.Published {
			http.Error(w, "the generated learning path is unavailable", http.StatusInternalServerError)
			return
		}
		result.RecommendedPathID = sql.NullInt64{Int64: p.ID, Valid: true}
	}
	if _, err := s.st.SavePlacementResultFor(ctx, by, result); err != nil {
		if errors.Is(err, store.ErrPlacementEnrollmentLocked) {
			http.Redirect(w, r, "/placement/result", http.StatusSeeOther)
			return
		}
		http.Error(w, "could not save your result", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/placement/result", http.StatusSeeOther)
}

func (s *Server) handlePlacementResult(w http.ResponseWriter, r *http.Request) {
	s.renderPlacementResult(w, r, "")
}

func (s *Server) renderPlacementResult(w http.ResponseWriter, r *http.Request, flash string) {
	ctx := r.Context()
	by := s.solver(r)
	if by == (store.Solver{}) {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	result, err := s.st.LatestPlacementFor(ctx, by)
	if err != nil || result == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	var enrollment *store.Enrollment
	if by.UserID != 0 {
		enrollment, err = s.st.EnsureEnrollment(ctx, by.UserID)
		if err != nil {
			http.Error(w, "could not load your enrollment", http.StatusInternalServerError)
			return
		}
		if enrollment.PlacementResultID.Valid {
			result, err = s.st.PlacementResultByIDFor(ctx, by, enrollment.PlacementResultID.Int64)
			if err != nil {
				http.Error(w, "could not load your enrolled placement", http.StatusInternalServerError)
				return
			}
		}
	}
	attempt, err := s.st.GetAttemptFor(ctx, by, result.AttemptID)
	if err != nil {
		http.Error(w, "could not load your result", http.StatusInternalServerError)
		return
	}

	type topicRow struct {
		Topic, Label string
		Score        int
		Core         bool
	}
	rows := make([]topicRow, 0, len(placement.AllTopics))
	for _, topic := range placement.AllTopics {
		rows = append(rows, topicRow{
			Topic: topic, Label: placement.TopicLabel(topic),
			Score: attempt.TopicScores[topic], Core: placement.IsCore(topic),
		})
	}

	type exemptRow struct{ Slug, Title, Emoji string }
	exemptions := make([]exemptRow, 0, len(result.Exemptions))
	for _, slug := range result.Exemptions {
		row := exemptRow{Slug: slug, Title: slug}
		if course, err := s.st.GetCourseBySlug(ctx, slug); err == nil {
			row.Title, row.Emoji = course.Title, course.Emoji
		}
		exemptions = append(exemptions, row)
	}

	var recommended *store.Path
	if result.RecommendedPathID.Valid {
		recommended, _ = s.st.GetPathByID(ctx, result.RecommendedPathID.Int64)
	}
	retry, _ := s.st.PlacementRetryFor(ctx, by)

	s.render(w, r, "placement_result.html", ViewData{Title: "Your placement", Flash: flash, Data: map[string]any{
		"bodyClass":            "placement-dark",
		"result":               result,
		"attempt":              attempt,
		"topics":               rows,
		"coreScore":            placement.MeanCore(attempt.TopicScores),
		"passThreshold":        placement.PassThreshold,
		"foundationsThreshold": placement.FoundationsThreshold,
		"recommended":          recommended,
		"exemptions":           exemptions,
		"enrollment":           enrollment,
		"bands":                cohort.Bands(),
		"signedIn":             by.UserID != 0,
		"retry":                retry,
	}})
}

func (s *Server) handlePlacementReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by := s.solver(r)
	if by == (store.Solver{}) {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	result, err := s.st.LatestPlacementFor(ctx, by)
	if err != nil || result == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	if by.UserID != 0 {
		enrollment, err := s.st.LiveEnrollment(ctx, by.UserID)
		if err != nil {
			http.Error(w, "could not load your enrollment", http.StatusInternalServerError)
			return
		}
		if enrollment != nil && enrollment.PlacementResultID.Valid {
			result, err = s.st.PlacementResultByIDFor(ctx, by, enrollment.PlacementResultID.Int64)
			if err != nil {
				http.Error(w, "could not load your enrolled placement", http.StatusInternalServerError)
				return
			}
		}
	}
	// Failed attempts receive topic-level feedback only. Revealing every answer
	// after each failure would turn retakes into a way to harvest the bank.
	if !result.Passed {
		http.Redirect(w, r, "/placement/result", http.StatusSeeOther)
		return
	}
	attempt, err := s.st.GetAttemptFor(ctx, by, result.AttemptID)
	if err != nil || !attempt.Submitted() {
		s.notFound(w, r)
		return
	}
	review, err := s.st.ReviewAttemptFor(ctx, by, attempt.ID)
	if err != nil {
		http.Error(w, "could not load your answers", http.StatusInternalServerError)
		return
	}
	s.render(w, r, "placement_review.html", ViewData{Title: "Your answers", Data: map[string]any{
		"bodyClass": "placement-dark",
		"review":    review,
		"attempt":   attempt,
	}})
}

// handleEnroll commits an authenticated learner to the exact path and
// exemptions generated by the displayed passing placement result. The client
// supplies only that result's identity and a timezone band; the store derives
// every curriculum field from the owned result in one transaction.
func (s *Server) handleEnroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)
	resultID, err := strconv.ParseInt(r.FormValue("placement_result"), 10, 64)
	if err != nil || resultID <= 0 {
		s.renderPlacementResult(w, r, "That placement result is no longer available. Review your current result and try again.")
		return
	}
	band := strings.TrimSpace(r.FormValue("band"))
	if !cohort.ValidBand(band) {
		s.renderPlacementResult(w, r, "Choose a valid timezone region so we can place you in the right cohort.")
		return
	}

	if _, err := s.st.EnrollFromPlacement(ctx, u.ID, resultID, band); err != nil {
		switch {
		case errors.Is(err, store.ErrEnrollmentLocked):
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		case errors.Is(err, store.ErrEnrollmentPlacementInvalid):
			s.renderPlacementResult(w, r, "That result cannot be used for enrollment. Review your current placement result and try again.")
		default:
			http.Error(w, "could not enroll you", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
