# Talking to Other APIs

So far you've been the one *serving* data — building an API for others to call. Now you'll flip roles and become the **client**: the program that reaches out and asks *someone else's* API for data. Want live weather, exchange rates, or a global game leaderboard your own server doesn't host? You fetch it over the network. This lesson shows you how, using Python's `requests` library.

## Being a Client Instead of a Server

Every API interaction has two sides. The **server** waits at a URL and answers requests — that was you in the REST lesson. The **client** is whoever sends the request and reads the reply. When your Flask app calls out to a game-stats service to grab a worldwide leaderboard, *your app* is now the client, and that service is the server.

The beautiful part: it's the same HTTP grammar you already know, just seen from the other end. You still send a GET or a POST to a URL, you still read a status code, you still parse JSON out of the response. The only difference is who initiates.

:::analogy
Building your own API was like running a shop and serving customers at the counter. Being a client is *you* walking into someone else's shop, asking for something, and carrying it home. Same transactions, opposite side of the counter.
:::

## The requests Library

Python's standard library can make HTTP calls, but the tooling is clunky. The community-standard library **`requests`** wraps it in something genuinely pleasant. You install it once:

```bash
pip install requests
```

Then a GET request — the workhorse — is a single, readable line:

```python
import requests

response = requests.get("https://api.example-games.com/leaderboard")
```

That call reaches across the network, sends an HTTP GET, waits for the reply, and hands you back a `response` object holding everything the server said.

:::tip
You'll also hear about **`httpx`**, a newer library with an almost identical API to `requests`. Its big advantage is that it can make **async** requests — firing off many calls concurrently without blocking — which matters when you're fetching from lots of APIs at once. Learn `requests` first; the concepts transfer directly to `httpx` when you need the extra power.
:::

## Reading the Response: status_code and json()

The `response` object has two attributes you'll reach for constantly.

`response.status_code` is the same three-digit code from the REST lesson, now seen from the client's side. **200** means success; **404** means the resource wasn't found; **500** means the server broke. Always check it before trusting the body:

```python
response = requests.get("https://api.example-games.com/leaderboard")
print(response.status_code)   # → 200
```

`response.json()` parses the response body from JSON text into Python objects — dicts and lists you can index like any other:

```python
if response.status_code == 200:
    data = response.json()
    for player in data["players"]:
        print(player["name"], player["score"])
```

`response.json()` does the reverse of the `jsonify` you used as a server: it takes JSON text off the wire and turns it back into Python data structures. `data["players"]` is now a normal list of dicts.

:::example
Fetching a game's top-10 board might return JSON like `{"players": [{"name": "Ada", "score": 900}, ...]}`. After `data = response.json()`, `data["players"][0]["name"]` is simply the string `"Ada"` — a plain Python dict lookup. The network round-trip is over; you're back in familiar territory.
:::

## Sending Query Parameters

APIs often let you *narrow* what you ask for using query parameters — the `?key=value` bits you see in URLs. You could paste them onto the URL by hand, but `requests` has a cleaner way: pass a `params` dict, and it builds the query string (and escapes it) for you:

```python
response = requests.get(
    "https://api.example-games.com/leaderboard",
    params={"region": "eu", "limit": 10},
)
# requests builds: .../leaderboard?region=eu&limit=10
print(response.url)   # handy for seeing exactly what was requested
```

Here we're asking for only the European region's top 10. Letting `requests` assemble the query string means you never have to worry about spaces or special characters breaking the URL.

## Sending Data with POST

Reading isn't the only thing clients do. To *create* something on another server — say, submit your own game's score to a shared leaderboard — you send a POST with a JSON body. The `json=` argument tells `requests` to serialize your dict to JSON and set the right headers:

```python
response = requests.post(
    "https://api.example-games.com/scores",
    json={"name": "Ada", "score": 950},
)

if response.status_code == 201:
    print("Score submitted!")
    print(response.json())   # the server echoes back what it created
```

Notice the **201** check — the same "Created" status you *returned* as a server now confirms, as a client, that your new score was accepted. Both sides of the conversation speak the identical vocabulary.

:::key
`json=` (POST body) and `params=` (query string) are different tools. `params` shapes *which* resource you're asking about and rides in the URL; `json` carries the *content* you're sending to be stored. Reading usually uses `params`; creating usually uses `json`.
:::

## Handling Errors and Timeouts

The network is unreliable in a way local code never is. A server can be down, slow, or return an error. Robust client code plans for this.

First, **never assume success** — check the status code, or call `raise_for_status()`, which throws an exception on any 4xx or 5xx response so failures don't slip by silently:

```python
response = requests.get("https://api.example-games.com/leaderboard")
response.raise_for_status()   # raises on 404, 500, etc.
data = response.json()
```

Second, **always set a timeout**. Without one, a hung server can freeze your program *forever* while it waits for a reply that never comes. The `timeout` argument caps the wait in seconds:

```python
try:
    response = requests.get(
        "https://api.example-games.com/leaderboard",
        timeout=5,   # give up after 5 seconds
    )
    response.raise_for_status()
    data = response.json()
except requests.exceptions.Timeout:
    print("The leaderboard service is taking too long.")
except requests.exceptions.RequestException as error:
    print(f"Request failed: {error}")
```

`requests.exceptions.RequestException` is the base class for every error `requests` can raise — connection failures, timeouts, bad status codes — so catching it is a solid safety net.

:::warning
A missing `timeout` is a silent time bomb. If the remote server hangs, your request will wait indefinitely, and in a web app that ties up a worker that could be serving other users. Treat `timeout=` as mandatory on every network call, not optional.
:::

## API Keys and Auth Headers

Most real APIs won't answer anonymous callers forever — they want to know *who* is asking, to enforce limits and, sometimes, to bill you. The common gatekeeper is an **API key**: a long secret string the API provider issues to you.

You usually send the key in an HTTP **header**, most often an `Authorization` header. In `requests`, headers are just a dict passed via `headers=`:

```python
import os

api_key = os.environ["GAME_API_KEY"]   # read from the environment, not hard-coded

response = requests.get(
    "https://api.example-games.com/leaderboard",
    headers={"Authorization": f"Bearer {api_key}"},
    timeout=5,
)
```

The `Bearer <token>` format is a widespread convention. The exact header an API wants is always in its documentation — read it. The key point is *conceptual*: the key identifies you, it travels in a header, and it must be kept **secret**.

:::warning
Never paste an API key directly into your code or commit it to git. Anyone who reads the file — or your public repo — can use your key, run up your bill, or hit your rate limits. Store secrets in **environment variables** (as `os.environ` above does) or a config file that git ignores, and keep them out of version control.
:::

## A Note on Where Network Calls Run

One practical caveat for this course. The interactive labs run in your browser's sandbox, which has **no live internet access** — it can't actually reach `api.example-games.com`. So the network examples here are for you to run on **your own machine**, where `pip install requests` and a real connection let them work for real.

In the in-browser labs you'll practice the *shape* of client code — building the request, reading `status_code`, parsing `response.json()` — against stand-in data, so the patterns are second nature by the time you make a genuine call from your own project.

:::predict
Q: You call `response = requests.get(url)`. What turns the reply into a Python dict you can index?
```python
data = response.____()
```
- text
- json*
- status_code
- params
E: `response.json()` parses the JSON body into Python objects (dicts and lists). `response.text` gives the raw string, and `status_code` is just the number — neither gives you an indexable dict.
:::

:::quiz
Q: Why should every `requests` call include a `timeout` argument?
- It makes the request run faster
- It stops your program from hanging forever if the remote server never responds*
- It is required to parse the JSON response
- It automatically retries the request on failure
E: Without a timeout, a hung or unreachable server leaves your request waiting indefinitely, freezing your program (and, in a web app, tying up a worker). A timeout caps the wait and lets you handle the failure.
:::

## Recap

- As a **client** you send requests to someone else's API and read the reply — the same HTTP grammar as building an API, seen from the other side of the counter.
- The **`requests`** library makes calls simple: `requests.get(url)` and `requests.post(url, json=...)`. `httpx` is the modern, async-capable equivalent with a nearly identical API.
- Read the reply with `response.status_code` (did it work?) and `response.json()` (turn the JSON body into Python data).
- Use `params={...}` to add query-string parameters and `json={...}` to send a JSON body on a POST.
- Handle failure deliberately: check the status or call `raise_for_status()`, and **always** pass `timeout=` so a slow server can't hang your program.
- Authenticate with an **API key**, usually sent in an `Authorization` header — and keep it secret via environment variables, never hard-coded or committed.
- Live network calls run on **your own machine**, not the in-browser labs, which practice the shape of client code against stand-in data.

**Next up:** you now have both halves — serving data and consuming it. From here you can wire your Flask app to fetch, store, and re-serve data from anywhere on the web.
