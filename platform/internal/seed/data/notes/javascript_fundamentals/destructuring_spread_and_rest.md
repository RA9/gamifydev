# Destructuring Spread and Rest

Destructuring, spread, and rest are three ES6+ features that make working with objects and arrays dramatically cleaner. They eliminate boilerplate, reduce intermediate variables, and show up in virtually every modern JavaScript codebase.

## Object Destructuring

Instead of accessing properties one by one, pull out the values you need in a single statement:

```js
const user = { name: "Ada", role: "Engineer", level: 5 };

// Without destructuring
const name = user.name;
const role = user.role;

// With destructuring
const { name, role } = user;

console.log(name);  // → "Ada"
console.log(role);  // → "Engineer"
```

The variable names must match the object's keys. You only need to extract the properties you care about — the rest are ignored.

### Renaming Variables

When the key name isn't what you want (or conflicts with an existing variable), rename it with a colon:

```js
const response = { data: [1, 2, 3], status: 200 };

const { data: items, status: httpStatus } = response;

console.log(items);       // → [1, 2, 3]
console.log(httpStatus);  // → 200
// console.log(data);     // ❌ ReferenceError — 'data' was renamed to 'items'
```

### Default Values

Provide fallbacks for properties that might be missing:

```js
const config = { theme: "dark" };

const { theme, language = "en", fontSize = 14 } = config;

console.log(theme);     // → "dark"     (from the object)
console.log(language);  // → "en"       (default — not in object)
console.log(fontSize);  // → 14         (default — not in object)
```

Defaults apply only when the value is `undefined`, not when it's `null` or `0`:

```js
const { count = 10 } = { count: 0 };
console.log(count);  // → 0  (0 is not undefined, so default doesn't apply)
```

## Array Destructuring

Works by position instead of by name:

```js
const rgb = [255, 128, 0];

const [red, green, blue] = rgb;

console.log(red);    // → 255
console.log(green);  // → 128
console.log(blue);   // → 0
```

### Skipping Elements

Use empty slots to skip values you don't need:

```js
const scores = [90, 85, 78, 92];

const [first, , third] = scores;

console.log(first);  // → 90
console.log(third);  // → 78
```

### Swapping Variables

A classic trick — no temporary variable needed:

```js
let a = 1;
let b = 2;

[a, b] = [b, a];

console.log(a);  // → 2
console.log(b);  // → 1
```

:::tip
Array destructuring with functions that return arrays is common: `const [value, setValue] = useState(0)` in React, or `const [key, value] = entry` when iterating `Object.entries`.
:::

## Nested Destructuring

You can destructure deeply nested structures in one statement:

```js
const user = {
  name: "Ada",
  address: {
    city: "Lagos",
    country: "Nigeria",
  },
};

const { name, address: { city, country } } = user;

console.log(name);     // → "Ada"
console.log(city);     // → "Lagos"
console.log(country);  // → "Nigeria"
// console.log(address); // ❌ address itself is not assigned
```

Nested array destructuring:

```js
const matrix = [[1, 2], [3, 4]];

const [[a, b], [c, d]] = matrix;

console.log(a, b, c, d);  // → 1 2 3 4
```

:::warning
Nested destructuring fails if the intermediate value is `undefined` or `null`. Use defaults or optional chaining to guard against missing data: `const { address: { city } = {} } = user;`
:::

## Destructuring in Function Parameters

Instead of accepting a whole object and accessing properties inside the function, destructure right in the parameter list:

```js
// Without destructuring
function greet(user) {
  return `Hello, ${user.name}! You're level ${user.level}.`;
}

// With destructuring
function greet({ name, level }) {
  return `Hello, ${name}! You're level ${level}.`;
}

greet({ name: "Ada", level: 5 });
// → "Hello, Ada! You're level 5."
```

Add defaults for optional properties:

```js
function createButton({ label = "Click", color = "blue", size = "md" } = {}) {
  return `<button class="${color} ${size}">${label}</button>`;
}

createButton({ label: "Submit", color: "green" });
// → '<button class="green md">Submit</button>'

createButton();  // works because of the = {} default
// → '<button class="blue md">Click</button>'
```

:::key
Destructuring in parameters makes function signatures self-documenting. You can see exactly which properties the function uses without reading the body.
:::

## Spread Syntax (...)

The spread operator expands an iterable (array, object) into individual elements.

### Spreading Arrays

```js
const front = [1, 2, 3];
const back = [4, 5, 6];

const combined = [...front, ...back];
console.log(combined);  // → [1, 2, 3, 4, 5, 6]

// Insert in the middle
const withMiddle = [0, ...front, 99, ...back];
console.log(withMiddle);  // → [0, 1, 2, 3, 99, 4, 5, 6]
```

Copy an array (shallow):

```js
const original = [1, 2, 3];
const copy = [...original];

copy.push(4);
console.log(original);  // → [1, 2, 3]  (unaffected)
```

### Spreading Objects

```js
const defaults = { theme: "light", fontSize: 14, lang: "en" };
const userPrefs = { theme: "dark", fontSize: 18 };

const config = { ...defaults, ...userPrefs };
console.log(config);
// → { theme: "dark", fontSize: 18, lang: "en" }
```

Later spreads override earlier ones for the same key. This is the foundation of immutable update patterns.

## Rest Parameters (...args)

Rest syntax collects remaining elements into an array. It looks like spread but works in reverse — it gathers instead of expanding.

### In Function Parameters

```js
function sum(...numbers) {
  let total = 0;
  for (let n of numbers) {
    total += n;
  }
  return total;
}

console.log(sum(1, 2, 3));       // → 6
console.log(sum(10, 20, 30, 40)); // → 100
```

Combine named parameters with rest:

```js
function log(level, ...messages) {
  console.log(`[${level}]`, ...messages);
}

log("INFO", "Server started", "on port 3000");
// → [INFO] Server started on port 3000
```

:::warning
The rest parameter must be the **last** parameter. `function bad(...rest, last)` is a syntax error.
:::

### In Destructuring

Collect the "remaining" properties or elements:

```js
// Object rest
const user = { name: "Ada", role: "admin", level: 5, score: 100 };
const { name, ...rest } = user;

console.log(name);  // → "Ada"
console.log(rest);  // → { role: "admin", level: 5, score: 100 }

// Array rest
const [first, second, ...remaining] = [1, 2, 3, 4, 5];

console.log(first);      // → 1
console.log(remaining);  // → [3, 4, 5]
```

:::tip
Object rest is great for separating known properties from "everything else" — a common pattern when passing props in component-based UI frameworks.
:::

## Immutable Update Patterns

Spread enables updating data without mutating the original — critical for predictable state management:

### Update a Property

```js
const state = { user: "Ada", theme: "light", count: 0 };

const nextState = { ...state, count: state.count + 1 };

console.log(state.count);      // → 0  (unchanged)
console.log(nextState.count);  // → 1
```

### Add to an Array

```js
const todos = ["Plan", "Build"];
const nextTodos = [...todos, "Test"];

console.log(todos);      // → ["Plan", "Build"]
console.log(nextTodos);  // → ["Plan", "Build", "Test"]
```

### Remove from an Array

```js
const items = ["a", "b", "c", "d"];
const without = items.filter(item => item !== "c");

console.log(without);  // → ["a", "b", "d"]
```

### Update a Nested Object

```js
const user = {
  name: "Ada",
  settings: { theme: "light", notifications: true },
};

const updated = {
  ...user,
  settings: {
    ...user.settings,
    theme: "dark",
  },
};

console.log(user.settings.theme);     // → "light"  (unchanged)
console.log(updated.settings.theme);  // → "dark"
```

:::key
Immutable updates create new objects instead of modifying existing ones. This makes state changes traceable and prevents bugs from shared references.
:::

:::quiz
Q: What does `const { data: items } = response` do?
- Creates a variable called `data` with value `response.items`
- Creates a variable called `items` with the value of `response.data` *
- Creates both `data` and `items` variables
E: The colon in destructuring renames: it reads the `data` property from `response` and assigns it to a new variable named `items`.
:::

:::quiz
Q: What is the difference between spread and rest?
- They are the same thing
- Spread expands elements out; rest gathers elements together *
- Spread is for arrays; rest is for objects
E: Spread (`...arr` in an expression) expands into individual elements. Rest (`...args` in a parameter or destructuring target) collects remaining elements into an array or object.
:::

:::quiz
Q: Why are immutable update patterns important?
- They make code run faster
- They prevent bugs from shared references and make state changes traceable *
- They are required by JavaScript
E: Mutating shared objects causes hard-to-track bugs. Immutable updates create new objects, so the original stays intact and changes are explicit.
:::

## Recap

- Object destructuring extracts properties by name: `const { a, b } = obj`. Rename with `:`, default with `=`.
- Array destructuring extracts by position: `const [x, y] = arr`. Skip with `,`, swap with `[a, b] = [b, a]`.
- Nested destructuring unpacks deep structures in one statement.
- Destructuring in function parameters makes signatures self-documenting.
- Spread (`...`) expands arrays and objects — use it for copying, combining, and immutable updates.
- Rest (`...`) gathers remaining items — use it for variadic functions and "collect the rest" destructuring.
- Immutable update patterns (spread + overwrite) create new objects without mutating the original.

**Next up:** Error Handling with Try/Catch — making your code resilient when things go wrong.
