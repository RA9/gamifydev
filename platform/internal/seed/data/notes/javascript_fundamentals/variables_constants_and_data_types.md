# JavaScript Basics

JavaScript is the programming language that brings web pages to life — clicks, animations, form checks, live updates. This lesson gives you the rock-solid foundation: variables, types, operators, decisions, and loops.

## What JavaScript Is and Where It Runs

HTML is the structure of a page and CSS is the styling, but JavaScript (JS) is the *behavior*. It is a full programming language that runs **inside the browser**. Every browser — Chrome, Firefox, Safari, Edge — ships with a JavaScript engine that reads your code and executes it.

That means when someone visits your page, *their* browser runs your JS. You don't need to install anything special to start; you already have a browser.

:::analogy
Think of a web page as a house. HTML is the walls and rooms, CSS is the paint and furniture, and JavaScript is the electricity and plumbing — the stuff that actually *does* things when you flip a switch.
:::

JavaScript can also run outside the browser (on servers, via Node.js), but in this Frontend course we live entirely in the browser.

## How to Add JavaScript to a Page

There are two common ways. The recommended pattern is a `<script>` tag placed at the **end of the body**, just before `</body>`:

```html
<!DOCTYPE html>
<html>
  <body>
    <h1>Hello</h1>

    <!-- Script goes last so the HTML above is ready first -->
    <script>
      console.log("The page loaded!");
    </script>
  </body>
</html>
```

Why at the end? The browser reads top to bottom. If your script runs before the HTML exists, it can't find the elements it wants to work with. Putting it last means everything above it is ready.

The cleaner approach for real projects is an **external file**:

```html
<body>
  <h1>Hello</h1>
  <script src="app.js"></script>
</body>
```

```js
// app.js
console.log("Code lives in its own file now.");
```

:::tip
Keep your JS in a separate `.js` file once a project grows. It keeps your HTML readable and lets the browser cache the script.
:::

## console.log and the Dev Console

`console.log()` prints a message to the browser's **developer console**. It is your single most-used tool for seeing what your code is doing.

Open the console: right-click the page → **Inspect** → **Console** tab (or press F12).

```js
console.log("Hello, world!");   // → Hello, world!
console.log(42);                // → 42
console.log(3 + 4);             // → 7
```

You can log several things at once:

```js
console.log("Score:", 10, "Lives:", 3);  // → Score: 10 Lives: 3
```

## Variables: let and const

A variable is a named box that holds a value. You create one with `let` or `const`.

```js
let score = 0;        // can change later
score = 10;           // ✅ allowed
console.log(score);   // → 10

const name = "Ada";   // cannot be reassigned
console.log(name);    // → Ada
```

Try reassigning a `const` and you get an error:

```js
const pi = 3.14;
pi = 3.15;   // ❌ TypeError: Assignment to constant variable.
```

**Rule of thumb:** use `const` by default. Only reach for `let` when you *know* the value will change (like a counter or a running total). This makes your intent obvious and prevents accidental changes.

:::warning
You'll see older code use `var`. It's legacy — it behaves in surprising ways with scope and is best avoided in new code. Stick to `let` and `const`.
:::

### Naming Variables

Names should describe what they hold. JavaScript uses **camelCase** by convention: first word lowercase, later words capitalized.

```js
let userAge = 25;
let isLoggedIn = true;
const maxAttempts = 3;
```

Rules: start with a letter, `$`, or `_` (not a number), no spaces, and don't use reserved words like `let` or `for`. Good names are worth the typing — `totalPrice` beats `tp`.

## Data Types

A value's *type* describes what kind of thing it is. The basic (primitive) types you'll use constantly:

```js
let title = "Level One";   // string  → text, in quotes
let lives = 3;             // number  → integers and decimals alike
let isPaused = false;      // boolean → true or false
let winner = null;         // null    → intentional "no value"
let prize;                 // undefined → declared but no value yet

console.log(typeof title);    // → string
console.log(typeof lives);    // → number
console.log(typeof isPaused); // → boolean
console.log(typeof prize);    // → undefined
```

Strings can use single or double quotes — just be consistent:

```js
let a = "double";
let b = 'single';
```

Numbers don't distinguish integers from decimals; it's all `number`:

```js
let count = 7;
let price = 9.99;
```

You'll also meet **objects** and **arrays** soon — they group many values together. A quick peek:

```js
let player = { name: "Ada", level: 5 };  // object: labeled values
let scores = [10, 20, 30];               // array: an ordered list
```

We cover those in depth in the next lesson.

## Template Literals

To build strings with values inside them, use **backticks** (`` ` ``) and `${...}`. This is far cleaner than gluing pieces with `+`.

```js
let name = "Ada";
let level = 5;

// The old way:
console.log("Player " + name + " is level " + level);
// → Player Ada is level 5

// Template literal — much nicer:
console.log(`Player ${name} is level ${level}`);
// → Player Ada is level 5
```

You can put any expression inside `${}`:

```js
let a = 6, b = 4;
console.log(`Total: ${a + b}`);  // → Total: 10
```

Backtick strings can also span multiple lines:

```js
let msg = `Line one
Line two`;
```

## Operators

### Arithmetic

```js
console.log(5 + 2);   // → 7   addition
console.log(5 - 2);   // → 3   subtraction
console.log(5 * 2);   // → 10  multiplication
console.log(5 / 2);   // → 2.5 division
console.log(5 % 2);   // → 1   remainder (modulo)
```

The `%` (modulo) operator gives the *remainder* of a division. It's incredibly handy — for example, a number is even if `n % 2 === 0`.

```js
console.log(10 % 2);  // → 0  (even)
console.log(7 % 2);   // → 1  (odd)
```

### Comparison: === vs ==

Comparison operators return a boolean.

```js
console.log(5 > 3);    // → true
console.log(5 < 3);    // → false
console.log(5 >= 5);   // → true
console.log(5 === 5);  // → true   (strictly equal)
console.log(5 !== 3);  // → true   (not equal)
```

Here's a crucial one. `===` checks value **and** type (strict). `==` checks value but will *convert types first* (loose), which leads to confusing results:

```js
console.log(5 === "5");  // → false  (number vs string)
console.log(5 == "5");   // → true   (== converts "5" to 5 first!)
console.log(0 == "");    // → true   😱
console.log(0 === "");   // → false
```

:::key
Always use `===` and `!==`. They are predictable. The loose `==` does silent type conversion that surprises even experienced developers.
:::

### Logical Operators

```js
console.log(true && false);  // → false  (AND: both must be true)
console.log(true || false);  // → true   (OR: at least one true)
console.log(!true);          // → false  (NOT: flips it)
```

Combine them with comparisons:

```js
let age = 20;
console.log(age >= 18 && age < 65);  // → true
```

## Type Coercion Gotchas

JavaScript sometimes auto-converts types, especially with `+`. Because `+` also means "join strings," a number next to a string becomes text:

```js
console.log(2 + 2);       // → 4    (both numbers)
console.log(2 + "2");     // → "22" (one is a string → concatenation)
console.log("3" * 2);     // → 6    (* has no string meaning → both become numbers)
console.log("5" - 1);     // → 4    (- forces numbers too)
```

:::warning
The classic trap: `+` with any string concatenates. If a value came from a form input or `prompt()`, it's a string — convert it with `Number(x)` before doing math.
:::

```js
let input = "10";
console.log(Number(input) + 5);  // → 15
console.log(String(42));         // → "42"
```

## Conditionals

### if / else if / else

Run code only when a condition is true.

```js
let score = 75;

if (score >= 90) {
  console.log("A");
} else if (score >= 70) {
  console.log("B");   // → B
} else {
  console.log("Keep going");
}
```

### Ternary Operator

A compact `if/else` that produces a value: `condition ? ifTrue : ifFalse`.

```js
let age = 20;
let status = age >= 18 ? "adult" : "minor";
console.log(status);  // → adult
```

### switch

When you compare one value against many fixed options, `switch` reads cleanly. Don't forget `break`.

```js
let day = "Mon";

switch (day) {
  case "Sat":
  case "Sun":
    console.log("Weekend");
    break;
  case "Mon":
    console.log("Monday");  // → Monday
    break;
  default:
    console.log("A weekday");
}
```

## Truthy and Falsy

Inside an `if`, any value is treated as either "truthy" or "falsy." The **falsy** values are: `false`, `0`, `""` (empty string), `null`, `undefined`, and `NaN`. Everything else is truthy.

```js
if ("hello") console.log("runs");  // → runs (non-empty string is truthy)
if (0) console.log("nope");        // (never runs — 0 is falsy)

let name = "";
if (!name) console.log("Name is empty");  // → Name is empty
```

:::tip
This lets you check "is there a value?" simply: `if (username) { ... }` runs only when `username` isn't empty/undefined.
:::

## Loops

Loops repeat code so you don't write it by hand.

### for

The workhorse. Three parts: start, condition, step.

```js
for (let i = 0; i < 3; i++) {
  console.log(i);  // → 0, then 1, then 2
}
```

`i++` adds 1 each time. The loop stops when `i < 3` becomes false.

### while

Repeats *while* a condition holds. Use it when you don't know the count ahead of time.

```js
let n = 3;
while (n > 0) {
  console.log(n);  // → 3, 2, 1
  n--;
}
```

:::warning
Always make sure the condition eventually becomes false. If `n--` were missing above, the loop would run forever and freeze the page — an infinite loop.
:::

### for...of

The cleanest way to walk through the items of an array (or any list).

```js
let fruits = ["apple", "pear", "kiwi"];
for (let fruit of fruits) {
  console.log(fruit);  // → apple, pear, kiwi
}
```

## Comments

Comments are notes the engine ignores. Use them to explain *why*, not just *what*.

```js
// This is a single-line comment

/* This is a
   multi-line comment */

let tax = 0.08;  // 8% sales tax
```

## Practice

:::predict
Q: What does this print? `console.log(2 + "2")`
- 4
- "22" *
- error
E: `+` with a string concatenates, so the number `2` becomes "2" and they join into "22".
:::

:::predict
Q: What does this print? `console.log(5 === "5")`
- true
- false *
- "5"
E: `===` checks value *and* type. A number and a string are different types, so it's `false`.
:::

:::predict
Q: After this loop, what is logged? `for (let i = 0; i < 2; i++) { console.log(i); }`
- 1, 2
- 0, 1 *
- 0, 1, 2
E: `i` starts at 0 and the loop runs while `i < 2`, so it logs 0 and 1, then stops.
:::

:::quiz
Q: Which keyword should you reach for by default when declaring a variable?
- var
- let
- const *
E: Prefer `const`; it signals the value won't be reassigned. Use `let` only when it must change.
:::

:::fill
Q: The operator that returns the remainder of a division is `___`.
A: %
E: `%` is the modulo operator. `7 % 2` is `1`.
:::

## Recap

- JavaScript runs in the browser and controls a page's behavior; add it with a `<script>` tag at the end of the body or via an external `.js` file.
- `console.log()` prints to the dev console — your main debugging tool.
- Use `const` by default, `let` when a value changes; avoid legacy `var`. Name variables clearly in camelCase.
- Core types: string, number, boolean, null, undefined (plus objects and arrays).
- Build strings with backtick template literals and `${}`.
- Prefer `===` over `==` to avoid surprising type coercion; `+` with a string concatenates.
- Make decisions with `if/else if/else`, the ternary, and `switch`; remember the falsy values.
- Repeat work with `for`, `while`, and `for...of` — and never write an infinite loop.

**Next up:** JavaScript Functions, Arrays and Objects — packaging code into reusable functions and working with the data shapes you'll use every day.
