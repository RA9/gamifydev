# DevTools: Sources and Breakpoints

`console.log` is a flashlight — it illuminates one spot at a time. Breakpoints are floodlights — they freeze everything so you can look at the entire scene. This lesson teaches you to debug with the Sources panel.

## Opening the Sources panel

Open DevTools (F12 or Cmd + Option + I), then click the **Sources** tab.

The panel has three areas:

1. **File navigator** (left) — a tree of all loaded files
2. **Code editor** (center) — the source of the selected file
3. **Debugger sidebar** (right) — scope, watch, call stack, breakpoints

## Finding your files

In the file navigator, your project files appear under the domain (e.g. `localhost:5500` or `127.0.0.1`). Expand the tree to find your `.js` files.

**Shortcut:** Press **Cmd + P** (Mac) or **Ctrl + P** (Windows) to open a quick-search dialog. Type a filename to jump directly to it.

:::tip
If your files aren't showing up, refresh the page with DevTools open. The Sources panel only captures files loaded during the current session.
:::

## Setting line breakpoints

Click a line number in the editor to set a breakpoint — a blue marker appears. The next time that line runs, execution **pauses** and you can inspect everything.

```js
function handleAnswer(selectedIndex) {
  if (answered) return;          // line 1
  answered = true;               // line 2 ← set breakpoint here
  const q = questions[currentIndex];  // line 3
  // ...
}
```

Set a breakpoint on line 2. Now click a quiz option. The page freezes at that line. Hover over any variable — its current value appears in a tooltip.

To remove a breakpoint, click the blue marker again.

## Conditional breakpoints

Sometimes a line runs hundreds of times and you only care about a specific case. Right-click a line number and choose **Add conditional breakpoint**.

Enter a condition — the breakpoint only pauses when the condition is true:

```js
// Only pause when the user selects option 0
selectedIndex === 0
```

```js
// Only pause on the third question
currentIndex === 2
```

```js
// Only pause when score exceeds 3
score > 3
```

:::key
Conditional breakpoints are one of the most underused debugging tools. Instead of adding `if` checks and `console.log` statements to your code, set a conditional breakpoint and leave your source clean.
:::

## The debugger keyword

You can also trigger a breakpoint from your code by writing `debugger;` on its own line:

```js
function renderQuestion() {
  debugger;  // execution pauses here when DevTools is open
  const q = questions[currentIndex];
  // ...
}
```

When DevTools is open, the browser stops at this line exactly like a manual breakpoint. When DevTools is closed, `debugger;` is ignored.

:::warning
Remove `debugger;` statements before deploying. They won't break production (they're silently ignored without DevTools), but they're a sign of unfinished debugging if someone reviews your code.
:::

## Stepping controls

Once paused at a breakpoint, use the controls at the top of the debugger sidebar:

| Button | Name | What it does |
|--------|------|-------------|
| ▶️ | **Resume** | Continue running until the next breakpoint |
| ⤵️ | **Step over** | Run the current line and pause on the next line |
| ⬇️ | **Step into** | If the line calls a function, jump *inside* that function |
| ⬆️ | **Step out** | Finish the current function and pause when it returns |

**Step over** is your default move — it goes line by line.
**Step into** is for when you need to see what's happening inside a function call.
**Step out** gets you back up if you stepped into a function you don't care about.

```js
function nextQuestion() {
  currentIndex++;                    // Step over: see the index change
  if (currentIndex < questions.length) {
    renderQuestion();                // Step into: go inside renderQuestion
  } else {
    showResults();                   // Step into: go inside showResults
  }
}
```

## Watching expressions

The **Watch** section in the debugger sidebar lets you track specific expressions that update as you step through code.

Click **+** and add expressions:

```
currentIndex
score
questions[currentIndex].question
answered
questions.length - currentIndex
```

Every time you step to a new line, these values update in real time. You don't have to hover over variables — the answers are always visible.

:::tip
Add `questions[currentIndex]` to Watch when debugging a quiz. As you step through `renderQuestion`, you can verify the correct question object is being used at every point.
:::

## The Scope pane

When paused, the **Scope** pane shows every variable available at the current execution point, organized by scope:

- **Local** — variables declared inside the current function
- **Closure** — variables from enclosing functions that this function "closes over"
- **Global** — window-level variables

```text
▾ Local
    selectedIndex: 2
    q: {question: "...", options: [...], correct: 1}
    buttons: NodeList(4)
▾ Closure (handleAnswer)
    answered: true
    score: 3
    currentIndex: 2
▾ Global
    window: Window
```

This tells you exactly what values are in play — no guessing, no logging.

## The Call Stack pane

The **Call Stack** shows you the chain of function calls that led to the current line.

```text
handleAnswer          ← you are here
(anonymous)           ← the click event handler
```

Click any frame in the stack to jump to that function's context. This is invaluable when you're deep inside a call chain and need to understand *who called this function and with what arguments*.

:::key
The Call Stack answers "how did I get here?" — it shows the entire chain of function calls leading to the current breakpoint. When a function receives an unexpected value, the Call Stack tells you which caller passed it.
:::

## Event listener breakpoints

Sometimes you need to pause when an event fires, but you don't know which function handles it.

In the debugger sidebar, expand **Event Listener Breakpoints**. You'll see categories like:

- **Mouse** — click, mousedown, mouseup, mouseover
- **Keyboard** — keydown, keyup, keypress
- **DOM Mutation** — node inserted, node removed, attribute modified
- **Timer** — setTimeout, setInterval

Check **Mouse → click**. Now click anything on the page — the debugger pauses at the first line of whatever click handler runs. This is perfect for finding "what code runs when I click this button?"

## DOM breakpoints

Right-click any element in the Elements panel and select **Break on...**:

- **Subtree modifications** — pauses when children are added or removed
- **Attribute modifications** — pauses when an attribute changes (like a class being added)
- **Node removal** — pauses when this element is removed from the DOM

This is useful when something on the page changes and you don't know what JavaScript is causing it.

```text
Example: A class "hidden" keeps appearing on your quiz screen
and you don't know what's adding it.

→ Right-click the element in Elements panel
→ Break on → Attribute modifications
→ The debugger pauses on the exact line of JS that changes the class
```

## Practical debugging workflow

1. **Reproduce the bug** — make it happen consistently
2. **Set a breakpoint** near where you think the problem is
3. **Step through** line by line, checking values in Scope and Watch
4. **Narrow down** — once you find the wrong value, set a breakpoint *earlier* to find where it went wrong
5. **Fix the code** in your editor, save, and verify

:::quiz
Q: What is the main advantage of a conditional breakpoint over a regular breakpoint?
- It runs faster
- It only pauses when a specific condition is true, so you can target the exact case you're debugging *
- It highlights the line in a different color
- It works without DevTools open
E: When a line runs many times (like inside a loop or repeated function calls), a conditional breakpoint lets you pause only on the specific iteration or state you care about, skipping all the noise.
:::

:::quiz
Q: What does "Step into" do when paused at a function call?
- It skips the function entirely
- It jumps inside the called function so you can step through its code line by line *
- It runs the function backward
- It deletes the function from memory
E: "Step into" enters the function being called, letting you debug its internal logic. "Step over" would treat the function as a single step and move to the next line after it returns.
:::

## Recap

- **Sources panel** shows loaded files with a code editor and debugger sidebar.
- **Line breakpoints** (click a line number) pause execution at that line.
- **Conditional breakpoints** only pause when an expression evaluates to `true`.
- **`debugger;`** in source code triggers a breakpoint when DevTools is open.
- **Step over/into/out** let you navigate through code line by line.
- **Watch expressions** track specific values as you step through execution.
- **Scope pane** shows all variables available at the current pause point.
- **Call Stack** reveals the chain of function calls that led to the current line.
- **Event listener breakpoints** pause when specific events (click, keydown) fire.
- **DOM breakpoints** pause when an element's attributes, children, or existence change.

**Next up:** DevTools Network and Performance — inspecting requests, responses, and page speed.
