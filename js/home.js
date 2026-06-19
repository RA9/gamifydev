// import { scratchPage } from "./modules.js";
// import QuizPage from "./quiz.js";
// import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
// import { randomID } from "./utils.js";

async function HomePage(htmlEl) {
  let state = await DB.states.where("name").equals("general").last();

  if (!state) {
    state = await createStorage("states", {
      id: randomID(),
      previous: null,
      current: "home",
      next: "user_info",
      name: "general",
    });
  }

  if (state.current === "home") {
    htmlEl.innerHTML = `
  <div class="space-y-8 animate-fade-up">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 via-brand-500 to-brand-700 text-white shadow-soft">
      <div class="absolute -top-10 -right-10 h-48 w-48 rounded-full bg-grass-400/30 blur-2xl"></div>
      <div class="absolute -bottom-12 -left-8 h-40 w-40 rounded-full bg-white/10 blur-2xl"></div>
      <div class="relative grid gap-8 p-8 sm:p-12 lg:grid-cols-2 lg:items-center">
        <div>
          <span class="gd-chip bg-white/15 text-white mb-4">🎮 Learn by playing</span>
          <h1 class="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-white leading-tight">
            Level up your<br/>coding skills.
          </h1>
          <p class="mt-4 text-lg text-white/90 max-w-md">
            GamifyDev turns learning to code into a game — earn XP, keep streaks,
            unlock badges, and build real projects from scratch.
          </p>
          <div class="mt-7 flex flex-wrap gap-3">
            <button id="get-started" class="gd-btn gd-btn-amber">Get Started — it's free</button>
            <a href="#test" class="gd-btn gd-btn-secondary">Take a quick quiz</a>
          </div>
          <div class="mt-6 flex flex-wrap gap-x-6 gap-y-2 text-sm text-white/80 font-bold">
            <span>⚡ XP & Levels</span>
            <span>🔥 Daily Streaks</span>
            <span>🏆 Badges</span>
            <span>📚 Guided Lessons</span>
          </div>
        </div>
        <div class="hidden lg:flex justify-center">
          <div class="animate-float text-[10rem] leading-none select-none">👩‍💻</div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="grid gap-4 sm:grid-cols-3">
      <div class="gd-card-sm flex items-start gap-4">
        <div class="text-3xl">🎯</div>
        <div>
          <h3 class="font-extrabold text-lg">Pick a path</h3>
          <p class="text-slate-600 text-sm mt-1">Frontend, backend, or fullstack — start where you are.</p>
        </div>
      </div>
      <div class="gd-card-sm flex items-start gap-4">
        <div class="text-3xl">🧩</div>
        <div>
          <h3 class="font-extrabold text-lg">Learn & practice</h3>
          <p class="text-slate-600 text-sm mt-1">Bite-size lessons and quizzes that actually stick.</p>
        </div>
      </div>
      <div class="gd-card-sm flex items-start gap-4">
        <div class="text-3xl">🚀</div>
        <div>
          <h3 class="font-extrabold text-lg">Build for real</h3>
          <p class="text-slate-600 text-sm mt-1">Apply the fundamentals to real-world projects.</p>
        </div>
      </div>
    </section>

    <!-- Creator note -->
    <section class="gd-card">
      <div class="flex items-center gap-3 mb-3">
        <div class="grid h-11 w-11 place-items-center rounded-full bg-brand-100 text-2xl">💬</div>
        <h2 class="text-xl font-extrabold">A message from the creator</h2>
      </div>
      <p class="text-slate-600 leading-relaxed">
        <b>Welcome, fellow web explorers!</b> I'm
        <a class="font-bold text-brand-600 underline decoration-2 underline-offset-2" href="https://twitter.com/rademejs" target="_blank" rel="noopener">Carlos S. Nah</a>,
        a software engineer passionate about empowering fellow developers. Navigating the
        web-dev world can be tough — but you're not alone. We're building a platform focused
        on boosting your mental agility and problem-solving skills. Ready to dive deeper?
        Check out
        <a href="https://kit.kwagei.com" class="font-bold text-brand-600 underline decoration-2 underline-offset-2" target="_blank" rel="noopener">Kwagei Innovators Training</a>
        for more.
      </p>
      <p class="mt-4 rounded-2xl bg-amber-50 border border-amber-200 px-4 py-3 text-sm text-amber-800">
        <b>Heads up:</b> this isn't a get-rich-quick scheme. We give you the fundamentals of
        programming and help you build real projects — what you do with that is up to you.
      </p>
    </section>
  </div>
    `;

    const GET_STARTED_BUTTON = document.querySelector("#get-started");

    // Show the preference card
    GET_STARTED_BUTTON.addEventListener("click", () => {
      state.current = "user_info";
      state.previous = "home";
      state.next = "preference";
      const page = document.querySelector("main");
      updateStorage("states", state);
      UserInfoSection(page);
    });
  } else if (state.current === "user_info") {
    const page = document.querySelector("main");
    UserInfoSection(page);
  } else if (state.current === "preference") {
    const page = document.querySelector("main");
    PreferenceSection(page);
  } else if (
    state.current === "quiz" ||
    state.current === "quiz-started" ||
    state.current === "quiz-completed"
  ) {
    // Legacy quiz states route into the working "Test Yourself" flow (tys.js).
    window.location.hash = "test";
  } else if (state.current === "scratch" || state.current === "note") {
    const page = document.querySelector("main");

    scratchPage(page);
  }
}

async function UserInfoSection(htmlEl) {
  let state = await DB.states.where("name").equals("general").last();

  if (state.current === "user_info") {
    htmlEl.innerHTML = `
  <div class="max-w-lg mx-auto animate-fade-up">
    <div class="gd-card">
      <div class="text-center mb-6">
        <div class="text-4xl mb-2">👋</div>
        <h1 class="text-2xl font-extrabold">Let's get you set up</h1>
        <p class="text-slate-500 mt-1">Tell us your name and what you'd like to focus on.</p>
      </div>
      <div class="space-y-4">
        <div>
          <label for="name" class="gd-label">Your name</label>
          <input type="text" id="name" placeholder="e.g. Ada Lovelace" class="gd-input" />
        </div>
        <div>
          <label for="preference" class="gd-label">Preferred field</label>
          <select id="preference" class="gd-select">
            <option value="">Choose your field…</option>
            <option value="fullstack">Fullstack</option>
            <option value="frontend">Frontend</option>
            <option value="backend">Backend</option>
            <option value="notsure">Not sure yet</option>
          </select>
        </div>
        <button id="submit-name" class="gd-btn gd-btn-primary gd-btn-block">Continue</button>
      </div>
    </div>
  </div>
  `;
  }

  const SUBMIT_NAME_BUTTON = document.querySelector("#submit-name");

  SUBMIT_NAME_BUTTON.addEventListener("click", async () => {
    const name = document.querySelector("#name").value;
    const preference = document.querySelector("#preference").value;
    let user = await DB.users.where("name").equals(name).last();

    if (name.length < 4) {
      return alert("Your name should be at least 4 characters!");
    }

    if (!user) {
      user = await createStorage("users", {
        id: randomID(),
        name,
        preference,
        // test: {},
      });
    }

    state.current = "preference";
    // user.name = name;
    // user.preference = preference;
    await updateStorage("states", state);
    // await createStorage("users", user);
    const page = document.querySelector("main");
    PreferenceSection(page);
  });
}

async function PreferenceSection(htmlEl) {
  const user = (await DB.users.toArray())[0];
  const state = await DB.states.where("name").equals("general").last();

  if (state.current === "preference") {
    htmlEl.innerHTML = `
  <div class="max-w-3xl mx-auto animate-fade-up">
    <div class="text-center mb-6">
      <span class="gd-chip gd-chip-brand mb-3">${user.preference.toUpperCase()} PATH</span>
      <h1 class="text-2xl sm:text-3xl font-extrabold">Welcome aboard, ${user.name}! 🎉</h1>
      <p class="text-slate-500 mt-2 max-w-xl mx-auto">
        You're all set on the <b>${user.preference.toUpperCase()}</b> path. How would you like to begin?
      </p>
    </div>
    <div class="grid gap-4 sm:grid-cols-2">
      <button id="scratch" class="gd-card text-left hover:shadow-soft hover:-translate-y-0.5 transition-all">
        <div class="text-4xl mb-3">📚</div>
        <h2 class="text-xl font-extrabold">Start from scratch</h2>
        <p class="text-slate-600 mt-1 text-sm">Follow guided lessons step by step and build up from the basics.</p>
        <span class="gd-chip gd-chip-grass mt-4">Recommended</span>
      </button>
      <button id="get-started" class="gd-card text-left hover:shadow-soft hover:-translate-y-0.5 transition-all">
        <div class="text-4xl mb-3">⚡</div>
        <h2 class="text-xl font-extrabold">Test yourself</h2>
        <p class="text-slate-600 mt-1 text-sm">Already know some things? Jump into a timed quiz and earn XP now.</p>
        <span class="gd-chip gd-chip-slate mt-4">Quick challenge</span>
      </button>
    </div>
  </div>
  `;

    const TAKE_QUIZ_BUTTON = document.querySelector("#get-started");
    const START_FROM_SCRATCH_BUTTON = document.querySelector("#scratch");

    TAKE_QUIZ_BUTTON.addEventListener("click", () => {
      // The interactive quiz lives in the "Test Yourself" flow (tys.js).
      window.location.hash = "test";
    });

    START_FROM_SCRATCH_BUTTON.addEventListener("click", () => {
      state.current = "scratch";
      state.previous = "preference";
      state.next = null;

      console.log({ state });

      // window.location.href = "scratch.html";
      updateStorage("states", state);
      const page = document.querySelector("main");
      scratchPage(page);
    });
  } else if (state.current === "scratch" || state.current === "note") {
    const page = document.querySelector("main");
    scratchPage(page);
  }
}

// export default HomePage;
