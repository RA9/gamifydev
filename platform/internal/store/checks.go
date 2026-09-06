package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RA9/gamifydev/platform/internal/runner"
)

// Automated checks on assignments: the machine's half of a checkpoint verdict.
//
// The design rule (see the peer-review design doc): checks answer "does it
// work?" and gate progression. Peers answer "is it any good?" and never gate.
// Nothing in this file records a human judgment.

// Runnable reports whether the sandbox can execute an assignment's language,
// and therefore whether its checks will ever actually run.
//
// Delegates to the runner so the two cannot disagree: when shell support was
// added, a hardcoded copy of this list here left every shell checkpoint marked
// ungradable.
func Runnable(language string) bool {
	_, ok := runner.NormalizeLang(language)
	// An empty language would normalize to Python, but an assignment with no
	// language set is not a deliberate Python checkpoint.
	return ok && language != ""
}

// Check is one authored assertion about a submission.
type Check struct {
	ID           int64
	AssignmentID int64
	Sort         int
	Label        string
	Test         string
	Hidden       bool
	Points       int
	// Stdin is the input fed to the program for this check. Only meaningful for
	// compiled languages, where assertions are made over what the program
	// printed rather than over a namespace it left behind.
	Stdin string
}

// CheckResult is a check joined with how one submission fared against it.
type CheckResult struct {
	Check
	Passed bool
	// Ran is false when the submission predates the check, or the harness never
	// reached it — the UI must not show an unrun check as a failure.
	Ran bool
}

// ListChecks returns an assignment's checks in author order.
func (s *Store) ListChecks(ctx context.Context, assignmentID int64) ([]Check, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, assignment_id, sort, label, test, hidden, points, stdin
		FROM assignment_checks WHERE assignment_id = ? ORDER BY sort, id`, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Check
	for rows.Next() {
		var c Check
		var hidden int
		if err := rows.Scan(&c.ID, &c.AssignmentID, &c.Sort, &c.Label, &c.Test, &hidden, &c.Points, &c.Stdin); err != nil {
			return nil, err
		}
		c.Hidden = hidden == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

// ReplaceChecks synchronizes an assignment's checks and recomputes whether it
// can be auto-graded. Existing rows are updated by position so startup reseeding
// preserves submission_checks associations.
//
// Auto-gradable requires a language the sandbox can execute — Python, C, or
// shell. An assignment in any other language may still carry checks (they document the
// spec) but they will never run, and it falls back to mentor grading.
func (s *Store) ReplaceChecks(ctx context.Context, assignmentID int64, checks []Check) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	rows, err := tx.QueryContext(ctx,
		`SELECT id FROM assignment_checks WHERE assignment_id = ? ORDER BY sort, id`, assignmentID)
	if err != nil {
		return err
	}
	var existing []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		existing = append(existing, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for i, c := range checks {
		hidden := 0
		if c.Hidden {
			hidden = 1
		}
		if c.Points < 0 {
			c.Points = 0
		}
		if i < len(existing) {
			if _, err := tx.ExecContext(ctx, `
				UPDATE assignment_checks
				SET sort=?, label=?, test=?, hidden=?, points=?, stdin=?
				WHERE id=?`,
				i, c.Label, c.Test, hidden, c.Points, c.Stdin, existing[i]); err != nil {
				return err
			}
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO assignment_checks (assignment_id, sort, label, test, hidden, points, stdin)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			assignmentID, i, c.Label, c.Test, hidden, c.Points, c.Stdin); err != nil {
			return err
		}
	}
	if len(existing) > len(checks) {
		for _, id := range existing[len(checks):] {
			if _, err := tx.ExecContext(ctx, `DELETE FROM submission_checks WHERE check_id = ?`, id); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `DELETE FROM assignment_checks WHERE id = ?`, id); err != nil {
				return err
			}
		}
	}
	// Which languages are runnable is the runner's fact, not this query's —
	// hardcoding the list here is how it silently drifted when shell was added.
	var lang string
	if err := tx.QueryRowContext(ctx,
		`SELECT language FROM assignments WHERE id = ?`, assignmentID).Scan(&lang); err != nil {
		return err
	}
	auto := 0
	if len(checks) > 0 && Runnable(lang) {
		auto = 1
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE assignments SET auto_gradable = ?, updated_at = datetime('now') WHERE id = ?`,
		auto, assignmentID); err != nil {
		return err
	}
	return tx.Commit()
}

// AutoGradable reports whether an assignment's checks can actually be executed.
func (s *Store) AutoGradable(ctx context.Context, assignmentID int64) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT auto_gradable FROM assignments WHERE id = ?`, assignmentID).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotFound
	}
	return n == 1, err
}

// PendingCheckRuns lists submissions whose checks have not been run yet, for
// auto-gradable assignments only.
func (s *Store) PendingCheckRuns(ctx context.Context, limit int) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT sub.id FROM submissions sub
		JOIN assignments a ON a.id = sub.assignment_id
		WHERE sub.checks_ran_at IS NULL AND a.auto_gradable = 1
		ORDER BY sub.created_at LIMIT ?`, limit)
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

// SubmissionForChecks is the minimum needed to run a submission's checks.
type SubmissionForChecks struct {
	ID           int64
	AssignmentID int64
	UserID       int64
	Code         string
	MaxPoints    int
	PassPoints   int
	Language     string
}

// GetSubmissionForChecks loads a submission and the assignment settings the
// harness needs.
func (s *Store) GetSubmissionForChecks(ctx context.Context, id int64) (*SubmissionForChecks, error) {
	var sub SubmissionForChecks
	err := s.db.QueryRowContext(ctx, `
		SELECT sub.id, sub.assignment_id, sub.user_id, sub.code,
		       a.max_points, a.pass_points, a.language
		FROM submissions sub JOIN assignments a ON a.id = sub.assignment_id
		WHERE sub.id = ?`, id).
		Scan(&sub.ID, &sub.AssignmentID, &sub.UserID, &sub.Code,
			&sub.MaxPoints, &sub.PassPoints, &sub.Language)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &sub, err
}

// RecordCheckRun stores per-check outcomes and, when every weighted check
// passed, grades the submission automatically.
//
// Grading works by setting the existing `status` and `score` columns rather
// than introducing a parallel notion of passing — so the progression gate in
// AssignmentPassState keeps working untouched.
//
// A failing run is deliberately NOT marked `returned`: that status means a
// human sent work back, and conflating the two would tell a learner a person
// looked at their code when nobody did. It stays `submitted`, with the failed
// checks visible, so they can fix it and resubmit.
func (s *Store) RecordCheckRun(ctx context.Context, subID int64, passedIDs map[int64]bool, output string, maxPoints, passPoints int) (passed bool, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx,
		`DELETE FROM submission_checks WHERE submission_id = ?`, subID); err != nil {
		return false, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT c.id, c.points FROM assignment_checks c
		JOIN submissions sub ON sub.assignment_id = c.assignment_id
		WHERE sub.id = ?`, subID)
	if err != nil {
		return false, err
	}
	type ck struct {
		id     int64
		points int
	}
	var checks []ck
	for rows.Next() {
		var c ck
		if err := rows.Scan(&c.id, &c.points); err != nil {
			rows.Close()
			return false, err
		}
		checks = append(checks, c)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return false, err
	}

	gotPoints, totalPoints, passedCount, counted := 0, 0, 0, 0
	for _, c := range checks {
		ok := passedIDs[c.id]
		v := 0
		if ok {
			v = 1
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO submission_checks (submission_id, check_id, passed) VALUES (?, ?, ?)`,
			subID, c.id, v); err != nil {
			return false, err
		}
		// Zero-point checks are advisory: shown, but not part of the verdict.
		if c.points > 0 {
			counted++
			totalPoints += c.points
			if ok {
				passedCount++
				gotPoints += c.points
			}
		}
	}

	score := 0
	if totalPoints > 0 {
		score = gotPoints * maxPoints / totalPoints
	}
	// Every weighted check must pass. A checkpoint is a gate, not a grade curve:
	// partial credit would let a learner through with a broken solution.
	passed = counted > 0 && passedCount == counted

	set := `checks_ran_at = datetime('now'), checks_passed = ?, checks_total = ?, checks_output = ?`
	args := []any{passedCount, counted, output}
	if passed {
		set += `, status = 'graded', score = ?, graded_at = datetime('now')`
		args = append(args, score)
	}
	args = append(args, subID)
	if _, err := tx.ExecContext(ctx, `UPDATE submissions SET `+set+` WHERE id = ?`, args...); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return passed, nil
}

// SubmissionCheckResults returns a submission's checks with their outcomes, for
// the learner-facing results panel.
//
// Hidden checks are returned with their label blanked unless the submission
// passed, so a learner cannot read the full spec off the failure list.
func (s *Store) SubmissionCheckResults(ctx context.Context, subID int64) ([]CheckResult, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.assignment_id, c.sort, c.label, c.test, c.hidden, c.points, c.stdin,
		       COALESCE(sc.passed, 0), sc.check_id IS NOT NULL
		FROM assignment_checks c
		JOIN submissions sub ON sub.assignment_id = c.assignment_id
		LEFT JOIN submission_checks sc ON sc.check_id = c.id AND sc.submission_id = sub.id
		WHERE sub.id = ? ORDER BY c.sort, c.id`, subID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CheckResult
	allPassed := true
	for rows.Next() {
		var r CheckResult
		var hidden, passed, ran int
		if err := rows.Scan(&r.ID, &r.AssignmentID, &r.Sort, &r.Label, &r.Test,
			&hidden, &r.Points, &r.Stdin, &passed, &ran); err != nil {
			return nil, err
		}
		r.Hidden, r.Passed, r.Ran = hidden == 1, passed == 1, ran == 1
		if r.Points > 0 && !r.Passed {
			allPassed = false
		}
		// Never send the test expression or its input to a learner; together
		// they are the answer key.
		r.Test, r.Stdin = "", ""
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !allPassed {
		for i := range out {
			if out[i].Hidden {
				out[i].Label = "Hidden check"
			}
		}
	}
	return out, nil
}

// CheckHealth reports how often each check fails, so an author can tell a badly
// worded checkpoint from learners genuinely getting it wrong.
type CheckHealth struct {
	Label string
	Runs  int
	Fails int
}

// AssignmentCheckHealth aggregates outcomes for one assignment's checks.
func (s *Store) AssignmentCheckHealth(ctx context.Context, assignmentID int64) ([]CheckHealth, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.label, COUNT(sc.check_id),
		       SUM(CASE WHEN sc.passed = 0 THEN 1 ELSE 0 END)
		FROM assignment_checks c
		LEFT JOIN submission_checks sc ON sc.check_id = c.id
		WHERE c.assignment_id = ?
		GROUP BY c.id ORDER BY c.sort`, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CheckHealth
	for rows.Next() {
		var h CheckHealth
		var fails sql.NullInt64
		if err := rows.Scan(&h.Label, &h.Runs, &fails); err != nil {
			return nil, err
		}
		h.Fails = int(fails.Int64)
		out = append(out, h)
	}
	return out, rows.Err()
}

// --- Fixture files -----------------------------------------------------------

// File is an authored fixture written into the sandbox before a submission runs.
type File struct {
	Name    string
	Content string
}

// ListFiles returns an assignment's fixture files in author order.
func (s *Store) ListFiles(ctx context.Context, assignmentID int64) ([]File, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT name, content FROM assignment_files WHERE assignment_id = ? ORDER BY sort, id`,
		assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.Name, &f.Content); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// FileMap returns an assignment's fixtures shaped for a runner request.
func (s *Store) FileMap(ctx context.Context, assignmentID int64) (map[string]string, error) {
	files, err := s.ListFiles(ctx, assignmentID)
	if err != nil || len(files) == 0 {
		return nil, err
	}
	out := make(map[string]string, len(files))
	for _, f := range files {
		out[f.Name] = f.Content
	}
	return out, nil
}

// ReplaceFiles swaps an assignment's fixture files for the given set.
func (s *Store) ReplaceFiles(ctx context.Context, assignmentID int64, files []File) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM assignment_files WHERE assignment_id = ?`, assignmentID); err != nil {
		return err
	}
	seen := map[string]bool{}
	for i, f := range files {
		if f.Name == "" || seen[f.Name] {
			continue // a blank row, or a duplicate name the unique index would reject
		}
		seen[f.Name] = true
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO assignment_files (assignment_id, sort, name, content) VALUES (?, ?, ?, ?)`,
			assignmentID, i, f.Name, f.Content); err != nil {
			return err
		}
	}
	return tx.Commit()
}
