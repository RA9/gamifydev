# DevTools: Network and Performance

The Network tab shows every file and API request your page makes. The Performance tab reveals how fast it loads. Together, they help you understand what your page is actually doing behind the scenes — and how to make it faster.

## Network tab overview

Open DevTools → **Network** tab. Refresh the page to capture all requests.

Each row is a request the browser made: HTML files, CSS, JavaScript, images, fonts, API calls. The bottom bar shows a summary:

```text
23 requests | 450 KB transferred | 1.2 MB resources | Finish: 1.4s
```

This tells you: the page made 23 requests, downloaded 450 KB of data, and finished loading in 1.4 seconds.

:::tip
Always refresh the page *after* opening the Network tab. It only records requests that happen while it's open. If the tab was closed during page load, you'll see an empty list.
:::

## Reading request rows

Each row shows key information at a glance:

| Column | What it tells you |
|--------|------------------|
| **Name** | The file or endpoint requested |
| **Status** | HTTP status code (200 = OK, 304 = cached, 404 = not found, 500 = server error) |
| **Type** | The resource type (document, stylesheet, script, img, xhr/fetch) |
| **Size** | How much data was transferred (smaller = faster) |
| **Time** | How long the request took from start to finish |

**Status codes to know:**

```text
200  OK — the file loaded successfully
301  Moved Permanently — URL has changed, browser follows redirect
304  Not Modified — using cached version (no download needed)
404  Not Found — file doesn't exist at that URL
403  Forbidden — server refuses to serve it
500  Internal Server Error — something broke on the server
```

:::key
A `404` on a CSS file means your stylesheet isn't loading. A `404` on an image means a broken `src` path. The Network tab reveals these issues instantly — you don't have to guess why something "isn't working."
:::

## Filtering by type

The filter bar at the top lets you narrow down requests by type:

- **All** — everything
- **JS** — JavaScript files only
- **CSS** — stylesheets only
- **Img** — images only
- **XHR/Fetch** — API calls and data requests
- **Font** — web fonts
- **Doc** — HTML documents

Click **XHR/Fetch** when debugging API calls. Click **Img** to see which images are large. Click **JS** to check if your scripts are loading.

## Inspecting request and response headers

Click any request row to open its details panel. The tabs you'll use most:

**Headers tab:**
- **Request Headers** — what the browser sent (including cookies, content type, authorization)
- **Response Headers** — what the server sent back (content type, cache settings, CORS headers)

```text
Request:
  GET /api/users HTTP/2
  Accept: application/json
  Authorization: Bearer eyJhbG...

Response:
  HTTP/2 200 OK
  Content-Type: application/json
  Cache-Control: max-age=3600
```

**Preview tab** — the response body rendered in a readable format (especially useful for JSON).

**Response tab** — the raw response body as text.

## Previewing JSON responses

When you click an API request (XHR/Fetch type), the **Preview** tab shows the JSON response as a collapsible, color-coded tree:

```json
{
  "users": [
    { "id": 1, "name": "Ada", "role": "engineer" },
    { "id": 2, "name": "Grace", "role": "architect" }
  ],
  "total": 2
}
```

This is the fastest way to verify that an API is returning the data you expect, in the structure you expect.

:::tip
Right-click the request row and choose **Copy → Copy as fetch** to get a ready-made `fetch()` call you can paste into the console and tweak. Great for testing API calls with different parameters.
:::

## Throttling — simulating slow connections

Your page loads instantly on your fast Wi-Fi. But what about a user on a slow mobile connection?

In the Network tab, find the **throttling** dropdown (usually says "No throttling" or "Online"). Change it to:

- **Fast 3G** — simulates a moderate mobile connection (~1.6 Mbps)
- **Slow 3G** — simulates a poor connection (~400 Kbps)
- **Offline** — no connection at all

Now refresh. Watch how long each resource takes to load. A 2 MB image that loads instantly on your connection takes 5+ seconds on Slow 3G.

:::warning
Always test on Slow 3G at least once before deploying. If your page is usable on a slow connection, it'll be great on a fast one. The reverse isn't true — a page that works on fast Wi-Fi can be unusable on 3G.
:::

## Disable cache

Check the **Disable cache** checkbox at the top of the Network tab. This forces the browser to download every file fresh instead of using cached versions.

Keep this checked during development. Without it, you might be testing old CSS or JavaScript that the browser cached from a previous visit.

## The Lighthouse tab

Lighthouse runs an automated audit of your page and scores it on four categories:

1. **Performance** — how fast the page loads
2. **Accessibility** — how usable it is for people with disabilities
3. **Best Practices** — security, modern APIs, console errors
4. **SEO** — search engine optimization basics

To run it:
1. Open DevTools → **Lighthouse** tab
2. Select **Mobile** (more demanding and realistic)
3. Check all categories
4. Click **Analyze page load**

Lighthouse gives you a score (0–100) for each category and specific, actionable recommendations.

```text
Performance:     92
Accessibility:   88  ← "Image elements do not have [alt] attributes"
Best Practices:  100
SEO:             78  ← "Document does not have a meta description"
```

Each recommendation links to documentation explaining why it matters and how to fix it.

## Performance basics — Core Web Vitals

Google uses three key metrics called **Core Web Vitals** to measure user experience:

### FCP — First Contentful Paint

**How long until the user sees something.** Measures when the first text or image appears on screen.

- **Good:** under 1.8 seconds
- **Needs improvement:** 1.8–3.0 seconds
- **Poor:** over 3.0 seconds

**How to improve:** reduce render-blocking CSS/JS, use system fonts or preload web fonts.

### LCP — Largest Contentful Paint

**How long until the main content is visible.** Measures when the largest image or text block finishes rendering.

- **Good:** under 2.5 seconds
- **Needs improvement:** 2.5–4.0 seconds
- **Poor:** over 4.0 seconds

**How to improve:** compress and resize hero images, preload critical resources.

### CLS — Cumulative Layout Shift

**How much the page jumps around while loading.** Measures visual stability — content shouldn't shift after it appears.

- **Good:** under 0.1
- **Needs improvement:** 0.1–0.25
- **Poor:** over 0.25

**How to improve:** set `width` and `height` on images, avoid inserting content above existing content.

```html
<!-- Good: dimensions prevent layout shift -->
<img src="hero.jpg" alt="Hero image" width="1200" height="800" />

<!-- Bad: no dimensions, image causes layout shift when it loads -->
<img src="hero.jpg" alt="Hero image" />
```

:::key
Core Web Vitals (FCP, LCP, CLS) aren't abstract metrics — they measure what users actually experience. A page that scores well on these feels fast, stable, and responsive. Google also uses them as search ranking signals.
:::

## The waterfall chart

In the Network tab, the colored bars on each row form a **waterfall chart** showing request timing:

```text
index.html    ████░░░░░░░░░░░░░░░░
style.css         ░░████░░░░░░░░░░░
app.js            ░░████░░░░░░░░░░░
hero.jpg              ░░░░████████░
font.woff                 ░░████░░░
```

Reading the waterfall:
- **Horizontal position** = when the request started (left = early, right = late)
- **Bar length** = how long the request took
- **Gaps** between bars = wasted time (opportunity to optimize)
- Resources that start **after** other resources finish are blocked — they couldn't start until a dependency loaded

## Practical optimization workflow

1. Open Network tab → refresh with **Disable cache** on
2. Check the bottom summary — how many requests, how much data, how long?
3. Sort by **Size** — find the biggest files
4. Sort by **Time** — find the slowest requests
5. Filter by **Img** — are there oversized images?
6. Run **Lighthouse** on mobile — follow its recommendations
7. Test on **Slow 3G** — is the page usable?

:::quiz
Q: You see a CSS file with status code 404 in the Network tab. What does this mean?
- The CSS file has a syntax error
- The browser cannot find the CSS file at the URL specified in your HTML *
- The CSS file is too large to load
- The server is down
E: A 404 status means "Not Found." The browser tried to load the CSS file but the URL was wrong — check the `href` path in your `<link>` tag. It's probably a typo or incorrect relative path.
:::

:::quiz
Q: What does CLS (Cumulative Layout Shift) measure?
- How long until the first pixel appears
- How much the page content visually shifts while loading *
- The total size of all CSS files
- How many API calls the page makes
E: CLS measures visual stability — how much elements move around as the page loads. A high CLS means content jumps unexpectedly (like text shifting down when an image loads above it), which frustrates users.
:::

## Recap

- The **Network tab** records every request your page makes — filter by type to focus.
- **Status codes:** 200 = OK, 304 = cached, 404 = not found, 500 = server error.
- **Inspect headers and responses** by clicking a request row — use Preview for JSON.
- **Throttling** simulates slow connections — always test on Slow 3G before shipping.
- **Disable cache** during development to avoid testing stale files.
- **Lighthouse** audits Performance, Accessibility, Best Practices, and SEO with actionable recommendations.
- **Core Web Vitals:** FCP (first paint), LCP (main content visible), CLS (visual stability).
- The **waterfall chart** shows request timing and reveals bottlenecks.

**Next up:** Branching, Merging, and Pull Requests — the Git workflow that professional teams use every day.
