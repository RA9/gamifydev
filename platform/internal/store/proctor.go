package store

import (
	"context"
	"database/sql"
	"errors"
)

// Proctoring for the placement diagnostic.
//
// The browser reports what it sees; this decides what it means. That split
// matters: everything the page can detect, a determined candidate can suppress,
// so the count, the limit and the void all live here where the candidate cannot
// reach them. A client that stops reporting simply stops accumulating strikes —
// it can never talk its way out of one already recorded.

// ProctorStrikeLimit is how many caught violations end a sitting. The first is
// a warning; the second closes the paper.
const ProctorStrikeLimit = 2

// VoidCheating is the reason recorded against an attempt the proctor closed.
const VoidCheating = "cheating"

// ProctorKinds are the events a browser may report. Anything else is rejected
// rather than stored, so the log stays a fixed vocabulary that an administrator
// reviewing a voided sitting can actually read.
var ProctorKinds = map[string]bool{
	"copy":    true, // copy, cut, or a right-click aimed at one
	"hidden":  true, // the tab stopped being the visible one
	"blur":    true, // the window lost focus
	"print":   true, // a print or save dialogue was invoked
	"capture": true, // a screenshot shortcut we were able to see
}

// ErrAttemptVoided means the sitting was closed by the proctor. Distinct from
// ErrAlreadySubmitted because the candidate is owed a different explanation.
var ErrAttemptVoided = errors.New("attempt voided by the proctor")

// ErrUnknownViolation means the client reported something not in ProctorKinds.
var ErrUnknownViolation = errors.New("unknown violation kind")

// ProctorOutcome is what the page is told after reporting a violation.
type ProctorOutcome struct {
	Strikes int  // how many have now been recorded for this sitting
	Limit   int  // how many end it
	Voided  bool // whether this one ended it
}

// RecordViolation logs one caught violation against a live attempt and reports
// where that leaves the candidate.
//
// Idempotency is deliberately not attempted. A browser that fires both a blur
// and a visibilitychange for one alt-tab would spend two strikes, so the client
// is responsible for reporting a single event once — see proctor.js, which
// coalesces them. Doing that here instead would mean guessing which two reports
// a second apart were the same act, and guessing wrong in the direction that
// fails an honest candidate.
func (s *Store) RecordViolation(ctx context.Context, by Solver, attemptID int64, kind string) (ProctorOutcome, error) {
	out := ProctorOutcome{Limit: ProctorStrikeLimit}
	if !ProctorKinds[kind] {
		return out, ErrUnknownViolation
	}
	if !by.valid() {
		return out, ErrNoSolver
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return out, err
	}
	defer tx.Rollback() //nolint:errcheck

	userID, guestID := by.cols()
	var submitted sql.NullString
	var violations int
	var voided string
	var expired int
	err = tx.QueryRowContext(ctx, `
		SELECT submitted_at, violations, voided_reason,
		       CASE WHEN datetime(expires_at) <= datetime('now') THEN 1 ELSE 0 END
		FROM assessment_attempts
		WHERE id = ? AND (user_id = ? OR guest_id = ?)`, attemptID, userID, guestID).
		Scan(&submitted, &violations, &voided, &expired)
	if errors.Is(err, sql.ErrNoRows) {
		return out, ErrNotFound
	}
	if err != nil {
		return out, err
	}
	// A finished sitting cannot be cheated on. Report where it stands rather
	// than erroring: a page closing down may well fire one last event, and that
	// is not a failure worth showing anybody.
	if voided != "" {
		return ProctorOutcome{Strikes: violations, Limit: ProctorStrikeLimit, Voided: true}, nil
	}
	if submitted.Valid || expired == 1 {
		return ProctorOutcome{Strikes: violations, Limit: ProctorStrikeLimit}, nil
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO attempt_violations (attempt_id, kind) VALUES (?, ?)`, attemptID, kind); err != nil {
		return out, err
	}
	violations++
	out.Strikes = violations

	if violations >= ProctorStrikeLimit {
		// Close the paper here rather than trusting the page to stop. The score
		// is zero and the sitting is spent, which is what "failed" has to mean
		// for the retry rules to treat it like any other failure.
		if _, err := tx.ExecContext(ctx, `
			UPDATE assessment_attempts
			SET violations = ?, voided_reason = ?, submitted_at = datetime('now'),
			    score = 0, topic_scores = '{}'
			WHERE id = ?`, violations, VoidCheating, attemptID); err != nil {
			return out, err
		}
		out.Voided = true
	} else if _, err := tx.ExecContext(ctx,
		`UPDATE assessment_attempts SET violations = ? WHERE id = ?`, violations, attemptID); err != nil {
		return out, err
	}
	return out, tx.Commit()
}

// AttemptViolations lists what was recorded against a sitting, oldest first —
// the evidence behind a void, for whoever has to answer an appeal.
func (s *Store) AttemptViolations(ctx context.Context, attemptID int64) ([]Violation, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT kind, created_at FROM attempt_violations WHERE attempt_id = ? ORDER BY id`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Violation
	for rows.Next() {
		var v Violation
		if err := rows.Scan(&v.Kind, &v.At); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// Violation is one logged proctoring event.
type Violation struct {
	Kind string
	At   string
}
