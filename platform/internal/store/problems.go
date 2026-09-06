package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Verdicts a judged submission can carry.
const (
	VerdictPending      = "pending"
	VerdictAccepted     = "accepted"
	VerdictWrongAnswer  = "wrong_answer"
	VerdictTimeLimit    = "time_limit"
	VerdictCompileError = "compile_error"
	VerdictRuntimeError = "runtime_error"
	VerdictError        = "error"
)

type Problem struct {
	ID              int64
	Slug            string
	Title           string
	Difficulty      string
	Topic           string
	Statement       string
	TimeLimitMs     int
	Sort            int
	Published       bool
	WorkloadMinutes int
	Mode            string
	Language        string
	CourseSort      int  // populated for course-integrated problems
	Required        bool // populated for course-integrated problems
	// Solved is filled by list queries that know who is asking.
	Solved bool
}

type ProblemSubmission struct {
	ID          int64
	ProblemID   int64
	Language    string
	Code        string
	Verdict     string
	Passed      int
	Total       int
	Output      string
	FailedLabel string
	RuntimeMs   int
	CreatedAt   string
	JudgedAt    sql.NullString
	// By is who wrote it. Carried on the row so a handler can prove the person
	// polling for a result is the one who submitted it — a submission holds
	// someone's code, and knowing its id is not permission to read it.
	By Solver
	// Problem fields the judge needs without a second query.
	ProblemSlug string
	TimeLimitMs int
}

// Solver identifies who is attempting a problem: a signed-in user, or a guest
// holding a cookie. Exactly one field is set.
//
// Guests are first-class here on purpose — being able to try the problems
// without an account is the point of the bank, and threading a nullable user id
// through every call site would make "no account" the exceptional path when it
// is the common one.
type Solver struct {
	UserID  int64
	GuestID int64
}

func (s Solver) valid() bool { return (s.UserID == 0) != (s.GuestID == 0) }

func (s Solver) cols() (userID, guestID any) {
	if s.UserID != 0 {
		return s.UserID, nil
	}
	return nil, s.GuestID
}

// Owns reports whether a submission belongs to this solver. A zero Solver owns
// nothing, so an unidentified visitor is refused rather than matched against
// rows whose author column is NULL.
func (s Solver) Owns(sub *ProblemSubmission) bool {
	return s.valid() && sub != nil && s == sub.By
}

// ErrNoSolver means neither a user nor a guest was identified — a programming
// error rather than something a visitor can cause.
var ErrNoSolver = errors.New("store: submission has no author")

// --- problems ---------------------------------------------------------------

// ListProblems returns the published bank in author order, flagging which the
// solver has already solved. A zero Solver simply leaves everything unsolved.
func (s *Store) ListProblems(ctx context.Context, by Solver) ([]Problem, error) {
	userID, guestID := by.cols()
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id, p.slug, p.title, p.difficulty, p.topic, p.time_limit_ms, p.sort,
		       p.workload_minutes, p.mode, p.language,
		       EXISTS (
		         SELECT 1 FROM problem_submissions ps
		         WHERE ps.problem_id = p.id AND ps.verdict = 'accepted'
		           AND (ps.user_id = ? OR ps.guest_id = ?)
		       ) AS solved
		FROM problems p
		WHERE p.published = 1
		ORDER BY p.sort, p.id`, userID, guestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.Topic,
			&p.TimeLimitMs, &p.Sort, &p.WorkloadMinutes, &p.Mode, &p.Language, &p.Solved); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetProblemBySlug(ctx context.Context, slug string) (*Problem, error) {
	var p Problem
	err := s.db.QueryRowContext(ctx, `
		SELECT id, slug, title, difficulty, topic, statement, time_limit_ms, sort,
		       published, workload_minutes, mode, language
		FROM problems WHERE slug = ? AND published = 1`, slug).
		Scan(&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.Topic, &p.Statement,
			&p.TimeLimitMs, &p.Sort, &p.Published, &p.WorkloadMinutes, &p.Mode, &p.Language)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &p, err
}

// Starters returns the languages a problem accepts, and their starter code.
func (s *Store) Starters(ctx context.Context, problemID int64) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT language, code FROM problem_starters WHERE problem_id = ? ORDER BY language`, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var lang, code string
		if err := rows.Scan(&lang, &code); err != nil {
			return nil, err
		}
		out[lang] = code
	}
	return out, rows.Err()
}

// ProblemTests returns a problem's tests as checks, so the judge can hand them
// straight to the same grader that runs checkpoint checks.
func (s *Store) ProblemTests(ctx context.Context, problemID int64) ([]Check, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, sort, label, test, stdin, hidden, points
		 FROM problem_tests WHERE problem_id = ? ORDER BY sort, id`, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Check
	for rows.Next() {
		var c Check
		var hidden int
		if err := rows.Scan(&c.ID, &c.Sort, &c.Label, &c.Test, &c.Stdin, &hidden, &c.Points); err != nil {
			return nil, err
		}
		c.Hidden = hidden == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpsertProblem writes a problem by slug and returns its id.
func (s *Store) UpsertProblem(ctx context.Context, p Problem) (int64, error) {
	if p.WorkloadMinutes <= 0 {
		p.WorkloadMinutes = 30
	}
	if p.Mode == "" {
		p.Mode = ModePractical
	}
	if p.Language == "" {
		p.Language = "python"
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO problems
			(slug, title, difficulty, topic, statement, time_limit_ms, sort, published,
			 workload_minutes, mode, language)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title=excluded.title, difficulty=excluded.difficulty, topic=excluded.topic,
			statement=excluded.statement, time_limit_ms=excluded.time_limit_ms,
			sort=excluded.sort, published=excluded.published,
			workload_minutes=excluded.workload_minutes, mode=excluded.mode,
			language=excluded.language`,
		p.Slug, p.Title, p.Difficulty, p.Topic, p.Statement, p.TimeLimitMs, p.Sort,
		boolToInt(p.Published), p.WorkloadMinutes, p.Mode, p.Language)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.db.QueryRowContext(ctx, `SELECT id FROM problems WHERE slug = ?`, p.Slug).Scan(&id)
	return id, err
}

// SetCourseProblemsBySlug replaces the ordered practice set integrated into a
// course. Problems remain available in the public bank; this table only gives
// them a curricular home and order.
func (s *Store) SetCourseProblemsBySlug(ctx context.Context, courseID int64, slugs []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `DELETE FROM course_problems WHERE course_id = ?`, courseID); err != nil {
		return err
	}
	for i, slug := range slugs {
		var problemID int64
		if err := tx.QueryRowContext(ctx,
			`SELECT id FROM problems WHERE slug = ? AND published = 1`, slug).Scan(&problemID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("course problem %q: %w", slug, ErrNotFound)
			}
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO course_problems (course_id, problem_id, sort, required)
			VALUES (?, ?, ?, 1)`, courseID, problemID, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// ListProblemsByCourse returns the ordered practice set embedded in a course.
func (s *Store) ListProblemsByCourse(ctx context.Context, courseID int64, by Solver) ([]Problem, error) {
	userID, guestID := by.cols()
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id, p.slug, p.title, p.difficulty, p.topic, p.time_limit_ms, p.sort,
		       p.workload_minutes, p.mode, p.language, cp.sort, cp.required,
		       EXISTS (
		         SELECT 1 FROM problem_submissions ps
		         WHERE ps.problem_id = p.id AND ps.verdict = 'accepted'
		           AND (ps.user_id = ? OR ps.guest_id = ?)
		       ) AS solved
		FROM course_problems cp JOIN problems p ON p.id = cp.problem_id
		WHERE cp.course_id = ? AND p.published = 1
		ORDER BY cp.sort, p.id`, userID, guestID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.Topic,
			&p.TimeLimitMs, &p.Sort, &p.WorkloadMinutes, &p.Mode, &p.Language,
			&p.CourseSort, &p.Required, &p.Solved); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ReplaceStarters and ReplaceProblemTests swap a problem's content wholesale,
// so re-seeding picks up an edit instead of stacking duplicates.
func (s *Store) ReplaceStarters(ctx context.Context, problemID int64, starters map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `DELETE FROM problem_starters WHERE problem_id = ?`, problemID); err != nil {
		return err
	}
	for lang, code := range starters {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO problem_starters (problem_id, language, code) VALUES (?, ?, ?)`,
			problemID, lang, code); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ReplaceProblemTests(ctx context.Context, problemID int64, tests []Check) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `DELETE FROM problem_tests WHERE problem_id = ?`, problemID); err != nil {
		return err
	}
	for i, t := range tests {
		hidden := 0
		if t.Hidden {
			hidden = 1
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO problem_tests (problem_id, sort, label, test, stdin, hidden, points)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			problemID, i, t.Label, t.Test, t.Stdin, hidden, t.Points); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// --- guests -----------------------------------------------------------------

// CreateGuestSession records a new anonymous visitor and returns its row id.
func (s *Store) CreateGuestSession(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO guest_sessions (token_hash) VALUES (?) RETURNING id`,
		HashResetToken(token)).Scan(&id)
	return id, err
}

// GuestByToken resolves a cookie to a guest row, refreshing last_seen_at. A
// claimed session is not returned: once its work belongs to an account, the
// cookie is spent.
func (s *Store) GuestByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`SELECT id FROM guest_sessions WHERE token_hash = ? AND claimed_by IS NULL`,
		HashResetToken(token)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	_, _ = s.db.ExecContext(ctx,
		`UPDATE guest_sessions SET last_seen_at = datetime('now') WHERE id = ?`, id)
	return id, nil
}

// ClaimGuestWork hands every guest-owned resource to an existing account.
// Problem submissions and placement attempts/results move together in one
// transaction, then the guest session is marked spent so the cookie cannot be
// replayed.
func (s *Store) ClaimGuestWork(ctx context.Context, guestID, userID int64) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck

	res, err := tx.ExecContext(ctx,
		`UPDATE problem_submissions SET user_id = ?, guest_id = NULL WHERE guest_id = ?`,
		userID, guestID)
	if err != nil {
		return 0, err
	}
	moved, _ := res.RowsAffected()
	if _, err := tx.ExecContext(ctx,
		`UPDATE assessment_attempts SET user_id = ?, guest_id = NULL WHERE guest_id = ?`,
		userID, guestID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE placement_results SET user_id = ?, guest_id = NULL WHERE guest_id = ?`,
		userID, guestID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE guest_sessions SET claimed_by = ? WHERE id = ?`, userID, guestID); err != nil {
		return 0, err
	}
	return int(moved), tx.Commit()
}

// PruneGuestSessions drops anonymous sessions that were never claimed and have
// gone quiet. Their submissions go with them via the foreign key.
func (s *Store) PruneGuestSessions(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan).UTC().Format("2006-01-02 15:04:05")
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM guest_sessions WHERE claimed_by IS NULL AND last_seen_at < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// --- submissions ------------------------------------------------------------

func (s *Store) CreateProblemSubmission(ctx context.Context, problemID int64, by Solver, lang, code string) (int64, error) {
	if !by.valid() {
		return 0, ErrNoSolver
	}
	userID, guestID := by.cols()
	var id int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO problem_submissions (problem_id, user_id, guest_id, language, code)
		VALUES (?, ?, ?, ?, ?) RETURNING id`,
		problemID, userID, guestID, lang, code).Scan(&id)
	return id, err
}

// PendingJudgeRuns lists submissions waiting on the judge.
func (s *Store) PendingJudgeRuns(ctx context.Context, limit int) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id FROM problem_submissions WHERE verdict = 'pending' ORDER BY id LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (s *Store) GetProblemSubmission(ctx context.Context, id int64) (*ProblemSubmission, error) {
	var sub ProblemSubmission
	var userID, guestID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT ps.id, ps.problem_id, ps.user_id, ps.guest_id, ps.language, ps.code,
		       ps.verdict, ps.passed, ps.total, ps.output, ps.failed_label, ps.runtime_ms,
		       ps.created_at, ps.judged_at, p.slug, p.time_limit_ms
		FROM problem_submissions ps JOIN problems p ON p.id = ps.problem_id
		WHERE ps.id = ?`, id).
		Scan(&sub.ID, &sub.ProblemID, &userID, &guestID, &sub.Language, &sub.Code,
			&sub.Verdict, &sub.Passed, &sub.Total, &sub.Output, &sub.FailedLabel,
			&sub.RuntimeMs, &sub.CreatedAt, &sub.JudgedAt, &sub.ProblemSlug, &sub.TimeLimitMs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	sub.By = Solver{UserID: userID.Int64, GuestID: guestID.Int64}
	return &sub, err
}

// LatestProblemSubmission returns the solver's most recent attempt, which is
// what the problem page shows them on return.
func (s *Store) LatestProblemSubmission(ctx context.Context, problemID int64, by Solver) (*ProblemSubmission, error) {
	if !by.valid() {
		return nil, ErrNotFound
	}
	userID, guestID := by.cols()
	var sub ProblemSubmission
	err := s.db.QueryRowContext(ctx, `
		SELECT ps.id, ps.problem_id, ps.language, ps.code, ps.verdict, ps.passed, ps.total,
		       ps.output, ps.failed_label, ps.runtime_ms, ps.created_at, ps.judged_at,
		       p.slug, p.time_limit_ms
		FROM problem_submissions ps JOIN problems p ON p.id = ps.problem_id
		WHERE ps.problem_id = ? AND (ps.user_id = ? OR ps.guest_id = ?)
		ORDER BY ps.id DESC LIMIT 1`, problemID, userID, guestID).
		Scan(&sub.ID, &sub.ProblemID, &sub.Language, &sub.Code, &sub.Verdict, &sub.Passed,
			&sub.Total, &sub.Output, &sub.FailedLabel, &sub.RuntimeMs, &sub.CreatedAt,
			&sub.JudgedAt, &sub.ProblemSlug, &sub.TimeLimitMs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	// The query already matched on the author, so this is who it is.
	sub.By = by
	return &sub, err
}

// Verdict is the judge's decision about one submission.
//
// A struct rather than a run of positional arguments because every field but
// the name is a number or a string, and a call site that swaps two of them
// would compile and quietly report the wrong result to a learner.
type Verdict struct {
	Name        string // one of the Verdict* constants
	Passed      int
	Total       int
	RuntimeMs   int
	Output      string
	FailedLabel string
}

// ProblemSubmissionsFor lists a solver's attempts at one problem, newest first
// — the history behind the Submissions tab.
func (s *Store) ProblemSubmissionsFor(ctx context.Context, problemID int64, by Solver, limit int) ([]ProblemSubmission, error) {
	if !by.valid() {
		return nil, nil
	}
	userID, guestID := by.cols()
	rows, err := s.db.QueryContext(ctx, `
		SELECT ps.id, ps.language, ps.verdict, ps.passed, ps.total, ps.failed_label,
		       ps.runtime_ms, ps.created_at
		FROM problem_submissions ps
		WHERE ps.problem_id = ? AND (ps.user_id = ? OR ps.guest_id = ?)
		ORDER BY ps.id DESC LIMIT ?`, problemID, userID, guestID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProblemSubmission
	for rows.Next() {
		var sub ProblemSubmission
		if err := rows.Scan(&sub.ID, &sub.Language, &sub.Verdict, &sub.Passed, &sub.Total,
			&sub.FailedLabel, &sub.RuntimeMs, &sub.CreatedAt); err != nil {
			return nil, err
		}
		sub.ProblemID = problemID
		sub.By = by
		out = append(out, sub)
	}
	return out, rows.Err()
}

// RecordVerdict stores the judge's decision.
func (s *Store) RecordVerdict(ctx context.Context, id int64, v Verdict) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE problem_submissions
		SET verdict = ?, passed = ?, total = ?, runtime_ms = ?, output = ?, failed_label = ?,
		    judged_at = datetime('now')
		WHERE id = ?`, v.Name, v.Passed, v.Total, v.RuntimeMs, v.Output, v.FailedLabel, id)
	return err
}

// SolvedCount is how many distinct problems a solver has accepted — the number
// a guest is shown, and the one they'd lose by not signing up.
func (s *Store) SolvedCount(ctx context.Context, by Solver) (int, error) {
	if !by.valid() {
		return 0, nil
	}
	userID, guestID := by.cols()
	var n int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT problem_id) FROM problem_submissions
		WHERE verdict = 'accepted' AND (user_id = ? OR guest_id = ?)`, userID, guestID).Scan(&n)
	return n, err
}
