# Project: Build a Notes App

This is the big one: a **full-stack app**. A frontend where you type notes, talking to a backend that stores them. When you finish, you'll have built every layer of a real application.

:::project
A Notes app: a webpage that lists notes and lets you add them, backed by a small Python API that stores the data. Frontend, API, and backend — the whole stack, working together.
:::

## The plan

We'll build two pieces that talk to each other:

1. **Backend** — a Flask API with `GET /notes` (list) and `POST /notes` (add).
2. **Frontend** — an HTML page that fetches the notes and sends new ones.

You've built each half before; now you connect them.

## Part 1 — The backend

Create `app.py`. This should look familiar from the To-Do API:

```python
from flask import Flask, jsonify, request
from flask_cors import CORS

app = Flask(__name__)
CORS(app)   # let the browser frontend call this API

notes = []
next_id = 1

@app.get("/notes")
def list_notes():
    return jsonify(notes)

@app.post("/notes")
def add_note():
    global next_id
    data = request.get_json()
    note = { "id": next_id, "text": data["text"] }
    notes.append(note)
    next_id += 1
    return jsonify(note), 201

if __name__ == "__main__":
    app.run(debug=True)
```

Run `python app.py`. Your API lives at `http://localhost:5000`.

:::warning
A browser will block a frontend from calling an API on a different address unless the server allows it. That's what `CORS(app)` does (install it with `pip install flask-cors`). It's a security feature — not a bug to fight, but a door to open deliberately.
:::

## Part 2 — The frontend

Create `index.html`:

```html
<input id="noteInput" placeholder="Write a note…" />
<button id="addBtn">Add</button>
<ul id="notes"></ul>

<script>
  const API = "http://localhost:5000";

  async function loadNotes() {
    const res = await fetch(API + "/notes");
    const notes = await res.json();
    document.getElementById("notes").innerHTML =
      notes.map(n => "<li>" + n.text + "</li>").join("");
  }

  document.getElementById("addBtn").addEventListener("click", async () => {
    const input = document.getElementById("noteInput");
    await fetch(API + "/notes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text: input.value })
    });
    input.value = "";
    loadNotes();          // refresh the list
  });

  loadNotes();            // load notes when the page opens
</script>
```

## Part 3 — Watch it work

With the backend running, open `index.html` in your browser. Type a note, hit **Add**, and it appears — because:

1. The frontend `POST`s your note to the API.
2. The backend stores it and replies.
3. The frontend re-fetches and re-renders the list.

That's a complete full-stack round trip. **You built an app.** 🎉

:::quiz
Q: When you click "Add", why does the frontend call `loadNotes()` again afterwards?
- To restart the server
- To re-fetch the updated list and show the new note *
- To delete the note
E: After POSTing, the frontend re-fetches from the API so the page reflects the new data — fetch, render, repeat.
:::

## You built it! 🚀

You connected a frontend to a backend over an API, with data flowing both ways. This exact shape — UI ⇄ API ⇄ server ⇄ data — is what real products are made of.

**Stretch goals:**

- Add a delete button that calls `DELETE /notes/<id>`.
- Save notes to a real database so they survive a restart.
- Deploy the backend and host the frontend so others can use it.

:::key
You built a **full-stack app**: a frontend using `fetch` to talk to a backend API that stores data. The round trip — UI sends a request, server stores and responds, UI re-renders — is the heartbeat of every application.
:::

## What's next

You've shipped a full-stack app. The final lesson, **Becoming a Full-Stack Developer**, maps where to go from here.
