// Package jobs is the platform's scheduled-work runner.
//
// Every mechanism in the cohort program is time-driven — standup windows open
// and close, absences accrue, sprints roll over — and none of it can hang off an
// HTTP request. This package provides the loop that makes that possible.
//
// Design notes:
//
//   - One in-process runner per app instance, started from main and stopped on
//     shutdown. No external scheduler to operate.
//   - Before running, a job takes a short lease in `job_locks`. If two instances
//     are deployed, only one runs a given job on a given tick. Leases carry an
//     expiry, so a crashed instance cannot wedge the schedule.
//   - Every attempt is recorded in `job_runs` (including failures), so the
//     schedule is observable from the admin UI rather than only from logs.
//   - Jobs must be idempotent. They can and will run twice: overlapping windows,
//     a restart mid-run, a lease that expired under a slow job.
package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

// Func is the work a job performs. The returned string is a short human-readable
// summary ("rolled up 12 days") stored with the run for observability.
type Func func(ctx context.Context) (string, error)

// Job is a registered unit of scheduled work.
type Job struct {
	Name  string        // stable identifier, also the lock key: "activity:rollup"
	Every time.Duration // how often to attempt it
	// Timeout bounds a single execution. The lease is taken for slightly longer
	// so a job cannot outlive its own lock.
	Timeout time.Duration
	Run     Func
}

// Store is the slice of the data layer the runner needs. Keeping it as an
// interface means jobs are testable without a real database.
type Store interface {
	AcquireJobLock(ctx context.Context, name, holder string, ttl time.Duration) (bool, error)
	ReleaseJobLock(ctx context.Context, name, holder string) error
	StartJobRun(ctx context.Context, name string) (int64, error)
	FinishJobRun(ctx context.Context, id int64, ok bool, detail string, dur time.Duration) error
}

// Runner executes registered jobs on their cadence until stopped.
type Runner struct {
	st     Store
	holder string // identifies this instance in the lock table
	jobs   []Job
	// tick is how often the runner wakes to see what is due. Jobs whose Every is
	// shorter than this effectively run at tick cadence.
	tick time.Duration

	mu   sync.Mutex
	last map[string]time.Time

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// New builds a runner. `holder` should identify the process; when empty a
// hostname/pid pair is used.
func New(st Store, holder string) *Runner {
	if holder == "" {
		host, _ := os.Hostname()
		holder = fmt.Sprintf("%s/%d", host, os.Getpid())
	}
	return &Runner{
		st:     st,
		holder: holder,
		tick:   time.Minute,
		last:   map[string]time.Time{},
	}
}

// Register adds a job. Must be called before Start.
func (r *Runner) Register(j Job) {
	if j.Timeout == 0 {
		j.Timeout = 2 * time.Minute
	}
	if j.Every == 0 {
		j.Every = time.Hour
	}
	r.jobs = append(r.jobs, j)
}

// Jobs returns the registered jobs, name-sorted (for the admin view).
func (r *Runner) Jobs() []Job {
	out := append([]Job(nil), r.jobs...)
	sort.Slice(out, func(i, k int) bool { return out[i].Name < out[k].Name })
	return out
}

// Start begins the loop in the background. Safe to call with no jobs registered.
func (r *Runner) Start(parent context.Context) {
	if len(r.jobs) == 0 {
		log.Printf("jobs: no jobs registered, runner idle")
		return
	}
	ctx, cancel := context.WithCancel(parent)
	r.cancel = cancel
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		t := time.NewTicker(r.tick)
		defer t.Stop()
		// One pass immediately so a restart doesn't wait a full tick.
		r.pass(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				r.pass(ctx)
			}
		}
	}()
	names := make([]string, 0, len(r.jobs))
	for _, j := range r.jobs {
		names = append(names, j.Name)
	}
	sort.Strings(names)
	log.Printf("jobs: runner started (%d jobs: %v)", len(r.jobs), names)
}

// Stop cancels the loop and waits for the in-flight pass to finish.
func (r *Runner) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()
}

// pass runs every job that is due. Jobs run sequentially: the workloads are
// small, and serial execution keeps SQLite write contention low.
func (r *Runner) pass(ctx context.Context) {
	for _, j := range r.jobs {
		if ctx.Err() != nil {
			return
		}
		if !r.due(j) {
			continue
		}
		r.runOne(ctx, j)
	}
}

func (r *Runner) due(j Job) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	last, seen := r.last[j.Name]
	return !seen || time.Since(last) >= j.Every
}

// RunNow executes a job by name immediately, bypassing its cadence but still
// taking the lease. Used by the admin "run now" control.
func (r *Runner) RunNow(ctx context.Context, name string) error {
	for _, j := range r.jobs {
		if j.Name == name {
			r.runOne(ctx, j)
			return nil
		}
	}
	return fmt.Errorf("unknown job %q", name)
}

func (r *Runner) runOne(ctx context.Context, j Job) {
	// Hold the lease a little longer than the job may run, so it cannot be
	// stolen mid-execution but still expires if we die.
	lease := j.Timeout + 30*time.Second
	got, err := r.st.AcquireJobLock(ctx, j.Name, r.holder, lease)
	if err != nil {
		log.Printf("jobs: %s: lock error: %v", j.Name, err)
		return
	}
	if !got {
		// Another instance has it; treat as attempted so we don't spin.
		r.mark(j.Name)
		return
	}
	defer func() {
		if err := r.st.ReleaseJobLock(context.WithoutCancel(ctx), j.Name, r.holder); err != nil {
			log.Printf("jobs: %s: unlock error: %v", j.Name, err)
		}
	}()
	r.mark(j.Name)

	runID, err := r.st.StartJobRun(ctx, j.Name)
	if err != nil {
		log.Printf("jobs: %s: could not record start: %v", j.Name, err)
		// Still run the work — observability failing shouldn't stop the program.
	}

	jctx, cancel := context.WithTimeout(ctx, j.Timeout)
	start := time.Now()
	detail, runErr := safeRun(jctx, j.Run)
	cancel()
	dur := time.Since(start)

	if runErr != nil {
		log.Printf("jobs: %s failed after %s: %v", j.Name, dur.Round(time.Millisecond), runErr)
		detail = runErr.Error()
	} else if detail != "" {
		log.Printf("jobs: %s ok in %s — %s", j.Name, dur.Round(time.Millisecond), detail)
	}
	if runID > 0 {
		// Use a detached context: the run record must be written even when the
		// parent context is being cancelled by shutdown.
		fctx, fcancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		if err := r.st.FinishJobRun(fctx, runID, runErr == nil, detail, dur); err != nil {
			log.Printf("jobs: %s: could not record finish: %v", j.Name, err)
		}
		fcancel()
	}
}

// safeRun converts a panic in job code into an error, so one bad job cannot take
// the process down.
func safeRun(ctx context.Context, fn Func) (detail string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()
	return fn(ctx)
}

func (r *Runner) mark(name string) {
	r.mu.Lock()
	r.last[name] = time.Now()
	r.mu.Unlock()
}

// ErrNoRows re-exports the sentinel so job packages don't import database/sql
// only for this.
var ErrNoRows = sql.ErrNoRows

// IsNoRows reports whether err is a missing-row error.
func IsNoRows(err error) bool { return errors.Is(err, sql.ErrNoRows) }
