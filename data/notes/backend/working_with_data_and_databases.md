# Working with Data and Databases

Variables vanish the moment your program stops. But a real app needs to *remember* — your account, your posts, your orders — even after the server restarts. That's the job of a **database**.

By the end of this lesson you'll understand what a database is and how to read and write data with **SQL**.

## What is a database?

A **database** is organised, permanent storage for your app's data. The most common kind — a *relational* database — stores data in **tables**, which look a lot like spreadsheets: rows and columns.

A `users` table might look like:

```text
 id |  name  |     email
----+--------+------------------
  1 |  Ada   |  ada@email.com
  2 |  Grace |  grace@email.com
```

Each **row** is one record (a user); each **column** is a field (name, email).

:::analogy
A database is a giant, well-organised filing cabinet. Tables are the drawers, rows are the folders, and columns are the labelled fields on each folder. SQL is how you ask the cabinet for exactly what you need.
:::

## SQL: talking to the database

**SQL** (Structured Query Language) is how you read and write data. The four everyday operations:

```sql
-- READ: get data
SELECT name, email FROM users WHERE id = 1;

-- CREATE: add a new row
INSERT INTO users (name, email) VALUES ('Linus', 'linus@email.com');

-- UPDATE: change existing data
UPDATE users SET email = 'new@email.com' WHERE id = 1;

-- DELETE: remove a row
DELETE FROM users WHERE id = 2;
```

These four — read, create, update, delete — are so common they have a nickname: **CRUD**. Almost every app is, under the hood, CRUD on some tables.

:::quiz
Q: Which SQL statement retrieves data from a table?
- `INSERT`
- `SELECT` *
- `DELETE`
E: `SELECT` reads data. `INSERT` adds rows, `UPDATE` changes them, and `DELETE` removes them.
:::

## Filtering with WHERE

`WHERE` is how you get *specific* rows instead of everything:

```sql
SELECT * FROM todos WHERE done = false;   -- only unfinished todos
```

The `*` means "all columns". Leaving off `WHERE` returns every row — handy, but be careful with `UPDATE` and `DELETE`!

:::warning
A `DELETE` or `UPDATE` without a `WHERE` clause affects **every row** in the table. Always double-check your `WHERE` before running a change.
:::

## Keys connect tables

Each row usually has a unique **primary key** (often `id`). Other tables refer to it with a **foreign key** — that's how a `posts` table knows which `user` wrote each post. These connections (relations) are why it's called a *relational* database.

:::quiz
Q: What does the nickname "CRUD" stand for?
- Copy, Run, Undo, Delete
- Create, Read, Update, Delete *
- Connect, Render, Upload, Deploy
E: CRUD = Create, Read, Update, Delete — the four basic operations on stored data, matching SQL's INSERT, SELECT, UPDATE, and DELETE.
:::

:::key
A **database** stores data permanently in **tables** (rows and columns). You read and write it with **SQL**, and the four core operations are **CRUD** — Create, Read, Update, Delete. Always use `WHERE` to target the right rows.
:::

## Talk about it

Explain out loud:

> "Why do we need a database instead of variables, and what are the four CRUD operations?"

Take the lesson quiz to test your SQL — then you'll be ready to put data behind an API.

## What's next

You can store data and write logic. Next, **Building a REST API** ties them together and exposes your data to the frontend.
