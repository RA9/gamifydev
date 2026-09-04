// Package cohort holds the grouping and scheduling policy: how learners are
// banded by timezone, how big a cohort should be, when a cohort is too small to
// survive, and when each day's standup window opens and closes.
//
// Pure, like internal/placement, because these are the numbers most likely to
// move once real enrollment exists — and because the interaction between "form
// at 6–8" and "drop aggressively" is exactly the tension the PRD flags (§03).
// Keeping the sizing rules in one testable place is what stops that tension
// from being rediscovered in a handler.
package cohort

import (
	"fmt"
	"time"
)

// Timezone bands. Coarse on purpose: the point is that a cohort's standup window
// overlaps everyone's waking hours, not that members share a city.
const (
	BandAmericas     = "americas"
	BandEuropeAfrica = "europe_africa"
	BandAsiaPacific  = "asia_pacific"
)

// Band describes one band's representative local time.
type Band struct {
	Key   string
	Label string
	// OffsetHours is the representative UTC offset used to turn "06:00 local"
	// into an instant. Members inside a band differ by a few hours; the window
	// is wide enough to absorb that.
	OffsetHours int
}

var bands = []Band{
	{BandAmericas, "Americas (UTC-8 to UTC-3)", -5},
	{BandEuropeAfrica, "Europe & Africa (UTC-1 to UTC+3)", 1},
	{BandAsiaPacific, "Asia & Pacific (UTC+4 to UTC+12)", 8},
}

// Bands returns every selectable band, in display order.
func Bands() []Band { return append([]Band(nil), bands...) }

// LookupBand returns a band by key, falling back to Europe/Africa for an unknown
// or empty value so a missing preference can never block cohort formation.
func LookupBand(key string) Band {
	for _, b := range bands {
		if b.Key == key {
			return b
		}
	}
	return bands[1]
}

// ValidBand reports whether key names a real band.
func ValidBand(key string) bool {
	for _, b := range bands {
		if b.Key == key {
			return true
		}
	}
	return false
}

// Sizing. Cohorts form larger than their steady state because attrition on a
// free platform is heavy early: forming at four and dropping aggressively
// predictably produces a cohort of one.
const (
	// TargetSize is the size we form at when enough learners are queued.
	TargetSize = 6
	// MaxSize caps a single cohort so the standup stays readable.
	MaxSize = 8
	// MinViable is the active-member count below which a cohort is repacked
	// into another. Three is the point where a group stops feeling like one.
	MinViable = 3
)

// MaxQueueWait is how long a learner may sit unplaced before we form a cohort
// with whoever is available — including a cohort of one. Waiting for a full
// group forever is worse than starting small with a mentor.
const MaxQueueWait = 48 * time.Hour

// ShouldForm decides whether to open a cohort now for a queue of `waiting`
// learners whose longest wait is `oldest`.
//
// Two ways to qualify: enough people, or someone has waited long enough that
// starting matters more than group size.
func ShouldForm(waiting int, oldest time.Duration) bool {
	if waiting <= 0 {
		return false
	}
	return waiting >= TargetSize || oldest >= MaxQueueWait
}

// TakeSize is how many of a queue to place into one new cohort.
func TakeSize(waiting int) int {
	if waiting > MaxSize {
		return MaxSize
	}
	return waiting
}

// Standup window. Deliberately wide — the PRD's guidance is to start wide and
// narrow it only with evidence, because a tight window silently excludes people
// whose day does not line up.
const (
	// OpenHourLocal is when the window opens in the band's local time.
	OpenHourLocal = 6
	// CloseHourLocal is when it closes. 23 keeps the whole evening available.
	CloseHourLocal = 23
)

// Window returns the UTC instants a cohort's standup opens and closes for a
// given local day.
//
// `day` is the cohort's local date. Both instants are returned in UTC because
// that is what the database stores and compares against.
func Window(day time.Time, band Band) (opens, closes time.Time) {
	off := time.Duration(band.OffsetHours) * time.Hour
	y, m, d := day.Date()
	localOpen := time.Date(y, m, d, OpenHourLocal, 0, 0, 0, time.UTC)
	localClose := time.Date(y, m, d, CloseHourLocal, 59, 0, 0, time.UTC)
	// Subtracting the offset converts "this wall-clock time in the band" to UTC.
	return localOpen.Add(-off), localClose.Add(-off)
}

// LocalDay returns the band-local date for a UTC instant. Used so a cohort's
// "today" matches its members' day rather than the server's.
func LocalDay(now time.Time, band Band) time.Time {
	return now.UTC().Add(time.Duration(band.OffsetHours) * time.Hour).Truncate(24 * time.Hour)
}

// prompts are the rotating daily standup questions the bot mentor posts.
//
// Templated and deterministic on purpose: facilitating a standup does not need a
// model, and a per-cohort-per-day model call is a recurring cost with no revenue
// behind it (PRD §06). The escalation path — a learner who reports a blocker —
// is where a model earns its keep.
var prompts = []string{
	"What did you get working yesterday, and what's the first thing you'll touch today?",
	"Where did you spend the most time yesterday? If it was longer than you expected, say why.",
	"What's one thing you understand better today than you did yesterday?",
	"What are you working on today, and what would make it go faster?",
	"Did anything block you yesterday? Even a small snag is worth naming.",
	"What's the next milestone you're aiming for, and what stands between you and it?",
	"What did you learn yesterday that you'd explain to someone starting this week?",
}

// Prompt returns the prompt for a given local day, rotating deterministically so
// every cohort in a band sees the same question and reruns are stable.
func Prompt(day time.Time) string {
	// Days since the Unix epoch keeps the rotation stable across restarts.
	idx := int(day.UTC().Unix()/86400) % len(prompts)
	if idx < 0 {
		idx += len(prompts)
	}
	return prompts[idx]
}

// Name builds a human cohort name, e.g. "Frontend Developer · Europe & Africa ·
// 3 Sep". Used in listings and the cohort header.
func Name(pathTitle string, band Band, start time.Time) string {
	return fmt.Sprintf("%s · %s · %s", pathTitle, shortBand(band), start.Format("2 Jan"))
}

func shortBand(b Band) string {
	switch b.Key {
	case BandAmericas:
		return "Americas"
	case BandAsiaPacific:
		return "Asia & Pacific"
	default:
		return "Europe & Africa"
	}
}
