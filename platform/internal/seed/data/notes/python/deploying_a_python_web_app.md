# Deploying a Python Web App

You've built an app that runs on your laptop. **Deploying** means putting it somewhere on the internet so anyone, anywhere, can reach it. This capstone lesson walks the path from `python main.py` on your machine to a live URL the whole world can visit.

## What "Deploying" Actually Means

While you develop, your app runs on your own computer and is reachable only by you, usually at an address like `http://127.0.0.1:8000`. That `127.0.0.1` (also called `localhost`) is a private loopback address — it never leaves your machine.

Deploying moves your code onto a **server**: a computer that stays on, has a public address, and runs your app around the clock. Someone else visits your URL, and *their* request travels across the internet to that server, which runs your code and sends a response back.

:::analogy
Developing on `localhost` is like cooking a great meal in your own kitchen — only you can taste it. Deploying is opening a restaurant: the same recipe, but now on a reliable stove, at a public address, ready to serve anyone who walks in.
:::

The rest of this lesson covers the pieces that turn "works on my machine" into "works for everyone."

## The Dev Server vs. a Production Server

The server you use while coding is built for *convenience*, not for crowds. FastAPI's `uvicorn main:app --reload` and Flask's `flask run` auto-restart on file changes and print helpful errors — but they are single-process and not hardened for real traffic. Both frameworks warn you not to use them in production.

Production uses a purpose-built server, and which one depends on your framework's interface:

- **WSGI** is the older, synchronous standard used by Flask and Django. The standard production server is **gunicorn**.
- **ASGI** is the modern, asynchronous standard used by FastAPI. The production server is **uvicorn** (often managed by gunicorn workers).

```bash
# Flask (WSGI) in production — "app" is the Flask object in app.py
gunicorn app:app --workers 4 --bind 0.0.0.0:8000
```

```bash
# FastAPI (ASGI) in production — "app" is the FastAPI object in main.py
uvicorn main:app --host 0.0.0.0 --port 8000
```

Note `--host 0.0.0.0`: unlike `127.0.0.1`, this tells the server to accept connections from *outside* the machine — essential once it's public. The `--workers` flag runs several copies of your app so it can handle many visitors at once.

:::key
Remember the pairing: **Flask → WSGI → gunicorn**, **FastAPI → ASGI → uvicorn**. And drop `--reload` in production — it's a development convenience that wastes resources and can leak details on a live server.
:::

## Dependencies: requirements.txt and Virtual Environments

Your app imports libraries — `fastapi`, `pydantic`, maybe `requests`. The server needs those exact libraries installed, or your code won't start. You capture them in a **`requirements.txt`** file.

```text
fastapi==0.115.0
uvicorn==0.30.6
pydantic==2.9.2
```

Pinning versions with `==` means the server installs the *same* versions you tested with, so it behaves the same everywhere. On any machine, one command installs them all:

```bash
pip install -r requirements.txt
```

To keep each project's libraries isolated from every other project (and from your system Python), you work inside a **virtual environment** — a self-contained folder of packages just for this app.

```bash
python -m venv .venv          # create the environment
source .venv/bin/activate     # activate it (Windows: .venv\Scripts\activate)
pip install -r requirements.txt
```

You generate the file from whatever you've installed with `pip freeze > requirements.txt`.

:::tip
Never commit your `.venv/` folder to git — it's large and machine-specific. Commit `requirements.txt` instead; it's the recipe, and any machine can rebuild the environment from it. Add `.venv/` to your `.gitignore`.
:::

## Environment Variables and Secrets

Real apps need configuration that *isn't* code: database URLs, API keys, passwords. These must never be hard-coded, because your source is shared, committed to git, and read by other people.

The universal solution is **environment variables** — values the operating system hands your program at startup. Your code reads them; the actual secret lives in the deployment platform, not the repo.

```python
import os

# Read a secret from the environment, with a safe fallback for local dev
DATABASE_URL = os.environ.get("DATABASE_URL", "sqlite:///dev.db")
API_KEY = os.environ["API_KEY"]   # required — raises KeyError if missing
```

Locally you might keep these in a `.env` file, but that file is for *your machine only*:

```text
API_KEY=sk-local-dev-key-123
DATABASE_URL=postgresql://localhost/myapp
```

:::warning
Never commit secrets to git. Add `.env` to your `.gitignore` from the very first commit. A leaked API key in a public repo can be scraped within minutes and run up huge bills or expose user data. On your hosting platform, set these values in its "Environment Variables" or "Secrets" settings instead.
:::

## The PORT a Host Provides

On your laptop you pick the port — `8000` is just a habit. In production, the **hosting platform** usually chooses the port for you and tells your app which one to use through a `PORT` environment variable. Your app must listen on *that* port, not a hard-coded one, or the host can't route traffic to it.

```python
import os
import uvicorn

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8000))   # use the host's PORT, else 8000 locally
    uvicorn.run("main:app", host="0.0.0.0", port=port)
```

Or straight from the start command, where many hosts expand `$PORT` for you:

```bash
uvicorn main:app --host 0.0.0.0 --port $PORT
```

The pattern is the same everywhere: read `PORT` from the environment, fall back to a local default, and always bind to `0.0.0.0` so the platform can reach you.

## The Start Command and the Procfile

The host needs one more thing: the exact command to launch your app. Some platforms detect a Python project and guess; most let you state it explicitly.

The classic form is a **`Procfile`** — a one-line file naming the process type (`web`) and its command:

```text
web: uvicorn main:app --host 0.0.0.0 --port $PORT
```

For a Flask app it would be:

```text
web: gunicorn app:app --bind 0.0.0.0:$PORT
```

Newer platforms often ask for the same thing as a **"Start Command"** field in their dashboard instead of a file. Either way, it's the single line that answers "how do I run this?" — and it's the same command you'd type by hand, just parameterized with `$PORT`.

:::example
A minimal deployable FastAPI project is often just four files: `main.py` (your app), `requirements.txt` (its dependencies), `Procfile` (how to start it), and `.gitignore` (what to leave out). Push those to a git repo and most platforms can take it from there.
:::

## A Tour of Hosting Options

You don't need to rent and configure a raw server anymore. Modern **platform-as-a-service** hosts read your git repo, install your `requirements.txt`, run your start command, and hand you a URL.

- **Railway** and **Render** — connect a GitHub repo, set your environment variables in the dashboard, and they build and deploy on every push. Great first choices; generous free or low-cost tiers.
- **Fly.io** — deploys your app close to your users in multiple regions; leans toward containers but has a smooth CLI.
- **Containers (Docker)** — you package your app *and* its exact environment into an image described by a `Dockerfile`. That image runs identically on your laptop, a teammate's machine, or any cloud. It's more setup, but it's the most portable and reproducible option, and it underlies most of the platforms above.

```text
# Dockerfile — the recipe for a self-contained, runnable image
FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

Start with a platform-as-a-service host to get live fast; reach for containers when you need full control or are deploying to a serious cloud like AWS or Google Cloud.

:::tip
Deploy *early* and *often*, even before your app is finished. A tiny "Hello World" API that's live on the internet teaches you the whole pipeline — repo, dependencies, environment variables, port, start command — while there's almost nothing that can go wrong. Then you're just pushing updates.
:::

## Practice

:::quiz
Q: Which production server pairs with a FastAPI (ASGI) app?
- gunicorn, because FastAPI uses WSGI
- uvicorn, because FastAPI uses ASGI *
- the built-in dev server with `--reload` left on
- Docker, which replaces the need for a server
E: FastAPI is an ASGI framework, and uvicorn is its ASGI production server. gunicorn pairs with WSGI frameworks like Flask; the dev server with `--reload` should never be used in production.
:::

:::predict
Q: A host sets a `PORT` environment variable, but your app hard-codes `port=8000`. What's the likely result?
- The app works fine; the host adapts to port 8000
- The host can't route traffic to your app and it appears unreachable *
- The app refuses to start with a syntax error
E: The platform routes incoming traffic to the port it assigned via `PORT`. If your app listens on a different, hard-coded port, nothing connects — you must read `PORT` from the environment.
:::

## Recap

- **Deploying** moves your app from `localhost` (only you) onto a public server (everyone).
- The **dev server** is for coding; production uses a real server: **gunicorn** for Flask (WSGI), **uvicorn** for FastAPI (ASGI). Bind to `0.0.0.0` and drop `--reload`.
- List your libraries in **`requirements.txt`** (pin versions) and develop inside a **virtual environment**; commit the file, not the `.venv/`.
- Keep secrets in **environment variables**, never in code, and never commit a `.env` file.
- Read the **`PORT`** the host provides instead of hard-coding one.
- Give the host a **start command** — as a `Procfile` line or a dashboard field.
- Start on a platform host like **Railway** or **Render**; reach for **containers** when you need portability and control.

**You've reached production.** From a single `dict` returned by a FastAPI endpoint to a live URL on the internet — that's the full journey of a modern Python web app.
