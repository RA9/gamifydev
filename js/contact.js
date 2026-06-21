function ContactPage(htmlEl) {
  htmlEl.innerHTML = `
  <div class="max-w-5xl mx-auto space-y-6 animate-fade-up">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-8 sm:p-10 text-center">
      <div class="absolute -bottom-10 -left-8 h-40 w-40 rounded-full bg-grass-400/30 blur-2xl"></div>
      <div class="relative">
        <span class="gd-chip bg-white/15 text-white mb-3">${icon("chat", "w-4 h-4")} Contact</span>
        <h1 class="text-3xl sm:text-4xl font-extrabold text-white">Let's talk</h1>
        <p class="mt-3 text-white/90 max-w-xl mx-auto">
          Questions, feedback, or just want to say hi? We're a small team and we'll get back to you as soon as we can.
        </p>
      </div>
    </section>

    <div class="grid gap-6 lg:grid-cols-5">
      <!-- Form -->
      <div class="lg:col-span-3">
        <form id="contact-form" class="gd-card space-y-4" novalidate>
          <h2 class="text-xl font-extrabold">Send us a message</h2>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label for="cf-name" class="gd-label">Name</label>
              <input id="cf-name" type="text" placeholder="Your name" class="gd-input" />
            </div>
            <div>
              <label for="cf-email" class="gd-label">Email</label>
              <input id="cf-email" type="email" placeholder="you@example.com" class="gd-input" />
            </div>
          </div>
          <div>
            <label for="cf-subject" class="gd-label">Subject</label>
            <input id="cf-subject" type="text" placeholder="What's this about?" class="gd-input" />
          </div>
          <div>
            <label for="cf-message" class="gd-label">Message</label>
            <textarea id="cf-message" rows="5" placeholder="Tell us what's on your mind…" class="gd-input resize-y"></textarea>
          </div>
          <p id="cf-error" class="hidden text-sm font-bold text-rose-500"></p>
          <button type="submit" class="gd-btn gd-btn-primary gd-btn-block">${icon("send", "w-4 h-4")} Send Message</button>
          <p class="text-xs text-slate-400 text-center">This opens your email app with the message ready to send.</p>
        </form>
      </div>

      <!-- Info sidebar -->
      <div class="lg:col-span-2 space-y-4">
        <a href="mailto:cnah27@gmail.com" class="gd-card-sm flex items-center gap-4 hover:shadow-soft hover:-translate-y-0.5 transition-all">
          <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("mail", "w-6 h-6")}</div>
          <div>
            <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Email</p>
            <p class="font-bold text-brand-600 break-all">cnah27@gmail.com</p>
          </div>
        </a>
        <div class="gd-card-sm flex items-center gap-4">
          <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("pin", "w-6 h-6")}</div>
          <div>
            <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Where we are</p>
            <p class="font-bold text-slate-700">We work remotely!</p>
          </div>
        </div>
        <div class="gd-card-sm flex items-center gap-4">
          <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("clock", "w-6 h-6")}</div>
          <div>
            <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Response time</p>
            <p class="font-bold text-slate-700">Usually within 1–2 days</p>
          </div>
        </div>
        <a href="https://twitter.com/gamifydev" target="_blank" rel="noopener" class="gd-card-sm flex items-center gap-4 hover:shadow-soft hover:-translate-y-0.5 transition-all">
          <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-slate-100 text-slate-600">${icon("twitter", "w-6 h-6")}</div>
          <div>
            <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Social</p>
            <p class="font-bold text-brand-600">@gamifydev</p>
          </div>
        </a>
      </div>
    </div>
  </div>
    `;

  // Compose an email from the form via a mailto: link — no backend needed.
  const form = document.querySelector("#contact-form");
  const error = document.querySelector("#cf-error");
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const name = document.querySelector("#cf-name").value.trim();
    const email = document.querySelector("#cf-email").value.trim();
    const subject = document.querySelector("#cf-subject").value.trim();
    const message = document.querySelector("#cf-message").value.trim();

    const showError = (msg) => {
      error.textContent = msg;
      error.classList.remove("hidden");
    };

    if (name.length < 2) return showError("Please enter your name.");
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
      return showError("Please enter a valid email address.");
    if (message.length < 5) return showError("Please enter a longer message.");

    error.classList.add("hidden");
    const body = `From: ${name} (${email})\n\n${message}`;
    window.location.href = `mailto:cnah27@gmail.com?subject=${encodeURIComponent(
      subject || "GamifyDev contact"
    )}&body=${encodeURIComponent(body)}`;
  });
}

// export default ContactPage;
