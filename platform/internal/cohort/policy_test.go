package cohort

import (
	"testing"
	"time"
)

func TestShouldFormOnSizeOrWait(t *testing.T) {
	// The two qualifying routes, and the case that must not form.
	if ShouldForm(TargetSize, 0) != true {
		t.Fatal("a full queue did not qualify")
	}
	if ShouldForm(TargetSize-1, 0) != false {
		t.Fatal("a short queue qualified with no wait")
	}
	// The wait route is what makes a cohort of one possible on day zero — the
	// alternative is leaving the first learner queued forever.
	if ShouldForm(1, MaxQueueWait) != true {
		t.Fatal("a long-waiting single learner did not qualify")
	}
	if ShouldForm(1, MaxQueueWait-time.Minute) != false {
		t.Fatal("a single learner qualified before the max wait")
	}
	// An empty queue never forms, whatever the age says.
	if ShouldForm(0, 100*time.Hour) != false {
		t.Fatal("an empty queue formed a cohort")
	}
}

func TestTakeSizeCaps(t *testing.T) {
	if got := TakeSize(3); got != 3 {
		t.Fatalf("TakeSize(3) = %d, want 3", got)
	}
	if got := TakeSize(MaxSize + 5); got != MaxSize {
		t.Fatalf("TakeSize over max = %d, want %d", got, MaxSize)
	}
}

func TestSizingIsInternallyConsistent(t *testing.T) {
	// The PRD's decision is that cohorts form above their steady state so decay
	// lands near four. If these ever invert, formation and repacking fight.
	if MinViable >= TargetSize {
		t.Fatalf("MinViable (%d) must be below TargetSize (%d)", MinViable, TargetSize)
	}
	if TargetSize > MaxSize {
		t.Fatalf("TargetSize (%d) must not exceed MaxSize (%d)", TargetSize, MaxSize)
	}
}

func TestLookupBandFallsBack(t *testing.T) {
	// A missing or unknown band must never block formation.
	if got := LookupBand("").Key; got != BandEuropeAfrica {
		t.Fatalf("empty band = %q, want the default", got)
	}
	if got := LookupBand("mars").Key; got != BandEuropeAfrica {
		t.Fatalf("unknown band = %q, want the default", got)
	}
	if got := LookupBand(BandAsiaPacific).Key; got != BandAsiaPacific {
		t.Fatalf("known band = %q", got)
	}
	if ValidBand("mars") {
		t.Fatal("ValidBand accepted a nonsense band")
	}
}

func TestWindowIsLocalMorningToEvening(t *testing.T) {
	day := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	for _, b := range Bands() {
		opens, closes := Window(day, b)
		if !opens.Before(closes) {
			t.Fatalf("%s: window opens (%s) after it closes (%s)", b.Key, opens, closes)
		}
		// Converting back to band-local time must give the configured hours.
		off := time.Duration(b.OffsetHours) * time.Hour
		if h := opens.Add(off).Hour(); h != OpenHourLocal {
			t.Fatalf("%s: opens at %d local, want %d", b.Key, h, OpenHourLocal)
		}
		if h := closes.Add(off).Hour(); h != CloseHourLocal {
			t.Fatalf("%s: closes at %d local, want %d", b.Key, h, CloseHourLocal)
		}
		// The window must be wide — a narrow one silently excludes people.
		if d := closes.Sub(opens); d < 12*time.Hour {
			t.Fatalf("%s: window is only %s wide", b.Key, d)
		}
	}
}

func TestPromptRotatesDeterministically(t *testing.T) {
	day := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	first := Prompt(day)
	if first == "" {
		t.Fatal("empty prompt")
	}
	// Stable for the same day, so a rerun or a second cohort in the band matches.
	for i := 0; i < 5; i++ {
		if got := Prompt(day); got != first {
			t.Fatalf("prompt changed between calls: %q then %q", first, got)
		}
	}
	// And it actually moves day to day.
	if Prompt(day.AddDate(0, 0, 1)) == first {
		t.Fatal("prompt did not change the next day")
	}
	// Rotation covers the whole set within its length.
	seen := map[string]bool{}
	for i := 0; i < len(prompts); i++ {
		seen[Prompt(day.AddDate(0, 0, i))] = true
	}
	if len(seen) != len(prompts) {
		t.Fatalf("rotation produced %d distinct prompts over %d days, want %d",
			len(seen), len(prompts), len(prompts))
	}
}

func TestLocalDayShiftsWithBand(t *testing.T) {
	// 23:00 UTC is already tomorrow in Asia/Pacific and still today in the
	// Americas — the whole reason cohorts carry a band.
	at := time.Date(2026, 9, 3, 23, 0, 0, 0, time.UTC)
	asia := LocalDay(at, LookupBand(BandAsiaPacific))
	americas := LocalDay(at, LookupBand(BandAmericas))
	if !asia.After(americas) {
		t.Fatalf("Asia/Pacific local day (%s) should be ahead of Americas (%s)",
			asia.Format("2006-01-02"), americas.Format("2006-01-02"))
	}
}
