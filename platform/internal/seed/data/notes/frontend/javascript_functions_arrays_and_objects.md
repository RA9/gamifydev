# JavaScript Functions, Arrays and Objects

Functions let you package code so you can reuse it; arrays and objects let you store and organize data. Together they are the heart of every real JavaScript program.

## Functions: Why They Exist

A function is a named, reusable block of code. Instead of copying the same lines over and over, you write them once and *call* the function whenever you need them. This is the **DRY** principle — Don't Repeat Yourself.

:::analogy
A function is like a coffee machine. You don't rebuild it each morning — you press a button (call it), optionally choose a size (an argument), and get coffee back (the return value).
:::

## Declaring and Calling Functions

```js
function greet() {
  console.log("Hello!");
}

greet();  // → Hello!
greet();  // → Hello!  (reuse it as many times as you like)
```

### Parameters and Arguments

A **parameter** is a placeholder named in the function. An **argument** is the actual value you pass when calling.

```js
function greet(name) {        // name is the parameter
  console.log(`Hello, ${name}!`);
}

greet("Ada");   // → Hello, Ada!   ("Ada" is the argument)
greet("Bob");   // → Hello, Bob!
```

Multiple parameters are separated by commas:

```js
function add(a, b) {
  console.log(a + b);
}
add(3, 4);  // → 7
```

### return

`console.log` only *prints*. To send a value *back* so the caller can use it, use `return`.

```js
function add(a, b) {
  return a + b;
}

let sum = add(3, 4);
console.log(sum);        // → 7
console.log(add(10, 5)); // → 15
```

Once `return` runs, the function stops immediately:

```js
function checkAge(age) {
  if (age < 18) {
    return "Too young";
  }
  return "Welcome";
}
console.log(checkAge(15));  // → Too young
console.log(checkAge(30));  // → Welcome
```

:::warning
A function with no `return` gives back `undefined`. If your `console.log(myFunc())` shows `undefined`, you probably forgot to `return`.
:::

### Default Parameters

Give a parameter a fallback for when no argument is passed.

```js
function greet(name = "friend") {
  return `Hi, ${name}!`;
}
console.log(greet("Ada"));  // → Hi, Ada!
console.log(greet());       // → Hi, friend!
```

### Arrow Functions

A shorter syntax for writing functions. These are equivalent:

```js
// Regular
function double(n) {
  return n * 2;
}

// Arrow
const double = (n) => {
  return n * 2;
};

// Arrow, short form: one expression auto-returns (no braces, no return)
const triple = (n) => n * 3;

console.log(double(5));  // → 10
console.log(triple(5));  // → 15
```

Arrow functions shine when passing a function to another function — you'll see this with array methods below.

## Scope (Briefly)

**Scope** is where a variable is visible. A variable declared *inside* a function only exists inside it (local). A variable declared outside any function is **global** and visible everywhere.

```js
let team = "Blue";       // global

function play() {
  let points = 10;       // local — only exists in here
  console.log(team);     // → Blue (can see global)
  console.log(points);   // → 10
}

play();
console.log(points);  // ❌ ReferenceError: points is not defined
```

`let` and `const` are also **block-scoped** — they only exist inside the `{ }` they're declared in (like an `if` or a `for`):

```js
if (true) {
  let secret = 42;
}
console.log(secret);  // ❌ not defined outside the block
```

## Arrays

An array is an ordered list of values, written with square brackets.

```js
let scores = [90, 85, 100];
let fruits = ["apple", "pear", "kiwi"];
let mixed  = [1, "two", true];
```

### Indexing and length

Items are numbered starting at **0**.

```js
let fruits = ["apple", "pear", "kiwi"];
console.log(fruits[0]);       // → apple
console.log(fruits[2]);       // → kiwi
console.log(fruits[10]);      // → undefined (no item there)
console.log(fruits.length);   // → 3

fruits[1] = "banana";         // change an item
console.log(fruits);          // → ["apple", "banana", "kiwi"]
```

:::warning
Off-by-one mistakes are common: a 3-item array has indexes 0, 1, 2 — the last index is `length - 1`, not `length`.
:::

### Adding and Removing Items

```js
let stack = ["a", "b"];

stack.push("c");      // add to the END
console.log(stack);   // → ["a", "b", "c"]

stack.pop();          // remove from the END (returns "c")
console.log(stack);   // → ["a", "b"]

stack.shift();        // remove from the START (returns "a")
console.log(stack);   // → ["b"]

stack.unshift("z");   // add to the START
console.log(stack);   // → ["z", "b"]
```

### Looping Over an Array

```js
let nums = [10, 20, 30];

// for...of — clean and readable
for (let n of nums) {
  console.log(n);  // → 10, 20, 30
}

// classic for — when you need the index
for (let i = 0; i < nums.length; i++) {
  console.log(i, nums[i]);  // → 0 10, then 1 20, then 2 30
}
```

### Key Array Methods

These are used constantly. Each takes a function that runs once per item.

**`forEach`** — do something with every item (no new array returned):

```js
let names = ["Ada", "Bob"];
names.forEach((name) => {
  console.log(`Hi ${name}`);  // → Hi Ada, then Hi Bob
});
```

**`map`** — make a *new* array by transforming each item:

```js
let nums = [1, 2, 3];
let doubled = nums.map((n) => n * 2);
console.log(doubled);  // → [2, 4, 6]
console.log(nums);     // → [1, 2, 3]  (original unchanged)
```

**`filter`** — make a *new* array of items that pass a test:

```js
let nums = [5, 12, 8, 20];
let big = nums.filter((n) => n > 10);
console.log(big);  // → [12, 20]
```

**`find`** — get the *first* item that passes a test (or `undefined`):

```js
let nums = [5, 12, 8, 20];
let firstBig = nums.find((n) => n > 10);
console.log(firstBig);  // → 12
```

**`includes`** — does the array contain a value? Returns a boolean:

```js
let fruits = ["apple", "pear"];
console.log(fruits.includes("pear"));   // → true
console.log(fruits.includes("mango"));  // → false
```

:::key
`map` and `filter` return **new** arrays and leave the original alone. `forEach` returns nothing — use it only for side effects like logging. Reach for `map` when you want a transformed copy.
:::

## Objects

An object groups related values under **labels** (keys). Use it when a list-by-position isn't enough and each piece of data deserves a name.

```js
let player = {
  name: "Ada",
  level: 5,
  isOnline: true,
};
```

### Dot vs Bracket Access

```js
console.log(player.name);      // → Ada   (dot notation — usual)
console.log(player["level"]);  // → 5     (bracket notation)

let key = "isOnline";
console.log(player[key]);      // → true  (bracket lets you use a variable)
```

Use dot notation normally. Use brackets when the key is stored in a variable or has spaces/special characters.

### Updating and Adding Properties

```js
let player = { name: "Ada", level: 5 };

player.level = 6;          // update existing
player.score = 100;        // add a new property
delete player.name;        // remove a property

console.log(player);       // → { level: 6, score: 100 }
```

### Nested Objects

Objects can hold other objects:

```js
let user = {
  name: "Ada",
  address: {
    city: "Lagos",
    zip: "10001",
  },
};

console.log(user.address.city);  // → Lagos
```

### Methods (Functions on Objects)

A property whose value is a function is called a **method**. Inside it, `this` refers to the object.

```js
let counter = {
  count: 0,
  increment() {
    this.count++;
  },
};

counter.increment();
counter.increment();
console.log(counter.count);  // → 2
```

### Object.keys and Object.values

Get an array of an object's keys or values — useful for looping.

```js
let scores = { math: 90, art: 80 };

console.log(Object.keys(scores));    // → ["math", "art"]
console.log(Object.values(scores));  // → [90, 80]

Object.keys(scores).forEach((subject) => {
  console.log(`${subject}: ${scores[subject]}`);
});
// → math: 90
// → art: 80
```

## Arrays of Objects: The Everyday Data Shape

Most real data is a *list of records* — and that's an array of objects. This is the single most common shape you'll work with (think: a list of users, products, or todos).

```js
let products = [
  { name: "Pen",      price: 2,  inStock: true  },
  { name: "Notebook", price: 5,  inStock: false },
  { name: "Backpack", price: 30, inStock: true  },
];
```

Loop over it to build output:

```js
products.forEach((p) => {
  console.log(`${p.name} costs $${p.price}`);
});
// → Pen costs $2
// → Notebook costs $5
// → Backpack costs $30
```

Combine with the methods you learned. Only the in-stock items:

```js
let available = products.filter((p) => p.inStock);
console.log(available.length);  // → 2
```

Just the names:

```js
let names = products.map((p) => p.name);
console.log(names);  // → ["Pen", "Notebook", "Backpack"]
```

Find one specific record:

```js
let bag = products.find((p) => p.name === "Backpack");
console.log(bag.price);  // → 30
```

:::tip
Chain methods together. To get the names of in-stock products: `products.filter(p => p.inStock).map(p => p.name)` → `["Pen", "Backpack"]`. Read it left to right: filter first, then transform.
:::

## Practice

:::predict
Q: What does this print? `function f(a, b) { return a + b; } console.log(f(2, 3));`
- undefined
- 5 *
- "23"
E: Both arguments are numbers, so `a + b` is `2 + 3 = 5`, and `return` sends it back to be logged.
:::

:::predict
Q: What does this print? `let a = [1, 2, 3]; let b = a.map(n => n + 1); console.log(b);`
- [1, 2, 3]
- [2, 3, 4] *
- [1, 2, 3, 4]
E: `map` builds a new array by adding 1 to each item, giving `[2, 3, 4]`. The original `a` is unchanged.
:::

:::predict
Q: What does this print? `let u = { name: "Ada" }; console.log(u.age);`
- error
- undefined *
- null
E: Reading a property that doesn't exist returns `undefined` — it does not throw an error.
:::

:::quiz
Q: Which method returns a NEW array containing only items that pass a test?
- forEach
- filter *
- push
E: `filter` keeps items where the test returns true and returns them in a new array.
:::

:::fill
Q: Array indexes start at the number `___`.
A: 0
E: The first item is at index 0, so the last index is `length - 1`.
:::

## Recap

- Functions package reusable code (DRY). Parameters are placeholders; arguments are the values you pass; `return` sends a value back.
- Default parameters supply fallbacks; arrow functions are a shorter syntax, great for passing into other functions.
- Variables are scoped: local to a function or block, or global; `let`/`const` are block-scoped.
- Arrays are ordered lists indexed from 0; use `push`/`pop`/`shift`/`unshift` to change them and `for...of` to loop.
- Master `forEach`, `map`, `filter`, `find`, and `includes` — `map`/`filter` return new arrays.
- Objects store labeled values; access with dot or bracket notation; methods are functions on objects; use `Object.keys`/`Object.values` to loop.
- Arrays of objects are the bread-and-butter data shape — loop, `filter`, `map`, and `find` over them to build output.

**Next up:** The DOM — using JavaScript to read and change the actual elements on a web page.
