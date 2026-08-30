# JavaScript Patterns: Scope, Closures and Array Methods

At some point, frontend development stops being about knowing syntax and starts being about controlling complexity.

This lesson is about the JavaScript patterns that make real interfaces maintainable:

- understanding **where values live**
- understanding **what a function remembers**
- transforming arrays into UI-ready data without creating a mess

These ideas show up everywhere in professional frontend work.

## Why this lesson matters

Beginner code often works once and then becomes painful to change.

Professional code tends to have a few traits:

- variables exist in the smallest useful scope
- functions are predictable
- data is transformed cleanly before it reaches the UI
- state changes are explicit, not accidental

That is exactly what scope, closures, and array methods help with.

## Scope decides where a variable exists

A variable is only available inside the area where it was created.

The main scopes you need to understand are:

- **global scope**
- **function scope**
- **block scope**

```js
const appName = "Frontend Lab"; // global

function startSession() {
  const sessionId = 42; // function scope

  if (true) {
    const step = "intro"; // block scope
    console.log(appName, sessionId, step);
  }

  console.log(appName, sessionId);
  // console.log(step); // ❌ step does not exist here
}
```

The safest rule is simple:

:::key
Keep variables in the **smallest scope that still makes the code clear**.
:::

That reduces bugs, makes refactoring easier, and avoids different parts of the app accidentally depending on each other.

## `let` and `const` are block-scoped

This is one reason modern JavaScript is easier to trust than old `var`-heavy code.

```js
for (let i = 0; i < 3; i++) {
  console.log(i);
}

// console.log(i); // ❌ not available here
```

With `let` and `const`, loop counters and temporary values stay contained.

Use:

- `const` when the variable should not be reassigned
- `let` when the variable must change

## Closures: functions remember their surroundings

A **closure** means a function keeps access to variables from the scope where it was created, even after that outer code has finished running.

```js
function makeCounter() {
  let count = 0;

  return function () {
    count += 1;
    return count;
  };
}

const nextCount = makeCounter();

console.log(nextCount()); // 1
console.log(nextCount()); // 2
console.log(nextCount()); // 3
```

Why does this work?

Because the returned function closes over `count`. That variable is not copied. It is remembered.

Closures are extremely useful in frontend work for:

- private internal state
- event handlers
- factory functions
- configuration helpers

## A practical closure example

Imagine you want buttons that log analytics with different labels.

```js
function makeTracker(label) {
  return function () {
    console.log(`Tracked: ${label}`);
  };
}

const trackSignup = makeTracker("signup");
const trackPurchase = makeTracker("purchase");

trackSignup();   // Tracked: signup
trackPurchase(); // Tracked: purchase
```

Each returned function keeps its own remembered value.

That pattern is common in UI code when you need behavior that is partly shared and partly customized.

## Array methods are how frontend data gets shaped

Most interfaces render lists:

- tasks
- products
- messages
- menu items
- search results

That means strong frontend developers spend a lot of time transforming arrays.

The core methods are:

- `map`
- `filter`
- `find`
- `some`
- `every`
- `reduce`

## `map` transforms items

Use `map` when you want the same number of items back, but in a new shape.

```js
const tasks = [
  { title: "Ship landing page", done: true },
  { title: "Fix mobile nav", done: false },
];

const labels = tasks.map((task) => task.title);
console.log(labels); // ["Ship landing page", "Fix mobile nav"]
```

A very common frontend use is to transform server data into display data.

```js
const cards = tasks.map((task) => ({
  label: task.title,
  status: task.done ? "Complete" : "Open",
}));
```

## `filter` narrows a list

Use `filter` when you want only the items that match a rule.

```js
const openTasks = tasks.filter((task) => !task.done);
```

This is the core of:

- search UIs
- category filters
- tabs like all / active / completed
- permission-based lists

## `find`, `some`, and `every`

These are small but powerful.

```js
const currentTask = tasks.find((task) => task.title.includes("mobile"));
const hasCompleted = tasks.some((task) => task.done);
const allCompleted = tasks.every((task) => task.done);
```

Use them when the UI needs a quick answer, not a whole new array.

## `reduce` combines many values into one

`reduce` is useful for totals, grouped data, and summary metrics.

```js
const cart = [
  { name: "Mouse", price: 25 },
  { name: "Keyboard", price: 60 },
  { name: "Monitor", price: 180 },
];

const total = cart.reduce((sum, item) => sum + item.price, 0);
console.log(total); // 265
```

A frontend dashboard might use `reduce` for:

- counts by status
- total revenue
- average rating
- grouped sections

## Destructuring keeps data access readable

When objects get larger, destructuring helps pull out the parts you need.

```js
const user = {
  name: "Ada",
  role: "Frontend Developer",
  location: "Remote",
};

const { name, role } = user;
console.log(name, role);
```

It also makes function parameters clearer.

```js
function renderProfile({ name, role }) {
  return `${name} — ${role}`;
}
```

## Spread syntax helps you update data without mutation

As projects grow, mutating data in place becomes risky because different parts of the app may share references to the same object.

```js
const state = {
  filter: "all",
  theme: "light",
};

const nextState = {
  ...state,
  theme: "dark",
};
```

For arrays:

```js
const tasks = ["Plan", "Build"];
const nextTasks = [...tasks, "Test"];
```

This pattern is valuable because it makes state changes easier to reason about.

:::tip
If a UI suddenly updates in strange ways, check whether you accidentally mutated an object or array that multiple parts of the app are sharing.
:::

## A practical frontend pipeline

Here is a realistic pattern:

```js
const visibleCards = products
  .filter((product) => product.inStock)
  .filter((product) => product.name.toLowerCase().includes(searchTerm))
  .map((product) => ({
    title: product.name,
    priceLabel: `$${product.price}`,
  }));
```

That is professional frontend JavaScript:

1. start with raw data
2. filter it based on state
3. transform it into exactly what the UI needs

## Common mistakes to avoid

### 1. Too much global state

```js
let currentUser;
let currentTheme;
let currentPage;
let currentFilter;
```

This gets hard to manage fast.

Prefer grouping related state into objects and passing values intentionally.

### 2. Mutating arrays by accident

Methods like `push`, `pop`, and `sort` change the original array.

```js
const nums = [3, 1, 2];
nums.sort(); // mutates nums
```

If you want a safer pattern:

```js
const sorted = [...nums].sort();
```

### 3. Using loops when a data method says the intent better

A loop is not wrong, but sometimes `filter` or `map` makes the code much clearer.

## Practice mindset

When you build your next project, ask yourself:

- is this value global, local, or block-local?
- does this function need to remember something?
- am I transforming data clearly before rendering it?
- am I mutating data accidentally?

That shift in thinking is part of moving from beginner scripts to maintainable frontend architecture.

:::quiz
Q: Why are closures useful in frontend work?
- They make CSS load faster
- They let functions keep access to variables from where they were created *
- They prevent arrays from changing size
E: Closures are what allow functions to remember values from outer scopes, which is useful for handlers, factories, and small pieces of private state.
:::

:::quiz
Q: Which array method is the best fit when you want only products that match a search term?
- `map`
- `filter` *
- `reduce`
E: `filter` returns only the items that pass a rule, which is exactly what a search or category filter needs.
:::

## Recap

- Scope controls where variables are visible.
- `let` and `const` are block-scoped and safer than legacy `var`.
- Closures let functions remember outer variables.
- Array methods like `map`, `filter`, `find`, and `reduce` are the backbone of real UI data work.
- Destructuring and spread syntax make state updates and data access cleaner.
- Strong frontend code transforms data deliberately before rendering it.

**Next up:** you'll use these ideas directly against the DOM, where state and data transformations start turning into visible UI.