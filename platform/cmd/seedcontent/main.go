// Command seedcontent migrates the legacy course content (data/app.json +
// data/notes/<world>/*.md) into the platform database. Idempotent: re-running
// upserts by slug. Run from the platform dir:
//
//	go run ./cmd/seedcontent              # seeds local gamifydev.db
//	DATABASE_URL=libsql://... go run ./cmd/seedcontent   # seeds Turso
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

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

var meta = map[string]struct{ Emoji, Tagline string }{
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

func main() {
	dataDir := flag.String("data", "../data", "path to the legacy data directory")
	flag.Parse()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "gamifydev.db"
	}
	st, err := store.Open(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer st.Close()
	ctx := context.Background()
	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(*dataDir, "app.json"))
	if err != nil {
		log.Fatalf("read app.json: %v", err)
	}
	var app appJSON
	if err := json.Unmarshal(raw, &app); err != nil {
		log.Fatalf("parse app.json: %v", err)
	}

	imgRewrite := strings.NewReplacer("/images/lessons/", "/static/img/lessons/")
	courses, lessons := 0, 0

	for ci, p := range app.Config.Paths {
		cslug := slugify(p.Name)
		m := meta[cslug]
		courseID, err := st.UpsertCourse(ctx, store.Course{
			Slug:      cslug,
			Title:     p.Name,
			Emoji:     m.Emoji,
			Tagline:   m.Tagline,
			Sort:      ci,
			Published: true,
		})
		if err != nil {
			log.Fatalf("course %s: %v", cslug, err)
		}
		courses++

		for li, mod := range p.Modules {
			lslug := slugify(mod.Title)
			body := ""
			notePath := filepath.Join(*dataDir, "notes", cslug, lslug+".md")
			if b, err := os.ReadFile(notePath); err == nil {
				body = imgRewrite.Replace(stripFirstH1(string(b)))
			}
			if err := st.UpsertLesson(ctx, store.Lesson{
				CourseID: courseID,
				Slug:     lslug,
				Title:    mod.Title,
				Summary:  mod.Description,
				Body:     body,
				Sort:     li,
			}); err != nil {
				log.Fatalf("lesson %s/%s: %v", cslug, lslug, err)
			}
			lessons++
		}
	}
	// Sample mentor-graded assignments, linked to courses by slug.
	type seedA struct {
		slug, title, course, lang, prompt, starter string
	}
	samples := []seedA{
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
	assignments := 0
	for i, a := range samples {
		var courseID sql.NullInt64
		if c, err := st.GetCourseBySlug(ctx, a.course); err == nil {
			courseID = sql.NullInt64{Int64: c.ID, Valid: true}
		}
		if err := st.UpsertAssignment(ctx, store.Assignment{
			CourseID: courseID, Slug: a.slug, Title: a.title, Language: a.lang,
			Prompt: a.prompt, Starter: a.starter, MaxPoints: 100, Published: true, Sort: i,
		}); err != nil {
			log.Fatalf("assignment %s: %v", a.slug, err)
		}
		assignments++
	}

	log.Printf("seeded %d courses, %d lessons, %d assignments", courses, lessons, assignments)
}

// stripFirstH1 removes a leading "# Title" line (the page shows the title
// separately) and any blank lines that follow it.
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
