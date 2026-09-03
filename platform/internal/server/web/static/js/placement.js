// Placement diagnostic: a countdown and an answered-count.
//
// The timer here is a courtesy, not the enforcement. The deadline is stored on
// the attempt and checked server-side at submit, so a stopped clock, a closed
// tab, or a edited DOM buys no extra time. When the clock runs out we submit
// what the learner has so far rather than silently discarding the sitting.
(function () {
  var paper = document.querySelector(".pl-paper");
  if (!paper) return;

  var form = document.getElementById("plForm");
  var timerEl = document.getElementById("plTimer");
  var answeredEl = document.getElementById("plAnswered");
  var progressEl = document.getElementById("plProgress");
  var questions = [].slice.call(paper.querySelectorAll("[data-q]"));

  // The server writes "YYYY-MM-DD HH:MM:SS" in UTC.
  var raw = paper.dataset.expires || "";
  var expires = new Date(raw.replace(" ", "T") + "Z").getTime();

  function refreshCount() {
    var done = questions.filter(function (q) {
      return q.querySelector("input:checked");
    }).length;
    answeredEl.textContent = done;
    progressEl.style.width = questions.length
      ? (done / questions.length) * 100 + "%"
      : "0%";
  }

  form.addEventListener("change", function (e) {
    if (e.target && e.target.type === "radio") {
      var q = e.target.closest("[data-q]");
      if (q) q.classList.add("is-answered");
      refreshCount();
    }
  });

  var submitted = false;
  form.addEventListener("submit", function () {
    submitted = true;
  });

  function tick() {
    if (!expires || isNaN(expires)) {
      timerEl.textContent = "";
      return;
    }
    var left = Math.floor((expires - Date.now()) / 1000);
    if (left <= 0) {
      timerEl.textContent = "0:00";
      timerEl.classList.add("is-out");
      if (!submitted) {
        submitted = true;
        form.submit(); // hand in whatever is answered
      }
      return;
    }
    var m = Math.floor(left / 60);
    var s = left % 60;
    timerEl.textContent = m + ":" + (s < 10 ? "0" : "") + s;
    // Warn in the last five minutes.
    timerEl.classList.toggle("is-low", left <= 300);
    setTimeout(tick, 1000);
  }

  refreshCount();
  tick();
})();
