# Handling Requests

Your routes can send responses — now let's make them *listen*. Real web apps react to what the client sends: a search term in the URL, the fields of a submitted form, a chunk of JSON from a mobile app. In this lesson you'll meet Flask's `request` object and learn to read every kind of incoming data, then send back clean JSON of your own.

## The request Object

Every time a request reaches your app, Flask packages all of its details — the URL, the headers, any submitted data — into a single object called `request`. You import it once and it's automatically filled in for whichever request is currently being handled:

```python
from flask import Flask, request

app = Flask(__name__)
```

Inside any route function, `request` refers to *this* request. It carries everything you might want to inspect:

- `request.method` — the HTTP verb (`"GET"`, `"POST"`, ...).
- `request.args` — values from the query string.
- `request.form` — fields from a submitted form.
- `request.get_json()` — a JSON body, parsed into a Python dict.

:::key
You don't create the `request` object — Flask does, fresh for each incoming request, and hands it to you ready to read. Just `from flask import request` and reach into it inside your route. It always reflects the request being handled right now.
:::

## Query Strings

Remember the query string from the URL lesson — the `?key=value` part on the end of a URL? Flask puts those pairs in `request.args`, which behaves like a dictionary. The safe way to read a value is `.get()`:

```python
@app.route("/greet")
def greet():
    name = request.args.get("name")
    return f"Hello, {name}!"
```

A request to `/greet?name=Ada` sets `name` to `"Ada"` and returns `Hello, Ada!`. The `?name=Ada` is read straight out of `request.args`.

Why `.get()` and not `request.args["name"]`? Because square-bracket access **crashes** with an error if the key is missing, while `.get()` politely returns `None`. Even better, you can supply a **default** for when the value isn't there:

```python
@app.route("/greet")
def greet():
    name = request.args.get("name", "stranger")
    return f"Hello, {name}!"
```

Now `/greet` with no query string returns `Hello, stranger!`, and `/greet?name=Ada` still returns `Hello, Ada!`. A default keeps your route from breaking when someone forgets a parameter.

:::warning
Everything in `request.args` arrives as a **string**, even numbers. `request.args.get("n")` for `?n=5` gives you the text `"5"`, not the integer `5`. If you need to do math, convert it yourself with `int(...)` — and be ready for the value to be missing or non-numeric.
:::

Here's that conversion in a route that adds two numbers from the query string:

```python
@app.route("/sum")
def add():
    a = int(request.args.get("a", 0))
    b = int(request.args.get("b", 0))
    return f"{a + b}"
```

A request to `/sum?a=3&b=4` reads `"3"` and `"4"`, converts each to an integer, and returns `7`. The `, 0` defaults mean a missing value counts as zero instead of crashing.

:::predict
Q: With the `add` route above, what does a request to `/sum?a=10&b=5` return?
- 105
- 15*
- error
E: Both values come in as strings (`"10"` and `"5"`), but `int(...)` converts them to real integers before adding, so `10 + 5` is `15`. Without `int(...)`, `"10" + "5"` would have concatenated into `"105"`.
:::

## Form Data

When a user submits an HTML `<form>` with the POST method, the fields don't ride in the URL — they travel in the request **body**. Flask parses them into `request.form`, which works just like `request.args`:

```python
@app.route("/signup", methods=["POST"])
def signup():
    username = request.form.get("username", "")
    return f"Welcome aboard, {username}!"
```

If a form posts a `username` field of `"Ada"`, this returns `Welcome aboard, Ada!`. The only difference from query strings is *where* the data lives: `request.args` reads the URL, `request.form` reads a posted form body. Both use the same friendly `.get("key", default)` pattern.

Notice `methods=["POST"]` on the route — submitting a form is a POST, since it sends data to be processed. We'll come back to restricting methods shortly.

## JSON Bodies

Forms are how browsers submit data, but mobile apps and other programs usually send **JSON** — the same key/value format you've seen everywhere. A JSON request body looks like this:

```text
POST /echo
Content-Type: application/json

{"message": "hi there"}
```

Flask reads and parses that body for you with `request.get_json()`, handing you an ordinary Python dictionary:

```python
@app.route("/echo", methods=["POST"])
def echo():
    data = request.get_json()
    message = data.get("message", "")
    return f"You said: {message}"
```

`request.get_json()` turns `{"message": "hi there"}` into the Python dict `{"message": "hi there"}`, so `data.get("message")` pulls out `"hi there"`. From here it's just normal dictionary access.

:::tip
`request.get_json()` expects the client to send a JSON body with the `Content-Type: application/json` header. If the body isn't valid JSON, it returns `None` — so a quick `if data is None:` check (or `request.get_json(silent=True)`) protects a route that might receive junk.
:::

## Returning JSON with jsonify

Reading JSON is half the story; often you want to *send* JSON back. You could never just `return` a Python dictionary as a raw string — the client wouldn't know it's JSON, and the format would be off. Flask gives you `jsonify` to do it properly:

```python
from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/echo", methods=["POST"])
def echo():
    data = request.get_json()
    message = data.get("message", "")
    return jsonify({"you_said": message})
```

`jsonify` takes a Python dict, converts it to a proper JSON string, and — importantly — sets the response's `Content-Type` header to `application/json`. Now a POST to `/echo` with `{"message": "hi"}` returns a real JSON response: `{"you_said": "hi"}`, correctly labelled so the client parses it as data.

:::key
Use **`request.get_json()`** to *read* an incoming JSON body into a Python dict, and **`jsonify(...)`** to *send* a Python dict back as a proper JSON response. One reads, one writes — and `jsonify` is what sets the right `Content-Type` so the client knows it's JSON.
:::

You can return a custom status alongside JSON using the same tuple trick from the last lesson:

```python
@app.route("/items", methods=["POST"])
def create_item():
    data = request.get_json()
    return jsonify({"created": data}), 201
```

That sends the JSON body *and* a `201 Created` status — exactly right for a route that makes something new.

## Restricting Methods

By default, a Flask route only answers **GET** requests. To accept other verbs, you list them — you saw `methods=["POST"]` above. Flask also offers shorthand decorators that read a little cleaner:

```python
@app.get("/status")
def status():
    return "all good"

@app.post("/echo")
def echo():
    data = request.get_json()
    return jsonify({"you_said": data.get("message", "")})
```

`@app.get(...)` means "GET only" and `@app.post(...)` means "POST only" — tidy equivalents to `@app.route(..., methods=[...])`. There are `@app.put`, `@app.patch`, and `@app.delete` too, matching the verbs from the HTTP lesson.

So what happens if a client uses the *wrong* verb — say, a GET request to a POST-only route? Flask handles it automatically: it responds with **405 Method Not Allowed**. You don't write any code for this; declaring the route's methods is enough, and Flask rejects everything else with the correct status.

:::example
Here's the full request-handling lab in one piece — a query-param greeting, a sum, and a POST endpoint that reads JSON and returns JSON:

```python
from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/greet")
def greet():
    name = request.args.get("name", "stranger")
    return f"Hello, {name}!"

@app.route("/sum")
def add():
    a = int(request.args.get("a", 0))
    b = int(request.args.get("b", 0))
    return f"{a + b}"

@app.post("/echo")
def echo():
    data = request.get_json()
    message = data.get("message", "")
    return jsonify({"you_said": message})
```

As always in these labs, the checker runs this app for you and sends the requests — you just write the routes.
:::

:::quiz
Q: A client sends a GET request to a route declared with `@app.post("/echo")`. What does Flask respond with?
- 200 OK
- 404 Not Found
- 405 Method Not Allowed*
E: The path `/echo` exists, but it only accepts POST. A GET to a POST-only route is the wrong verb for a valid path, so Flask automatically returns 405 Method Not Allowed.
:::

## Recap

- Flask fills a `request` object for every incoming request; `from flask import request` and read it inside your route.
- **Query strings** live in `request.args`; use `request.args.get("key", default)` — values arrive as strings, so `int(...)` them for math.
- **Form fields** live in `request.form`, read the same way; they travel in the body of a POST.
- **JSON bodies** are parsed by `request.get_json()` into a Python dict.
- Send JSON back with **`jsonify(...)`**, which formats the dict and sets `Content-Type: application/json`; add `, 201` for a custom status.
- Routes are GET-only by default; declare others with `methods=[...]` or shorthands like **`@app.post(...)`**.
- A wrong verb on a valid path yields an automatic **405 Method Not Allowed** — no extra code needed.

**Next up:** you'll combine everything — routes, dynamic paths, requests, and JSON — to build a small working API of your own.
