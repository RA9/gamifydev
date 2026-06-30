# Deploying Your Website

You've built something great — but right now it only exists on your laptop. **Deploying** is the step that puts your files on a public server so anyone in the world can visit your site at a real URL.

By the end of this lesson you'll have deployed a real site with both GitHub Pages and Netlify, and you'll have a checklist to run before every launch.

## What "deploying" actually means

When you open a `.html` file on your computer, the address bar shows something like `file:///Users/you/site/index.html`. Only you can see that.

A **web server** is a computer that's always on and connected to the internet, ready to send your files to anyone who asks. **Deploying** means copying your files onto such a server and getting a public address like `https://yourname.github.io`.

:::analogy
Building your site is like cooking a meal in your own kitchen. Deploying is opening a restaurant — now the food is somewhere the public can actually come and order it.
:::

## Static hosting — free and perfect for now

Your site is **static**: plain HTML, CSS, and JavaScript files with no server-side code. Static sites are cheap (often free), fast, and easy to host. The big three free options:

- **GitHub Pages** — hosts directly from a GitHub repository. Great if your code already lives on GitHub.
- **Netlify** — drag-and-drop a folder, or connect a repo for automatic deploys.
- **Vercel** — similar to Netlify, especially polished for modern frameworks.

All three are free for personal static sites and give you HTTPS automatically. We'll walk through the first two.

## Deploy with GitHub Pages, step by step

This assumes your project is already a Git repo pushed to GitHub (from the Git lesson).

**1. Make sure your main page is named `index.html`** and lives at the root of the repo. The server looks for `index.html` by default.

**2. Push your code to GitHub** if you haven't:

```bash
git add .
git commit -m "Prepare site for deployment"
git push
```

**3. Enable Pages.** On your repo's GitHub page, go to **Settings -> Pages**. Under **Source**, choose **Deploy from a branch**, pick the `main` branch and the `/ (root)` folder, then **Save**.

**4. Wait a minute, then visit your URL.** GitHub gives you an address like:

```bash
https://yourusername.github.io/your-repo-name/
```

That's it — your site is live. Every time you `git push` new commits, GitHub Pages updates automatically.

:::warning
On GitHub Pages your site often lives in a subfolder (`/your-repo-name/`). That means **absolute paths** like `/css/app.css` break, because they point to the domain root, not your subfolder. Use **relative paths** like `css/app.css` or `./css/app.css` instead. This is the single most common "it worked locally but not deployed" bug.
:::

## Deploy with Netlify

Netlify gives you two easy paths.

**Option A — drag and drop (fastest).**

1. Go to [app.netlify.com](https://app.netlify.com) and sign up (free).
2. Find the deploy area that says "drag and drop your site folder."
3. Drag your whole project folder onto it.

Seconds later you get a live URL like `https://shiny-otter-123abc.netlify.app`. No Git required. The catch: to update the site you have to drag the folder again.

**Option B — connect a repo (auto-deploys).**

1. In Netlify, click **Add new site -> Import an existing project**.
2. Connect your GitHub account and pick your repo.
3. Leave the build settings empty for a plain static site, and set the **publish directory** to the folder containing `index.html` (usually the root).
4. Click **Deploy**.

Now every `git push` triggers a fresh deploy automatically. This is the professional setup: your Git history *is* your deploy history.

```bash
# After the repo is connected, deploying is just:
git add .
git commit -m "Update hero copy"
git push
# Netlify rebuilds and goes live on its own
```

## A quick word on build steps

A **build step** is a command that transforms your source files into the final files the browser gets — bundling JavaScript, compiling CSS, optimizing images. Right now your hand-written HTML/CSS/JS needs no build step, so you can ignore this.

Later, when you use tools like Vite or a framework, you'll add a build command. Hosts let you configure it:

```bash
# Example build settings you'd enter later (not needed yet)
Build command:   npm run build
Publish directory: dist
```

For now: plain files, no build, just deploy.

## Custom domains, briefly

The free `*.github.io` and `*.netlify.app` URLs are perfectly real, but you can attach your own domain like `janedoe.com`.

1. Buy a domain from a registrar (Namecheap, Cloudflare, Google Domains, etc.) — usually around $10–15/year.
2. In your host's dashboard, add the custom domain.
3. Update the domain's **DNS records** (the host gives you exact values to paste at your registrar).
4. Wait for it to propagate (minutes to a few hours). HTTPS is set up for you automatically.

You don't need this to launch — but it's a nice upgrade for a portfolio.

## Pre-deploy checklist

Run through this every time before you ship:

- **Links work.** Click every internal and external link. No `404`s.
- **Paths are relative.** Image `src` and CSS/JS paths use `images/x.png`, not `/images/x.png`, so they survive in a subfolder.
- **Mobile looks right.** Check in DevTools' device toolbar at a few widths.
- **Console is clean.** Open DevTools -> Console. No red errors.
- **Images load and aren't huge.** A 5 MB hero image makes the page crawl — resize it.
- **Title and favicon set.** Your `<title>` and a small favicon make the browser tab look finished.
- **Spelling and content.** Read it once more. Typos in a portfolio are costly.

:::tip
Test the production build by visiting the *deployed* URL on your actual phone, not just the simulator. Real devices catch issues emulators miss.
:::

:::quiz
Q: Your site works locally but on GitHub Pages all the images are broken. What's the most likely cause?
- GitHub Pages doesn't support images
- You used absolute paths like /images/photo.png, which break inside a project subfolder *
- The images are too colorful
E: GitHub Pages serves project sites from a subfolder, so an absolute path starting with `/` points to the wrong place. Relative paths like `images/photo.png` fix it.
:::

:::quiz
Q: What is the advantage of connecting a Git repo to Netlify instead of using drag-and-drop?
- It makes the site load faster for visitors
- Every git push automatically redeploys the site, so updates are hands-off *
- It is the only way to get an HTTPS URL
E: A connected repo gives you continuous deploys — push your code and Netlify rebuilds and publishes on its own. Drag-and-drop requires manually re-uploading each time.
:::

:::fill
For a plain hand-written HTML/CSS/JS site with no bundler, the build step you need is: ______.
- none *
- npm run build
- vite deploy
E: Static, hand-written sites require no build step — the files you wrote are the files the browser gets. Build steps come into play later with tools like Vite.
:::

## Recap

- **Deploying** means putting your files on an always-on public server so the world can reach them at a real URL.
- Your site is **static**, so free hosts like **GitHub Pages**, **Netlify**, and **Vercel** are perfect.
- **GitHub Pages:** push the repo, enable Settings -> Pages on `main`, get a `*.github.io` URL.
- **Netlify:** drag-and-drop for instant deploys, or connect a repo for automatic deploys on every push.
- A **build step** transforms source into final files — you don't need one yet.
- **Custom domains** let you use your own name; optional and inexpensive.
- Run the **pre-deploy checklist** (links, relative paths, mobile, clean console) before every launch.

**Next up:** Building Your Portfolio Website — pulling everything together into a site that shows off your work, then shipping it.
