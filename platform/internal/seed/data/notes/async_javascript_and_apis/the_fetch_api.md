# The Fetch API

`fetch` is the modern way to make HTTP requests from JavaScript. This lesson covers the full API: making requests, reading responses, sending data with POST, building URLs safely, and cancelling requests that are no longer needed.

## Basic GET request

At its simplest, `fetch` takes a URL and returns a promise that resolves to a Response object:

```js
const response = await fetch("https://jsonplaceholder.typicode.com/posts");
```

That gives you the **Response** — the envelope. To read the body, you call one of its methods:

```js
const data = await response.json();  // parse body as JSON
const text = await response.text();  // read body as plain text
const blob = await response.blob();  // read as binary (for images, files)
```

:::key
`fetch` returns a Response. `response.json()` returns the parsed data. Both are promises — you need two `await` calls.
:::

## The Response object

The Response carries useful information beyond the body:

```js
const response = await fetch("/api/users");

console.log(response.ok);         // true if status is 200-299
console.log(response.status);     // 200, 404, 500, etc.
console.log(response.statusText); // "OK", "Not Found", etc.
console.log(response.headers.get("Content-Type")); // "application/json"
console.log(response.url);        // the final URL (after redirects)
```

### Checking for errors

```js
const response = await fetch("/api/data");

if (!response.ok) {
  throw new Error(`Request failed: ${response.status} ${response.statusText}`);
}

const data = await response.json();
```

:::warning
Always check `response.ok` before reading the body. A 404 or 500 response still fulfills the fetch promise — it does not throw. Without this check, you might try to parse an error page as JSON and get a confusing error.
:::

## Request options

The second argument to `fetch` is an options object:

```js
const response = await fetch("/api/data", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "Authorization": "Bearer my-token-here",
  },
  body: JSON.stringify({ name: "Ada", role: "engineer" }),
});
```

The most common options:

| Option | What it does |
|---|---|
| `method` | HTTP method: `"GET"`, `"POST"`, `"PUT"`, `"PATCH"`, `"DELETE"` |
| `headers` | Object of request headers |
| `body` | The request body (string, FormData, Blob, etc.) |
| `signal` | An AbortSignal for cancellation |

## Sending JSON with POST

The most common pattern for sending data to an API:

```js
async function createUser(userData) {
  const response = await fetch("/api/users", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    throw new Error(`Failed to create user: ${response.status}`);
  }

  return await response.json(); // the server usually returns the created object
}

const newUser = await createUser({ name: "Ada", email: "ada@example.com" });
```

:::key
When sending JSON, you must do two things: set the `Content-Type` header to `"application/json"` and stringify the body with `JSON.stringify`. Missing either one is a common source of "400 Bad Request" errors.
:::

## Sending form data

For file uploads or multipart forms, use `FormData` directly as the body:

```js
const form = document.querySelector("#upload-form");

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const formData = new FormData(form);

  const response = await fetch("/api/upload", {
    method: "POST",
    body: formData, // no Content-Type header — the browser sets it automatically
  });

  const result = await response.json();
  console.log(result);
});
```

:::tip
Do **not** set the `Content-Type` header when sending `FormData`. The browser automatically sets it to `multipart/form-data` with the correct boundary string. Setting it yourself breaks the upload.
:::

## Building URLs with URLSearchParams

Never concatenate query strings by hand — special characters break the URL:

```js
// BAD — breaks if searchTerm contains & or spaces
const url = "/api/search?q=" + searchTerm + "&limit=10";

// GOOD — URLSearchParams handles encoding
const params = new URLSearchParams({
  q: searchTerm,
  limit: 10,
  category: "design",
});

const url = `/api/search?${params}`;
// /api/search?q=web+design&limit=10&category=design
```

Adding and modifying parameters:

```js
const params = new URLSearchParams();
params.set("page", 2);
params.set("sort", "date");
params.append("tag", "javascript");
params.append("tag", "css"); // duplicate keys are fine

console.log(params.toString());
// "page=2&sort=date&tag=javascript&tag=css"
```

## AbortController — cancelling requests

When a user navigates away or types a new search query, you should cancel the previous request. `AbortController` does this:

```js
let controller = null;

async function search(query) {
  // Cancel the previous request if it is still running
  if (controller) controller.abort();

  controller = new AbortController();

  try {
    const response = await fetch(`/api/search?q=${query}`, {
      signal: controller.signal,
    });
    const results = await response.json();
    renderResults(results);
  } catch (error) {
    if (error.name === "AbortError") {
      // Request was cancelled — this is expected, not an error
      return;
    }
    console.error("Search failed:", error);
  }
}
```

:::tip
When you abort a request, `fetch` rejects with an `AbortError`. Check `error.name === "AbortError"` to distinguish intentional cancellations from real failures. Do not show an error message for cancelled requests.
:::

## A reusable fetch wrapper

Many projects define a helper that handles JSON parsing, error checking, and headers:

```js
async function api(endpoint, options = {}) {
  const { method = "GET", body, headers = {} } = options;

  const config = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const response = await fetch(`/api${endpoint}`, config);

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || `HTTP ${response.status}`);
  }

  return response.json();
}

// Usage
const users = await api("/users");
const newPost = await api("/posts", {
  method: "POST",
  body: { title: "Hello", content: "World" },
});
```

## Practice

:::quiz
Q: You are sending a JSON body with `fetch` but the server returns "400 Bad Request." What is the most likely cause?
- The URL is wrong
- You forgot to set `Content-Type: application/json` and/or forgot to `JSON.stringify` the body *
- The server is down
- You used the wrong HTTP method
E: Servers expect JSON bodies to be sent as a string with the `Content-Type: application/json` header. Missing either one causes the server to misinterpret the body, resulting in a 400 error.
:::

:::quiz
Q: A user types a new search query while the previous search request is still in progress. What should you do?
- Let both requests complete and show whichever finishes last
- Abort the previous request with AbortController before starting the new one *
- Block the input until the first request finishes
- Ignore the new query
E: Aborting the previous request prevents race conditions where an older, slower response overwrites newer results. `AbortController` cancels the in-flight request cleanly, and you can ignore the resulting `AbortError`.
:::

## Recap

- `fetch(url)` returns a promise resolving to a **Response** object. Read the body with `.json()`, `.text()`, or `.blob()`.
- Check **`response.ok`** before reading the body — fetch does not reject on 404/500.
- Set `method`, `headers`, and `body` in the **options object** (second argument).
- For JSON bodies: set `Content-Type: application/json` and use **`JSON.stringify`**.
- For FormData/file uploads: do **not** set Content-Type — the browser handles it.
- Build URLs safely with **`URLSearchParams`** — it handles encoding automatically.
- Cancel requests with **`AbortController`** — pass `controller.signal` to fetch and call `controller.abort()`.
- Build a **reusable wrapper** to handle JSON, headers, and error checking in one place.

**Next up:** HTTP Methods, Headers, and CORS — understanding the protocol behind every request.
