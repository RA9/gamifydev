# Performance SEO and Launch

Your portfolio is built, the content is strong, and the projects are linked. Before you launch, optimize for performance, make it findable by search engines, and prepare it for social sharing. These finishing touches separate a student project from a professional site.

## Optimize images

Images are the biggest performance bottleneck on most portfolios. A single uncompressed screenshot can be larger than all your HTML, CSS, and JavaScript combined.

### Step 1: Resize

Your project screenshots don't need to be 3000px wide. Resize to the maximum display width — usually 600–1200px.

### Step 2: Compress

Run every image through a compression tool:
- **Squoosh** (squoosh.app) — free, runs in browser, great controls
- **TinyPNG** — batch compression, drag-and-drop
- **ImageOptim** (Mac) — desktop app, automatic compression

### Step 3: Use WebP format

WebP is 25–35% smaller than JPEG at the same quality. Use it as the primary format with a JPEG fallback:

```html
<picture>
  <source srcset="images/quiz-game.webp" type="image/webp" />
  <img
    src="images/quiz-game.jpg"
    alt="Quiz game showing a question with four answer options"
    width="600"
    height="400"
    loading="lazy"
  />
</picture>
```

### Step 4: Set width and height

Always include `width` and `height` attributes on `<img>` tags. This prevents **layout shift** (CLS) — the browser reserves space before the image loads, so content doesn't jump around.

```html
<!-- Good: dimensions set, lazy loading enabled -->
<img src="quiz.jpg" alt="Quiz game" width="600" height="400" loading="lazy" />

<!-- Bad: no dimensions, causes layout shift -->
<img src="quiz.jpg" alt="Quiz game" />
```

### Target file sizes

| Image type | Target size |
|-----------|-------------|
| Hero / about photo | < 150 KB |
| Project screenshots | < 100 KB each |
| Icons / logos | < 20 KB |

:::warning
A portfolio with five 2 MB screenshots downloads 10 MB of images. On a mobile connection, that's 10+ seconds of loading. Compress everything — there's no excuse for shipping unoptimized images.
:::

## Add the `<meta name="description">` tag

The meta description appears in search engine results below your page title. It's your elevator pitch to anyone who finds you through Google.

```html
<head>
  <meta
    name="description"
    content="Jane Doe — Frontend developer portfolio. Projects built with HTML, CSS, and JavaScript including responsive landing pages, interactive quiz games, and more."
  />
</head>
```

**Rules:**
- 150–160 characters maximum
- Include your name and role
- Mention key technologies
- Make it read like a sentence, not keyword stuffing

## Open Graph tags for social sharing

When someone shares your portfolio link on LinkedIn, Twitter, or Slack, **Open Graph tags** control what the preview card looks like.

```html
<head>
  <!-- Primary Meta Tags -->
  <title>Jane Doe — Frontend Developer</title>
  <meta name="description" content="Frontend developer portfolio featuring responsive landing pages, interactive JavaScript projects, and more." />

  <!-- Open Graph / Social -->
  <meta property="og:type" content="website" />
  <meta property="og:title" content="Jane Doe — Frontend Developer" />
  <meta property="og:description" content="Portfolio featuring responsive, accessible web projects built with HTML, CSS, and JavaScript." />
  <meta property="og:image" content="https://janedoe.dev/images/og-preview.png" />
  <meta property="og:url" content="https://janedoe.dev" />

  <!-- Twitter -->
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="Jane Doe — Frontend Developer" />
  <meta name="twitter:description" content="Portfolio featuring responsive, accessible web projects." />
  <meta name="twitter:image" content="https://janedoe.dev/images/og-preview.png" />
</head>
```

### Creating the OG image

The `og:image` should be a 1200×630px image that represents your portfolio. Options:

- Screenshot of your hero section
- A designed card with your name, title, and a visual
- Use a tool like **Figma** or **Canva** to create one

:::key
Open Graph tags take 5 minutes to add and dramatically improve how your portfolio looks when shared. A link with a rich preview card gets significantly more clicks than a plain URL. This is especially important when sharing on LinkedIn.
:::

## The title tag

The `<title>` tag appears in browser tabs and search results. Make it clear and specific:

```html
<!-- Good -->
<title>Jane Doe — Frontend Developer Portfolio</title>

<!-- Bad -->
<title>My Website</title>
<title>Portfolio</title>
<title>index.html</title>
```

## Favicon

A missing favicon makes the browser tab look unfinished. Add one:

```html
<link rel="icon" type="image/png" sizes="32x32" href="favicon-32x32.png" />
<link rel="icon" type="image/png" sizes="16x16" href="favicon-16x16.png" />
<link rel="apple-touch-icon" sizes="180x180" href="apple-touch-icon.png" />
```

Use **favicon.io** to generate all sizes from a letter, emoji, or image.

## Lighthouse audit — aim for 90+

Run Lighthouse in Chrome DevTools (Lighthouse tab → Mobile → Analyze) and target:

| Category | Target | Common fixes |
|----------|--------|-------------|
| Performance | 90+ | Compress images, add dimensions, lazy-load below-fold images |
| Accessibility | 95+ | Alt text, heading order, contrast, focus states |
| Best Practices | 95+ | HTTPS, no console errors, correct image aspect ratios |
| SEO | 90+ | Meta description, title tag, mobile viewport, crawlable links |

Address the specific recommendations Lighthouse gives you. Each one links to documentation explaining why it matters and how to fix it.

:::tip
Run Lighthouse on the **deployed** URL, not `localhost`. Local results can be misleading because there's no network latency and the browser has cached everything.
:::

## Checking mobile rendering

Open DevTools → device toolbar and check these widths:

- **320px** — smallest phones (iPhone SE)
- **375px** — standard phones (iPhone 12–15)
- **768px** — tablets
- **1024px** — small laptops

At each width, verify:

- Text is readable without zooming
- Project cards stack properly
- Navigation is usable (hamburger menu works if applicable)
- No horizontal scrollbar
- Touch targets are at least 44px
- Images don't overflow their containers

Then test on your **actual phone**. Real devices reveal font rendering, scroll behavior, and touch target issues that emulation misses.

## Final accessibility pass

Before launch, run through this checklist one more time:

```text
✅ Every image has descriptive alt text
✅ Headings follow h1 → h2 → h3 order (no skipping)
✅ All text meets 4.5:1 contrast ratio
✅ Tab through the entire page — every link and button is reachable
✅ Focus indicators are visible
✅ Links have descriptive text (not "click here")
✅ The page has a lang attribute: <html lang="en">
✅ Interactive elements have hover AND focus states
```

Install the **axe DevTools** browser extension for an automated accessibility audit that catches issues Lighthouse might miss.

## Deploying

Choose your platform and ship:

**GitHub Pages:**
```bash
git add .
git commit -m "Final portfolio with optimized images and meta tags"
git push
# Settings → Pages → Deploy from main branch
```

**Netlify:**
- Connect your repo → auto-deploys on every push
- Or drag-and-drop your folder for instant deployment

**Vercel:**
- Import repo → auto-deploys with preview URLs on branches

All three give you HTTPS and a free subdomain. Add a custom domain if you have one.

## Sharing your portfolio

Your portfolio is live — now make sure people see it.

### LinkedIn

Write a short post announcing your portfolio:

```text
I just launched my frontend developer portfolio! 🚀

Built with HTML, CSS, and vanilla JavaScript, featuring:
→ A responsive landing page
→ An interactive quiz game
→ A task tracker with localStorage

Check it out: [link]

I'd love any feedback from the dev community.

#frontend #webdevelopment #portfolio
```

### GitHub

- **Pin** your portfolio repo on your GitHub profile
- **Add a description** and **website URL** to the repo
- Write a clear **README** with a screenshot and live link

### Other places to share

- **Twitter/X** — post with a screenshot and link
- **Dev.to** — write a short "I built my portfolio" post
- **Reddit** — r/webdev and r/learnprogramming have portfolio review threads

## The maintenance mindset

A portfolio is not "done forever." Treat it as a living document:

- **Update projects** as you build stronger ones
- **Replace screenshots** when you improve a project's design
- **Check links** periodically — broken links are embarrassing
- **Update your positioning statement** as your skills grow
- **Remove projects** that no longer represent your best work

:::warning
A portfolio with broken links, outdated screenshots, or placeholder text actively hurts your credibility. If you can't maintain it, keep it minimal — fewer projects that work perfectly are better than many that are half-broken.
:::

Set a quarterly reminder to review your portfolio for 15 minutes. Update links, swap in better projects, refresh the copy.

:::quiz
Q: Why are Open Graph tags important for a portfolio?
- They improve page load speed
- They control how your portfolio looks when shared on LinkedIn, Twitter, and other platforms *
- They are required for the page to render in browsers
- They encrypt your personal information
E: Open Graph tags define the title, description, and image that appear in link preview cards on social platforms. A rich preview with a professional image gets significantly more clicks than a bare URL.
:::

:::quiz
Q: What is the most impactful performance optimization for most portfolio sites?
- Minifying HTML tags
- Compressing and resizing images *
- Using a faster JavaScript framework
- Adding more CSS animations
E: Images are typically the largest assets on a portfolio page. A single uncompressed 3 MB screenshot dwarfs all your HTML, CSS, and JS combined. Compressing, resizing, and using modern formats like WebP can reduce total page weight by 80% or more.
:::

## Launch checklist

- [ ] Images compressed and sized (< 100-150 KB each)
- [ ] Width and height set on all `<img>` tags
- [ ] `<title>` tag is descriptive
- [ ] `<meta name="description">` is set (150 chars)
- [ ] Open Graph tags added for social previews
- [ ] Favicon added
- [ ] Lighthouse mobile scores: Performance 90+, Accessibility 95+
- [ ] Tested on a real phone
- [ ] All project links work (live demo + source)
- [ ] Accessibility pass complete (alt text, headings, contrast, keyboard)
- [ ] Deployed to a public URL
- [ ] Shared on LinkedIn and GitHub

## Recap

- **Optimize images:** resize, compress, use WebP, set dimensions, lazy-load.
- **Meta description** and **title tag** make your site findable in search results.
- **Open Graph tags** control social sharing previews — add `og:title`, `og:description`, `og:image`.
- **Favicon** finishes the browser tab.
- **Lighthouse audit** on mobile — aim for 90+ across all categories.
- **Test on real devices** — emulation misses font rendering, touch, and scroll issues.
- **Final accessibility pass** — alt text, heading order, contrast, keyboard navigation.
- **Deploy** to GitHub Pages, Netlify, or Vercel.
- **Share** on LinkedIn, GitHub, and developer communities.
- **Maintain** your portfolio — update projects, check links, refresh content quarterly.

**Your portfolio is live.** You've planned it, built it, written strong case studies, optimized it, and shipped it. That's the full professional workflow — and the portfolio itself is proof you can do it.
