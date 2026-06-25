# Profile Card

Almost every app on the web shows people: a face, a name, a short bio, a button to connect. Today you'll build that classic little card from scratch, and it'll look genuinely good.

:::project
You'll have a single HTML file that displays a polished, centered profile card with a circular avatar, a name, a bio, and a Follow button that reacts when you hover over it.
:::

To follow along you only need a text editor and a web browser. Create a file called `index.html`, and after each step save it and refresh the browser to watch your card take shape.

## Step 1 — Build the page skeleton

Start with a valid HTML document and drop in the card's markup: a container, an avatar image, a name, a short bio, and a Follow button. We're using `https://i.pravatar.cc/120` for a ready-made avatar, so you don't need an image file.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Profile Card</title>
</head>
<body>
  <div class="card">
    <img class="avatar" src="https://i.pravatar.cc/120" alt="Profile avatar" />
    <h1 class="name">Pixel Park</h1>
    <p class="bio">Frontend tinkerer. Coffee enthusiast. Probably debugging right now.</p>
    <button class="follow">Follow</button>
  </div>
</body>
</html>
```

## Step 2 — Center the card with flexbox

Right now the card hugs the top-left corner. Add a `<style>` block in the `<head>` and turn the whole page into a flex container so the card sits dead-center, both horizontally and vertically.

```html
<style>
  body {
    margin: 0;
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #f0f2f5;
    font-family: system-ui, sans-serif;
  }
</style>
```

## Step 3 — Style the card surface

Give the card itself a look: a white background, comfortable padding, rounded corners, a soft shadow, and centered text. This is what turns a stack of elements into a "card".

```css
.card {
  background: #ffffff;
  padding: 32px;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.12);
  text-align: center;
  width: 280px;
}
```

## Step 4 — Shape the avatar and text

Make the avatar a perfect circle with `border-radius: 50%`, and tidy up the spacing and color of the name and bio so they read cleanly.

```css
.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #f0f2f5;
}

.name {
  margin: 16px 0 4px;
  font-size: 22px;
  color: #1a1a1a;
}

.bio {
  margin: 0 0 24px;
  font-size: 14px;
  line-height: 1.5;
  color: #65676b;
}
```

## Step 5 — Style the button with a hover effect

Finish with a friendly Follow button. Remove the default border, give it color and padding, and add a `:hover` rule so it darkens when the mouse is over it. The `transition` makes that change feel smooth.

```css
.follow {
  background: #4f46e5;
  color: #ffffff;
  border: none;
  padding: 10px 28px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.follow:hover {
  background: #4338ca;
}
```
