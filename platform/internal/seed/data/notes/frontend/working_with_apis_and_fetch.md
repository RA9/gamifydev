# Working with APIs and Fetch

Most interesting apps don't keep all their data in the browser — they ask a server for it. This lesson shows you how to request data over the internet with `fetch`, turn the response into something usable, and put it on the page.

## What Is an API?

An **API** (Application Programming Interface) is, for our purposes, *a server you can ask for data*. You send a request to a web address (a URL), and the server sends back a response. You don't care how the server works inside — you just know the "menu" of things you can ask for.

:::analogy
An API is like a restaurant kitchen. You (the browser) don't walk into the kitchen and cook. You read the menu, hand the waiter an order (a request), and a plate comes back (a response). The API is the menu plus the waiter — a defined way to ask for things and get them back.
:::

Most web APIs speak over HTTP, the same protocol your browser uses to load pages. The data they send back is almost always in a format called JSON.

## JSON: Data as Text

**JSON** (JavaScript Object Notation) is a way to write objects and arrays as plain text so they can travel across the internet. It looks almost exactly like JavaScript objects and arrays:

```json
{
  "id": 1,
  "name": "Leanne Graham",
  "email": "leanne@example.com",
  "hobbies": ["reading", "coding"]
}
```

The catch: JSON is *text*, not a live JavaScript object. You can't use the dots and brackets on a string. Two built-in functions convert between them:

```js
// JSON.parse: text  ->  real JavaScript object
const text = '{"name": "Ada", "age": 36}';
const obj = JSON.parse(text);
console.log(obj.name); // "Ada"

// JSON.stringify: real object  ->  text (to send or store)
const user = { name: "Ada", age: 36 };
const out = JSON.stringify(user);
console.log(out); // '{"name":"Ada","age":36}'
```

:::warning
JSON is strict: keys must be in double quotes, no trailing commas, no comments. `{name: 'Ada'}` is valid JavaScript but **invalid** JSON. The error you'll see is usually "Unexpected token."
:::

You'll rarely call `JSON.parse` by hand for fetch results, though — `fetch` gives you a shortcut, as you'll see next.

## fetch Returns a Promise

`fetch(url)` starts a network request and returns a **promise** (covered in the Async JavaScript lesson). The data isn't ready immediately — the promise settles when the server responds.

```js
const promise = fetch("https://jsonplaceholder.typicode.com/users");
console.log(promise); // Promise { <pending> }  -- not the data yet!
```

Because it's a promise, you use `await` (inside an `async` function) to pause until the response arrives.

## Reading JSON with response.json()

What `fetch` resolves to is a **Response** object — think of it as the sealed envelope, not the letter inside. To get the actual data, call `response.json()`, which *also* returns a promise, so you `await` it too:

```js
async function loadUsers() {
  const response = await fetch("https://jsonplaceholder.typicode.com/users");
  const users = await response.json(); // parse the JSON body into real objects
  console.log(users[0].name); // "Leanne Graham"
}

loadUsers();
```

:::warning
You need *two* awaits: one for `fetch` (wait for the server), one for `.json()` (wait for the body to be read and parsed). Forgetting the second is a classic mistake — `users` would be a Promise, not your array.
:::

## A Complete Worked Example

Let's fetch a list of users and render them into the page. Assume this HTML exists:

```html
<ul id="user-list"></ul>
```

And this JavaScript:

```js
async function renderUsers() {
  const list = document.getElementById("user-list");
  const response = await fetch("https://jsonplaceholder.typicode.com/users");
  const users = await response.json();

  // Build a list item for each user
  for (const user of users) {
    const li = document.createElement("li");
    li.textContent = `${user.name} (${user.email})`;
    list.appendChild(li);
  }
}

renderUsers();
// The <ul> fills with: Leanne Graham (Sincere@april.biz), etc.
```

This works on a good day. But networks fail, servers go down, and URLs get typos. A real app must handle that.

## Checking response.ok and Catching Errors

Here's a surprise: `fetch` does **not** reject for HTTP errors like 404 (not found) or 500 (server error). The promise still fulfills — you just got an error *response*. `fetch` only rejects when the request couldn't be made at all (no internet, bad domain).

So you must check `response.ok` (true for status codes 200–299) yourself, and wrap everything in `try/catch` for network failures:

```js
async function loadUsers() {
  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users");

    if (!response.ok) {
      // e.g. 404 or 500 — throw so the catch block handles it
      throw new Error(`Server responded with ${response.status}`);
    }

    const users = await response.json();
    return users;
  } catch (error) {
    // network down, bad URL, or the thrown error above
    console.log("Could not load users:", error.message);
    return [];
  }
}
```

:::key
`fetch` rejects only on network failure. For 404/500 you must check `response.ok` and throw yourself. Always combine `response.ok` with `try/catch` to cover both cases.
:::

## Loading and Error States

Users hate staring at a blank screen wondering if anything is happening. Show a **loading state** while you wait, and an **error state** if it fails. Here's the full pattern:

```html
<div id="status"></div>
<ul id="user-list"></ul>
```

```js
async function renderUsers() {
  const status = document.getElementById("status");
  const list = document.getElementById("user-list");

  status.textContent = "Loading users...";   // 1. loading state
  list.innerHTML = "";                        // clear any old results

  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users");
    if (!response.ok) {
      throw new Error(`Status ${response.status}`);
    }
    const users = await response.json();

    status.textContent = "";                  // 2. success: clear status
    for (const user of users) {
      const li = document.createElement("li");
      li.textContent = `${user.name} (${user.email})`;
      list.appendChild(li);
    }
  } catch (error) {
    status.textContent = "Sorry, we couldn't load users. Try again."; // 3. error state
    console.log(error.message);
  }
}

renderUsers();
```

Three clear phases: loading, success, error. Every fetch you write in a real app should account for all three.

:::tip
Set the loading message *before* the `await`, and clear or replace it in both the success path and the `catch`. That way the user always knows what's going on.
:::

## Query Parameters

Often you don't want *all* the data — you want a filtered slice. APIs accept **query parameters** in the URL: a `?`, then `key=value` pairs joined by `&`.

```js
// Get only the posts written by user 1:
const url = "https://jsonplaceholder.typicode.com/posts?userId=1";
const response = await fetch(url);
const posts = await response.json();
```

To build URLs safely (especially with user input that might contain spaces or symbols), use `URLSearchParams` instead of gluing strings together:

```js
const params = new URLSearchParams({ userId: 1, _limit: 5 });
const url = `https://jsonplaceholder.typicode.com/posts?${params}`;
// "https://jsonplaceholder.typicode.com/posts?userId=1&_limit=5"
const response = await fetch(url);
```

:::tip
`URLSearchParams` automatically *encodes* special characters. Searching for "ada lovelace" becomes `ada%20lovelace` for you — no broken URLs.
:::

## A Note on CORS and API Keys

Two things you'll bump into with real APIs:

**CORS** (Cross-Origin Resource Sharing) is a browser security rule. By default, a page can only fetch from its *own* domain. To allow other sites to call it, a server must send special permission headers. If you see a console error mentioning "CORS policy" or "Access-Control-Allow-Origin," it means the server hasn't allowed your page — and you usually can't fix that from the browser. It's the server's call.

**API keys** are secret passwords many APIs require to identify you and track usage. They're often sent in a header:

```js
const response = await fetch("https://api.example.com/data", {
  headers: { "Authorization": "Bearer YOUR_API_KEY" },
});
```

:::warning
Never put a secret API key in front-end JavaScript that ships to the browser — anyone can open dev tools and read it. Real apps keep secret keys on their *own* server and have the browser talk to that. For now, only use keys that are meant to be public.
:::

## Practice

:::quiz
Q: After `const response = await fetch(url)`, how do you get the actual JSON data?
- `response.data`
- `await response.json()` *
- `JSON.parse(response)`
E: A fetch resolves to a Response object (the envelope). Call `await response.json()` to read and parse the body into real JavaScript objects.
:::

:::predict
Q: A request returns HTTP 404 (not found). What does `fetch` do?
- Rejects the promise, jumping to catch
- Fulfills with a response where `response.ok` is false *
- Returns `null`
E: `fetch` only rejects on network failure. A 404 is a successful *response* with `ok === false`, so you must check `response.ok` yourself and throw if needed.
:::

:::predict
Q: You want only the first 5 posts. Which URL is correct?
- `https://.../posts/5`
- `https://.../posts?_limit=5` *
- `https://.../posts&limit=5`
E: Query parameters start with `?` and use `key=value`. `posts/5` would fetch the single post with id 5, and `&` before the first param is wrong (it joins params, the first uses `?`).
:::

## Recap

- An API is a server you ask for data; web APIs usually send back JSON.
- JSON is objects/arrays written as text; `JSON.parse` reads it, `JSON.stringify` writes it.
- `fetch(url)` returns a promise that resolves to a Response; `await response.json()` gets the data (two awaits total).
- `fetch` rejects only on network failure — check `response.ok` for 404/500 and wrap in `try/catch`.
- Show loading, success, and error states so users always know what's happening.
- Add query parameters with `?key=value`; build them safely with `URLSearchParams`.
- CORS is a browser rule controlled by the server; keep secret API keys off the front end.

**Next up:** You'll combine fetch with everything you've learned to build a small data-driven app from scratch.
