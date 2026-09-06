// Package placement holds the diagnostic policy: given per-topic scores, decide
// whether a learner passed, which path they start on, and which courses they may
// skip in the schedule.
//
// This is deliberately pure — no database, no HTTP — because it is the piece
// most likely to be re-tuned once real attempt data exists, and it needs to be
// cheap to test at the boundaries.
//
// The governing decision (PRD §03) is that the diagnostic *routes* and the
// course checkpoint *exempts*. Nothing here grants credit for a course; an
// exemption only removes a course from the learner's schedule, and the content
// stays readable. That is what makes cheating the test self-defeating rather
// than rewarding.
package placement

import "sort"

// Topics the item bank is tagged with. The first five are "core CS" and decide
// whether the attempt passes and whether foundations is required; the last two
// only inform which specialization to recommend.
const (
	TopicProgramming    = "programming"
	TopicDataStructures = "data_structures"
	TopicComplexity     = "complexity"
	TopicAlgorithms     = "algorithms"
	TopicComputers      = "computers"
	TopicWeb            = "web"
	TopicBackend        = "backend"
)

// CoreTopics decide the pass and foundations gates. Specialization topics are
// excluded on purpose: not knowing HTTP says nothing about whether you understand
// a hash table, and it must not fail an attempt or require remedial CS.
var CoreTopics = []string{
	TopicProgramming,
	TopicDataStructures,
	TopicComplexity,
	TopicAlgorithms,
	TopicComputers,
}

// AllTopics is every topic the bank uses, in presentation order.
var AllTopics = append(append([]string{}, CoreTopics...), TopicWeb, TopicBackend)

// Thresholds. Expressed as percentages so they read the way we talk about them.
const (
	// PassThreshold is the minimum mean core-topic score for a passing attempt.
	// Failed attempts are not routed and cannot earn exemptions.
	PassThreshold = 50

	// FoundationsThreshold is the mean core-topic score at or above which a
	// passing learner may enter a specialization. Passing learners below it must
	// start with CS Foundations.
	FoundationsThreshold = 70

	// ExemptionThreshold is the per-topic score at or above which the courses
	// backed by that topic are proposed as skippable in the schedule.
	// Deliberately well above the pass mark: skipping should require clear
	// command of the topic, not a bare pass.
	ExemptionThreshold = 85

	// SpecializationMargin is how far ahead one specialization score must be
	// before we recommend it over the other. Inside the margin we read the
	// learner as balanced and suggest full-stack.
	SpecializationMargin = 15
)

// topicCourses maps a topic to the course slugs an exemption in that topic may
// skip. Linux is intentionally absent: nothing in the diagnostic tests it, so
// it is never exempted.
var topicCourses = map[string][]string{
	TopicProgramming:    {"c"},
	TopicDataStructures: {"data_structures"},
	TopicComplexity:     {"complexity_and_analysis"},
	TopicAlgorithms:     {"algorithms"},
	TopicComputers:      {"how_computers_work"},
}

// Path slugs this policy can route to. They must exist in the seeded paths.
const (
	PathFoundations = "cs-foundations"
	PathFrontend    = "frontend-developer"
	PathBackend     = "backend-developer"
	PathFullstack   = "fullstack-developer"
)

// Decision is the outcome of evaluating and routing an attempt.
type Decision struct {
	// CoreScore is the mean of the core-topic scores, 0-100.
	CoreScore int
	// Passed reports whether CoreScore cleared PassThreshold.
	Passed bool
	// FoundationsRequired locks a passing learner to CS Foundations.
	FoundationsRequired bool
	// RecommendedPath is the path slug to start on. It is empty for failed
	// attempts and PathFoundations when FoundationsRequired is true.
	RecommendedPath string
	// Exemptions are course slugs a passing learner may skip in the schedule.
	// Empty unless the attempt passed and the corresponding topic cleared
	// ExemptionThreshold.
	Exemptions []string
	// Summary is a short human explanation shown on the result page.
	Summary string
}

// Decide turns per-topic percentage scores into a routing decision.
//
// A topic missing from `scores` counts as zero: an unanswered topic cannot be
// evidence of competence in it.
func Decide(scores map[string]int) Decision {
	d := Decision{CoreScore: MeanCore(scores)}
	d.Passed = d.CoreScore >= PassThreshold
	if !d.Passed {
		d.Summary = "Your core score did not meet the passing threshold. Review the fundamentals and try the diagnostic again."
		return d
	}

	d.FoundationsRequired = d.CoreScore < FoundationsThreshold

	// Exemptions are earned per topic only after the overall attempt passes.
	seen := map[string]bool{}
	for _, topic := range CoreTopics {
		if scores[topic] < ExemptionThreshold {
			continue
		}
		for _, course := range topicCourses[topic] {
			if !seen[course] {
				seen[course] = true
				d.Exemptions = append(d.Exemptions, course)
			}
		}
	}
	sort.Strings(d.Exemptions)

	if d.FoundationsRequired {
		d.RecommendedPath = PathFoundations
		d.Summary = "Start with Computer Science Foundations to build the fundamentals the other paths assume."
		return d
	}

	web, backend := scores[TopicWeb], scores[TopicBackend]
	switch {
	case web-backend > SpecializationMargin:
		d.RecommendedPath = PathFrontend
		d.Summary = "Your fundamentals are solid and you're strongest on the browser side."
	case backend-web > SpecializationMargin:
		d.RecommendedPath = PathBackend
		d.Summary = "Your fundamentals are solid and you're strongest on the server side."
	default:
		d.RecommendedPath = PathFullstack
		d.Summary = "Your fundamentals are solid and you're balanced across the stack."
	}
	return d
}

// MeanCore averages the core-topic scores. Missing topics count as zero.
func MeanCore(scores map[string]int) int {
	if len(CoreTopics) == 0 {
		return 0
	}
	total := 0
	for _, t := range CoreTopics {
		total += scores[t]
	}
	return total / len(CoreTopics)
}

// ExemptibleCourses reports every course slug the policy is capable of
// exempting. Used to validate that the mapping matches seeded content.
func ExemptibleCourses() []string {
	var out []string
	for _, courses := range topicCourses {
		out = append(out, courses...)
	}
	sort.Strings(out)
	return out
}

// TopicLabel renders a topic for display.
func TopicLabel(topic string) string {
	switch topic {
	case TopicProgramming:
		return "Programming basics"
	case TopicDataStructures:
		return "Data structures"
	case TopicComplexity:
		return "Complexity"
	case TopicAlgorithms:
		return "Algorithms"
	case TopicComputers:
		return "How computers work"
	case TopicWeb:
		return "Web and the browser"
	case TopicBackend:
		return "Servers and data"
	}
	return topic
}

// IsCore reports whether a topic counts toward the foundations gate.
func IsCore(topic string) bool {
	for _, t := range CoreTopics {
		if t == topic {
			return true
		}
	}
	return false
}
