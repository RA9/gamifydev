# Arrow Functions and Callbacks

Arrow functions are the modern, compact syntax for writing functions in JavaScript. Callbacks are the pattern of passing a function *into* another function to be called later. Together they form the backbone of array methods, event handlers, and asynchronous code.

## Arrow Function Syntax

An arrow function replaces the `function` keyword with `=>`:

```js
// Regular function expression
const add = function (a, b) {
  return a + b;
};

// Arrow function
const addArrow = (a, b) => {
  return a + b;
};

console.log(addArrow(3, 4));  // → 7
```

### Shorthand Rules

**One parameter — parentheses are optional:**

```js
const double = n => {
  return n * 2;
};
```

**Zero or multiple parameters — parentheses required:**

```js
const greet = () => {
  return "Hello!";
};

const multiply = (a, b) => {
  return a * b;
};
```

## Implicit Return

When the body is a **single expression**, you can drop the braces and the `return` keyword. The expression's result is returned automatically:

```js
const double = n => n * 2;
const add = (a, b) => a + b;
const greet = name => `Hello, ${name}!`;

console.log(double(5));       // → 10
console.log(add(3, 4));       // → 7
console.log(greet("Ada"));    // → "Hello, Ada!"
```

:::key
Implicit return only works with a single expression — no braces. The moment you add `{ }`, you must write `return` explicitly.
:::

To return an object literal implicitly, wrap it in parentheses (otherwise JS thinks the braces are a function body):

```js
// ❌ Bug — JS thinks { } is a block, not an object
const makeUser = name => { name: name, role: "viewer" };

// ✅ Correct — parentheses tell JS it's an expression
const makeUser = name => ({ name: name, role: "viewer" });

console.log(makeUser("Ada"));  // → { name: "Ada", role: "viewer" }
```

## When to Use Arrows vs Declarations

| Use Case | Preferred Style |
|----------|----------------|
| Top-level named functions | `function` declaration |
| Callbacks and inline functions | Arrow function |
| Array method callbacks | Arrow function |
| Object methods that use `this` | `function` declaration or method shorthand |

```js
// Top-level utility — declaration is fine
function calculateTax(amount, rate) {
  return amount * rate;
}

// Callback — arrow is cleaner
const prices = [10, 20, 30];
const withTax = prices.map(price => price * 1.08);
```

:::warning
Arrow functions do **not** have their own `this`. They inherit `this` from the surrounding scope. That's great for callbacks inside methods, but means you should avoid arrows as object methods when you need `this` to refer to the object.
:::

```js
const player = {
  name: "Ada",

  // ❌ Arrow — `this` is NOT the player object
  greetArrow: () => {
    console.log(`Hi, I'm ${this.name}`);  // `this` is the outer scope
  },

  // ✅ Method shorthand — `this` IS the player object
  greet() {
    console.log(`Hi, I'm ${this.name}`);
  },
};

player.greet();       // → "Hi, I'm Ada"
player.greetArrow();  // → "Hi, I'm undefined"
```

## Callbacks: Passing Functions as Arguments

A **callback** is a function you pass to another function, which calls it at the right time. This pattern is everywhere in JavaScript.

```js
function doTwice(action) {
  action();
  action();
}

doTwice(() => console.log("Hello!"));
// → Hello!
// → Hello!
```

The function receiving the callback decides *when* and *how* to call it. You just provide the "what."

### setTimeout and setInterval

`setTimeout` calls a function once after a delay (in milliseconds):

```js
console.log("Start");

setTimeout(() => {
  console.log("Delayed (2 seconds later)");
}, 2000);

console.log("End");
// → Start → End → Delayed (2 seconds later)
```

Notice that "End" prints before "Delayed." JavaScript doesn't pause — it schedules the callback and keeps going.

`setInterval` calls a function repeatedly at a fixed interval:

```js
let count = 0;

const timer = setInterval(() => {
  count++;
  console.log(`Tick ${count}`);
  if (count >= 3) {
    clearInterval(timer);  // stop after 3 ticks
  }
}, 1000);
// → Tick 1 (after 1s) → Tick 2 (after 2s) → Tick 3 (after 3s)
```

:::tip
Always store the return value of `setInterval` so you can call `clearInterval` to stop it. Forgetting this causes memory leaks and unexpected behavior.
:::

## Array Method Callbacks

Array methods are the most common place you'll write callbacks. Each method takes a function that runs once per item.

### forEach — Do Something with Each Item

```js
let names = ["Ada", "Bob", "Cal"];

names.forEach((name, index) => {
  console.log(`${index + 1}. ${name}`);
});
// → 1. Ada, 2. Bob, 3. Cal
```

`forEach` returns nothing — use it for side effects (logging, DOM updates), not for producing a new array.

### map — Transform Each Item into a New Array

```js
let prices = [10, 20, 30];
let discounted = prices.map(price => price * 0.9);

console.log(discounted);  // → [9, 18, 27]
console.log(prices);      // → [10, 20, 30]  (original unchanged)
```

The callback receives each item and returns the transformed version:

```js
let users = [
  { name: "Ada", age: 28 },
  { name: "Bob", age: 34 },
];

let names = users.map(user => user.name);
console.log(names);  // → ["Ada", "Bob"]
```

### filter — Keep Items that Pass a Test

```js
let numbers = [1, 2, 3, 4, 5, 6, 7, 8];

let evens = numbers.filter(n => n % 2 === 0);
console.log(evens);  // → [2, 4, 6, 8]

let adults = users.filter(user => user.age >= 18);
```

The callback returns `true` to keep an item, `false` to exclude it.

### Chaining Methods

Because `map` and `filter` return new arrays, you can chain them:

```js
let products = [
  { name: "Pen", price: 2, inStock: true },
  { name: "Notebook", price: 5, inStock: false },
  { name: "Backpack", price: 30, inStock: true },
];

let availableNames = products
  .filter(p => p.inStock)
  .map(p => p.name);

console.log(availableNames);  // → ["Pen", "Backpack"]
```

:::key
`map` transforms, `filter` narrows, `forEach` performs side effects. Read chains left to right: filter first, then transform.
:::

## Higher-Order Functions

A **higher-order function** is any function that takes a function as an argument or returns a function. `map`, `filter`, `forEach`, `setTimeout` — they're all higher-order functions.

You can write your own:

```js
function repeat(times, action) {
  for (let i = 0; i < times; i++) {
    action(i);
  }
}

repeat(3, i => console.log(`Iteration ${i}`));
// → Iteration 0, Iteration 1, Iteration 2
```

Returning a function:

```js
function createGreeter(greeting) {
  return name => `${greeting}, ${name}!`;
}

const sayHello = createGreeter("Hello");
const sayHey = createGreeter("Hey");

console.log(sayHello("Ada"));  // → "Hello, Ada!"
console.log(sayHey("Bob"));    // → "Hey, Bob!"
```

:::tip
If you understand that functions are values, higher-order functions are just functions that work with other functions — nothing magical. This is the foundation for advanced patterns like middleware, decorators, and composition.
:::

## The Callback Parameters You Get for Free

Array method callbacks receive up to three arguments: the **item**, the **index**, and the **full array**:

```js
["a", "b", "c"].forEach((item, index, array) => {
  console.log(`${index}: ${item} (of ${array.length})`);
});
// → 0: a (of 3)
// → 1: b (of 3)
// → 2: c (of 3)
```

You only need to declare the ones you use. Most of the time, you only need the item.

:::quiz
Q: What does an arrow function with no braces do with the expression after `=>`?
- Logs it to the console
- Returns it implicitly *
- Assigns it to a variable
E: A single expression without braces is an implicit return. `n => n * 2` automatically returns `n * 2`.
:::

:::quiz
Q: Why does `setTimeout` not pause the rest of your code?
- It runs on a separate thread
- JavaScript schedules the callback and continues executing — the callback runs later via the event loop *
- It's a bug in the language
E: JavaScript is single-threaded but non-blocking. `setTimeout` schedules a callback to run after a delay, but the main thread continues immediately.
:::

:::quiz
Q: What is a higher-order function?
- A function defined inside another function
- A function that takes a function as an argument or returns one *
- A function with more than two parameters
E: Higher-order functions operate on other functions — either by accepting them as arguments (like `map` and `filter`) or by returning them (like factory functions).
:::

## Recap

- Arrow functions (`=>`) are a shorter syntax: `const fn = (a, b) => a + b`.
- Single expressions can use implicit return (no braces, no `return` keyword).
- Arrows don't have their own `this` — don't use them as object methods.
- Callbacks are functions passed as arguments, called later by the receiving function.
- `setTimeout` and `setInterval` are timer-based callbacks — always clean up intervals.
- Array methods (`forEach`, `map`, `filter`) take callbacks that run per item. Chain `map` and `filter` for clean data pipelines.
- Higher-order functions take or return other functions — `map`, `filter`, `setTimeout`, and custom factory functions are all examples.

**Next up:** Objects and Object Patterns — structuring data with labeled properties and powerful access patterns.
