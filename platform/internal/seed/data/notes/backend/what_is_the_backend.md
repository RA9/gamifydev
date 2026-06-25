# What is the Backend?

You've seen the **frontend** — the part of an app you can see and click. But when you log in, post a comment, or check your balance, something behind the scenes remembers your data and enforces the rules. That something is the **backend**.

By the end of this lesson you'll understand what a backend does and how the pieces — browser, server, and database — work together.

## Frontend vs backend

- The **frontend** runs in *your browser*. It's what you see: buttons, text, layout. (HTML, CSS, JavaScript.)
- The **backend** runs on a *server* — a computer somewhere else. It stores data, checks passwords, and makes decisions. You never see it directly.

:::analogy
A restaurant is the perfect picture. The **frontend** is the dining room — menus, tables, the food on your plate. The **backend** is the kitchen — where the real work happens, out of sight. You only see the result.
:::

## The request–response cycle

The whole web runs on a simple conversation: the browser **asks**, the server **answers**.

![How the backend works: browser, server, and database](/images/lessons/client-server.svg)

1. You click something. The **browser** sends a **request** to the server (e.g. "give me my to-do list").
2. The **server** runs your backend code. It may **query the database** for data.
3. The server sends a **response** back — usually data, which the frontend then displays.

This happens constantly, often many times per page.

:::quiz
Q: Where does backend code run?
- In the user's browser
- On a server — a separate computer *
- Inside the CSS file
E: The backend runs on a server, not in the browser. The frontend (browser) sends requests to it and shows the responses.
:::

## What a backend is responsible for

Backends typically handle:

- **Data storage** — saving and retrieving information in a database.
- **Business logic** — the rules ("only the owner can delete this post").
- **Authentication** — who is this user, and what are they allowed to do?
- **Talking to other services** — payments, email, maps, and more.

:::key
The **backend** runs on a server and handles data, rules, and security. The web works as a **request → response** cycle: the browser asks, the server (often using a database) answers.
:::

## The tools you'll meet

To build a backend you need three things, and you'll learn each next:

1. A **programming language** to write the logic — we'll use **Python**.
2. A **database** to store data — we'll use **SQL**.
3. A way to expose your data to the frontend — a **REST API**.

## Talk about it

Explain it out loud:

> "What's the difference between the frontend and the backend, and what are the three steps of the request–response cycle?"

If the restaurant analogy makes it click, you're ready to start writing backend code.

## What's next

Next, **Programming with Python** — the friendly language most beginners use to build their first backend.
