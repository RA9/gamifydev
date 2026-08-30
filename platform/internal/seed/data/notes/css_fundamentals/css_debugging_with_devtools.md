# CSS Debugging with DevTools

Every developer — beginner and expert alike — spends significant time debugging CSS. The difference between a frustrated beginner and a productive professional is knowing how to use DevTools. Browser DevTools are the most powerful CSS debugging instrument available, and learning them well will save you hours of trial and error.

This lesson walks you through a systematic debugging workflow using Chrome DevTools (Firefox and Safari have nearly identical features).

## Opening DevTools

- **Right-click → Inspect** on any element (opens directly to that element).
- **Keyboard**: `Cmd + Option + I` (Mac) or `Ctrl + Shift + I` (Windows/Linux).
- **Elements panel**: `Cmd + Option + C` (Mac) or `Ctrl + Shift + C` to enter inspect mode — hover over the page and click to select any element.

:::tip
Get in the habit of right-clicking the exact element you want to debug and choosing "Inspect." It jumps you straight to that element in the DOM tree, saving you from scrolling through HTML to find it.
:::

## The Elements Panel

The Elements panel shows two things side by side:

- **Left**: the live DOM tree (your HTML structure).
- **Right**: the CSS rules applied to the selected element.

### The Styles Pane

When you select an element, the Styles pane shows every CSS rule that targets it, in order of **specificity** (highest first):

```
element.style {                    ← inline styles
}

.card-title {                      ← your class rule
  font-size: 1.25rem;
  color: #111827;
}

h3 {                               ← element rule
  font-size: 1.17em;               ← struck through (overridden)
  font-weight: bold;
}

user agent stylesheet              ← browser defaults
h3 {
  display: block;
  font-size: 1.17em;               ← struck through
  margin-block-start: 1em;
}
```

**Struck-through** declarations are overridden by a higher-specificity rule. This is the most important diagnostic: you can instantly see *which rule won* and *which rules lost*.

## Tracing Why a Style Isn't Working

When a style doesn't apply, there are only a few possible causes. Check them in this order:

### 1. Is the rule reaching the element?

Select the element and search in the Styles pane. If your rule doesn't appear at all:
- The selector doesn't match (typo in class name, wrong nesting structure).
- The stylesheet isn't loaded (check the Network panel).
- The rule is inside a media query that isn't active.

### 2. Is the declaration overridden?

If the rule appears but the declaration is struck through:
- Another rule with higher specificity wins. The Styles pane shows the winner above the loser.
- Look for the winning rule and decide: should you increase your specificity (not ideal) or reduce the other rule's specificity?

### 3. Is the declaration invalid?

If the declaration has a yellow warning triangle or the property name is struck through with no competing rule:
- There's a syntax error (missing semicolon, invalid value).
- The property doesn't exist (typo).
- The value is invalid for that property.

### 4. Does the property apply to this element type?

Some properties only work in certain contexts:
- `margin: auto` only centers if the element has a defined width.
- `vertical-align` only works on inline/table-cell elements.
- `gap` only works inside flex/grid containers.
- `z-index` only works on positioned elements (not `static`).

:::key
When a CSS property "isn't working," 90% of the time it's one of four things: the selector doesn't match, a higher-specificity rule overrides it, the value is invalid, or the property doesn't apply to that element/context. Check them in that order.
:::

## Toggling Rules On and Off

Every declaration in the Styles pane has a **checkbox** on the left. Uncheck it to disable that declaration without deleting it. This lets you:

- See what an element looks like without a specific style.
- Isolate which rule is causing a layout problem.
- Test alternative values by disabling one and adding another.

This is faster than editing your CSS file, saving, and reloading.

## Editing Styles Live

### Change values

Click any value in the Styles pane and type a new one. Changes apply instantly. For numeric values, use **arrow keys** to increment/decrement:

- `↑`/`↓` — change by 1.
- `Shift + ↑`/`↓` — change by 10.
- `Alt + ↑`/`↓` — change by 0.1.

This is incredibly powerful for dialing in spacing, font sizes, and colors.

### Add new declarations

Click the empty space below the last declaration in a rule and start typing. DevTools autocompletes property names and values.

### Add a new rule

Click the **+** button at the top of the Styles pane to create a new rule. It pre-fills with a selector matching the selected element.

:::tip
Use the live editor to experiment freely. Once you find values that work, transfer them to your actual CSS file. Don't close DevTools before copying — changes are lost on reload.
:::

## Forcing Element States

Some styles only appear on `:hover`, `:focus`, or `:active` — but those states disappear the moment you move your mouse to DevTools.

Force a state: in the Styles pane, click **:hov** (toggle element state) and check the states you want to force:

- `:hover` — force hover styles.
- `:focus` — force focus styles.
- `:focus-visible` — force keyboard-focus styles.
- `:active` — force active/pressed styles.
- `:visited` — force visited-link styles.
- `:focus-within` — force focus-within styles.

The element stays in that state until you uncheck it, letting you inspect and edit hover/focus styles at your leisure.

## The Computed Tab

The **Computed** tab shows the *final, resolved values* for every CSS property on the selected element — after the cascade, specificity, and inheritance have all been applied.

This is where you go to answer: "What is the *actual* font-size/color/margin of this element right now?"

Each property shows:
- The final value.
- A disclosure arrow to expand and see *which rule* provided that value.

### The Box Model Diagram

At the top of the Computed tab is the interactive **box model diagram**:

```
┌─────── margin ───────┐
│  ┌──── border ────┐  │
│  │ ┌── padding ─┐ │  │
│  │ │  content   │ │  │
│  │ │  300 × 200 │ │  │
│  │ └────────────┘ │  │
│  └────────────────┘  │
└──────────────────────┘
```

Each layer shows its pixel values. Hover over a layer, and the browser highlights that exact area on the page. This is the definitive tool for debugging spacing issues.

## Spotting Margin Collapse

Margin collapse (covered in The Box Model lesson) is invisible in the Styles pane — the margin value *looks* correct. You'll spot it in the box model diagram or by hovering:

1. Select the first element — note its bottom margin (e.g. 24px).
2. Select the second element — note its top margin (e.g. 16px).
3. The actual gap between them is 24px (the larger one), not 40px.

You can confirm by hovering over each element — the margin highlights will *overlap* instead of stacking.

:::warning
If the gap between two elements isn't what you expect, always check for margin collapse. Select both elements one at a time and compare their margin highlights. If they overlap, the margins are collapsing.
:::

## Checking Specificity

When two rules conflict and you need to understand why one wins:

1. Find both rules in the Styles pane (the loser has a struck-through declaration).
2. Count each selector's specificity mentally: IDs → classes/attributes/pseudo-classes → elements.
3. The rule with the higher score wins. If tied, the one appearing later in source order wins.

DevTools in some browsers now show a specificity tooltip when you hover over a selector — look for it.

## Responsive Mode

Test how your site looks on different screen sizes:

- **Toggle device toolbar**: `Cmd + Shift + M` (Mac) or `Ctrl + Shift + M`.
- Choose a device preset (iPhone, iPad, Pixel) or set custom dimensions.
- Drag the viewport edges to resize freely.
- Test **orientation** (portrait/landscape).
- **Throttle network** and **CPU** to simulate slow mobile connections.

### Testing Media Queries

In responsive mode, a bar at the top shows colored segments for your media query breakpoints. Click a segment to jump to that viewport width.

## Step-by-Step Debugging Workflow

When something looks wrong in CSS, follow this systematic approach:

**Step 1: Identify the element.** Right-click → Inspect the exact element that looks wrong.

**Step 2: Check the box model.** Open the Computed tab and look at the box model diagram. Is the spacing what you expect? Are there unexpected margins, padding, or borders?

**Step 3: Check the Styles pane.** Look for struck-through declarations. These tell you another rule is winning. Find the winner and understand why.

**Step 4: Force states if needed.** If the problem is with hover, focus, or other interactive states, force those states with the :hov toggle.

**Step 5: Toggle declarations off.** Systematically disable declarations to isolate which one is causing the problem. When disabling a rule "fixes" the issue, you've found the culprit.

**Step 6: Experiment live.** Edit values in DevTools until it looks right, then copy those values to your actual CSS.

**Step 7: Check responsive.** If the issue only appears at certain sizes, use responsive mode to find the breakpoint where it breaks.

:::tip
Before asking "why isn't this working?", open DevTools and spend 60 seconds inspecting the element. The answer is almost always visible in the Styles pane or box model diagram. This habit alone will make you dramatically more productive.
:::

## Copying Styles

Right-click a rule in the Styles pane to:
- **Copy declaration** — copies the single property: value pair.
- **Copy rule** — copies the entire rule including the selector.
- **Copy all declarations** — copies every declaration in the block.

You can also right-click an element in the DOM tree and choose **Copy → Copy styles** to get all computed styles.

## Common Gotchas DevTools Reveals

| Symptom | What to look for in DevTools |
|---|---|
| Element wider than expected | Box model shows large padding/border, or check `box-sizing` |
| Margin not appearing | Look for margin collapse (overlapping margin highlights) |
| `z-index` not working | Check if element has `position: static` (the default) |
| Color wrong | Struck-through declaration — a higher-specificity rule wins |
| Hover style not visible | Force `:hover` state with the :hov toggle |
| Element invisible | Check `display: none`, `opacity: 0`, `visibility: hidden`, or zero dimensions |

:::quiz
Q: In the DevTools Styles pane, what does a struck-through declaration mean?
- The declaration has a syntax error
- The declaration is being overridden by a higher-specificity rule *
- The declaration is deprecated
- The declaration only applies on hover
E: A struck-through declaration means another rule with higher specificity (or later source order) is overriding it. Look above in the Styles pane to find the winning rule.
:::

:::quiz
Q: How can you inspect an element's `:hover` styles without the hover state disappearing when you move to DevTools?
- You can't — hover styles can't be inspected
- Use the :hov toggle in the Styles pane to force the hover state *
- Hold Shift while moving the mouse
- Type `:hover` in the Console
E: The :hov toggle in the Styles pane lets you force pseudo-class states like `:hover`, `:focus`, and `:active`. The state stays active even when your mouse leaves the element, letting you inspect and edit hover styles at your leisure.
:::

## Recap

- **Right-click → Inspect** jumps you directly to an element. Make this your reflex.
- The **Styles pane** shows every rule targeting an element, in specificity order. **Struck-through** = overridden.
- The **Computed tab** and **box model diagram** show final values and actual spacing.
- **Toggle declarations** on/off to isolate problems. **Edit live** to experiment, then copy to your CSS.
- **Force states** (`:hover`, `:focus`) with the :hov toggle to inspect interactive styles.
- Use **responsive mode** to test at different viewport sizes.
- Follow a systematic workflow: inspect → check box model → check styles → toggle → experiment → copy.

**Next up:** CSS Best Practices — writing CSS that's clean, scalable, and maintainable.
