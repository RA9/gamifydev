# Building a REST API

You have built pages that render HTML. Now you'll build something different: an API — a service that speaks **data**, not pages. Instead of returning a web page for a human to read, a REST API returns clean chunks of JSON for *other programs* to consume. In this lesson you'll design and build one with Flask, using a to-do list as our resource.

## What an API Really Is

API stands for Application Programming Interface. Strip away the jargon and it means one thing: a doorway that lets one program talk to another. When your phone's weather app shows the forecast, it didn't invent that data — it asked a weather API for it and got a tidy reply back.

A **web** API does this over HTTP, the same protocol your browser already speaks. Your Flask app listens at certain URLs and answers with data. The caller might be a JavaScript front-end, a mobile app, or another server — none care what your HTML looks like; they want the raw facts.

:::analogy
A regular web page is a restaurant that serves you a plated, garnished meal ready to eat. An API is the kitchen's supply window — it hands out raw ingredients in labeled containers so *anyone* can cook their own dish. Same food, different packaging.
:::

## What Makes an API "REST"

REST is a *style* for designing web APIs, and it rests on two simple ideas.

**First: every thing is a resource, and every resource has a URL.** A "thing" is a noun in your app — a task, a user, a player. The URL is its address:

```
/tasks          → the whole collection of tasks
/tasks/7        → one specific task, the one with id 7
```

Notice the URLs are **nouns**, not verbs. You never write `/getTasks` or `/deleteTask`. The URL names *what* you're talking about; the *action* comes from somewhere else.

**Second: the HTTP method is the verb.** HTTP already ships with a set of action words, and REST maps them onto the four things you do with data:

```
GET     /tasks       → read the whole list
POST    /tasks       → create a new task
GET     /tasks/7     → read task 7
DELETE  /tasks/7     → delete task 7
```

Same nouns, different verbs, different outcomes. This is the heart of REST: a small, predictable grammar. Learn the pattern once and you can guess how the rest of the API behaves.

:::key
A REST endpoint = an HTTP **method** + a **URL path**. The path names the resource (a noun); the method names the action (GET/POST/PUT/DELETE). Keep paths as nouns and let the verb do the acting.
:::

## Our Resource: a Task

For this lesson, our resource is a task in a to-do list. Each task is a small object with three fields:

```python
{
    "id": 1,
    "title": "Beat level one",
    "done": False
}
```

`id` uniquely identifies the task, `title` is the human-readable text, and `done` is a boolean flag. To keep things simple we'll store our tasks in an ordinary Python list in memory — the same shape the lab uses. (In the next lesson we'll swap that list for a real database so the data survives a restart.)

```python
from flask import Flask, jsonify, request, abort

app = Flask(__name__)

tasks = [
    {"id": 1, "title": "Beat level one", "done": True},
    {"id": 2, "title": "Collect 100 coins", "done": False},
]
```

## Reading the Collection: GET /tasks

The simplest endpoint returns the whole list. In Flask, `jsonify` turns a Python list or dict into a proper JSON response with the right `Content-Type` header:

```python
@app.route("/tasks", methods=["GET"])
def get_tasks():
    return jsonify(tasks)
```

Visit `/tasks` and you get back a JSON array:

```json
[
  {"id": 1, "title": "Beat level one", "done": true},
  {"id": 2, "title": "Collect 100 coins", "done": false}
]
```

Note how Python's `True` became JSON's lowercase `true` — `jsonify` handles that translation for you.

:::tip
Always return data through `jsonify`, not with `str(tasks)`. `jsonify` sets the `Content-Type: application/json` header, so the caller's code knows it's receiving JSON and can parse it automatically.
:::

## Status Codes: the Response's Report Card

Every HTTP response carries a three-digit **status code** that tells the caller how things went — before they even look at the body. You've met `404` in the wild ("page not found"). There's a whole vocabulary, and a good API uses it precisely:

- **200 OK** — the default success. The request worked, here's your data.
- **201 Created** — a new resource was successfully created (the response to a good POST).
- **204 No Content** — success, but there's nothing to send back (a good DELETE).
- **404 Not Found** — you asked for a resource that doesn't exist.

Codes in the 200s mean success, 400s mean *you* (the caller) erred, 500s mean *the server* broke. Returning the right code is part of being a well-behaved API.

:::example
A caller sends `DELETE /tasks/99` but no task 99 exists. A lazy API returns `200 OK` with an empty body, and the caller wrongly thinks it deleted something. A correct API returns `404 Not Found` — an honest signal that there was nothing to delete.
:::

## Creating a Task: POST /tasks

To create a task, the caller sends a POST request with a JSON body describing the new task. On the server we read that body with `request.get_json()`, build the task, append it, and — crucially — return status **201**.

```python
@app.route("/tasks", methods=["POST"])
def create_task():
    data = request.get_json()
    new_id = max((t["id"] for t in tasks), default=0) + 1
    task = {
        "id": new_id,
        "title": data["title"],
        "done": data.get("done", False),
    }
    tasks.append(task)
    return jsonify(task), 201
```

Two details worth pausing on. We compute `new_id` from the current highest id, so ids never collide even after deletions. And we return a **tuple**: `(jsonify(task), 201)`. Flask reads the second element as the status code, and we hand the freshly created task back so the caller learns its assigned `id`.

We also use `data.get("done", False)` instead of `data["done"]`. If the caller omits `done`, `.get` supplies a default instead of crashing with a `KeyError`.

:::predict
Q: A POST creates a resource. Which status code should the response carry?
```python
return jsonify(task), ___
```
- 200
- 201*
- 204
- 404
E: 201 Created is the specific "success, and I made a new thing" code. It signals a resource was born, which 200 (generic OK) doesn't convey.
:::

## Fetching One Task: GET /tasks/<id>

To read a single task, we put a variable in the URL path. Flask captures it and passes it to our function. The `<int:task_id>` converter tells Flask this part of the URL is an integer:

```python
@app.route("/tasks/<int:task_id>", methods=["GET"])
def get_task(task_id):
    for task in tasks:
        if task["id"] == task_id:
            return jsonify(task)
    abort(404)
```

We loop through looking for a matching id. Find it? Return it. Fall off the end of the loop without a match? Then the task doesn't exist, and we call `abort(404)`.

`abort(404)` is Flask's shortcut for "stop here and send back a Not Found error." It immediately ends the function and produces a proper 404 response — you don't build one by hand. There are siblings like `abort(400)` and `abort(403)`, but 404 is the one you'll reach for most.

:::key
`abort(404)` short-circuits the request: it raises an exception Flask catches and turns into a 404 response. Any code after `abort()` never runs, so you can treat it like a `return` that also sets the error status.
:::

## Deleting a Task: DELETE /tasks/<id>

Deletion combines everything so far: find the task, remove it, and answer with the right code. If it isn't there, `abort(404)`.

```python
@app.route("/tasks/<int:task_id>", methods=["DELETE"])
def delete_task(task_id):
    for task in tasks:
        if task["id"] == task_id:
            tasks.remove(task)
            return "", 204
    abort(404)
```

On success we return `"", 204` — an empty body with the **204 No Content** status. There's nothing useful to send back after a delete, and 204 says exactly that: "done, and there's no body to read." A 200 would work too, but 204 is the precise, self-documenting choice.

:::tip
When you're unsure which success code fits: return **201** if you created something, **204** if you did something but have nothing to return, and **200** for everything else. Those three cover the vast majority of REST responses.
:::

## Putting It Together

Here's the complete little API. Save it, run `python app.py`, and you have a working REST service with full create-read-delete over the `tasks` collection:

```python
from flask import Flask, jsonify, request, abort

app = Flask(__name__)

tasks = [
    {"id": 1, "title": "Beat level one", "done": True},
    {"id": 2, "title": "Collect 100 coins", "done": False},
]

@app.route("/tasks", methods=["GET"])
def get_tasks():
    return jsonify(tasks)

@app.route("/tasks", methods=["POST"])
def create_task():
    data = request.get_json()
    new_id = max((t["id"] for t in tasks), default=0) + 1
    task = {"id": new_id, "title": data["title"], "done": data.get("done", False)}
    tasks.append(task)
    return jsonify(task), 201

@app.route("/tasks/<int:task_id>", methods=["GET"])
def get_task(task_id):
    for task in tasks:
        if task["id"] == task_id:
            return jsonify(task)
    abort(404)

@app.route("/tasks/<int:task_id>", methods=["DELETE"])
def delete_task(task_id):
    for task in tasks:
        if task["id"] == task_id:
            tasks.remove(task)
            return "", 204
    abort(404)

if __name__ == "__main__":
    app.run(debug=True)
```

:::quiz
Q: In a REST API, why is `DELETE /tasks/7` preferred over a URL like `/deleteTask/7`?
- Because `/deleteTask` is too long to type
- Because the HTTP method should be the verb and the path should stay a noun*
- Because Flask cannot route URLs that contain the word "delete"
- Because DELETE is faster than GET
E: REST keeps paths as nouns (the resource) and lets the HTTP method carry the action. `DELETE /tasks/7` reads as "delete the tasks/7 resource" using the grammar every REST client already understands.
:::

## Recap

- An API is a doorway for programs to exchange **data**; a web API does it over HTTP and answers with JSON, not HTML.
- REST rests on two ideas: **resources live at URLs** (nouns like `/tasks/7`), and the **HTTP method is the verb** (GET reads, POST creates, DELETE removes).
- Return data with `jsonify` so the `Content-Type` header is correct and Python types convert cleanly to JSON.
- Use status codes precisely: **200** OK, **201** Created (after POST), **204** No Content (after DELETE), **404** Not Found.
- Return a `(body, status)` tuple to set the code, e.g. `return jsonify(task), 201`.
- Capture URL variables with `<int:task_id>`, and use `abort(404)` to short-circuit a request when a resource doesn't exist.

**Next up:** Databases with SQL — trading our in-memory list for storage that survives a restart, so those tasks don't vanish the moment the server stops.
