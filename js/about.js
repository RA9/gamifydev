function AboutPage(htmlEl) {
  const features = [
    { icon: "zap", tint: "brand", title: "XP & Levels", desc: "Earn points for every correct answer and watch your level climb." },
    { icon: "flame", tint: "amber", title: "Daily Streaks", desc: "Build a habit — keep your streak alive by practicing each day." },
    { icon: "trophy", tint: "grass", title: "Badges", desc: "Unlock achievements as you hit milestones and master new skills." },
    { icon: "book", tint: "brand", title: "Guided Lessons", desc: "Bite-size, project-based lessons that take you from zero to building." },
  ];

  const steps = [
    { n: "1", icon: "target", title: "Pick a path", desc: "Frontend, backend, or fullstack — start where you are." },
    { n: "2", icon: "puzzle", title: "Learn & practice", desc: "Work through lessons and quizzes that actually stick." },
    { n: "3", icon: "rocket", title: "Build for real", desc: "Apply the fundamentals to real-world projects." },
  ];

  const tints = {
    brand: "bg-brand-100 text-brand-600",
    grass: "bg-brand-100 text-brand-600",
    amber: "bg-brand-100 text-brand-600",
  };

  htmlEl.innerHTML = `
  <div class="max-w-5xl mx-auto space-y-6 animate-fade-up">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-8 sm:p-10 text-center">
      <div class="absolute -top-10 -right-8 h-40 w-40 rounded-full bg-brand-300/30 blur-2xl"></div>
      <div class="relative">
        <span class="gd-chip bg-white/15 text-white mb-3">${icon("gamepad", "w-4 h-4")} About GamifyDev</span>
        <h1 class="text-3xl sm:text-4xl font-extrabold text-white">Learn to code, the fun way</h1>
        <p class="mt-3 text-white/90 max-w-2xl mx-auto">
          A project-based learning platform that turns the hardest part of coding — getting started — into a game you actually want to play.
        </p>
      </div>
    </section>

    <!-- Intro cards -->
    <div class="grid gap-5 md:grid-cols-2">
      <div class="gd-card">
        <div class="flex items-center gap-3 mb-2">
          <span class="grid h-11 w-11 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("gamepad", "w-6 h-6")}</span>
          <h2 class="text-xl font-extrabold">What is GamifyDev?</h2>
        </div>
        <p class="text-slate-600 leading-relaxed">
          GamifyDev teaches you how to code by building real-world projects. Through an interactive,
          game-like environment, you'll pick up the fundamentals in a fun way — no prior programming
          knowledge and very little internet bandwidth required.
        </p>
      </div>
      <div class="gd-card">
        <div class="flex items-center gap-3 mb-2">
          <span class="grid h-11 w-11 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("target", "w-6 h-6")}</span>
          <h2 class="text-xl font-extrabold">Our Mission</h2>
        </div>
        <p class="text-slate-600 leading-relaxed">
          Learning something new is hardest at the very start. We make that first step lower,
          friendlier, and a lot more rewarding — so you actually keep going and build momentum.
        </p>
      </div>
    </div>

    <!-- What you get -->
    <section class="gd-card">
      <h2 class="text-xl font-extrabold mb-1">What you'll get</h2>
      <p class="text-slate-500 mb-5">Everything is designed to keep you motivated and moving forward.</p>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        ${features
          .map(
            (f) => `
          <div class="rounded-2xl border-2 border-slate-100 bg-slate-50 p-5 text-center transition-transform hover:-translate-y-0.5">
            <div class="grid h-12 w-12 mx-auto mb-3 place-items-center rounded-2xl ${tints[f.tint]}">${icon(f.icon, "w-6 h-6")}</div>
            <h3 class="font-extrabold">${f.title}</h3>
            <p class="text-sm text-slate-600 mt-1">${f.desc}</p>
          </div>`
          )
          .join("")}
      </div>
    </section>

    <!-- How it works -->
    <section class="gd-card">
      <h2 class="text-xl font-extrabold mb-5">How it works</h2>
      <div class="grid gap-4 sm:grid-cols-3">
        ${steps
          .map(
            (s) => `
          <div class="relative rounded-2xl border-2 border-slate-100 p-5">
            <span class="absolute -top-3 -left-2 grid h-8 w-8 place-items-center rounded-full bg-brand-500 text-white font-extrabold text-sm shadow-card">${s.n}</span>
            <div class="grid h-12 w-12 mb-2 place-items-center rounded-2xl bg-brand-50 text-brand-600">${icon(s.icon, "w-6 h-6")}</div>
            <h3 class="font-extrabold">${s.title}</h3>
            <p class="text-sm text-slate-600 mt-1">${s.desc}</p>
          </div>`
          )
          .join("")}
      </div>
    </section>

    <!-- Stats band -->
    <section class="grid grid-cols-3 gap-4">
      <div class="gd-card-sm text-center">
        <p class="text-3xl font-extrabold text-brand-600">7</p>
        <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Languages</p>
      </div>
      <div class="gd-card-sm text-center">
        <p class="text-3xl font-extrabold text-grass-600">100%</p>
        <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Free</p>
      </div>
      <div class="gd-card-sm text-center">
        <p class="text-3xl font-extrabold text-amber-500">∞</p>
        <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Practice</p>
      </div>
    </section>

    <!-- CTA -->
    <section class="rounded-3xl bg-gradient-to-r from-grass-500 to-grass-600 text-white shadow-soft p-8 text-center">
      <h2 class="text-2xl font-extrabold text-white">Ready to level up?</h2>
      <p class="text-white/90 mt-1 mb-5">Jump into a quick quiz or start a guided lesson — it's free.</p>
      <div class="flex flex-wrap justify-center gap-3">
        <a href="#test" class="gd-btn gd-btn-primary">Take a quiz</a>
        <a href="#" class="gd-btn gd-btn-secondary">Start learning</a>
      </div>
    </section>
  </div>
      `;
}

// export default AboutPage;
