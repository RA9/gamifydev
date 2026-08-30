# Portfolio Strategy and Planning

Your portfolio is not a gallery of everything you've built. It's a strategic tool designed to convince a specific audience — in 30 seconds — that you're worth talking to. This lesson covers how to plan that tool before writing any code.

## Who your portfolio is for

Your primary audience is a **hiring manager or recruiter** scanning dozens of portfolios in an afternoon. Here's what their process looks like:

1. **Open the link** (2 seconds) — does the page load fast? Does it look professional?
2. **Scan the hero** (5 seconds) — who is this person? What do they do?
3. **Check 1–2 projects** (15 seconds) — do they have real, working projects? Can I click a live demo?
4. **Decide** (5 seconds) — is this person worth interviewing?

Total: **30 seconds.** If your portfolio doesn't communicate clearly in that window, it fails — no matter how good your code is underneath.

:::key
You are not designing your portfolio for yourself. You are designing it for someone who has 30 seconds, is slightly impatient, and is comparing you to other candidates. Every decision — layout, copy, project selection — should serve that person.
:::

## Choosing 3–5 best projects

More projects is not better. Five polished projects outperform fifteen half-finished ones.

**Selection criteria:**

- **Is it deployed and working?** Dead links disqualify instantly.
- **Does it demonstrate a real skill?** Layout, interactivity, state management, API integration.
- **Is it visually presentable?** Can you take a screenshot that looks good?
- **Can you explain it?** If someone asks "what was hard about this?", can you answer?

**Aim for variety** — show range, not repetition:

```text
✅ Strong mix:
  1. Responsive landing page      → layout, design, responsive CSS
  2. Quiz game                    → DOM manipulation, state management
  3. Weather app                  → API integration, async/await
  4. Task tracker                 → CRUD, localStorage persistence
  5. Multi-step form              → validation, UX flow

❌ Weak mix:
  1. Landing page
  2. Another landing page
  3. A third landing page
  4. A landing page template
  5. Yet another landing page
```

:::tip
If you only have 2–3 strong projects, show 2–3. A portfolio with two excellent projects is better than one with five mediocre ones. Quality signals credibility.
:::

## Your positioning statement

The hero section needs a short statement that answers: "Who are you and what do you do?"

**Format:** [Role] + [what you build] + [what makes you interesting]

**Strong examples:**

```text
Frontend developer building accessible, responsive
interfaces with HTML, CSS, and JavaScript.

Junior web developer focused on clean code,
thoughtful UX, and shipping finished projects.

Self-taught frontend engineer with a background in
design — I care about how things look AND how they work.
```

**Weak examples:**

```text
Welcome to my portfolio!              ← says nothing
I'm passionate about coding.          ← generic
Full-stack AI blockchain developer.   ← probably not true
```

Keep it one to two sentences. No buzzwords. Be honest about your level — hiring managers appreciate self-awareness over exaggeration.

## Deciding the sections

Your portfolio needs these sections, in this order:

1. **Hero** — name, positioning statement, CTA (view projects or contact)
2. **About** — 3–4 sentences about who you are, what you enjoy building, what you're learning
3. **Projects** — the heart of the site (3–5 entries with screenshots, descriptions, and links)
4. **Skills** — honest list of technologies you can discuss in an interview
5. **Contact** — email, GitHub, LinkedIn, optional contact form

```html
<header><!-- nav --></header>
<main>
  <section id="hero"><!-- name, statement, CTA --></section>
  <section id="about"><!-- short bio --></section>
  <section id="projects"><!-- project cards --></section>
  <section id="skills"><!-- tech list --></section>
  <section id="contact"><!-- email, links --></section>
</main>
<footer><!-- copyright --></footer>
```

Each section earns its place. If it doesn't help the visitor make a decision, remove it.

## Wireframing your portfolio

Before designing in the browser, sketch the layout. Focus on:

- **Hero:** full-width? Split with an image? Just text?
- **Projects grid:** 2 columns? 3 columns? Cards or full-width rows?
- **About:** text-only? Text + photo?
- **Contact:** inline form? Just links?

```text
┌──────────────────────────────────────┐
│  Logo                  Nav links     │
├──────────────────────────────────────┤
│                                      │
│   Hi, I'm Jane Doe                   │
│   Frontend developer building...     │
│   [View My Work]                     │
│                                      │
├──────────────────────────────────────┤
│  About Me                            │
│  Short paragraph + photo             │
├──────────────────────────────────────┤
│  Projects                            │
│  ┌──────────┐  ┌──────────┐         │
│  │  Project  │  │  Project  │        │
│  │  1        │  │  2        │        │
│  └──────────┘  └──────────┘         │
│  ┌──────────┐  ┌──────────┐         │
│  │  Project  │  │  Project  │        │
│  │  3        │  │  4        │        │
│  └──────────┘  └──────────┘         │
├──────────────────────────────────────┤
│  Skills: HTML CSS JS Git ...         │
├──────────────────────────────────────┤
│  Contact: email | GitHub | LinkedIn  │
├──────────────────────────────────────┤
│  Footer                             │
└──────────────────────────────────────┘
```

Use paper, Excalidraw, or Figma. Keep it rough — the goal is structure, not pixel perfection.

## Choosing a visual identity

Your portfolio should feel intentional, not random. Make three decisions:

### Colors — 2 to 3 max

```css
:root {
  --color-primary: #2563eb;   /* accent — links, buttons, highlights */
  --color-text: #1e293b;      /* headings and body text */
  --color-muted: #64748b;     /* secondary text */
  --color-bg: #ffffff;        /* page background */
  --color-surface: #f8fafc;   /* card backgrounds, section alternation */
}
```

:::warning
Avoid trendy color schemes that sacrifice readability. A dark-mode portfolio with low-contrast gray text might look "cool" but fails the 30-second scan if the text is hard to read. Contrast first, style second.
:::

### Fonts — 1 to 2

A single font family at different weights is perfectly professional:

```css
body {
  font-family: 'Inter', system-ui, sans-serif;
}
```

Or pair a serif heading font with a sans-serif body:

```css
h1, h2, h3 { font-family: 'DM Serif Display', serif; }
body { font-family: 'Inter', system-ui, sans-serif; }
```

### Spacing — consistent scale

```css
:root {
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 2rem;
  --space-xl: 4rem;
}
```

## The content-first approach

Write the content before designing the layout. Open a text file and write:

- Your positioning statement
- Your about paragraph
- Each project's name, description, tech, and one "what was hard" sentence
- Your skills list
- Your contact info

```text
HERO:
Hi, I'm Jane Doe.
Frontend developer building accessible, responsive interfaces.

ABOUT:
I'm a self-taught frontend developer with a background in
graphic design. I love building things that look good and
work for everyone. Currently learning React and TypeScript.

PROJECT 1: Quiz Game
A timed quiz game with keyboard shortcuts and score tracking.
Built with HTML, CSS, vanilla JavaScript.
Hardest part: managing state transitions between questions.
Live: [url] | Code: [url]

...
```

:::key
Content-first means you write every word before opening your code editor. This prevents the most common portfolio mistake: spending hours on CSS while the project descriptions say "Lorem ipsum." The content IS the portfolio — the design is just the frame.
:::

Once you have the content, the design decisions become obvious: how much space does the about section need? How many cards fit in a row? Should the hero have a background image or just text?

:::quiz
Q: Why should you limit your portfolio to 3–5 projects instead of showing everything?
- Browsers can only load 5 images
- Fewer, polished projects are more convincing than many mediocre ones, and reviewers only spend 30 seconds scanning *
- GitHub only allows 5 repositories
- CSS Grid only supports 5 columns
E: Hiring managers scan quickly. Five well-presented projects with working demos, clear descriptions, and live links tell a stronger story than fifteen entries where half have broken links and no explanations.
:::

:::quiz
Q: What does "content-first approach" mean for portfolio planning?
- Design the visual layout before writing any text
- Write all the text content (bio, project descriptions, skills) before designing the layout *
- Focus on CSS animations first
- Let the browser auto-generate the content
E: Content-first means deciding what to say before deciding how to display it. This ensures your design serves the content (which is what visitors actually read) rather than forcing content into a premade design that doesn't fit.
:::

## Recap

- Your audience is a **hiring manager with 30 seconds** — design for fast scanning.
- Choose **3–5 best projects** that show variety: layout, interactivity, data, persistence.
- Write a **positioning statement** — 1–2 sentences about who you are and what you build.
- Include **five sections:** hero, about, projects, skills, contact.
- **Wireframe** the layout before coding — on paper or with a simple tool.
- Choose a **visual identity:** 2–3 colors, 1–2 fonts, consistent spacing scale.
- **Content-first:** write every word before opening the editor. The words are the portfolio; the design is the frame.

**Next up:** Building Your Portfolio Website — turning this plan into a real, coded site.
