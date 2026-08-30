# Promise.all and Promise.race

When you need to run multiple async operations, JavaScript gives you four tools that handle them differently. Knowing which one to use — and why — is the difference between fast, resilient code and code that is slow or fragile.

## Promise.all — run in parallel, fail fast

`Promise.all` takes an array of promises, runs them all simultaneously, and resolves with an array of results **in the same order**:

```js
async function loadDashboard() {
  const [users, posts, comments] = await Promise.all([
    fetch("/api/users").then(r => r.json()),
    fetch("/api/posts").then(r => r.json()),
    fetch("/api/comments").then(r => r.json()),
  ]);

  console.log(users.length, posts.length, comments.length);
}
```

All three fetches run at the same time. Total time = the **slowest** request, not the sum.

### The fail-fast catch

If **any** promise rejects, `Promise.all` rejects immediately with that error. The other promises keep running in the background, but you never get their results:

```js
try {
  const [a, b, c] = await Promise.all([
    fetchUsers(),    // succeeds
    fetchPosts(),    // REJECTS
    fetchComments(), // succeeds, but we never see the result
  ]);
} catch (error) {
  // Only the first rejection is caught
  console.log(error.message);
}
```

:::key
Use `Promise.all` when you need **all results to succeed** and want to bail immediately if anything fails. Think of it as "all or nothing."
:::

## Promise.allSettled — always get every result

`Promise.allSettled` waits for **all** promises to settle (fulfill or reject) and gives you the outcome of each:

```js
const results = await Promise.allSettled([
  fetch("/api/users").then(r => r.json()),
  fetch("/api/posts").then(r => r.json()),
  fetch("/api/bad-url").then(r => r.json()),
]);

results.forEach((result, i) => {
  if (result.status === "fulfilled") {
    console.log(`Request ${i} succeeded:`, result.value);
  } else {
    console.log(`Request ${i} failed:`, result.reason.message);
  }
});
```

Each result is an object with:
- `{ status: "fulfilled", value: ... }` — on success
- `{ status: "rejected", reason: ... }` — on failure

### When to use it

```js
// Load widgets for a dashboard — show what you can, skip what fails
const results = await Promise.allSettled([
  fetchWeather(),
  fetchNews(),
  fetchStockPrices(),
]);

const weather = results[0].status === "fulfilled" ? results[0].value : null;
const news = results[1].status === "fulfilled" ? results[1].value : null;
const stocks = results[2].status === "fulfilled" ? results[2].value : null;

renderDashboard({ weather, news, stocks });
```

:::tip
Use `Promise.allSettled` when individual failures should not block the others. Dashboards, multi-widget pages, and batch operations where partial results are acceptable.
:::

## Promise.race — first to settle wins

`Promise.race` resolves or rejects with the **first** promise to settle, whether it fulfilled or rejected:

```js
const result = await Promise.race([
  fetch("/api/primary-server").then(r => r.json()),
  fetch("/api/backup-server").then(r => r.json()),
]);
// whichever server responds first wins
```

### Practical use: timeout

The most common use of `Promise.race` is implementing a timeout:

```js
function timeout(ms) {
  return new Promise((_, reject) =>
    setTimeout(() => reject(new Error("Request timed out")), ms)
  );
}

async function fetchWithTimeout(url, ms = 5000) {
  const response = await Promise.race([
    fetch(url),
    timeout(ms),
  ]);
  return response.json();
}

try {
  const data = await fetchWithTimeout("/api/data", 3000);
} catch (error) {
  console.log(error.message); // "Request timed out" if the fetch takes > 3s
}
```

:::warning
`Promise.race` settles with the **first** promise to finish, whether it succeeds or fails. If a fast promise rejects before a slower one resolves, you get the rejection — even if the slower one would have succeeded.
:::

## Promise.any — first success wins

`Promise.any` resolves with the **first fulfilled** promise, ignoring rejections. It only rejects if *all* promises reject:

```js
const fastest = await Promise.any([
  fetch("https://cdn1.example.com/data.json").then(r => r.json()),
  fetch("https://cdn2.example.com/data.json").then(r => r.json()),
  fetch("https://cdn3.example.com/data.json").then(r => r.json()),
]);
// The first CDN to respond successfully wins
```

If all three fail, `Promise.any` rejects with an `AggregateError` containing all the individual errors:

```js
try {
  const result = await Promise.any([failingPromise1, failingPromise2]);
} catch (error) {
  console.log(error.errors); // array of all individual errors
}
```

## Choosing the right one

| Method | Resolves when | Rejects when | Use case |
|---|---|---|---|
| `Promise.all` | All fulfill | Any rejects (fail-fast) | Need every result |
| `Promise.allSettled` | All settle | Never (always fulfills) | Partial results OK |
| `Promise.race` | First settles | First settles (if rejected) | Timeouts, fastest source |
| `Promise.any` | First fulfills | All reject | Fallback servers, fastest success |

## Practical examples

### Loading page data with Promise.all

```js
async function loadPageData() {
  try {
    const [user, notifications, feed] = await Promise.all([
      fetchJSON("/api/me"),
      fetchJSON("/api/notifications"),
      fetchJSON("/api/feed"),
    ]);

    renderPage({ user, notifications, feed });
  } catch (error) {
    showError("Failed to load page data");
  }
}
```

### Batch processing with Promise.allSettled

```js
async function sendInvites(emails) {
  const results = await Promise.allSettled(
    emails.map(email =>
      fetch("/api/invite", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      })
    )
  );

  const sent = results.filter(r => r.status === "fulfilled").length;
  const failed = results.filter(r => r.status === "rejected").length;

  showMessage(`${sent} invites sent, ${failed} failed`);
}
```

### Image loading with Promise.any

```js
async function loadProfilePicture(userId) {
  try {
    const img = await Promise.any([
      loadImage(`https://cdn1.example.com/avatars/${userId}.jpg`),
      loadImage(`https://cdn2.example.com/avatars/${userId}.jpg`),
    ]);
    container.append(img);
  } catch {
    container.append(createPlaceholderAvatar());
  }
}
```

## Combining with async/await

All four methods return promises, so they work seamlessly with `await`:

```js
async function loadWithFallbacks() {
  // Try multiple sources, take the fastest success
  const data = await Promise.any([
    fetchJSON("/api/v2/data"),
    fetchJSON("/api/v1/data"),
  ]).catch(() => ({ items: [] })); // fallback if all fail

  // Load related data in parallel
  const [comments, tags] = await Promise.all([
    fetchJSON(`/api/comments?dataId=${data.id}`),
    fetchJSON("/api/tags"),
  ]);

  return { data, comments, tags };
}
```

## Practice

:::quiz
Q: You are building a dashboard with three widgets (weather, news, stocks). If the news API fails, you still want to show weather and stocks. Which method should you use?
- `Promise.all`
- `Promise.allSettled` *
- `Promise.race`
- `Promise.any`
E: `Promise.allSettled` waits for all promises and gives you the result of each one individually, whether it succeeded or failed. You can render the widgets that loaded and show an error state for the ones that did not.
:::

:::quiz
Q: You want to fetch from a server but cancel the request if it takes longer than 3 seconds. Which method helps here?
- `Promise.all`
- `Promise.allSettled`
- `Promise.race` *
- `Promise.any`
E: `Promise.race` settles with the first promise to finish. Race your fetch against a timeout promise that rejects after 3 seconds. If the timeout wins, you get the rejection and can show an error.
:::

## Recap

- **`Promise.all`** — runs in parallel, resolves when all succeed, rejects if any fail. "All or nothing."
- **`Promise.allSettled`** — runs in parallel, always resolves with status of each. "Give me everything, success or failure."
- **`Promise.race`** — resolves/rejects with the first to settle. Use for timeouts and fastest-source patterns.
- **`Promise.any`** — resolves with the first success, ignores rejections. Use for fallback servers and redundancy.
- All four return promises and work with `await`.
- Choose based on your resilience needs: do you need all results, partial results, the fastest, or the first success?

**Next up:** The Fetch API — the full toolbox for making HTTP requests.
