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

func TestPassAndFoundationsBoundaries(t *testing.T) {
	tests := []struct {
		name                string
		score               int
		passed              bool
		foundationsRequired bool
		recommendedPath     string
	}{
		{
			name:            "49 fails without a route",
			score:           PassThreshold - 1,
			recommendedPath: "",
		},
		{
			name:                "50 passes into foundations",
			score:               PassThreshold,
			passed:              true,
			foundationsRequired: true,
			recommendedPath:     PathFoundations,
		},
		{
			name:                "69 passes into foundations",
			score:               FoundationsThreshold - 1,
			passed:              true,
			foundationsRequired: true,
			recommendedPath:     PathFoundations,
		},
		{
			name:            "70 passes into a specialization",
			score:           FoundationsThreshold,
			passed:          true,
			recommendedPath: PathFullstack,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Decide(core(tt.score))
			if d.Passed != tt.passed {
				t.Errorf("Passed = %t, want %t", d.Passed, tt.passed)
			}
			if d.FoundationsRequired != tt.foundationsRequired {
				t.Errorf("FoundationsRequired = %t, want %t", d.FoundationsRequired, tt.foundationsRequired)
			}
			if d.RecommendedPath != tt.recommendedPath {
				t.Errorf("RecommendedPath = %q, want %q", d.RecommendedPath, tt.recommendedPath)
			}
			if len(d.Exemptions) != 0 {
				t.Errorf("Exemptions = %v, want none", d.Exemptions)
			}
		})
	}
}

func TestEmptyScoresAreTreatedAsZero(t *testing.T) {
	// An unanswered topic is not evidence of competence. A learner who submits a
	// blank paper must not be routed as if they had passed.
	d := Decide(map[string]int{})
	if d.Passed {
		t.Fatal("an empty score map passed")
	}
	if d.FoundationsRequired {
		t.Fatal("a failed empty attempt required foundations instead of remaining unrouted")
	}
	if d.RecommendedPath != "" {
		t.Fatalf("RecommendedPath = %q, want no route", d.RecommendedPath)
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

func TestExemptionsRequirePassedAttemptAndTopicMastery(t *testing.T) {
	// A bare pass must not exempt anything — skipping is meant to require clear
	// command of the topic.
	m := core(PassThreshold)
	if d := Decide(m); len(d.Exemptions) != 0 {
		t.Fatalf("a bare pass produced exemptions: %v", d.Exemptions)
	}

	// A passed learner can earn an exemption while still being routed through
	// foundations.
	m = core(PassThreshold)
	m[TopicDataStructures] = ExemptionThreshold
	d := Decide(m)
	if !d.Passed {
		t.Fatal("expected attempt with a passing core mean to pass")
	}
	if !d.FoundationsRequired {
		t.Fatal("expected foundations to still be required")
	}
	if !reflect.DeepEqual(d.Exemptions, []string{"data_structures"}) {
		t.Fatalf("exemptions = %v, want [data_structures]", d.Exemptions)
	}

	// Topic mastery cannot earn an exemption when the overall attempt failed.
	m = core(0)
	m[TopicDataStructures] = 100
	d = Decide(m)
	if d.Passed {
		t.Fatal("expected attempt with one mastered topic and a failing core mean to fail")
	}
	if d.RecommendedPath != "" {
		t.Fatalf("failed attempt recommended %q, want no route", d.RecommendedPath)
	}
	if len(d.Exemptions) != 0 {
		t.Fatalf("failed attempt produced exemptions: %v", d.Exemptions)
	}
}

func TestProgrammingExemptsBothLanguageCourses(t *testing.T) {
	m := core(PassThreshold)
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

func TestSpecializationTopicsDoNotAffectTheGates(t *testing.T) {
	// Not knowing HTTP says nothing about whether you understand a hash table,
	// and must not fail an attempt or push someone into remedial CS.
	strong := core(90)
	strong[TopicWeb], strong[TopicBackend] = 0, 0
	if d := Decide(strong); !d.Passed || d.FoundationsRequired {
		t.Fatalf("zero specialization scores changed strong core outcome: %+v", d)
	}

	weak := core(PassThreshold - 1)
	weak[TopicWeb], weak[TopicBackend] = 100, 100
	if d := Decide(weak); d.Passed || d.RecommendedPath != "" {
		t.Fatalf("specialization scores changed failing core outcome: %+v", d)
	}
}
