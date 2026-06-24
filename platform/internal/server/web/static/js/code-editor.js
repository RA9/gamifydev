// <gd-code-editor> — an interactive code editor built with tan-compose
// (the project's web-components library).
//
// Usage (the slotted <textarea> stays in light DOM, so it submits with the form
// natively — and remains usable even if this script fails to load):
//
//   <gd-code-editor lang="python">
//     <textarea name="code" class="gd-editor-ta">starter</textarea>
//   </gd-code-editor>
//
// The component adds a toolbar (language + line count), a line-number gutter,
// tab-to-indent, and auto-grow. It wires listeners directly on the light-DOM
// textarea in afterMount, so no re-render (and no cursor jumps) on every keypress.
import { describe, build, html } from "https://cdn.jsdelivr.net/gh/RA9/tan-compose@v1.3.0/dist/mod.js";

build(
  "gd-code-editor",
  describe({
    props: { lang: { type: "string", default: "code" } },
    template: ({ props }) => html`
      <style>
        :host { display:block; }
        .ed { border:2px solid #e2e8f0; border-radius:14px; overflow:hidden; background:#fff;
          font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace; transition:border-color .12s; }
        .ed:focus-within { border-color:#8b6df3; }
        .ed-bar { display:flex; align-items:center; justify-content:space-between; padding:7px 14px;
          background:#f8fafc; border-bottom:1px solid #eef2f7; }
        .ed-lang { font-size:.72rem; font-weight:800; letter-spacing:.04em; text-transform:uppercase; color:#7857f7; }
        .ed-lines { font-size:.72rem; font-weight:700; color:#94a3b8; }
        .ed-body { display:flex; align-items:stretch; max-height:460px; overflow:auto; }
        .ed-gutter { padding:14px 10px; margin:0; text-align:right; color:#cbd5e1; background:#fcfcfe;
          white-space:pre; line-height:1.65; font-size:.88rem; user-select:none; border-right:1px solid #f1f5f9; min-width:2.2em; }
        ::slotted(textarea) { flex:1; border:none !important; outline:none !important; resize:none !important;
          padding:14px !important; margin:0; line-height:1.65; font-size:.88rem; font-family:inherit;
          color:#0f172a; background:#fff; min-height:0 !important; overflow:hidden; }
      </style>
      <div class="ed">
        <div class="ed-bar"><span class="ed-lang">${props.lang}</span><span class="ed-lines"></span></div>
        <div class="ed-body"><pre class="ed-gutter"></pre><slot></slot></div>
      </div>`,
    afterMount: function () {
      const host = this;
      const sr = host.shadowRoot;
      const ta = host.querySelector("textarea");
      if (!ta || !sr) return;
      const gutter = sr.querySelector(".ed-gutter");
      const lines = sr.querySelector(".ed-lines");

      const refresh = () => {
        const n = ta.value.split("\n").length || 1;
        let g = "";
        for (let i = 1; i <= n; i++) g += i + "\n";
        gutter.textContent = g;
        lines.textContent = n + (n === 1 ? " line" : " lines");
        ta.style.height = "auto";
        ta.style.height = ta.scrollHeight + "px";
      };

      ta.addEventListener("input", refresh);
      ta.addEventListener("keydown", (e) => {
        if (e.key === "Tab") {
          e.preventDefault();
          const s = ta.selectionStart;
          const en = ta.selectionEnd;
          ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(en);
          ta.selectionStart = ta.selectionEnd = s + 2;
          refresh();
        }
      });
      // Defer once so layout is settled before measuring scrollHeight.
      requestAnimationFrame(refresh);
    },
  })
);
