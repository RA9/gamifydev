package placement

import (
	"reflect"
	"testing"
)

// core builds a score map where every core topic has the same value.
func core(v int) map[string]int {
	m := map[string]int{}
	for _, t := range CoreTopics {
		m[t] = v
	}
	return m
}

func TestFoundationsGateBoundary(t *testing.T) {
	// The threshold is inclusive: exactly at the mark must open the paths, since
	// the copy tells learners "70% opens the other paths".
	if d := Decide(core(FoundationsThreshold)); d.FoundationsRequired {
		t.Fatalf("score of exactly %d required foundations; the threshold should be inclusive", FoundationsThreshold)
	}
	if d := Decide(core(FoundationsThreshold - 1)); !d.FoundationsRequired {
		t.Fatalf("score of %d did not require foundations", FoundationsThreshold-1)
	}
	if d := Decide(core(FoundationsThreshold)); d.RecommendedPath == PathFoundations {
		t.Fatal("a passing learner was still routed to foundations")
	}
}

func TestFailingRoutesToFoundations(t *testing.T) {
	d := Decide(core(20))
	if !d.FoundationsRequired {
		t.Fatal("a weak score did not require foundations")
	}
	if d.RecommendedPath != PathFoundations {
		t.Fatalf("recommended %q, want %q", d.RecommendedPath, PathFoundations)
	}
}

func TestEmptyScoresAreTreatedAsZero(t *testing.T) {
	// An unanswered topic is not evidence of competence. A learner who submits a
	// blank paper must not be routed as if they had passed.
	d := Decide(map[string]int{})
	if !d.FoundationsRequired {
		t.Fatal("an empty score map did not require foundations")
	}
	if d.CoreScore != 0 {
		t.Fatalf("CoreScore = %d, want 0", d.CoreScore)
	}
	if len(d.Exemptions) != 0 {
		t.Fatalf("empty scores produced exemptions: %v", d.Exemptions)
	}
}

func TestSpecializationRouting(t *testing.T) {
	pass := FoundationsThreshold + 10
	mk := func(web, backend int) map[string]int {
		m := core(pass)
		m[TopicWeb], m[TopicBackend] = web, backend
		return m
	}

	if got := Decide(mk(90, 40)).RecommendedPath; got != PathFrontend {
		t.Fatalf("web-dominant routed to %q, want %q", got, PathFrontend)
	}
	if got := Decide(mk(40, 90)).RecommendedPath; got != PathBackend {
		t.Fatalf("backend-dominant routed to %q, want %q", got, PathBackend)
	}
	// Inside the margin reads as balanced, whichever side is nominally ahead.
	if got := Decide(mk(80, 80)).RecommendedPath; got != PathFullstack {
		t.Fatalf("balanced routed to %q, want %q", got, PathFullstack)
	}
	if got := Decide(mk(80, 80-SpecializationMargin)).RecommendedPath; got != PathFullstack {
		t.Fatalf("exactly at the margin routed to %q, want %q (margin is exclusive)", got, PathFullstack)
	}
	if got := Decide(mk(80, 80-SpecializationMargin-1)).RecommendedPath; got != PathFrontend {
		t.Fatalf("just past the margin routed to %q, want %q", got, PathFrontend)
	}
}

func TestExemptionsRequireMasteryNotAPass(t *testing.T) {
	// A bare pass must not exempt anything — skipping is meant to require clear
	// command of the topic.
	m := core(FoundationsThreshold)
	if d := Decide(m); len(d.Exemptions) != 0 {
		t.Fatalf("a bare pass produced exemptions: %v", d.Exemptions)
	}

	m = core(0)
	m[TopicDataStructures] = ExemptionThreshold
	d := Decide(m)
	if !reflect.DeepEqual(d.Exemptions, []string{"data_structures"}) {
		t.Fatalf("exemptions = %v, want [data_structures]", d.Exemptions)
	}
	// Crucially, exemptions are earned even when foundations are required: a
	// learner sent to foundations who clearly knows one topic shouldn't resit it.
	if !d.FoundationsRequired {
		t.Fatal("expected foundations to still be required")
	}
}

func TestProgrammingExemptsBothLanguageCourses(t *testing.T) {
	m := core(0)
	m[TopicProgramming] = 100
	got := Decide(m).Exemptions
	want := []string{"c", "java"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("exemptions = %v, want %v", got, want)
	}
}

func TestLinuxIsNeverExempted(t *testing.T) {
	// Nothing in the diagnostic tests the command line, so it must never be
	// skipped on the strength of it.
	full := map[string]int{}
	for _, t := range AllTopics {
		full[t] = 100
	}
	for _, c := range Decide(full).Exemptions {
		if c == "linux" {
			t.Fatal("linux was exempted, but no diagnostic topic covers it")
		}
	}
}

func TestExemptionsAreDeterministicAndUnique(t *testing.T) {
	full := map[string]int{}
	for _, t := range AllTopics {
		full[t] = 100
	}
	first := Decide(full).Exemptions
	for i := 0; i < 20; i++ {
		if got := Decide(full).Exemptions; !reflect.DeepEqual(got, first) {
			t.Fatalf("exemptions vary between calls: %v then %v", first, got)
		}
	}
	seen := map[string]bool{}
	for _, c := range first {
		if seen[c] {
			t.Fatalf("duplicate exemption %q in %v", c, first)
		}
		seen[c] = true
	}
}

func TestSpecializationTopicsDoNotAffectTheGate(t *testing.T) {
	// Not knowing HTTP says nothing about whether you understand a hash table,
	// and must not push someone into remedial CS.
	strong := core(90)
	strong[TopicWeb], strong[TopicBackend] = 0, 0
	if d := Decide(strong); d.FoundationsRequired {
		t.Fatal("zero specialization scores forced a strong learner into foundations")
	}
}
