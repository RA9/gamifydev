# Planning a Landing Page

Before you write a single line of code, plan the page. A landing page that feels polished always starts with clear decisions about structure, color, typography, and spacing — made *before* the editor opens. This lesson walks you through the planning process professionals use.

## Start with the product brief

Every landing page sells or explains something. Before designing, answer these questions:

- **What is the product or service?** (e.g. a coffee subscription)
- **Who is the audience?** (busy professionals who want great coffee at home)
- **What is the one action you want visitors to take?** (sign up for a trial)

That last question is the most important. Everything on the page should guide the visitor toward that single action.

:::key
A landing page has one job: get the visitor to take one specific action. Every section, headline, and button should support that goal. If a section doesn't serve the goal, cut it.
:::

## Identify the core sections

Most effective landing pages follow a predictable pattern. Visitors expect it, and it works:

1. **Hero** — headline, subheadline, primary call-to-action button
2. **Features / Benefits** — 3–4 cards explaining what makes the product valuable
3. **Social proof / Testimonials** — quotes from real users or trust signals
4. **CTA (call-to-action) block** — a final push to convert, often with strong contrast
5. **Footer** — copyright, legal links, secondary navigation

```text
┌────────────────────────────────────┐
│           HEADER / NAV             │
├────────────────────────────────────┤
│            HERO                    │
│   Headline + CTA button           │
├────────────────────────────────────┤
│          FEATURES                  │
│   Card   Card   Card              │
├────────────────────────────────────┤
│        TESTIMONIALS                │
│   "Quote" — Name                  │
├────────────────────────────────────┤
│       FINAL CTA BLOCK             │
│   Strong headline + button        │
├────────────────────────────────────┤
│           FOOTER                   │
└────────────────────────────────────┘
```

:::tip
Sketch this on paper or a whiteboard first. A five-minute wireframe saves an hour of aimless CSS later.
:::

## Sketching wireframes

A wireframe is a rough layout — no colors, no real images, just boxes and labels. You can use:

- **Paper and pen** — fastest and most flexible
- **Excalidraw** (excalidraw.com) — free, simple, shareable
- **Figma** — professional-grade, free tier available

Keep wireframes low-fidelity. The point is layout decisions, not visual design:

- How wide is the hero? Full-width or contained?
- Is the feature grid 2 columns or 3?
- Where does the testimonial sit relative to the CTA?
- Does the nav stick to the top on scroll?

Answer these with rough boxes before writing HTML.

## Choosing a color palette — the 3-color rule

Too many colors make a page feel chaotic. Start with exactly three:

| Role          | Purpose                            | Example      |
|---------------|-------------------------------------|-------------|
| **Primary**   | Brand color, buttons, accents       | `#2563eb`   |
| **Neutral**   | Text, backgrounds, borders          | `#1e293b`   |
| **Background**| Page background, card surfaces      | `#f8fafc`   |

Add one or two shades of each (lighter and darker) for hover states and subtle variations.

```css
:root {
  --color-primary: #2563eb;
  --color-primary-dark: #1d4ed8;
  --color-neutral: #1e293b;
  --color-muted: #64748b;
  --color-bg: #f8fafc;
  --color-card: #ffffff;
}
```

:::warning
Check contrast ratios before committing to your palette. Text must have at least a **4.5:1** contrast ratio against its background to meet WCAG AA. Use a tool like WebAIM's contrast checker.
:::

## Selecting fonts — heading + body pair

Pick two fonts maximum: one for headings, one for body text.

**Safe, proven combinations:**

- **Inter** (body) + **Inter** (headings, bold weight) — clean and modern
- **DM Sans** (body) + **DM Serif Display** (headings) — warm and editorial
- **System font stack** (body + headings) — zero load time, always readable

```css
body {
  font-family: 'Inter', system-ui, sans-serif;
  font-size: 1rem;    /* 16px base */
  line-height: 1.6;
}

h1, h2, h3 {
  font-family: 'DM Serif Display', serif;
  line-height: 1.2;
}
```

**Typography hierarchy** means your heading sizes follow a clear scale:

```css
h1 { font-size: 2.5rem; }   /* 40px — hero headline */
h2 { font-size: 1.75rem; }  /* 28px — section titles */
h3 { font-size: 1.25rem; }  /* 20px — card titles */
p  { font-size: 1rem; }     /* 16px — body text */
```

:::tip
If using Google Fonts, load only the weights you need (e.g. 400 and 700). Every extra weight adds load time.
:::

## Building a spacing scale

Inconsistent spacing is the number one reason amateur sites look "off." Professional designs use a **spacing scale** — a fixed set of values based on a base unit.

**The 4px / 8px system:**

```css
:root {
  --space-xs: 4px;     /* 0.25rem */
  --space-sm: 8px;     /* 0.5rem  */
  --space-md: 16px;    /* 1rem    */
  --space-lg: 32px;    /* 2rem    */
  --space-xl: 64px;    /* 4rem    */
  --space-2xl: 128px;  /* 8rem    */
}
```

Every margin, padding, and gap should come from this scale. No magic numbers.

```css
/* Good — values from the scale */
.hero { padding: var(--space-2xl) var(--space-md); }
.feature-grid { gap: var(--space-lg); }

/* Bad — random values */
.hero { padding: 73px 19px; }
.feature-grid { gap: 23px; }
```

:::key
Consistent spacing is what separates "looks professional" from "looks like a school project." Pick a scale, write it as CSS custom properties, and use nothing else.
:::

## Creating a simple style guide

Before coding, write down your decisions in one place. This is your **style guide** — a reference you check while building.

```text
PROJECT STYLE GUIDE — Brewly Landing Page
──────────────────────────────────────────
Colors:
  Primary:    #b5651d (warm brown)
  Dark:       #7c3f10 (hover states)
  Text:       #2b2118
  Muted:      #6b5d4f
  Background: #fdf8f3
  Card:       #ffffff

Fonts:
  Headings: DM Serif Display, serif
  Body:     Inter, system-ui, sans-serif

Type scale:
  h1: 2.5rem / 1.2 line-height
  h2: 1.75rem / 1.2
  h3: 1.25rem / 1.3
  body: 1rem / 1.6

Spacing scale (base 8px):
  xs: 4px | sm: 8px | md: 16px | lg: 32px | xl: 64px

Border radius: 12px (cards, buttons)

Shadows:
  Card: 0 4px 14px rgba(0, 0, 0, 0.05)
  Hover: 0 8px 24px rgba(0, 0, 0, 0.1)
```

This might feel like overkill for a small project. It isn't. Ten minutes of planning prevents an hour of guessing.

## Putting it all together

Here's the planning workflow, in order:

1. **Read the brief** — understand the product, audience, and goal
2. **List the sections** — hero, features, testimonials, CTA, footer
3. **Sketch wireframes** — rough boxes on paper or Excalidraw
4. **Choose 3 colors** — primary, neutral, background
5. **Pick 2 fonts** — heading + body pair
6. **Define a spacing scale** — based on 4px or 8px multiples
7. **Write a style guide** — document everything in one place
8. **Then open the editor**

:::quiz
Q: Why should you sketch a wireframe before writing any code?
- It makes the CSS render faster
- It helps you make layout decisions without getting distracted by code details *
- Wireframes are required by HTML standards
- It generates the HTML for you automatically
E: Wireframes separate structural decisions from implementation. You figure out what goes where before worrying about how to code it, which saves time and reduces rework.
:::

:::quiz
Q: What is the "3-color rule" for landing page design?
- Use at least 30 colors for visual variety
- Limit your palette to three core roles: primary, neutral, and background *
- Only use red, green, and blue
- Apply three different colors to every heading
E: A tight three-color palette keeps the page visually cohesive. You can add shades and tints of those three, but the base set stays small and intentional.
:::

## Recap

- **Plan before you code.** Read the brief, identify the audience, define the one action you want.
- **Sections follow a pattern:** hero → features → social proof → CTA → footer.
- **Wireframe first** with pen and paper or a simple tool like Excalidraw.
- **3 colors** (primary, neutral, background) keep the palette clean.
- **2 fonts** (heading + body) keep typography readable and consistent.
- **Spacing scale** based on 4px/8px multiples eliminates guesswork.
- **Write a style guide** that documents colors, fonts, sizes, and spacing before opening the editor.

**Next up:** Building the HTML Structure — turning this plan into semantic, well-organized markup.
