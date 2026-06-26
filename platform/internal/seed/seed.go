// Package seed loads the initial course content (and a few sample assignments)
// into the database. The content is embedded into the binary, so production
// (Railway/Turso) can seed itself on first boot with no external files.
package seed

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

//go:embed data
var contentFS embed.FS

type appJSON struct {
	Config struct {
		Paths []struct {
			Name    string `json:"name"`
			Modules []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"modules"`
		} `json:"paths"`
	} `json:"config"`
}

var courseMeta = map[string]struct{ Emoji, Tagline string }{
	"frontend":  {"🎨", "Build what users see and touch"},
	"backend":   {"🗄️", "Power apps from behind the scenes"},
	"fullstack": {"🔗", "Connect front and back into apps"},
	"c":         {"⚙️", "Program close to the metal"},
	"java":      {"☕", "Robust, portable applications"},
	"python":    {"🐍", "Clear, friendly, everywhere"},
	"linux":     {"🐧", "Command the machine directly"},
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	return strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(s), "_"), "_")
}

// Counts returned by Run.
type Result struct{ Courses, Lessons, Assignments, Paths int }

// seedPaths are the initial career paths, each an ordered list of course slugs.
var seedPaths = []struct {
	Slug, Title, Tagline, Emoji, Level, Description string
	Courses                                         []string
}{
	{"frontend-developer", "Frontend Developer", "Build what users see and touch", "🎨", "Beginner",
		"Go from zero to building modern, responsive, interactive web interfaces.",
		[]string{"frontend"}},
	{"backend-developer", "Backend Developer", "Power apps from behind the scenes", "🗄️", "Intermediate",
		"Design APIs, model data, and ship reliable server-side systems — with the Linux skills to run them.",
		[]string{"backend", "python", "linux"}},
	{"fullstack-developer", "Full-Stack Developer", "Own the whole stack, front to back", "🔗", "Intermediate",
		"Combine frontend and backend skills to build complete applications end to end.",
		[]string{"frontend", "backend", "fullstack"}},
	{"cs-foundations", "Computer Science Foundations", "Program close to the metal", "⚙️", "Beginner",
		"Build durable fundamentals with C, Java, and the Linux command line.",
		[]string{"c", "java", "linux"}},
}

// RunIfEmpty seeds content only when there are no courses yet (a fresh DB).
func RunIfEmpty(ctx context.Context, st *store.Store) (Result, bool, error) {
	courses, err := st.ListCourses(ctx, true)
	if err != nil {
		return Result{}, false, err
	}
	if len(courses) > 0 {
		return Result{}, false, nil
	}
	r, err := Run(ctx, st)
	return r, true, err
}

// Run (idempotently) upserts all embedded courses, lessons, and sample
// assignments into the store.
func Run(ctx context.Context, st *store.Store) (Result, error) {
	var res Result
	raw, err := contentFS.ReadFile("data/app.json")
	if err != nil {
		return res, err
	}
	var app appJSON
	if err := json.Unmarshal(raw, &app); err != nil {
		return res, err
	}
	imgRewrite := strings.NewReplacer("/images/lessons/", "/static/img/lessons/")

	for ci, p := range app.Config.Paths {
		cslug := slugify(p.Name)
		m := courseMeta[cslug]
		courseID, err := st.UpsertCourse(ctx, store.Course{
			Slug: cslug, Title: p.Name, Emoji: m.Emoji, Tagline: m.Tagline, Sort: ci, Published: true,
		})
		if err != nil {
			return res, err
		}
		res.Courses++
		for li, mod := range p.Modules {
			lslug := slugify(mod.Title)
			body := ""
			if b, err := contentFS.ReadFile("data/notes/" + cslug + "/" + lslug + ".md"); err == nil {
				body = imgRewrite.Replace(stripFirstH1(string(b)))
			}
			if err := st.UpsertLesson(ctx, store.Lesson{
				CourseID: courseID, Slug: lslug, Title: mod.Title, Summary: mod.Description, Body: body, Sort: li,
			}); err != nil {
				return res, err
			}
			res.Lessons++
		}
	}

	for i, a := range sampleAssignments {
		var courseID sql.NullInt64
		if c, err := st.GetCourseBySlug(ctx, a.course); err == nil {
			courseID = sql.NullInt64{Int64: c.ID, Valid: true}
		}
		if err := st.UpsertAssignment(ctx, store.Assignment{
			CourseID: courseID, Slug: a.slug, Title: a.title, Language: a.lang,
			Prompt: a.prompt, Starter: a.starter, MaxPoints: 100, Published: true, Sort: i,
		}); err != nil {
			return res, err
		}
		res.Assignments++
	}

	for i, p := range seedPaths {
		id, err := st.UpsertPath(ctx, store.Path{
			Slug: p.Slug, Title: p.Title, Tagline: p.Tagline, Description: p.Description,
			Emoji: p.Emoji, Level: p.Level, Sort: i, Published: true,
		})
		if err != nil {
			return res, err
		}
		if err := st.SetPathCoursesBySlug(ctx, id, p.Courses); err != nil {
			return res, err
		}
		res.Paths++
	}
	return res, nil
}

type seedA struct{ slug, title, course, lang, prompt, starter string }

var sampleAssignments = []seedA{
	{"py-sum-list", "Sum a List", "python", "python",
		"Write a function `sum_list(nums)` that returns the sum of all numbers in the list `nums`. An empty list should return `0`.\n\nExplain your approach in the note to your mentor.",
		"def sum_list(nums):\n    # your code here\n    pass\n"},
	{"js-reverse", "Reverse a String", "frontend", "javascript",
		"Write a function `reverse(str)` that returns the characters of `str` in reverse order, **without** using the built-in `.reverse()`. Walk a loop yourself.",
		"function reverse(str) {\n  // your code here\n}\n"},
	{"c-max-three", "Largest of Three", "c", "c",
		"Read three integers and print the largest. Handle negative numbers correctly.",
		"#include <stdio.h>\n\nint main() {\n    int a, b, c;\n    scanf(\"%d %d %d\", &a, &b, &c);\n    // print the largest\n    return 0;\n}\n"},
}

// stripFirstH1 removes a leading "# Title" line (the page shows the title
// separately) and the blank lines around it.
func stripFirstH1(s string) string {
	lines := strings.Split(s, "\n")
	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	if i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "# ") {
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		return strings.Join(lines[i:], "\n")
	}
	return s
}
