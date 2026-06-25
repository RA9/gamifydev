# Project: Build a To-Do API

You've learned Python, databases, and what an API is. Now you'll **build a real one** — a backend that a frontend could actually use to manage a to-do list.

:::project
A REST API for to-dos with endpoints to list, add, and delete tasks — written in Python with Flask. You'll see your own server respond with real JSON.
:::

## Setup

You'll need Python installed, then Flask:

```bash
pip install flask
```

Create a file called `app.py`. We'll build it up endpoint by endpoint, running it after each step.

## Step 1 — A server that runs

```python
from flask import Flask, jsonify, request

app = Flask(__name__)

todos = [
    { "id": 1, "task": "Learn Python", "done": True }
]
next_id = 2

@app.get("/todos")
def list_todos():
    return jsonify(todos)

if __name__ == "__main__":
    app.run(debug=True)
```

Run it with `python app.py`, then open `http://localhost:5000/todos`. Your server answers with JSON. You just ran a backend! 🎉

## Step 2 — Add a todo (POST)

A real app needs to *create* data. This endpoint reads JSON from the request and adds it:

```python
@app.post("/todos")
def add_todo():
    global next_id
    data = request.get_json()
    todo = { "id": next_id, "task": data["task"], "done": False }
    todos.append(todo)
    next_id += 1
    return jsonify(todo), 201
```

The `201` is an HTTP status code meaning "Created". Status codes are how an API tells the frontend how things went (200 OK, 404 Not Found, and so on).

## Step 3 — Delete a todo (DELETE)

```python
@app.delete("/todos/<int:todo_id>")
def delete_todo(todo_id):
    global todos
    todos = [t for t in todos if t["id"] != todo_id]
    return jsonify({ "deleted": todo_id })
```

The `<int:todo_id>` part captures the id from the URL — so `DELETE /todos/1` removes todo #1.

## Step 4 — Try it out

Test your endpoints with `curl` (or a tool like Postman):

```bash
# list
curl http://localhost:5000/todos

# add
curl -X POST http://localhost:5000/todos \
  -H "Content-Type: application/json" \
  -d '{"task": "Build an API"}'

# delete
curl -X DELETE http://localhost:5000/todos/1
```

Watch the list change as you add and remove items. That's a working backend. 🚀

:::quiz
Q: In `DELETE /todos/1`, what does the `1` refer to?
- The number of todos to delete
- The id of the specific todo to delete *
- The server port
E: The `1` is captured from the URL as `todo_id` — it identifies which specific todo to delete.
:::

## You built it! 🚀

You created a REST API with **GET**, **POST**, and **DELETE** endpoints returning real JSON — exactly how the apps you use every day store and serve data.

**Stretch goals:**

- Add a `PUT /todos/<id>` endpoint to mark a todo done.
- Return a `404` when someone asks for an id that doesn't exist.
- Swap the in-memory list for a real database (SQLite is a great first step).

:::key
You built a working REST API: endpoints mapped to HTTP methods, JSON in and out, and status codes to report results. Connect a frontend to this and you have a full app.
:::

## What's next

Your API works on your machine — but nobody else can reach it yet. Next, **Deploying and Next Steps** puts your backend online.
