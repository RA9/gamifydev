# Operators and Expressions

An **expression** is any piece of code that produces a value. Operators are the symbols that combine, compare, and transform those values. Mastering them is non-negotiable — they show up in every condition, calculation, and assignment you'll ever write.

## Arithmetic Operators

The basics work exactly like math class:

```js
console.log(10 + 3);   // → 13   addition
console.log(10 - 3);   // → 7    subtraction
console.log(10 * 3);   // → 30   multiplication
console.log(10 / 3);   // → 3.333...  division
console.log(10 % 3);   // → 1    remainder (modulo)
console.log(10 ** 3);  // → 1000 exponentiation
```

The `%` (modulo) operator returns the remainder after division. It's surprisingly useful:

```js
// Check if a number is even or odd
console.log(7 % 2);   // → 1  (odd)
console.log(8 % 2);   // → 0  (even)

// Wrap around a fixed range (e.g., cycling through 3 items)
let index = 5 % 3;    // → 2
```

The `**` operator (ES2016) replaces `Math.pow()`:

```js
console.log(2 ** 10);         // → 1024
console.log(Math.pow(2, 10)); // → 1024  (same thing, older syntax)
```

## Assignment Operators

Assignment stores a value. Compound assignment operators do math and assign in one step:

```js
let score = 100;    // basic assignment

score += 10;   // score = score + 10  → 110
score -= 5;    // score = score - 5   → 105
score *= 2;    // score = score * 2   → 210
score /= 3;    // score = score / 3   → 70
score %= 9;    // score = score % 9   → 7
score **= 2;   // score = score ** 2  → 49
```

Increment and decrement are also common:

```js
let count = 0;
count++;    // count is now 1
count--;    // count is now 0
```

## Comparison Operators

Comparisons return a boolean (`true` or `false`):

```js
console.log(5 > 3);    // → true
console.log(5 < 3);    // → false
console.log(5 >= 5);   // → true
console.log(5 <= 4);   // → false
console.log(5 === 5);  // → true   (strict equality)
console.log(5 !== 3);  // → true   (strict inequality)
```

## Why === Over ==

This is one of JavaScript's most important rules. `===` (strict equality) checks value **and** type. `==` (loose equality) converts types first, producing surprising results:

```js
console.log(5 === "5");   // → false  (number vs string — different types)
console.log(5 == "5");    // → true   (== converts "5" to 5 first)

console.log(0 == "");     // → true   😱
console.log(0 == false);  // → true   😱
console.log("" == false); // → true   😱

console.log(0 === "");    // → false  ✅ predictable
console.log(0 === false); // → false  ✅ predictable
```

:::key
Always use `===` and `!==`. The loose `==` does silent type coercion that introduces subtle bugs. There is almost never a reason to use `==` in modern JavaScript.
:::

## Logical Operators

Logical operators combine boolean expressions:

```js
console.log(true && true);    // → true   (AND: both must be true)
console.log(true && false);   // → false
console.log(true || false);   // → true   (OR: at least one true)
console.log(false || false);  // → false
console.log(!true);           // → false  (NOT: flips it)
```

In practice, you combine them with comparisons:

```js
let age = 25;
let hasID = true;

if (age >= 18 && hasID) {
  console.log("Entry allowed");
}

if (age < 13 || age > 65) {
  console.log("Discounted ticket");
}
```

### Short-Circuit Evaluation

`&&` and `||` don't always return `true`/`false` — they return the value that decided the outcome:

```js
console.log("hello" && "world");   // → "world"  (both truthy → returns last)
console.log(null && "world");      // → null      (first is falsy → returns it)
console.log("" || "default");      // → "default" (first falsy → returns second)
console.log("value" || "default"); // → "value"   (first truthy → returns it)
```

This enables the classic default-value pattern:

```js
let username = inputName || "Guest";
```

:::warning
The `||` default pattern treats `0`, `""`, and `false` as "missing." If those are valid values, use `??` (nullish coalescing) instead.
:::

## Nullish Coalescing (??)

The `??` operator (ES2020) returns the right side only when the left is `null` or `undefined` — not just any falsy value:

```js
let count = 0;
console.log(count || 10);   // → 10   (0 is falsy — wrong!)
console.log(count ?? 10);   // → 0    (0 is not null/undefined — correct!)

let name = "";
console.log(name || "Guest");  // → "Guest" (empty string is falsy)
console.log(name ?? "Guest");  // → ""      (empty string is not null/undefined)

let missing = undefined;
console.log(missing ?? "fallback");  // → "fallback"
```

:::key
Use `??` when `0`, `""`, or `false` are valid values. Use `||` when any falsy value should trigger the fallback.
:::

## Optional Chaining (?.)

Accessing a property on `null` or `undefined` throws an error. Optional chaining (ES2020) short-circuits to `undefined` instead:

```js
let user = { profile: { name: "Ada" } };
console.log(user.profile.name);     // → "Ada"
console.log(user.settings?.theme);  // → undefined (no crash)

let guest = null;
console.log(guest?.profile?.name);  // → undefined (no crash)
// Without ?.: guest.profile.name → TypeError!
```

It works with methods and bracket notation too:

```js
let arr = null;
console.log(arr?.[0]);         // → undefined
console.log(arr?.length);      // → undefined

let obj = {};
console.log(obj.doStuff?.());  // → undefined (method doesn't exist)
```

:::tip
Combine `?.` with `??` for safe access with a fallback: `user.settings?.theme ?? "light"` gives you `"light"` when the path doesn't exist.
:::

## The Ternary Operator

A compact `if/else` that produces a value: `condition ? ifTrue : ifFalse`.

```js
let age = 20;
let label = age >= 18 ? "adult" : "minor";
console.log(label);  // → "adult"

// Often used inline
console.log(`Status: ${score > 50 ? "pass" : "fail"}`);
```

:::warning
Don't nest ternaries. `a ? b ? c : d : e` is unreadable. Use `if/else` when logic has more than two branches.
:::

## Operator Precedence Basics

JavaScript evaluates operators in a specific order, just like math:

```js
console.log(2 + 3 * 4);     // → 14  (multiplication first)
console.log((2 + 3) * 4);   // → 20  (parentheses override)

console.log(!true || false); // → false  (! runs first, then ||)
console.log(!(true || false)); // → false (parentheses change meaning)
```

The practical rule:

1. `()` — parentheses (highest)
2. `!`, `**` — unary NOT, exponentiation
3. `*`, `/`, `%` — multiplication, division, remainder
4. `+`, `-` — addition, subtraction
5. `<`, `>`, `<=`, `>=` — comparisons
6. `===`, `!==` — equality
7. `&&` — logical AND
8. `||`, `??` — logical OR, nullish coalescing
9. `=`, `+=`, `-=` — assignment (lowest)

:::tip
When in doubt, use parentheses. `(a > 5) && (b < 10)` is clearer than `a > 5 && b < 10`, even though they're equivalent.
:::

## Type Coercion Traps

JavaScript auto-converts types in certain operations. The `+` operator is the biggest offender because it doubles as string concatenation:

```js
console.log(2 + 2);       // → 4      (both numbers)
console.log(2 + "2");     // → "22"   (string wins → concatenation)
console.log("5" - 1);     // → 4      (- forces numbers)
console.log("5" * 2);     // → 10     (* forces numbers)
console.log(true + 1);    // → 2      (true becomes 1)
console.log(null + 5);    // → 5      (null becomes 0)
```

Defend against it by converting explicitly:

```js
let input = "42";
let num = Number(input);     // → 42
let str = String(100);       // → "100"
let bool = Boolean("");      // → false
```

:::quiz
Q: What does `0 ?? 10` evaluate to?
- 10
- 0 *
- undefined
E: `??` only falls through on `null` or `undefined`. Since `0` is neither, it returns `0`. This is the key difference from `||`, which would return `10`.
:::

:::quiz
Q: What does `user.address?.city` do when `user.address` is `undefined`?
- Throws a TypeError
- Returns `undefined` without crashing *
- Returns `null`
E: Optional chaining (`?.`) short-circuits to `undefined` when the left side is `null` or `undefined`, preventing the TypeError you'd get with regular dot access.
:::

:::quiz
Q: Why should you prefer `===` over `==`?
- `===` is faster
- `===` checks both value and type, avoiding silent type coercion *
- `==` has been deprecated
E: `==` performs type coercion before comparing, which leads to surprises like `0 == ""` being `true`. `===` is predictable — it never converts types.
:::

## Recap

- Arithmetic: `+`, `-`, `*`, `/`, `%` (remainder), `**` (exponent).
- Assignment: `=` and compound forms (`+=`, `-=`, `*=`, etc.).
- Comparison: always use `===` and `!==` — never loose `==`.
- Logical: `&&` (AND), `||` (OR), `!` (NOT) — they short-circuit and return the deciding value.
- Nullish coalescing (`??`) falls through only on `null`/`undefined`, unlike `||` which treats all falsy values as missing.
- Optional chaining (`?.`) safely accesses nested properties without crashing on `null`/`undefined`.
- Ternary (`? :`) is a concise inline conditional — don't nest it.
- When in doubt about precedence, use parentheses.
- Watch for `+` with strings — it concatenates. Convert explicitly with `Number()`, `String()`, `Boolean()`.

**Next up:** Control Flow: Conditionals — making your programs branch and decide.
