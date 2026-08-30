# Control Flow: Loops

Loops let you repeat code without writing it out by hand. Need to process 1,000 records? Sum a list of prices? Search for a name? A loop handles it in a few lines.

## The for Loop

The classic workhorse. Three parts separated by semicolons: **initialization**, **condition**, and **update**.

```js
for (let i = 0; i < 5; i++) {
  console.log(i);
}
// → 0, 1, 2, 3, 4
```

How it works, step by step:

1. `let i = 0` — runs once at the start
2. `i < 5` — checked before each iteration; if `false`, the loop stops
3. The loop body runs
4. `i++` — runs after each iteration
5. Back to step 2

### Counting in Different Ways

```js
// Count down
for (let i = 5; i > 0; i--) {
  console.log(i);  // → 5, 4, 3, 2, 1
}

// Count by twos
for (let i = 0; i <= 10; i += 2) {
  console.log(i);  // → 0, 2, 4, 6, 8, 10
}

// Loop through an array by index
let colors = ["red", "green", "blue"];
for (let i = 0; i < colors.length; i++) {
  console.log(`${i}: ${colors[i]}`);
}
// → 0: red, 1: green, 2: blue
```

:::warning
The condition `i <= colors.length` is a classic off-by-one bug. An array of 3 items has indexes 0, 1, 2 — so you want `i < colors.length`, not `<=`.
:::

## while Loop

Repeats *while* a condition remains `true`. Use it when you don't know ahead of time how many iterations you need.

```js
let attempts = 0;

while (attempts < 3) {
  console.log(`Attempt ${attempts + 1}`);
  attempts++;
}
// → Attempt 1, Attempt 2, Attempt 3
```

A practical example — keep halving until a number drops below 1:

```js
let value = 100;

while (value >= 1) {
  console.log(value);
  value = value / 2;
}
// → 100, 50, 25, 12.5, 6.25, 3.125, 1.5625, 1
```

## do...while Loop

Like `while`, but the body runs **at least once** because the condition is checked *after* the first iteration:

```js
let input;

do {
  input = prompt("Enter 'yes' to continue:");
} while (input !== "yes");
```

:::tip
`do...while` is uncommon in modern JS, but it's the right tool when you need at least one execution before checking a condition — like prompting for user input.
:::

## for...of (Arrays and Strings)

The cleanest way to iterate over the *values* of an iterable (arrays, strings, Maps, Sets):

```js
let fruits = ["apple", "pear", "kiwi"];

for (let fruit of fruits) {
  console.log(fruit);
}
// → apple, pear, kiwi
```

Works with strings too — each character is a value:

```js
for (let char of "Hello") {
  console.log(char);
}
// → H, e, l, l, o
```

:::key
Use `for...of` when you care about the **values**. Use a classic `for` loop when you also need the **index**. If you need both, use `forEach` or `entries()`.
:::

```js
// When you need both index and value
let tasks = ["Plan", "Build", "Test"];

for (let [index, task] of tasks.entries()) {
  console.log(`${index + 1}. ${task}`);
}
// → 1. Plan, 2. Build, 3. Test
```

## for...in (Objects)

`for...in` iterates over the **keys** of an object:

```js
let player = { name: "Ada", level: 5, score: 120 };

for (let key in player) {
  console.log(`${key}: ${player[key]}`);
}
// → name: Ada, level: 5, score: 120
```

:::warning
Don't use `for...in` on arrays. It iterates over keys (which are string indexes), can include inherited properties, and doesn't guarantee order in all edge cases. Use `for...of` for arrays.
:::

```js
// ❌ Avoid this
let nums = [10, 20, 30];
for (let i in nums) {
  console.log(typeof i);  // → "string" — not a number!
}

// ✅ Do this instead
for (let num of nums) {
  console.log(num);  // → 10, 20, 30
}
```

## break and continue

### break — Exit the Loop Entirely

```js
let numbers = [3, 7, 2, 9, 1, 5];

for (let num of numbers) {
  if (num === 9) {
    console.log("Found 9!");
    break;  // stop looping
  }
  console.log(`Checking ${num}...`);
}
// → Checking 3..., Checking 7..., Checking 2..., Found 9!
```

### continue — Skip to the Next Iteration

```js
for (let i = 0; i < 6; i++) {
  if (i % 2 !== 0) continue;  // skip odd numbers
  console.log(i);
}
// → 0, 2, 4
```

:::tip
`break` is perfect for search loops — stop as soon as you find what you want. `continue` is useful for skipping invalid items without nesting everything in an `if`.
:::

## Infinite Loop Prevention

An infinite loop freezes the browser tab (or the Node process). The fix is simple: **always make sure the condition will eventually become `false`**.

```js
// ❌ Infinite loop — i never changes
for (let i = 0; i < 10; ) {
  console.log(i);
  // forgot i++!
}

// ❌ Infinite loop — condition is always true
while (true) {
  console.log("stuck");
  // no break condition
}

// ✅ Safe — clear exit condition
let count = 0;
while (count < 100) {
  count++;
}
```

:::warning
If your page freezes and the tab becomes unresponsive, you likely have an infinite loop. Close the tab, find the loop, and check its exit condition.
:::

## Choosing the Right Loop

| Loop | Best For |
|------|----------|
| `for` | Known iteration count, need the index |
| `while` | Unknown count, condition-based termination |
| `do...while` | Must execute at least once |
| `for...of` | Iterating array/string values |
| `for...in` | Iterating object keys |

In practice, `for...of` and the classic `for` cover the vast majority of cases. Array methods like `forEach`, `map`, and `filter` replace many loops entirely.

## Practical Patterns

### Summing Values

```js
let prices = [12.99, 5.50, 8.75, 3.25];
let total = 0;

for (let price of prices) {
  total += price;
}
console.log(total.toFixed(2));  // → "30.49"
```

### Searching for an Item

```js
let users = ["Ada", "Bob", "Cal", "Diana"];
let target = "Cal";
let found = false;

for (let user of users) {
  if (user === target) {
    found = true;
    break;
  }
}
console.log(found);  // → true
```

### Building a String

```js
let words = ["JavaScript", "is", "awesome"];
let sentence = "";

for (let word of words) {
  sentence += word + " ";
}
console.log(sentence.trim());  // → "JavaScript is awesome"
```

### Filtering Without filter()

```js
let numbers = [1, 2, 3, 4, 5, 6, 7, 8];
let evens = [];

for (let n of numbers) {
  if (n % 2 === 0) {
    evens.push(n);
  }
}
console.log(evens);  // → [2, 4, 6, 8]
```

### Nested Loops

```js
let matrix = [
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9],
];

for (let row of matrix) {
  for (let cell of row) {
    process(cell);
  }
}
```

:::warning
Nested loops multiply iterations. A loop inside a loop over 1,000 items each runs 1,000,000 times. Be mindful of performance with large data sets.
:::

:::quiz
Q: What is the key difference between `for...of` and `for...in`?
- `for...of` is faster
- `for...of` iterates values, `for...in` iterates keys *
- `for...in` only works with arrays
E: `for...of` gives you the values of an iterable (array items, string characters). `for...in` gives you the keys (property names) of an object.
:::

:::quiz
Q: What does `continue` do inside a loop?
- Exits the loop entirely
- Skips the rest of the current iteration and moves to the next one *
- Restarts the loop from the beginning
E: `continue` jumps to the next iteration, skipping any remaining code in the current loop body. `break` is what exits entirely.
:::

:::quiz
Q: Which loop guarantees the body runs at least once?
- `for`
- `while`
- `do...while` *
E: `do...while` checks its condition after the body runs, so the body always executes at least once.
:::

## Recap

- `for` loops have three parts: init, condition, update. Use them when you know the count or need an index.
- `while` repeats while a condition is true — ideal for unknown iteration counts.
- `do...while` runs the body first, then checks — use it when at least one execution is required.
- `for...of` is the cleanest way to iterate array and string values.
- `for...in` iterates object keys — never use it on arrays.
- `break` exits a loop early; `continue` skips to the next iteration.
- Always ensure your loop's condition will eventually become `false` to prevent infinite loops.
- In modern JS, array methods (`forEach`, `map`, `filter`) often replace loops entirely.

**Next up:** Functions and Parameters — packaging reusable logic with inputs and outputs.
