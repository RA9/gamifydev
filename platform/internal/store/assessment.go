package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Placement diagnostic persistence: the item bank, a sitting's lifecycle, and
// the routing decision it produces.
//
// Answers never leave the server. AttemptItem has no answer field, while each
// sitting snapshots the complete drawn item (including its answer) in
// assessment_attempt_items. Reseeding the bank therefore cannot change a live
// paper, its score, or its review history.

// Assessment is a named item bank.
type Assessment struct {
	ID         int64
	Slug       string
	Title      string
	Kind       string
	TimeLimitS int
	PerTopic   int
}

// Item is a bank item, including its answer. Server-side only.
type Item struct {
	ID          int64
	Topic       string
	Difficulty  int
	Prompt      string
	Code        string
	Lang        string
	Options     []string
	Answer      int
	Explanation string
}

// AttemptItem is one question as presented to a learner. It deliberately has no
// answer field.
type AttemptItem struct {
	ID       int64
	Sort     int
	Topic    string
	Prompt   string
	Code     string
	Lang     string
	Options  []string
	Response sql.NullInt64
}

// Attempt is one sitting. By is the canonical owner; UserID and GuestID are
// retained for compatibility with existing user-only callers.
type Attempt struct {
	ID           int64
	UserID       int64
	GuestID      int64
	By           Solver
	AssessmentID int64
	StartedAt    string
	ExpiresAt    string
	SubmittedAt  sql.NullString
	Score        sql.NullInt64
	TopicScores  map[string]int
	Abandoned    bool
}

// Submitted reports whether the attempt has been scored.
func (a Attempt) Submitted() bool { return a.SubmittedAt.Valid }

// OwnedBy reports whether this attempt belongs to by.
func (a Attempt) OwnedBy(by Solver) bool { return by.valid() && a.By == by }

// GetAssessment loads a published assessment by slug.
func (s *Store) GetAssessment(ctx context.Context, slug string) (*Assessment, error) {
	var a Assessment
	err := s.db.QueryRowContext(ctx,
		`SELECT id, slug, title, kind, time_limit_s, per_topic
		 FROM assessments WHERE slug = ? AND published = 1`, slug).
		Scan(&a.ID, &a.Slug, &a.Title, &a.Kind, &a.TimeLimitS, &a.PerTopic)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &a, err
}

// UpsertAssessment creates or updates a bank definition (used by the seeder).
func (s *Store) UpsertAssessment(ctx context.Context, a Assessment) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO assessments (slug, title, kind, time_limit_s, per_topic)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title=excluded.title, kind=excluded.kind,
			time_limit_s=excluded.time_limit_s, per_topic=excluded.per_topic
		RETURNING id`,
		a.Slug, a.Title, a.Kind, a.TimeLimitS, a.PerTopic).Scan(&id)
	return id, err
}

// ReplaceItems swaps an assessment's bank for the given items. Attempt item rows
// contain immutable snapshots, so replacing or deleting bank rows cannot alter
// an in-progress or historical sitting.
func (s *Store) ReplaceItems(ctx context.Context, assessmentID int64, items []Item) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM assessment_items WHERE assessment_id = ?`, assessmentID); err != nil {
		return err
	}
	for _, it := range items {
		opts, err := json.Marshal(it.Options)
		if err != nil {
			return err
		}
		if it.Answer < 0 || it.Answer >= len(it.Options) {
			return fmt.Errorf("item %q: answer index %d out of range", it.Prompt, it.Answer)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO assessment_items
				(assessment_id, topic, difficulty, prompt, code, lang, options, answer, explanation)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			assessmentID, it.Topic, it.Difficulty, it.Prompt, it.Code, it.Lang,
			string(opts), it.Answer, it.Explanation); err != nil {
			return err
		}
	}
	return tx.Commit()
}

const attemptCols = `id, user_id, guest_id, assessment_id, started_at, expires_at,
	submitted_at, score, topic_scores, abandoned`

// LiveAttempt is the user-only compatibility wrapper for LiveAttemptFor.
func (s *Store) LiveAttempt(ctx context.Context, userID int64) (*Attempt, error) {
	return s.LiveAttemptFor(ctx, Solver{UserID: userID})
}

// LiveAttemptFor returns the owner's in-progress attempt, or nil.
func (s *Store) LiveAttemptFor(ctx context.Context, by Solver) (*Attempt, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	userID, guestID := by.cols()
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT `+attemptCols+`
		FROM assessment_attempts
		WHERE (user_id = ? OR guest_id = ?)
		  AND submitted_at IS NULL AND abandoned = 0
		  AND datetime(expires_at) > datetime('now')
		ORDER BY datetime(started_at) DESC, id DESC LIMIT 1`, userID, guestID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

// LatestSubmittedAttempt is the user-only compatibility wrapper.
func (s *Store) LatestSubmittedAttempt(ctx context.Context, userID int64) (*Attempt, error) {
	return s.LatestSubmittedAttemptFor(ctx, Solver{UserID: userID})
}

// LatestSubmittedAttemptFor returns the owner's most recent scored attempt.
func (s *Store) LatestSubmittedAttemptFor(ctx context.Context, by Solver) (*Attempt, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	userID, guestID := by.cols()
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT `+attemptCols+`
		FROM assessment_attempts
		WHERE (user_id = ? OR guest_id = ?) AND submitted_at IS NOT NULL
		ORDER BY datetime(submitted_at) DESC, id DESC LIMIT 1`, userID, guestID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

func (s *Store) scanAttempt(row interface{ Scan(...any) error }) (*Attempt, error) {
	var a Attempt
	var userID, guestID sql.NullInt64
	var topics string
	err := row.Scan(&a.ID, &userID, &guestID, &a.AssessmentID, &a.StartedAt, &a.ExpiresAt,
		&a.SubmittedAt, &a.Score, &topics, &a.Abandoned)
	if err != nil {
		return nil, err
	}
	a.UserID, a.GuestID = userID.Int64, guestID.Int64
	a.By = Solver{UserID: a.UserID, GuestID: a.GuestID}
	a.TopicScores = map[string]int{}
	_ = json.Unmarshal([]byte(topics), &a.TopicScores)
	return &a, nil
}

const (
	PlacementRetryDelay  = 7 * 24 * time.Hour
	PlacementRetryWindow = 30 * 24 * time.Hour
	PlacementMaxAttempts = 3
)

var (
	// ErrAttemptInProgress prevents rerolling a live paper.
	ErrAttemptInProgress = errors.New("assessment attempt already in progress")
	// ErrPlacementRetryTooSoon means seven days have not elapsed since a failure.
	ErrPlacementRetryTooSoon = errors.New("placement retry is not available yet")
	// ErrPlacementAttemptLimit means three sittings have started in 30 days.
	ErrPlacementAttemptLimit = errors.New("placement attempt limit reached")
	// ErrPlacementAlreadyPassed means the owner already has a passing result.
	ErrPlacementAlreadyPassed = errors.New("placement already passed")
)

// PlacementRetryStatus explains whether an owner may start another placement.
type PlacementRetryStatus struct {
	Allowed        bool
	AttemptsLast30 int
	RetryAt        string
	Reason         error
}

type assessmentQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func placementRetryStatus(ctx context.Context, q assessmentQuerier, by Solver) (PlacementRetryStatus, error) {
	var status PlacementRetryStatus
	userID, guestID := by.cols()
	if err := q.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM assessment_attempts a
		JOIN assessments bank ON bank.id = a.assessment_id
		WHERE bank.kind = 'placement'
		  AND (a.user_id = ? OR a.guest_id = ?)
		  AND datetime(a.started_at) >= datetime('now', '-30 days')`, userID, guestID).
		Scan(&status.AttemptsLast30); err != nil {
		return status, err
	}

	var passed int
	var retryAt string
	err := q.QueryRowContext(ctx, `
		SELECT passed, datetime(created_at, '+7 days')
		FROM placement_results
		WHERE user_id = ? OR guest_id = ?
		ORDER BY datetime(created_at) DESC, id DESC LIMIT 1`, userID, guestID).
		Scan(&passed, &retryAt)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return status, err
	}
	if err == nil && passed == 1 {
		status.Reason = ErrPlacementAlreadyPassed
		return status, nil
	}
	if status.AttemptsLast30 >= PlacementMaxAttempts {
		status.Reason = ErrPlacementAttemptLimit
		return status, nil
	}
	if err == nil {
		var waiting int
		if err := q.QueryRowContext(ctx,
			`SELECT CASE WHEN datetime(?) > datetime('now') THEN 1 ELSE 0 END`, retryAt).
			Scan(&waiting); err != nil {
			return status, err
		}
		if waiting == 1 {
			status.RetryAt = retryAt
			status.Reason = ErrPlacementRetryTooSoon
			return status, nil
		}
	}
	status.Allowed = true
	return status, nil
}

// PlacementRetryFor reports the persisted retry state for an owner.
func (s *Store) PlacementRetryFor(ctx context.Context, by Solver) (PlacementRetryStatus, error) {
	if !by.valid() {
		return PlacementRetryStatus{}, ErrNoSolver
	}
	return placementRetryStatus(ctx, s.db, by)
}

// StartAttempt is the user-only compatibility wrapper for StartAttemptFor.
func (s *Store) StartAttempt(ctx context.Context, userID int64, a *Assessment) (*Attempt, error) {
	return s.StartAttemptFor(ctx, Solver{UserID: userID}, a)
}

// StartAttemptFor draws a fresh paper and opens a sitting for exactly one owner.
// Placement attempts enforce a seven-day wait after failure and a rolling limit
// of three starts per 30 days.
func (s *Store) StartAttemptFor(ctx context.Context, by Solver, a *Assessment) (*Attempt, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	if a == nil {
		return nil, ErrNotFound
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	if err := solverExists(ctx, tx, by); err != nil {
		return nil, err
	}
	userID, guestID := by.cols()
	var live int
	err = tx.QueryRowContext(ctx, `
		SELECT 1 FROM assessment_attempts
		WHERE (user_id = ? OR guest_id = ?)
		  AND submitted_at IS NULL AND abandoned = 0
		  AND datetime(expires_at) > datetime('now') LIMIT 1`, userID, guestID).Scan(&live)
	if err == nil {
		return nil, ErrAttemptInProgress
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if a.Kind == "placement" {
		status, err := placementRetryStatus(ctx, tx, by)
		if err != nil {
			return nil, err
		}
		if !status.Allowed {
			return nil, status.Reason
		}
	}

	expires := time.Now().UTC().Add(time.Duration(a.TimeLimitS) * time.Second).Format("2006-01-02 15:04:05")
	var attemptID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO assessment_attempts (user_id, guest_id, assessment_id, expires_at)
		VALUES (?, ?, ?, ?) RETURNING id`, userID, guestID, a.ID, expires).Scan(&attemptID); err != nil {
		return nil, err
	}

	// Stratified sample: per_topic random items from each topic. The selected bank
	// id remains the form key, while every value needed to render and grade is
	// copied into the sitting.
	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM (
			SELECT id, topic,
			       ROW_NUMBER() OVER (PARTITION BY topic ORDER BY RANDOM()) AS rn
			FROM assessment_items
			WHERE assessment_id = ? AND retired = 0
		)
		WHERE rn <= ?
		ORDER BY topic, rn`, a.ID, a.PerTopic)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("assessment has no items")
	}
	for i, id := range ids {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO assessment_attempt_items
				(attempt_id, item_id, sort, topic, difficulty, prompt, code, lang, options, answer, explanation)
			SELECT ?, id, ?, topic, difficulty, prompt, code, lang, options, answer, explanation
			FROM assessment_items WHERE id = ?`, attemptID, i, id); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetAttemptFor(ctx, by, attemptID)
}

func solverExists(ctx context.Context, q assessmentQuerier, by Solver) error {
	var one int
	var err error
	if by.UserID != 0 {
		err = q.QueryRowContext(ctx, `SELECT 1 FROM users WHERE id = ?`, by.UserID).Scan(&one)
	} else {
		err = q.QueryRowContext(ctx,
			`SELECT 1 FROM guest_sessions WHERE id = ? AND claimed_by IS NULL`, by.GuestID).Scan(&one)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// GetAttempt loads an attempt by id. New request paths should use GetAttemptFor;
// this unscoped form remains for existing trusted user-only call sites.
func (s *Store) GetAttempt(ctx context.Context, id int64) (*Attempt, error) {
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx,
		`SELECT `+attemptCols+` FROM assessment_attempts WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

// GetAttemptFor loads an attempt only when it belongs to by.
func (s *Store) GetAttemptFor(ctx context.Context, by Solver, id int64) (*Attempt, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	userID, guestID := by.cols()
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT `+attemptCols+` FROM assessment_attempts
		WHERE id = ? AND (user_id = ? OR guest_id = ?)`, id, userID, guestID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

// AttemptItems returns the snapshotted paper without answers. New request paths
// should use AttemptItemsFor.
func (s *Store) AttemptItems(ctx context.Context, attemptID int64) ([]AttemptItem, error) {
	return s.attemptItems(ctx, attemptID)
}

// AttemptItemsFor returns the paper only to its owner.
func (s *Store) AttemptItemsFor(ctx context.Context, by Solver, attemptID int64) ([]AttemptItem, error) {
	if _, err := s.GetAttemptFor(ctx, by, attemptID); err != nil {
		return nil, err
	}
	return s.attemptItems(ctx, attemptID)
}

func (s *Store) attemptItems(ctx context.Context, attemptID int64) ([]AttemptItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT item_id, sort, topic, prompt, code, lang, options, response
		FROM assessment_attempt_items
		WHERE attempt_id = ? ORDER BY sort`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AttemptItem
	for rows.Next() {
		var it AttemptItem
		var opts string
		if err := rows.Scan(&it.ID, &it.Sort, &it.Topic, &it.Prompt, &it.Code, &it.Lang, &opts, &it.Response); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(opts), &it.Options)
		out = append(out, it)
	}
	return out, rows.Err()
}

// ErrAttemptExpired is returned when a submission arrives after the time limit.
var ErrAttemptExpired = errors.New("attempt expired")

// ErrAlreadySubmitted guards against double submission.
var ErrAlreadySubmitted = errors.New("attempt already submitted")

// ScoreAttempt is the compatibility form for trusted call sites that already
// obtained the attempt through an owner-scoped lookup.
func (s *Store) ScoreAttempt(ctx context.Context, attemptID int64, responses map[int64]int) (*Attempt, error) {
	return s.scoreAttempt(ctx, Solver{}, false, attemptID, responses)
}

// ScoreAttemptFor grades an attempt only when it belongs to by.
func (s *Store) ScoreAttemptFor(ctx context.Context, by Solver, attemptID int64, responses map[int64]int) (*Attempt, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	return s.scoreAttempt(ctx, by, true, attemptID, responses)
}

func (s *Store) scoreAttempt(ctx context.Context, by Solver, scoped bool, attemptID int64, responses map[int64]int) (*Attempt, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	var submitted sql.NullString
	var expired int
	if scoped {
		userID, guestID := by.cols()
		err = tx.QueryRowContext(ctx, `
			SELECT submitted_at, CASE WHEN datetime(expires_at) <= datetime('now') THEN 1 ELSE 0 END
			FROM assessment_attempts
			WHERE id = ? AND (user_id = ? OR guest_id = ?)`, attemptID, userID, guestID).
			Scan(&submitted, &expired)
	} else {
		err = tx.QueryRowContext(ctx, `
			SELECT submitted_at, CASE WHEN datetime(expires_at) <= datetime('now') THEN 1 ELSE 0 END
			FROM assessment_attempts WHERE id = ?`, attemptID).Scan(&submitted, &expired)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if submitted.Valid {
		return nil, ErrAlreadySubmitted
	}
	if expired == 1 {
		if _, err := tx.ExecContext(ctx,
			`UPDATE assessment_attempts SET abandoned = 1 WHERE id = ?`, attemptID); err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, ErrAttemptExpired
	}

	type graded struct {
		id      int64
		topic   string
		correct bool
		resp    *int
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT item_id, topic, answer
		FROM assessment_attempt_items
		WHERE attempt_id = ?`, attemptID)
	if err != nil {
		return nil, err
	}
	var all []graded
	for rows.Next() {
		var id int64
		var topic string
		var answer int
		if err := rows.Scan(&id, &topic, &answer); err != nil {
			rows.Close()
			return nil, err
		}
		g := graded{id: id, topic: topic}
		if r, ok := responses[id]; ok {
			rr := r
			g.resp = &rr
			g.correct = r == answer
		}
		all = append(all, g)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	perTopicTotal := map[string]int{}
	perTopicRight := map[string]int{}
	right := 0
	for _, g := range all {
		perTopicTotal[g.topic]++
		if g.correct {
			right++
			perTopicRight[g.topic]++
		}
		correct := 0
		if g.correct {
			correct = 1
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE assessment_attempt_items SET response = ?, correct = ?
			WHERE attempt_id = ? AND item_id = ?`, g.resp, correct, attemptID, g.id); err != nil {
			return nil, err
		}
	}

	topicScores := map[string]int{}
	for topic, total := range perTopicTotal {
		if total > 0 {
			topicScores[topic] = perTopicRight[topic] * 100 / total
		}
	}
	overall := 0
	if len(all) > 0 {
		overall = right * 100 / len(all)
	}
	ts, _ := json.Marshal(topicScores)
	if _, err := tx.ExecContext(ctx, `
		UPDATE assessment_attempts
		SET submitted_at = datetime('now'), score = ?, topic_scores = ?
		WHERE id = ?`, overall, string(ts), attemptID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	if scoped {
		return s.GetAttemptFor(ctx, by, attemptID)
	}
	return s.GetAttempt(ctx, attemptID)
}

// PlacementResult is the stored routing decision. Passed is explicit because
// score-band interpretation belongs to placement policy, not persistence.
type PlacementResult struct {
	ID                  int64
	AttemptID           int64
	UserID              int64
	GuestID             int64
	By                  Solver
	RecommendedPathID   sql.NullInt64
	FoundationsRequired bool
	Passed              bool
	Exemptions          []string
	CreatedAt           string
}

// OwnedBy reports whether this result belongs to by.
func (r PlacementResult) OwnedBy(by Solver) bool { return by.valid() && r.By == by }

var ErrPlacementResultExists = errors.New("placement result already exists for attempt")

// SavePlacementResult is the user-only compatibility wrapper.
func (s *Store) SavePlacementResult(ctx context.Context, r PlacementResult) (int64, error) {
	by := r.By
	if !by.valid() {
		by = Solver{UserID: r.UserID, GuestID: r.GuestID}
	}
	return s.SavePlacementResultFor(ctx, by, r)
}

// SavePlacementResultFor records one routing decision for an owned attempt.
func (s *Store) SavePlacementResultFor(ctx context.Context, by Solver, r PlacementResult) (int64, error) {
	if !by.valid() {
		return 0, ErrNoSolver
	}
	if r.Exemptions == nil {
		r.Exemptions = []string{}
	}
	ex, err := json.Marshal(r.Exemptions)
	if err != nil {
		return 0, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck
	userID, guestID := by.cols()
	var one int
	if err := tx.QueryRowContext(ctx, `
		SELECT 1 FROM assessment_attempts
		WHERE id = ? AND submitted_at IS NOT NULL
		  AND (user_id = ? OR guest_id = ?)`, r.AttemptID, userID, guestID).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO placement_results
			(attempt_id, user_id, guest_id, recommended_path_id,
			 foundations_required, passed, exemptions)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		r.AttemptID, userID, guestID, r.RecommendedPathID,
		boolToInt(r.FoundationsRequired), boolToInt(r.Passed), string(ex)).Scan(&id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return 0, ErrPlacementResultExists
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

// LatestPlacement is the user-only compatibility wrapper.
func (s *Store) LatestPlacement(ctx context.Context, userID int64) (*PlacementResult, error) {
	return s.LatestPlacementFor(ctx, Solver{UserID: userID})
}

// LatestPlacementFor returns the owner's most recent routing decision.
func (s *Store) LatestPlacementFor(ctx context.Context, by Solver) (*PlacementResult, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	userID, guestID := by.cols()
	return scanPlacementResult(s.db.QueryRowContext(ctx, `
		SELECT id, attempt_id, user_id, guest_id, recommended_path_id,
		       foundations_required, passed, exemptions, created_at
		FROM placement_results
		WHERE user_id = ? OR guest_id = ?
		ORDER BY datetime(created_at) DESC, id DESC LIMIT 1`, userID, guestID), true)
}

func scanPlacementResult(row interface{ Scan(...any) error }, nilWhenMissing bool) (*PlacementResult, error) {
	var r PlacementResult
	var userID, guestID sql.NullInt64
	var ex string
	var req, passed int
	err := row.Scan(&r.ID, &r.AttemptID, &userID, &guestID, &r.RecommendedPathID,
		&req, &passed, &ex, &r.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		if nilWhenMissing {
			return nil, nil
		}
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	r.UserID, r.GuestID = userID.Int64, guestID.Int64
	r.By = Solver{UserID: r.UserID, GuestID: r.GuestID}
	r.FoundationsRequired = req == 1
	r.Passed = passed == 1
	_ = json.Unmarshal([]byte(ex), &r.Exemptions)
	return &r, nil
}

// AttemptReview is one graded question, for the result page.
type AttemptReview struct {
	Topic       string
	Prompt      string
	Options     []string
	Answer      int
	Response    sql.NullInt64
	Correct     bool
	Explanation string
}

// ReviewAttempt returns the graded snapshotted paper. New request paths should
// use ReviewAttemptFor.
func (s *Store) ReviewAttempt(ctx context.Context, attemptID int64) ([]AttemptReview, error) {
	return s.reviewAttempt(ctx, Solver{}, false, attemptID)
}

// ReviewAttemptFor returns answers only for a submitted attempt owned by by.
func (s *Store) ReviewAttemptFor(ctx context.Context, by Solver, attemptID int64) ([]AttemptReview, error) {
	if !by.valid() {
		return nil, ErrNoSolver
	}
	return s.reviewAttempt(ctx, by, true, attemptID)
}

func (s *Store) reviewAttempt(ctx context.Context, by Solver, scoped bool, attemptID int64) ([]AttemptReview, error) {
	var one int
	var err error
	if scoped {
		userID, guestID := by.cols()
		err = s.db.QueryRowContext(ctx, `
			SELECT 1 FROM assessment_attempts
			WHERE id = ? AND submitted_at IS NOT NULL
			  AND (user_id = ? OR guest_id = ?)`, attemptID, userID, guestID).Scan(&one)
	} else {
		err = s.db.QueryRowContext(ctx,
			`SELECT 1 FROM assessment_attempts WHERE id = ? AND submitted_at IS NOT NULL`, attemptID).Scan(&one)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT topic, prompt, options, answer, response, COALESCE(correct, 0), explanation
		FROM assessment_attempt_items
		WHERE attempt_id = ? ORDER BY sort`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AttemptReview
	for rows.Next() {
		var r AttemptReview
		var opts string
		var correct int
		if err := rows.Scan(&r.Topic, &r.Prompt, &opts, &r.Answer, &r.Response, &correct, &r.Explanation); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(opts), &r.Options)
		r.Correct = correct == 1
		out = append(out, r)
	}
	return out, rows.Err()
}

var (
	ErrGuestAlreadyClaimed = errors.New("guest session already claimed")
	ErrPlacementNotPassed  = errors.New("guest has no passing placement result")
)

// CreateLearnerFromGuestPlacement creates a learner account and atomically
// transfers all of an unclaimed passing guest's placement work and problem
// submissions. The role is intentionally not an argument: this public signup
// path can never mint an admin or grader.
func (s *Store) CreateLearnerFromGuestPlacement(ctx context.Context, guestID int64, email, passwordHash, name string) (*User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	var claimedBy sql.NullInt64
	if err := tx.QueryRowContext(ctx,
		`SELECT claimed_by FROM guest_sessions WHERE id = ?`, guestID).Scan(&claimedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if claimedBy.Valid {
		return nil, ErrGuestAlreadyClaimed
	}

	var passing int
	err = tx.QueryRowContext(ctx, `
		SELECT 1
		FROM placement_results pr
		JOIN assessment_attempts a ON a.id = pr.attempt_id
		WHERE pr.guest_id = ? AND pr.user_id IS NULL AND pr.passed = 1
		  AND a.guest_id = ? AND a.user_id IS NULL AND a.submitted_at IS NOT NULL
		LIMIT 1`, guestID, guestID).Scan(&passing)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPlacementNotPassed
	}
	if err != nil {
		return nil, err
	}

	u, err := s.scanUser(tx.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, name, role) VALUES (?, ?, ?, 'learner')
		 RETURNING `+userCols,
		strings.ToLower(strings.TrimSpace(email)), passwordHash, name))
	if err != nil {
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE guest_sessions SET claimed_by = ?
		WHERE id = ? AND claimed_by IS NULL`, u.ID, guestID)
	if err != nil {
		return nil, err
	}
	claimed, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if claimed != 1 {
		return nil, ErrGuestAlreadyClaimed
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE assessment_attempts SET user_id = ?, guest_id = NULL
		WHERE guest_id = ? AND user_id IS NULL`, u.ID, guestID); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE placement_results SET user_id = ?, guest_id = NULL
		WHERE guest_id = ? AND user_id IS NULL`, u.ID, guestID); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE problem_submissions SET user_id = ?, guest_id = NULL
		WHERE guest_id = ? AND user_id IS NULL`, u.ID, guestID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

// ExpireAbandonedAttempts marks unsubmitted attempts past their deadline.
func (s *Store) ExpireAbandonedAttempts(ctx context.Context) (int, error) {
	res, err := s.db.ExecContext(ctx, `
		UPDATE assessment_attempts SET abandoned = 1
		WHERE submitted_at IS NULL AND abandoned = 0
		  AND datetime(expires_at) <= datetime('now')`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
