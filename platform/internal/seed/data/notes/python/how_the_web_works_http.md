# How the Web Works (HTTP)

Before you write a single line of web code, it helps to understand what actually happens when you open a website. Behind every page you load is a short, polite conversation between two computers. This lesson explains that conversation — the client, the server, requests, responses, and the language they speak: **HTTP**.

## The Client and the Server

Every interaction on the web has two sides. The **client** is whoever *asks* — usually a web browser like Chrome or Safari, but it could be a phone app or another program. The **server** is whoever *answers* — a computer, somewhere in the world, running software that waits for requests and sends back replies.

You type an address, hit Enter, and your browser (the client) sends a message to a server. The server reads it, figures out what you wanted, and sends something back — usually a web page. Then the connection closes. The whole web is millions of these tiny back-and-forth exchanges.

:::analogy
Think of ordering at a restaurant. You (the client) tell the waiter what you want. The kitchen (the server) prepares it and sends it back out. You don't walk into the kitchen — you just make a request and receive a response. The web works exactly like this, and HTTP is the language you both agree to speak.
:::

One crucial detail: the server does the *work*. It runs the code, looks things up in databases, and decides what to send. This is where **Python** lives on the web — on the server side, deciding how to answer each request. Your browser never runs your Python; the server does.

## HTTP: The Language of the Web

**HTTP** stands for **HyperText Transfer Protocol**. A "protocol" is just an agreed-upon format — a set of rules so both sides understand each other. When the client and server talk, they wrap their messages in this format.

Every exchange has exactly two parts:

- A **request**, sent by the client. ("Please give me the home page.")
- A **response**, sent by the server. ("Here it is, and everything went fine.")

You'll also see **HTTPS** — the same protocol with an **S** for **Secure**. It encrypts the conversation so nobody in between can read it. Every serious website uses HTTPS today, but the shape of a request and response is identical.

:::key
HTTP is *request-then-response*. The client always speaks first, the server always answers, and each exchange stands on its own. The server doesn't remember you between requests unless the app is specifically built to (with cookies or sessions) — a property called being **stateless**.
:::

## URLs: The Address of a Request

A **URL** (Uniform Resource Locator) is the address the client sends its request to. It looks like one string, but it's made of parts that each mean something:

```text
https://shop.example.com/products/42?color=blue
└─┬─┘   └──────┬───────┘└────┬─────┘└────┬────┘
scheme       host          path       query
```

- **Scheme** — `https` — which protocol to use.
- **Host** — `shop.example.com` — *which server* to contact.
- **Path** — `/products/42` — *what* you want from that server.
- **Query string** — `?color=blue` — extra details, written as `key=value` pairs joined by `&`.

When you build a web app, you'll spend a lot of time deciding which **paths** your app answers to — `/`, `/about`, `/products/42` — and what each one sends back. Those paths are called **routes**, and they're the heart of every web framework.

## HTTP Methods: The Verb of a Request

Every request carries a **method** (also called a *verb*) that says what kind of action you want. The path says *what* resource; the method says *what to do* with it. There are a handful you'll use constantly:

- **GET** — *read* something. Fetching a page or some data. GET requests should never change anything on the server; they just look.
- **POST** — *create* something new, or submit data. Signing up, posting a comment, placing an order.
- **PUT** — *replace* an existing resource entirely with new data.
- **PATCH** — *partially update* a resource — change just one or two fields.
- **DELETE** — *remove* a resource.

:::tip
The safe mental model: **GET reads, POST creates, PUT replaces, PATCH edits, DELETE removes.** When you're unsure, ask "does this change data on the server?" If yes, it should almost never be a GET.
:::

Here's why the distinction matters. Because GET only reads, browsers feel free to pre-load and cache GET requests. If you made "delete my account" a GET, a browser or search-engine crawler could trigger it just by visiting a link. Using the right verb keeps your app predictable and safe.

## What a Request and Response Look Like

Let's peek at the actual messages. They're just text. Here's a simple request asking a server for a page:

```text
GET /products/42 HTTP/1.1
Host: shop.example.com
Accept: text/html
User-Agent: Chrome/126.0
```

The first line is the important one: the **method** (`GET`), the **path** (`/products/42`), and the HTTP version. The lines below are **headers** — extra information about the request. Here they say "I want HTML back" and "I'm the Chrome browser."

The server reads that and sends back a **response**:

```text
HTTP/1.1 200 OK
Content-Type: text/html
Content-Length: 47

<html><body><h1>Blue Widget</h1></body></html>
```

The first line carries the **status code** (`200 OK`) — did it work? Then more **headers** describing the reply. Then a blank line, and finally the **body**: the actual content, in this case the HTML page.

:::example
A request and a response share the same shape: a first line, then headers, then a blank line, then an optional body. Requests lead with a *method and path*; responses lead with a *status code*. Learn to read that first line and you can read any HTTP message.
:::

## Status Codes: Did It Work?

The **status code** is a three-digit number the server sends back to summarize what happened. They're grouped by their first digit, and once you know the five groups, any code is easy to place:

- **2xx — Success.** The request worked.
  - **200 OK** — the standard "here's what you asked for."
  - **201 Created** — a POST succeeded and made a new resource.
- **3xx — Redirection.** "What you want is somewhere else."
  - **301 Moved Permanently** and **302 Found** — go to a different URL instead.
- **4xx — Client error.** *You* (the client) made a mistake.
  - **400 Bad Request** — the request was malformed.
  - **401 Unauthorized** — you need to log in.
  - **403 Forbidden** — you're logged in but not allowed.
  - **404 Not Found** — no resource at that path. The famous one.
  - **405 Method Not Allowed** — right path, wrong verb (e.g. POST to a GET-only route).
- **5xx — Server error.** The *server* broke while handling a valid request.
  - **500 Internal Server Error** — something crashed on the server side.

:::warning
A **4xx** means the client sent something wrong; a **5xx** means the server itself failed. Mixing these up sends debugging in the wrong direction. If you see a 404, check the URL you requested. If you see a 500, check the server's code and logs — the request was fine, your app threw an error.
:::

The pattern is worth memorizing: **2 = worked, 3 = go elsewhere, 4 = your fault, 5 = my fault.**

:::quiz
Q: A user POSTs a new comment and the server successfully saves it. Which status code best fits the response?
- 200 OK
- 201 Created*
- 404 Not Found
E: 201 Created is the precise code for "the request succeeded and made a new resource." A plain 200 works too, but 201 tells the client exactly what happened.
:::

## Headers and Body

We've mentioned headers and bodies — let's pin them down, because you'll set both when you build responses.

**Headers** are `Name: Value` lines carrying *information about* the message rather than the content itself. A few you'll meet often:

- `Content-Type` — what format the body is in: `text/html`, `application/json`, `image/png`.
- `Content-Length` — how many bytes the body is.
- `Authorization` — credentials proving who's making the request.

**The body** is the actual payload. On a *request*, the body carries data the client is sending up — the fields of a form you submitted, or a chunk of JSON. GET requests usually have no body (they're just asking). On a *response*, the body is what you came for: the HTML page, the JSON data, the image.

:::key
The **first line** answers "what action / did it work?", the **headers** answer "what format, how big, who are you?", and the **body** carries the actual content. Every HTTP message is built from these three ideas.
:::

## Where Python Fits

Now the picture comes together. When a request arrives at a server, some program has to read it, decide what to do, and build a response. That program is **your web app** — and you'll write it in Python.

A Python web framework like **Flask** handles the fiddly HTTP details for you. It parses the incoming request into a tidy object you can inspect, lets you match paths to functions, and turns whatever you return into a proper HTTP response with the right status code and headers. You focus on the logic — "when someone GETs `/products/42`, look up product 42 and send it back" — and Flask speaks HTTP on your behalf.

So the flow of a Python-powered web app looks like this:

```text
Browser  ──GET /products/42──►  Server (your Python/Flask app)
                                     │  runs your route function
                                     ▼
Browser  ◄──200 OK + HTML────  Server
```

Your code runs on the server, once per request, and its whole job is to turn a request into a response.

:::predict
Q: In a Python web app, where does your Python code actually run?
- In the visitor's browser, like JavaScript
- On the server, once for each incoming request*
- On both the browser and the server at the same time
E: Python runs on the *server*. For each request that arrives, the server runs your code to build a response and sends the result back. The browser never executes your Python.
:::

## Recap

- The web is a conversation between a **client** (asks) and a **server** (answers), spoken in **HTTP**.
- Every exchange is one **request** and one **response**; the client always speaks first, and each exchange is independent (**stateless**).
- A **URL** breaks into scheme, host, path, and query string; the **path** is the route your app answers.
- **Methods** are verbs: GET reads, POST creates, PUT replaces, PATCH edits, DELETE removes.
- **Status codes** group by first digit: 2xx success, 3xx redirect, 4xx client error, 5xx server error — with 200, 201, 404, and 500 the ones to know first.
- **Headers** describe the message (like `Content-Type`); the **body** carries the actual content.
- **Python runs on the server**, turning each request into a response — which is exactly what Flask helps you do.

**Next up:** Your First Web Server (Flask) — you'll write real Python that answers requests and sends back your very first response.
