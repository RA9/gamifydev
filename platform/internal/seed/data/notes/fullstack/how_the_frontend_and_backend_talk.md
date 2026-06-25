# How the Frontend and Backend Talk

A frontend and a backend are two separate programs, often on two different computers. So how do they communicate? Through **HTTP requests** — and on the frontend, the tool for making them is `fetch`.

By the end of this lesson you'll be able to read JavaScript that calls an API and uses the response.

## The conversation, recapped

The frontend **asks**, the backend **answers**, and the data travels as **JSON**.

![How the backend works: request and response](/images/lessons/client-server.svg)

Your job on the frontend is to *send the request* and *do something with the response*.

## Fetching data with `fetch`

`fetch` is JavaScript's built-in way to call an API. It's **asynchronous** — the response takes time to arrive, so we wait for it with `await`:

```js
async function loadTodos() {
  const response = await fetch("https://api.example.com/todos");
  const todos = await response.json();   // parse the JSON body
  console.log(todos);
}
```

Two `await`s, two waits: one for the server to respond, one to read the JSON out of the response. The result is plain JavaScript data you can drop into the page.

:::analogy
`fetch` is like texting a question to a friend. You send it (`fetch`), then *wait* for the reply (`await`) — you don't freeze staring at your phone; the rest of life continues until the answer pings back. That "carry on until the reply arrives" is what *asynchronous* means.
:::

## Sending data (POST)

To create something, send a request with a method, headers, and a JSON body:

```js
async function addTodo(task) {
  const response = await fetch("https://api.example.com/todos", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ task: task })
  });
  return await response.json();
}
```

`JSON.stringify` turns a JavaScript object into a JSON string to send; the server turns it back into data. Same JSON you saw on the backend — now you're on the other end of it.

:::quiz
Q: Why do API calls use `await`?
- To make the code run faster
- Because the response takes time to arrive, so we wait for it *
- To convert HTML into CSS
E: Network requests are asynchronous — they take time. `await` pauses until the response (and then the JSON) is ready, without freezing the page.
:::

## Putting it together

The full loop on the frontend looks like this:

```js
async function showTodos() {
  const todos = await loadTodos();           // 1. fetch from the API
  const list = document.getElementById("list");
  list.innerHTML = todos
    .map(t => "<li>" + t.task + "</li>")     // 2. build HTML from the data
    .join("");                                // 3. put it on the page
}
```

Fetch the data → turn it into HTML → update the DOM. That's how a frontend displays anything that lives in a backend.

:::key
The frontend and backend talk over **HTTP**, exchanging **JSON**. Use **`fetch`** with **`await`** to call an API, `response.json()` to read the data, and `JSON.stringify` to send it. Then render the result into the DOM.
:::

## Talk about it

Explain out loud:

> "What does `fetch` do, and why do we need `await` when calling an API?"

## What's next

You know how the two halves connect. Now prove it — next is **Project: Build a Notes App**, a full-stack app from frontend to data.
