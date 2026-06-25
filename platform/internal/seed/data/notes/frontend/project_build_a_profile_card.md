# Project: Build a Profile Card

Enough theory — time to **build something real**. You've learned HTML structure and CSS styling; now you'll combine them into a polished profile card you can actually show people.

:::project
A clean, centered profile card with an avatar, a name, a short bio, and a button — built from scratch with just HTML and CSS. No installs, no frameworks: one file, opened in your browser.
:::

## Before you start

You need nothing but a text editor and a browser. Create a file called `index.html`, and after each step, **save it and refresh your browser** to watch it come together.

:::tip
Build in small steps and refresh often. Seeing each change land is the fastest way to learn — and to catch mistakes early.
:::

## Step 1 — The skeleton

Start with the basic page and a container for the card:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>My Profile Card</title>
    <style>
      /* our CSS will go here */
    </style>
  </head>
  <body>
    <div class="card">
      <img class="avatar" src="https://i.pravatar.cc/120" alt="My avatar" />
      <h1 class="name">Ada Lovelace</h1>
      <p class="bio">Aspiring developer learning one project at a time.</p>
      <button class="btn">Follow</button>
    </div>
  </body>
</html>
```

Open it now. It works — but it's plain. That's where CSS comes in.

## Step 2 — Center it on the page

Add this inside the `<style>` block. Flexbox centers the card in the middle of the screen:

```css
body {
  margin: 0;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  font-family: sans-serif;
}
```

Refresh — the card jumps to the centre. 🎯

## Step 3 — Style the card

Now make the card itself look like a card — padding, rounded corners, and a soft shadow:

```css
.card {
  background: white;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  text-align: center;
  max-width: 280px;
}
```

Remember the **box model**? `padding` is the breathing room inside; `border-radius` rounds the corners; the shadow lifts it off the page.

## Step 4 — Polish the details

```css
.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
}
.name { margin: 16px 0 4px; }
.bio { color: #64748b; }
.btn {
  margin-top: 16px;
  padding: 10px 24px;
  border: none;
  border-radius: 999px;
  background: #6740e8;
  color: white;
  font-weight: bold;
  cursor: pointer;
}
.btn:hover { background: #5731c8; }
```

`border-radius: 50%` turns the square photo into a circle, and `:hover` darkens the button when the mouse is over it. Refresh — you've got a real profile card. 🎉

:::quiz
Q: Which CSS turns the square avatar into a circle?
- `object-fit: cover`
- `border-radius: 50%` *
- `text-align: center`
E: `border-radius: 50%` rounds every corner by half the element's size, making a square into a circle.
:::

## You built it! 🚀

You just combined HTML structure with CSS styling to ship a real, polished component. That's exactly what front-end work is.

**Stretch goals** — make it yours:

- Swap in your own name, bio, and photo.
- Add a row of social links under the bio.
- Try a dark background and light card.

:::key
You built a complete UI component from scratch: structure in HTML, looks in CSS, centred with flexbox and finished with the box model. The same moves scale up to entire pages.
:::

## What's next

You can structure and style. Next, **JavaScript Basics** adds behaviour — and soon you'll build something interactive.
