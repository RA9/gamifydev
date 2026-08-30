# Error Handling with Try Catch

Errors happen. APIs return unexpected data, users enter invalid input, properties don't exist. Professional code doesn't ignore these problems — it anticipates them, catches them, and recovers gracefully.

## try / catch / finally

Wrap risky code in a `try` block. If an error occurs, execution jumps to `catch` instead of crashing:

```js
try {
  let data = JSON.parse('{ invalid json }');
} catch (error) {
  console.log("Parsing failed:", error.message);
}
// → "Parsing failed: Expected property name or '}' ..."
// Program keeps running — no crash
```

Without `try/catch`, that same error would halt your entire script.

### The finally Block

`finally` runs **no matter what** — whether the `try` succeeded or the `catch` fired:

```js
function loadData() {
  let connection = openConnection();

  try {
    let data = connection.fetch("/api/users");
    return data;
  } catch (error) {
    console.error("Fetch failed:", error.message);
    return [];
  } finally {
    connection.close();  // always runs — cleanup is guaranteed
  }
}
```

:::key
Use `finally` for cleanup that must happen regardless of success or failure — closing connections, hiding loading spinners, resetting state.
:::

## The Error Object

When JavaScript throws an error, it creates an `Error` object with three useful properties:

```js
try {
  undeclaredVariable;
} catch (error) {
  console.log(error.name);     // → "ReferenceError"
  console.log(error.message);  // → "undeclaredVariable is not defined"
  console.log(error.stack);    // → full stack trace (file, line numbers)
}
```

| Property | What it tells you |
|----------|------------------|
| `name` | The error type (e.g., `TypeError`, `ReferenceError`) |
| `message` | A human-readable description of what went wrong |
| `stack` | The call stack at the point of the error — invaluable for debugging |

## Common Error Types

JavaScript has several built-in error types. Recognizing them helps you debug faster:

### TypeError

Happens when a value isn't the type you expected:

```js
let user = null;
user.name;          // TypeError: Cannot read properties of null
"hello".toFixed(2); // TypeError: "hello".toFixed is not a function
```

### ReferenceError

Accessing a variable that doesn't exist:

```js
console.log(myVar);  // ReferenceError: myVar is not defined
```

### SyntaxError

Code that can't be parsed (often caught at load time, not at runtime):

```js
JSON.parse("not json");  // SyntaxError: Unexpected token 'o'
```

### RangeError

A value is outside an allowed range:

```js
let arr = new Array(-1);      // RangeError: Invalid array length
function forever() { forever(); }
forever();                    // RangeError: Maximum call stack size exceeded
```

:::tip
When debugging, read the error name and message first. A `TypeError: Cannot read properties of undefined` means something is `undefined` when you expected an object. Trace backwards to find where the value went missing.
:::

## The throw Keyword

You can throw your own errors to signal problems explicitly:

```js
function divide(a, b) {
  if (b === 0) {
    throw new Error("Cannot divide by zero");
  }
  return a / b;
}

try {
  let result = divide(10, 0);
} catch (error) {
  console.log(error.message);  // → "Cannot divide by zero"
}
```

Always throw `Error` objects (not strings) — they include the stack trace:

```js
// ❌ Avoid — no stack trace, hard to debug
throw "something went wrong";

// ✅ Correct — full Error object
throw new Error("something went wrong");
```

## Custom Errors

Create specific error types for different failure categories:

```js
class ValidationError extends Error {
  constructor(field, message) {
    super(message);
    this.name = "ValidationError";
    this.field = field;
  }
}

class NotFoundError extends Error {
  constructor(resource) {
    super(`${resource} not found`);
    this.name = "NotFoundError";
    this.resource = resource;
  }
}

function getUser(id) {
  if (typeof id !== "number") {
    throw new ValidationError("id", "ID must be a number");
  }
  // ... lookup logic
  throw new NotFoundError("User");
}
```

Catch specific types with `instanceof`:

```js
try {
  getUser("abc");
} catch (error) {
  if (error instanceof ValidationError) {
    console.log(`Invalid ${error.field}: ${error.message}`);
  } else if (error instanceof NotFoundError) {
    console.log(error.message);
  } else {
    throw error;  // re-throw unexpected errors
  }
}
```

:::key
Custom errors make your code self-documenting. `throw new ValidationError("email", "Invalid format")` is far more informative than `throw new Error("bad input")`.
:::

## When to Catch vs When to Let It Throw

Not every error should be caught. The question is: **can this code meaningfully recover?**

### Catch When You Can Recover

```js
// Parse user input — might fail, but we have a fallback
function safeParseJSON(text) {
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

// API call — might fail, but we can show an error state
try {
  const data = await fetchUserProfile(userId);
  renderProfile(data);
} catch (error) {
  renderError("Could not load profile. Please try again.");
}
```

### Let It Throw When You Can't

```js
// If the config file is missing, the app simply can't start
// Don't catch this — let it crash with a clear error
const config = JSON.parse(readFile("config.json"));

// If a required dependency is broken, catching hides the real problem
import { critical } from "./core.js";
```

:::warning
Catching errors you can't handle is worse than not catching them. Silent failures hide bugs. If you catch an error, either recover from it or log it and re-throw.
:::

```js
// ❌ Silent swallow — the bug disappears
try {
  riskyOperation();
} catch (error) {
  // nothing — bug silently vanishes
}

// ✅ Log and re-throw if you can't recover
try {
  riskyOperation();
} catch (error) {
  console.error("Unexpected error in riskyOperation:", error);
  throw error;
}
```

## Defensive Coding Patterns

### Guard with Early Validation

Validate inputs before doing work — fail fast with clear messages:

```js
function processPayment(amount, currency) {
  if (typeof amount !== "number" || amount <= 0) {
    throw new Error(`Invalid amount: ${amount}. Must be a positive number.`);
  }
  if (!["USD", "EUR", "GBP"].includes(currency)) {
    throw new Error(`Unsupported currency: ${currency}`);
  }

  // Safe to proceed
  return chargeCard(amount, currency);
}
```

### Safe Property Access

```js
// Optional chaining + nullish coalescing
const city = user?.address?.city ?? "Unknown";
const firstItem = data?.items?.[0] ?? null;
```

### Default Values for Missing Data

```js
function renderCard(product) {
  const name = product.name ?? "Untitled";
  const price = product.price ?? 0;
  const image = product.image ?? "/placeholder.png";

  return `<div>${name} — $${price}</div>`;
}
```

### Type Checking Before Operations

```js
function formatPrice(value) {
  if (typeof value !== "number") {
    throw new TypeError(`Expected number, got ${typeof value}`);
  }
  return `$${value.toFixed(2)}`;
}

formatPrice(9.99);    // → "$9.99"
formatPrice("ten");   // → TypeError: Expected number, got string
```

## Error Handling in Practice

A realistic pattern combining multiple techniques:

```js
class APIError extends Error {
  constructor(status, message) {
    super(message);
    this.name = "APIError";
    this.status = status;
  }
}

async function fetchUser(id) {
  if (!id) throw new ValidationError("id", "User ID is required");

  try {
    const response = await fetch(`/api/users/${id}`);

    if (!response.ok) {
      throw new APIError(response.status, `HTTP ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    if (error instanceof APIError && error.status === 404) {
      return null;  // user not found — recoverable
    }
    throw error;  // network failures, 500s — re-throw
  }
}

// Usage
try {
  const user = await fetchUser(42);
  if (user) {
    renderProfile(user);
  } else {
    renderNotFound();
  }
} catch (error) {
  renderError("Something went wrong. Please try again.");
  console.error(error);
}
```

:::quiz
Q: When should you use `finally`?
- Only when an error occurs
- For cleanup code that must run regardless of success or failure *
- To retry the failed operation
E: `finally` always runs — whether `try` succeeds or `catch` fires. It's the right place for cleanup like closing connections, hiding loaders, or resetting state.
:::

:::quiz
Q: Why should you throw `new Error("msg")` instead of `throw "msg"`?
- Strings are invalid throw values
- Error objects include a stack trace, making debugging much easier *
- There is no practical difference
E: Error objects contain `name`, `message`, and `stack` properties. A thrown string has none of these, making it much harder to trace the source of the problem.
:::

:::quiz
Q: What is the danger of an empty `catch` block?
- It causes a syntax error
- It silently swallows errors, hiding bugs *
- It makes the code run slower
E: An empty `catch` makes errors disappear with no trace. The bug still exists — you just can't see it anymore. Always log, handle, or re-throw caught errors.
:::

## Recap

- `try/catch/finally` prevents errors from crashing your program. `finally` always runs.
- Error objects have `name`, `message`, and `stack` — read them when debugging.
- Common types: `TypeError` (wrong type), `ReferenceError` (undefined variable), `SyntaxError` (bad parsing), `RangeError` (out of bounds).
- Use `throw new Error("message")` to signal problems explicitly. Never throw raw strings.
- Custom errors (`class MyError extends Error`) make error handling precise and self-documenting.
- Catch errors only when you can recover. Silent swallowing hides bugs — log and re-throw if you can't fix it.
- Defend proactively: validate inputs early, use optional chaining and defaults, and check types before operations.

**Next up:** Asynchronous JavaScript — handling promises, async/await, and working with APIs.
