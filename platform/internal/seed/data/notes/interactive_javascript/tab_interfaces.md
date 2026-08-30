# Tab Interfaces

A tabbed interface shows one content panel at a time, switched by clicking a row of buttons. It is a small state machine with clean rules: exactly one tab is active, and the visible panel matches the active tab. This lesson builds a fully accessible tab component.

## The state — one string

The entire tab system is driven by a single piece of state:

```js
let activeTab = "home";
```

That string determines which button looks selected and which panel is visible. Everything else is derived.

## HTML structure

```html
<div class="tabs" role="tablist">
  <button role="tab" data-tab="home" aria-selected="true" id="tab-home" aria-controls="panel-home">Home</button>
  <button role="tab" data-tab="profile" aria-selected="false" id="tab-profile" aria-controls="panel-profile">Profile</button>
  <button role="tab" data-tab="settings" aria-selected="false" id="tab-settings" aria-controls="panel-settings">Settings</button>
</div>

<div role="tabpanel" id="panel-home" data-panel="home" aria-labelledby="tab-home">
  <p>Welcome home.</p>
</div>
<div role="tabpanel" id="panel-profile" data-panel="profile" aria-labelledby="tab-profile" hidden>
  <p>Your profile info.</p>
</div>
<div role="tabpanel" id="panel-settings" data-panel="settings" aria-labelledby="tab-settings" hidden>
  <p>App settings.</p>
</div>
```

The ARIA roles are not decoration — they tell screen readers this is a tab interface:
- `role="tablist"` on the container
- `role="tab"` on each button
- `role="tabpanel"` on each content section
- `aria-selected` on the active tab
- `aria-controls` / `aria-labelledby` to link tabs and panels

## CSS

```css
[role="tabpanel"][hidden] { display: none; }

[role="tab"] {
  padding: 0.5rem 1rem;
  border: none;
  background: transparent;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}

[role="tab"][aria-selected="true"] {
  border-bottom-color: teal;
  font-weight: bold;
}
```

## The render function

```js
const tabs = document.querySelectorAll("[role='tab']");
const panels = document.querySelectorAll("[role='tabpanel']");

function render() {
  tabs.forEach(tab => {
    const isActive = tab.dataset.tab === activeTab;
    tab.setAttribute("aria-selected", isActive);
    tab.tabIndex = isActive ? 0 : -1;
  });

  panels.forEach(panel => {
    panel.hidden = panel.dataset.panel !== activeTab;
  });
}
```

:::key
`tabIndex = -1` on inactive tabs removes them from the Tab key sequence. Only the active tab is focusable with Tab, while arrow keys navigate between tabs. This matches the WAI-ARIA tab pattern.
:::

## Click handling with delegation

One listener on the tab container handles every tab:

```js
const tabContainer = document.querySelector("[role='tablist']");

tabContainer.addEventListener("click", (event) => {
  const tab = event.target.closest("[role='tab']");
  if (!tab) return;

  activeTab = tab.dataset.tab;
  render();
  tab.focus();
});
```

No need to attach individual listeners or query each tab's state — delegation and the `data-tab` attribute handle everything.

## Keyboard navigation

The WAI-ARIA authoring practices specify that arrow keys should move between tabs:

```js
tabContainer.addEventListener("keydown", (event) => {
  const tabElements = Array.from(tabs);
  const currentIndex = tabElements.findIndex(t => t.dataset.tab === activeTab);

  let nextIndex;

  if (event.key === "ArrowRight") {
    nextIndex = (currentIndex + 1) % tabElements.length;
  } else if (event.key === "ArrowLeft") {
    nextIndex = (currentIndex - 1 + tabElements.length) % tabElements.length;
  } else if (event.key === "Home") {
    nextIndex = 0;
  } else if (event.key === "End") {
    nextIndex = tabElements.length - 1;
  } else {
    return; // not a key we handle
  }

  event.preventDefault();
  activeTab = tabElements[nextIndex].dataset.tab;
  render();
  tabElements[nextIndex].focus();
});
```

The modulo math wraps around: pressing Right on the last tab goes to the first, pressing Left on the first goes to the last.

:::tip
Home jumps to the first tab, End jumps to the last. These are standard keyboard shortcuts in the tab pattern and are easy to add once you have the arrow key logic.
:::

## The full component

Here is everything wired together:

```js
const tabContainer = document.querySelector("[role='tablist']");
const tabs = document.querySelectorAll("[role='tab']");
const panels = document.querySelectorAll("[role='tabpanel']");

let activeTab = "home";

function render() {
  tabs.forEach(tab => {
    const isActive = tab.dataset.tab === activeTab;
    tab.setAttribute("aria-selected", isActive);
    tab.tabIndex = isActive ? 0 : -1;
  });

  panels.forEach(panel => {
    panel.hidden = panel.dataset.panel !== activeTab;
  });
}

tabContainer.addEventListener("click", (event) => {
  const tab = event.target.closest("[role='tab']");
  if (!tab) return;
  activeTab = tab.dataset.tab;
  render();
  tab.focus();
});

tabContainer.addEventListener("keydown", (event) => {
  const tabElements = Array.from(tabs);
  const currentIndex = tabElements.findIndex(t => t.dataset.tab === activeTab);
  let nextIndex;

  switch (event.key) {
    case "ArrowRight": nextIndex = (currentIndex + 1) % tabElements.length; break;
    case "ArrowLeft":  nextIndex = (currentIndex - 1 + tabElements.length) % tabElements.length; break;
    case "Home":       nextIndex = 0; break;
    case "End":        nextIndex = tabElements.length - 1; break;
    default: return;
  }

  event.preventDefault();
  activeTab = tabElements[nextIndex].dataset.tab;
  render();
  tabElements[nextIndex].focus();
});

render(); // show the default tab
```

## Toggling .active class instead of hidden

An alternative to the `hidden` attribute is toggling a CSS class:

```css
.panel { display: none; }
.panel.active { display: block; }
```

```js
panels.forEach(panel => {
  panel.classList.toggle("active", panel.dataset.panel === activeTab);
});
```

Both approaches work. The `hidden` attribute is semantically clearer and works without any CSS, but `.active` classes give you more control over transitions and animations.

## Adding animation

A simple fade when switching panels:

```css
[role="tabpanel"] {
  opacity: 0;
  transition: opacity 0.2s ease;
}

[role="tabpanel"]:not([hidden]) {
  opacity: 1;
}
```

:::warning
The `hidden` attribute sets `display: none`, which blocks transitions. For animated tabs, use a class-based approach instead:

```css
.panel { opacity: 0; height: 0; overflow: hidden; transition: opacity 0.2s ease; }
.panel.active { opacity: 1; height: auto; }
```
:::

## Practice

:::quiz
Q: In an accessible tab interface, what key should move focus from the first tab to the second?
- Tab
- Enter
- ArrowRight *
- Space
E: In the WAI-ARIA tab pattern, arrow keys navigate between tabs within the tablist. The Tab key moves focus out of the tablist entirely, and Enter/Space activate the focused tab (though many implementations activate on focus).
:::

:::quiz
Q: The `activeTab` state is `"profile"`. What should the render function do to the Home tab button?
- Add the `active` class
- Set `aria-selected="true"`
- Set `aria-selected="false"` and `tabIndex = -1` *
- Remove it from the DOM
E: Only the active tab gets `aria-selected="true"` and `tabIndex = 0`. Inactive tabs get `aria-selected="false"` and `tabIndex = -1`, removing them from the Tab key sequence while keeping them navigable with arrow keys.
:::

## Recap

- Tab state is a **single string** — the id of the active tab.
- Use `data-tab` and `data-panel` attributes to link buttons to their content.
- **ARIA roles:** `tablist`, `tab`, `tabpanel`, `aria-selected`, `aria-controls`, `aria-labelledby`.
- Set **`tabIndex = -1`** on inactive tabs and `tabIndex = 0` on the active tab.
- Handle **arrow keys** for keyboard navigation; wrap around with modulo.
- Use **event delegation** — one click listener on the tablist container.
- Render by toggling `hidden` or an `.active` class based on the state.

**Next up:** Modal Dialogs — showing overlays with proper focus management.
