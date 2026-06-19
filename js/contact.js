function ContactPage(htmlEl) {
  htmlEl.innerHTML = `
  <div class="max-w-2xl mx-auto animate-fade-up">
    <div class="text-center mb-6">
      <span class="gd-chip gd-chip-brand mb-3">Contact</span>
      <h1 class="text-3xl font-extrabold">Let's talk 👋</h1>
      <p class="text-slate-500 mt-2">
        Thanks for reaching out! We're a small team and we'll get back to you as soon as we can.
      </p>
    </div>
    <div class="grid gap-4 sm:grid-cols-2">
      <div class="gd-card-sm flex items-center gap-4">
        <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-brand-100 text-2xl">🌍</div>
        <div>
          <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Address</p>
          <p class="font-bold text-slate-700">We work remotely!</p>
        </div>
      </div>
      <a href="mailto:cnah27@gmail.com" class="gd-card-sm flex items-center gap-4 hover:shadow-soft hover:-translate-y-0.5 transition-all">
        <div class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-grass-100 text-2xl">✉️</div>
        <div>
          <p class="text-xs font-extrabold uppercase tracking-wide text-slate-500">Email</p>
          <p class="font-bold text-brand-600">cnah27@gmail.com</p>
        </div>
      </a>
    </div>
  </div>
    `;
}

// export default ContactPage;
