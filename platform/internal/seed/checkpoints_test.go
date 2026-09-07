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

	"ds-hash-map-c": {
		correct: `
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUCKET_COUNT 8

typedef struct Entry {
    int key;
    int value;
    struct Entry *next;
} Entry;

typedef struct {
    Entry *buckets[BUCKET_COUNT];
    int size;
} HashMap;

static int bucket_index(int key) {
    int index = key % BUCKET_COUNT;
    return index < 0 ? index + BUCKET_COUNT : index;
}

static void put(HashMap *map, int key, int value) {
    int index = bucket_index(key);
    for (Entry *entry = map->buckets[index]; entry; entry = entry->next) {
        if (entry->key == key) {
            entry->value = value;
            return;
        }
    }
    Entry *entry = malloc(sizeof(*entry));
    if (!entry) exit(1);
    entry->key = key;
    entry->value = value;
    entry->next = map->buckets[index];
    map->buckets[index] = entry;
    map->size++;
}

static int get(const HashMap *map, int key, int *value) {
    for (Entry *entry = map->buckets[bucket_index(key)]; entry; entry = entry->next) {
        if (entry->key == key) {
            *value = entry->value;
            return 1;
        }
    }
    return 0;
}

int main(void) {
    HashMap map = {0};
    int operations;
    if (scanf("%d", &operations) != 1) return 1;
    for (int i = 0; i < operations; i++) {
        char command[8];
        if (scanf("%7s", command) != 1) return 1;
        if (strcmp(command, "PUT") == 0) {
            int key, value;
            if (scanf("%d %d", &key, &value) != 2) return 1;
            put(&map, key, value);
        } else if (strcmp(command, "GET") == 0) {
            int key, value;
            if (scanf("%d", &key) != 1) return 1;
            if (get(&map, key, &value)) printf("%d\n", value);
            else printf("NOT_FOUND\n");
        } else if (strcmp(command, "SIZE") == 0) {
            printf("%d\n", map.size);
        }
    }
    return 0;
}
`,
		// Inserts duplicates at the end of a chain instead of replacing them, so
		// GET sees the stale value and SIZE counts the same key twice.
		wrong: `
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUCKET_COUNT 8

typedef struct Entry {
    int key;
    int value;
    struct Entry *next;
} Entry;

typedef struct {
    Entry *buckets[BUCKET_COUNT];
    int size;
} HashMap;

static int bucket_index(int key) {
    int index = key % BUCKET_COUNT;
    return index < 0 ? index + BUCKET_COUNT : index;
}

static void put(HashMap *map, int key, int value) {
    int index = bucket_index(key);
    Entry *entry = malloc(sizeof(*entry));
    if (!entry) exit(1);
    entry->key = key;
    entry->value = value;
    entry->next = NULL;
    Entry **tail = &map->buckets[index];
    while (*tail) tail = &(*tail)->next;
    *tail = entry;
    map->size++;
}

static int get(const HashMap *map, int key, int *value) {
    for (Entry *entry = map->buckets[bucket_index(key)]; entry; entry = entry->next) {
        if (entry->key == key) {
            *value = entry->value;
            return 1;
        }
    }
    return 0;
}

int main(void) {
    HashMap map = {0};
    int operations;
    if (scanf("%d", &operations) != 1) return 1;
    for (int i = 0; i < operations; i++) {
        char command[8];
        if (scanf("%7s", command) != 1) return 1;
        if (strcmp(command, "PUT") == 0) {
            int key, value;
            if (scanf("%d %d", &key, &value) != 2) return 1;
            put(&map, key, value);
        } else if (strcmp(command, "GET") == 0) {
            int key, value;
            if (scanf("%d", &key) != 1) return 1;
            if (get(&map, key, &value)) printf("%d\n", value);
            else printf("NOT_FOUND\n");
        } else if (strcmp(command, "SIZE") == 0) {
            printf("%d\n", map.size);
        }
    }
    return 0;
}
`,
		wrongFails: "putting a key twice replaces the value instead of adding a second copy",
	},

	"cx-pair-sum-c": {
		correct: `
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

static size_t hash_int(int value, size_t capacity) {
    uint32_t x = (uint32_t)value;
    x ^= x >> 16;
    x *= 0x7feb352dU;
    x ^= x >> 15;
    return x & (capacity - 1);
}

static int contains(const int *keys, const unsigned char *used, size_t capacity, int value) {
    size_t index = hash_int(value, capacity);
    while (used[index]) {
        if (keys[index] == value) return 1;
        index = (index + 1) & (capacity - 1);
    }
    return 0;
}

static void insert(int *keys, unsigned char *used, size_t capacity, int value) {
    size_t index = hash_int(value, capacity);
    while (used[index] && keys[index] != value) {
        index = (index + 1) & (capacity - 1);
    }
    keys[index] = value;
    used[index] = 1;
}

int main(void) {
    int n, target;
    if (scanf("%d %d", &n, &target) != 2) return 1;
    size_t capacity = 1;
    while (capacity < (size_t)(n > 0 ? n : 1) * 2) capacity <<= 1;
    int *keys = malloc(capacity * sizeof(*keys));
    unsigned char *used = calloc(capacity, sizeof(*used));
    if (!keys || !used) return 1;

    int found = 0;
    for (int i = 0; i < n; i++) {
        int value;
        if (scanf("%d", &value) != 1) return 1;
        if (!found && contains(keys, used, capacity, target - value)) found = 1;
        insert(keys, used, capacity, value);
    }
    printf(found ? "YES\n" : "NO\n");
    return 0;
}
`,
		// Produces the right answers on ordinary inputs, but compares every pair
		// and cannot finish the deliberately worst-case 200,000-value check.
		wrong: `
#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, target;
    if (scanf("%d %d", &n, &target) != 2) return 1;
    int *values = malloc((size_t)n * sizeof(*values));
    if (n > 0 && !values) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%d", &values[i]) != 1) return 1;
    }
    for (int i = 0; i < n; i++) {
        for (int j = i + 1; j < n; j++) {
            if (values[i] + values[j] == target) {
                printf("YES\n");
                return 0;
            }
        }
    }
    printf("NO\n");
    return 0;
}
`,
		wrongFails: "finishes on 200,000 numbers — a nested loop won't",
	},

	"algo-shortest-path-c": {
		correct: `
#include <stdio.h>

#define MAX_VERTICES 100

int main(void) {
    int n, m, start, goal;
    if (scanf("%d %d %d %d", &n, &m, &start, &goal) != 4) return 1;
    int edges[MAX_VERTICES][MAX_VERTICES] = {{0}};
    for (int i = 0; i < m; i++) {
        int from, to;
        if (scanf("%d %d", &from, &to) != 2) return 1;
        edges[from][to] = 1;
    }

    int queue[MAX_VERTICES], parent[MAX_VERTICES], seen[MAX_VERTICES] = {0};
    int front = 0, back = 0;
    for (int i = 0; i < n; i++) parent[i] = -1;
    queue[back++] = start;
    seen[start] = 1;
    while (front < back && !seen[goal]) {
        int node = queue[front++];
        for (int next = 0; next < n; next++) {
            if (edges[node][next] && !seen[next]) {
                seen[next] = 1;
                parent[next] = node;
                queue[back++] = next;
            }
        }
    }

    if (!seen[goal]) {
        printf("NONE\n");
        return 0;
    }
    int path[MAX_VERTICES], length = 0;
    for (int node = goal; node != -1; node = parent[node]) path[length++] = node;
    for (int i = length - 1; i >= 0; i--) {
        if (i != length - 1) putchar(' ');
        printf("%d", path[i]);
    }
    putchar('\n');
    return 0;
}
`,
		// Depth-first search remembers visited vertices and returns a valid route,
		// but takes the first long branch instead of the direct shortest edge.
		wrong: `
#include <stdio.h>

#define MAX_VERTICES 100

static int n;
static int edges[MAX_VERTICES][MAX_VERTICES];
static int seen[MAX_VERTICES];
static int path[MAX_VERTICES];
static int path_length;

static int find_path(int node, int goal) {
    seen[node] = 1;
    path[path_length++] = node;
    if (node == goal) return 1;
    for (int next = 0; next < n; next++) {
        if (edges[node][next] && !seen[next] && find_path(next, goal)) return 1;
    }
    path_length--;
    return 0;
}

int main(void) {
    int m, start, goal;
    if (scanf("%d %d %d %d", &n, &m, &start, &goal) != 4) return 1;
    for (int i = 0; i < m; i++) {
        int from, to;
        if (scanf("%d %d", &from, &to) != 2) return 1;
        edges[from][to] = 1;
    }
    if (!find_path(start, goal)) {
        printf("NONE\n");
        return 0;
    }
    for (int i = 0; i < path_length; i++) {
        if (i) putchar(' ');
        printf("%d", path[i]);
    }
    putchar('\n');
    return 0;
}
`,
		wrongFails: "takes the short route, not the first one it finds",
	},

	"hcw-twos-complement-c": {
		correct: `
#include <stdio.h>
#include <stdint.h>
#include <string.h>

int main(void) {
    char operation[8];
    if (scanf("%7s", operation) != 1) return 1;
    if (strcmp(operation, "ENCODE") == 0) {
        int bits;
        long long value;
        if (scanf("%d %lld", &bits, &value) != 2) return 1;
        uint64_t encoded = (uint64_t)value & ((UINT64_C(1) << bits) - 1);
        for (int bit = bits - 1; bit >= 0; bit--) {
            putchar((encoded & (UINT64_C(1) << bit)) ? '1' : '0');
        }
        putchar('\n');
    } else if (strcmp(operation, "DECODE") == 0) {
        char pattern[33];
        if (scanf("%32s", pattern) != 1) return 1;
        int bits = (int)strlen(pattern);
        uint64_t encoded = 0;
        for (int i = 0; i < bits; i++) {
            encoded = encoded * 2 + (uint64_t)(pattern[i] - '0');
        }
        int64_t value = (int64_t)encoded;
        if (pattern[0] == '1') value -= (INT64_C(1) << bits);
        printf("%lld\n", (long long)value);
    }
    return 0;
}
`,
		// Uses magnitude bits for ENCODE and treats DECODE as unsigned: positive
		// examples work, but a negative value never receives two's complement.
		wrong: `
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(void) {
    char operation[8];
    if (scanf("%7s", operation) != 1) return 1;
    if (strcmp(operation, "ENCODE") == 0) {
        int bits;
        long long value;
        if (scanf("%d %lld", &bits, &value) != 2) return 1;
        unsigned long long magnitude = (unsigned long long)llabs(value);
        for (int bit = bits - 1; bit >= 0; bit--) {
            putchar((magnitude & (1ULL << bit)) ? '1' : '0');
        }
        putchar('\n');
    } else if (strcmp(operation, "DECODE") == 0) {
        char pattern[33];
        if (scanf("%32s", pattern) != 1) return 1;
        unsigned long long value = 0;
        for (size_t i = 0; pattern[i]; i++) value = value * 2 + (unsigned long long)(pattern[i] - '0');
        printf("%llu\n", value);
    }
    return 0;
}
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
			if a.lang != "c" {
				t.Parallel()
			}
			// Share the package's compile budget. Go resumes a parent's paused
			// parallel subtests alongside the *next* top-level test, so the
			// problem bank's three hundred sandboxed builds overlap these — and
			// unbounded, they starve a checkpoint past its limit, which then
			// reads as a correct solution failing.
			defer compileSlot()()
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
			if a.lang != "c" {
				t.Parallel()
			}
			defer compileSlot()()
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
			wantLanguage := runner.LangC
			if course == "linux" {
				// This checkpoint exercises a real command pipeline over a log file,
				// so shell is the curriculum rather than an implementation shortcut.
				wantLanguage = runner.LangShell
			}
			if a.lang != wantLanguage {
				t.Errorf("Foundations checkpoint %q uses %q; want %q", a.slug, a.lang, wantLanguage)
			}
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
// them. The last course in a path is exempt: its gate holds nothing back.
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
