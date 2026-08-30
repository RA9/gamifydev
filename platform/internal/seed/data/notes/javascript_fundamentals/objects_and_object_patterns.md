# Objects and Object Patterns

Objects are JavaScript's way of grouping related data under named keys. A user has a name, email, and role. A product has a title, price, and stock count. Whenever you need labeled values instead of a positional list, you reach for an object.

## Object Literals

Create an object with curly braces and `key: value` pairs:

```js
let player = {
  name: "Ada",
  level: 5,
  isOnline: true,
};

console.log(player);
// → { name: "Ada", level: 5, isOnline: true }
```

Keys are strings by default (quotes are optional if the key is a valid identifier). Values can be any type — strings, numbers, booleans, arrays, other objects, even functions.

## Dot vs Bracket Notation

Two ways to access properties:

```js
let user = { name: "Ada", age: 28 };

// Dot notation — the usual way
console.log(user.name);     // → "Ada"

// Bracket notation — for dynamic keys or special characters
console.log(user["age"]);   // → 28
```

Brackets are essential when the key is stored in a variable:

```js
let field = "name";
console.log(user[field]);   // → "Ada"
console.log(user.field);    // → undefined (looks for a literal key "field")
```

:::key
Use dot notation by default. Use bracket notation when the key is dynamic (in a variable) or contains spaces/special characters.
:::

## Adding, Updating, and Removing Properties

Objects are mutable — you can change them after creation:

```js
let product = { name: "Pen", price: 2 };

// Update
product.price = 2.50;

// Add new property
product.inStock = true;
product.tags = ["office", "writing"];

// Remove a property
delete product.tags;

console.log(product);
// → { name: "Pen", price: 2.5, inStock: true }
```

## Methods

A property whose value is a function is called a **method**. Use the shorthand syntax:

```js
let counter = {
  count: 0,

  increment() {
    this.count++;
  },

  decrement() {
    this.count--;
  },

  reset() {
    this.count = 0;
  },
};

counter.increment();
counter.increment();
counter.increment();
console.log(counter.count);  // → 3

counter.decrement();
console.log(counter.count);  // → 2
```

### `this` in Methods (Briefly)

Inside a regular method, `this` refers to the object the method was called on:

```js
let user = {
  name: "Ada",
  greet() {
    console.log(`Hi, I'm ${this.name}`);
  },
};

user.greet();  // → "Hi, I'm Ada"
```

:::warning
`this` can be tricky. It depends on *how* a function is called, not where it's defined. Arrow functions don't have their own `this`, and extracting a method into a variable loses the `this` binding. For now, just use `this` inside methods defined with shorthand syntax.
:::

## Nested Objects

Objects can contain other objects — this is how you model real-world data:

```js
let company = {
  name: "Acme Corp",
  address: {
    street: "123 Main St",
    city: "Lagos",
    country: "Nigeria",
  },
  ceo: {
    name: "Ada",
    since: 2020,
  },
};

console.log(company.address.city);   // → "Lagos"
console.log(company.ceo.name);       // → "Ada"
```

Access deeply nested values carefully — if an intermediate property is `undefined`, you'll crash:

```js
// ❌ Crashes if company.address is undefined
console.log(company.address.zip);

// ✅ Safe with optional chaining
console.log(company.address?.zip);    // → undefined (no crash)
```

## Checking If a Key Exists

### The `in` Operator

```js
let config = { theme: "dark", fontSize: 14 };

console.log("theme" in config);     // → true
console.log("language" in config);  // → false
```

### hasOwnProperty

```js
console.log(config.hasOwnProperty("theme"));     // → true
console.log(config.hasOwnProperty("toString"));  // → false (inherited)
```

:::tip
`in` checks the entire prototype chain. `hasOwnProperty` checks only the object itself. For plain objects, both usually give the same result, but `hasOwnProperty` is more precise.
:::

### Checking for Undefined

```js
if (config.language !== undefined) {
  // key exists and has a value
}

// But beware — a key can exist with value undefined
let obj = { a: undefined };
console.log(obj.a !== undefined);  // → false
console.log("a" in obj);          // → true
```

## Object.keys, Object.values, Object.entries

These static methods give you arrays from an object — essential for iteration:

```js
let scores = { math: 90, science: 85, english: 78 };

console.log(Object.keys(scores));
// → ["math", "science", "english"]

console.log(Object.values(scores));
// → [90, 85, 78]

console.log(Object.entries(scores));
// → [["math", 90], ["science", 85], ["english", 78]]
```

## Looping Over Objects

Objects aren't iterable with `for...of`, but you have several options:

### for...in

```js
let config = { theme: "dark", fontSize: 14, lang: "en" };

for (let key in config) {
  console.log(`${key}: ${config[key]}`);
}
// → theme: dark, fontSize: 14, lang: en
```

### Object.keys + forEach

```js
Object.keys(config).forEach(key => {
  console.log(`${key}: ${config[key]}`);
});
```

### Object.entries + for...of

```js
for (let [key, value] of Object.entries(config)) {
  console.log(`${key}: ${value}`);
}
```

:::key
`Object.entries` with destructuring (`[key, value]`) is the most readable way to loop over an object's properties.
:::

## Copying Objects

### The Problem with Direct Assignment

Objects are assigned **by reference**, not by value:

```js
let original = { name: "Ada", level: 5 };
let copy = original;        // NOT a copy — same object!

copy.level = 99;
console.log(original.level); // → 99  (both point to the same object)
```

### Spread Syntax (Shallow Copy)

```js
let original = { name: "Ada", level: 5 };
let copy = { ...original };

copy.level = 99;
console.log(original.level); // → 5  (original is safe)
```

### Object.assign (Shallow Copy)

```js
let copy = Object.assign({}, original);
```

:::warning
Both spread and `Object.assign` create **shallow** copies. Nested objects are still shared references. If you need a deep copy, use `structuredClone(obj)` (modern browsers) or a library.
:::

```js
let user = {
  name: "Ada",
  settings: { theme: "dark" },
};

let shallow = { ...user };
shallow.settings.theme = "light";
console.log(user.settings.theme);  // → "light" (nested object was shared!)

let deep = structuredClone(user);
deep.settings.theme = "blue";
console.log(user.settings.theme);  // → "light" (truly independent copy)
```

## Practical Patterns

### Config Objects

```js
function createButton(options) {
  const label = options.label ?? "Click me";
  const color = options.color ?? "blue";
  const size = options.size ?? "medium";

  console.log(`[${size}] ${label} (${color})`);
}

createButton({ label: "Submit", color: "green" });
// → [medium] Submit (green)
```

### Lookup Maps

```js
const statusLabels = {
  pending: "⏳ Pending",
  active: "✅ Active",
  archived: "📦 Archived",
};

let currentStatus = "active";
console.log(statusLabels[currentStatus]);  // → "✅ Active"
```

### Counting Occurrences

```js
let words = ["apple", "banana", "apple", "cherry", "banana", "apple"];
let counts = {};

for (let word of words) {
  counts[word] = (counts[word] ?? 0) + 1;
}

console.log(counts);
// → { apple: 3, banana: 2, cherry: 1 }
```

:::quiz
Q: When should you use bracket notation instead of dot notation?
- Always — it's the preferred style
- When the property name is stored in a variable or contains special characters *
- Only when accessing nested objects
E: Bracket notation accepts any expression as the key, so it works with variables (`obj[key]`) and keys with spaces or special characters (`obj["full name"]`). Dot notation is simpler for static, valid-identifier keys.
:::

:::quiz
Q: What does `{ ...original }` create?
- A deep copy of the object
- A shallow copy of the object *
- A reference to the same object
E: Spread creates a new object with the same top-level properties. Nested objects are still shared references — it's a shallow copy, not a deep one.
:::

:::quiz
Q: What does `Object.entries(obj)` return?
- An array of the object's keys
- An array of the object's values
- An array of `[key, value]` pairs *
E: `Object.entries` returns an array of two-element arrays, each containing a key and its corresponding value. This makes it ideal for destructured iteration.
:::

## Recap

- Object literals group data under named keys: `{ name: "Ada", level: 5 }`.
- Use dot notation by default; bracket notation for dynamic keys or special characters.
- Add, update, or remove properties freely — objects are mutable.
- Methods are functions on objects. Inside a regular method, `this` refers to the object.
- Nested objects model real-world data. Use optional chaining (`?.`) for safe deep access.
- Check key existence with `in` or `hasOwnProperty`.
- `Object.keys`, `Object.values`, `Object.entries` convert objects to arrays for iteration.
- Objects are assigned by reference. Use spread (`{ ...obj }`) or `Object.assign` for shallow copies, `structuredClone` for deep copies.

**Next up:** Destructuring, Spread, and Rest — the modern syntax for unpacking and combining data.
