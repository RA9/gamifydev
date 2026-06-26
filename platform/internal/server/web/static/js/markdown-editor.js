// <gd-markdown-editor> — a split markdown editor built with tan-compose (the
// project's web-components library): a live, server-rendered preview on the left
// and a markdown textarea with a formatting toolbar on the right. Responsive:
// the two panes stack on narrow screens (editor first).
//
// The preview is rendered by the server (POST to the `endpoint` attribute) so it
// matches published output exactly — including :::tip / :::warning blocks.
//
// Usage (slotted light-DOM nodes stay in the form, so the textarea submits
// natively and the editor degrades to a plain textarea if this script fails):
//
//   <gd-markdown-editor endpoint="/admin/blog/preview">
//     <textarea slot="editor" name="body" class="md-input">...</textarea>
//     <div slot="preview" class="prose md-preview"></div>
//   </gd-markdown-editor>
import { describe, build, html } from "https://cdn.jsdelivr.net/gh/RA9/tan-compose@v1.3.0/dist/mod.js";

build(
  "gd-markdown-editor",
  describe({
    props: { endpoint: { type: "string", default: "/admin/blog/preview" } },
    template: ({ props }) => html`
      <style>
        :host { display:block; }
        .wrap { border:2px solid #e2e8f0; border-radius:14px; overflow:hidden; background:#fff; }
        .wrap:focus-within { border-color:#8b6df3; }
        .toolbar { display:flex; flex-wrap:wrap; gap:4px; padding:8px 10px;
          background:#f8fafc; border-bottom:1px solid #eef2f7; }
        .tb { font:inherit; font-size:.8rem; font-weight:700; color:#475569; cursor:pointer;
          background:#fff; border:1px solid #e2e8f0; border-radius:8px; padding:5px 9px; line-height:1; }
        .tb:hover { background:#f1f5f9; border-color:#cbd5e1; color:#1e293b; }
        .tb-sep { width:1px; align-self:stretch; background:#e2e8f0; margin:2px 4px; }
        .panes { display:grid; grid-template-columns:1fr 1fr; min-height:380px; }
        .pane { display:flex; flex-direction:column; min-width:0; }
        .pane-preview { border-right:1px solid #eef2f7; }
        .pane-label { font-size:.68rem; font-weight:800; letter-spacing:.05em; text-transform:uppercase;
          color:#94a3b8; padding:8px 14px 4px; }
        .scroll { flex:1; overflow:auto; max-height:62vh; }
        .scroll-preview { padding:4px 16px 18px; }
        ::slotted(textarea) { display:block; width:100%; box-sizing:border-box; flex:1; border:none !important;
          outline:none !important; resize:none !important; padding:8px 16px 16px !important; margin:0;
          min-height:340px; line-height:1.65; font-size:.9rem; color:#0f172a; background:#fff;
          font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace; }
        ::slotted(.md-preview) { min-height:340px; }
        @media (max-width:820px) {
          .panes { grid-template-columns:1fr; }
          .pane-preview { border-right:none; border-top:1px solid #eef2f7; order:2; }
          .pane-editor { order:1; }
          .scroll { max-height:none; }
        }
      </style>
      <div class="wrap">
        <div class="toolbar">
          <button class="tb" type="button" data-act="bold" title="Bold (Ctrl/Cmd+B)"><strong>B</strong></button>
          <button class="tb" type="button" data-act="italic" title="Italic (Ctrl/Cmd+I)"><em>I</em></button>
          <button class="tb" type="button" data-act="h2" title="Heading">H2</button>
          <button class="tb" type="button" data-act="h3" title="Subheading">H3</button>
          <span class="tb-sep"></span>
          <button class="tb" type="button" data-act="link" title="Link">Link</button>
          <button class="tb" type="button" data-act="code" title="Inline code">Code</button>
          <button class="tb" type="button" data-act="codeblock" title="Code block">{ }</button>
          <button class="tb" type="button" data-act="ul" title="Bullet list">List</button>
          <button class="tb" type="button" data-act="quote" title="Quote">Quote</button>
          <span class="tb-sep"></span>
          <button class="tb" type="button" data-act="tip" title="Tip callout">Tip</button>
          <button class="tb" type="button" data-act="warning" title="Warning callout">Warn</button>
        </div>
        <div class="panes">
          <div class="pane pane-preview">
            <span class="pane-label">Preview</span>
            <div class="scroll scroll-preview"><slot name="preview"></slot></div>
          </div>
          <div class="pane pane-editor">
            <span class="pane-label">Markdown</span>
            <div class="scroll"><slot name="editor"></slot></div>
          </div>
        </div>
      </div>`,
    afterMount: function () {
      const host = this;
      const sr = host.shadowRoot;
      const ta = host.querySelector("textarea");
      const preview = host.querySelector('[slot="preview"]');
      if (!ta || !sr) return;

      const endpoint = host.getAttribute("endpoint") || "/admin/blog/preview";

      // --- live preview (server-rendered, debounced) ---
      let timer = null;
      let lastSent = null;
      const renderPreview = () => {
        if (!preview || ta.value === lastSent) return;
        lastSent = ta.value;
        const params = new URLSearchParams();
        params.set("body", ta.value);
        fetch(endpoint, {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: params.toString(),
          credentials: "same-origin",
        })
          .then((r) => (r.ok ? r.text() : Promise.reject(r.status)))
          .then((htmlText) => {
            preview.innerHTML = htmlText;
          })
          .catch(() => {});
      };
      const schedulePreview = () => {
        clearTimeout(timer);
        timer = setTimeout(renderPreview, 350);
      };

      // --- toolbar / shortcuts ---
      const surround = (before, after) => {
        const s = ta.selectionStart;
        const e = ta.selectionEnd;
        const sel = ta.value.slice(s, e);
        ta.value = ta.value.slice(0, s) + before + sel + after + ta.value.slice(e);
        ta.focus();
        ta.selectionStart = s + before.length;
        ta.selectionEnd = s + before.length + sel.length;
      };
      const linePrefix = (prefix) => {
        const s = ta.selectionStart;
        const lineStart = ta.value.lastIndexOf("\n", s - 1) + 1;
        ta.value = ta.value.slice(0, lineStart) + prefix + ta.value.slice(lineStart);
        ta.focus();
        ta.selectionStart = ta.selectionEnd = s + prefix.length;
      };
      const block = (open, close) => {
        const s = ta.selectionStart;
        const e = ta.selectionEnd;
        const sel = ta.value.slice(s, e) || "Your content here";
        const snippet = `\n${open}\n${sel}\n${close}\n`;
        ta.value = ta.value.slice(0, s) + snippet + ta.value.slice(e);
        ta.focus();
        ta.selectionStart = ta.selectionEnd = s + snippet.length;
      };
      const actions = {
        bold: () => surround("**", "**"),
        italic: () => surround("*", "*"),
        code: () => surround("`", "`"),
        h2: () => linePrefix("## "),
        h3: () => linePrefix("### "),
        ul: () => linePrefix("- "),
        quote: () => linePrefix("> "),
        link: () => surround("[", "](https://)"),
        codeblock: () => block("```", "```"),
        tip: () => block(":::tip", ":::"),
        warning: () => block(":::warning", ":::"),
      };
      sr.querySelectorAll(".tb").forEach((btn) => {
        btn.addEventListener("click", () => {
          const fn = actions[btn.dataset.act];
          if (fn) {
            fn();
            schedulePreview();
          }
        });
      });

      ta.addEventListener("input", schedulePreview);
      ta.addEventListener("keydown", (e) => {
        const mod = e.metaKey || e.ctrlKey;
        if (mod && e.key.toLowerCase() === "b") { e.preventDefault(); actions.bold(); schedulePreview(); }
        else if (mod && e.key.toLowerCase() === "i") { e.preventDefault(); actions.italic(); schedulePreview(); }
        else if (e.key === "Tab") {
          e.preventDefault();
          const s = ta.selectionStart, en = ta.selectionEnd;
          ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(en);
          ta.selectionStart = ta.selectionEnd = s + 2;
        }
      });

      // Initial render.
      renderPreview();
    },
  })
);
