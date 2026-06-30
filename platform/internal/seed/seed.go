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
				Section     string `json:"section"`
				Kind        string `json:"kind"`
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
			kind := mod.Kind
			if kind == "" {
				kind = "theory"
			}
			if err := st.UpsertLesson(ctx, store.Lesson{
				CourseID: courseID, Slug: lslug, Title: mod.Title, Summary: mod.Description, Body: body, Sort: li,
				Section: mod.Section, Kind: kind,
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
			Required: true, // act as the course checkpoint that gates the next course
		}); err != nil {
			return res, err
		}
		res.Assignments++
	}

	// Interactive steps for the demo workshop lesson.
	if c, err := st.GetCourseBySlug(ctx, "frontend"); err == nil {
		if l, err := st.GetLesson(ctx, c.ID, "workshop_build_a_cat_photo_app"); err == nil {
			if err := st.ReplaceLessonSteps(ctx, l.ID, catPhotoSteps); err != nil {
				return res, err
			}
		}
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

// catPhotoSteps is the interactive "Build a Cat Photo App" workshop. Each step's
// starter is the expected result of the previous one, so the page builds up.
// Check tests are tiny JS boolean expressions run against the preview document
// (`doc`); they use single quotes so they sit cleanly in an HTML data attribute.
var catPhotoSteps = []store.Step{
	{
		Instruction: "Every page starts with a main heading. Add an `<h1>` element with the text **CatPhotoApp**.",
		Starter:     "<!-- Add your h1 below -->\n",
		Checks:      `[{"text":"You should have an h1 element.","test":"doc.querySelector('h1')"},{"text":"Your h1 should say CatPhotoApp.","test":"doc.querySelector('h1') && /catphotoapp/i.test(doc.querySelector('h1').textContent)"}]`,
	},
	{
		Instruction: "The main content of a page belongs in a `<main>` element. Add a `<main>` element below your `h1`.",
		Starter:     "<h1>CatPhotoApp</h1>\n<!-- Add a main element below -->\n",
		Checks:      `[{"text":"Keep your h1 with the text CatPhotoApp.","test":"doc.querySelector('h1') && /catphotoapp/i.test(doc.querySelector('h1').textContent)"},{"text":"You should have a main element.","test":"doc.querySelector('main')"}]`,
	},
	{
		Instruction: "Inside `main`, add an `<h2>` with the text **Cat Photos** and a `<p>` that says something about cats.",
		Starter:     "<h1>CatPhotoApp</h1>\n<main>\n  <!-- Add an h2 and a p here -->\n</main>\n",
		Checks:      `[{"text":"main should contain an h2.","test":"doc.querySelector('main h2')"},{"text":"Your h2 should say Cat Photos.","test":"doc.querySelector('main h2') && /cat photos/i.test(doc.querySelector('main h2').textContent)"},{"text":"main should contain a non-empty paragraph (p).","test":"doc.querySelector('main p') && doc.querySelector('main p').textContent.trim().length > 0"}]`,
	},
	{
		Instruction: "Now add a photo. Add an `<img>` inside `main` with a `src` (any image URL) and a descriptive `alt` attribute.",
		Starter:     "<h1>CatPhotoApp</h1>\n<main>\n  <h2>Cat Photos</h2>\n  <p>Everybody loves cute cats online!</p>\n  <!-- Add an img with src and alt -->\n</main>\n",
		Checks:      `[{"text":"You should have an img element.","test":"doc.querySelector('img')"},{"text":"Your img needs a non-empty src.","test":"doc.querySelector('img') && (doc.querySelector('img').getAttribute('src')||'').length > 0"},{"text":"Your img needs a descriptive alt attribute.","test":"doc.querySelector('img') && (doc.querySelector('img').getAttribute('alt')||'').trim().length > 0"}]`,
	},
	{
		Instruction: "Add a link. Below the image add an `<a>` element whose `href` points anywhere (use `#` for now) and whose text mentions cats — e.g. *See more cat photos*.",
		Starter:     "<h1>CatPhotoApp</h1>\n<main>\n  <h2>Cat Photos</h2>\n  <p>Everybody loves cute cats online!</p>\n  <img src=\"https://cdn.freecodecamp.org/curriculum/cat-photo-app/relaxing-cat.jpg\" alt=\"A relaxing cat\">\n  <!-- Add an anchor link -->\n</main>\n",
		Checks:      `[{"text":"You should have an a (anchor) element.","test":"doc.querySelector('a')"},{"text":"Your link should have an href.","test":"doc.querySelector('a') && doc.querySelector('a').getAttribute('href') !== null"},{"text":"Your link text should mention cats.","test":"doc.querySelector('a') && /cat/i.test(doc.querySelector('a').textContent)"}]`,
	},
	{
		Instruction: "Finish with a list. Add an unordered list `<ul>` with at least **three** `<li>` items of things cats love.",
		Starter:     "<h1>CatPhotoApp</h1>\n<main>\n  <h2>Cat Photos</h2>\n  <p>Everybody loves cute cats online!</p>\n  <img src=\"https://cdn.freecodecamp.org/curriculum/cat-photo-app/relaxing-cat.jpg\" alt=\"A relaxing cat\">\n  <a href=\"#\">See more cat photos</a>\n  <!-- Add a ul with at least three li items -->\n</main>\n",
		Checks:      `[{"text":"You should have a ul element.","test":"doc.querySelector('ul')"},{"text":"Your list should have at least 3 li items.","test":"doc.querySelectorAll('ul li').length >= 3"}]`,
	},
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
