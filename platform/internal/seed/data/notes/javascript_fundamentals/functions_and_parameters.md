# Functions and Parameters

Functions are the primary unit of reusable code in JavaScript. Instead of copying logic every time you need it, you write it once, give it a name, and call it wherever you want. Understanding how they work — declarations, parameters, return values, and the call stack — is essential.

## Function Declarations

A function declaration uses the `function` keyword, followed by a name, parentheses for parameters, and a body in braces:

```js
function greet() {
  console.log("Hello!");
}

greet();  // → Hello!
greet();  // → Hello!  (call it as many times as you need)
```

The name should describe what the function *does*. JavaScript convention is camelCase, and verbs make the best function names:

```js
function calculateTotal() { /* ... */ }
function validateEmail() { /* ... */ }
function formatDate() { /* ... */ }
function isEligible() { /* ... */ }    // "is" prefix for boolean returns
function getUser() { /* ... */ }       // "get" prefix for retrievals
```

## Function Expressions

You can also assign a function to a variable:

```js
const greet = function () {
  console.log("Hello!");
};

greet();  // → Hello!
```

The key difference: **declarations are hoisted** (available before they appear in the code), expressions are not.

```js
// Works — declarations are hoisted
sayHi();
function sayHi() {
  console.log("Hi!");
}

// ❌ Error — expressions are NOT hoisted
sayBye();  // TypeError: sayBye is not a function
const sayBye = function () {
  console.log("Bye!");
};
```

:::tip
For top-level utility functions, declarations are fine — hoisting makes them available everywhere in the file. For inline callbacks and local helpers, expressions (usually arrow functions) are preferred.
:::

## Parameters vs Arguments

A **parameter** is the placeholder in the function definition. An **argument** is the actual value you pass when calling.

```js
function add(a, b) {    // a and b are parameters
  return a + b;
}

add(3, 4);   // 3 and 4 are arguments
```

Multiple parameters are separated by commas:

```js
function introduce(name, role, company) {
  console.log(`${name} is a ${role} at ${company}`);
}

introduce("Ada", "Engineer", "Acme");
// → "Ada is a Engineer at Acme"
```

If you pass fewer arguments than parameters, the missing ones are `undefined`:

```js
function greet(name) {
  console.log(`Hello, ${name}!`);
}

greet();  // → "Hello, undefined!"
```

## Return Values

`console.log` only **prints**. To send a value *back* to the caller so it can be stored or used, use `return`:

```js
function multiply(a, b) {
  return a * b;
}

let result = multiply(6, 7);
console.log(result);            // → 42
console.log(multiply(3, 10));   // → 30
```

:::key
`return` does two things: it sends a value back to the caller, and it immediately exits the function. No code after `return` will run.
:::

```js
function checkAge(age) {
  if (age < 0) return "Invalid age";
  if (age < 18) return "Minor";
  return "Adult";
  console.log("This never runs");  // dead code
}

console.log(checkAge(25));  // → "Adult"
console.log(checkAge(-1));  // → "Invalid age"
```

A function with no `return` (or a bare `return;`) gives back `undefined`:

```js
function doSomething() {
  let x = 10;
  // no return statement
}

let result = doSomething();
console.log(result);  // → undefined
```

:::warning
If `console.log(myFunction())` shows `undefined`, you probably forgot to `return` a value. `console.log` inside a function is for debugging — `return` is for producing a result.
:::

## Default Parameters

Give parameters a fallback value for when no argument is passed:

```js
function greet(name = "friend") {
  return `Hello, ${name}!`;
}

console.log(greet("Ada"));  // → "Hello, Ada!"
console.log(greet());       // → "Hello, friend!"
```

Defaults can be expressions and can reference earlier parameters:

```js
function createUser(name, role = "viewer", createdAt = Date.now()) {
  return { name, role, createdAt };
}

console.log(createUser("Ada"));
// → { name: "Ada", role: "viewer", createdAt: 1698234567890 }

console.log(createUser("Bob", "admin"));
// → { name: "Bob", role: "admin", createdAt: 1698234567891 }
```

## The Call Stack

When a function calls another function, JavaScript keeps track using a **call stack** — a list of "where we are" that grows and shrinks as functions are entered and exited.

```js
function first() {
  console.log("first start");
  second();
  console.log("first end");
}

function second() {
  console.log("second start");
  third();
  console.log("second end");
}

function third() {
  console.log("third");
}

first();
// → first start → second start → third → second end → first end
```

The stack grows: `first` → `second` → `third`. Then unwinds: `third` finishes → `second` finishes → `first` finishes.

:::warning
If functions call each other without stopping (infinite recursion), the stack overflows: `RangeError: Maximum call stack size exceeded`. Always have a base case that stops recursion.
:::

## Pure Functions

A **pure function** has two properties:

1. Same inputs always produce the same output
2. No side effects (doesn't modify anything outside itself)

```js
// ✅ Pure — depends only on inputs, changes nothing external
function add(a, b) {
  return a + b;
}

// ❌ Impure — depends on external state
let multiplier = 3;
function multiply(n) {
  return n * multiplier;  // result changes if multiplier changes
}

// ❌ Impure — has a side effect (modifies external array)
let log = [];
function addToLog(message) {
  log.push(message);  // mutates external state
}
```

:::key
Pure functions are easier to test, debug, and reason about. Prefer them when possible. Not everything can be pure (UI updates, API calls), but your data-processing logic should be.
:::

## Functions as First-Class Values

In JavaScript, functions are values — just like numbers or strings. You can store them in variables, put them in arrays, pass them to other functions, and return them from functions.

```js
// Store in a variable
const double = function (n) { return n * 2; };

// Put in an array
const operations = [
  function (n) { return n + 1; },
  function (n) { return n * 2; },
  function (n) { return n ** 2; },
];

// Pass as an argument
function applyOperation(value, operation) {
  return operation(value);
}

console.log(applyOperation(5, double));  // → 10

// Return from a function
function createMultiplier(factor) {
  return function (n) {
    return n * factor;
  };
}

const triple = createMultiplier(3);
console.log(triple(7));  // → 21
```

This "first-class" nature is what makes callbacks, array methods, and event handlers possible — all topics you'll explore next.

:::quiz
Q: What is the difference between a parameter and an argument?
- There is no difference
- A parameter is the placeholder in the definition; an argument is the actual value passed *
- An argument is the placeholder; a parameter is the value
E: Parameters are the names listed in the function definition. Arguments are the real values provided when calling the function.
:::

:::quiz
Q: What does a function return if it has no `return` statement?
- `null`
- `0`
- `undefined` *
E: A function without an explicit `return` (or with a bare `return;`) automatically returns `undefined`.
:::

:::quiz
Q: What makes a function "pure"?
- It uses arrow syntax
- Same inputs always produce the same output, with no side effects *
- It takes no parameters
E: Pure functions are deterministic (same input → same output) and have no side effects (they don't modify external state). This makes them easy to test and reason about.
:::

## Recap

- Function declarations are hoisted; function expressions are not.
- Parameters are placeholders; arguments are the values you pass.
- `return` sends a value back and exits the function. No `return` means `undefined`.
- Default parameters provide fallback values: `function greet(name = "friend")`.
- The call stack tracks nested function calls; infinite recursion overflows it.
- Name functions with verbs in camelCase: `calculateTotal`, `isValid`, `getUser`.
- Pure functions (same input → same output, no side effects) are easier to test and debug.
- Functions are first-class values — you can store, pass, and return them like any other value.

**Next up:** Arrow Functions and Callbacks — the modern syntax and the patterns that power array methods and event handling.
