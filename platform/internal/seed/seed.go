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

const productCatalogScaffold = "<label for=\"search\">Search</label>\n<input id=\"search\" placeholder=\"Search products\">\n<ul id=\"products\"></ul>\n<p id=\"empty\" hidden>No matching products.</p>\n"

var productCatalogSteps = []store.Step{
	{
		Lang:        "js",
		Scaffold:    productCatalogScaffold,
		Instruction: "Create a `products` array with at least **4 objects** (`name` and `category`), then render one `<li>` per product into `#products`.",
		Starter:     "// Create a products array and render each product name into #products as an <li>\n",
		Checks:      `[{"text":"products should be an array of objects with name and category.","test":"Array.isArray(products) && products.length >= 4 && products.every(function(p){return p && typeof p.name === 'string' && typeof p.category === 'string';})"},{"text":"The product list should render one li per product.","test":"document.querySelectorAll('#products li').length === products.length"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    productCatalogScaffold,
		Instruction: "Extract the rendering into a function `renderProducts(items)` that clears the list and renders whichever items it receives.",
		Starter:     "const list = document.querySelector('#products');\n\nconst products = [\n  { name: 'Notebook', category: 'Office' },\n  { name: 'Monitor', category: 'Tech' },\n  { name: 'Mouse Pad', category: 'Tech' },\n  { name: 'Water Bottle', category: 'Lifestyle' },\n];\n\nproducts.forEach((product) => {\n  const item = document.createElement('li');\n  item.textContent = product.name;\n  list.appendChild(item);\n});\n\n// Write renderProducts(items) and call it with products\n",
		Checks:      `[{"text":"renderProducts should be a function.","test":"typeof renderProducts === 'function'"},{"text":"renderProducts should clear the list and render the provided items.","test":"(function(){renderProducts(products.slice(0, 2));return document.querySelectorAll('#products li').length === 2;})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    productCatalogScaffold,
		Instruction: "Add an `input` listener on `#search` that filters products by name, **case-insensitive**, and calls `renderProducts(filtered)`.",
		Starter:     "const search = document.querySelector('#search');\nconst list = document.querySelector('#products');\n\nconst products = [\n  { name: 'Notebook', category: 'Office' },\n  { name: 'Monitor', category: 'Tech' },\n  { name: 'Mouse Pad', category: 'Tech' },\n  { name: 'Water Bottle', category: 'Lifestyle' },\n];\n\nfunction renderProducts(items) {\n  list.innerHTML = '';\n  items.forEach((product) => {\n    const item = document.createElement('li');\n    item.textContent = product.name;\n    list.appendChild(item);\n  });\n}\n\nrenderProducts(products);\n\n// Add an input listener that filters products by name and re-renders\n",
		Checks:      `[{"text":"Typing into search should filter the rendered products by name.","test":"(function(){var input=document.querySelector('#search');if(!input)return false;input.value='note';input.dispatchEvent(new Event('input',{bubbles:true}));var items=Array.from(document.querySelectorAll('#products li'));return items.length > 0 && items.every(function(li){return /note/i.test(li.textContent);});})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    productCatalogScaffold,
		Instruction: "Show the `#empty` message when nothing matches, and hide it again when results exist.",
		Starter:     "const search = document.querySelector('#search');\nconst list = document.querySelector('#products');\nconst empty = document.querySelector('#empty');\n\nconst products = [\n  { name: 'Notebook', category: 'Office' },\n  { name: 'Monitor', category: 'Tech' },\n  { name: 'Mouse Pad', category: 'Tech' },\n  { name: 'Water Bottle', category: 'Lifestyle' },\n];\n\nfunction renderProducts(items) {\n  list.innerHTML = '';\n  items.forEach((product) => {\n    const item = document.createElement('li');\n    item.textContent = product.name;\n    list.appendChild(item);\n  });\n}\n\nrenderProducts(products);\n\nsearch.addEventListener('input', () => {\n  const term = search.value.toLowerCase();\n  const filtered = products.filter((product) => product.name.toLowerCase().includes(term));\n  renderProducts(filtered);\n\n  // Toggle the empty-state message here\n});\n",
		Checks:      `[{"text":"The empty-state message should show with no matches and hide when matches return.","test":"(function(){var input=document.querySelector('#search'),empty=document.querySelector('#empty');if(!input||!empty)return false;input.value='zzzzz';input.dispatchEvent(new Event('input',{bubbles:true}));var shown=empty.hidden===false;input.value='note';input.dispatchEvent(new Event('input',{bubbles:true}));var hiddenAgain=empty.hidden===true;return shown && hiddenAgain;})()"}]`,
	},
}

const passwordStrengthScaffold = "<label for=\"password\">Password</label>\n<input id=\"password\" type=\"password\">\n<p id=\"strength\">Enter a password</p>\n<div id=\"meter\" style=\"width:220px;height:10px;background:#e5e7eb;border-radius:999px;overflow:hidden;\"><span id=\"fill\" style=\"display:block;height:100%;width:0%;background:#ef4444;\"></span></div>\n"

var passwordStrengthSteps = []store.Step{
	{
		Lang:        "js",
		Scaffold:    passwordStrengthScaffold,
		Instruction: "Add an `input` listener on `#password` that updates `#strength` to show the current password length.",
		Starter:     "// Select #password and #strength, then show the current password length while typing\n",
		Checks:      `[{"text":"Typing should update the strength text with the current length.","test":"(function(){var input=document.querySelector('#password'),strength=document.querySelector('#strength');if(!input||!strength)return false;input.value='abc';input.dispatchEvent(new Event('input',{bubbles:true}));return strength.textContent.indexOf('3') !== -1;})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    passwordStrengthScaffold,
		Instruction: "Also update `#fill` so its width is `password length × 10`, capped at `100%`.",
		Starter:     "const password = document.querySelector('#password');\nconst strength = document.querySelector('#strength');\n\npassword.addEventListener('input', () => {\n  strength.textContent = 'Length: ' + password.value.length;\n});\n\n// Also update #fill width here\n",
		Checks:      `[{"text":"The meter fill should grow as the password gets longer.","test":"(function(){var input=document.querySelector('#password'),fill=document.querySelector('#fill');if(!input||!fill)return false;input.value='hello';input.dispatchEvent(new Event('input',{bubbles:true}));return fill.style.width === '50%';})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    passwordStrengthScaffold,
		Instruction: "If the password is shorter than `8` characters, show `Too short`. Otherwise show `Good length`.",
		Starter:     "const password = document.querySelector('#password');\nconst strength = document.querySelector('#strength');\nconst fill = document.querySelector('#fill');\n\npassword.addEventListener('input', () => {\n  const width = Math.min(password.value.length * 10, 100);\n  fill.style.width = width + '%';\n\n  // Update the message based on the password length\n});\n",
		Checks:      `[{"text":"Short passwords should show a warning, and longer ones should show a positive message.","test":"(function(){var input=document.querySelector('#password'),strength=document.querySelector('#strength');if(!input||!strength)return false;input.value='abc';input.dispatchEvent(new Event('input',{bubbles:true}));var shortOk=/short/i.test(strength.textContent);input.value='abcdefgh';input.dispatchEvent(new Event('input',{bubbles:true}));var longOk=/good length/i.test(strength.textContent) || /strong/i.test(strength.textContent);return shortOk && longOk;})()"}]`,
	},
	{
		Lang:        "js",
		Scaffold:    passwordStrengthScaffold,
		Instruction: "If the password has at least `10` characters, **one number**, and **one special character**, show `Strong password` and make sure the meter reaches `100%`.",
		Starter:     "const password = document.querySelector('#password');\nconst strength = document.querySelector('#strength');\nconst fill = document.querySelector('#fill');\n\npassword.addEventListener('input', () => {\n  const value = password.value;\n  const width = Math.min(value.length * 10, 100);\n  fill.style.width = width + '%';\n\n  if (value.length < 8) {\n    strength.textContent = 'Too short';\n    return;\n  }\n\n  strength.textContent = 'Good length';\n\n  // Upgrade the message for a strong password\n});\n",
		Checks:      `[{"text":"A strong password should show a strong message and fill the meter completely.","test":"(function(){var input=document.querySelector('#password'),strength=document.querySelector('#strength'),fill=document.querySelector('#fill');if(!input||!strength||!fill)return false;input.value='abc123!xyz';input.dispatchEvent(new Event('input',{bubbles:true}));return /strong/i.test(strength.textContent) && fill.style.width === '100%';})()"}]`,
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

// chatFilterSteps — string methods via a chat command filter.
var chatFilterSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Chat needs a hype command. Write `shout(msg)` that returns `msg` in UPPERCASE with a `!` added to the end.",
		Starter:     "def shout(msg):\n    pass\n",
		Checks:      `[{"text":"shout should be a function.","test":"callable(shout)"},{"text":"shout('hi') should be 'HI!'.","test":"shout('hi') == 'HI!'"},{"text":"shout('go') should be 'GO!'.","test":"shout('go') == 'GO!'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `is_command(msg)` that returns `True` when `msg` starts with a slash (`/`), otherwise `False`. Try the `.startswith()` string method.",
		Starter:     "def is_command(msg):\n    pass\n",
		Checks:      `[{"text":"'/help' is a command.","test":"is_command('/help') == True"},{"text":"'hello' is not a command.","test":"is_command('hello') == False"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `censor(msg, bad)` that returns `msg` with every occurrence of the word `bad` replaced by `***`. The `.replace()` method does this in one call.",
		Starter:     "def censor(msg, bad):\n    pass\n",
		Checks:      `[{"text":"The bad word gets starred out.","test":"censor('you noob', 'noob') == 'you ***'"},{"text":"Clean messages are untouched.","test":"censor('gg wp', 'noob') == 'gg wp'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Finally, write `word_count(msg)` that returns how many words are in `msg`. An empty string has `0` words. Hint: `.split()` breaks text into a list of words.",
		Starter:     "def word_count(msg):\n    pass\n",
		Checks:      `[{"text":"Three words counts as 3.","test":"word_count('hello there friend') == 3"},{"text":"An empty message has 0 words.","test":"word_count('') == 0"}]`,
	},
}

// tallySteps — dictionaries, aggregation, and sorting via a scoreboard.
var tallySteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Match results arrive as `(name, points)` pairs. Write `tally(matches)` that returns a dictionary mapping each name to their **total** points across all pairs. An empty list gives an empty dict.",
		Starter:     "def tally(matches):\n    # matches is a list of (name, points) tuples\n    pass\n",
		Checks:      `[{"text":"Points add up per player.","test":"tally([('a', 5), ('b', 3), ('a', 2)]) == {'a': 7, 'b': 3}"},{"text":"No matches means an empty scoreboard.","test":"tally([]) == {}"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `leader(totals)` that returns the name with the highest total from a `totals` dictionary.",
		Starter:     "def leader(totals):\n    pass\n",
		Checks:      `[{"text":"'a' leads with 7.","test":"leader({'a': 7, 'b': 3}) == 'a'"},{"text":"'y' leads with 9.","test":"leader({'x': 1, 'y': 9}) == 'y'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `above(totals, n)` that returns the list of names whose total is **greater than** `n`, sorted alphabetically.",
		Starter:     "def above(totals, n):\n    pass\n",
		Checks:      `[{"text":"Only players above 5 make the cut.","test":"above({'a': 7, 'b': 3, 'c': 10}, 5) == ['a', 'c']"},{"text":"A high bar can exclude everyone.","test":"above({'a': 7, 'b': 3}, 100) == []"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Print the scoreboard: loop over a `totals` dict and print each `name: points` line.",
		Starter:     "totals = {'Nova': 7, 'Rex': 3}\n# Print each player and their score, like: Nova: 7\n",
		Checks:      `[{"text":"Nova's score is printed.","test":"'Nova: 7' in _out"},{"text":"Rex's score is printed.","test":"'Rex: 3' in _out"}]`,
	},
}

// healthBarSteps — loops, integer math, and conditionals via a health bar.
var healthBarSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Draw a health bar. Write `bar(filled, total)` that returns a string of `filled` `#` characters followed by enough `-` characters to reach `total` length.",
		Starter:     "def bar(filled, total):\n    pass\n",
		Checks:      `[{"text":"3 of 5 looks like '###--'.","test":"bar(3, 5) == '###--'"},{"text":"An empty bar is all dashes.","test":"bar(0, 4) == '----'"},{"text":"A full bar is all hashes.","test":"bar(5, 5) == '#####'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `percent(hp, maxhp)` that returns the health percentage as a whole number, rounded **down**. Use integer division (`//`).",
		Starter:     "def percent(hp, maxhp):\n    pass\n",
		Checks:      `[{"text":"50 of 200 is 25%.","test":"percent(50, 200) == 25"},{"text":"Full health is 100%.","test":"percent(200, 200) == 100"},{"text":"No health is 0%.","test":"percent(0, 10) == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `status(hp, maxhp)` that returns `'DEAD'` when hp is `0`, `'LOW'` when the percentage is under 25, and `'OK'` otherwise. Check for `0` first.",
		Starter:     "def status(hp, maxhp):\n    pass\n",
		Checks:      `[{"text":"0 HP is DEAD.","test":"status(0, 100) == 'DEAD'"},{"text":"10% is LOW.","test":"status(10, 100) == 'LOW'"},{"text":"80% is OK.","test":"status(80, 100) == 'OK'"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Put it together: for `hp = 30`, `maxhp = 100`, print a line that shows the bar and the percentage, like `[###-------] 30%`.",
		Starter:     "hp = 30\nmaxhp = 100\n# Print a health bar line ending in the percentage, e.g. [###-------] 30%\n",
		Checks:      `[{"text":"You show the percentage.","test":"'30%' in _out"},{"text":"You draw a bar with # and -.","test":"'#' in _out and '-' in _out"}]`,
	},
}

// combatMathSteps — default args, *args, and **kwargs via combat helpers.
var combatMathSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Write `damage(base, multiplier=1)` that returns `base * multiplier`. The `multiplier` should **default** to `1` when it isn't given.",
		Starter:     "def damage(base, multiplier=1):\n    pass\n",
		Checks:      `[{"text":"With no multiplier, damage is the base.","test":"damage(10) == 10"},{"text":"A 2x multiplier doubles it.","test":"damage(10, 2) == 20"},{"text":"You can pass the multiplier by name.","test":"damage(10, multiplier=3) == 30"}]`,
	},
	{
		Lang:        "python",
		Instruction: "A combo lands several hits. Write `total_damage(*hits)` that returns the sum of **any number** of hit values. No hits totals `0`.",
		Starter:     "def total_damage(*hits):\n    pass\n",
		Checks:      `[{"text":"Three hits add up.","test":"total_damage(1, 2, 3) == 6"},{"text":"No hits deal 0.","test":"total_damage() == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `buff(base, **mods)` that returns `base` plus the sum of all the keyword modifier values passed in.",
		Starter:     "def buff(base, **mods):\n    pass\n",
		Checks:      `[{"text":"Modifiers stack onto the base.","test":"buff(100, atk=10, spd=5) == 115"},{"text":"No modifiers leaves the base.","test":"buff(100) == 100"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Finally, write `strongest(*hits)` that returns the biggest hit — or `0` when there are no hits.",
		Starter:     "def strongest(*hits):\n    pass\n",
		Checks:      `[{"text":"It finds the biggest hit.","test":"strongest(3, 9, 4) == 9"},{"text":"No hits returns 0.","test":"strongest() == 0"}]`,
	},
}

// safeParserSteps — try/except and specific exceptions via input parsing.
var safeParserSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Player input is messy. Write `to_int(text)` that returns `int(text)` when it can, but returns `0` if the text isn't a number. Wrap it in `try` / `except ValueError`.",
		Starter:     "def to_int(text):\n    pass\n",
		Checks:      `[{"text":"Numbers convert.","test":"to_int('42') == 42"},{"text":"Junk becomes 0.","test":"to_int('abc') == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `safe_div(a, b)` that returns `a / b`, but returns `None` when `b` is `0`. Catch `ZeroDivisionError`.",
		Starter:     "def safe_div(a, b):\n    pass\n",
		Checks:      `[{"text":"Normal division works.","test":"safe_div(10, 2) == 5"},{"text":"Dividing by zero is safe.","test":"safe_div(5, 0) is None"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `get_stat(stats, key)` that returns `stats[key]`, or `0` if the key is missing. Catch `KeyError`.",
		Starter:     "def get_stat(stats, key):\n    pass\n",
		Checks:      `[{"text":"Existing stats come back.","test":"get_stat({'hp': 10}, 'hp') == 10"},{"text":"Missing stats default to 0.","test":"get_stat({}, 'mp') == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `parse_all(texts)` that turns a list of strings into a list of integers, **skipping** any that aren't numbers. Reuse the idea from `to_int`.",
		Starter:     "def parse_all(texts):\n    pass\n",
		Checks:      `[{"text":"Only the numbers survive.","test":"parse_all(['1', 'x', '3']) == [1, 3]"},{"text":"An empty list stays empty.","test":"parse_all([]) == []"}]`,
	},
}

// numberStreamSteps — generators and lazy iteration.
var numberStreamSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Write a **generator** `countdown(n)` that `yield`s the numbers from `n` down to `1`.",
		Starter:     "def countdown(n):\n    pass\n",
		Checks:      `[{"text":"countdown(3) yields 3, 2, 1.","test":"list(countdown(3)) == [3, 2, 1]"},{"text":"countdown(1) yields just 1.","test":"list(countdown(1)) == [1]"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write a generator `evens(limit)` that yields the even numbers starting at `0` that are **less than** `limit`.",
		Starter:     "def evens(limit):\n    pass\n",
		Checks:      `[{"text":"evens(6) yields 0, 2, 4.","test":"list(evens(6)) == [0, 2, 4]"},{"text":"evens(1) yields just 0.","test":"list(evens(1)) == [0]"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Generators are lazy, so you can pull just a few items. Write `take(gen, k)` that returns a list of the first `k` items from any iterator `gen` (or fewer, if it runs out).",
		Starter:     "def countdown(n):\n    while n > 0:\n        yield n\n        n -= 1\n\ndef take(gen, k):\n    pass\n",
		Checks:      `[{"text":"It grabs the first 3.","test":"take(countdown(100), 3) == [100, 99, 98]"},{"text":"It stops early if the stream is short.","test":"take(countdown(2), 5) == [2, 1]"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write an **infinite** generator `naturals()` that yields `1, 2, 3, …` forever. Then your `take` can safely grab the first few.",
		Starter:     "def take(gen, k):\n    out = []\n    for i, x in enumerate(gen):\n        if i >= k:\n            break\n        out.append(x)\n    return out\n\ndef naturals():\n    pass\n",
		Checks:      `[{"text":"The first 5 naturals are 1..5.","test":"take(naturals(), 5) == [1, 2, 3, 4, 5]"},{"text":"The first natural is 1.","test":"take(naturals(), 1) == [1]"}]`,
	},
}

// decoratorSteps — decorators that wrap and augment functions.
var decoratorSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "A **decorator** wraps a function to change its behavior. Write `double_result(func)` that returns a wrapper which calls `func` and returns **double** its result.",
		Starter:     "def double_result(func):\n    pass\n\n@double_result\ndef total(a, b):\n    return a + b\n",
		Checks:      `[{"text":"total is still callable.","test":"callable(total)"},{"text":"2 + 3, doubled, is 10.","test":"total(2, 3) == 10"},{"text":"0 + 0, doubled, is 0.","test":"total(0, 0) == 0"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `count_calls(func)` — a decorator whose wrapper also tracks how many times it was called on a `.calls` attribute (starting at `0`).",
		Starter:     "def count_calls(func):\n    pass\n\n@count_calls\ndef ping():\n    return 'pong'\n",
		Checks:      `[{"text":"The wrapped function still works.","test":"ping() == 'pong'"},{"text":"Calling it 3 times sets .calls to 3.","test":"[p := count_calls(lambda: None), p(), p(), p(), p.calls][4] == 3"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `ensure_positive(func)` whose wrapper returns `0` when the first argument is negative, and otherwise calls `func` normally.",
		Starter:     "def ensure_positive(func):\n    pass\n\n@ensure_positive\ndef square(n):\n    return n * n\n",
		Checks:      `[{"text":"A positive input runs normally.","test":"square(4) == 16"},{"text":"A negative input is blocked, returning 0.","test":"square(-3) == 0"}]`,
	},
}

// recursionSteps — base cases and recursive cases.
var recursionSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Write `factorial(n)` **recursively**: `factorial(0)` is `1`, and otherwise `n * factorial(n - 1)`.",
		Starter:     "def factorial(n):\n    pass\n",
		Checks:      `[{"text":"factorial(5) is 120.","test":"factorial(5) == 120"},{"text":"factorial(0) is 1 (the base case).","test":"factorial(0) == 1"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `fib(n)` recursively for the Fibonacci sequence, where `fib(0)` is `0`, `fib(1)` is `1`, and each later value is the sum of the two before it.",
		Starter:     "def fib(n):\n    pass\n",
		Checks:      `[{"text":"fib(0) is 0.","test":"fib(0) == 0"},{"text":"fib(1) is 1.","test":"fib(1) == 1"},{"text":"fib(7) is 13.","test":"fib(7) == 13"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `power(base, exp)` recursively — `base` multiplied by itself `exp` times. Anything to the power `0` is `1`.",
		Starter:     "def power(base, exp):\n    pass\n",
		Checks:      `[{"text":"2 to the 10th is 1024.","test":"power(2, 10) == 1024"},{"text":"5 to the 0th is 1.","test":"power(5, 0) == 1"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `deep_sum(items)` that adds up every number in a **nested** list (lists inside lists, any depth). Hint: if an item is a list, recurse into it.",
		Starter:     "def deep_sum(items):\n    pass\n",
		Checks:      `[{"text":"It sums through the nesting.","test":"deep_sum([1, [2, [3, 4]], 5]) == 15"},{"text":"An empty list sums to 0.","test":"deep_sum([]) == 0"}]`,
	},
}

// regexSteps — the re module for parsing text.
var regexSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Write `find_numbers(text)` that returns a list of every run of digits in `text` as strings. Use `re.findall` with the pattern `\\d+`.",
		Starter:     "import re\n\ndef find_numbers(text):\n    pass\n",
		Checks:      `[{"text":"It finds both numbers.","test":"find_numbers('hp 30 mp 5') == ['30', '5']"},{"text":"No digits means an empty list.","test":"find_numbers('none here') == []"}]`,
	},
	{
		Lang:        "python",
		Instruction: "A player tag looks like `#ABC1234` — a `#`, then exactly 3 uppercase letters, then exactly 4 digits. Write `is_valid_tag(s)` that returns `True`/`False` (try `re.fullmatch`, and wrap it in `bool(...)`).",
		Starter:     "import re\n\ndef is_valid_tag(s):\n    pass\n",
		Checks:      `[{"text":"A well-formed tag passes.","test":"is_valid_tag('#ABC1234') == True"},{"text":"Lowercase fails.","test":"is_valid_tag('#abc1234') == False"},{"text":"Junk fails.","test":"is_valid_tag('nope') == False"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Given a line like `Zed hit Rex for 42 damage`, write `extract_damage(line)` that returns the damage as an **int**. Use `re.search` with a capture group `(\\d+)`.",
		Starter:     "import re\n\ndef extract_damage(line):\n    pass\n",
		Checks:      `[{"text":"It pulls out 42.","test":"extract_damage('Zed hit Rex for 42 damage') == 42"},{"text":"It works on any line.","test":"extract_damage('Boss hit you for 7 damage') == 7"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Write `redact(text)` that replaces every digit with a `#`. Use `re.sub`.",
		Starter:     "import re\n\ndef redact(text):\n    pass\n",
		Checks:      `[{"text":"Digits become hashes.","test":"redact('lvl 12') == 'lvl ##'"},{"text":"Text with no digits is unchanged.","test":"redact('no digits') == 'no digits'"}]`,
	},
}

// dataclassSteps — @dataclass for tidy data objects.
var dataclassSteps = []store.Step{
	{
		Lang:        "python",
		Instruction: "Dataclasses give you tidy data objects with no boilerplate. Make an `@dataclass` called `Item` with two fields: `name` (a `str`) and `price` (an `int`).",
		Starter:     "from dataclasses import dataclass\n\n@dataclass\nclass Item:\n    pass\n",
		Checks:      `[{"text":"Item can be created.","test":"callable(Item)"},{"text":"It keeps the name.","test":"Item('sword', 100).name == 'sword'"},{"text":"It keeps the price.","test":"Item('sword', 100).price == 100"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Add a third field `qty` (an `int`) that **defaults** to `1`. Dataclasses also compare by value for free.",
		Starter:     "from dataclasses import dataclass\n\n@dataclass\nclass Item:\n    name: str\n    price: int\n    # add qty with a default of 1\n",
		Checks:      `[{"text":"qty defaults to 1.","test":"Item('a', 5).qty == 1"},{"text":"You can still set qty.","test":"Item('a', 5, 3).qty == 3"},{"text":"Equal fields mean equal items.","test":"Item('a', 5) == Item('a', 5)"}]`,
	},
	{
		Lang:        "python",
		Instruction: "Give `Item` a method `total(self)` that returns `price * qty` — the cost of the whole stack.",
		Starter:     "from dataclasses import dataclass\n\n@dataclass\nclass Item:\n    name: str\n    price: int\n    qty: int = 1\n\n    def total(self):\n        pass\n",
		Checks:      `[{"text":"3 items at 5 each is 15.","test":"Item('a', 5, 3).total() == 15"},{"text":"A single item's total is its price.","test":"Item('a', 10).total() == 10"}]`,
	},
}

// --- Server-side labs (Lang "pyserver"): the learner writes a real web app or
// database program that runs on the server; the checks exercise it in-process
// via the framework's test client. ---

// flaskRouteSteps — a first Flask web server.
var flaskRouteSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Time to build a real web server! Create a Flask app named `app`, then add a route for `/` that returns `Hello, GamifyDev!`. The checker runs your app for you — you don't need `app.run()`.",
		Starter:     "from flask import Flask\n\napp = Flask(__name__)\n\n# Add a route for '/' that returns 'Hello, GamifyDev!'\n",
		Checks:      `[{"text":"The home page loads (status 200).","test":"client.get('/').status_code == 200"},{"text":"It returns 'Hello, GamifyDev!'.","test":"b'Hello, GamifyDev!' in client.get('/').data"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Add a second route, `/ping`, that returns the text `pong`.",
		Starter:     "from flask import Flask\n\napp = Flask(__name__)\n\n@app.route('/')\ndef home():\n    return 'Hello, GamifyDev!'\n\n# Add a /ping route that returns 'pong'\n",
		Checks:      `[{"text":"/ping loads (status 200).","test":"client.get('/ping').status_code == 200"},{"text":"/ping returns exactly 'pong'.","test":"client.get('/ping').data == b'pong'"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Routes can capture parts of the URL. Add `/hi/<name>` that returns `Hi, <name>!` — so `/hi/Ada` returns `Hi, Ada!`.",
		Starter:     "from flask import Flask\n\napp = Flask(__name__)\n\n# Add /hi/<name> that returns 'Hi, <name>!'\n",
		Checks:      `[{"text":"/hi/Ada greets Ada.","test":"b'Hi, Ada!' in client.get('/hi/Ada').data"},{"text":"/hi/Zed greets Zed.","test":"b'Hi, Zed!' in client.get('/hi/Zed').data"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "A route can set the HTTP status code by returning a tuple of `(body, status)`. Add `/created` that returns the body `Created` with status `201`.",
		Starter:     "from flask import Flask\n\napp = Flask(__name__)\n\n# Add /created that returns 'Created' with status code 201\n",
		Checks:      `[{"text":"/created returns status 201.","test":"client.get('/created').status_code == 201"},{"text":"Its body is 'Created'.","test":"b'Created' in client.get('/created').data"}]`,
	},
}

// flaskRequestSteps — reading query strings, form/JSON, and methods.
var flaskRequestSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Read query strings with `request.args`. Add `/greet` that returns `Hello, <name>!` using the `name` query parameter, defaulting to `friend` when it's missing.",
		Starter:     "from flask import Flask, request\n\napp = Flask(__name__)\n\n# /greet -> 'Hello, <name>!', reading ?name=... (default 'friend')\n",
		Checks:      `[{"text":"No name defaults to 'friend'.","test":"b'Hello, friend!' in client.get('/greet').data"},{"text":"?name=Ada greets Ada.","test":"b'Hello, Ada!' in client.get('/greet?name=Ada').data"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Add `/sum` that reads two integer query params `a` and `b` and returns their sum as text. Query values are strings — convert them with `int(...)`.",
		Starter:     "from flask import Flask, request\n\napp = Flask(__name__)\n\n# /sum?a=2&b=3 -> '5'\n",
		Checks:      `[{"text":"2 + 3 is 5.","test":"client.get('/sum?a=2&b=3').data == b'5'"},{"text":"10 + 15 is 25.","test":"client.get('/sum?a=10&b=15').data == b'25'"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Handle JSON. Add a **POST** route `/echo` that reads a JSON body `{\"msg\": ...}` and returns JSON `{\"echo\": <msg>}`. Use `request.get_json()` and `jsonify`.",
		Starter:     "from flask import Flask, request, jsonify\n\napp = Flask(__name__)\n\n# POST /echo: read {'msg': ...} and return {'echo': <msg>}\n",
		Checks:      `[{"text":"It echoes the message back as JSON.","test":"client.post('/echo', json={'msg': 'hi'}).get_json() == {'echo': 'hi'}"},{"text":"It works for any message.","test":"client.post('/echo', json={'msg': 'yo'}).get_json() == {'echo': 'yo'}"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Flask knows which methods a route allows. Confirm your `/echo` only accepts POST — a `GET` to it should return status `405` (Method Not Allowed).",
		Starter:     "from flask import Flask, request, jsonify\n\napp = Flask(__name__)\n\n@app.post('/echo')\ndef echo():\n    return jsonify({'echo': request.get_json()['msg']})\n\n# Nothing to add — run the check to confirm GET is rejected with 405.\n",
		Checks:      `[{"text":"A GET to /echo is rejected (405).","test":"client.get('/echo').status_code == 405"},{"text":"A POST to /echo still works.","test":"client.post('/echo', json={'msg': 'x'}).get_json() == {'echo': 'x'}"}]`,
	},
}

// restApiSteps — a full CRUD REST resource.
var restApiSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Build a REST API for tasks. Using the `tasks` list provided, add `GET /tasks` that returns the whole list as JSON.",
		Starter:     "from flask import Flask, request, jsonify\n\napp = Flask(__name__)\ntasks = []\n\n# GET /tasks -> the tasks list as JSON\n",
		Checks:      `[{"text":"GET /tasks returns 200.","test":"client.get('/tasks').status_code == 200"},{"text":"It starts empty.","test":"client.get('/tasks').get_json() == []"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Add `POST /tasks` that reads `{\"title\": ...}`, appends `{\"id\": <n>, \"title\": <title>, \"done\": False}` to `tasks`, and returns the new task with status `201`. A task's position (`len(tasks) + 1`) makes a fine id.",
		Starter:     "from flask import Flask, request, jsonify\n\napp = Flask(__name__)\ntasks = []\n\n@app.get('/tasks')\ndef list_tasks():\n    return jsonify(tasks)\n\n# POST /tasks: create a task from {'title': ...} and return it with status 201\n",
		Checks:      `[{"text":"Creating a task returns 201.","test":"client.post('/tasks', json={'title': 'Learn Flask'}).status_code == 201"},{"text":"The new task echoes its title.","test":"client.post('/tasks', json={'title': 'Ship it'}).get_json()['title'] == 'Ship it'"},{"text":"A new task starts not done.","test":"client.post('/tasks', json={'title': 'x'}).get_json()['done'] == False"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Add `GET /tasks/<int:id>` that returns a single task by its id, or a `404` if there's no such task. `abort(404)` is the easy way to bail out.",
		Starter:     "from flask import Flask, request, jsonify, abort\n\napp = Flask(__name__)\ntasks = []\n\n@app.post('/tasks')\ndef create_task():\n    task = {'id': len(tasks) + 1, 'title': request.get_json()['title'], 'done': False}\n    tasks.append(task)\n    return jsonify(task), 201\n\n# GET /tasks/<int:id> -> the matching task, or abort(404)\n",
		Checks:      `[{"text":"A created task can be fetched by id.","test":"(tid := client.post('/tasks', json={'title': 'Find me'}).get_json()['id']) and client.get(f'/tasks/{tid}').get_json()['title'] == 'Find me'"},{"text":"An unknown id returns 404.","test":"client.get('/tasks/9999').status_code == 404"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Complete the CRUD set: add `DELETE /tasks/<int:id>` that removes a task and returns status `204`. Deleting a task then fetching it should give `404`.",
		Starter:     "from flask import Flask, request, jsonify, abort\n\napp = Flask(__name__)\ntasks = []\n\n@app.post('/tasks')\ndef create_task():\n    task = {'id': len(tasks) + 1, 'title': request.get_json()['title'], 'done': False}\n    tasks.append(task)\n    return jsonify(task), 201\n\n@app.get('/tasks/<int:id>')\ndef get_task(id):\n    for t in tasks:\n        if t['id'] == id:\n            return jsonify(t)\n    abort(404)\n\n# DELETE /tasks/<int:id> -> remove it, return status 204\n",
		Checks:      `[{"text":"Deleting a task returns 204.","test":"(tid := client.post('/tasks', json={'title': 'Bye'}).get_json()['id']) and client.delete(f'/tasks/{tid}').status_code == 204"},{"text":"After deleting, it's gone (404).","test":"(tid := client.post('/tasks', json={'title': 'Gone'}).get_json()['id']) and (client.delete(f'/tasks/{tid}'), client.get(f'/tasks/{tid}').status_code)[1] == 404"}]`,
	},
}

// sqliteSteps — persisting and querying data with SQLite (stdlib).
var sqliteSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Databases hold your app's data. Write `setup_db()` that returns a new in-memory SQLite connection with a `players` table (`name TEXT, score INT`) seeded with `('Ada', 90)` and `('Zed', 70)`.",
		Starter:     "import sqlite3\n\ndef setup_db():\n    conn = sqlite3.connect(':memory:')\n    # create the players table and insert Ada (90) and Zed (70)\n    return conn\n",
		Checks:      `[{"text":"The table has 2 players.","test":"setup_db().execute('SELECT COUNT(*) FROM players').fetchone()[0] == 2"},{"text":"Ada's score is stored.","test":"setup_db().execute(\"SELECT score FROM players WHERE name = 'Ada'\").fetchone()[0] == 90"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Write `top_player(conn)` that returns the **name** of the highest-scoring player. `ORDER BY score DESC LIMIT 1` does the work in SQL.",
		Starter:     "import sqlite3\n\ndef setup_db():\n    conn = sqlite3.connect(':memory:')\n    conn.execute('CREATE TABLE players (name TEXT, score INT)')\n    conn.executemany('INSERT INTO players VALUES (?, ?)', [('Ada', 90), ('Zed', 70)])\n    return conn\n\ndef top_player(conn):\n    pass\n",
		Checks:      `[{"text":"Ada is on top.","test":"top_player(setup_db()) == 'Ada'"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Write `add_player(conn, name, score)` that inserts a new player, and `count_players(conn)` that returns how many players there are.",
		Starter:     "import sqlite3\n\ndef setup_db():\n    conn = sqlite3.connect(':memory:')\n    conn.execute('CREATE TABLE players (name TEXT, score INT)')\n    conn.executemany('INSERT INTO players VALUES (?, ?)', [('Ada', 90), ('Zed', 70)])\n    return conn\n\ndef add_player(conn, name, score):\n    pass\n\ndef count_players(conn):\n    pass\n",
		Checks:      `[{"text":"A fresh database has 2 players.","test":"count_players(setup_db()) == 2"},{"text":"Adding a player bumps the count to 3.","test":"(c := setup_db(), add_player(c, 'Nova', 80), count_players(c))[2] == 3"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Write `leaderboard(conn)` that returns a list of `(name, score)` tuples, highest score first.",
		Starter:     "import sqlite3\n\ndef setup_db():\n    conn = sqlite3.connect(':memory:')\n    conn.execute('CREATE TABLE players (name TEXT, score INT)')\n    conn.executemany('INSERT INTO players VALUES (?, ?)', [('Ada', 90), ('Zed', 70)])\n    return conn\n\ndef leaderboard(conn):\n    pass\n",
		Checks:      `[{"text":"The leaderboard is sorted high to low.","test":"leaderboard(setup_db()) == [('Ada', 90), ('Zed', 70)]"}]`,
	},
}

// fastapiSteps — a modern, type-hinted API with FastAPI + Pydantic.
var fastapiSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "FastAPI is a modern, type-hinted web framework. Create `app = FastAPI()` and a `GET /` route that returns the JSON `{'msg': 'Hello from FastAPI'}`.",
		Starter:     "from fastapi import FastAPI\n\napp = FastAPI()\n\n# GET / -> {'msg': 'Hello from FastAPI'}\n",
		Checks:      `[{"text":"The root returns 200.","test":"client.get('/').status_code == 200"},{"text":"It returns the greeting JSON.","test":"client.get('/').json() == {'msg': 'Hello from FastAPI'}"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "FastAPI reads typed path parameters for you. Add `GET /square/{n}` with `n: int` that returns `{'n': n, 'square': n * n}`.",
		Starter:     "from fastapi import FastAPI\n\napp = FastAPI()\n\n# GET /square/{n} with n: int -> {'n': n, 'square': n * n}\n",
		Checks:      `[{"text":"square of 4 is 16.","test":"client.get('/square/4').json() == {'n': 4, 'square': 16}"},{"text":"square of 9 is 81.","test":"client.get('/square/9').json()['square'] == 81"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Add a typed request body with Pydantic. A `Msg` model with a `text: str` field is provided. Add `POST /echo` taking a `Msg` that returns `{'echo': <text>, 'length': <len of text>}`.",
		Starter:     "from fastapi import FastAPI\nfrom pydantic import BaseModel\n\napp = FastAPI()\n\nclass Msg(BaseModel):\n    text: str\n\n# POST /echo taking a Msg -> {'echo': <text>, 'length': <len>}\n",
		Checks:      `[{"text":"It echoes and measures the text.","test":"client.post('/echo', json={'text': 'hi'}).json() == {'echo': 'hi', 'length': 2}"},{"text":"It works for longer text.","test":"client.post('/echo', json={'text': 'hello'}).json()['length'] == 5"}]`,
	},
}

// --- Testing labs (Lang "pyserver"): the learner writes assert-based tests
// (the shape pytest discovers and runs) or fixes buggy code to make a provided
// test pass. A passing test function returns None; a failing assert raises, so
// `test_x() is None` is True exactly when the learner's test is green. `_code`
// (the learner's source) lets a check confirm an `assert` was really written. ---

// writeTestsSteps — writing your first assert-based tests.
var writeTestsSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Good code comes with tests. The `total()` function is written for you — now write `test_total()` that **asserts** `total([1, 2, 3])` is `6` and `total([])` is `0`. (A test function whose asserts all hold is a green test.)",
		Starter:     "def total(prices):\n    return sum(prices)\n\ndef test_total():\n    # assert total([1, 2, 3]) == 6, and total([]) == 0\n    pass\n",
		Checks:      `[{"text":"test_total is a function.","test":"callable(test_total)"},{"text":"Your test uses assert.","test":"'assert' in _code"},{"text":"It actually calls total().","test":"'total(' in _code"},{"text":"Your test passes (no assertion fails).","test":"test_total() is None"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "Good tests cover the tricky cases. Add a second test, `test_total_edge()`, that asserts `total([5])` is `5` and `total([-1, 1])` is `0`.",
		Starter:     "def total(prices):\n    return sum(prices)\n\ndef test_total():\n    assert total([1, 2, 3]) == 6\n    assert total([]) == 0\n\ndef test_total_edge():\n    # assert total([5]) == 5, and total([-1, 1]) == 0\n    pass\n",
		Checks:      `[{"text":"test_total_edge is a function.","test":"callable(test_total_edge)"},{"text":"You now have two test functions.","test":"_code.count('def test_') >= 2"},{"text":"Both tests pass.","test":"test_total() is None and test_total_edge() is None"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "A test's job is to **catch mistakes**. Write `test_total_mixed()` that asserts `total([10, -3, 3])` is `10` — a mix of positive and negative numbers.",
		Starter:     "def total(prices):\n    return sum(prices)\n\ndef test_total_mixed():\n    # assert total([10, -3, 3]) == 10\n    pass\n",
		Checks:      `[{"text":"test_total_mixed is a function.","test":"callable(test_total_mixed)"},{"text":"It uses assert on total().","test":"'assert' in _code and 'total(' in _code"},{"text":"It passes.","test":"test_total_mixed() is None"}]`,
	},
}

// findTheBugSteps — read a failing test, fix the code to turn it green.
var findTheBugSteps = []store.Step{
	{
		Lang:        "pyserver",
		Instruction: "Here's a function with a bug and a test that catches it. `average()` crashes on an empty list. **Fix `average`** so the provided `test_average` passes — an empty list should return `0`. (Don't change the test.)",
		Starter:     "def average(nums):\n    return sum(nums) / len(nums)   # bug: crashes on an empty list\n\ndef test_average():\n    assert average([2, 4]) == 3\n    assert average([]) == 0\n",
		Checks:      `[{"text":"The provided test now passes.","test":"test_average() is None"},{"text":"average([2, 4]) is still 3.","test":"average([2, 4]) == 3"},{"text":"average([]) returns 0 instead of crashing.","test":"average([]) == 0"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "`biggest()` returns the **first** number, not the biggest. Fix it so `test_biggest` passes.",
		Starter:     "def biggest(nums):\n    return nums[0]   # bug: returns the first, not the biggest\n\ndef test_biggest():\n    assert biggest([3, 9, 4]) == 9\n    assert biggest([1]) == 1\n",
		Checks:      `[{"text":"test_biggest passes.","test":"test_biggest() is None"},{"text":"It finds the max of several numbers.","test":"biggest([3, 9, 4]) == 9"},{"text":"A single number is its own max.","test":"biggest([1]) == 1"}]`,
	},
	{
		Lang:        "pyserver",
		Instruction: "`is_even()` has its logic backwards. Fix it so `test_is_even` passes.",
		Starter:     "def is_even(n):\n    return n % 2 == 1   # bug: this is the test for ODD\n\ndef test_is_even():\n    assert is_even(4) is True\n    assert is_even(3) is False\n    assert is_even(0) is True\n",
		Checks:      `[{"text":"test_is_even passes.","test":"test_is_even() is None"},{"text":"4 is even.","test":"is_even(4) is True"},{"text":"3 is not even.","test":"is_even(3) is False"}]`,
	},
}

// pythonLabs are the interactive Python labs (run via Pyodide/WASM).
var pythonLabs = map[string][]store.Step{
	"workshop_python_warm_up":             pythonWarmupSteps,
	"workshop_build_an_xp_calculator":     xpCalcSteps,
	"workshop_build_a_loot_inventory":     lootInventorySteps,
	"workshop_score_multipliers":          comboSteps,
	"workshop_model_a_player":             playerClassSteps,
	"workshop_build_a_chat_filter":        chatFilterSteps,
	"workshop_tally_a_scoreboard":         tallySteps,
	"workshop_render_a_health_bar":        healthBarSteps,
	"workshop_combat_math":                combatMathSteps,
	"workshop_safe_stat_parser":           safeParserSteps,
	"workshop_build_a_number_stream":      numberStreamSteps,
	"workshop_write_a_decorator":          decoratorSteps,
	"workshop_recursion_puzzles":          recursionSteps,
	"workshop_parse_game_logs_with_regex": regexSteps,
	"workshop_model_with_dataclasses":     dataclassSteps,
	"workshop_your_first_flask_route":     flaskRouteSteps,
	"workshop_handle_request_data":        flaskRequestSteps,
	"workshop_build_a_rest_api":           restApiSteps,
	"workshop_query_a_database":           sqliteSteps,
	"workshop_a_fastapi_endpoint":         fastapiSteps,
	"workshop_write_your_first_tests":     writeTestsSteps,
	"workshop_find_the_bug":               findTheBugSteps,
}

// labSteps maps a frontend lesson slug to its interactive steps.
var labSteps = map[string][]store.Step{
	"workshop_code_a_level_up_system":           levelUpSteps,
	"workshop_build_an_inventory":               inventorySteps,
	"workshop_filter_a_product_catalog":         productCatalogSteps,
	"workshop_build_a_like_button":              likeButtonSteps,
	"workshop_build_a_password_strength_meter":  passwordStrengthSteps,
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
