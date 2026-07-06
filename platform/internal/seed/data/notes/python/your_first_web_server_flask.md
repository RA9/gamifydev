# Your First Web Server (Flask)

You know the theory now — a client sends a request, a server sends back a response, and your Python runs on the server. This lesson turns that into working code. By the end you'll have written a real web server that answers several routes, greets people by name, and even sends back a custom status code.

## What Flask Is

**Flask** is a small, friendly Python framework for building web servers. Its job is to handle all the messy HTTP details — reading requests, matching URLs, formatting responses — so you can focus on *what to send back* for each path.

Flask is called a **micro-framework** because it starts tiny and stays out of your way. You don't need to learn a huge system to get going; a complete Flask app can be four lines long. As your project grows, Flask grows with it, but the core idea never changes: you map **routes** (URL paths) to **functions**, and each function returns a response.

:::analogy
Think of Flask as a receptionist for your server. Requests arrive at the front desk, and the receptionist checks the address on each one — "ah, this is for `/about`" — and hands it to the right person (your function) to handle. You write the people; Flask runs the front desk.
:::

## Installing Flask

Flask isn't built into Python, so you install it with **pip**, Python's package installer. In a terminal you'd run:

```text
pip install flask
```

That downloads Flask and everything it needs. You typically do this inside a **virtual environment** — an isolated folder for one project's packages — so different projects don't clash. Creating one looks like:

```text
python -m venv venv
source venv/bin/activate    # on Windows: venv\Scripts\activate
pip install flask
```

:::tip
In these labs the environment is already set up for you — Flask is installed and ready. You can focus entirely on writing the app. But knowing `pip install flask` matters the moment you build something on your own machine.
:::

## A Minimal App

Here is a complete, working Flask application. It's worth reading slowly, because every Flask app you ever write starts from this exact shape:

```python
from flask import Flask

app = Flask(__name__)

@app.route("/")
def home():
    return "Hello, GamifyDev!"
```

Four ideas, top to bottom:

1. `from flask import Flask` — bring in the framework.
2. `app = Flask(__name__)` — create your application. Passing `__name__` tells Flask where your code lives so it can find related files. You'll write this line, unchanged, in every app.
3. `@app.route("/")` — a **decorator** that says "when a request comes in for the path `/`, run the function below."
4. `def home(): return "Hello, GamifyDev!"` — the function that handles that route. Whatever it **returns** becomes the response body.

When someone visits `/`, Flask runs `home()`, takes the returned string, and sends it back as a `200 OK` response. You just answered an HTTP request in Python.

:::key
The pattern is always **route + function**. The `@app.route(path)` decorator connects a URL path to the function right below it, and that function's **return value** is what the visitor receives. Everything else in Flask builds on this.
:::

## Returning Strings

The simplest thing a route can return is a string, and that's often all you need. Add as many routes as you like — each is a decorator plus a function:

```python
from flask import Flask

app = Flask(__name__)

@app.route("/")
def home():
    return "Welcome to the GamifyDev server!"

@app.route("/ping")
def ping():
    return "pong"
```

Now the server answers two paths. A request to `/` returns the welcome message; a request to `/ping` returns `pong`. The `/ping` route is a classic **health check** — a tiny endpoint you can hit to confirm the server is alive and responding.

Each function name (`home`, `ping`) just has to be unique; Flask uses the **decorator's path**, not the function name, to route requests. Pick names that describe what the route does.

## Running the Dev Server

To see your app in a browser, Flask needs to actually *run* and listen for requests. There are two common ways.

The recommended way from a terminal is the `flask run` command:

```text
flask run
```

That starts a small development web server (usually at `http://127.0.0.1:5000`) that watches for requests and hands them to your app. The other way is to start it from inside your Python file by calling `app.run()`:

```python
if __name__ == "__main__":
    app.run(debug=True)
```

The `if __name__ == "__main__":` guard means "only start the server when this file is run directly." Setting `debug=True` gives you helpful error pages and auto-reloads the server when you save a change — great while developing.

:::tip
Good news for these labs: **you never have to start the server yourself.** The checker runs your app for you using a built-in test client, sends requests to your routes, and inspects the responses. Just write your routes and return the right thing — no `flask run`, no `app.run()`, no ports to worry about. Out in the real world you'd use one of the two methods above.
:::

## Dynamic Routes

So far every route has a fixed path. But often the path itself contains data — a username, a product id, a page number. Flask lets you capture part of the URL with angle brackets `<...>` and receive it as a function argument:

```python
@app.route("/hi/<name>")
def hi(name):
    return f"Hi, {name}!"
```

Now the `<name>` part of the path is a variable. A request to `/hi/Ada` runs `hi("Ada")` and returns `Hi, Ada!`. A request to `/hi/Sam` returns `Hi, Sam!`. One route, endless possibilities — the value between the slashes is captured and passed straight into your function.

The name inside the brackets must match the function's parameter name. `<name>` pairs with `def hi(name):`. You can even tell Flask what *type* to expect:

```python
@app.route("/square/<int:n>")
def square(n):
    return f"{n} squared is {n * n}"
```

The `<int:n>` **converter** means "only match if this part is a whole number, and give it to me as an `int`." Now `n` is a real integer, so `n * n` does math instead of gluing strings together. Visiting `/square/5` returns `5 squared is 25`.

:::warning
Without the `int:` converter, captured URL values arrive as **strings**. A request to `/square/5` would give you the text `"5"`, and `"5" * "5"` is an error. When a route parameter is a number you'll do math with, use `<int:n>` so Flask hands you an actual integer.
:::

:::predict
Q: With `@app.route("/hi/<name>")` and `def hi(name): return f"Hi, {name}!"`, what does a request to `/hi/Bo` return?
- Hi, name!
- Hi, Bo!*
- Hi, <name>!
E: The `<name>` in the path captures whatever comes after `/hi/` and passes it in as the `name` argument. So `name` is `"Bo"`, and the f-string produces `Hi, Bo!`.
:::

## Returning a Custom Status Code

By default, when your function returns a string, Flask sends it with a **200 OK** status. But sometimes you want a different code — remember from the last lesson that a newly created resource should return **201 Created**.

To control the status, return a **tuple**: the body first, then the status code as a number.

```python
@app.route("/created")
def created():
    return "New resource made!", 201
```

Now a request to `/created` gets the body `"New resource made!"` *and* a `201` status code. The pattern is simply `return body, status`. You can use it for any code:

```python
@app.route("/missing")
def missing():
    return "Nothing here.", 404
```

This returns your own message with a `404 Not Found` status. Being able to set the status code is what makes your responses honest — a client reading `201` knows something was created, and a `404` clearly signals "not found," exactly as we covered in the HTTP lesson.

:::example
Putting it all together, here's the full progression this lab builds — a home page, a health check, a dynamic greeting, and a route with a custom status:

```python
from flask import Flask

app = Flask(__name__)

@app.route("/")
def home():
    return "Welcome to the GamifyDev server!"

@app.route("/ping")
def ping():
    return "pong"

@app.route("/hi/<name>")
def hi(name):
    return f"Hi, {name}!"

@app.route("/created")
def created():
    return "New resource made!", 201
```

Four routes, four ideas — and a real web server.
:::

:::quiz
Q: How do you make a Flask route respond with a 201 status code instead of the default 200?
- return "made it", status=201
- return "made it", 201*
- return 201, "made it"
E: Return a tuple of `(body, status)` — the body first, then the status code, as in `return "made it", 201`. Flask reads the second item as the HTTP status.
:::

## Recap

- **Flask** is a small Python framework that maps URL **routes** to **functions** and turns their return values into HTTP responses.
- Install it with `pip install flask` (already done for you in these labs).
- Every app starts the same: `app = Flask(__name__)`, then `@app.route("/")` above a function that **returns** the response.
- Returning a plain **string** sends it as a `200 OK` body; add routes freely, each a decorator plus a function.
- Run a dev server with `flask run` or `app.run()` — but in these labs **the checker runs it for you**, so you only write routes.
- **Dynamic routes** capture parts of the URL with `<name>`, passed into your function; use `<int:n>` when you need a number.
- Return `body, status` to send a **custom status code**, like `return "made", 201`.

**Next up:** Handling Requests — reading query strings, form data, and JSON so your routes can respond to the details a client sends.
