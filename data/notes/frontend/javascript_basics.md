# JavaScript Basics

HTML gives your page structure. CSS gives it presentation. **JavaScript gives it behavior**.

That means JavaScript is the part of the frontend that can:

- react to user input
- make decisions
- store and update data
- change what appears on screen without reloading the page

By the end of this lesson, you should understand the core programming building blocks that power browser interactivity.

## What JavaScript is doing on a page

JavaScript runs in response to things that happen:

- a user clicks a button
- a form is submitted
- text is typed into a field
- data is loaded

The basic mental model is:

> something happens → JavaScript runs → state or the page changes

That loop is the foundation of interactive frontend work.

:::key
JavaScript is not "magic page movement." It is logic reacting to events and updating data or UI.
:::

## Variables: storing data

Variables are named containers for values.

```js
const name = "Ada";
let score = 0;
```

Two keywords matter most:

- `const` — use this by default when the variable should not be reassigned
- `let` — use this when the value needs to change later

```js
let likes = 0;
likes = likes + 1;
```

Here `likes` changes over time, so `let` makes sense.

:::quiz
Q: Which keyword should usually be your default choice when the variable will not be reassigned?
- `var`
- `const` *
- `let`
E: `const` is the safest default. Switch to `let` only when the value really needs to change.
:::

## Common JavaScript data types

The most useful starter types are:

- **string** — text, like `"hello"`
- **number** — numeric values, like `42`
- **boolean** — `true` or `false`
- **array** — ordered list of values
- **object** — grouped named values

```js
const title = "Sprintboard";
const taskCount = 3;
const isOpen = true;
const tags = ["frontend", "html", "css"];
const user = { name: "Maya", role: "designer" };
```

## Arrays: ordered collections

Arrays are perfect for lists of related data.

```js
const tasks = ["Ship landing page", "Fix form labels", "Write test cases"];
```

You can access items by index:

```js
console.log(tasks[0]); // "Ship landing page"
```

You can also add new items:

```js
tasks.push("Improve mobile nav");
```

Arrays are a huge part of frontend work because so many UIs display lists.

## Objects: structured data

Objects store related data under named keys.

```js
const task = {
  text: "Ship landing page",
  done: false,
  priority: "high"
};
```

You can read values like this:

```js
console.log(task.text);
console.log(task.done);
```

Objects are useful because real UI items usually have more than one piece of information.

:::fill
Q: Complete the expression to read the `name` property from the `user` object.
`user.___`
- name *
- value
- key
E: Dot notation like `user.name` reads the named property from an object.
:::

## Functions: reusable behavior

Functions let you package up a set of steps and run them when needed.

```js
function greet(name) {
  return "Hello, " + name;
}

const message = greet("Ada");
```

This function:

- takes an input called `name`
- returns a new string
- can be reused with different values

Functions are how you name behavior clearly.

## Conditionals: making decisions

Conditionals let JavaScript choose between paths.

```js
const passwordLength = 10;

if (passwordLength >= 8) {
  console.log("Password is long enough");
} else {
  console.log("Password is too short");
}
```

Frontend UIs constantly make conditional decisions:

- show an error or success message
- open or hide a panel
- disable or enable a button

## Loops and repeated work

If you want to repeat a step for each item in a list, use a loop pattern.

A very common frontend pattern is `forEach`:

```js
const skills = ["HTML", "CSS", "JavaScript"];

skills.forEach(function (skill) {
  console.log(skill);
});
```

This is especially useful when rendering lists into the page.

## State: data that changes over time

In frontend work, people often talk about **state**.

State is just the data that describes what your UI currently knows.

Examples:

- how many likes a post has
- whether a modal is open
- which tab is active
- what a user typed into a field

```js
let count = 0;
```

That `count` variable is state.

If the user clicks a button, the state changes.
If the state changes, the UI should usually update too.

:::predict
```js
let score = 0;
score = score + 5;
score = score + 10;
console.log(score);
```
Q: What gets logged?
- 0
- 15 *
- "510"
- 510
E: The numeric additions happen in order: 0 → 5 → 15.
:::

## Putting the pieces together

A tiny interactive feature often combines all of this:

- a variable stores state
- a function changes that state
- a conditional checks something
- the page updates afterward

Example:

```js
let count = 0;

function increment() {
  count = count + 1;
  console.log(count);
}
```

This is small, but it contains the same ingredients bigger apps use.

## Common mistakes to avoid early

- using `let` for everything instead of preferring `const`
- mixing up arrays and objects
- forgetting that array indexes start at `0`
- trying to memorize syntax without understanding the data flow
- writing long code without breaking repeated behavior into functions

:::warning
If you don't know what your data looks like, JavaScript quickly feels confusing. Always ask: what values am I storing, and how should they change?
:::

## What good looks like

You should now be able to:

- explain what JavaScript adds to the frontend
- choose between `const` and `let`
- recognize strings, numbers, booleans, arrays, and objects
- write a basic function
- understand `if / else` decisions
- read simple state-changing code

## What's next

In **DOM Manipulation & Events**, you'll connect these programming basics to the actual browser — selecting elements, listening for events, and updating the page in response.