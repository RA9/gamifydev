# What is Full-Stack Development?

You've met the two halves of the web: the **frontend** (what users see) and the **backend** (the server and data). A **full-stack developer** can build *both* — and, crucially, connect them into one working app.

By the end of this lesson you'll see the whole picture of how a complete application fits together.

## The whole stack

"The stack" is just all the layers of an app, from what the user touches down to where the data lives.

![Full-stack architecture: frontend, API, backend, database](/images/lessons/fullstack-architecture.svg)

- **Frontend** — HTML, CSS, JavaScript running in the browser. The face of the app.
- **API** — the bridge: requests and responses carrying **JSON** over **HTTP**.
- **Backend** — your server code (Python, Node…) with the logic and rules.
- **Database** — permanent storage for the data.

A full-stack developer is comfortable moving across all four.

:::analogy
Think of a full-stack developer as someone who can both *design the dining room* and *run the kitchen* — and make sure orders flow smoothly between them. You don't have to be the world's best at either to build something great end to end.
:::

## A feature, end to end

Watch one ordinary action travel through the whole stack — clicking "Add note":

1. **Frontend** — you type a note and click a button; JavaScript sends it to the API.
2. **API request** — a `POST /notes` with your note as JSON travels to the server.
3. **Backend** — the server validates it and saves it to the database.
4. **Database** — the note is stored permanently.
5. **Response** — the server replies "saved!", and the frontend updates the list.

Every feature in every app you use is some version of this round trip.

:::quiz
Q: In full-stack terms, what is "the stack"?
- Only the visual design of an app
- All the layers of an app, from frontend to database *
- A pile of programming books
E: The stack is every layer — frontend, API, backend, and database. A full-stack developer works across all of them.
:::

## You're closer than you think

Here's the encouraging part: you've already learned most of the pieces. HTML, CSS, and JavaScript on the front; the idea of servers, APIs, and databases on the back. Full-stack development is mostly about **connecting** them — which is exactly what this path focuses on.

:::key
**Full-stack** means working across the whole app: **frontend → API → backend → database**. Every feature is a round trip through those layers. The core skill is *connecting* the front and back you've already started learning.
:::

## Talk about it

Explain out loud:

> "What are the four layers of the stack, and what happens at each one when a user adds a note?"

## What's next

Next, **How the Frontend and Backend Talk** — the exact mechanics of connecting the two halves with `fetch` and JSON.
