package jobs

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeStore records what the runner asked of the data layer, and can simulate
// another instance already holding a lease.
type fakeStore struct {
	mu       sync.Mutex
	locked   map[string]string // name -> holder
	refuse   bool              // pretend someone else holds every lease
	starts   int32
	finishes int32
	lastOK   bool
	lastDeta string
	released int32
}

func newFake() *fakeStore { return &fakeStore{locked: map[string]string{}} }

func (f *fakeStore) AcquireJobLock(_ context.Context, name, holder string, _ time.Duration) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.refuse {
		return false, nil
	}
	if h, ok := f.locked[name]; ok && h != holder {
		return false, nil
	}
	f.locked[name] = holder
	return true, nil
}

func (f *fakeStore) ReleaseJobLock(_ context.Context, name, holder string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.locked[name] == holder {
		delete(f.locked, name)
		atomic.AddInt32(&f.released, 1)
	}
	return nil
}

func (f *fakeStore) StartJobRun(_ context.Context, _ string) (int64, error) {
	return int64(atomic.AddInt32(&f.starts, 1)), nil
}

func (f *fakeStore) FinishJobRun(_ context.Context, _ int64, ok bool, detail string, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	atomic.AddInt32(&f.finishes, 1)
	f.lastOK, f.lastDeta = ok, detail
	return nil
}

func TestRunNowRecordsSuccess(t *testing.T) {
	f := newFake()
	r := New(f, "test")
	var ran int32
	r.Register(Job{Name: "ok", Every: time.Hour, Run: func(context.Context) (string, error) {
		atomic.AddInt32(&ran, 1)
		return "did 3 things", nil
	}})

	if err := r.RunNow(context.Background(), "ok"); err != nil {
		t.Fatalf("RunNow: %v", err)
	}
	if ran != 1 {
		t.Fatalf("job ran %d times, want 1", ran)
	}
	if f.starts != 1 || f.finishes != 1 {
		t.Fatalf("history: %d starts / %d finishes, want 1/1", f.starts, f.finishes)
	}
	if !f.lastOK || f.lastDeta != "did 3 things" {
		t.Fatalf("recorded ok=%v detail=%q", f.lastOK, f.lastDeta)
	}
	if f.released != 1 {
		t.Fatalf("lease released %d times, want 1", f.released)
	}
}

func TestFailureIsRecordedAndLeaseReleased(t *testing.T) {
	f := newFake()
	r := New(f, "test")
	r.Register(Job{Name: "bad", Every: time.Hour, Run: func(context.Context) (string, error) {
		return "", errors.New("boom")
	}})

	_ = r.RunNow(context.Background(), "bad")
	if f.lastOK {
		t.Fatal("failure recorded as success")
	}
	if f.lastDeta != "boom" {
		t.Fatalf("detail = %q, want the error text", f.lastDeta)
	}
	// A failing job must not hold its lease — otherwise one bad run blocks the
	// job until the lease expires.
	if f.released != 1 {
		t.Fatalf("lease released %d times after failure, want 1", f.released)
	}
}

func TestPanicDoesNotKillTheRunner(t *testing.T) {
	f := newFake()
	r := New(f, "test")
	r.Register(Job{Name: "panicky", Every: time.Hour, Run: func(context.Context) (string, error) {
		panic("bad pointer")
	}})

	// The point of the test: this call returns instead of taking the process down.
	_ = r.RunNow(context.Background(), "panicky")
	if f.lastOK {
		t.Fatal("panicking job recorded as success")
	}
	if f.lastDeta == "" {
		t.Fatal("panic produced no recorded detail")
	}
	if f.released != 1 {
		t.Fatalf("lease released %d times after panic, want 1", f.released)
	}
}

func TestSkipsWhenAnotherInstanceHoldsTheLease(t *testing.T) {
	f := newFake()
	f.refuse = true
	r := New(f, "test")
	var ran int32
	r.Register(Job{Name: "held", Every: time.Hour, Run: func(context.Context) (string, error) {
		atomic.AddInt32(&ran, 1)
		return "", nil
	}})

	_ = r.RunNow(context.Background(), "held")
	if ran != 0 {
		t.Fatal("job ran despite another instance holding the lease")
	}
	if f.starts != 0 {
		t.Fatal("recorded a run that never happened")
	}
}

func TestTimeoutCancelsTheJobContext(t *testing.T) {
	f := newFake()
	r := New(f, "test")
	done := make(chan struct{})
	r.Register(Job{
		Name:    "slow",
		Every:   time.Hour,
		Timeout: 30 * time.Millisecond,
		Run: func(ctx context.Context) (string, error) {
			<-ctx.Done() // must be cancelled by the timeout, not hang forever
			close(done)
			return "", ctx.Err()
		},
	})

	_ = r.RunNow(context.Background(), "slow")
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("job context was never cancelled by its timeout")
	}
	if f.lastOK {
		t.Fatal("timed-out job recorded as success")
	}
}

func TestStartStopRunsDueJobs(t *testing.T) {
	f := newFake()
	r := New(f, "test")
	r.tick = 5 * time.Millisecond
	var ran int32
	r.Register(Job{Name: "tick", Every: time.Millisecond, Run: func(context.Context) (string, error) {
		atomic.AddInt32(&ran, 1)
		return "", nil
	}})

	r.Start(context.Background())
	// The runner does one pass immediately, so this should not need to wait long.
	deadline := time.After(2 * time.Second)
	for atomic.LoadInt32(&ran) < 2 {
		select {
		case <-deadline:
			t.Fatalf("job ran %d times, expected repeat execution", ran)
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}
	r.Stop() // must return promptly and not deadlock

	after := atomic.LoadInt32(&ran)
	time.Sleep(30 * time.Millisecond)
	if atomic.LoadInt32(&ran) != after {
		t.Fatal("job kept running after Stop")
	}
}

func TestUnknownJob(t *testing.T) {
	r := New(newFake(), "test")
	if err := r.RunNow(context.Background(), "nope"); err == nil {
		t.Fatal("RunNow on an unregistered job should error")
	}
}
