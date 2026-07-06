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

	// Interactive step-based "lab" lessons, keyed by course slug.
	for cslug, labs := range courseLabs {
		c, err := st.GetCourseBySlug(ctx, cslug)
		if err != nil {
			continue
		}
		for slug, steps := range labs {
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

// achievementBadgeSteps is an original CSS lab: a circular achievement badge.
var achievementBadgeSteps = []store.Step{
	{
		Instruction: "Give the `.badge` a fixed size — a `width` and `height` of about `96px`.",
		Starter:     "<style>\n  .badge {\n    /* give it a size */\n  }\n</style>\n<div class=\"badge\">🏆</div>\n",
		Checks:      `[{"text":"The badge should be at least 80x80px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.badge')).width) >= 80 && parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.badge')).height) >= 80"}]`,
	},
	{
		Instruction: "Make it a perfect **circle** with `border-radius: 50%`.",
		Starter:     "<style>\n  .badge {\n    width: 96px;\n    height: 96px;\n    /* make it round */\n  }\n</style>\n<div class=\"badge\">🏆</div>\n",
		Checks:      `[{"text":"The badge should be a circle (border-radius: 50%).","test":"doc.defaultView.getComputedStyle(doc.querySelector('.badge')).borderTopLeftRadius === '50%' || parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.badge')).borderTopLeftRadius) >= 20"}]`,
	},
	{
		Instruction: "Centre the emoji: make `.badge` a **flex** container with `align-items: center` and `justify-content: center`.",
		Starter:     "<style>\n  .badge {\n    width: 96px;\n    height: 96px;\n    border-radius: 50%;\n    /* center the icon */\n  }\n</style>\n<div class=\"badge\">🏆</div>\n",
		Checks:      `[{"text":"The badge should center its content with flexbox.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.badge')).display === 'flex' && doc.defaultView.getComputedStyle(doc.querySelector('.badge')).alignItems === 'center' && doc.defaultView.getComputedStyle(doc.querySelector('.badge')).justifyContent === 'center'"}]`,
	},
	{
		Instruction: "Give it a **background** — a solid colour or a `linear-gradient`.",
		Starter:     "<style>\n  .badge {\n    width: 96px;\n    height: 96px;\n    border-radius: 50%;\n    display: flex;\n    align-items: center;\n    justify-content: center;\n    /* add a background */\n  }\n</style>\n<div class=\"badge\">🏆</div>\n",
		Checks:      `[{"text":"The badge should have a background colour or gradient.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.badge')).backgroundImage !== 'none' || doc.defaultView.getComputedStyle(doc.querySelector('.badge')).backgroundColor !== 'rgba(0, 0, 0, 0)'"}]`,
	},
	{
		Instruction: "Make the emoji **bigger** — set `font-size` to at least `40px`.",
		Starter:     "<style>\n  .badge {\n    width: 96px;\n    height: 96px;\n    border-radius: 50%;\n    display: flex;\n    align-items: center;\n    justify-content: center;\n    background: linear-gradient(135deg, #6d4aff, #a855f7);\n    /* enlarge the emoji */\n  }\n</style>\n<div class=\"badge\">🏆</div>\n",
		Checks:      `[{"text":"The emoji should be at least 40px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.badge')).fontSize) >= 40"}]`,
	},
}

// xpBarSteps is an original CSS lab: an XP progress bar (a fill inside a track).
var xpBarSteps = []store.Step{
	{
		Instruction: "Style the track: give `.bar` a **height** of at least `14px` and a light **background**.",
		Starter:     "<style>\n  .bar {\n    width: 300px;\n    /* add a height and a background */\n  }\n  .fill {\n    /* you'll style this next */\n  }\n</style>\n<div class=\"bar\"><div class=\"fill\"></div></div>\n",
		Checks:      `[{"text":"The .bar should have a height of at least 14px and a background.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.bar')).height) >= 14 && doc.defaultView.getComputedStyle(doc.querySelector('.bar')).backgroundColor !== 'rgba(0, 0, 0, 0)'"}]`,
	},
	{
		Instruction: "Round the track and clip the fill: add `border-radius` and `overflow: hidden` to `.bar`.",
		Starter:     "<style>\n  .bar {\n    width: 300px;\n    height: 18px;\n    background: #e2e8f0;\n    /* round it and hide overflow */\n  }\n  .fill {\n    /* you'll style this next */\n  }\n</style>\n<div class=\"bar\"><div class=\"fill\"></div></div>\n",
		Checks:      `[{"text":"The .bar should have rounded corners and overflow: hidden.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.bar')).borderTopLeftRadius) > 0 && doc.defaultView.getComputedStyle(doc.querySelector('.bar')).overflow === 'hidden'"}]`,
	},
	{
		Instruction: "Style the fill: give `.fill` a `height` of `100%` and a bright **background** colour.",
		Starter:     "<style>\n  .bar {\n    width: 300px;\n    height: 18px;\n    background: #e2e8f0;\n    border-radius: 999px;\n    overflow: hidden;\n  }\n  .fill {\n    /* height 100% and a colour */\n  }\n</style>\n<div class=\"bar\"><div class=\"fill\"></div></div>\n",
		Checks:      `[{"text":"The .fill should have height and a background colour.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.fill')).height) >= 14 && doc.defaultView.getComputedStyle(doc.querySelector('.fill')).backgroundColor !== 'rgba(0, 0, 0, 0)'"}]`,
	},
	{
		Instruction: "Set the progress: give `.fill` a **`width`** of `70%` so it only fills part of the track.",
		Starter:     "<style>\n  .bar {\n    width: 300px;\n    height: 18px;\n    background: #e2e8f0;\n    border-radius: 999px;\n    overflow: hidden;\n  }\n  .fill {\n    height: 100%;\n    background: #6d4aff;\n    /* set a partial width */\n  }\n</style>\n<div class=\"bar\"><div class=\"fill\"></div></div>\n",
		Checks:      `[{"text":"The .fill should be a partial width (wider than 0, narrower than the track).","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.fill')).width) > 0 && parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.fill')).width) < 300"}]`,
	},
}

// leaderboardRowSteps is an original Flexbox lab: a leaderboard entry.
var leaderboardRowSteps = []store.Step{
	{
		Instruction: "Make `.row` a **flex** container with `display: flex`.",
		Starter:     "<style>\n  .row {\n    /* make a flex row */\n  }\n  .rank { font-weight: 800; }\n</style>\n<div class=\"row\">\n  <span class=\"rank\">#1</span>\n  <span class=\"name\">Ada</span>\n  <span class=\"score\">1200 XP</span>\n</div>\n",
		Checks:      `[{"text":".row should be a flex container.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.row')).display === 'flex'"}]`,
	},
	{
		Instruction: "Vertically centre the items with `align-items: center`.",
		Starter:     "<style>\n  .row {\n    display: flex;\n    /* vertically center */\n  }\n  .rank { font-weight: 800; }\n</style>\n<div class=\"row\">\n  <span class=\"rank\">#1</span>\n  <span class=\"name\">Ada</span>\n  <span class=\"score\">1200 XP</span>\n</div>\n",
		Checks:      `[{"text":".row should use align-items: center.","test":"doc.defaultView.getComputedStyle(doc.querySelector('.row')).alignItems === 'center'"}]`,
	},
	{
		Instruction: "Add a **`gap`** of at least `8px` between the items.",
		Starter:     "<style>\n  .row {\n    display: flex;\n    align-items: center;\n    /* add a gap */\n  }\n  .rank { font-weight: 800; }\n</style>\n<div class=\"row\">\n  <span class=\"rank\">#1</span>\n  <span class=\"name\">Ada</span>\n  <span class=\"score\">1200 XP</span>\n</div>\n",
		Checks:      `[{"text":".row should have a gap of at least 8px.","test":"parseInt(doc.defaultView.getComputedStyle(doc.querySelector('.row')).gap) >= 8"}]`,
	},
	{
		Instruction: "Push the score to the far right: give `.name` **`flex: 1`** so it grows and fills the space.",
		Starter:     "<style>\n  .row {\n    display: flex;\n    align-items: center;\n    gap: 12px;\n  }\n  .rank { font-weight: 800; }\n  .name {\n    /* let this grow */\n  }\n</style>\n<div class=\"row\">\n  <span class=\"rank\">#1</span>\n  <span class=\"name\">Ada</span>\n  <span class=\"score\">1200 XP</span>\n</div>\n",
		Checks:      `[{"text":".name should grow to fill the space (flex: 1).","test":"doc.defaultView.getComputedStyle(doc.querySelector('.name')).flexGrow === '1'"}]`,
	},
}

// levelUpSteps is a JavaScript lab (Lang "js"). The learner's code runs in a
// sandbox; checks call their functions and inspect `logs` (console output).
var levelUpSteps = []store.Step{
	{
		Lang:        "js",
		Instruction: "Create a variable `xp`, set it to `150`, and print it with `console.log(xp)`.",
		Starter:     "// Create xp and log it\n",
		Checks:      `[{"text":"xp should be the number 150.","test":"typeof xp === 'number' && xp === 150"},{"text":"You should log something to the console.","test":"logs.length > 0"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Write a function `levelUp(xp)` that **returns** `xp + 100`.",
		Starter:     "let xp = 150;\nconsole.log(xp);\n\n// Write a function levelUp(xp) that returns xp + 100\n",
		Checks:      `[{"text":"levelUp should be a function.","test":"typeof levelUp === 'function'"},{"text":"levelUp(200) should return 300.","test":"levelUp(200) === 300"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Call `levelUp(150)` and log the result — the console should show `250`.",
		Starter:     "let xp = 150;\nconsole.log(xp);\n\nfunction levelUp(xp) {\n  return xp + 100;\n}\n\n// Call levelUp(150) and log the result\n",
		Checks:      `[{"text":"The console should print 250.","test":"logs.some(function(l){return l.trim() === '250'})"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Write a function `isMaxLevel(level)` that **returns** `true` when `level` is 100 or more, otherwise `false`.",
		Starter:     "let xp = 150;\nconsole.log(xp);\n\nfunction levelUp(xp) {\n  return xp + 100;\n}\nconsole.log(levelUp(150));\n\n// Write isMaxLevel(level)\n",
		Checks:      `[{"text":"isMaxLevel should be a function.","test":"typeof isMaxLevel === 'function'"},{"text":"isMaxLevel(100) should be true.","test":"isMaxLevel(100) === true"},{"text":"isMaxLevel(50) should be false.","test":"isMaxLevel(50) === false"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Total up an array. Create `scores = [10, 25, 5]` and log the sum — it should print `40`. (Try `.reduce()` or a loop.)",
		Starter:     "let xp = 150;\nconsole.log(xp);\n\nfunction levelUp(xp) {\n  return xp + 100;\n}\nconsole.log(levelUp(150));\n\nfunction isMaxLevel(level) {\n  return level >= 100;\n}\n\n// Create scores and log the total\n",
		Checks:      `[{"text":"scores should be an array of 3 numbers.","test":"Array.isArray(scores) && scores.length === 3"},{"text":"The console should print 40.","test":"logs.some(function(l){return l.trim() === '40'})"}]`,
	},
}

// inventorySteps is a JavaScript lab on arrays & objects (pure logic + console).
var inventorySteps = []store.Step{
	{
		Lang:        "js",
		Instruction: "Create an array `items` of **3 objects**, each with a `name` and a `price`. Then log `items.length`.",
		Starter:     "// Create an array `items` of 3 objects (name, price)\n// then log items.length\n",
		Checks:      `[{"text":"items should be an array of 3 objects with name and price.","test":"Array.isArray(items) && items.length===3 && typeof items[0]==='object' && 'name' in items[0] && 'price' in items[0]"},{"text":"Log the number of items (3).","test":"logs.some(function(l){return l.trim()==='3'})"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Add a 4th item with `items.push(...)`, then log the new `items.length` (4).",
		Starter:     "const items = [\n  { name: 'Sword', price: 30 },\n  { name: 'Shield', price: 20 },\n  { name: 'Potion', price: 5 },\n];\nconsole.log(items.length);\n\n// Push a 4th item, then log items.length\n",
		Checks:      `[{"text":"items should now have 4 entries.","test":"Array.isArray(items) && items.length===4"},{"text":"Log the new length (4).","test":"logs.some(function(l){return l.trim()==='4'})"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Make a list of just the names with `items.map(i => i.name)`; store it in `names` and log it.",
		Starter:     "const items = [\n  { name: 'Sword', price: 30 },\n  { name: 'Shield', price: 20 },\n  { name: 'Potion', price: 5 },\n];\nconsole.log(items.length);\nitems.push({ name: 'Bow', price: 25 });\nconsole.log(items.length);\n\n// Create names with map, then log it\n",
		Checks:      `[{"text":"names should be an array of every item's name.","test":"Array.isArray(names) && names.length===items.length && names.indexOf(items[0].name)!==-1"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Add up the prices with `reduce`: store the sum in `total` and log it.",
		Starter:     "const items = [\n  { name: 'Sword', price: 30 },\n  { name: 'Shield', price: 20 },\n  { name: 'Potion', price: 5 },\n];\nitems.push({ name: 'Bow', price: 25 });\nconst names = items.map(i => i.name);\nconsole.log(names);\n\n// Create total with reduce, then log it\n",
		Checks:      `[{"text":"total should equal the sum of all prices.","test":"typeof total==='number' && total===items.reduce(function(s,i){return s+i.price},0)"}]`,
	},
	{
		Lang:        "js",
		Instruction: "Find the **cheapest** item (lowest price). Store it in `cheapest` and log its name.",
		Starter:     "const items = [\n  { name: 'Sword', price: 30 },\n  { name: 'Shield', price: 20 },\n  { name: 'Potion', price: 5 },\n];\nitems.push({ name: 'Bow', price: 25 });\nconst total = items.reduce((s, i) => s + i.price, 0);\nconsole.log(total);\n\n// Find the cheapest item; store in `cheapest` and log its name\n",
		Checks:      `[{"text":"cheapest should be the lowest-priced item.","test":"cheapest && typeof cheapest==='object' && cheapest.price===Math.min.apply(null, items.map(function(i){return i.price}))"}]`,
	},
}

// likeScaffold is the markup the DOM lab's JavaScript runs against.
const likeScaffold = "<button id=\"like\">👍 Like</button>\n<p>Likes: <span id=\"count\">0</span></p>\n<button id=\"reset\">Reset</button>\n"

// likeButtonSteps is a JavaScript DOM + events lab; checks simulate clicks.
var likeButtonSteps = []store.Step{
	{
		Lang:        "js",
		Scaffold:    likeScaffold,
		Instruction: "Wire up the Like button: add a **click** listener on `#like` that adds `1` to the number in `#count` each time it's clicked.",
		Starter:     "// Select #like and #count, then add a click listener that increments the count\n",
		Checks:      `[{"text":"Clicking Like should add 1 to the count each time.","test":"(function(){var b=document.querySelector('#like'),c=document.querySelector('#count');if(!b||!c)return false;var s=parseInt(c.textContent)||0;b.click();b.click();b.click();return (parseInt(c.textContent)||0)===s+3;})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    likeScaffold,
		Instruction: "Also reflect the count in the **button text** — after a click it should read like `👍 3`.",
		Starter:     "const like = document.querySelector('#like');\nconst count = document.querySelector('#count');\nlike.addEventListener('click', () => {\n  count.textContent = Number(count.textContent) + 1;\n});\n\n// Also update the button's text to include the count\n",
		Checks:      `[{"text":"After clicking, the Like button text should include the count.","test":"(function(){var b=document.querySelector('#like'),c=document.querySelector('#count');b.click();b.click();return b.textContent.indexOf(c.textContent.trim())!==-1;})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    likeScaffold,
		Instruction: "Wire the **Reset** button (`#reset`) so it sets the count back to `0`.",
		Starter:     "const like = document.querySelector('#like');\nconst count = document.querySelector('#count');\nconst reset = document.querySelector('#reset');\nlike.addEventListener('click', () => {\n  count.textContent = Number(count.textContent) + 1;\n  like.textContent = '👍 ' + count.textContent;\n});\n\n// Wire the Reset button to set the count back to 0\n",
		Checks:      `[{"text":"Reset should set the count back to 0.","test":"(function(){var l=document.querySelector('#like'),c=document.querySelector('#count'),r=document.querySelector('#reset');l.click();l.click();r.click();return c.textContent.trim()==='0';})()"}]`,
	},
}

// pythonWarmupSteps is the first interactive Python lab — variables, print,
// functions, and lists — run through Pyodide (Python compiled to WASM).
var pythonWarmupSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Let's warm up. Create a variable `name` set to your name (a string), then `print` a greeting that uses it — something like `Hello, Ada!`.",
		Starter:     "# Set `name` to your name, then print a greeting that uses it\n",
		Checks:      `[{"text":"name should be a non-empty string.","test":"isinstance(name, str) and len(name) > 0"},{"text":"You should print a greeting that includes your name.","test":"name in _out"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Define a function `double(n)` that **returns** `n * 2`. (Return it — don't print it.)",
		Starter:     "name = \"Ada\"\nprint(\"Hello, \" + name + \"!\")\n\n# Define a function double(n) that returns n * 2\n",
		Checks:      `[{"text":"double should be a function.","test":"callable(double)"},{"text":"double(5) should return 10.","test":"double(5) == 10"},{"text":"double(0) should return 0.","test":"double(0) == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Now **call** `double(21)` and `print` the result — the output should show `42`.",
		Starter:     "def double(n):\n    return n * 2\n\n# Call double(21) and print the result\n",
		Checks:      `[{"text":"The printed output should include 42.","test":"\"42\" in _out"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Make a list called `scores` holding the numbers `4, 8, 15, 16, 23, 42`, then `print` how many items it has using `len(scores)`.",
		Starter:     "# Create the list `scores`, then print its length\n",
		Checks:      `[{"text":"scores should be a list of 6 numbers.","test":"isinstance(scores, list) and len(scores) == 6"},{"text":"You should print the length, 6.","test":"\"6\" in _out"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Finish by printing the **total** of every score using Python's built-in `sum()`. The total of the list should be `108`.",
		Starter:     "scores = [4, 8, 15, 16, 23, 42]\n\n# Print the total of all the scores using sum()\n",
		Checks:      `[{"text":"You should print the total using sum().","test":"\"108\" in _out"},{"text":"The list should still hold all 6 scores.","test":"sum(scores) == 108"}]`,
	},
}

// courseLabs maps a course slug to that course's interactive lab lessons
// (lesson slug → steps).
var courseLabs = map[string]map[string][]store.Step{
	"frontend": labSteps,
	"python":   pythonLabs,
}

// xpCalcSteps — numbers, functions, and f-strings via an XP calculator.
var xpCalcSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Different foes give different XP. Write `xp_reward(kind)` that returns `10` for `\"easy\"`, `25` for `\"normal\"`, `100` for `\"boss\"`, and `0` for anything else.",
		Starter:     "def xp_reward(kind):\n    # return the XP for this foe kind\n    pass\n",
		Checks:      `[{"text":"xp_reward should be a function.","test":"callable(xp_reward)"},{"text":"xp_reward('easy') should be 10.","test":"xp_reward('easy') == 10"},{"text":"xp_reward('boss') should be 100.","test":"xp_reward('boss') == 100"},{"text":"An unknown kind should give 0.","test":"xp_reward('dragon') == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `total_xp(rewards)` that returns the sum of a list of XP numbers. An empty list should return `0`.",
		Starter:     "def xp_reward(kind):\n    return {'easy': 10, 'normal': 25, 'boss': 100}.get(kind, 0)\n\ndef total_xp(rewards):\n    # add up every number in the list\n    pass\n",
		Checks:      `[{"text":"total_xp([10, 25, 100]) should be 135.","test":"total_xp([10, 25, 100]) == 135"},{"text":"An empty list should total 0.","test":"total_xp([]) == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Every 100 XP is one level. Write `level_for(xp)` that returns how many **full** levels that XP is worth — use integer division (`//`).",
		Starter:     "def level_for(xp):\n    # 100 XP per level\n    pass\n",
		Checks:      `[{"text":"0 XP is level 0.","test":"level_for(0) == 0"},{"text":"100 XP is level 1.","test":"level_for(100) == 1"},{"text":"250 XP is level 2.","test":"level_for(250) == 2"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Finish with a summary. For `xp = 350`, use an **f-string** to print `Level 3 · 350 XP` (call your `level_for` to get the level).",
		Starter:     "def level_for(xp):\n    return xp // 100\n\nxp = 350\n# Print a summary line like: Level 3 · 350 XP\n",
		Checks:      `[{"text":"Your summary should include the level, 3.","test":"'Level 3' in _out"},{"text":"Your summary should include the XP, 350.","test":"'350' in _out"}]`,
	},
}

// lootInventorySteps — lists, dictionaries, and loops via a loot bag.
var lootInventorySteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Start an empty `inventory` list, then add `'sword'` and `'shield'` to it using `.append()`.",
		Starter:     "# Make an empty list called inventory, then append 'sword' and 'shield'\n",
		Checks:      `[{"text":"inventory should be a list.","test":"isinstance(inventory, list)"},{"text":"It should contain 'sword' and 'shield'.","test":"'sword' in inventory and 'shield' in inventory"},{"text":"It should hold exactly 2 items.","test":"len(inventory) == 2"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Make a `prices` dictionary that maps `'sword'` → `100`, `'shield'` → `50`, and `'potion'` → `10`.",
		Starter:     "# Create the prices dictionary\n",
		Checks:      `[{"text":"prices should be a dict.","test":"isinstance(prices, dict)"},{"text":"A sword should cost 100.","test":"prices['sword'] == 100"},{"text":"There should be 3 prices.","test":"len(prices) == 3"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `total_value(items, prices)` that adds up the price of every item in the `items` list (items can repeat). Return `0` for an empty list.",
		Starter:     "def total_value(items, prices):\n    # sum the price of each item\n    pass\n",
		Checks:      `[{"text":"Two potions and a sword cost 120.","test":"total_value(['potion', 'potion', 'sword'], {'sword': 100, 'potion': 10}) == 120"},{"text":"An empty bag is worth 0.","test":"total_value([], {}) == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Loop over `inventory` and print each item next to its price, like `sword — 100` (use the `prices` dict).",
		Starter:     "inventory = ['sword', 'shield']\nprices = {'sword': 100, 'shield': 50, 'potion': 10}\n# Print each item with its price, like: sword — 100\n",
		Checks:      `[{"text":"You should print the sword's price.","test":"'sword — 100' in _out"},{"text":"You should print the shield's price.","test":"'shield — 50' in _out"}]`,
	},
}

// comboSteps — loops and conditionals via an original FizzBuzz-style combo meter.
var comboSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Our combo meter labels each hit number: divisible by 3 → `COMBO`, divisible by 5 → `CRIT`, divisible by **both** → `COMBOCRIT`, otherwise the number itself as text. Write `label(n)` that returns the right **string**.",
		Starter:     "def label(n):\n    pass\n",
		Checks:      `[{"text":"3 is a COMBO.","test":"label(3) == 'COMBO'"},{"text":"5 is a CRIT.","test":"label(5) == 'CRIT'"},{"text":"15 is a COMBOCRIT.","test":"label(15) == 'COMBOCRIT'"},{"text":"7 is just '7'.","test":"label(7) == '7'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `combo(n)` that returns a **list** of the labels for every hit from `1` to `n` inclusive.",
		Starter:     "def label(n):\n    if n % 15 == 0:\n        return 'COMBOCRIT'\n    if n % 3 == 0:\n        return 'COMBO'\n    if n % 5 == 0:\n        return 'CRIT'\n    return str(n)\n\ndef combo(n):\n    pass\n",
		Checks:      `[{"text":"combo(5) has five labels.","test":"len(combo(5)) == 5"},{"text":"combo(5) is ['1', '2', 'COMBO', '4', 'CRIT'].","test":"combo(5) == ['1', '2', 'COMBO', '4', 'CRIT']"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Now print the full combo for `n = 15`, one label per line.",
		Starter:     "def label(n):\n    if n % 15 == 0: return 'COMBOCRIT'\n    if n % 3 == 0: return 'COMBO'\n    if n % 5 == 0: return 'CRIT'\n    return str(n)\n\n# Print label(1) through label(15), one per line\n",
		Checks:      `[{"text":"Hit 15 should print COMBOCRIT.","test":"'COMBOCRIT' in _out"},{"text":"You should print 15 lines.","test":"len([l for l in _out.strip().split(chr(10)) if l]) == 15"}]`,
	},
}

// playerClassSteps — classes and objects via a Player model.
var playerClassSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Define a `Player` class. Its `__init__` should take a `name`, store it on `self.name`, and set `self.hp` to `100`.",
		Starter:     "class Player:\n    def __init__(self, name):\n        pass\n",
		Checks:      `[{"text":"Player should be a class you can call.","test":"callable(Player)"},{"text":"A new player keeps their name.","test":"Player('Zed').name == 'Zed'"},{"text":"A new player starts at 100 HP.","test":"Player('Zed').hp == 100"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Add a `take_damage(self, amount)` method that subtracts `amount` from `self.hp` — but never lets HP drop below `0`.",
		Starter:     "class Player:\n    def __init__(self, name):\n        self.name = name\n        self.hp = 100\n\n    def take_damage(self, amount):\n        pass\n",
		Checks:      `[{"text":"30 damage leaves 70 HP.","test":"[p := Player('Zed'), p.take_damage(30), p.hp][2] == 70"},{"text":"HP never goes negative.","test":"[p := Player('Zed'), p.take_damage(999), p.hp][2] == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Add an `is_alive(self)` method that returns `True` while HP is above `0`, and `False` otherwise.",
		Starter:     "class Player:\n    def __init__(self, name):\n        self.name = name\n        self.hp = 100\n\n    def take_damage(self, amount):\n        self.hp = max(0, self.hp - amount)\n\n    def is_alive(self):\n        pass\n",
		Checks:      `[{"text":"A fresh player is alive.","test":"Player('Zed').is_alive() == True"},{"text":"A defeated player is not.","test":"[p := Player('Zed'), p.take_damage(999), p.is_alive()][2] == False"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Give `Player` a `__str__(self)` that returns text like `Zed: 100 HP`. Then create a player and `print` it.",
		Starter:     "class Player:\n    def __init__(self, name):\n        self.name = name\n        self.hp = 100\n\n    def take_damage(self, amount):\n        self.hp = max(0, self.hp - amount)\n\n    def is_alive(self):\n        return self.hp > 0\n\n    def __str__(self):\n        pass\n\n# Create a Player and print it\n",
		Checks:      `[{"text":"str(player) reads like 'Zed: 100 HP'.","test":"str(Player('Zed')) == 'Zed: 100 HP'"},{"text":"You printed a player.","test":"'HP' in _out"}]`,
	},
}

// pythonLabs are the interactive Python labs (run via Pyodide/WASM).
var pythonLabs = map[string][]store.Step{
	"workshop_python_warm_up":         pythonWarmupSteps,
	"workshop_build_an_xp_calculator": xpCalcSteps,
	"workshop_build_a_loot_inventory": lootInventorySteps,
	"workshop_score_multipliers":      comboSteps,
	"workshop_model_a_player":         playerClassSteps,
}

// labSteps maps a frontend lesson slug to its interactive steps.
var labSteps = map[string][]store.Step{
	"workshop_code_a_level_up_system":           levelUpSteps,
	"workshop_build_an_inventory":               inventorySteps,
	"workshop_build_a_like_button":              likeButtonSteps,
	"workshop_build_a_game_launch_page":         gameLaunchSteps,
	"workshop_build_a_sign_up_form":             signupFormSteps,
	"workshop_style_a_business_card":            businessCardSteps,
	"workshop_build_a_pricing_card":             pricingCardSteps,
	"workshop_build_an_achievement_badge":       achievementBadgeSteps,
	"workshop_build_an_xp_progress_bar":         xpBarSteps,
	"workshop_build_a_nav_bar_with_flexbox":     navBarSteps,
	"workshop_build_a_leaderboard_row":          leaderboardRowSteps,
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
