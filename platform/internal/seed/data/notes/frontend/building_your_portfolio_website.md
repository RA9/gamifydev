# Building your portfolio website

A portfolio is the single most useful thing you can build as a new developer. It proves you can ship real work, and it's the link you'll paste into every application and message. In this lesson you'll plan, build, and prepare to deploy your own.

## Why a portfolio matters

Anyone can *say* they know HTML and CSS. A portfolio *shows* it. A good one does three jobs at once:

- It demonstrates your skills by being, itself, a thing you built.
- It collects your projects in one place with context and links.
- It gives people an easy way to learn about you and reach you.

:::key
Your portfolio is both your résumé and a sample of your work at the same time. The site itself is evidence. Sloppy code or a broken link undercuts everything you say — polish counts here more than anywhere.
:::

## What to include

Keep it focused. A strong beginner portfolio has these sections:

- **Hero / intro** — your name, a one-line description of what you do, and a clear call to action (view work, contact me).
- **About** — a short, human paragraph: who you are, what you're into, what you're learning.
- **Projects** — the centerpiece. Each project gets a screenshot, a short description, the tech used, and links to the live site and source code.
- **Skills** — the languages and tools you're comfortable with. Keep it honest.
- **Contact** — an email link and your GitHub/LinkedIn. Make it effortless to reach you.

:::tip
Three solid, finished projects beat ten half-built ones. Quality and a clear write-up matter more than quantity. Lead with your best.
:::

## Planning the sections

Before writing any code, sketch the page top to bottom on paper or in notes. A simple, proven order:

1. Sticky nav with links that jump to each section.
2. Hero.
3. About.
4. Projects.
5. Skills.
6. Contact + footer.

Deciding the structure first means your HTML writes itself.

## Structuring it with semantic HTML

Use semantic tags so the page is meaningful and accessible (recall the HTML lessons). Reach for `<header>`, `<nav>`, `<main>`, `<section>`, `<article>`, and `<footer>` instead of a pile of `<div>`s.

```html
<header>
  <nav>
    <a href="#projects">Work</a>
    <a href="#about">About</a>
    <a href="#contact">Contact</a>
  </nav>
</header>

<main>
  <section id="hero">
    <h1>Jane Doe</h1>
    <p>Frontend developer building clean, accessible websites.</p>
    <a class="btn" href="#projects">See my work</a>
  </section>
  <!-- about, projects, skills, contact below -->
</main>

<footer>
  <p>&copy; 2026 Jane Doe</p>
</footer>
```

## A starter projects section

This is the part recruiters actually read, so make it shine. Here's a skeleton using a responsive card grid (recall CSS Grid).

```html
<section id="projects">
  <h2>Projects</h2>
  <div class="project-grid">

    <article class="project-card">
      <img src="images/quiz-game.png" alt="Screenshot of the quiz game app" />
      <h3>Quiz Game</h3>
      <p>A timed multiple-choice quiz built with vanilla JavaScript and the DOM.</p>
      <ul class="tags"><li>HTML</li><li>CSS</li><li>JavaScript</li></ul>
      <div class="links">
        <a href="https://yourname.github.io/quiz-game/">Live</a>
        <a href="https://github.com/yourname/quiz-game">Code</a>
      </div>
    </article>

    <article class="project-card">
      <img src="images/profile-card.png" alt="Screenshot of the profile card component" />
      <h3>Profile Card</h3>
      <p>A responsive profile card practicing flexbox, spacing, and hover states.</p>
      <ul class="tags"><li>HTML</li><li>CSS</li></ul>
      <div class="links">
        <a href="https://yourname.github.io/profile-card/">Live</a>
        <a href="https://github.com/yourname/profile-card">Code</a>
      </div>
    </article>

  </div>
</section>
```

```css
.project-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.project-card {
  border: 1px solid #e2e2e2;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.project-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.project-card img { width: 100%; display: block; }
.project-card h3, .project-card p { padding: 0 1rem; }

.tags { display: flex; gap: 0.5rem; list-style: none; padding: 0 1rem; }
.tags li { font-size: 0.8rem; background: #f0f0f0; padding: 2px 8px; border-radius: 999px; }

.links { display: flex; gap: 1rem; padding: 1rem; }
```

That `repeat(auto-fit, minmax(280px, 1fr))` line is the magic: cards flow into as many columns as fit and stack to one column on phones — responsive with zero media queries.

:::project
**Build it.** Create a `portfolio/` folder with `index.html`, `css/style.css`, and an `images/` folder. Build all six sections. Add at least three real projects you've made in this course (the profile card and quiz game are great starts). Take a clean screenshot of each, drop them in `images/`, and wire up the live + code links.
:::

## Writing good project case studies

A screenshot and a link aren't enough — say something about each project. A tiny case study answers:

- **What is it?** One sentence.
- **What did you build / learn?** "Practiced the DOM by handling clicks and updating the score live."
- **What was tricky?** A sentence on a challenge you solved shows real thinking.

```html
<p>
  A timed quiz that tracks your score as you answer. The tricky part was
  resetting the timer between questions without leaking old intervals —
  I solved it with clearInterval before each new round.
</p>
```

That extra sentence is what makes you sound like a developer, not a copier of tutorials.

## Performance and polish

Small touches make a portfolio feel professional:

- **Optimize images.** Resize screenshots to the size they display and export as compressed JPG/WebP. Don't ship a 4 MB PNG.
- **Add `alt` text** to every image for accessibility.
- **Consistent spacing.** Pick a spacing scale (8px, 16px, 24px) and stick to it.
- **Readable type.** Comfortable line length, decent contrast, sane font sizes.
- **Hover and focus states.** Links and buttons should react. Don't forget `:focus` for keyboard users.
- **A real `<title>` and favicon.** It's the first thing seen in the browser tab.

:::warning
Test on a real phone, not just the desktop browser. Tiny tap targets, text running off the edge, and images that overflow are the most common portfolio problems — and all of them only show up at small widths.
:::

## Then deploy it

A portfolio nobody can visit isn't doing its job. Once it looks good, push it to GitHub and deploy with GitHub Pages or Netlify (recall the deploy lesson).

```bash
git init
git add .
git commit -m "Initial portfolio site"
git remote add origin https://github.com/yourname/portfolio.git
git push -u origin main
# then enable GitHub Pages in Settings -> Pages
```

Run the pre-deploy checklist — relative paths, working links, clean console, mobile check — and share your URL.

:::quiz
Q: What is the most important section of a developer portfolio?
- A long autobiography
- The projects section, with screenshots, descriptions, and live + code links *
- A list of every tutorial you've watched
E: Projects are the evidence that you can build real things. Screenshots, short write-ups, and working links make that evidence easy to evaluate.
:::

:::quiz
Q: Why use semantic tags like <section> and <article> instead of plain <div>s?
- They render faster than divs
- They make the page more meaningful and accessible, and the code easier to read *
- They are required by the browser
E: Semantic elements describe what content *is*, which helps screen readers, search engines, and other developers (including future you) understand the page.
:::

:::fill
The CSS that creates a responsive card grid with no media queries is `grid-template-columns: repeat(auto-fit, minmax(280px, ______));`.
- 1fr *
- 100px
- auto
E: `minmax(280px, 1fr)` lets each card be at least 280px and grow to fill an equal share, while `auto-fit` flows them into as many columns as fit.
:::

## Recap

- A portfolio proves your skills and centralizes your work — it's the link you'll share everywhere.
- Include a hero, about, projects, skills, and contact section; plan the order before coding.
- Build with semantic HTML and a responsive CSS Grid; cards with `auto-fit` + `minmax` go responsive without media queries.
- Write a short case study per project: what it is, what you built/learned, what was tricky.
- Polish with optimized images, `alt` text, consistent spacing, hover/focus states, a title, and a favicon.
- Test on a real phone, then deploy with GitHub Pages or Netlify and run the pre-deploy checklist.

**Next up:** Next Steps — where to take your skills from here.
