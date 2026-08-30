# Deployment Platforms and CI Basics

You've built projects and tracked them with Git. Now learn how to ship them to the world — and how professional teams automate the entire process.

## GitHub Pages — deploy from a branch

GitHub Pages is the simplest way to deploy a static site directly from a GitHub repository.

### Setup (step by step)

1. Push your project to a GitHub repository
2. Make sure your main page is named `index.html` at the root
3. Go to **Settings → Pages**
4. Under **Source**, select **Deploy from a branch**
5. Choose the `main` branch and `/ (root)` folder
6. Click **Save**

Your site goes live at:

```text
https://yourusername.github.io/your-repo-name/
```

Every future `git push` to `main` triggers a rebuild. Changes go live in about 60 seconds.

```bash
# Deploying is just pushing
git add .
git commit -m "Update hero section copy"
git push
# Site updates automatically
```

:::tip
If your site uses relative paths for CSS, JS, and images (e.g. `style.css` not `/style.css`), GitHub Pages works without any path issues. Absolute paths break because the site lives in a subfolder.
:::

### Using GitHub Actions for Pages

For more control, you can deploy using **GitHub Actions** instead of the branch-based method:

1. In Settings → Pages, change Source to **GitHub Actions**
2. GitHub will suggest a workflow — accept the default for static sites
3. It creates a `.github/workflows/static.yml` file that deploys on every push

This is the professional setup — it gives you build steps, environment variables, and deployment previews.

## Netlify — drag-and-drop or connect repo

Netlify is a dedicated deployment platform with more features than GitHub Pages.

### Option A — drag and drop

1. Go to **app.netlify.com** and sign up (free)
2. Find the deploy area on your dashboard
3. Drag your project folder onto it
4. Get a live URL in seconds: `https://random-name-123.netlify.app`

This is the fastest path from code to live site. The trade-off: updating requires re-dragging.

### Option B — connect a repo (recommended)

1. Click **Add new site → Import an existing project**
2. Connect your GitHub account
3. Select the repository
4. Configure:
   - **Build command:** leave empty for static sites (no build step)
   - **Publish directory:** `.` or the folder containing `index.html`
5. Click **Deploy site**

Now every `git push` auto-deploys. Netlify also provides:

- **Rollback** — one-click revert to any previous deploy
- **Deploy log** — see exactly what happened during deployment
- **Branch deploys** — preview branches before merging

:::key
Connecting a repo to Netlify means your deploy process is: write code → commit → push → site updates. No manual steps, no FTP, no copying files. This is the professional workflow.
:::

## Vercel — connect repo, auto-deploy

Vercel is similar to Netlify, with a focus on frontend frameworks (Next.js, SvelteKit, etc.). It works great for static sites too.

### Setup

1. Go to **vercel.com** and sign up with your GitHub account
2. Click **Add New → Project**
3. Import your GitHub repository
4. For a static site, leave the framework preset as **Other**
5. Click **Deploy**

Your site is live at `https://your-project.vercel.app`.

```text
Push to main   → production deploy  (your-project.vercel.app)
Push to branch → preview deploy     (your-project-git-branch-name.vercel.app)
```

## Preview deploys from branches

Both Netlify and Vercel create **preview deploys** when you push a non-`main` branch. This is incredibly useful:

```bash
git switch -c feature/new-hero
# make changes
git push -u origin feature/new-hero
```

Netlify/Vercel automatically builds and deploys that branch to a temporary URL:

```text
https://feature-new-hero--your-site.netlify.app
```

You can share this URL with a teammate or client to review the changes before merging. When the branch merges, the preview URL stops mattering and the production site updates.

:::tip
Preview deploys are one of the best tools for collaboration. Instead of saying "check out branch X and run it locally," you send a link. Anyone with a browser can review.
:::

## Custom domains

All three platforms (GitHub Pages, Netlify, Vercel) support custom domains.

### The process

1. **Buy a domain** from a registrar (Namecheap, Cloudflare, Google Domains) — usually $10–15/year
2. **Add the domain** in your platform's dashboard
3. **Update DNS records** at your registrar — the platform tells you exactly what to add
4. **Wait for propagation** — usually minutes, sometimes up to an hour
5. **HTTPS is automatic** — all three platforms issue SSL certificates for free

```text
Before: https://random-name-123.netlify.app
After:  https://janedoe.dev
```

A custom domain is optional for learning projects but strongly recommended for your portfolio.

## Environment variables

When your project needs secrets (API keys, configuration) that shouldn't be in source code, deployment platforms let you set **environment variables**.

**Netlify:** Site settings → Environment variables → Add

**Vercel:** Settings → Environment Variables → Add

```text
Variable name:  API_KEY
Value:          sk-abc123-real-secret-key
```

These are injected at build time or runtime, depending on your setup. They never appear in your source code or Git history.

:::warning
Environment variables on static sites are tricky. Since your HTML/CSS/JS files are sent directly to the browser, any variable embedded in JavaScript during build is visible in the source. For true server-side secrets, you need a backend or serverless functions.
:::

## What CI/CD means

**CI/CD** stands for **Continuous Integration / Continuous Deployment.**

- **Continuous Integration (CI):** every time you push code, automated checks run — tests, linting, type checking. If anything fails, you're notified before merging.

- **Continuous Deployment (CD):** when code passes all checks and merges to `main`, it's automatically deployed to production. No manual steps.

```text
Developer pushes code
    ↓
CI runs tests automatically
    ↓
Tests pass? → Merge allowed
    ↓
Merge to main triggers CD
    ↓
Site is live
```

When you connect a repo to Netlify or Vercel, you already have **CD** — pushes to `main` auto-deploy. To add **CI**, you configure automated checks.

## A simple CI with GitHub Actions

GitHub Actions lets you run automated checks on every push or pull request. Here's a minimal workflow that runs a linter:

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npm run lint
```

This runs on every push to `main` and every pull request targeting `main`. If the lint check fails, you see a red ❌ on the PR.

```text
Pull Request: "Add quiz timer feature"
  ✅ CI / lint — All checks passed
  → Safe to merge
```

:::key
CI/CD removes human error from the deployment process. Instead of remembering to run tests and manually uploading files, the system does it automatically on every push. The more automated your workflow, the fewer mistakes reach production.
:::

## Choosing a platform — comparison

| Feature | GitHub Pages | Netlify | Vercel |
|---------|-------------|---------|--------|
| **Static sites** | ✅ | ✅ | ✅ |
| **Auto-deploy on push** | ✅ | ✅ | ✅ |
| **Preview deploys** | ❌ | ✅ | ✅ |
| **Custom domains** | ✅ | ✅ | ✅ |
| **Free HTTPS** | ✅ | ✅ | ✅ |
| **Serverless functions** | ❌ | ✅ | ✅ |
| **Drag-and-drop** | ❌ | ✅ | ❌ |
| **Best for** | Simple static sites | General frontend | Framework-heavy projects |

For your projects right now, any of them works. GitHub Pages is simplest. Netlify gives you the most flexibility. Vercel shines if you later adopt Next.js or SvelteKit.

## The deployment workflow summary

```text
1. Code → commit → push (you do this)
2. Platform detects the push (automatic)
3. CI runs checks if configured (automatic)
4. Site builds and deploys (automatic)
5. Live at your URL (automatic)
```

Your entire job is step 1. Everything else is handled.

:::quiz
Q: What does "Continuous Deployment" mean in practice?
- You deploy once and never update the site
- Every merge to the main branch automatically deploys to production without manual steps *
- You continuously write deployment scripts
- You deploy to multiple servers at the same time
E: Continuous Deployment automates the release process. When code merges to `main` and passes all checks, it goes live automatically. This removes manual deployment steps and ensures the live site always matches the latest code.
:::

:::quiz
Q: What is the main advantage of preview deploys on Netlify or Vercel?
- They make the site load faster
- They let you (and reviewers) see branch changes at a live URL before merging to production *
- They backup your code automatically
- They replace the need for Git branches
E: Preview deploys create a temporary live URL for each branch, allowing anyone to review changes in a real browser without cloning the repo or running code locally. This makes code review faster and more collaborative.
:::

## Recap

- **GitHub Pages:** Settings → Pages → deploy from `main`. Free, simple, auto-deploys on push.
- **Netlify:** connect a repo for auto-deploy, or drag-and-drop for instant publishing. Preview deploys on branches.
- **Vercel:** connect a repo, auto-deploy with preview deploys. Excellent for frameworks.
- **Preview deploys** let reviewers see branch changes at a live URL before merging.
- **Custom domains** are easy to configure and cost $10–15/year.
- **Environment variables** keep secrets out of source code.
- **CI/CD** automates testing (CI) and deployment (CD) — push code and everything else happens automatically.
- **GitHub Actions** lets you add CI checks to every push and pull request.

**Next up:** your projects — take everything you've learned about tooling and workflow and apply it to your portfolio.
