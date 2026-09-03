package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/placement"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// The placement diagnostic and the enrollment gate it feeds.
//
// The rule (PRD §03): the diagnostic *routes* — it picks a path and proposes
// which courses to skip in the schedule — but it never grants credit. Advancing
// past a course still requires its checkpoint. So nothing in this file unlocks
// content; it only decides what a learner is enrolled in.

// handlePlacement shows the diagnostic: an intro before starting, the paper
// while a sitting is live, and a pointer to the result once taken.
func (s *Server) handlePlacement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	live, err := s.st.LiveAttempt(ctx, u.ID)
	if err != nil {
		http.Error(w, "could not load your diagnostic", http.StatusInternalServerError)
		return
	}
	if live != nil {
		items, err := s.st.AttemptItems(ctx, live.ID)
		if err != nil {
			http.Error(w, "could not load your questions", http.StatusInternalServerError)
			return
		}
		s.render(w, r, "placement.html", ViewData{Title: "Placement diagnostic", Data: map[string]any{
			"bodyClass": "placement-dark",
			"attempt":   live,
			"items":     items,
			"expiresAt": live.ExpiresAt,
			"total":     len(items),
		}})
		return
	}

	// No sitting in progress: show the intro, or the existing result.
	prev, _ := s.st.LatestPlacement(ctx, u.ID)
	a, err := s.st.GetAssessment(ctx, "placement")
	if err != nil {
		http.Error(w, "the diagnostic is not available", http.StatusServiceUnavailable)
		return
	}
	s.render(w, r, "placement_intro.html", ViewData{Title: "Placement diagnostic", Data: map[string]any{
		"bodyClass":  "placement-dark",
		"assessment": a,
		"minutes":    a.TimeLimitS / 60,
		"done":       prev != nil,
		"topics":     placement.AllTopics,
	}})
}

// handlePlacementStart opens a sitting and draws a paper.
func (s *Server) handlePlacementStart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	// An existing live sitting wins — restarting would be a way to reroll the
	// paper until an easy one turns up.
	if live, _ := s.st.LiveAttempt(ctx, u.ID); live != nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	a, err := s.st.GetAssessment(ctx, "placement")
	if err != nil {
		http.Error(w, "the diagnostic is not available", http.StatusServiceUnavailable)
		return
	}
	if _, err := s.st.StartAttempt(ctx, u.ID, a); err != nil {
		http.Error(w, "could not start the diagnostic", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/placement", http.StatusSeeOther)
}

// handlePlacementSubmit grades a sitting and records the routing decision.
func (s *Server) handlePlacementSubmit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	live, err := s.st.LiveAttempt(ctx, u.ID)
	if err != nil || live == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not read your answers", http.StatusBadRequest)
		return
	}

	// Answers arrive as q<itemID>=<optionIndex>. Anything unparseable is simply
	// left unanswered rather than failing the submission — a learner should not
	// lose a sitting to a malformed field.
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
		if err != nil {
			continue
		}
		responses[itemID] = choice
	}

	attempt, err := s.st.ScoreAttempt(ctx, live.ID, responses)
	if errors.Is(err, store.ErrAttemptExpired) {
		s.render(w, r, "placement_expired.html", ViewData{Title: "Time's up", Data: map[string]any{
			"bodyClass": "placement-dark",
		}})
		return
	}
	if err != nil {
		http.Error(w, "could not score your diagnostic", http.StatusInternalServerError)
		return
	}

	// Route, then record the decision.
	decision := placement.Decide(attempt.TopicScores)
	res := store.PlacementResult{
		AttemptID:           attempt.ID,
		UserID:              u.ID,
		FoundationsRequired: decision.FoundationsRequired,
		Exemptions:          decision.Exemptions,
	}
	if p, err := s.st.GetPathBySlug(ctx, decision.RecommendedPath); err == nil {
		res.RecommendedPathID = sql.NullInt64{Int64: p.ID, Valid: true}
	}
	if _, err := s.st.SavePlacementResult(ctx, res); err != nil {
		http.Error(w, "could not save your result", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/placement/result", http.StatusSeeOther)
}

// handlePlacementResult shows the routing decision and the paths now open.
func (s *Server) handlePlacementResult(w http.ResponseWriter, r *http.Request) {
	s.renderPlacementResult(w, r, "")
}

// renderPlacementResult renders the result page, optionally with a message
// explaining why an action was refused.
func (s *Server) renderPlacementResult(w http.ResponseWriter, r *http.Request, flash string) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	res, err := s.st.LatestPlacement(ctx, u.ID)
	if err != nil || res == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	attempt, err := s.st.GetAttempt(ctx, res.AttemptID)
	if err != nil {
		http.Error(w, "could not load your result", http.StatusInternalServerError)
		return
	}

	// Per-topic breakdown in a stable order, flagged for which count toward the
	// foundations gate — the learner should see exactly what decided their route.
	type topicRow struct {
		Topic, Label string
		Score        int
		Core         bool
	}
	var rows []topicRow
	for _, t := range placement.AllTopics {
		rows = append(rows, topicRow{
			Topic: t, Label: placement.TopicLabel(t),
			Score: attempt.TopicScores[t], Core: placement.IsCore(t),
		})
	}

	// Exemptions are stored as course slugs; show learners the course names.
	type exemptRow struct{ Slug, Title, Emoji string }
	var exemptions []exemptRow
	for _, slug := range res.Exemptions {
		row := exemptRow{Slug: slug, Title: slug}
		if c, err := s.st.GetCourseBySlug(ctx, slug); err == nil {
			row.Title, row.Emoji = c.Title, c.Emoji
		}
		exemptions = append(exemptions, row)
	}

	paths, _ := s.st.ListPaths(ctx, false)
	var recommended *store.Path
	for i := range paths {
		if res.RecommendedPathID.Valid && paths[i].ID == res.RecommendedPathID.Int64 {
			recommended = &paths[i]
		}
	}
	enr, _ := s.st.LiveEnrollment(ctx, u.ID)

	s.render(w, r, "placement_result.html", ViewData{Title: "Your placement", Flash: flash, Data: map[string]any{
		"bodyClass":   "placement-dark",
		"result":      res,
		"attempt":     attempt,
		"topics":      rows,
		"coreScore":   placement.MeanCore(attempt.TopicScores),
		"threshold":   placement.FoundationsThreshold,
		"paths":       paths,
		"recommended": recommended,
		"exemptions":  exemptions,
		"enrollment":  enr,
		"review":      nil, // review is a separate, opt-in page
	}})
}

// handlePlacementReview shows the graded paper with explanations. Available only
// once submitted, since it contains the answers.
func (s *Server) handlePlacementReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	res, err := s.st.LatestPlacement(ctx, u.ID)
	if err != nil || res == nil {
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}
	attempt, err := s.st.GetAttempt(ctx, res.AttemptID)
	if err != nil || attempt.UserID != u.ID || !attempt.Submitted() {
		s.notFound(w, r)
		return
	}
	review, err := s.st.ReviewAttempt(ctx, attempt.ID)
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

// handleEnroll commits a learner to a path.
//
// This is the gate. Two rules are enforced here rather than in the template,
// because a hidden button is not a control:
//
//  1. No enrollment without a completed diagnostic.
//  2. If the diagnostic said foundations are required, only that path is
//     selectable until it is completed.
func (s *Server) handleEnroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.CurrentUser(ctx)

	res, err := s.st.LatestPlacement(ctx, u.ID)
	if err != nil {
		http.Error(w, "could not check your placement", http.StatusInternalServerError)
		return
	}
	if res == nil {
		// Rule 1.
		http.Redirect(w, r, "/placement", http.StatusSeeOther)
		return
	}

	slug := r.FormValue("path")
	p, err := s.st.GetPathBySlug(ctx, slug)
	if err != nil || !p.Published {
		s.notFound(w, r)
		return
	}

	// Rule 2. Refused with an explanation rather than a bare redirect, so the
	// learner understands why the choice was not honoured.
	if res.FoundationsRequired && p.Slug != placement.PathFoundations {
		s.renderPlacementResult(w, r,
			"Computer Science Foundations comes first — the diagnostic showed gaps that the other paths build on.")
		return
	}

	enr, err := s.st.EnsureEnrollment(ctx, u.ID)
	if err != nil {
		http.Error(w, "could not enroll you", http.StatusInternalServerError)
		return
	}
	// Only an unplaced or paused enrollment can be (re)placed; an active one is
	// already committed and changing paths mid-cohort is a phase 3 concern.
	if enr.State != store.EnrollUnplaced && enr.State != store.EnrollPaused {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	if err := s.st.PlaceEnrollment(ctx, enr.ID, p.ID, "placement diagnostic"); err != nil {
		http.Error(w, "could not enroll you", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
