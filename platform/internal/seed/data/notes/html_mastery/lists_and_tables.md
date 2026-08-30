# Lists and Tables

Lists and tables are how HTML structures grouped and tabular data. Using the right element gives screen readers the context they need ("list, 5 items" or "table, 3 columns, 10 rows") and makes your content machine-readable. In this lesson you'll learn every list type, build accessible tables, and understand when to reach for a table versus CSS layout.

## Unordered Lists

Use `<ul>` when the order of items doesn't matter — a set of features, a navigation menu, a shopping list.

```html
<ul>
  <li>Semantic HTML</li>
  <li>Accessible forms</li>
  <li>Responsive images</li>
</ul>
```

Browsers show bullet points by default. You can change or remove them with CSS — the semantics stay the same regardless of styling.

## Ordered Lists

Use `<ol>` when sequence matters — steps in a recipe, a ranked list, instructions.

```html
<ol>
  <li>Preheat the oven to 180°C.</li>
  <li>Mix flour, sugar, and eggs.</li>
  <li>Bake for 25 minutes.</li>
</ol>
```

Useful attributes on `<ol>`:

- **`start`** — begin numbering at a different value: `<ol start="5">`.
- **`reversed`** — count downward: `<ol reversed>`.
- **`type`** — change the marker: `1` (default), `a`, `A`, `i`, `I`.

```html
<ol type="a" start="3">
  <li>Third item (shown as "c")</li>
  <li>Fourth item (shown as "d")</li>
</ol>
```

:::tip
If you're debating between `<ul>` and `<ol>`, ask: "Would swapping two items change the meaning?" If yes, use `<ol>`. If the items are interchangeable, use `<ul>`.
:::

## Nested Lists

Lists can be nested inside other list items to create hierarchies.

```html
<ul>
  <li>Frontend
    <ul>
      <li>HTML</li>
      <li>CSS</li>
      <li>JavaScript</li>
    </ul>
  </li>
  <li>Backend
    <ul>
      <li>Node.js</li>
      <li>Python</li>
    </ul>
  </li>
</ul>
```

The inner `<ul>` sits inside the parent `<li>`, not after it. Screen readers announce the nesting: "list, 2 items. Item 1: Frontend. List, 3 items…"

:::warning
A common mistake is placing the nested list *outside* the `<li>`. The nested `<ul>` or `<ol>` must be a child of an `<li>` — never a direct child of the parent list.
:::

## Description Lists

Use `<dl>` for key-value pairs — glossaries, metadata, FAQs, or any list of terms and definitions.

```html
<dl>
  <dt>HTML</dt>
  <dd>HyperText Markup Language — the structure of web pages.</dd>

  <dt>CSS</dt>
  <dd>Cascading Style Sheets — the presentation layer.</dd>

  <dt>JavaScript</dt>
  <dd>A programming language for interactivity and logic.</dd>
</dl>
```

- **`<dt>`** — the term (description term).
- **`<dd>`** — the definition (description details).

A single term can have multiple definitions, and multiple terms can share one definition:

```html
<dl>
  <dt>a11y</dt>
  <dt>Accessibility</dt>
  <dd>Making websites usable by everyone, regardless of ability.</dd>
</dl>
```

:::tip
Description lists are underused. Anytime you see yourself building a two-column layout of labels and values (profile info, product specs, setting descriptions), consider a `<dl>`. It's semantic, accessible, and easy to style.
:::

## Tables: The Basics

Tables display two-dimensional data — spreadsheets, schedules, comparisons. A proper table uses these elements:

```html
<table>
  <caption>Q1 2025 Sales by Region</caption>
  <thead>
    <tr>
      <th scope="col">Region</th>
      <th scope="col">Revenue</th>
      <th scope="col">Growth</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>North America</td>
      <td>$1.2M</td>
      <td>+12%</td>
    </tr>
    <tr>
      <td>Europe</td>
      <td>$860K</td>
      <td>+8%</td>
    </tr>
    <tr>
      <td>Asia Pacific</td>
      <td>$640K</td>
      <td>+22%</td>
    </tr>
  </tbody>
</table>
```

### Element Breakdown

| Element | Purpose |
|---------|---------|
| `<table>` | The table container |
| `<caption>` | A visible title/description for the table |
| `<thead>` | Groups the header row(s) |
| `<tbody>` | Groups the body rows |
| `<tfoot>` | Groups footer rows (totals, summaries) |
| `<tr>` | A table row |
| `<th>` | A header cell (bold and centered by default) |
| `<td>` | A data cell |

:::key
Every data table should have a `<caption>` and use `<th>` with `scope` attributes. Without these, screen readers can't associate data cells with their headers, and the table becomes a wall of disconnected numbers.
:::

## The `scope` Attribute

`scope` tells assistive technology which cells a header applies to.

```html
<!-- Column header -->
<th scope="col">Price</th>

<!-- Row header -->
<th scope="row">Premium Plan</th>
```

- `scope="col"` — this header labels the column below it.
- `scope="row"` — this header labels the row to its right.

When a screen reader reaches a data cell, it reads the associated headers: "Revenue, Europe, $860K." Without `scope`, the user hears just "$860K" with no context.

## Table Footer

`<tfoot>` holds summary rows — totals, averages, notes.

```html
<table>
  <caption>Monthly Expenses</caption>
  <thead>
    <tr>
      <th scope="col">Category</th>
      <th scope="col">Amount</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Hosting</td>
      <td>$120</td>
    </tr>
    <tr>
      <td>Domain</td>
      <td>$15</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <th scope="row">Total</th>
      <td>$135</td>
    </tr>
  </tfoot>
</table>
```

## Spanning Rows and Columns

`colspan` and `rowspan` let a cell stretch across multiple columns or rows.

```html
<table>
  <caption>Class Schedule</caption>
  <thead>
    <tr>
      <th scope="col">Time</th>
      <th scope="col">Monday</th>
      <th scope="col">Tuesday</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>9:00 AM</td>
      <td colspan="2">Team Standup (both days)</td>
    </tr>
    <tr>
      <td>10:00 AM</td>
      <td>Design Review</td>
      <td rowspan="2">Deep Work Block</td>
    </tr>
    <tr>
      <td>11:00 AM</td>
      <td>Code Review</td>
    </tr>
  </tbody>
</table>
```

- `colspan="2"` makes a cell span two columns.
- `rowspan="2"` makes a cell span two rows.

:::warning
Complex spanning makes tables harder for screen readers. Keep spans to a minimum and test with a screen reader when you use them. If your table requires heavy spanning, consider whether it can be split into simpler tables.
:::

## Tables vs. CSS Grid/Flexbox

This is a question that comes up constantly: when should you use a table, and when should you use CSS layout?

**Use `<table>` for:**
- Data that has rows and columns with headers — spreadsheets, comparisons, schedules, statistics.
- Content where a screen reader should announce column/row headers as context.

**Use CSS Grid or Flexbox for:**
- Page layout (headers, sidebars, content areas).
- Card grids, image galleries, navigation bars.
- Anything that isn't genuinely tabular data.

```html
<!-- Correct: tabular data in a table -->
<table>
  <tr><th>Name</th><th>Role</th></tr>
  <tr><td>Alice</td><td>Engineer</td></tr>
</table>

<!-- Wrong: using a table for page layout -->
<table>
  <tr>
    <td>Navigation sidebar</td>
    <td>Main content</td>
  </tr>
</table>
```

:::key
If the content makes sense as a spreadsheet, it belongs in a `<table>`. If you're using a table just because it gives you columns, use CSS instead. Layout tables were common in the early web and are now considered an anti-pattern.
:::

:::quiz
Q: What is the purpose of the `scope` attribute on `<th>` elements?
- It limits which users can see the table
- It tells screen readers whether the header labels a column or a row *
- It controls the width of the column
E: `scope="col"` or `scope="row"` lets screen readers associate header cells with their data cells, so users hear context like "Revenue, Europe, $860K" instead of just "$860K."
:::

:::quiz
Q: When should you choose a `<dl>` (description list) over a `<ul>`?
- When the list has more than 10 items
- When items are key-value pairs or term-definition pairs *
- When you want numbered items
E: `<dl>` is designed for pairs — terms (`<dt>`) and their descriptions (`<dd>`). It's ideal for glossaries, metadata, specs, and FAQs. Use `<ul>` for simple sets and `<ol>` for ordered sequences.
:::

## Recap

- Use `<ul>` for unordered sets, `<ol>` for ordered sequences, and `<dl>` for term-definition pairs.
- Nested lists go inside `<li>`, not as siblings of list items.
- Data tables need `<caption>`, `<thead>`/`<tbody>`/`<tfoot>`, `<th>` with `scope`, and `<td>`.
- `colspan` and `rowspan` let cells span across rows or columns — use sparingly.
- Use `<table>` for tabular data only. Use CSS Grid or Flexbox for page layout.

**Next up:** Form Validation and Best Practices — making forms bulletproof with built-in browser validation.
