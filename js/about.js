function AboutPage(htmlEl) {
  htmlEl.innerHTML = `
  <div class="max-w-3xl mx-auto space-y-5 animate-fade-up">
    <div class="text-center">
      <span class="gd-chip gd-chip-brand mb-3">About</span>
      <h1 class="text-3xl font-extrabold">Learn to code, the fun way</h1>
    </div>
    <div class="gd-card">
      <div class="flex items-center gap-3 mb-2">
        <span class="text-2xl">🎮</span>
        <h2 class="text-xl font-extrabold">What is GamifyDev?</h2>
      </div>
      <p class="text-slate-600 leading-relaxed">
        GamifyDev is a project-based learning platform that teaches you how to code by building
        real-world projects. Through an interactive, game-like environment, you'll pick up the
        fundamentals in a fun way — no prior programming knowledge and very little internet
        bandwidth required.
      </p>
    </div>
    <div class="gd-card">
      <div class="flex items-center gap-3 mb-2">
        <span class="text-2xl">🎯</span>
        <h2 class="text-xl font-extrabold">Our Mission</h2>
      </div>
      <p class="text-slate-600 leading-relaxed">
        Learning something new is hardest at the very start. We make that first step lower,
        friendlier, and a lot more rewarding — so you actually keep going.
      </p>
    </div>
  </div>
      `;
}

// export default AboutPage;
