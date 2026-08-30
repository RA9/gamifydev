# Promises Deep Dive

You have seen promises used with `fetch`, but there is more to understand. This lesson takes you inside the Promise machinery: the three states, how to create your own promises, chaining, returning values between `.then` calls, and the mistakes that trip up everyone.

## The three states of a promise

A promise is always in one of three states:

1. **Pending** — the operation has not finished yet. No result.
2. **Fulfilled** — the operation succeeded. The promise has a **value**.
3. **Rejected** — the operation failed. The promise has a **reason** (an error).

A promise starts pending and then *settles* into either fulfilled or rejected, exactly once. Once settled, it never changes again.

```js
const p = fetch("https://jsonplaceholder.typicode.com/users/1");
console.log(p); // Promise { <pending> }

// later, when the server responds, p settles to fulfilled
// (or rejected, if the network fails)
```

:::key
A settled promise is locked. If it fulfilled with the value `42`, it will be `42` forever — no callback can change it. This immutability makes promises predictable.
:::

## Creating a promise with new Promise

Most of the time you use promises that other APIs give you (`fetch`, `setTimeout` wrappers, etc.). But you can create your own:

```js
const myPromise = new Promise((resolve, reject) => {
  const success = true;

  if (success) {
    resolve("It worked!"); // fulfills the promise with this value
  } else {
    reject(new Error("Something went wrong")); // rejects with this reason
  }
});
```

The function you pass to `new Promise` is called the **executor**. It receives two callbacks: `resolve` (fulfill) and `reject` (reject). Call exactly one of them, once.

### Wrapping setTimeout in a promise

```js
function wait(ms) {
  return new Promise(resolve => {
    setTimeout(resolve, ms);
  });
}

// Usage
wait(2000).then(() => console.log("2 seconds later"));
```

This is how you turn callback-based APIs into promise-based ones.

### Wrapping a real-world callback

```js
function loadImage(src) {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve(img);
    img.onerror = () => reject(new Error(`Failed to load: ${src}`));
    img.src = src;
  });
}

loadImage("/photo.jpg")
  .then(img => document.body.append(img))
  .catch(err => console.error(err.message));
```

## .then — reacting to success

`.then(callback)` registers a function to run when the promise fulfills. It receives the resolved value:

```js
fetch("https://jsonplaceholder.typicode.com/users/1")
  .then(response => response.json())
  .then(user => console.log(user.name));
```

## .catch — handling errors

`.catch(callback)` runs when the promise rejects (or when any `.then` above it throws):

```js
fetch("https://bad-url.example.com")
  .then(response => response.json())
  .catch(error => console.log("Failed:", error.message));
```

`.catch` at the end of a chain handles errors from *any* step above it.

## .finally — cleanup either way

`.finally(callback)` runs whether the promise fulfills or rejects. It does not receive the value or error — it is for cleanup:

```js
showSpinner();

fetch(url)
  .then(response => response.json())
  .then(data => renderData(data))
  .catch(error => showError(error))
  .finally(() => hideSpinner()); // always runs
```

## Chaining — the real power

Each `.then` returns a **new promise** that resolves with whatever the callback returns. This lets you chain steps in a flat sequence:

```js
fetch("https://jsonplaceholder.typicode.com/users/1")
  .then(response => response.json())       // returns parsed JSON
  .then(user => user.name)                  // returns just the name string
  .then(name => name.toUpperCase())         // returns "LEANNE GRAHAM"
  .then(upper => console.log(upper));       // logs it
```

:::key
Whatever you `return` from a `.then` callback becomes the value the *next* `.then` receives. If you return a promise, the chain waits for it to settle before continuing.
:::

### Returning a promise in the chain

```js
function getUserPosts(userId) {
  return fetch(`https://jsonplaceholder.typicode.com/users/${userId}`)
    .then(res => res.json())
    .then(user => {
      console.log(`Fetching posts for ${user.name}`);
      return fetch(`https://jsonplaceholder.typicode.com/posts?userId=${userId}`);
    })
    .then(res => res.json()); // this .then waits for the second fetch
}
```

The chain stays flat. Each `.then` either returns a plain value or another promise.

## Common mistakes

### Mistake 1: Forgetting to return

```js
// BROKEN — the second .then gets undefined
fetch(url)
  .then(response => {
    response.json(); // no return!
  })
  .then(data => {
    console.log(data); // undefined
  });

// FIXED
fetch(url)
  .then(response => {
    return response.json(); // return the promise
  })
  .then(data => {
    console.log(data); // actual data
  });
```

With arrow functions and implicit return, you avoid this:

```js
fetch(url)
  .then(response => response.json()) // implicit return
  .then(data => console.log(data));
```

:::warning
If your `.then` callback has curly braces `{}`, you must explicitly write `return`. Without braces, the arrow function returns the expression automatically. This is the most common promise bug.
:::

### Mistake 2: Nesting .then calls

```js
// BAD — callback hell returns in promise clothing
fetch(url)
  .then(response => {
    response.json().then(data => {
      processData(data).then(result => {
        console.log(result);
      });
    });
  });

// GOOD — flat chain
fetch(url)
  .then(response => response.json())
  .then(data => processData(data))
  .then(result => console.log(result));
```

### Mistake 3: Missing .catch

```js
// DANGEROUS — if the fetch fails, you get an unhandled rejection
fetch(url)
  .then(response => response.json())
  .then(data => render(data));

// SAFE — errors are handled
fetch(url)
  .then(response => response.json())
  .then(data => render(data))
  .catch(error => console.error("Failed:", error.message));
```

Always put a `.catch` at the end of your chain or use `try/catch` with `async/await`.

## Promise.resolve and Promise.reject

Shortcuts for creating already-settled promises:

```js
const fulfilled = Promise.resolve(42);
fulfilled.then(val => console.log(val)); // 42

const rejected = Promise.reject(new Error("Oops"));
rejected.catch(err => console.log(err.message)); // "Oops"
```

These are useful when a function needs to return a promise but you already have the value:

```js
function getUser(id) {
  if (cache[id]) return Promise.resolve(cache[id]); // no fetch needed
  return fetch(`/api/users/${id}`).then(res => res.json());
}
```

## Practice

:::quiz
Q: What does the next `.then` receive if the previous `.then` does not return anything?
- The original promise
- null
- undefined *
- An error
E: A function without a `return` statement returns `undefined` in JavaScript. The next `.then` in the chain receives `undefined` as its argument. Always return values you want to pass forward.
:::

:::quiz
Q: Where should `.catch` go in a promise chain?
- Before the first `.then`
- Between every `.then`
- At the end of the chain to catch errors from any step above *
- Only when using `async/await`
E: A `.catch` at the end of the chain handles rejections from any `.then` above it. This is the simplest and most common pattern — one `.catch` covers the whole sequence.
:::

## Recap

- A promise has three states: **pending**, **fulfilled**, or **rejected**. It settles once and never changes.
- Create promises with `new Promise((resolve, reject) => {...})` to wrap callback-based APIs.
- `.then` handles success; `.catch` handles errors; `.finally` runs either way.
- **Chaining** works because each `.then` returns a new promise. Return a value to pass it forward; return a promise to make the chain wait.
- Common mistakes: **forgetting to return**, **nesting .then instead of chaining**, and **missing .catch**.
- `Promise.resolve(val)` and `Promise.reject(err)` create already-settled promises.

**Next up:** Async Await — the modern syntax that makes promise code read like synchronous code.
