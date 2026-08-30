# Images and Media

The web would be a wall of text without images, video, and audio. But media is also the leading cause of slow pages, layout shifts, and accessibility failures. In this lesson you'll learn to embed media that loads fast, looks sharp on every screen, and works for everyone.

## The `<img>` Element

The `<img>` tag is self-closing. Its two required attributes are `src` (the image source) and `alt` (alternative text).

```html
<img src="/images/hero.jpg" alt="A developer writing code at a sunny desk">
```

The `alt` text is critical:

- **Screen readers** read it aloud to blind users.
- **Search engines** use it to understand the image.
- **Broken images** show the alt text instead, so users know what was supposed to be there.

### Writing Good Alt Text

- Describe the **meaning**, not just what you see. "Bar chart showing revenue doubled in Q4" is better than "chart image."
- Keep it concise — typically one sentence.
- Don't start with "Image of…" or "Photo of…" — the screen reader already announces it as an image.
- For **decorative** images (borders, swooshes), use an empty `alt=""` so screen readers skip them entirely.

```html
<!-- Informative image -->
<img src="results.png"
     alt="Bar chart: signups grew from 1,000 to 4,500 over 6 months">

<!-- Decorative image -->
<img src="divider.svg" alt="">
```

:::warning
Never omit the `alt` attribute entirely. A missing `alt` is different from `alt=""`. Without `alt`, some screen readers will read the file name ("slash images slash D-S-C-zero-four-two-seven dot jpeg"), which is worse than useless.
:::

## Preventing Layout Shift with `width` and `height`

When the browser loads a page, it doesn't know how big an image is until the file downloads. Without dimensions, the page content jumps around as images pop in — this is called **Cumulative Layout Shift (CLS)**, and it's a Core Web Vital that affects your search ranking.

The fix is simple: always set `width` and `height` attributes.

```html
<img src="photo.jpg" alt="Team photo" width="800" height="600">
```

These attributes tell the browser the image's aspect ratio *before* it loads, so the browser reserves the right amount of space. Combined with CSS, the image stays responsive:

```css
img {
  max-width: 100%;
  height: auto;
}
```

:::key
Always include `width` and `height` on images. The browser uses them to calculate the aspect ratio and reserve space, eliminating layout shift — one of the easiest performance wins available.
:::

## Responsive Images with `srcset` and `sizes`

A single image can't serve every device well. A 2400px hero image wastes bandwidth on a phone; a 400px thumbnail looks blurry on a retina display. `srcset` and `sizes` solve this.

```html
<img
  src="photo-800.jpg"
  srcset="
    photo-400.jpg   400w,
    photo-800.jpg   800w,
    photo-1600.jpg 1600w
  "
  sizes="(max-width: 600px) 100vw,
         (max-width: 1200px) 50vw,
         33vw"
  alt="Mountain landscape at sunset"
  width="1600"
  height="900"
>
```

How it works:

- **`srcset`** lists available image files and their widths in pixels (the `w` descriptor).
- **`sizes`** tells the browser how wide the image will be displayed, using media conditions. "On screens up to 600px, this image takes 100% of the viewport width; up to 1200px, 50%; above that, 33%."
- The browser picks the best file based on the viewport width, device pixel ratio, and `sizes`.
- **`src`** is the fallback for browsers that don't support `srcset`.

:::tip
You don't pick which image loads — the browser does. Your job is to provide the options and describe how big the image will render. The browser handles device pixel ratios and network conditions automatically.
:::

## The `<picture>` Element

`<picture>` gives you full control over which image loads, based on media queries or format support. Use it for art direction (showing different crops on different screens) or serving modern formats with fallbacks.

```html
<!-- Art direction: different crop for mobile -->
<picture>
  <source media="(max-width: 600px)" srcset="hero-mobile.jpg">
  <source media="(min-width: 601px)" srcset="hero-desktop.jpg">
  <img src="hero-desktop.jpg" alt="Product showcase" width="1200" height="600">
</picture>

<!-- Format negotiation: serve WebP with JPEG fallback -->
<picture>
  <source type="image/avif" srcset="photo.avif">
  <source type="image/webp" srcset="photo.webp">
  <img src="photo.jpg" alt="Sunset over the ocean" width="800" height="600">
</picture>
```

The `<img>` inside `<picture>` is required — it's the fallback and where the `alt` text lives.

## Figure and Figcaption

When an image needs a visible caption, wrap it in `<figure>` with a `<figcaption>`.

```html
<figure>
  <img src="diagram.png"
       alt="Flowchart showing the HTTP request-response cycle"
       width="600" height="400">
  <figcaption>
    Figure 1: The HTTP request-response cycle between browser and server.
  </figcaption>
</figure>
```

`<figure>` semantically groups the image and its caption. Screen readers associate the caption with the image automatically. You can also use `<figure>` for code listings, charts, or any self-contained illustrative content.

:::tip
`<figcaption>` is the *visible* caption; `alt` is the *accessible* description. They serve different purposes and can coexist. The alt text describes what's in the image for people who can't see it; the caption adds context for everyone.
:::

## Lazy Loading

Images below the fold don't need to load until the user scrolls to them. Native lazy loading is built into HTML:

```html
<img src="photo.jpg"
     alt="Office workspace"
     width="800" height="600"
     loading="lazy">
```

The `loading="lazy"` attribute tells the browser to defer loading until the image is near the viewport. This can dramatically cut initial page load time on image-heavy pages.

:::warning
Don't lazy-load your hero image or any image visible in the initial viewport (above the fold). Those should load immediately. Use `loading="eager"` (the default) or `fetchpriority="high"` for critical above-the-fold images.
:::

```html
<!-- Above the fold: load immediately, high priority -->
<img src="hero.jpg" alt="Welcome banner"
     width="1200" height="600"
     fetchpriority="high">

<!-- Below the fold: lazy load -->
<img src="team.jpg" alt="Team photo"
     width="800" height="600"
     loading="lazy">
```

## Video

The `<video>` element embeds video with built-in playback controls.

```html
<video controls width="640" height="360" preload="metadata">
  <source src="demo.webm" type="video/webm">
  <source src="demo.mp4" type="video/mp4">
  <p>Your browser doesn't support HTML video.
     <a href="demo.mp4">Download the video</a>.</p>
</video>
```

Key attributes:

- **`controls`** — shows play/pause, volume, and seek. Always include this unless you're building a custom player.
- **`preload="metadata"`** — downloads only enough to show duration and dimensions, not the whole file.
- **`poster="thumbnail.jpg"`** — a preview image shown before the video plays.
- **`autoplay`** — plays automatically. Browsers block autoplay with sound, so pair with `muted` if you need autoplay.

Multiple `<source>` elements let you offer different formats; the browser picks the first one it supports.

:::warning
Auto-playing videos with sound violate WCAG accessibility guidelines and annoy users. If you must autoplay, add `muted` and provide visible controls to pause.
:::

## Audio

`<audio>` works like `<video>` but without the visual.

```html
<audio controls preload="metadata">
  <source src="podcast.ogg" type="audio/ogg">
  <source src="podcast.mp3" type="audio/mpeg">
  <p>Your browser doesn't support HTML audio.
     <a href="podcast.mp3">Download the episode</a>.</p>
</audio>
```

Always provide a **transcript** for audio content and **captions** for video to meet accessibility standards.

## Embedding with `<iframe>`

`<iframe>` embeds another HTML page inside yours — commonly used for YouTube videos, maps, and third-party widgets.

```html
<iframe
  src="https://www.youtube.com/embed/dQw4w9WgXcQ"
  width="560"
  height="315"
  title="Product demo video"
  loading="lazy"
  allow="accelerometer; autoplay; clipboard-write; encrypted-media"
  allowfullscreen
></iframe>
```

Always include:

- **`title`** — describes the iframe content for screen readers ("Product demo video" is far better than silence).
- **`loading="lazy"`** — defer loading off-screen iframes.
- **`allow`** — a permissions policy that restricts what the embedded page can do.

:::key
Every `<iframe>` must have a `title` attribute. Without it, screen readers announce "frame" with no description, leaving users unable to understand what the embedded content is.
:::

:::quiz
Q: What causes Cumulative Layout Shift (CLS) with images, and how do you prevent it?
- Large file sizes; use compression
- Missing `width` and `height` attributes; add them so the browser reserves space *
- Using JPEG instead of PNG; switch formats
E: Without `width` and `height`, the browser doesn't know the image's aspect ratio until it loads. Content shifts as images pop in. Setting these attributes lets the browser reserve the exact space needed.
:::

:::quiz
Q: When should you use `loading="lazy"` on an image?
- On every image for best performance
- Only on images below the fold that aren't visible in the initial viewport *
- Only on decorative images
E: Lazy loading defers images the user hasn't scrolled to yet. Above-the-fold images should load immediately (with default `loading="eager"` or `fetchpriority="high"`) to avoid a blank hero section.
:::

## Recap

- Every `<img>` needs `src`, `alt`, `width`, and `height`. Include all four, always.
- Use `srcset`/`sizes` for responsive resolution switching; use `<picture>` for art direction and format negotiation.
- Wrap captioned images in `<figure>` + `<figcaption>`.
- Use `loading="lazy"` for below-the-fold images and iframes; keep above-the-fold assets eager.
- Provide `controls` on `<video>` and `<audio>`; never autoplay with sound.
- Every `<iframe>` needs a `title` for accessibility.

**Next up:** Lists and Tables — structuring data clearly and semantically.
