# Async Await

Promise chains work, but they still read a bit backwards. `async/await` lets you write asynchronous code that reads top-to-bottom like regular synchronous code, while still being non-blocking under the hood.

## The async keyword

Mark a function `async` and it automatically returns a promise:

```js
async function greet() {
  return "Hello";
}

greet().then(msg => console.log(msg)); // "Hello"
```

Even though `greet` just returns a string, the `async` keyword wraps it in a promise. You never need to write `return Promise.resolve(...)` inside an async function.

## The await keyword

Inside an `async` function, `await` pauses execution until a promise settles, then gives you the resolved value:

```js
async function loadUser() {
  const response = await fetch("https://jsonplaceholder.typicode.com/users/1");
  const user = await response.json();
  console.log(user.name); // "Leanne Graham"
}

loadUser();
```

Two awaits, two pauses — but only inside this function. The rest of the page (clicks, animations, other scripts) keeps running.

:::key
`await` pauses the *function*, not the *browser*. JavaScript stays non-blocking. The page remains responsive while the function waits.
:::

## Converting a .then chain to async/await

The translation is mechanical:

```js
// .then chain
function getUser() {
  return fetch("https://jsonplaceholder.typicode.com/users/1")
    .then(response => response.json())
    .then(user => {
      console.log(user.name);
      return user;
    });
}

// async/await version
async function getUser() {
  const response = await fetch("https://jsonplaceholder.typicode.com/users/1");
  const user = await response.json();
  console.log(user.name);
  return user;
}
```

Same behavior, but the async version reads like a recipe: do this, then this, then this.

## Error handling with try/catch

With `.then` you use `.catch`. With `async/await` you use the same `try/catch` you already know from synchronous code:

```js
async function loadUser() {
  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users/1");

    if (!response.ok) {
      throw new Error(`Server error: ${response.status}`);
    }

    const user = await response.json();
    console.log(user.name);
  } catch (error) {
    console.error("Failed to load user:", error.message);
  }
}
```

If anything inside `try` throws or if an awaited promise rejects, execution jumps straight to `catch`. One block handles the whole sequence.

:::tip
Always check `response.ok` inside your try block. Remember: `fetch` does not reject on 404 or 500 — it only rejects on network failure. Throw manually for HTTP errors so they are caught by your `catch` block.
:::

## Async arrow functions

You can make arrow functions async too:

```js
const loadUser = async () => {
  const response = await fetch("https://jsonplaceholder.typicode.com/users/1");
  const user = await response.json();
  return user;
};
```

This works anywhere you use arrow functions — event handlers, array callbacks, or standalone functions.

### Async event handlers

```js
document.querySelector("#load-btn").addEventListener("click", async () => {
  const response = await fetch("/api/data");
  const data = await response.json();
  render(data);
});
```

## Top-level await

In JavaScript modules (`<script type="module">`), you can use `await` outside of any function:

```html
<script type="module">
  const response = await fetch("/api/config");
  const config = await response.json();
  console.log(config);
</script>
```

:::warning
Top-level `await` only works in modules, not in regular scripts. If you try it in a normal `<script>` tag, you get a syntax error. Add `type="module"` to your script tag to enable it.
:::

## Common patterns

### Pattern 1: Fetch and render

```js
async function renderUsers() {
  const container = document.querySelector("#users");
  container.textContent = "Loading...";

  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users");
    if (!response.ok) throw new Error(`Status: ${response.status}`);
    const users = await response.json();

    container.innerHTML = users
      .map(u => `<div class="card">${u.name}</div>`)
      .join("");
  } catch (error) {
    container.textContent = "Failed to load users. Try again.";
  }
}
```

### Pattern 2: Sequential operations

When step B depends on step A's result:

```js
async function loadUserAndPosts(userId) {
  const userRes = await fetch(`https://jsonplaceholder.typicode.com/users/${userId}`);
  const user = await userRes.json();

  const postsRes = await fetch(`https://jsonplaceholder.typicode.com/posts?userId=${userId}`);
  const posts = await postsRes.json();

  return { user, posts };
}
```

### Pattern 3: Parallel with await and Promise.all

When operations are independent:

```js
async function loadDashboard() {
  const [usersRes, postsRes, commentsRes] = await Promise.all([
    fetch("/api/users"),
    fetch("/api/posts"),
    fetch("/api/comments"),
  ]);

  const users = await usersRes.json();
  const posts = await postsRes.json();
  const comments = await commentsRes.json();

  return { users, posts, comments };
}
```

:::key
Use sequential `await` when step B needs step A's result. Use `Promise.all` with `await` when operations are independent — they run in parallel, cutting total time to the slowest request instead of the sum.
:::

## Common gotchas

### Forgetting await

```js
async function broken() {
  const response = fetch("/api/data"); // missing await!
  const data = await response.json();  // TypeError: response.json is not a function
}
```

Without `await`, `response` is a Promise object, not a Response. The symptom is usually a cryptic error or logging `Promise { <pending> }`.

### Using await in a non-async function

```js
// SyntaxError
function loadData() {
  const data = await fetch("/api/data"); // can't await here!
}

// Fix: make it async
async function loadData() {
  const data = await fetch("/api/data"); // works
}
```

### await in forEach (does not work as expected)

```js
// BROKEN — forEach does not wait for async callbacks
items.forEach(async (item) => {
  await processItem(item); // these all fire at once
});

// FIXED — use for...of for sequential processing
for (const item of items) {
  await processItem(item); // processes one at a time
}

// Or Promise.all for parallel
await Promise.all(items.map(item => processItem(item)));
```

:::warning
`forEach` ignores the promises returned by async callbacks — it fires all of them at once and moves on. Use `for...of` for sequential processing or `Promise.all(items.map(...))` for parallel.
:::

## Practice

:::quiz
Q: What does `await` actually do inside an async function?
- Blocks the entire browser until the promise settles
- Pauses only the async function until the promise settles, while the rest of the app continues *
- Converts the promise to a synchronous call
- Cancels the promise if it takes too long
E: `await` pauses execution of the containing async function and yields control back to the event loop. Other code, event handlers, and animations keep running. When the promise settles, the function resumes.
:::

:::quiz
Q: You need to fetch user data and post data, but posts do not depend on the user result. What is the fastest approach?
- `const user = await fetchUser(); const posts = await fetchPosts();`
- `const [user, posts] = await Promise.all([fetchUser(), fetchPosts()]);` *
- `fetchUser().then(() => fetchPosts())`
- `setTimeout(() => fetchPosts(), 0)`
E: `Promise.all` starts both fetches simultaneously and waits for both. Sequential `await` would fetch one, wait, then start the other — taking twice as long when they are independent.
:::

## Recap

- `async` makes a function return a promise; `await` pauses the function until a promise settles.
- `async/await` reads top-to-bottom — easier to follow than `.then` chains.
- Handle errors with **`try/catch`**, just like synchronous code.
- **Async arrow functions** work everywhere regular arrow functions do.
- **Top-level `await`** works in `<script type="module">` only.
- Use sequential `await` for dependent steps; **`Promise.all`** for parallel independent operations.
- Watch out for: **missing `await`**, using `await` in non-async functions, and `await` inside `forEach`.

**Next up:** Error Handling in Async Code — building resilient applications that handle failure gracefully.
