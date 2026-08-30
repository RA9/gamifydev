# Code Review and Ship

Your landing page is built, styled, responsive, and polished. Before you share it with the world, run through a professional review process — validate the code, check accessibility, audit performance, and deploy.

## Validate your HTML

The W3C HTML validator catches structural errors that browsers silently tolerate but that hurt accessibility and SEO.

Go to **validator.w3.org**, paste your URL or upload your HTML file, and review the results.

Common issues it catches:

- Missing `alt` attributes on images
- Unclosed tags
- Duplicate `id` attributes
- Invalid nesting (e.g. `<a>` inside `<a>`)
- Missing `lang` attribute on `<html>`

```html
<!-- Before: missing alt and lang -->
<html>
  <img src="hero.jpg">

<!-- After: validated and clean -->
<html lang="en">
  <img src="hero.jpg" alt="A steaming cup of fresh coffee on a wooden table">
```

:::tip
Fix **errors** (red) first — they indicate real problems. **Warnings** (yellow) are worth reviewing but don't always need action. Aim for zero errors.
:::

## Check for accessibility

Accessibility isn't an optional nice-to-have. It's a professional standard. Run through this checklist:

### Alt text on every image

Every `<img>` needs an `alt` attribute. Describe what the image shows, not what it is.

```html
<!-- Weak -->
<img src="photo.jpg" alt="image">

<!-- Strong -->
<img src="photo.jpg" alt="Three bags of coffee beans from Brazil, Ethiopia, and Colombia">

<!-- Decorative images get empty alt -->
<img src="decoration.svg" alt="">
```

### Heading hierarchy

Headings must follow a logical order: `h1` → `h2` → `h3`. Never skip levels.

```html
<!-- Wrong: skips h2 -->
<h1>Brewly</h1>
<h3>Why people love us</h3>

<!-- Correct -->
<h1>Brewly</h1>
<h2>Why people love us</h2>
```

### Color contrast

Verify every text/background pair meets the 4.5:1 ratio (covered in the previous lesson). Use the axe DevTools extension or Lighthouse for an automated check.

### Keyboard navigation

Tab through the entire page without touching your mouse:

- Can you reach every link and button?
- Is the focus indicator visible?
- Does the tab order make sense (top to bottom, left to right)?

### Link text

Avoid generic link text that means nothing to a screen reader.

```html
<!-- Bad -->
<a href="/pricing">Click here</a>

<!-- Good -->
<a href="/pricing">View pricing plans</a>
```

:::key
Accessibility is not extra work — it's part of building a professional website. If your page isn't accessible, it isn't finished.
:::

## Run a Lighthouse audit

Lighthouse is built into Chrome DevTools. It scores your page on Performance, Accessibility, Best Practices, and SEO.

1. Open DevTools → **Lighthouse** tab
2. Select **Mobile** (more demanding than desktop)
3. Check all categories
4. Click **Analyze page load**

**Target scores:**

| Category        | Target |
|-----------------|--------|
| Performance     | 90+    |
| Accessibility   | 95+    |
| Best Practices  | 95+    |
| SEO             | 90+    |

Lighthouse gives you specific, actionable recommendations. Common wins:

- Add `width` and `height` to images (prevents layout shift)
- Add `<meta name="description">` for SEO
- Compress images (see next section)
- Ensure sufficient color contrast

## Compress images

Large images are the single biggest performance killer on landing pages. A 3 MB hero image makes the page feel slow on any connection.

**Steps:**

1. **Resize** — does your hero image need to be 4000px wide? Probably not. 1200–1600px is enough for most screens.
2. **Compress** — use tools like Squoosh (squoosh.app), TinyPNG, or ImageOptim.
3. **Use modern formats** — WebP is 25–35% smaller than JPEG at the same quality.

```html
<!-- Modern: WebP with JPEG fallback -->
<picture>
  <source srcset="hero.webp" type="image/webp" />
  <img src="hero.jpg" alt="Coffee beans" width="1200" height="800" />
</picture>
```

**Target file sizes:**
- Hero image: under 200 KB
- Feature icons: under 20 KB each
- Testimonial avatars: under 30 KB each

:::warning
Never deploy a page with uncompressed images. A single 4 MB photo can add 5+ seconds of load time on a mobile connection. Run every image through a compressor before shipping.
:::

## Deploy to GitHub Pages

If your code is already on GitHub:

1. Go to your repo → **Settings** → **Pages**
2. Under Source, choose **Deploy from a branch**
3. Select the `main` branch and `/ (root)` folder
4. Click **Save**
5. Wait 1–2 minutes, then visit `https://yourusername.github.io/repo-name/`

```bash
# Make sure everything is pushed
git add .
git commit -m "Final polish: transitions, accessibility, image compression"
git push
```

## Deploy to Netlify

**Option A — connect the repo (recommended):**

1. Go to app.netlify.com → **Add new site** → **Import an existing project**
2. Connect your GitHub account, select the repo
3. Leave build settings empty (no build step for a static site)
4. Click **Deploy**

Every future `git push` auto-deploys.

**Option B — drag and drop:**

1. Go to app.netlify.com
2. Drag your project folder onto the deploy area
3. Get a live URL in seconds

## Test on a real phone

DevTools device emulation is useful but imperfect. Real devices reveal issues emulators miss:

- Touch targets that feel too small
- Fonts that render differently
- Scroll performance
- Viewport quirks on specific browsers (Safari on iOS is notorious)

Open your deployed URL on your actual phone. Check:

- Does the hero look right?
- Can you tap every button comfortably?
- Does the hamburger menu work?
- Is the text readable without zooming?
- Does smooth scroll work?

:::tip
Send the deployed link to a friend and ask them to open it on their phone. Fresh eyes catch things you've gone blind to after hours of building.
:::

## Share the link

Your landing page is live — now make it visible:

- **Pin the repo** on your GitHub profile
- **Add a description and URL** to the GitHub repo (top of the repo page → gear icon)
- **Add the project to your portfolio** with a screenshot, live link, and source link
- **Share on LinkedIn** with a short post about what you built and learned

Write a strong GitHub README:

```markdown
# Brewly — Landing Page

A responsive landing page for a fictional coffee subscription service.

**Live demo:** https://yourusername.github.io/brewly/

## Built with
- Semantic HTML5
- CSS Custom Properties, Grid, Flexbox
- Mobile-first responsive design
- IntersectionObserver for scroll reveals

## What I learned
- Planning a page with a style guide before coding
- Building a responsive feature grid with CSS Grid
- Adding accessible, performant animations
```

## What makes this portfolio-worthy

Not every project deserves a spot on your portfolio. This one does if:

- **It's deployed and working** — live link, no broken images
- **It's responsive** — genuinely looks good on mobile, not just "doesn't break"
- **It demonstrates planning** — you can talk about your style guide and wireframe decisions
- **The code is clean** — semantic HTML, organized CSS with custom properties, no dead code
- **It passes audits** — Lighthouse scores above 90, HTML validates, contrast passes
- **You can explain it** — what you built, what was tricky, what you'd improve

:::quiz
Q: Why should you test your deployed site on a real phone instead of just using DevTools emulation?
- Real phones have better color accuracy
- Real devices reveal touch, scroll, and rendering issues that emulators can miss *
- DevTools can't simulate any mobile device
- Phones automatically optimize your CSS
E: DevTools emulation simulates screen size and pixel ratio, but it can't perfectly replicate touch interactions, real rendering engines (especially Safari on iOS), scroll performance, or how fonts actually look on a small screen.
:::

:::quiz
Q: What is the recommended approach for optimizing a large hero image?
- Delete the image entirely
- Resize it to a reasonable width, compress it, and consider using WebP format *
- Set its CSS width to 100px
- Move it to a separate server
E: Image optimization is a three-step process: resize to the dimensions you actually need, compress to reduce file size, and use modern formats like WebP that deliver smaller files at the same visual quality.
:::

## Final checklist

Before calling this project done:

- [ ] HTML validates (validator.w3.org)
- [ ] Every image has meaningful `alt` text
- [ ] Headings follow a logical h1 → h2 → h3 order
- [ ] All text meets 4.5:1 contrast ratio
- [ ] Keyboard navigation works (Tab through the page)
- [ ] Lighthouse scores: Performance 90+, Accessibility 95+
- [ ] Images are compressed (hero < 200 KB)
- [ ] Tested on a real phone
- [ ] Deployed to a public URL
- [ ] GitHub repo has a README with a live link
- [ ] Project added to your portfolio

## Recap

- **Validate HTML** at validator.w3.org — fix all errors.
- **Check accessibility:** alt text, heading order, contrast, keyboard navigation, link text.
- **Run Lighthouse** on mobile — aim for 90+ across all categories.
- **Compress images** — resize, compress, use WebP where possible.
- **Deploy** to GitHub Pages or Netlify.
- **Test on a real phone** — emulation doesn't catch everything.
- **Share the link** — pin the repo, update your portfolio, write a README.
- A portfolio-worthy project is deployed, responsive, accessible, and explainable.

**Congratulations** — you've planned, built, styled, polished, and shipped a complete landing page. That's the full professional workflow.
