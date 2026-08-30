# Control Flow: Conditionals

Programs need to make decisions. Should we show a login button or a dashboard? Is the form valid? Did the player win? Conditionals let your code take different paths based on what's true right now.

## if / else

The most fundamental branching construct. If the condition is `true`, the block runs. Otherwise, the `else` block runs (if provided).

```js
let temperature = 35;

if (temperature > 30) {
  console.log("It's hot outside");
} else {
  console.log("It's not too bad");
}
// → "It's hot outside"
```

You can omit `else` when there's no alternate path:

```js
let loggedIn = true;

if (!loggedIn) {
  console.log("Please sign in");
}
// Nothing prints — loggedIn is true
```

## else-if Chains

When there are more than two branches, use `else if`:

```js
let score = 78;

if (score >= 90) {
  console.log("A — Excellent");
} else if (score >= 80) {
  console.log("B — Good");
} else if (score >= 70) {
  console.log("C — Passing");   // → this one runs
} else {
  console.log("F — Below passing");
}
```

:::warning
Order matters. Conditions are checked top-to-bottom, and the **first** match wins. If you checked `score >= 70` before `score >= 90`, a score of 95 would hit the wrong branch.
:::

## Truthy and Falsy Values

Inside an `if`, JavaScript doesn't require an actual boolean. It coerces any value to `true` or `false`. The **falsy** values are:

| Value | Type |
|-------|------|
| `false` | boolean |
| `0` | number |
| `-0` | number |
| `""` | empty string |
| `null` | null |
| `undefined` | undefined |
| `NaN` | number |

**Everything else is truthy** — including `"0"`, `" "` (space), `[]`, and `{}`.

```js
if ("hello") console.log("truthy");   // runs
if (42) console.log("truthy");        // runs
if ([]) console.log("truthy");        // runs — empty array is truthy!

if (0) console.log("nope");           // skipped
if ("") console.log("nope");          // skipped
if (null) console.log("nope");        // skipped
```

:::key
The six falsy values: `false`, `0`, `""`, `null`, `undefined`, `NaN`. Memorize these — you'll use this knowledge in every conditional you write.
:::

This lets you write concise existence checks:

```js
let username = "";

if (!username) {
  console.log("Username is required");  // runs — empty string is falsy
}

let items = null;
if (items) {
  console.log(items.length);  // safely skipped — items is null
}
```

## Combining Conditions with && and ||

Build complex conditions by combining simple ones:

```js
let age = 25;
let hasLicense = true;
let isBanned = false;

// AND — all must be true
if (age >= 16 && hasLicense && !isBanned) {
  console.log("Can drive");  // → Can drive
}

// OR — at least one must be true
let isAdmin = false;
let isModerator = true;

if (isAdmin || isModerator) {
  console.log("Has elevated access");  // → Has elevated access
}
```

Practical pattern — range check:

```js
let value = 50;

if (value >= 1 && value <= 100) {
  console.log("In valid range");
}
```

## The Ternary Operator

A compact conditional that **produces a value**: `condition ? valueIfTrue : valueIfFalse`.

```js
let age = 20;
let status = age >= 18 ? "adult" : "minor";
console.log(status);  // → "adult"
```

Ternaries shine when you need a value inline:

```js
let score = 85;
console.log(`Result: ${score >= 60 ? "PASS" : "FAIL"}`);
// → "Result: PASS"

let items = 3;
let label = `${items} item${items === 1 ? "" : "s"}`;
console.log(label);  // → "3 items"
```

:::warning
Ternaries are for simple two-branch decisions. If you need `else if` or side effects, use a proper `if/else` block. Never nest ternaries.
:::

## switch / case

When you're comparing one value against several fixed options, `switch` reads cleaner than a long `else if` chain:

```js
let role = "editor";

switch (role) {
  case "admin":
    console.log("Full access");
    break;
  case "editor":
    console.log("Can edit content");  // → this runs
    break;
  case "viewer":
    console.log("Read only");
    break;
  default:
    console.log("Unknown role");
}
```

### Don't Forget break

Without `break`, execution **falls through** to the next case:

```js
let fruit = "apple";

// BUG — missing break
switch (fruit) {
  case "apple":
    console.log("Apple");     // runs
  case "banana":
    console.log("Banana");    // also runs! (fall-through)
    break;
  case "cherry":
    console.log("Cherry");
    break;
}
// Prints both "Apple" and "Banana"
```

Sometimes fall-through is intentional — grouping cases that share the same logic:

```js
let day = "Saturday";

switch (day) {
  case "Saturday":
  case "Sunday":
    console.log("Weekend");   // → Weekend
    break;
  default:
    console.log("Weekday");
}
```

:::tip
`switch` uses strict comparison (`===`). `case "5"` will **not** match the number `5`.
:::

## Guard Clauses (Early Returns)

Deeply nested `if/else` blocks are hard to read. Guard clauses flip the pattern: check for invalid cases first and return early.

```js
// Nested — hard to follow
function processOrder(order) {
  if (order) {
    if (order.items.length > 0) {
      if (order.isPaid) {
        // finally, the real logic
        return shipOrder(order);
      } else {
        return "Payment required";
      }
    } else {
      return "Cart is empty";
    }
  } else {
    return "No order provided";
  }
}

// Guard clauses — flat and scannable
function processOrder(order) {
  if (!order) return "No order provided";
  if (order.items.length === 0) return "Cart is empty";
  if (!order.isPaid) return "Payment required";

  return shipOrder(order);
}
```

:::key
Guard clauses reduce nesting and put the "happy path" at the bottom. Each guard handles one failure case and exits immediately.
:::

## Practical Examples

### Form Validation

```js
function validateEmail(email) {
  if (!email) return "Email is required";
  if (!email.includes("@")) return "Must contain @";
  if (email.length < 5) return "Too short";
  return null;  // null means no error
}

let error = validateEmail("test");
if (error) {
  console.log(`Error: ${error}`);  // → "Error: Must contain @"
}
```

### Feature Flags

```js
const features = { darkMode: true, betaSearch: false };

if (features.darkMode) {
  document.body.classList.add("dark");
}

let searchVersion = features.betaSearch ? "v2" : "v1";
```

### Pricing Logic

```js
function getPrice(plan, isAnnual) {
  switch (plan) {
    case "free":
      return 0;
    case "pro":
      return isAnnual ? 8 : 12;
    case "team":
      return isAnnual ? 20 : 30;
    default:
      return null;
  }
}

console.log(getPrice("pro", true));  // → 8
```

:::quiz
Q: Which of these values is truthy?
- `0`
- `""`
- `[]` *
- `null`
E: An empty array `[]` is truthy in JavaScript. Only `false`, `0`, `""`, `null`, `undefined`, and `NaN` are falsy.
:::

:::quiz
Q: What happens in a `switch` statement when you forget `break`?
- The program throws an error
- Execution falls through to the next case *
- The `default` block runs instead
E: Without `break`, JavaScript continues executing the next case's code regardless of whether it matches — this is called fall-through.
:::

:::quiz
Q: What is a guard clause?
- A special syntax for protecting variables
- An early return that handles invalid cases before the main logic *
- A type of loop that checks conditions
E: Guard clauses check for failure cases at the top of a function and return early, keeping the happy path flat and readable.
:::

## Recap

- `if/else` is the foundation of branching. `else if` chains handle multiple conditions — order matters.
- Know the six falsy values: `false`, `0`, `""`, `null`, `undefined`, `NaN`. Everything else is truthy.
- Combine conditions with `&&` (all true) and `||` (at least one true).
- The ternary `condition ? a : b` produces a value — great for assignments and inline expressions, but don't nest it.
- `switch` compares one value against fixed cases using `===`. Always include `break` unless fall-through is intentional.
- Guard clauses (early returns) flatten deeply nested `if/else` and make functions scannable.

**Next up:** Control Flow: Loops — repeating work efficiently with `for`, `while`, and `for...of`.
