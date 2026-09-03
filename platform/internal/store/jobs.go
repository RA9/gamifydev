package store

import (
	"context"
	"database/sql"
	"time"
)

// Job leases and run history. These back internal/jobs.Runner; see that package
// for why the runner needs them.

// AcquireJobLock takes a lease on `name` for `ttl`, returning false if another
// holder currently owns it. An expired lease is stolen — that is what stops a
// crashed instance from wedging the schedule.
//
// The whole decision happens in one statement so two instances racing on the
// same tick cannot both win: SQLite serializes the writes, and the WHERE clause
// on the upsert only lets the row through if it is unheld, expired, or already
// ours (re-entrant for the same holder).
func (s *Store) AcquireJobLock(ctx context.Context, name, holder string, ttl time.Duration) (bool, error) {
	secs := int(ttl / time.Second)
	if secs < 1 {
		secs = 1
	}
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO job_locks (name, holder, acquired_at, expires_at)
		VALUES (?, ?, datetime('now'), datetime('now', '+' || ? || ' seconds'))
		ON CONFLICT(name) DO UPDATE SET
			holder      = excluded.holder,
			acquired_at = excluded.acquired_at,
			expires_at  = excluded.expires_at
		WHERE job_locks.expires_at <= datetime('now')
		   OR job_locks.holder = excluded.holder`,
		name, holder, secs)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

// ReleaseJobLock drops a lease we hold. Releasing a lease we no longer own (it
// expired and was stolen) is a no-op, not an error.
func (s *Store) ReleaseJobLock(ctx context.Context, name, holder string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM job_locks WHERE name = ? AND holder = ?`, name, holder)
	return err
}

// StartJobRun opens a run record and returns its id.
func (s *Store) StartJobRun(ctx context.Context, name string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO job_runs (name, started_at) VALUES (?, datetime('now')) RETURNING id`, name).
		Scan(&id)
	return id, err
}

// FinishJobRun closes a run record with its outcome.
func (s *Store) FinishJobRun(ctx context.Context, id int64, ok bool, detail string, dur time.Duration) error {
	okv := 0
	if ok {
		okv = 1
	}
	if len(detail) > 500 {
		detail = detail[:500]
	}
	_, err := s.db.ExecContext(ctx, `
		UPDATE job_runs
		SET finished_at = datetime('now'), ok = ?, detail = ?, duration_ms = ?
		WHERE id = ?`,
		okv, detail, dur.Milliseconds(), id)
	return err
}

// JobRun is one recorded execution, for the admin view.
type JobRun struct {
	ID         int64
	Name       string
	StartedAt  string
	FinishedAt sql.NullString
	OK         sql.NullBool
	Detail     string
	DurationMS sql.NullInt64
}

// Running reports whether this run never recorded an outcome — either it is
// in flight, or the process died mid-job.
func (r JobRun) Running() bool { return !r.OK.Valid }

// Failed reports a recorded failure.
func (r JobRun) Failed() bool { return r.OK.Valid && !r.OK.Bool }

// RecentJobRuns returns the newest runs across all jobs.
func (s *Store) RecentJobRuns(ctx context.Context, limit int) ([]JobRun, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, started_at, finished_at, ok, detail, duration_ms
		FROM job_runs ORDER BY started_at DESC, id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []JobRun
	for rows.Next() {
		var r JobRun
		if err := rows.Scan(&r.ID, &r.Name, &r.StartedAt, &r.FinishedAt, &r.OK, &r.Detail, &r.DurationMS); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// JobHealth summarises one job's recent history for the admin overview.
type JobHealth struct {
	Name     string
	LastRun  sql.NullString
	LastOK   sql.NullBool
	Runs24h  int
	Fails24h int
}

// JobHealthSummary aggregates the last 24 hours per job name.
func (s *Store) JobHealthSummary(ctx context.Context) ([]JobHealth, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT name,
		       MAX(started_at) AS last_run,
		       SUM(CASE WHEN started_at >= datetime('now','-1 day') THEN 1 ELSE 0 END) AS runs,
		       SUM(CASE WHEN started_at >= datetime('now','-1 day') AND ok = 0 THEN 1 ELSE 0 END) AS fails
		FROM job_runs
		GROUP BY name ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []JobHealth
	for rows.Next() {
		var h JobHealth
		if err := rows.Scan(&h.Name, &h.LastRun, &h.Runs24h, &h.Fails24h); err != nil {
			return nil, err
		}
		// The most recent run's outcome, for the status dot.
		_ = s.db.QueryRowContext(ctx,
			`SELECT ok FROM job_runs WHERE name = ? ORDER BY started_at DESC, id DESC LIMIT 1`, h.Name).
			Scan(&h.LastOK)
		out = append(out, h)
	}
	return out, rows.Err()
}

// PruneJobRuns drops run history older than `days`, keeping the table bounded.
func (s *Store) PruneJobRuns(ctx context.Context, days int) (int, error) {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM job_runs WHERE started_at < datetime('now', '-' || ? || ' days')`, days)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
