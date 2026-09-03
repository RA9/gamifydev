package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Placement diagnostic persistence: the item bank, a sitting's lifecycle, and
// the routing decision it produces.
//
// Two properties this file is responsible for:
//
//   - Answers never leave the server. Item rows carry the correct index, but the
//     type used to render a paper (AttemptItem) has no answer field at all, so
//     it cannot leak into a template by accident.
//   - The paper is fixed at start. The drawn items are recorded in
//     assessment_attempt_items, and scoring reads that table rather than trusting
//     anything the client posts back.

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

// Attempt is one sitting.
type Attempt struct {
	ID           int64
	UserID       int64
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

// ReplaceItems swaps an assessment's bank for the given items.
//
// Items are replaced rather than retired because the bank is seeded content; a
// learner's finished attempt keeps its own copy of what was asked via
// assessment_attempt_items, so history is not disturbed.
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

// LiveAttempt returns the user's in-progress attempt (started, not submitted,
// not expired), or nil.
func (s *Store) LiveAttempt(ctx context.Context, userID int64) (*Attempt, error) {
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT id, user_id, assessment_id, started_at, expires_at, submitted_at, score, topic_scores, abandoned
		FROM assessment_attempts
		WHERE user_id = ? AND submitted_at IS NULL AND abandoned = 0
		  AND datetime(expires_at) > datetime('now')
		ORDER BY started_at DESC LIMIT 1`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

// LatestSubmittedAttempt returns the user's most recent scored attempt, or nil.
func (s *Store) LatestSubmittedAttempt(ctx context.Context, userID int64) (*Attempt, error) {
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT id, user_id, assessment_id, started_at, expires_at, submitted_at, score, topic_scores, abandoned
		FROM assessment_attempts
		WHERE user_id = ? AND submitted_at IS NOT NULL
		ORDER BY submitted_at DESC LIMIT 1`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

func (s *Store) scanAttempt(row interface{ Scan(...any) error }) (*Attempt, error) {
	var a Attempt
	var topics string
	err := row.Scan(&a.ID, &a.UserID, &a.AssessmentID, &a.StartedAt, &a.ExpiresAt,
		&a.SubmittedAt, &a.Score, &topics, &a.Abandoned)
	if err != nil {
		return nil, err
	}
	a.TopicScores = map[string]int{}
	_ = json.Unmarshal([]byte(topics), &a.TopicScores)
	return &a, nil
}

// StartAttempt draws a fresh paper and opens a sitting.
//
// Items are sampled per topic rather than served whole and in order: two
// learners rarely see the same paper, which is the practical mitigation against
// a leaked bank. `perTopic` items are drawn from each topic present in the bank.
func (s *Store) StartAttempt(ctx context.Context, userID int64, a *Assessment) (*Attempt, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	expires := time.Now().UTC().Add(time.Duration(a.TimeLimitS) * time.Second).Format("2006-01-02 15:04:05")
	var attemptID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO assessment_attempts (user_id, assessment_id, expires_at)
		VALUES (?, ?, ?) RETURNING id`, userID, a.ID, expires).Scan(&attemptID); err != nil {
		return nil, err
	}

	// Stratified sample: `per_topic` random items from each topic. Ordering the
	// paper by topic keeps related questions together, which reads better than
	// a fully shuffled paper.
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
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO assessment_attempt_items (attempt_id, item_id, sort) VALUES (?, ?, ?)`,
			attemptID, id, i); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetAttempt(ctx, attemptID)
}

// GetAttempt loads an attempt by id.
func (s *Store) GetAttempt(ctx context.Context, id int64) (*Attempt, error) {
	a, err := s.scanAttempt(s.db.QueryRowContext(ctx, `
		SELECT id, user_id, assessment_id, started_at, expires_at, submitted_at, score, topic_scores, abandoned
		FROM assessment_attempts WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

// AttemptItems returns the paper drawn for an attempt, without answers.
func (s *Store) AttemptItems(ctx context.Context, attemptID int64) ([]AttemptItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT i.id, ai.sort, i.topic, i.prompt, i.code, i.lang, i.options, ai.response
		FROM assessment_attempt_items ai
		JOIN assessment_items i ON i.id = ai.item_id
		WHERE ai.attempt_id = ? ORDER BY ai.sort`, attemptID)
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

// ScoreAttempt grades a submission and records per-topic percentages.
//
// `responses` maps item id to the chosen option index; items absent from the map
// are treated as unanswered and wrong. Only items actually drawn for this
// attempt are considered, so a client cannot smuggle in extra answers.
//
// The time limit is enforced here, against the database clock.
func (s *Store) ScoreAttempt(ctx context.Context, attemptID int64, responses map[int64]int) (*Attempt, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	var submitted sql.NullString
	var expired int
	if err := tx.QueryRowContext(ctx, `
		SELECT submitted_at, CASE WHEN datetime(expires_at) <= datetime('now') THEN 1 ELSE 0 END
		FROM assessment_attempts WHERE id = ?`, attemptID).Scan(&submitted, &expired); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
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

	// Grade against the stored paper and stored answers.
	rows, err := tx.QueryContext(ctx, `
		SELECT i.id, i.topic, i.answer
		FROM assessment_attempt_items ai
		JOIN assessment_items i ON i.id = ai.item_id
		WHERE ai.attempt_id = ?`, attemptID)
	if err != nil {
		return nil, err
	}
	type graded struct {
		id      int64
		topic   string
		correct bool
		resp    (*int)
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
		c := 0
		if g.correct {
			c = 1
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE assessment_attempt_items SET response = ?, correct = ? WHERE attempt_id = ? AND item_id = ?`,
			g.resp, c, attemptID, g.id); err != nil {
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
	return s.GetAttempt(ctx, attemptID)
}

// PlacementResult is the stored routing decision.
type PlacementResult struct {
	ID                  int64
	AttemptID           int64
	UserID              int64
	RecommendedPathID   sql.NullInt64
	FoundationsRequired bool
	Exemptions          []string
	CreatedAt           string
}

// SavePlacementResult records a routing decision for an attempt.
func (s *Store) SavePlacementResult(ctx context.Context, r PlacementResult) (int64, error) {
	// A nil slice marshals to `null`; this column is documented as a JSON array,
	// so normalise the empty case rather than storing two shapes.
	if r.Exemptions == nil {
		r.Exemptions = []string{}
	}
	ex, _ := json.Marshal(r.Exemptions)
	req := 0
	if r.FoundationsRequired {
		req = 1
	}
	var id int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO placement_results
			(attempt_id, user_id, recommended_path_id, foundations_required, exemptions)
		VALUES (?, ?, ?, ?, ?) RETURNING id`,
		r.AttemptID, r.UserID, r.RecommendedPathID, req, string(ex)).Scan(&id)
	return id, err
}

// LatestPlacement returns the user's most recent routing decision, or nil if
// they have never completed the diagnostic.
func (s *Store) LatestPlacement(ctx context.Context, userID int64) (*PlacementResult, error) {
	var r PlacementResult
	var ex string
	var req int
	err := s.db.QueryRowContext(ctx, `
		SELECT id, attempt_id, user_id, recommended_path_id, foundations_required, exemptions, created_at
		FROM placement_results WHERE user_id = ? ORDER BY created_at DESC, id DESC LIMIT 1`, userID).
		Scan(&r.ID, &r.AttemptID, &r.UserID, &r.RecommendedPathID, &req, &ex, &r.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.FoundationsRequired = req == 1
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

// ReviewAttempt returns the graded paper. Only safe to call once the attempt is
// submitted — answers are included.
func (s *Store) ReviewAttempt(ctx context.Context, attemptID int64) ([]AttemptReview, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT i.topic, i.prompt, i.options, i.answer, ai.response, COALESCE(ai.correct, 0), i.explanation
		FROM assessment_attempt_items ai
		JOIN assessment_items i ON i.id = ai.item_id
		WHERE ai.attempt_id = ? ORDER BY ai.sort`, attemptID)
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

// ExpireAbandonedAttempts marks unsubmitted attempts past their deadline. Run by
// the placement:expire job so the "you have a sitting in progress" state cannot
// stick forever.
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
