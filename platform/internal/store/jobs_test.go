package store

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestJobLockIsSingleFlight(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	got, err := st.AcquireJobLock(ctx, "j", "a", time.Minute)
	if err != nil || !got {
		t.Fatalf("first acquire = %v, %v; want true, nil", got, err)
	}

	// A different holder must be refused while the lease is live — this is the
	// property that makes running two app instances safe.
	got, err = st.AcquireJobLock(ctx, "j", "b", time.Minute)
	if err != nil {
		t.Fatalf("second acquire: %v", err)
	}
	if got {
		t.Fatal("a second holder acquired a live lease")
	}

	// The same holder is re-entrant (a re-run extends its own lease).
	if got, err := st.AcquireJobLock(ctx, "j", "a", time.Minute); err != nil || !got {
		t.Fatalf("same-holder re-acquire = %v, %v; want true, nil", got, err)
	}

	// Releasing as a non-owner must not steal the lock away from its owner.
	if err := st.ReleaseJobLock(ctx, "j", "b"); err != nil {
		t.Fatalf("release by non-owner: %v", err)
	}
	if got, _ := st.AcquireJobLock(ctx, "j", "c", time.Minute); got {
		t.Fatal("non-owner release freed another holder's lease")
	}

	if err := st.ReleaseJobLock(ctx, "j", "a"); err != nil {
		t.Fatalf("release: %v", err)
	}
	if got, err := st.AcquireJobLock(ctx, "j", "b", time.Minute); err != nil || !got {
		t.Fatalf("acquire after release = %v, %v; want true, nil", got, err)
	}
}

func TestJobLockExpires(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	// A lease that has already elapsed must be stealable, otherwise a crashed
	// instance would wedge the job forever.
	if got, err := st.AcquireJobLock(ctx, "j", "dead", time.Second); err != nil || !got {
		t.Fatalf("acquire: %v %v", got, err)
	}
	mustExec(t, st, `UPDATE job_locks SET expires_at = datetime('now','-1 minute') WHERE name='j'`)

	got, err := st.AcquireJobLock(ctx, "j", "alive", time.Minute)
	if err != nil {
		t.Fatalf("steal: %v", err)
	}
	if !got {
		t.Fatal("expired lease was not stealable")
	}
	var holder string
	if err := st.db.QueryRow(`SELECT holder FROM job_locks WHERE name='j'`).Scan(&holder); err != nil {
		t.Fatalf("read holder: %v", err)
	}
	if holder != "alive" {
		t.Fatalf("holder = %q, want alive", holder)
	}
}

func TestJobLockConcurrentRace(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	// Many holders contend for one lease at once; exactly one may win. This is
	// the case the single-statement upsert exists to handle.
	const n = 12
	var wg sync.WaitGroup
	var mu sync.Mutex
	wins := 0
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			ok, err := st.AcquireJobLock(ctx, "contended", string(rune('a'+i)), time.Minute)
			if err != nil {
				return // busy/locked is acceptable; double-winning is not
			}
			if ok {
				mu.Lock()
				wins++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	if wins != 1 {
		t.Fatalf("%d holders acquired the same lease, want exactly 1", wins)
	}
}

func TestJobRunHistory(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	id, err := st.StartJobRun(ctx, "activity:rollup")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	runs, _ := st.RecentJobRuns(ctx, 10)
	if len(runs) != 1 || !runs[0].Running() {
		t.Fatalf("in-flight run not reported as running: %+v", runs)
	}

	if err := st.FinishJobRun(ctx, id, true, "recorded 3 active day(s)", 120*time.Millisecond); err != nil {
		t.Fatalf("finish: %v", err)
	}
	runs, _ = st.RecentJobRuns(ctx, 10)
	if runs[0].Running() || runs[0].Failed() {
		t.Fatalf("finished ok run reads wrong: %+v", runs[0])
	}
	if !runs[0].DurationMS.Valid || runs[0].DurationMS.Int64 != 120 {
		t.Fatalf("duration = %+v, want 120", runs[0].DurationMS)
	}

	// A failure must be recorded, not swallowed — the admin view depends on it.
	id2, _ := st.StartJobRun(ctx, "activity:rollup")
	if err := st.FinishJobRun(ctx, id2, false, "boom", time.Second); err != nil {
		t.Fatalf("finish fail: %v", err)
	}
	health, err := st.JobHealthSummary(ctx)
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if len(health) != 1 {
		t.Fatalf("health rows = %d, want 1", len(health))
	}
	if health[0].Runs24h != 2 || health[0].Fails24h != 1 {
		t.Fatalf("health = %+v, want 2 runs / 1 fail", health[0])
	}
	if !health[0].LastOK.Valid || health[0].LastOK.Bool {
		t.Fatalf("LastOK = %+v, want the most recent run (failed)", health[0].LastOK)
	}
}
