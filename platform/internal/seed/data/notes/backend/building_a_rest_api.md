# Building a REST API

You can write logic (Python) and store data (SQL). Now you need a way to let the **frontend** ask for that data. The bridge between them is an **API** — and the most common style on the web is a **REST API**.

By the end of this lesson you'll understand what an API is, the HTTP methods, and what an endpoint looks like.

## What is an API?

An **API** (Application Programming Interface) is a set of clear "doors" into your backend. The frontend doesn't touch your database directly — it makes requests to your API, and the API decides what to do.

:::analogy
An API is like a restaurant menu. You don't walk into the kitchen and cook — you order from the menu (the API), and the kitchen (your backend) prepares it. The menu defines exactly what you're allowed to ask for.
:::

## Endpoints and HTTP methods

A REST API is a collection of **endpoints** — URLs that represent your data. Each request also has an HTTP **method** that says what you want to *do*. The methods map neatly onto CRUD:

```text
METHOD   ENDPOINT        WHAT IT DOES
GET      /todos          read all todos
GET      /todos/1        read todo #1
POST     /todos          create a new todo
PUT      /todos/1        update todo #1
DELETE   /todos/1        delete todo #1
```

Notice the pattern: the **URL** names the resource (`/todos`), and the **method** says the action. That consistency is what makes REST predictable.

:::quiz
Q: Which HTTP method would you use to create a new item?
- GET
- POST *
- DELETE
E: POST creates new data. GET reads, PUT/PATCH update, and DELETE removes — together they cover CRUD.
:::

## What an endpoint looks like

Here's a tiny API in Python using **Flask**, a popular beginner-friendly framework:

```python
from flask import Flask, jsonify

app = Flask(__name__)

todos = [
    { "id": 1, "task": "Learn Python", "done": True },
    { "id": 2, "task": "Build an API", "done": False }
]

@app.get("/todos")
def get_todos():
    return jsonify(todos)
```

Visit `/todos` and the server responds with the list. The `jsonify` part matters — it sends the data as **JSON**.

## JSON: the language of APIs

APIs almost always speak **JSON** (JavaScript Object Notation) — a simple, universal text format for data that both Python and JavaScript understand:

```json
{
  "id": 2,
  "task": "Build an API",
  "done": false
}
```

It looks just like a Python dictionary or a JavaScript object — which is exactly why it's perfect for passing data between backend and frontend.

:::quiz
Q: What format do REST APIs usually use to send data?
- HTML
- JSON *
- CSS
E: JSON is the standard data format for APIs — lightweight, readable, and understood by virtually every language.
:::

:::key
A **REST API** exposes your data through **endpoints** (URLs) and **HTTP methods** (GET, POST, PUT, DELETE) that map onto CRUD. Data travels as **JSON**. The frontend talks to the API; the API talks to the database.
:::

## Talk about it

Explain out loud:

> "What is an API, and how do HTTP methods map onto CRUD operations?"

## What's next

Time to build one for real. Next up: **Project: Build a To-Do API**, where you'll create working endpoints step by step.
