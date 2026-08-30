# Error Handling in Async Code

Network requests fail. Servers return errors. JSON parsing breaks. If your code does not handle these cases, your app either crashes silently or shows a blank screen. This lesson covers every layer of error handling you need for real async JavaScript.

## The two sources of failure

When you make a fetch request, two different things can go wrong:

1. **Network failure** — no internet, DNS error, server unreachable. `fetch` rejects the promise.
2. **HTTP error** — the server responds, but with an error status like 404 or 500. `fetch` *fulfills* the promise — you just get a bad response.

This distinction catches many beginners off guard.

:::key
`fetch` only rejects on network failure. A 404 or 500 response is still a fulfilled promise — you must check `response.ok` yourself.
:::

## try/catch with await

The standard pattern for async/await error handling:

```js
async function loadUsers() {
  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users");

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const users = await response.json();
    return users;
  } catch (error) {
    console.error("Failed to load users:", error.message);
    return [];
  }
}
```

The `try` block covers both the network call and the JSON parsing. If *anything* fails — network down, bad status, corrupt JSON — the `catch` block handles it.

## .catch on promise chains

The `.then` chain equivalent:

```js
function loadUsers() {
  return fetch("https://jsonplaceholder.typicode.com/users")
    .then(response => {
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      return response.json();
    })
    .catch(error => {
      console.error("Failed:", error.message);
      return [];
    });
}
```

A single `.catch` at the end of the chain handles errors from every step above it.

## Checking response.ok properly

Do not skip this step. A 404 response will happily pass through your `.then` and blow up when you try to use the data:

```js
// BAD — no status check
async function getUser(id) {
  const response = await fetch(`/api/users/${id}`);
  const user = await response.json(); // if 404, this might parse an error page as JSON
  return user;
}

// GOOD — explicit check
async function getUser(id) {
  const response = await fetch(`/api/users/${id}`);

  if (!response.ok) {
    throw new Error(`User ${id} not found (${response.status})`);
  }

  return await response.json();
}
```

:::tip
`response.ok` is `true` for status codes 200–299. For more specific handling, check `response.status` directly:

```js
if (response.status === 404) {
  return null; // not found is not a crash, it's data
}
if (response.status === 401) {
  redirectToLogin();
  return;
}
if (!response.ok) {
  throw new Error(`Unexpected error: ${response.status}`);
}
```
:::

## Custom error messages

Raw error messages are for developers, not users. Transform them:

```js
async function loadPosts() {
  try {
    const response = await fetch("/api/posts");
    if (!response.ok) throw new Error(`Status: ${response.status}`);
    return await response.json();
  } catch (error) {
    // Log the technical detail for debugging
    console.error("loadPosts failed:", error);

    // Show a human-readable message
    showErrorMessage("We couldn't load posts right now. Please try again.");
    return [];
  }
}
```

:::key
Separate the developer message (logged to console) from the user message (displayed in the UI). Users do not need to see "TypeError: Failed to fetch" — they need "Something went wrong. Please try again."
:::

## Fallback values

When a request fails, return a sensible default instead of letting the error propagate and crash the render:

```js
async function loadConfig() {
  try {
    const response = await fetch("/api/config");
    if (!response.ok) throw new Error("Config not available");
    return await response.json();
  } catch {
    // Fall back to defaults
    return {
      theme: "light",
      language: "en",
      pageSize: 20,
    };
  }
}
```

Your app keeps working with defaults even if the config endpoint is down.

## Retry patterns

Some errors are temporary — a flaky connection, a server restarting. Retrying can recover without user intervention:

```js
async function fetchWithRetry(url, options = {}, retries = 3) {
  for (let attempt = 1; attempt <= retries; attempt++) {
    try {
      const response = await fetch(url, options);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      return await response.json();
    } catch (error) {
      console.warn(`Attempt ${attempt}/${retries} failed:`, error.message);

      if (attempt === retries) {
        throw new Error(`Failed after ${retries} attempts: ${error.message}`);
      }

      // Exponential backoff: 1s, 2s, 4s...
      const delay = Math.pow(2, attempt - 1) * 1000;
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
}
```

:::tip
Exponential backoff increases the wait between retries: 1 second, 2 seconds, 4 seconds. This gives overwhelmed servers time to recover instead of hammering them with rapid retries.
:::

## User-facing error messages in the UI

Integrate error handling into your render cycle:

```js
const state = {
  users: [],
  isLoading: false,
  error: null,
};

async function loadUsers() {
  state.isLoading = true;
  state.error = null;
  render();

  try {
    const response = await fetch("/api/users");
    if (!response.ok) throw new Error(`Server error (${response.status})`);
    state.users = await response.json();
  } catch (error) {
    state.error = "Unable to load users. Please check your connection and try again.";
    console.error(error);
  } finally {
    state.isLoading = false;
    render();
  }
}

function render() {
  const container = document.querySelector("#content");

  if (state.isLoading) {
    container.innerHTML = '<p class="loading">Loading...</p>';
    return;
  }

  if (state.error) {
    container.innerHTML = `
      <div class="error-state">
        <p>${state.error}</p>
        <button id="retry-btn">Try Again</button>
      </div>
    `;
    document.querySelector("#retry-btn").addEventListener("click", loadUsers);
    return;
  }

  container.innerHTML = state.users
    .map(u => `<div class="card">${u.name}</div>`)
    .join("");
}
```

## Error handling in Promise.all

`Promise.all` rejects as soon as *any* promise rejects, and you lose the results of the others:

```js
// If posts fetch fails, you lose users too
try {
  const [users, posts] = await Promise.all([fetchUsers(), fetchPosts()]);
} catch (error) {
  // Which one failed? Hard to tell
}
```

Use `Promise.allSettled` when you want all results regardless of individual failures:

```js
const results = await Promise.allSettled([fetchUsers(), fetchPosts()]);

const users = results[0].status === "fulfilled" ? results[0].value : [];
const posts = results[1].status === "fulfilled" ? results[1].value : [];
```

## Practice

:::quiz
Q: You fetch a URL and get a 404 response. What does `fetch` do?
- Rejects the promise, jumping to catch
- Fulfills the promise with a Response where `response.ok` is false *
- Returns null
- Throws a SyntaxError
E: `fetch` only rejects on network failure (no internet, DNS error). A 404 is a successful HTTP response — the promise fulfills, but `response.ok` is `false` and `response.status` is 404. You must check this yourself.
:::

:::quiz
Q: Your app fetches config from an API. The endpoint is down. What is the best approach?
- Crash and show a stack trace
- Log the error and return sensible default values so the app keeps working *
- Retry infinitely until the server responds
- Ignore the error completely
E: Falling back to default values lets the app remain functional even when a non-critical endpoint is unavailable. Log the technical error for debugging, but do not crash the user's experience.
:::

## Recap

- `fetch` rejects only on **network failure**. Check `response.ok` for HTTP errors (404, 500).
- Use **`try/catch`** with `await` or **`.catch`** on chains — always handle errors.
- Separate **developer logs** (`console.error`) from **user messages** (displayed in the UI).
- Return **fallback values** when a request fails so the app does not crash.
- **Retry with exponential backoff** for transient failures: 1s, 2s, 4s between attempts.
- Show a **retry button** in the UI so users can recover without refreshing the page.
- Use **`Promise.allSettled`** when you need all results even if some requests fail.

**Next up:** Promise.all and Promise.race — running multiple async operations together.
