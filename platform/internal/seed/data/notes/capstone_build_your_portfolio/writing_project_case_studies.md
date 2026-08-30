# Writing Project Case Studies

The projects section is the most important part of your portfolio. But a screenshot and a title aren't enough — you need to tell the story of each project in a way that shows how you think, not just what you built.

## What a case study is (and isn't)

A portfolio case study is **not** a 2,000-word blog post. It's **3–5 focused sentences** that answer the questions a hiring manager actually has.

```text
❌ Too short:
   Quiz Game
   [screenshot]

❌ Too long:
   [3 paragraphs about the history of quiz games,
    your personal learning journey, every tutorial
    you watched, and a philosophical reflection on
    the nature of interactive learning...]

✅ Just right:
   Quiz Game
   A timed quiz app with keyboard shortcuts, visual feedback, and score tracking.
   Built with HTML, CSS, and vanilla JavaScript.
   Used a state → render loop to manage question flow and score state.
   Hardest part: handling the timer countdown and auto-advance without race conditions.
   [Live Demo] [Source Code]
```

:::key
A case study shows engineering thinking. Listing features proves you built something. Explaining *how* you solved a hard problem proves you can think through challenges — which is what employers are actually hiring for.
:::

## The five questions to answer

For each project, write answers to these five questions:

### 1. What is it?

One sentence describing what the project does from a user's perspective.

```text
A responsive landing page for a fictional coffee subscription service.
A quiz game that tests web development knowledge with timed questions.
A task tracker that persists data across browser sessions using localStorage.
```

### 2. What did you build?

One sentence about what you specifically created — not what the tutorial covered, but what *you* did.

```text
Built a mobile-first layout with CSS Grid, a hamburger menu, and smooth scroll navigation.
Created a state-driven question engine that renders quiz content from a data array.
Implemented CRUD operations with DOM rendering and localStorage for data persistence.
```

### 3. What tech did you use?

Keep it to the actual technologies. Be specific but honest.

```text
HTML, CSS (Grid, Flexbox, custom properties), vanilla JavaScript
HTML, CSS, JavaScript, IntersectionObserver API
HTML, CSS, JavaScript, localStorage, Fetch API
```

:::warning
Only list technologies you can discuss in an interview. If someone asks "Why did you use CSS Grid instead of Flexbox for the feature section?", you should be able to answer. Listing React when you followed one tutorial will backfire.
:::

### 4. What was the hardest part?

This is the most valuable sentence in your case study. It shows problem-solving ability.

```text
Hardest part: getting the responsive nav to work without JavaScript
using the checkbox hack, while keeping it accessible with proper
ARIA attributes and keyboard support.

Hardest part: managing the quiz timer — clearInterval had to fire
in the right order to prevent the timer from continuing after an
answer was selected.

Hardest part: keeping the UI and localStorage in sync — I learned
to always re-render from the stored data rather than trying to
update both independently.
```

:::tip
If you can't identify a hard part, you either didn't build it yourself or you've forgotten the struggle. Revisit your git history — the commits where you fixed bugs are usually the hardest parts.
:::

### 5. What did you learn?

One sentence about a skill or concept you gained from this project.

```text
Learned that separating data from presentation makes code much easier
to extend — I added a question category feature in 20 minutes because
the architecture supported it.

Learned to test responsive layouts at real device widths instead of
just resizing the browser — iOS Safari behaves differently than Chrome's
device emulation.

Learned that state management is the core challenge in interactive UIs,
not DOM manipulation. This gave me a head start understanding React.
```

## Structuring a project card in HTML

```html
<article class="project-card">
  <img
    src="images/quiz-game.png"
    alt="Screenshot of the quiz game showing a question with four answer options"
    width="600"
    height="400"
    loading="lazy"
  />
  <div class="project-info">
    <h3>Quiz Game</h3>
    <p class="project-description">
      A timed quiz app with keyboard shortcuts, visual feedback,
      and score tracking. Built with a state → render architecture.
    </p>
    <p class="project-tech">HTML · CSS · JavaScript</p>
    <p class="project-insight">
      Used a state → render loop to manage question flow, which
      made adding a timer feature straightforward.
    </p>
    <div class="project-links">
      <a href="https://janedoe.github.io/quiz-game" target="_blank" rel="noopener">
        Live Demo
      </a>
      <a href="https://github.com/janedoe/quiz-game" target="_blank" rel="noopener">
        Source Code
      </a>
    </div>
  </div>
</article>
```

Every project card should have:
- **Screenshot** — real, not a placeholder
- **Title** — clear, short name
- **Description** — what it is and what it does
- **Tech** — what you used
- **Insight** — the "hard part" or "what I learned" sentence
- **Links** — live demo AND source code

## Screenshot + live link + source link

All three matter:

| Element | Why it matters |
|---------|---------------|
| **Screenshot** | First impression — visual proof the project looks professional |
| **Live link** | Proves it works — the reviewer can click and interact |
| **Source link** | Proves you wrote the code — they can check quality |

### Taking good screenshots

- Use a clean browser window (no bookmarks bar, no extensions showing)
- Capture at a standard resolution (1200×800 or similar)
- Show the most interesting state (a filled-out quiz, not the loading screen)
- Compress the image to under 200 KB (use Squoosh or TinyPNG)

:::tip
For quiz games and interactive projects, capture a screenshot of the "in-progress" state — a question with one option selected and feedback visible. This is more interesting than the start screen and shows the interaction design.
:::

## 3–5 sentences, not 3 paragraphs

The biggest mistake in portfolio case studies is writing too much. Reviewers scan — they don't read essays.

**The rule:** each project gets 3–5 sentences maximum in the card view. If you want to write more, create a separate project detail page.

```text
✅ Card-length case study:
   A responsive landing page for a fictional coffee subscription service.
   Built mobile-first with CSS Grid and custom properties for consistent spacing.
   The hardest part was getting the hamburger nav to work accessibly without JavaScript.

❌ Too long for a card:
   This project started when I was learning about responsive design...
   [4 more paragraphs]
```

## Show engineering thinking, not just features

There's a difference between describing features and demonstrating thinking:

```text
Features (what everyone writes):
  "Has a timer, score tracking, and keyboard shortcuts."

Thinking (what stands out):
  "I used a state → render loop so adding features like
   the timer only required new state variables and a render
   update — the architecture scaled without refactoring."
```

The second version tells a reviewer: "This person understands software design, not just DOM methods."

**More examples of engineering thinking:**

```text
"I chose CSS Grid over Flexbox for the feature section because the
2D layout needed both row and column control."

"I used IntersectionObserver instead of scroll events for reveal
animations because it's natively optimized and doesn't cause jank."

"I stored the quiz state in three variables (currentIndex, score,
answered) which fully described the UI at any point — this made
debugging trivial."
```

:::key
Features are what the project does. Engineering thinking is *why you built it the way you did.* Hiring managers see hundreds of quiz games. The one that stands out explains its architecture.
:::

## Styling the project cards

```css
.project-card {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
  background: var(--color-surface);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.25s ease, transform 0.25s ease;
}

.project-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  transform: translateY(-4px);
}

.project-card img {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.project-info {
  padding: 1.25rem;
}

.project-tech {
  font-size: 0.875rem;
  color: var(--color-muted);
  margin: 0.5rem 0;
}

.project-insight {
  font-size: 0.9rem;
  font-style: italic;
  color: var(--color-text);
  margin-bottom: 1rem;
}

.project-links a {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  text-decoration: none;
  margin-right: 0.5rem;
}
```

## Case study template

Use this template for each project:

```text
[Project Name]
[Screenshot — in-progress state]

What: [One sentence — what does it do?]
Built: [One sentence — what did you specifically build?]
Tech: [List of technologies used]
Challenge: [One sentence — hardest part or key decision]

[Live Demo]  [Source Code]
```

:::quiz
Q: Why is the "hardest part" sentence the most valuable part of a portfolio case study?
- It makes the project sound more difficult
- It demonstrates problem-solving ability and engineering thinking, which is what employers hire for *
- It fills space on the page
- It's required by HTML standards
E: Anyone can list features. Explaining a specific challenge you solved shows that you can think through problems, debug issues, and make architectural decisions — skills that matter far more than the feature list.
:::

:::quiz
Q: How long should a project case study be on the portfolio card?
- At least 500 words to be thorough
- 3–5 sentences that cover what, how, and why *
- One word: the project name
- As long as possible to show effort
E: Portfolio reviewers scan, they don't read essays. 3–5 concise sentences covering what the project is, what tech you used, and what was challenging communicate everything a reviewer needs in the time they'll actually spend looking.
:::

## Recap

- Each project case study answers five questions: What is it? What did you build? What tech? What was hard? What did you learn?
- Include a **screenshot** (in-progress state), **live link**, and **source link** for every project.
- Keep it **3–5 sentences** — short enough to scan, long enough to show thinking.
- **Engineering thinking** ("I chose X because Y") stands out more than feature lists.
- Use a consistent **project card template** for visual uniformity.
- Only list technologies **you can discuss in an interview**.

**Next up:** Performance, SEO, and Launch — optimizing your portfolio for speed, search engines, and social sharing before going live.
