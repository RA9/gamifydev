// proctor.js — invigilates a live placement sitting.
//
// What this can and cannot do, plainly, because the difference matters:
//
//   Copying      — copy, cut, right-click and the keyboard shortcuts are
//                  blocked and reported. Someone with devtools or "view source"
//                  can still read the text. This raises the cost, not a wall.
//   Leaving      — reliably detected. A tab switch, another window, or a
//                  minimise all surface as visibilitychange or blur.
//   Screenshots  — CANNOT be prevented by any web page. The operating system
//                  takes them; the browser is never asked. PrintScreen reaches
//                  us on Windows and is reported; macOS captures Cmd-Shift-3/4
//                  before the page sees anything. Printing and Save-As are
//                  blocked, which closes the easy routes.
//
// Detection is the browser's job; the count and the verdict are the server's.
// Everything here is suppressible by a determined candidate — which only ever
// means fewer strikes recorded, never one removed.
(function () {
  var paper = document.querySelector("[data-proctored]");
  if (!paper) return;

  var limit = parseInt(paper.dataset.strikeLimit, 10) || 2;
  var overlay = document.getElementById("plProctor");
  var titleEl = document.getElementById("plProctorTitle");
  var bodyEl = document.getElementById("plProctorBody");
  var dismissEl = document.getElementById("plProctorDismiss");
  var strikeEl = document.getElementById("plStrikes");
  var closed = false;

  var REASONS = {
    copy: "Copying the questions is not allowed.",
    hidden: "You left the test.",
    blur: "You left the test.",
    print: "Printing or saving the test is not allowed.",
    capture: "Screen capture is not allowed.",
  };

  // One act, one strike. A single alt-tab fires blur and visibilitychange
  // together, and a candidate should not lose their sitting to two names for
  // the same second. Reports inside this window collapse into the first.
  var COALESCE_MS = 1200;
  var lastAt = 0;

  function show(title, body, dismissable) {
    if (!overlay) return;
    titleEl.textContent = title;
    bodyEl.textContent = body;
    dismissEl.hidden = !dismissable;
    overlay.hidden = false;
    if (dismissable) dismissEl.focus();
  }

  function report(kind) {
    if (closed) return;
    var now = Date.now();
    if (now - lastAt < COALESCE_MS) return;
    lastAt = now;

    var body = new URLSearchParams();
    body.set("kind", kind);
    fetch("/placement/violation", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: body.toString(),
      credentials: "same-origin",
      keepalive: true, // so a report survives the tab being closed
    })
      .then(function (res) {
        return res.ok ? res.json() : null;
      })
      .then(function (out) {
        if (!out) return;
        if (strikeEl) strikeEl.textContent = out.strikes;
        if (out.voided) {
          closed = true;
          show(
            "Test closed",
            REASONS[kind] +
              " That was strike " + out.strikes + " of " + out.limit +
              ". Your placement has been closed and marked as failed for cheating.",
            false
          );
          // The server has already scored and closed the sitting; this only
          // takes the candidate to the page that explains it.
          setTimeout(function () {
            window.location.href = "/placement/result";
          }, 2600);
          return;
        }
        show(
          "Warning " + out.strikes + " of " + out.limit,
          REASONS[kind] +
            " This is a warning. Do it once more and the test closes and is" +
            " marked as failed.",
          true
        );
      })
      .catch(function () {
        /* A report that cannot be delivered is a strike we do not get. */
      });
  }

  if (dismissEl) {
    dismissEl.addEventListener("click", function () {
      overlay.hidden = true;
    });
  }

  // --- copying -------------------------------------------------------------
  ["copy", "cut"].forEach(function (evt) {
    paper.addEventListener(evt, function (e) {
      e.preventDefault();
      report("copy");
    });
  });
  paper.addEventListener("contextmenu", function (e) {
    e.preventDefault();
    report("copy");
  });
  // Dragging selected text out is a copy by another name.
  paper.addEventListener("dragstart", function (e) {
    e.preventDefault();
  });

  // --- leaving -------------------------------------------------------------
  document.addEventListener("visibilitychange", function () {
    if (document.visibilityState === "hidden") report("hidden");
  });
  window.addEventListener("blur", function () {
    report("blur");
  });

  // --- printing and capture ------------------------------------------------
  window.addEventListener("beforeprint", function () {
    report("print");
  });
  document.addEventListener("keydown", function (e) {
    var key = (e.key || "").toLowerCase();
    var meta = e.ctrlKey || e.metaKey;

    // Print and Save-As: the two routes to a clean copy of the paper.
    if (meta && (key === "p" || key === "s")) {
      e.preventDefault();
      report("print");
      return;
    }
    // The macOS capture shortcuts. The system usually swallows these before we
    // see them, so treat a hit as a bonus rather than a guarantee.
    if (e.metaKey && e.shiftKey && ["3", "4", "5"].indexOf(key) !== -1) {
      e.preventDefault();
      report("capture");
      return;
    }
    // Copy, select-all and view-source, blocked at the keyboard as well as the
    // event, since a blocked copy event still lets the selection be dragged.
    if (meta && (key === "c" || key === "a" || key === "u")) {
      e.preventDefault();
      if (key === "c") report("copy");
    }
  });
  // Windows PrintScreen reports on keyup, with no key event we can cancel.
  document.addEventListener("keyup", function (e) {
    if (e.key === "PrintScreen") report("capture");
  });
})();
