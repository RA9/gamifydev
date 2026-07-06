# Databases with SQL

In the last lesson our tasks lived in a Python list. That list has a fatal flaw: the moment the server stops, every task vanishes. A database fixes this — it stores your data on disk so it survives restarts, crashes, and deploys. In this lesson you'll learn just enough SQL to run a leaderboard, and you'll drive it all from Python with the built-in `sqlite3` module.

## Why a Database and Not a List

A Python list lives in your program's memory (RAM). RAM is fast but *temporary* — it's wiped clean every time the process ends. Store your players in a list, restart the server, and they're gone.

A database writes to disk, so the data is **persistent**: still there tomorrow. Databases also do things a plain list can't at scale — search millions of rows in milliseconds, sort by any field, count matches, and let many programs read and write safely at once.

The most common kind is a **relational** database, which organizes data into tables. A table is like a spreadsheet: named columns across the top, one row per record.

:::analogy
A Python list is a whiteboard — quick to scribble on, but wiped every time the room is cleaned. A database is a filing cabinet: slower to open a drawer, but your files are still there next week, neatly labeled and easy to search.
:::

We'll use **SQLite**, a tiny database that stores everything in a single file. There's no separate server to install — Python already ships with it. Perfect for learning, and genuinely used in real apps.

## SQL: the Language of Tables

SQL (Structured Query Language) is how you talk to a relational database. You write short English-like statements, and the database carries them out. Let's build a `players` table for a game leaderboard, where each player has a `name` and a `score`.

**CREATE TABLE** defines the shape of a table — its columns and their types:

```sql
CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    score INTEGER NOT NULL
);
```

`INTEGER` and `TEXT` are column types. `PRIMARY KEY` makes `id` a unique identifier that the database auto-fills for you. `NOT NULL` means the column can't be left empty.

**INSERT** adds a row:

```sql
INSERT INTO players (name, score) VALUES ('Ada', 900);
```

**SELECT** reads rows back out. `*` means "all columns":

```sql
SELECT * FROM players;
```

## Filtering, Sorting, and Counting

A `SELECT` becomes powerful once you add clauses to shape the result.

**WHERE** filters to only the rows that match a condition:

```sql
SELECT name, score FROM players WHERE score > 500;
```

**ORDER BY** sorts the results. Add `DESC` for descending (highest first); the default is `ASC`, ascending:

```sql
SELECT name, score FROM players ORDER BY score DESC;
```

**LIMIT** caps how many rows come back. Combine `ORDER BY` with `LIMIT 1` and you've asked for the single top player:

```sql
SELECT name FROM players ORDER BY score DESC LIMIT 1;
```

**COUNT** is an aggregate — it collapses many rows into one number. `COUNT(*)` counts rows:

```sql
SELECT COUNT(*) FROM players WHERE score > 500;
```

:::predict
Q: What does this SQL return?
```sql
SELECT name FROM players ORDER BY score DESC LIMIT 1;
```
- every player's name
- the highest scorer's name*
- the lowest scorer's name
E: ORDER BY score DESC puts the top score first; LIMIT 1 keeps only that row, so you get the single highest scorer's name.
:::

## Running SQL from Python with sqlite3

SQL statements have to run *somewhere*. Python's built-in `sqlite3` module is your bridge: it opens the database file, sends SQL, and hands back the results as Python values. Nothing to install.

You start by **connecting**, which opens (or creates) the database file, then getting a **cursor** — the object you actually run statements through:

```python
import sqlite3

conn = sqlite3.connect("game.db")   # opens/creates the file
cursor = conn.cursor()

cursor.execute("""
    CREATE TABLE IF NOT EXISTS players (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        score INTEGER NOT NULL
    )
""")
conn.commit()
```

`execute` runs one SQL statement. `IF NOT EXISTS` keeps the CREATE from erroring if the table already exists. And `conn.commit()` **saves** your changes to disk — without it, writes can be lost. Reads don't need a commit; writes (INSERT, DELETE, CREATE) do.

To insert many rows at once, `executemany` runs the same statement for each item in a list — cleaner than a loop:

```python
players = [("Ada", 900), ("Grace", 750), ("Linus", 820)]
cursor.executemany(
    "INSERT INTO players (name, score) VALUES (?, ?)",
    players,
)
conn.commit()
```

## Getting Results Back: fetchone and fetchall

After a `SELECT`, you pull the rows out with `fetchone` (the next single row) or `fetchall` (every remaining row, as a list). Each row comes back as a **tuple** of column values, in the order you selected them.

`fetchone` is ideal when you expect one row — like the top player:

```python
cursor.execute("SELECT name, score FROM players ORDER BY score DESC LIMIT 1")
top = cursor.fetchone()
print(f"Top player: {top[0]} with {top[1]} points")
# → Top player: Ada with 900 points
```

`top` is the tuple `("Ada", 900)`, so `top[0]` is the name and `top[1]` is the score. If no rows match, `fetchone` returns `None`, so guard against that before indexing.

`fetchall` gives you the whole leaderboard to loop over:

```python
cursor.execute("SELECT name, score FROM players ORDER BY score DESC")
leaderboard = cursor.fetchall()

for rank, (name, score) in enumerate(leaderboard, start=1):
    print(f"{rank}. {name} — {score}")
# → 1. Ada — 900
# → 2. Linus — 820
# → 3. Grace — 750
```

:::tip
`fetchall()` loads *every* matching row into memory at once. For a leaderboard of a few hundred that's fine, but always narrow your query with `WHERE`, `ORDER BY ... LIMIT`, or `COUNT` rather than fetching a huge table and filtering in Python. Let the database do the heavy lifting.
:::

## Parameterized Queries: the `?` Placeholder

Here is the single most important habit in this lesson. When your SQL includes a value that came from *outside* — a player's name typed into a form, a search term from a URL — you must **never** glue that value into the SQL string yourself. Instead you leave a `?` placeholder and pass the value separately:

```python
name = "Ada"
cursor.execute(
    "SELECT score FROM players WHERE name = ?",
    (name,),          # values go in a tuple, matched to each ?
)
row = cursor.fetchone()
```

The `?` marks a hole; the tuple fills it. `sqlite3` inserts the value safely, treating it strictly as *data* — never as SQL to run. Use one `?` per value, passed in a tuple in the same order. (Note `(name,)` — that trailing comma makes it a one-element tuple, not just parentheses.)

Contrast the two styles:

```python
# ❌ NEVER do this — building SQL by string formatting
cursor.execute(f"SELECT score FROM players WHERE name = '{name}'")

# ✅ ALWAYS do this — a parameter placeholder
cursor.execute("SELECT score FROM players WHERE name = ?", (name,))
```

Why does it matter so much? Because of an attack called **SQL injection**.

:::warning
**SQL injection.** If you build queries with f-strings or `+`, a malicious value can *become code*. Suppose a user types their name as `'; DROP TABLE players; --`. With string formatting, your query turns into `SELECT score FROM players WHERE name = ''; DROP TABLE players; --'` — and your entire table is deleted. The `?` placeholder makes this impossible: the input is always treated as a value, never as SQL. **Always** use `?` for outside data. No exceptions.
:::

## Connecting a Query to a Web Route

Now the two lessons meet. Instead of returning a hard-coded list, a Flask route can run a SQL query and return the results as JSON. Conceptually, the route is just a thin wrapper around a query:

```python
import sqlite3
from flask import Flask, jsonify

app = Flask(__name__)

@app.route("/leaderboard")
def leaderboard():
    conn = sqlite3.connect("game.db")
    cursor = conn.cursor()
    cursor.execute(
        "SELECT name, score FROM players ORDER BY score DESC LIMIT 10"
    )
    rows = cursor.fetchall()
    conn.close()

    # turn each (name, score) tuple into a labeled dict for clean JSON
    players = [{"name": name, "score": score} for name, score in rows]
    return jsonify(players)
```

The request comes in, the route runs a SELECT, converts the row tuples into dictionaries, and `jsonify` sends them back. The data lives on disk in `game.db`, so it's the same leaderboard for every visitor and still there after a restart. Remember to `conn.close()` when you're done to release the file.

:::key
A web route that reads data is usually three steps: **connect** to the database, **execute** a query with any user input passed via `?`, and **shape** the rows into JSON. The database is the source of truth; the route just fetches and formats.
:::

:::quiz
Q: Why should a query use `WHERE name = ?` with the value in a tuple instead of an f-string?
- The `?` version runs faster on large tables
- The `?` version prevents SQL injection by treating the input as data, never as executable SQL*
- f-strings don't work inside the `sqlite3` module
- There is no real difference; it's just a style preference
E: The placeholder guarantees the value is inserted as data, so a malicious input like `'; DROP TABLE players; --` can never break out and run as SQL. String formatting offers no such protection.
:::

## Recap

- Databases store data on **disk**, so it persists across restarts — unlike a Python list, which lives in RAM and is wiped when the program ends.
- SQL is the language of relational tables: **CREATE TABLE** defines columns, **INSERT** adds rows, **SELECT** reads, **WHERE** filters, **ORDER BY** sorts, `LIMIT` caps rows, and **COUNT** aggregates to a number.
- Python's built-in **`sqlite3`** bridges Python and SQL: `connect` opens the file, a cursor's `execute`/`executemany` runs statements, and `commit` saves writes to disk.
- Pull results with **`fetchone`** (one row, or `None`) and **`fetchall`** (a list of rows); each row is a tuple of column values.
- **Always use `?` placeholders** for any value from outside your program. Building SQL with f-strings or `+` opens the door to SQL injection.
- A data-reading web route is just connect → execute → shape rows into JSON; the database is the source of truth.

**Next up:** Talking to Other APIs — instead of serving data, you'll become a *client* and fetch data from someone else's API over the network.
