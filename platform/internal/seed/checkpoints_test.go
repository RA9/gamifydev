package seed

import (
	"os/exec"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// A checkpoint gates the next course in a path a learner cannot opt out of, so
// two things have to be true of every one we ship, and neither is obvious by
// reading it: a correct solution passes all of it, and a wrong solution is
// actually caught. A check that nothing can fail is decoration.
//
// These are the reference solutions. They live in the test rather than beside
// the prompts so there is no path by which they reach a learner.
type refSolution struct {
	correct string
	// wrong is a plausible near-miss — the mistake this checkpoint exists to
	// catch — and wrongFails names the check that must catch it.
	wrong      string
	wrongFails string
}

var refSolutions = map[string]refSolution{
	"py-sum-list": {
		correct: `
def sum_list(nums):
    total = 0
    for n in nums:
        total += n
    return total
`,
		// Sums by draining the list, which is right until someone needs their
		// list afterwards.
		wrong: `
def sum_list(nums):
    total = 0
    while nums:
        total += nums.pop()
    return total
`,
		wrongFails: "leaves the caller's list alone",
	},

	"c-max-three": {
		correct: `
#include <stdio.h>

int main(void) {
    int a, b, c;
    if (scanf("%d %d %d", &a, &b, &c) != 3) return 1;
    int m = a;
    if (b > m) m = b;
    if (c > m) m = c;
    printf("%d\n", m);
    return 0;
}
`,
		// Seeding the running maximum with 0 instead of the first input: right
		// for every positive example, wrong for every negative one.
		wrong: `
#include <stdio.h>

int main(void) {
    int a, b, c;
    if (scanf("%d %d %d", &a, &b, &c) != 3) return 1;
    int m = 0;
    if (a > m) m = a;
    if (b > m) m = b;
    if (c > m) m = c;
    printf("%d\n", m);
    return 0;
}
`,
		wrongFails: "handles three negative numbers",
	},

	"ds-hash-map": {
		correct: `
class HashMap:
    def __init__(self, size=8):
        self.buckets = [[] for _ in range(size)]
        self.count = 0

    def _bucket(self, key):
        return self.buckets[hash(key) % len(self.buckets)]

    def put(self, key, value):
        b = self._bucket(key)
        for i, (k, _) in enumerate(b):
            if k == key:
                b[i] = (key, value)
                return
        b.append((key, value))
        self.count += 1

    def get(self, key):
        for k, v in self._bucket(key):
            if k == key:
                return v
        return None

    def __len__(self):
        return self.count
`,
		// Appends unconditionally, so a repeated key leaves two entries and get
		// keeps returning the stale one.
		wrong: `
class HashMap:
    def __init__(self, size=8):
        self.buckets = [[] for _ in range(size)]
        self.count = 0

    def _bucket(self, key):
        return self.buckets[hash(key) % len(self.buckets)]

    def put(self, key, value):
        self._bucket(key).append((key, value))
        self.count += 1

    def get(self, key):
        for k, v in self._bucket(key):
            if k == key:
                return v
        return None

    def __len__(self):
        return self.count
`,
		wrongFails: "putting a key twice replaces the value instead of adding a second copy",
	},

	"cx-pair-sum": {
		correct: `
def has_pair_summing_to(nums, target):
    seen = set()
    for n in nums:
        if target - n in seen:
            return True
        seen.add(n)
    return False
`,
		// The whole point of the checkpoint: correct, and far too slow.
		wrong: `
def has_pair_summing_to(nums, target):
    for i in range(len(nums)):
        for j in range(i + 1, len(nums)):
            if nums[i] + nums[j] == target:
                return True
    return False
`,
		wrongFails: "finishes on 200,000 numbers — a nested loop won't",
	},

	"algo-shortest-path": {
		correct: `
from collections import deque


def shortest_path(graph, start, goal):
    if start == goal:
        return [start]
    seen = {start}
    queue = deque([[start]])
    while queue:
        path = queue.popleft()
        for nxt in graph.get(path[-1], []):
            if nxt == goal:
                return path + [nxt]
            if nxt not in seen:
                seen.add(nxt)
                queue.append(path + [nxt])
    return None
`,
		// Depth-first: finds a route, not the shortest one.
		wrong: `
def shortest_path(graph, start, goal):
    def go(node, path, seen):
        if node == goal:
            return path
        for nxt in graph.get(node, []):
            if nxt not in seen:
                found = go(nxt, path + [nxt], seen | {nxt})
                if found:
                    return found
        return None

    return go(start, [start], {start})
`,
		wrongFails: "takes the short route, not the first one it finds",
	},

	"hcw-twos-complement": {
		correct: `
def to_twos_complement(n, bits):
    if n < 0:
        n += 1 << bits
    return format(n, "0" + str(bits) + "b")[-bits:]


def from_twos_complement(s):
    n = int(s, 2)
    if s[0] == "1":
        n -= 1 << len(s)
    return n
`,
		// Sign-magnitude: the intuitive encoding, and not the one hardware uses.
		wrong: `
def to_twos_complement(n, bits):
    return format(abs(n), "0" + str(bits) + "b")


def from_twos_complement(s):
    return int(s, 2)
`,
		wrongFails: "encodes a negative number",
	},

	"linux-top-talkers": {
		correct: `#!/usr/bin/env bash
awk '{print $1}' access.log | sort | uniq -c | sort -rn | head -3
`,
		// Reads the expected output off the page and prints it back. Every
		// visible check passes; only the hidden one notices.
		wrong: `#!/usr/bin/env bash
printf '5 10.0.0.1\n3 10.0.0.2\n2 10.0.0.3\n'
`,
		wrongFails: "counts the log instead of printing the answer",
	},
}

func newCheckpointExecutor(t *testing.T) runner.Executor {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3 on PATH")
	}
	if _, err := exec.LookPath("cc"); err != nil {
		t.Skip("no C compiler on PATH")
	}
	e, err := runner.New(runner.Config{
		Mode: "local", PythonPath: "python3", CCPath: "cc", AllowUnsafe: true,
	})
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	return e
}

// materialize turns an authored checkpoint into the shape the grader takes,
// giving each check the row id it would have had in the database.
func materialize(a seedA) ([]store.Check, map[string]string) {
	checks := make([]store.Check, len(a.checks))
	for i, c := range a.checks {
		checks[i] = store.Check{
			ID: int64(i + 1), Sort: i, Label: c.label, Test: c.test,
			Stdin: c.stdin, Hidden: c.hidden, Points: 1,
		}
	}
	files := map[string]string{}
	for _, f := range a.files {
		files[f.Name] = f.Content
	}
	return checks, files
}

func TestSeededCheckpointsAcceptACorrectSolution(t *testing.T) {
	e := newCheckpointExecutor(t)
	for _, a := range sampleAssignments {
		if len(a.checks) == 0 {
			continue
		}
		t.Run(a.slug, func(t *testing.T) {
			t.Parallel()
			ref, ok := refSolutions[a.slug]
			if !ok {
				t.Fatalf("checkpoint %q has checks but no reference solution — "+
					"nothing proves a learner can pass it", a.slug)
			}
			checks, files := materialize(a)
			res, err := jobs.Grade(t.Context(), e, a.lang, ref.correct, files, checks)
			if err != nil {
				t.Fatalf("grade: %v", err)
			}
			passed, output := res.Passed, res.Output
			for _, c := range checks {
				if !passed[c.ID] {
					t.Errorf("a correct solution failed %q\nprogram output:\n%s", c.Label, output)
				}
			}
		})
	}
}

func TestSeededCheckpointsRejectTheMistakeTheyExistToCatch(t *testing.T) {
	e := newCheckpointExecutor(t)
	for _, a := range sampleAssignments {
		if len(a.checks) == 0 {
			continue
		}
		t.Run(a.slug, func(t *testing.T) {
			t.Parallel()
			ref := refSolutions[a.slug]
			checks, files := materialize(a)
			res, err := jobs.Grade(t.Context(), e, a.lang, ref.wrong, files, checks)
			if err != nil {
				t.Fatalf("grade: %v", err)
			}
			passed, output := res.Passed, res.Output
			var target *store.Check
			for i := range checks {
				if checks[i].Label == ref.wrongFails {
					target = &checks[i]
				}
			}
			if target == nil {
				t.Fatalf("no check labelled %q — the reference solution names a "+
					"check that was renamed or removed", ref.wrongFails)
			}
			if passed[target.ID] {
				t.Errorf("%q passed a solution that should fail it\nprogram output:\n%s",
					target.Label, output)
			}
			// The near-miss must be a near miss. If it fails everything it is
			// just a broken program, and proves nothing about this one check.
			anyPassed := false
			for _, c := range checks {
				if passed[c.ID] {
					anyPassed = true
				}
			}
			if !anyPassed && !res.CompileFailed {
				t.Errorf("the wrong solution failed every check, so %q wasn't "+
					"shown to be what caught it\nprogram output:\n%s", ref.wrongFails, output)
			}
		})
	}
}

// checkpointsByCourse indexes the seeded checkpoints by the course they belong
// to. A slice, not a single value, so a course that grows a second checkpoint
// doesn't silently hide one from these checks.
func checkpointsByCourse() map[string][]seedA {
	byCourse := map[string][]seedA{}
	for _, a := range sampleAssignments {
		byCourse[a.course] = append(byCourse[a.course], a)
	}
	return byCourse
}

// gradesItself mirrors what ReplaceChecks decides: a checkpoint grades itself
// only when the sandbox can run its language and it actually carries checks.
func gradesItself(a seedA) bool {
	return store.Runnable(a.lang) && len(a.checks) > 0
}

// Every course in the compulsory path needs a checkpoint. A path a learner is
// locked into, with a course that lets anyone through, is the one gap that
// matters.
func TestEveryCSFoundationsCourseHasACheckpoint(t *testing.T) {
	var foundations []string
	for _, p := range seedPaths {
		if p.Slug == "cs-foundations" {
			foundations = p.Courses
		}
	}
	if len(foundations) == 0 {
		t.Fatal("the cs-foundations path is gone or has no courses")
	}
	byCourse := checkpointsByCourse()
	for _, course := range foundations {
		as, ok := byCourse[course]
		if !ok {
			t.Errorf("course %q in the compulsory path has no checkpoint", course)
			continue
		}
		for _, a := range as {
			// Where the sandbox can run the language, the checks have to exist —
			// a runnable course left to mentor grading is work handed to a human
			// for no reason.
			if store.Runnable(a.lang) && len(a.checks) == 0 {
				t.Errorf("checkpoint %q is in %s, which the sandbox can run, but has no checks",
					a.slug, a.lang)
			}
		}
	}
}

// A gating checkpoint holds back every course after it in the path. If that
// gate can't grade itself, one mentor's availability stands in front of all of
// them — which is what put Java, at position two of seven, in front of five
// auto-graded courses in the one path nobody can opt out of.
//
// The last course in a path is exempt: its gate holds nothing back.
func TestNoMidPathGateWaitsOnAHuman(t *testing.T) {
	byCourse := checkpointsByCourse()
	for _, p := range seedPaths {
		for i, course := range p.Courses {
			behind := len(p.Courses) - i - 1
			if behind == 0 {
				continue
			}
			for _, a := range byCourse[course] {
				if a.elective || gradesItself(a) {
					continue
				}
				t.Errorf("%s: %q gates on %s, which the sandbox can't grade, so a "+
					"mentor stands in front of the %d course(s) behind it — make it "+
					"elective, or give it checks in a language the sandbox runs",
					p.Slug, a.slug, a.lang, behind)
			}
		}
	}
}
