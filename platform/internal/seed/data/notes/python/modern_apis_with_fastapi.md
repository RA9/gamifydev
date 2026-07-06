# Modern APIs with FastAPI

An API is how one program talks to another over the web. FastAPI is a Python framework for building those APIs quickly — and it uses ordinary Python type hints to do validation, documentation, and editor autocomplete for free. In this lesson you'll build a real, running API.

## What an API Is and Where FastAPI Fits

When your phone app fetches a weather forecast, it isn't loading a web page — it's asking a server for *data*, usually as JSON. That request-and-response contract is an **API** (Application Programming Interface). The server exposes some URLs (called **endpoints**), and clients send requests to them to read or change data.

FastAPI is one of the most popular Python tools for writing that server. You describe your endpoints as plain Python functions, and FastAPI turns them into a working web API — parsing incoming requests, validating them, calling your function, and serializing whatever you return back into JSON.

:::analogy
Think of your API as the counter at a coffee shop. The **endpoint** is the counter itself, the **request** is the customer's order, and the **response** is the finished drink. FastAPI is the well-trained barista who checks the order makes sense before starting, then hands back exactly what was asked for.
:::

## How FastAPI Differs from Flask

If you've seen Flask, FastAPI will feel familiar — both map URLs to functions. The big difference is what your **type hints** do.

In Flask, a path parameter arrives as a string and you convert and check it yourself:

```python
# Flask style — you do the validation by hand
@app.route("/square/<n>")
def square(n):
    n = int(n)          # will crash on "abc"
    return {"result": n * n}
```

In FastAPI, you *annotate* the parameter with a type, and the framework enforces it before your code ever runs:

```python
# FastAPI — the type hint IS the validation
@app.get("/square/{n}")
def square(n: int):
    return {"result": n * n}
```

Send `/square/abc` and FastAPI never calls your function — it returns a clean `422` error explaining that `n` must be an integer. Those same hints also generate interactive documentation and power your editor's autocomplete.

:::key
The core idea of FastAPI: **type hints are not just documentation, they are behavior.** One annotation gives you validation, error messages, docs, and editor support all at once.
:::

## A Minimal App

Every FastAPI project starts with an app object and at least one endpoint. Here is a complete, runnable file:

```python
# main.py
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
    return {"message": "Hello from FastAPI"}
```

Three things are happening:

- `app = FastAPI()` creates the application — the central object that collects all your endpoints.
- `@app.get("/")` is a **decorator** that says "run the function below when someone sends a GET request to `/`."
- The function returns a plain Python `dict`. FastAPI automatically converts it to JSON, so the client receives `{"message": "Hello from FastAPI"}`.

You never call `json.dumps` yourself. Return a dict, a list, or a Pydantic model, and FastAPI serializes it for you.

:::tip
`GET` is the HTTP method for *reading* data; `POST` is for *sending* data to create or process something. FastAPI gives you a decorator for each: `@app.get`, `@app.post`, `@app.put`, `@app.delete`.
:::

## Path Parameters (Typed)

A **path parameter** is a piece of the URL that changes — like the `n` in `/square/5`. You mark it with braces in the route and declare it as a function argument with a type hint.

```python
@app.get("/square/{n}")
def square(n: int):
    return {"input": n, "result": n * n}
```

Now `/square/5` returns `{"input": 5, "result": 25}`. Because `n` is typed as `int`, FastAPI converts the text `"5"` from the URL into a real integer before handing it to you. Requesting `/square/hello` returns a validation error automatically — your function is never even called with bad data.

The type you choose matters. `n: int` accepts whole numbers; `n: float` would accept decimals. The annotation is doing real work every request.

## Query Parameters

**Query parameters** are the `?key=value` pairs at the end of a URL, like `/greet?name=Ada&times=3`. In FastAPI, any function argument that *isn't* part of the path becomes a query parameter automatically.

```python
@app.get("/greet")
def greet(name: str, times: int = 1):
    return {"greeting": f"Hello {name}! " * times}
```

- `name: str` has no default, so it's **required** — leaving it off returns a `422`.
- `times: int = 1` has a default, so it's **optional**; if the caller omits it, `times` is `1`.

Calling `/greet?name=Ada&times=2` returns `{"greeting": "Hello Ada! Hello Ada! "}`. Notice `times` arrives as a proper `int` even though URLs are all text — the type hint converts it.

:::example
The rule of thumb: put values that *identify a thing* in the path (`/users/42`), and values that *filter, sort, or tweak* the result in the query string (`/users?active=true&sort=name`).
:::

## Request Bodies with Pydantic

For anything more than a couple of simple values — especially when a client is *sending* structured data — you use a **request body**. You describe the expected shape with a Pydantic `BaseModel`, and FastAPI validates the incoming JSON against it.

```python
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class Msg(BaseModel):
    text: str

@app.post("/echo")
def echo(msg: Msg):
    return {"echo": msg.text}
```

When a client sends `POST /echo` with a JSON body of `{"text": "hi"}`, FastAPI:

1. Reads the JSON body.
2. Validates it against `Msg` — the `text` field must be present and must be a string.
3. Builds a `Msg` object and passes it in as `msg`.

Inside the function you work with a real, typed object: `msg.text` is a string you can trust. Send `{"text": 123}` or leave `text` out entirely, and Pydantic rejects the request with a precise error before your code runs.

This lesson uses **Pydantic v2**, the current version. A model is just a class that inherits from `BaseModel` with annotated fields. You can add more fields, defaults, and constraints as your API grows:

```python
from pydantic import BaseModel

class Msg(BaseModel):
    text: str
    shout: bool = False   # optional, defaults to False
```

:::warning
Return values also get serialized, but they are *not* validated the way request bodies are unless you set a `response_model`. Never assume the client sent good data without a Pydantic model on the request — validation is exactly what saves you from crashes and bad input.
:::

## Automatic Interactive Docs

Here's the payoff for all those type hints. Because FastAPI knows the exact shape of every request and response, it generates **live, interactive documentation** with no extra work from you.

Start your server and open `http://127.0.0.1:8000/docs` in a browser. You'll see every endpoint listed, the parameters each expects, the models involved, and a **Try it out** button that sends real requests from the page itself.

```text
GET   /             read_root
GET   /square/{n}   square
POST  /echo         echo
```

There's a second docs view at `/redoc` if you prefer a reference-style layout. This is generated from an OpenAPI schema FastAPI builds out of your code — so your docs can never drift out of sync with your actual API.

:::tip
When you're stuck, open `/docs` first. It shows you exactly what your API expects and lets you poke at each endpoint without writing any client code.
:::

## Running the Server (ASGI + async)

FastAPI is built on **ASGI**, the modern asynchronous standard for Python web apps, which means it can handle many requests concurrently and lets you write `async def` endpoints when you need them. You run it with an ASGI server called **uvicorn**:

```bash
pip install fastapi uvicorn
uvicorn main:app --reload
```

`main:app` means "the `app` object inside `main.py`," and `--reload` restarts the server whenever you save a file — perfect for development. Your API is then live at `http://127.0.0.1:8000`.

Endpoints can be plain `def` (as above) or `async def` when they do I/O like calling a database or another API — FastAPI handles both.

## Practice

:::predict
Q: A client requests `/square/hello`, but the endpoint is defined as `def square(n: int)`. What happens?
- The function runs with `n` set to `"hello"`
- FastAPI returns a 422 validation error and never calls the function *
- The server crashes with an unhandled exception
E: The `int` type hint makes FastAPI validate the path parameter *before* your function runs. Non-integer input is rejected with a clean 422 error automatically.
:::

:::quiz
Q: What is the main job of a Pydantic `BaseModel` in a FastAPI endpoint?
- It styles the interactive documentation page
- It describes and validates the shape of request data *
- It starts the uvicorn server
- It replaces the need for the `app` object
E: A `BaseModel` declares the expected fields and their types, so FastAPI can validate incoming JSON and hand your function a typed object.
:::

## Recap

- An API exposes **endpoints** that clients call to read or send data, usually as JSON.
- Create an app with `app = FastAPI()` and map URLs to functions with decorators like `@app.get("/")`.
- Return a `dict` (or list, or Pydantic model) and FastAPI serializes it to JSON automatically.
- **Type hints drive behavior**: `n: int` on a path param `/square/{n}` gives you validation, conversion, and docs for free — the key difference from Flask.
- Function arguments not in the path become **query parameters**; a default makes them optional.
- Describe request bodies with a Pydantic v2 `BaseModel` (`class Msg(BaseModel): text: str`) and receive a typed, validated object.
- Interactive docs live at `/docs`, generated straight from your code.
- FastAPI is ASGI/async-capable; run it with `uvicorn main:app --reload`.

**Next up:** Deploying a Python Web App — taking the API you just built and putting it on the real internet.
