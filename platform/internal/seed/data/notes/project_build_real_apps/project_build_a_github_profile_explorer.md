# Project: Build a GitHub Profile Explorer

This project turns fetch, async JavaScript, rendering, and product thinking into one strong portfolio-ready build.

:::project
**Goal:** Build an app where a user searches for a GitHub username, sees the profile details, and browses recent repositories. Your app should handle loading, success, empty, and error states cleanly.
:::

## Why this project matters

A lot of beginner API projects stop at “fetch data and print it.”

A stronger frontend project asks better questions:

- what happens while the request is loading?
- what happens if the username does not exist?
- what happens if one request succeeds and another fails?
- how do we present data clearly once it arrives?

This project teaches that full product flow.

## What the app should include

A solid first version should support:

- a search form for GitHub usernames
- profile rendering: avatar, name, bio, followers, following
- repository rendering: name, stars, language, link
- loading and error states
- an empty state for no repositories
- a saved recent-search list in Local Storage

## Step 1 — Build the interface shell

```html
<main class="explorer-shell">
  <section class="explorer-card">
    <h1>GitHub Profile Explorer</h1>

    <form id="search-form">
      <label for="username">GitHub username</label>
      <div class="search-row">
        <input id="username" type="text" placeholder="e.g. gaearon" />
        <button type="submit">Search</button>
      </div>
    </form>

    <p id="status"></p>

    <section id="profile" hidden></section>

    <section>
      <h2>Repositories</h2>
      <ul id="repo-list"></ul>
      <p id="repo-empty" hidden>No repositories to show.</p>
    </section>

    <section>
      <h2>Recent searches</h2>
      <ul id="recent-list"></ul>
    </section>
  </section>
</main>
```

This gives your JavaScript clear places to render each state.

## Step 2 — Style for readability first

Use a clean card layout and keep the content easy to scan.

```css
body {
  margin: 0;
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f8fafc;
  font-family: system-ui, sans-serif;
  color: #0f172a;
}

.explorer-card {
  width: min(100% - 2rem, 52rem);
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 14px 36px rgba(15, 23, 42, 0.08);
}

.search-row {
  display: flex;
  gap: 0.75rem;
}

#repo-list {
  list-style: none;
  padding: 0;
  margin: 1rem 0 0;
}
```

You can refine the visual system once the loading states and data flow are correct.

## Step 3 — Know the API endpoints

For a username like `gaearon`, the profile endpoint is:

```text
https://api.github.com/users/gaearon
```

The repositories endpoint is:

```text
https://api.github.com/users/gaearon/repos?sort=updated&per_page=6
```

Those can be fetched in parallel.

## Step 4 — Fetch profile and repositories together

This is a perfect `Promise.all()` project.

```js
async function loadUser(username) {
  const [profileResponse, reposResponse] = await Promise.all([
    fetch(`https://api.github.com/users/${username}`),
    fetch(`https://api.github.com/users/${username}/repos?sort=updated&per_page=6`),
  ]);

  if (!profileResponse.ok) {
    throw new Error("User not found");
  }

  if (!reposResponse.ok) {
    throw new Error("Could not load repositories");
  }

  const profile = await profileResponse.json();
  const repos = await reposResponse.json();

  return { profile, repos };
}
```

This is stronger than making the requests one after another because the profile and repos do not depend on each other.

## Step 5 — Render the loading, success, and error states

This is where product quality shows up.

A good flow is:

1. clear old results
2. show a loading message
3. fetch data
4. render profile and repos on success
5. show a helpful error on failure

```js
async function handleSearch(username) {
  status.textContent = "Loading profile…";
  profileSection.hidden = true;
  repoList.innerHTML = "";
  repoEmpty.hidden = true;

  try {
    const data = await loadUser(username);
    status.textContent = "";
    renderProfile(data.profile);
    renderRepos(data.repos);
  } catch (error) {
    status.textContent = error.message;
  }
}
```

## Step 6 — Render the profile card

A useful profile summary can include:

- avatar
- display name
- username
- bio
- followers / following / public repos
- a link to the GitHub profile

Render it from the returned data object.

```js
function renderProfile(profile) {
  profileSection.hidden = false;
  profileSection.innerHTML = `
    <article class="profile-card">
      <img src="${profile.avatar_url}" alt="Avatar for ${profile.login}" />
      <div>
        <h2>${profile.name || profile.login}</h2>
        <p>@${profile.login}</p>
        <p>${profile.bio || "No bio provided."}</p>
        <p>${profile.followers} followers · ${profile.following} following</p>
      </div>
    </article>
  `;
}
```

## Step 7 — Render repositories thoughtfully

Do more than dump names in a list.

Each repo card can show:

- repository name
- description
- language
- star count
- link

```js
function renderRepos(repos) {
  repoList.innerHTML = "";

  if (repos.length === 0) {
    repoEmpty.hidden = false;
    return;
  }

  repoEmpty.hidden = true;

  repos.forEach((repo) => {
    const item = document.createElement("li");
    item.innerHTML = `
      <article class="repo-card">
        <h3><a href="${repo.html_url}" target="_blank" rel="noreferrer">${repo.name}</a></h3>
        <p>${repo.description || "No description provided."}</p>
        <p>${repo.language || "Unknown language"} · ★ ${repo.stargazers_count}</p>
      </article>
    `;
    repoList.appendChild(item);
  });
}
```

## Step 8 — Save recent searches

This makes the app feel more like a real tool.

Store the last few searched usernames in Local Storage.

```js
function saveRecent(username) {
  const current = JSON.parse(localStorage.getItem("recent-github-users") || "[]");
  const next = [username, ...current.filter((name) => name !== username)].slice(0, 5);
  localStorage.setItem("recent-github-users", JSON.stringify(next));
}
```

Then render the recent list and let users click a past search to re-run it.

## Step 9 — Handle real-world edge cases

A stronger final version should think about:

- blank usernames
- users with no bio
- users with no public repos
- GitHub rate limits
- slow network conditions
- keeping the UI readable during rapid searches

:::warning
GitHub’s public API can rate-limit unauthenticated requests. That is not a reason to avoid the project. It is part of learning how real APIs behave.
:::

## Stretch goals

Once the core version works, extend it with:

- tabs for overview / repositories / starred stats
- sorting repositories by stars
- filtering repositories by language
- a dark-mode toggle
- a "favorite developers" saved list
- skeleton loading placeholders

## What this project demonstrates

A polished version of this app shows that you can:

- work with external APIs
- coordinate multiple async requests
- build clear loading and error states
- render structured, real-world data
- persist useful UI history

That is excellent portfolio material.

:::quiz
Q: Why is `Promise.all()` a good fit for loading both the profile and the repositories?
- Because it automatically caches the results
- Because those requests are independent and can run in parallel *
- Because GitHub requires two requests at the same time
E: The profile and repository requests do not depend on one another, so starting them together reduces waiting time and keeps the app responsive.
:::

## What “done” looks like

A finished version should feel like a real developer-facing product:

- the search flow is fast and clear
- results are readable
- errors are handled without breaking the page
- recent searches make the tool feel useful
- the UI works well on both narrow and wide screens