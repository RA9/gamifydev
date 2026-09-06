package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// A candidate is whoever is sitting the placement diagnostic, identified by
// their email address rather than by the cookie in front of them.
//
// The distinction is the point. Solver answers "which browser is this?", which
// is the right question for keeping one visitor's work separate from another's.
// It is the wrong question for a gate: the retry cooldown, the attempt limit
// and a sitting voided for cheating all have to survive a private window, and a
// cookie cannot carry that because the candidate owns it.
type Candidate struct {
	ID       int64
	Email    string // normalized — the identity
	RawEmail string // as typed — the address to write to
	Name     string
	UserID   int64 // non-zero once they hold an account
}

// ErrInvalidCandidate means the details given are not usable as an identity.
var ErrInvalidCandidate = errors.New("candidate name or email is not usable")

// NormalizeEmail reduces an address to the identity it represents.
//
// Lower-cased, and any +tag dropped from the local part. Plus-addressing is the
// one-keystroke way to look like a new person to a system that keys on email,
// and closing it costs nothing legitimate: mail to alice+test@x.com is
// delivered to alice@x.com by every provider that offers the feature.
//
// Deliberately no dot-stripping. That is a Gmail rule, not an email rule, and
// applying it everywhere would merge genuinely different people at providers
// where first.last@ and firstlast@ are two accounts.
func NormalizeEmail(raw string) string {
	e := strings.ToLower(strings.TrimSpace(raw))
	at := strings.LastIndex(e, "@")
	if at <= 0 {
		return e
	}
	local, domain := e[:at], e[at:]
	if plus := strings.Index(local, "+"); plus > 0 {
		local = local[:plus]
	}
	return local + domain
}

// validEmail is the shallow check worth making here: something before an @,
// something after it, and a dot in the domain. Anything stricter rejects real
// addresses, and anything at all only filters typos — the address is not
// verified, which is what the comment on the migration says out loud.
func validEmail(e string) bool {
	at := strings.LastIndex(e, "@")
	if at <= 0 || at == len(e)-1 {
		return false
	}
	return strings.Contains(e[at+1:], ".") && !strings.ContainsAny(e, " \t\n")
}

// UpsertCandidate resolves a name and email to the person behind them, creating
// the record on first sight and refreshing the name on later ones.
func (s *Store) UpsertCandidate(ctx context.Context, name, rawEmail string) (*Candidate, error) {
	name = strings.TrimSpace(name)
	rawEmail = strings.TrimSpace(rawEmail)
	email := NormalizeEmail(rawEmail)
	if len(name) < 2 || !validEmail(email) {
		return nil, ErrInvalidCandidate
	}
	// A later sitting may give a different spelling of the same person's name.
	// Take the newer one: it is the one they just typed.
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO candidates (email, email_raw, name) VALUES (?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			name = excluded.name, email_raw = excluded.email_raw,
			last_seen_at = datetime('now')`, email, rawEmail, name)
	if err != nil {
		return nil, err
	}
	return s.CandidateByEmail(ctx, email)
}

// CandidateByEmail looks up an identity. The address is normalized first, so
// any spelling of the same identity finds the same record.
func (s *Store) CandidateByEmail(ctx context.Context, rawEmail string) (*Candidate, error) {
	var c Candidate
	var userID sql.NullInt64
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, email_raw, name, user_id FROM candidates WHERE email = ?`,
		NormalizeEmail(rawEmail)).Scan(&c.ID, &c.Email, &c.RawEmail, &c.Name, &userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c.UserID = userID.Int64
	return &c, nil
}

// CandidateForAttempt returns whoever sat a given paper, or nil for a sitting
// taken before candidates existed.
func (s *Store) CandidateForAttempt(ctx context.Context, attemptID int64) (*Candidate, error) {
	var c Candidate
	var userID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT c.id, c.email, c.email_raw, c.name, c.user_id
		FROM candidates c JOIN assessment_attempts a ON a.candidate_id = c.id
		WHERE a.id = ?`, attemptID).Scan(&c.ID, &c.Email, &c.RawEmail, &c.Name, &userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.UserID = userID.Int64
	return &c, nil
}

// LinkCandidateToUser records that an identity now holds an account, so a later
// sitting from any browser sees their enrollment rather than starting fresh.
//
// Only ever fills an empty link. Repointing a candidate at a second account
// would be the way to launder a spent identity into a clean one.
func (s *Store) LinkCandidateToUser(ctx context.Context, candidateID, userID int64) error {
	if candidateID == 0 || userID == 0 {
		return nil
	}
	_, err := s.db.ExecContext(ctx,
		`UPDATE candidates SET user_id = ? WHERE id = ? AND user_id IS NULL`, userID, candidateID)
	return err
}
