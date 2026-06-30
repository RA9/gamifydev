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

	// Interactive step-based "lab" lessons in the frontend course.
	if c, err := st.GetCourseBySlug(ctx, "frontend"); err == nil {
		for slug, steps := range labSteps {
			if l, err := st.GetLesson(ctx, c.ID, slug); err == nil {
				if err := st.ReplaceLessonSteps(ctx, l.ID, steps); err != nil {
					return res, err
				}
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

// gameLaunchSteps is an original HTML-structure workshop: a launch page for a
// fictional indie game. Each step's starter is the result of the previous one so
// the page builds up. Checks are tiny JS boolean expressions over the preview
// document (`doc`), single-quoted so they sit cleanly in an HTML data attribute.
var gameLaunchSteps = []store.Step{
	{
		Instruction: "Every launch page needs a title. Add an `<h1>` with the name of your game — invent one!",
		Starter:     "<!-- Add your game's title in an h1 -->\n",
		Checks:      `[{"text":"You should have an h1 element.","test":"doc.querySelector('h1')"},{"text":"Your h1 should have a title in it.","test":"doc.querySelector('h1') && doc.querySelector('h1').textContent.trim().length > 0"}]`,
	},
	{
		Instruction: "Wrap the page content in a `<main>` element below the title.",
		Starter:     "<h1>Pixel Quest</h1>\n<!-- Add a main element below -->\n",
		Checks:      `[{"text":"Keep your h1 title.","test":"doc.querySelector('h1') && doc.querySelector('h1').textContent.trim().length > 0"},{"text":"You should have a main element.","test":"doc.querySelector('main')"}]`,
	},
	{
		Instruction: "Inside `main`, add an `<h2>` that says **About the Game** and a `<p>` describing what it's about.",
		Starter:     "<h1>Pixel Quest</h1>\n<main>\n  <!-- Add an h2 and a p here -->\n</main>\n",
		Checks:      `[{"text":"main should contain an h2.","test":"doc.querySelector('main h2')"},{"text":"Your h2 should mention 'About'.","test":"doc.querySelector('main h2') && /about/i.test(doc.querySelector('main h2').textContent)"},{"text":"main should contain a non-empty paragraph.","test":"doc.querySelector('main p') && doc.querySelector('main p').textContent.trim().length > 0"}]`,
	},
	{
		Instruction: "Show off the game with a cover image. Add an `<img>` inside `main` with a `src` and a descriptive `alt`. _Tip: use `https://placehold.co/600x300` as a placeholder._",
		Starter:     "<h1>Pixel Quest</h1>\n<main>\n  <h2>About the Game</h2>\n  <p>A retro platformer where every jump counts.</p>\n  <!-- Add a cover img with src and alt -->\n</main>\n",
		Checks:      `[{"text":"You should have an img element.","test":"doc.querySelector('img')"},{"text":"Your img needs a non-empty src.","test":"doc.querySelector('img') && (doc.querySelector('img').getAttribute('src')||'').length > 0"},{"text":"Your img needs a descriptive alt attribute.","test":"doc.querySelector('img') && (doc.querySelector('img').getAttribute('alt')||'').trim().length > 0"}]`,
	},
	{
		Instruction: "List what makes it fun. Add a `<ul>` with at least **three** `<li>` feature bullets.",
		Starter:     "<h1>Pixel Quest</h1>\n<main>\n  <h2>About the Game</h2>\n  <p>A retro platformer where every jump counts.</p>\n  <img src=\"https://placehold.co/600x300\" alt=\"Pixel Quest cover art\">\n  <!-- Add a ul with at least three feature items -->\n</main>\n",
		Checks:      `[{"text":"You should have a ul element.","test":"doc.querySelector('ul')"},{"text":"Your list should have at least 3 li items.","test":"doc.querySelectorAll('ul li').length >= 3"}]`,
	},
	{
		Instruction: "Finish with a call to action. Add an `<a>` link (use `#` for the `href`) whose text says **Play now**.",
		Starter:     "<h1>Pixel Quest</h1>\n<main>\n  <h2>About the Game</h2>\n  <p>A retro platformer where every jump counts.</p>\n  <img src=\"https://placehold.co/600x300\" alt=\"Pixel Quest cover art\">\n  <ul>\n    <li>30 hand-crafted levels</li>\n    <li>Original chiptune soundtrack</li>\n    <li>Local co-op</li>\n  </ul>\n  <!-- Add a Play now link -->\n</main>\n",
		Checks:      `[{"text":"You should have an a (anchor) element with an href.","test":"doc.querySelector('a') && doc.querySelector('a').getAttribute('href') !== null"},{"text":"Your link text should say 'Play'.","test":"doc.querySelector('a') && /play/i.test(doc.querySelector('a').textContent)"}]`,
	},
}

// pricingCardSteps is an original CSS workshop: style a product pricing card.
var pricingCardSteps = []store.Step{
	{
		Instruction: "Give the `.price-card` some breathing room — add `padding` of at least `24px`.",
		Starter:     "<style>\n  .price-card {\n    width: 240px;\n    /* add padding */\n  }\n  .price { font-size: 1rem; }\n  .buy { }\n</style>\n<div class=\"price-card\">\n  <h2>GamifyDev Pro</h2>\n  <p class=\"price\">$9/mo</p>\n  <ul>\n    <li>Unlimited courses</li>\n    <li>Mentor code reviews</li>\n    <li>Shareable certificates</li>\n  </ul>\n  <button class=\"buy\">Get Pro</button>\n</div>\n",
		Checks:      `[{"text":"The .price-card should have padding of at least 24px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.price-card')).paddingTop) >= 24"}]`,
	},
	{
		Instruction: "Lift the card off the page with a **`box-shadow`**.",
		Starter:     "<style>\n  .price-card {\n    width: 240px;\n    padding: 28px;\n    /* add a box-shadow */\n  }\n  .price { font-size: 1rem; }\n  .buy { }\n</style>\n<div class=\"price-card\">\n  <h2>GamifyDev Pro</h2>\n  <p class=\"price\">$9/mo</p>\n  <ul>\n    <li>Unlimited courses</li>\n    <li>Mentor code reviews</li>\n    <li>Shareable certificates</li>\n  </ul>\n  <button class=\"buy\">Get Pro</button>\n</div>\n",
		Checks:      `[{"text":"The .price-card should have a box-shadow.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.price-card')).boxShadow !== 'none'"}]`,
	},
	{
		Instruction: "Soften the card with rounded corners — add **`border-radius`**.",
		Starter:     "<style>\n  .price-card {\n    width: 240px;\n    padding: 28px;\n    box-shadow: 0 10px 25px rgba(0,0,0,0.12);\n    /* round the corners */\n  }\n  .price { font-size: 1rem; }\n  .buy { }\n</style>\n<div class=\"price-card\">\n  <h2>GamifyDev Pro</h2>\n  <p class=\"price\">$9/mo</p>\n  <ul>\n    <li>Unlimited courses</li>\n    <li>Mentor code reviews</li>\n    <li>Shareable certificates</li>\n  </ul>\n  <button class=\"buy\">Get Pro</button>\n</div>\n",
		Checks:      `[{"text":"The .price-card should have rounded corners.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.price-card')).borderTopLeftRadius) > 0"}]`,
	},
	{
		Instruction: "Make the price pop — give `.price` a larger **`font-size`** (at least `2rem`).",
		Starter:     "<style>\n  .price-card {\n    width: 240px;\n    padding: 28px;\n    box-shadow: 0 10px 25px rgba(0,0,0,0.12);\n    border-radius: 16px;\n  }\n  .price {\n    /* make this big */\n  }\n  .buy { }\n</style>\n<div class=\"price-card\">\n  <h2>GamifyDev Pro</h2>\n  <p class=\"price\">$9/mo</p>\n  <ul>\n    <li>Unlimited courses</li>\n    <li>Mentor code reviews</li>\n    <li>Shareable certificates</li>\n  </ul>\n  <button class=\"buy\">Get Pro</button>\n</div>\n",
		Checks:      `[{"text":"The .price should be at least 2rem (32px).","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.price')).fontSize) >= 30"}]`,
	},
	{
		Instruction: "Finish the **`.buy` button**: add `padding` (at least `10px`) and `cursor: pointer` so it feels clickable.",
		Starter:     "<style>\n  .price-card {\n    width: 240px;\n    padding: 28px;\n    box-shadow: 0 10px 25px rgba(0,0,0,0.12);\n    border-radius: 16px;\n  }\n  .price { font-size: 2rem; }\n  .buy {\n    /* add padding and a pointer cursor */\n  }\n</style>\n<div class=\"price-card\">\n  <h2>GamifyDev Pro</h2>\n  <p class=\"price\">$9/mo</p>\n  <ul>\n    <li>Unlimited courses</li>\n    <li>Mentor code reviews</li>\n    <li>Shareable certificates</li>\n  </ul>\n  <button class=\"buy\">Get Pro</button>\n</div>\n",
		Checks:      `[{"text":"The button should have padding of at least 10px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.buy')).paddingTop) >= 10"},{"text":"The button should use cursor: pointer.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.buy')).cursor === 'pointer'"}]`,
	},
}

// businessCardSteps is a CSS lab. The learner edits a <style> block; checks read
// computed styles via doc.defaultView.getComputedStyle. Each starter pre-fills
// the previous step's answer so the page builds up.
var businessCardSteps = []store.Step{
	{
		Instruction: "Give the `.card` a **background colour** — add `background` inside the `.card` rule.",
		Starter:     "<style>\n  .card {\n    /* add a background colour */\n  }\n</style>\n<div class=\"card\">\n  <h1>Ada Lovelace</h1>\n  <p>Pioneer of Programming</p>\n</div>\n",
		Checks:      `[{"text":"The .card should have a background colour.","test":"doc.querySelector('.card') && doc.defaultView.getComputedStyle(doc.querySelector('.card')).backgroundColor !== 'rgba(0, 0, 0, 0)'"}]`,
	},
	{
		Instruction: "Add **padding** of at least `16px` to the `.card` so the text isn't cramped.",
		Starter:     "<style>\n  .card {\n    background: #eef2ff;\n    /* add padding */\n  }\n</style>\n<div class=\"card\">\n  <h1>Ada Lovelace</h1>\n  <p>Pioneer of Programming</p>\n</div>\n",
		Checks:      `[{"text":"The .card should have padding of at least 16px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.card')).paddingTop) >= 16"}]`,
	},
	{
		Instruction: "Round the corners with **`border-radius`**.",
		Starter:     "<style>\n  .card {\n    background: #eef2ff;\n    padding: 24px;\n    /* round the corners */\n  }\n</style>\n<div class=\"card\">\n  <h1>Ada Lovelace</h1>\n  <p>Pioneer of Programming</p>\n</div>\n",
		Checks:      `[{"text":"The .card should have rounded corners (border-radius > 0).","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.card')).borderTopLeftRadius) > 0"}]`,
	},
	{
		Instruction: "**Centre** the text with `text-align: center`.",
		Starter:     "<style>\n  .card {\n    background: #eef2ff;\n    padding: 24px;\n    border-radius: 14px;\n    /* center the text */\n  }\n</style>\n<div class=\"card\">\n  <h1>Ada Lovelace</h1>\n  <p>Pioneer of Programming</p>\n</div>\n",
		Checks:      `[{"text":"The .card text should be centred.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.card')).textAlign === 'center'"}]`,
	},
	{
		Instruction: "Finally, give the **`h1` a colour** other than black. Add a `.card h1` rule.",
		Starter:     "<style>\n  .card {\n    background: #eef2ff;\n    padding: 24px;\n    border-radius: 14px;\n    text-align: center;\n  }\n  .card h1 {\n    /* set a colour */\n  }\n</style>\n<div class=\"card\">\n  <h1>Ada Lovelace</h1>\n  <p>Pioneer of Programming</p>\n</div>\n",
		Checks:      `[{"text":"The h1 should have a colour other than black.","test":"doc.querySelector('.card h1') && doc.defaultView.getComputedStyle(doc.querySelector('.card h1')).color !== 'rgb(0, 0, 0)'"}]`,
	},
}

// navBarSteps is a Flexbox lab.
var navBarSteps = []store.Step{
	{
		Instruction: "Make `.nav` a **flex container** with `display: flex`.",
		Starter:     "<style>\n  .nav {\n    /* make this a flex container */\n  }\n  .links { display: flex; }\n</style>\n<nav class=\"nav\">\n  <div class=\"brand\">GamifyDev</div>\n  <div class=\"links\">\n    <a href=\"#\">Home</a>\n    <a href=\"#\">Courses</a>\n    <a href=\"#\">Login</a>\n  </div>\n</nav>\n",
		Checks:      `[{"text":".nav should be a flex container.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.nav')).display === 'flex'"}]`,
	},
	{
		Instruction: "Push the brand and links to opposite ends with **`justify-content: space-between`**.",
		Starter:     "<style>\n  .nav {\n    display: flex;\n    /* push items apart */\n  }\n  .links { display: flex; }\n</style>\n<nav class=\"nav\">\n  <div class=\"brand\">GamifyDev</div>\n  <div class=\"links\">\n    <a href=\"#\">Home</a>\n    <a href=\"#\">Courses</a>\n    <a href=\"#\">Login</a>\n  </div>\n</nav>\n",
		Checks:      `[{"text":".nav should use justify-content: space-between.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.nav')).justifyContent === 'space-between'"}]`,
	},
	{
		Instruction: "Vertically centre the items with **`align-items: center`**.",
		Starter:     "<style>\n  .nav {\n    display: flex;\n    justify-content: space-between;\n    /* vertically center */\n  }\n  .links { display: flex; }\n</style>\n<nav class=\"nav\">\n  <div class=\"brand\">GamifyDev</div>\n  <div class=\"links\">\n    <a href=\"#\">Home</a>\n    <a href=\"#\">Courses</a>\n    <a href=\"#\">Login</a>\n  </div>\n</nav>\n",
		Checks:      `[{"text":".nav should use align-items: center.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.nav')).alignItems === 'center'"}]`,
	},
	{
		Instruction: "Space the links out: add a **`gap`** of at least `12px` to `.links`.",
		Starter:     "<style>\n  .nav {\n    display: flex;\n    justify-content: space-between;\n    align-items: center;\n  }\n  .links {\n    display: flex;\n    /* add a gap */\n  }\n</style>\n<nav class=\"nav\">\n  <div class=\"brand\">GamifyDev</div>\n  <div class=\"links\">\n    <a href=\"#\">Home</a>\n    <a href=\"#\">Courses</a>\n    <a href=\"#\">Login</a>\n  </div>\n</nav>\n",
		Checks:      `[{"text":".links should have a gap of at least 12px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.links')).gap) >= 12"}]`,
	},
}

// signupFormSteps is an HTML forms lab. Checks are DOM queries.
var signupFormSteps = []store.Step{
	{
		Instruction: "Add a **text input** for the name: an `<input>` with `type=\"text\"` inside the form.",
		Starter:     "<form>\n  <!-- Add a text input for the name -->\n</form>\n",
		Checks:      `[{"text":"You should have a text input.","test":"doc.querySelector('form input[type=text], form input:not([type])')"}]`,
	},
	{
		Instruction: "Add an **email field**: an `<input>` with `type=\"email\"`.",
		Starter:     "<form>\n  <input type=\"text\" placeholder=\"Name\">\n  <!-- Add an email input -->\n</form>\n",
		Checks:      `[{"text":"You should have an email input.","test":"doc.querySelector('form input[type=email]')"}]`,
	},
	{
		Instruction: "Add a **password field**: an `<input>` with `type=\"password\"`.",
		Starter:     "<form>\n  <input type=\"text\" placeholder=\"Name\">\n  <input type=\"email\" placeholder=\"Email\">\n  <!-- Add a password input -->\n</form>\n",
		Checks:      `[{"text":"You should have a password input.","test":"doc.querySelector('form input[type=password]')"}]`,
	},
	{
		Instruction: "Make the email field **required** — add the `required` attribute to your email input.",
		Starter:     "<form>\n  <input type=\"text\" placeholder=\"Name\">\n  <input type=\"email\" placeholder=\"Email\">\n  <input type=\"password\" placeholder=\"Password\">\n</form>\n",
		Checks:      `[{"text":"Your email input should be required.","test":"doc.querySelector('form input[type=email]') && doc.querySelector('form input[type=email]').hasAttribute('required')"}]`,
	},
	{
		Instruction: "Finish with a **submit button** — add a `<button>` (or `<input type=\"submit\">`).",
		Starter:     "<form>\n  <input type=\"text\" placeholder=\"Name\">\n  <input type=\"email\" placeholder=\"Email\" required>\n  <input type=\"password\" placeholder=\"Password\">\n  <!-- Add a submit button -->\n</form>\n",
		Checks:      `[{"text":"You should have a submit button.","test":"doc.querySelector('form button, form input[type=submit]')"}]`,
	},
}

// gridGallerySteps is a CSS Grid lab. Checks read computed styles.
var gridGallerySteps = []store.Step{
	{
		Instruction: "Make `.gallery` a **grid container** with `display: grid`.",
		Starter:     "<style>\n  .gallery {\n    /* make this a grid */\n  }\n  .gallery div { background: #c7d2fe; height: 70px; border-radius: 8px; }\n</style>\n<div class=\"gallery\">\n  <div></div><div></div><div></div>\n  <div></div><div></div><div></div>\n</div>\n",
		Checks:      `[{"text":".gallery should be a grid container.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.gallery')).display === 'grid'"}]`,
	},
	{
		Instruction: "Lay out **3 columns** with `grid-template-columns: repeat(3, 1fr);`.",
		Starter:     "<style>\n  .gallery {\n    display: grid;\n    /* add 3 columns */\n  }\n  .gallery div { background: #c7d2fe; height: 70px; border-radius: 8px; }\n</style>\n<div class=\"gallery\">\n  <div></div><div></div><div></div>\n  <div></div><div></div><div></div>\n</div>\n",
		Checks:      `[{"text":"Your grid should have 3 columns.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.gallery')).gridTemplateColumns.split(' ').length === 3"}]`,
	},
	{
		Instruction: "Add a **gap** of at least `10px` between the items.",
		Starter:     "<style>\n  .gallery {\n    display: grid;\n    grid-template-columns: repeat(3, 1fr);\n    /* add a gap */\n  }\n  .gallery div { background: #c7d2fe; height: 70px; border-radius: 8px; }\n</style>\n<div class=\"gallery\">\n  <div></div><div></div><div></div>\n  <div></div><div></div><div></div>\n</div>\n",
		Checks:      `[{"text":"Your grid should have a gap of at least 10px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.gallery')).gap) >= 10"}]`,
	},
	{
		Instruction: "Make the **first item wider** — give `.gallery div:first-child` `grid-column: span 2;`.",
		Starter:     "<style>\n  .gallery {\n    display: grid;\n    grid-template-columns: repeat(3, 1fr);\n    gap: 12px;\n  }\n  .gallery div { background: #c7d2fe; height: 70px; border-radius: 8px; }\n  .gallery div:first-child {\n    /* span 2 columns */\n  }\n</style>\n<div class=\"gallery\">\n  <div></div><div></div><div></div>\n  <div></div><div></div><div></div>\n</div>\n",
		Checks:      `[{"text":"The first item should span 2 columns.","test":"/span 2/.test(doc.defaultView.getComputedStyle(doc.querySelector('.gallery div')).gridColumn)"}]`,
	},
}

// labSteps maps a frontend lesson slug to its interactive steps.
var labSteps = map[string][]store.Step{
	"workshop_build_a_game_launch_page":         gameLaunchSteps,
	"workshop_build_a_sign_up_form":             signupFormSteps,
	"workshop_style_a_business_card":            businessCardSteps,
	"workshop_build_a_pricing_card":             pricingCardSteps,
	"workshop_build_a_nav_bar_with_flexbox":     navBarSteps,
	"workshop_build_an_image_gallery_with_grid": gridGallerySteps,
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
