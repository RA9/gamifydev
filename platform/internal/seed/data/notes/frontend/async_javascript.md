# Async JavaScript

So far most of your code has run top to bottom, one line finishing before the next begins. But real apps need to *wait* — for a timer, for a file, for a server far away. This lesson is about how JavaScript handles waiting without freezing your whole page.

## Synchronous vs Asynchronous

**Synchronous** code runs in order, one step at a time. Each line blocks the next until it finishes.

```js
console.log("first");
console.log("second");
console.log("third");
// first
// second
// third
```

That is fine when every step is instant. The problem is *slow* steps. Imagine asking a server for data took 2 seconds and JavaScript just stopped everything until the answer came back:

```js
const data = downloadSomethingSlow(); // imagine this blocks for 2 seconds
console.log("Done!");
// During those 2 seconds, the page is FROZEN:
// no clicks, no scrolling, no typing. The user thinks it crashed.
```

That frozen feeling is why blocking is bad. **Asynchronous** code is the fix: you start a slow task, and JavaScript moves on to other work. When the slow task finishes later, your code gets notified and picks up where it left off.

:::analogy
Think of cooking. Synchronous cooking: you put pasta on, then *stare at the pot* until it boils, refusing to do anything else. Asynchronous cooking: you put the pasta on, set a timer, and go chop vegetables. When the timer beeps, you come back. Same kitchen, one cook, but you don't waste time staring.
:::

## A Mental Model: One Thing at a Time

Here's the key idea that trips up beginners: **JavaScript only does one thing at a time.** There is a single "lane" for running your code. So how can it wait *and* stay responsive?

The trick is that JavaScript hands off the waiting to the browser. When you start a timer or a network request, the browser holds onto it. JavaScript's single lane is free immediately. When the wait is over, the browser drops your follow-up code into a queue, and JavaScript runs it as soon as the lane is clear.

This handoff-and-queue system is called the **event loop**. You don't have to memorize its internals. Just hold this picture:

:::key
JavaScript does one thing at a time, but it can *hand off waiting* to the browser and keep going. When the wait finishes, your follow-up code is queued and runs as soon as the current code is done.
:::

## setTimeout: Your First Async Example

`setTimeout` runs a function *later*, after a delay in milliseconds. It is the simplest way to see async behavior.

```js
console.log("A");
setTimeout(function () {
  console.log("B");
}, 1000); // wait 1000ms = 1 second
console.log("C");
// A
// C
// B   <- one second later
```

Notice the order: `A`, then `C`, then `B`. Even though `setTimeout` appears in the middle, its function is handed off and runs *after* the rest of the current code finishes. This surprises almost everyone the first time.

A delay of `0` does **not** mean "right now" — it means "as soon as the current code is done":

```js
console.log("start");
setTimeout(() => console.log("timeout"), 0);
console.log("end");
// start
// end
// timeout
```

## Callbacks and "Callback Hell"

A **callback** is just a function you hand to another function to be called later. `setTimeout` takes a callback. This pattern works, but it gets ugly when steps depend on each other:

```js
setTimeout(() => {
  console.log("Step 1 done");
  setTimeout(() => {
    console.log("Step 2 done");
    setTimeout(() => {
      console.log("Step 3 done");
    }, 1000);
  }, 1000);
}, 1000);
```

See that staircase drifting to the right? Each step nests inside the last. This is nicknamed **callback hell** (or the "pyramid of doom"). It is hard to read, hard to add error handling to, and easy to break. Promises were invented to flatten this out.

## Promises

A **Promise** is an object that represents a value you don't have *yet* — it's an IOU for a future result. A promise is always in one of three states:

- **pending** — still waiting, no result yet.
- **fulfilled** — finished successfully, has a value.
- **rejected** — failed, has an error.

A promise starts pending and then settles into exactly one of the other two, once. You react to the result with `.then` (success), `.catch` (failure), and `.finally` (runs either way):

```js
const promise = fetch("https://jsonplaceholder.typicode.com/users");

promise
  .then((response) => response.json()) // runs when fulfilled
  .then((data) => console.log(data))   // chain another step
  .catch((error) => console.log("Something failed:", error))
  .finally(() => console.log("Done, success or not"));
```

Because each `.then` returns a new promise, you can *chain* steps in a flat line instead of nesting. That's the staircase from callback hell, straightened out.

:::tip
Whatever you `return` inside a `.then` becomes the value the *next* `.then` receives. Return a promise and the chain waits for it before continuing.
:::

## async / await: Nicer Syntax

Promises are good, but chaining `.then` still reads a bit backwards. `async`/`await` lets you write async code that *looks* synchronous — top to bottom — while still not blocking the page.

Two rules: mark a function `async`, and use `await` in front of a promise to pause until it settles and get its value.

```js
async function getUsers() {
  const response = await fetch("https://jsonplaceholder.typicode.com/users");
  const data = await response.json();
  console.log(data);
}

getUsers();
```

Compare that to the `.then` chain above — same work, but it reads like a normal recipe. `await` only pauses *inside* its async function; the rest of your page keeps running.

:::warning
You can only use `await` inside an `async` function (or at the top level of a module). Using `await` in a plain function is a syntax error.
:::

## Error Handling with try / catch

With `.then` you use `.catch`. With `async`/`await` you use the regular `try/catch` you already know:

```js
async function getUsers() {
  try {
    const response = await fetch("https://jsonplaceholder.typicode.com/users");
    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.log("Failed to load users:", error);
  }
}
```

If anything inside `try` rejects or throws, control jumps straight to `catch`. One block catches the whole sequence — much cleaner than wiring up error handling at every nested step.

## Sequence vs Parallel: Promise.all

When you `await` things one after another, they run **in sequence** — each waits for the one before it:

```js
async function loadSequential() {
  const a = await fetch("https://jsonplaceholder.typicode.com/users/1");
  const b = await fetch("https://jsonplaceholder.typicode.com/users/2");
  // total time = time(a) + time(b)
}
```

If `a` and `b` don't depend on each other, that's wasteful. Start them **in parallel** and wait for both with `Promise.all`:

```js
async function loadParallel() {
  const [resA, resB] = await Promise.all([
    fetch("https://jsonplaceholder.typicode.com/users/1"),
    fetch("https://jsonplaceholder.typicode.com/users/2"),
  ]);
  // total time = the SLOWER of the two, not the sum
  const userA = await resA.json();
  const userB = await resB.json();
}
```

`Promise.all` takes an array of promises and gives you back an array of results in the same order. Note: if *any* promise rejects, `Promise.all` rejects immediately.

:::tip
Rule of thumb: if step B needs step A's result, use sequential `await`. If they're independent, fire them together with `Promise.all` and wait once.
:::

## A Common Gotcha: Forgetting await

If you forget `await`, you get the *promise object* instead of the value inside it:

```js
async function broken() {
  const response = fetch("https://jsonplaceholder.typicode.com/users"); // no await!
  const data = await response.json(); // ERROR: response is a Promise, not a Response
}
```

The symptom is often a confusing error like "response.json is not a function," or logging `Promise { <pending> }` when you expected real data. When something async behaves weirdly, the first thing to check is: did I `await` it?

:::predict
Q: What does this print? `console.log(1); setTimeout(() => console.log(2), 0); console.log(3);`
- 1, 2, 3
- 1, 3, 2 *
- 2, 1, 3
E: `setTimeout` hands its callback off to be run *after* the current code finishes, even with a `0` delay. So `1` and `3` print first (synchronous), then `2` from the queued callback.
:::

:::predict
Q: Two independent fetches that each take 1 second. Which approach finishes in about 1 second total?
- Two `await` calls on separate lines, one after another
- `await Promise.all([...])` with both fetches *
- Wrapping both in setTimeout
E: Sequential `await` runs them back to back (about 2 seconds). `Promise.all` starts both at once and waits for the slower one, so the total is about 1 second.
:::

:::quiz
Q: Inside an `async` function, what is the cleanest way to handle a fetch that might fail?
- An `if` statement around the promise
- `try/catch` around the `await` calls *
- A second setTimeout
E: With `async`/`await`, errors from awaited promises can be caught with ordinary `try/catch` — the same syntax you use for synchronous errors.
:::

## Recap

- Synchronous code blocks; asynchronous code lets JavaScript start a slow task and keep going.
- JavaScript does one thing at a time but hands waiting off to the browser — the event loop queues your follow-up code.
- `setTimeout(fn, ms)` runs a callback later; even `0` ms runs *after* the current code.
- Nested callbacks become "callback hell"; Promises flatten this with `.then`/`.catch`/`.finally`.
- A promise is pending, then settles once into fulfilled or rejected.
- `async`/`await` makes promise code read top-to-bottom; use `try/catch` for errors.
- Sequential `await` runs steps in order; `Promise.all` runs independent tasks in parallel.
- The most common bug is forgetting `await` — you get a Promise instead of the value.

**Next up:** You'll put these skills to work talking to real servers in "Working with APIs and Fetch."
