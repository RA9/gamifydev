# HTTP Methods Headers and CORS

Every `fetch` call speaks HTTP — the protocol that powers the web. Understanding methods, headers, and the CORS security model lets you interact with any API confidently and debug the errors that confuse most beginners.

## HTTP methods — what you are asking the server to do

Each request includes a **method** that tells the server what action you want:

| Method | Purpose | Has body? | Example |
|---|---|---|---|
| `GET` | Read data | No | Fetch a list of users |
| `POST` | Create a new resource | Yes | Submit a new blog post |
| `PUT` | Replace a resource entirely | Yes | Update a full user profile |
| `PATCH` | Update part of a resource | Yes | Change just the user's email |
| `DELETE` | Remove a resource | Usually no | Delete a comment |

```js
// GET — the default, no body
const response = await fetch("/api/users");

// POST — create
await fetch("/api/users", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ name: "Ada", email: "ada@example.com" }),
});

// PUT — full replacement
await fetch("/api/users/1", {
  method: "PUT",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ name: "Ada", email: "ada@new.com", role: "admin" }),
});

// PATCH — partial update
await fetch("/api/users/1", {
  method: "PATCH",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ email: "ada@updated.com" }),
});

// DELETE
await fetch("/api/users/1", { method: "DELETE" });
```

:::tip
`PUT` replaces the entire resource — if you omit a field, it may be deleted on the server. `PATCH` updates only the fields you send. When in doubt, check the API's documentation for which it expects.
:::

## Essential request headers

Headers are metadata sent with the request. The most common ones:

### Content-Type

Tells the server what format the body is in:

```js
headers: { "Content-Type": "application/json" }     // JSON body
headers: { "Content-Type": "application/x-www-form-urlencoded" } // form data
// For FormData, let the browser set this automatically
```

### Authorization

Sends credentials — most commonly a Bearer token:

```js
headers: {
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

This is how APIs know *who* is making the request. The token usually comes from a login flow.

### Accept

Tells the server what response format you prefer:

```js
headers: {
  "Accept": "application/json"  // "I want JSON back"
}
```

Most JSON APIs default to JSON, but setting `Accept` explicitly is good practice.

### Combining headers

```js
const response = await fetch("/api/posts", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "Authorization": "Bearer my-token",
    "Accept": "application/json",
  },
  body: JSON.stringify({ title: "New Post" }),
});
```

:::warning
**Never put secret API keys in frontend JavaScript.** Anyone can open DevTools and see every header your code sends. API keys that grant write access, cost money, or access private data belong on your server — your frontend should talk to *your* server, which then talks to the external API with the secret key.
:::

## Response headers

The server sends headers back too:

```js
const response = await fetch("/api/data");

response.headers.get("Content-Type");        // "application/json; charset=utf-8"
response.headers.get("X-RateLimit-Remaining"); // "98"
response.headers.get("Cache-Control");        // "max-age=300"
```

Common response headers:
- `Content-Type` — format of the response body
- `X-RateLimit-*` — how many requests you have left
- `Cache-Control` — how long the browser can cache the response
- `Set-Cookie` — sets a cookie (handled automatically by the browser)

## What CORS is and why it blocks you

**CORS** (Cross-Origin Resource Sharing) is a browser security mechanism. It prevents JavaScript on one origin (like `localhost:3000`) from reading responses from a different origin (like `api.example.com`) unless that server explicitly allows it.

An **origin** is the combination of protocol + domain + port:
- `http://localhost:3000` and `http://localhost:5000` are different origins
- `https://myapp.com` and `https://api.myapp.com` are different origins
- `https://myapp.com/page1` and `https://myapp.com/page2` are the same origin

When you fetch a different origin, the browser checks the response for an `Access-Control-Allow-Origin` header. If it is missing or does not match your origin, the browser blocks the response:

```
Access to fetch at 'https://api.example.com/data' from origin 'http://localhost:3000'
has been blocked by CORS policy
```

:::key
CORS is enforced by the **browser**, not the server. The request actually reaches the server — the server even sends back a response. But the browser refuses to let your JavaScript read that response unless the server explicitly permits your origin.
:::

## Preflight requests

For "complex" requests (anything that is not a simple GET with basic headers), the browser sends an automatic **preflight** request using the `OPTIONS` method before the real request. This asks the server: "Will you accept a POST with these headers from this origin?"

```
OPTIONS /api/data HTTP/1.1
Origin: http://localhost:3000
Access-Control-Request-Method: POST
Access-Control-Request-Headers: Content-Type, Authorization
```

The server must respond with:

```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: POST, GET, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
```

Only then does the browser send the real request.

Simple requests (GET with no custom headers) skip the preflight — the browser sends them directly and checks the response headers.

## How to fix CORS

CORS is a **server-side** setting. You cannot fix it from frontend JavaScript. The solutions:

1. **Configure the server** to send the correct `Access-Control-Allow-Origin` header:

```
Access-Control-Allow-Origin: https://myapp.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
```

2. **Use a proxy** during development. Your dev server (Vite, Webpack) can proxy `/api` requests to the real server, avoiding CORS entirely since the browser thinks it is talking to the same origin.

3. **Build your own backend** that calls the external API. Your frontend talks to your server (same origin, no CORS), and your server talks to the external API (server-to-server, no CORS).

:::tip
If you see a CORS error, do not try to hack around it with browser extensions or disabling security. The fix belongs on the server. If you do not control the server, use option 2 or 3.
:::

## API keys — a security note

Many APIs require a key for identification and rate limiting:

```js
// Public API key — designed for browser use (limited permissions)
const response = await fetch(`https://api.weather.com/forecast?key=PUBLIC_KEY&city=London`);

// Secret API key — NEVER put this in frontend code
// This belongs on your server
const response = await fetch("https://api.stripe.com/charges", {
  headers: { "Authorization": "Bearer sk_live_SECRET_KEY" }, // WRONG in frontend!
});
```

The distinction:
- **Public keys** — rate-limited, read-only, designed for browser use. OK in frontend code.
- **Secret keys** — full access, can charge money, read private data. Must stay on your server.

:::warning
A common beginner mistake: putting a secret API key in a `.js` file that ships to the browser. Even if you "hide" it in a variable, anyone can see it in the browser's DevTools Network tab or by reading your JavaScript source.
:::

## Reading API documentation

Every API documents:
1. **Base URL** — the root of all endpoints (`https://api.example.com/v1`)
2. **Endpoints** — specific URLs for each resource (`GET /users`, `POST /posts`)
3. **Required headers** — authentication, content type
4. **Request body format** — what fields to send, their types
5. **Response format** — what the response JSON looks like
6. **Error codes** — what specific status codes mean

Learning to read API docs quickly is one of the most practical skills in web development.

## Practice

:::quiz
Q: You send a POST request to `https://api.example.com` from your app at `http://localhost:3000` and get a CORS error. Where is the fix?
- In your frontend JavaScript
- In your browser settings
- On the server at api.example.com (or use a proxy) *
- In the HTML file
E: CORS is enforced by the browser but configured on the server. The server must send `Access-Control-Allow-Origin` headers that include your origin. If you do not control the server, use a development proxy or your own backend.
:::

:::quiz
Q: What is the difference between PUT and PATCH?
- PUT is for creating, PATCH is for deleting
- PUT replaces the entire resource, PATCH updates only the fields you send *
- They are the same but with different names
- PATCH is deprecated
E: PUT expects you to send the complete resource — omitted fields may be removed. PATCH sends only the fields you want to change. Check the API docs to know which one to use.
:::

## Recap

- **HTTP methods** describe the action: GET (read), POST (create), PUT (replace), PATCH (update), DELETE (remove).
- **Content-Type** tells the server the body format; **Authorization** sends credentials; **Accept** requests a response format.
- **Never put secret API keys in frontend code.** Use public keys for browser requests and keep secrets on your server.
- **CORS** is a browser security rule — the server must send `Access-Control-Allow-Origin` to permit cross-origin requests.
- **Preflight requests** (`OPTIONS`) are sent automatically for complex requests; the server must respond with allowed methods and headers.
- Fix CORS on the **server side**, or use a **dev proxy** / your own backend as a relay.
- Read **API documentation** to find endpoints, required headers, body formats, and error codes.

**Next up:** Loading, Error, and Empty States — building UIs that handle every phase of an async operation.
