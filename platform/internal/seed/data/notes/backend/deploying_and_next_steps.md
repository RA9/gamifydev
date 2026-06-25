# Deploying and Next Steps

You've built a backend that runs on your computer. The final step to making it *real* is **deploying** it — putting it on a server the whole world can reach. Then we'll map out where to go from here.

## What "deploying" means

When you run `python app.py`, your API lives at `localhost` — your machine only. **Deploying** means uploading your code to a server that's always on and has a public address, so anyone (and any frontend) can call it.

:::analogy
Building locally is like cooking a great meal in your own kitchen. Deploying is opening a restaurant — now anyone can come and order. Same food; suddenly the world can taste it.
:::

## How deployment usually goes

The modern flow is refreshingly simple, and you already know the first step:

1. **Push your code to GitHub** (remember `git push`?).
2. **Connect a hosting platform** to your repo — services like **Render**, **Railway**, **Fly.io**, or **PythonAnywhere** are beginner-friendly and have free tiers.
3. The platform installs your dependencies, runs your app, and gives you a public URL.

After that, every time you push to GitHub, your live site can update automatically.

:::tip
List your dependencies in a `requirements.txt` file (e.g. `flask`). Hosting platforms read it to install exactly what your app needs. It's the backend equivalent of a recipe's ingredient list.
:::

:::quiz
Q: Why can't other people use your API while it only runs on `localhost`?
- localhost is too slow
- localhost only refers to your own computer — it isn't public *
- APIs can't be shared
E: `localhost` points to your own machine. Deploying puts your code on a public server with an address others can reach.
:::

## Where to go next as a backend developer

You've got the foundations. Strong next moves:

- **Use a real database** — move from in-memory lists to **PostgreSQL** or **SQLite**, and learn an ORM (like SQLAlchemy) to talk to it from Python.
- **Add authentication** — let users sign up, log in, and protect their own data.
- **Learn about security** — validating input, hashing passwords, and avoiding common vulnerabilities.
- **Explore other tools** — Node.js/Express is the other hugely popular backend stack worth knowing.

:::key
**Deploying** moves your backend from `localhost` to a public server — usually by pushing to **GitHub** and connecting a hosting platform. From here, grow by adding a real **database**, **authentication**, and stronger **security**.
:::

## Talk about it

Explain out loud:

> "What does it mean to deploy a backend, and why won't a `localhost` API work for real users?"

## What's next

You can now build *and* ship a backend. If you want to connect a frontend to a backend and build complete apps, the **Fullstack** path is your next adventure. Keep building! 🚀
