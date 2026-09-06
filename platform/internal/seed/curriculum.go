package seed

import (
	"context"

	"github.com/RA9/gamifydev/platform/internal/store"
)

var foundationCourses = map[string]bool{
	"c":                       true,
	"data_structures":         true,
	"complexity_and_analysis": true,
	"algorithms":              true,
	"how_computers_work":      true,
	"linux":                   true,
}

func courseLanguageMetadata(slug string) (primary, policy, exception string) {
	switch slug {
	case "c", "data_structures", "complexity_and_analysis", "algorithms", "how_computers_work":
		return "c", store.LanguagePolicyCOnly, ""
	case "linux":
		return "shell", store.LanguagePolicyCWithException, "Linux command-line labs and the final operations checkpoint require shell; all general programming work remains C."
	case "python", "backend":
		return "python", store.LanguagePolicyMixed, ""
	case "javascript_fundamentals", "javascript_and_the_dom", "interactive_javascript", "async_javascript_and_apis":
		return "javascript", store.LanguagePolicyMixed, ""
	case "java":
		return "java", store.LanguagePolicyMixed, ""
	default:
		return "", store.LanguagePolicyMixed, ""
	}
}

func lessonCurriculumMetadata(courseSlug, kind string, lessonIndex int) (workload int, mode, language string) {
	if foundationCourses[courseSlug] {
		// The 58 reading lessons contribute exactly 600 of the path's 2,000
		// workload minutes: 30%. The first lesson carries the orientation time.
		workload = 10
		if courseSlug == "c" && lessonIndex == 0 {
			workload = 30
		}
		language = "c"
		if courseSlug == "linux" {
			language = "shell"
		}
		return workload, store.ModeTheoretical, language
	}
	if kind == "lab" || kind == "workshop" || kind == "project" {
		return 60, store.ModePractical, store.LanguageNone
	}
	return 30, store.ModeTheoretical, store.LanguageNone
}

type foundationLab struct {
	course, slug, title, summary, language, mode string
	workload                                     int
	steps                                        []store.Step
}

var foundationLabs = []foundationLab{
	{
		course: "c", slug: "lab_decode_a_secret_signal", title: "Lab: Decode a Secret Signal",
		summary:  "Use loops and arithmetic to turn a tiny encoded transmission back into text.",
		language: "c", mode: store.ModeFun, workload: 100,
		steps: []store.Step{{
			Lang:        "c",
			Instruction: "A rover sent the character codes `72 69 76 76 79`. Complete the loop so the program prints `HELLO` followed by a newline.",
			Starter:     "#include <stdio.h>\n\nint main(void) {\n    int codes[] = {72, 69, 76, 76, 79};\n    for (int i = 0; i < 5; i++) {\n        /* print this character */\n    }\n    putchar('\\n');\n    return 0;\n}\n",
			Checks:      `[{"text":"The decoded signal should be HELLO.","test":"_out.strip() == 'HELLO'"},{"text":"The program should exit cleanly.","test":"_exit == 0"}]`,
		}},
	},
	{
		course: "data_structures", slug: "lab_build_a_stack_machine", title: "Lab: Build a Stack Machine",
		summary:  "Represent a stack with an array and trace last-in, first-out operations.",
		language: "c", mode: store.ModePractical, workload: 120,
		steps: []store.Step{{
			Lang:        "c",
			Instruction: "Push `4`, `8`, and `15` into the array-backed stack, then pop and print all three values. The output must demonstrate LIFO order: `15 8 4`.",
			Starter:     "#include <stdio.h>\n\nint main(void) {\n    int stack[8];\n    int top = 0;\n    /* push 4, 8, 15; then pop them into printf */\n    return 0;\n}\n",
			Checks:      `[{"text":"Values should leave the stack in LIFO order.","test":"_out.split() == ['15', '8', '4']"},{"text":"Use an explicit stack top.","test":"'top' in _code"}]`,
		}},
	},
	{
		course: "complexity_and_analysis", slug: "lab_measure_linear_growth", title: "Lab: Measure Linear Growth",
		summary:  "Count operations and connect measurements to an O(n) model.",
		language: "c", mode: store.ModePractical, workload: 70,
		steps: []store.Step{{
			Lang:        "c",
			Instruction: "Complete `count_work` so one unit of work is counted per array element. Print the counts for arrays of length 10 and 100 as `10 100`.",
			Starter:     "#include <stdio.h>\n\nint count_work(int n) {\n    int operations = 0;\n    /* one operation for each of n elements */\n    return operations;\n}\n\nint main(void) {\n    printf(\"%d %d\\n\", count_work(10), count_work(100));\n    return 0;\n}\n",
			Checks:      `[{"text":"Work should grow directly with input size.","test":"_out.split() == ['10', '100']"},{"text":"Use a loop to model repeated work.","test":"'for' in _code or 'while' in _code"}]`,
		}},
	},
	{
		course: "algorithms", slug: "lab_escape_the_maze", title: "Lab: Escape the Maze",
		summary:  "Use breadth-first search to race through a compact maze.",
		language: "c", mode: store.ModeFun, workload: 100,
		steps: []store.Step{{
			Lang:        "c",
			Instruction: "The graph below has edges `0→1`, `0→2`, `1→3`, and `2→4→3`. Complete the breadth-first search and print the shortest hop count from 0 to 3. It is `2`.",
			Starter:     "#include <stdio.h>\n\nint main(void) {\n    int graph[5][5] = {{0}};\n    graph[0][1] = graph[0][2] = graph[1][3] = graph[2][4] = graph[4][3] = 1;\n    int distance[5] = {-1, -1, -1, -1, -1};\n    int queue[5], head = 0, tail = 0;\n    /* breadth-first search from node 0 */\n    printf(\"%d\\n\", distance[3]);\n    return 0;\n}\n",
			Checks:      `[{"text":"The shortest route should take two hops.","test":"_out.strip() == '2'"},{"text":"Use a queue with separate head and tail positions.","test":"'head' in _code and 'tail' in _code"}]`,
		}},
	},
	{
		course: "how_computers_work", slug: "lab_inspect_binary_flags", title: "Lab: Inspect Binary Flags",
		summary:  "Use masks and shifts to inspect the bits inside an integer.",
		language: "c", mode: store.ModePractical, workload: 70,
		steps: []store.Step{{
			Lang:        "c",
			Instruction: "The byte `0b10110100` is stored as hexadecimal `0xB4`. Print its bits from most significant to least significant so the output is `10110100`.",
			Starter:     "#include <stdint.h>\n#include <stdio.h>\n\nint main(void) {\n    uint8_t value = 0xB4;\n    /* inspect masks 0x80 down to 0x01 */\n    putchar('\\n');\n    return 0;\n}\n",
			Checks:      `[{"text":"The byte should be printed as eight bits.","test":"_out.strip() == '10110100'"},{"text":"Use bitwise masking.","test":"'&' in _code"}]`,
		}},
	},
	{
		course: "linux", slug: "lab_trace_a_log_pipeline", title: "Lab: Trace a Log Pipeline",
		summary:  "Compose shell tools to count and rank repeated events.",
		language: "shell", mode: store.ModePractical, workload: 70,
		steps: []store.Step{{
			Lang:        "shell",
			Instruction: "Use a shell pipeline to count the repeated words below and print the most frequent word first as `3 build`, followed by `2 test` and `1 deploy`.",
			Starter:     "#!/usr/bin/env bash\nprintf '%s\\n' build test build deploy test build | # finish the pipeline\n",
			Checks:      `[{"text":"The pipeline should rank all three words.","test":"[line.split() for line in _out.strip().splitlines()] == [['3', 'build'], ['2', 'test'], ['1', 'deploy']]"},{"text":"Compose multiple shell tools with pipes.","test":"_code.count('|') >= 2"}]`,
		}},
	},
}

func seedFoundationLabs(ctx context.Context, st *store.Store) (int, error) {
	for i, lab := range foundationLabs {
		course, err := st.GetCourseBySlug(ctx, lab.course)
		if err != nil {
			return 0, err
		}
		if err := st.UpsertLesson(ctx, store.Lesson{
			CourseID: course.ID, Slug: lab.slug, Title: lab.title, Summary: lab.summary,
			Sort: 100 + i, Section: "Hands-on labs", Kind: "lab",
			WorkloadMinutes: lab.workload, Mode: lab.mode, Language: lab.language,
		}); err != nil {
			return 0, err
		}
		lesson, err := st.GetLesson(ctx, course.ID, lab.slug)
		if err != nil {
			return 0, err
		}
		if err := st.ReplaceLessonSteps(ctx, lesson.ID, lab.steps); err != nil {
			return 0, err
		}
	}
	return len(foundationLabs), nil
}
