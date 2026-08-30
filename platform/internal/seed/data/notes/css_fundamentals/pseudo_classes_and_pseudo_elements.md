# Pseudo-classes and Pseudo-elements

Pseudo-classes and pseudo-elements are what make CSS feel like it has superpowers. Pseudo-classes let you style elements based on *state* — is the user hovering? Is this the third item? Is the input valid? Pseudo-elements let you inject *visual content* without adding HTML. Together, they eliminate the need for many JavaScript interactions and extra markup.

## Pseudo-classes: Styling States

A pseudo-class starts with a single colon `:` and targets an element in a specific state.

### User Interaction States

#### `:hover`

Styles an element when the mouse pointer is over it:

```css
.btn {
  background: #2563eb;
  color: white;
}
.btn:hover {
  background: #1d4ed8;
}
```

:::warning
`:hover` doesn't work on touch devices — users can't "hover" with a finger. Never hide critical information or functionality behind hover-only interactions. Use hover to *enhance*, not to *enable*.
:::

#### `:focus`

Styles an element when it receives focus (via Tab key or click):

```css
input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
}
```

#### `:focus-visible`

Like `:focus`, but only when the focus is *visible* — typically keyboard navigation, not mouse clicks:

```css
button:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px;
}
```

This is the modern best practice: mouse users don't see a focus ring (they don't need it), but keyboard users do.

:::key
Use `:focus-visible` instead of `:focus` for focus rings. It gives keyboard users the visual indicator they need without cluttering the experience for mouse users.
:::

#### `:focus-within`

Matches a parent element when *any* of its descendants have focus:

```css
.search-bar:focus-within {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
}
```

This is perfect for highlighting a form group or container when the user is interacting with any input inside it.

#### `:active`

Styles the element while being clicked: `transform: scale(0.97)` gives a "press" feel.

#### `:visited`

Styles visited links. For privacy, only color-related properties can change.

### Link State Order

Style link pseudo-classes in **L**o**V**e **HA**te order to avoid specificity issues: `:link`, `:visited`, `:hover`, `:active`.

## Structural Pseudo-classes

### `:first-child` and `:last-child`

```css
li:first-child { font-weight: bold; }
li:last-child  { border-bottom: none; }
```

### `:nth-child()`

Target elements by position. Accepts a number, keyword, or formula:

```css
/* Specific position */
li:nth-child(3) { color: red; }

/* Even / odd */
tr:nth-child(even) { background: #f9fafb; }  /* zebra stripes */
tr:nth-child(odd)  { background: white; }

```css
/* Formula: every 3rd element */
li:nth-child(3n) { font-weight: bold; }
```

The formula `An+B` means: starting at position `B`, select every `A`th element.

### `:nth-of-type()`

Like `:nth-child()`, but counts only elements of the same *type*. `:nth-child(2)` means "the second child, if it's the right type." `:nth-of-type(2)` means "the second child of this type, regardless of other siblings."

### `:not()` — Negation

Exclude elements matching a selector:

```css
/* All inputs except submit buttons */
input:not([type="submit"]) {
  border: 1px solid #d1d5db;
}

/* All children except the last one get a bottom border */
.list-item:not(:last-child) {
  border-bottom: 1px solid #e5e7eb;
}
```

### `:empty`

Matches elements with no children (no text, no elements, no whitespace):

```css
.message:empty {
  display: none;
}

.container:empty::before {
  content: "Nothing here yet.";
  color: #9ca3af;
}
```

## Form Pseudo-classes

```css
input:required  { border-left: 3px solid #f59e0b; }
input:disabled  { background: #f3f4f6; cursor: not-allowed; opacity: 0.6; }
input:checked + label { font-weight: 700; color: #2563eb; }
input:valid     { border-color: #16a34a; }
input:invalid   { border-color: #dc2626; }
```

:::tip
Be careful with `:invalid` — empty required fields are invalid on page load, showing red borders immediately. Combine with `:not(:placeholder-shown)` to only show validation after the user types:
```css
input:invalid:not(:placeholder-shown):not(:focus) { border-color: #dc2626; }
```
:::

## Pseudo-elements: Generated Content

Pseudo-elements start with double colons `::` and create virtual elements you can style — without adding HTML.

### `::before` and `::after`

Insert content before or after an element's actual content:

```css
.required-label::after {
  content: " *";
  color: #dc2626;
}

.external-link::after {
  content: " ↗";
  font-size: 0.75em;
}
```

The `content` property is **required** — without it, the pseudo-element doesn't render. Use `content: ""` for decorative pseudo-elements.

#### Decorative Elements

```css
.section-title::before {
  content: "";
  display: inline-block;
  width: 4px;
  height: 1em;
  background: #2563eb;
  margin-right: 0.5rem;
  vertical-align: middle;
}
```

:::key
`::before` and `::after` are invisible to screen readers when used decoratively (with `content: ""`). If the content is meaningful (like a required-field asterisk), consider adding `aria-hidden="true"` and providing the information another way, or adding it as actual text.
:::

### `::placeholder`

Styles the placeholder text inside inputs:

```css
input::placeholder {
  color: #9ca3af;
  font-style: italic;
}
```

### `::selection`

Styles text that the user highlights with their cursor:

```css
::selection {
  background: #2563eb;
  color: white;
}
```

### `::marker`

Styles the bullet or number of a list item:

```css
li::marker {
  color: #2563eb;
  font-weight: bold;
}

ol li::marker {
  font-size: 1.2em;
}
```

## Practical UI Patterns

### Custom Checkbox

```css
.checkbox {
  appearance: none;
  width: 20px; height: 20px;
  border: 2px solid #d1d5db;
  border-radius: 4px;
  cursor: pointer;
  position: relative;
}
.checkbox:checked {
  background: #2563eb; border-color: #2563eb;
}
.checkbox:checked::after {
  content: "✓";
  position: absolute; top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  color: white; font-size: 14px;
}
```

### Separator lines between items

```css
.nav-item:not(:last-child)::after {
  content: "|";
  margin-left: 0.75rem;
  color: #d1d5db;
}
```

:::quiz
Q: What is the difference between `:focus` and `:focus-visible`?
- There is no difference
- `:focus` only works on inputs, `:focus-visible` works on all elements
- `:focus-visible` only shows when focus comes from keyboard navigation, not mouse clicks *
- `:focus-visible` has higher specificity than `:focus`
E: `:focus` triggers on any focus event (click, tab, programmatic). `:focus-visible` triggers only when the browser determines the user would benefit from a visible indicator — typically during keyboard navigation. This lets mouse users avoid seeing focus rings.
:::

:::quiz
Q: What is required for `::before` and `::after` pseudo-elements to render?
- A `display` property
- A `position` property
- A `content` property *
- A `width` and `height`
E: The `content` property is mandatory for `::before` and `::after` to appear. Without it, the pseudo-element doesn't exist. Use `content: ""` for decorative elements or `content: "text"` for visible text.
:::

## Recap

- **Pseudo-classes** (`:hover`, `:focus`, `:nth-child()`) style elements based on state or position.
- Use `:focus-visible` for focus rings — keyboard users see them, mouse users don't.
- `:focus-within` highlights parents when children have focus.
- `:nth-child()` and `:nth-of-type()` target by position with formulas.
- Form pseudo-classes (`:checked`, `:disabled`, `:required`) enable CSS-only form styling.
- **Pseudo-elements** (`::before`, `::after`) inject content without extra HTML. `content` is required.
- `::placeholder`, `::selection`, and `::marker` style parts of native elements.

**Next up:** CSS Debugging with DevTools — learning to diagnose and fix CSS problems efficiently.
