// Responsive code editor wrapper.
//
// Prefers CodeMirror 6 (vendored offline as js/external/codemirror.bundle.js,
// exposed on window.CM). If that bundle isn't present for any reason, it falls
// back to a styled <textarea> with tab handling — so coding sessions always
// work, even offline and even before the bundle is built. Both paths expose the
// same tiny interface: { getValue, setValue, focus, destroy }.

function gdLangExtension(language) {
  const CM = window.CM;
  if (!CM) return null;
  const l = (language || "javascript").toLowerCase();
  if (l === "python" || l === "py") return CM.langs.python();
  if (l === "c" || l === "cpp" || l === "c++") return CM.langs.cpp();
  if (l === "java") return CM.langs.java();
  return CM.langs.javascript();
}

// Create an editor inside `parent`. opts: { doc, language, onChange }.
function createCodeEditor(parent, opts = {}) {
  const doc = opts.doc != null ? String(opts.doc) : "";
  const CM = window.CM;

  if (CM && CM.EditorView) {
    const exts = [
      CM.basicSetup,
      CM.keymap.of([CM.indentWithTab]),
      gdLangExtension(opts.language),
      CM.EditorView.theme({
        "&": { fontSize: "14px", borderRadius: "16px", overflow: "hidden", backgroundColor: "#ffffff" },
        ".cm-scroller": {
          fontFamily: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace",
          maxHeight: "420px",
        },
        ".cm-content": { padding: "12px 0" },
        "&.cm-focused": { outline: "none" },
        ".cm-gutters": { backgroundColor: "#f7f7fd", border: "none", color: "#94a3b8" },
      }),
    ].filter(Boolean);

    if (typeof opts.onChange === "function") {
      exts.push(
        CM.EditorView.updateListener.of((v) => {
          if (v.docChanged) opts.onChange(view.state.doc.toString());
        })
      );
    }

    const view = new CM.EditorView({
      state: CM.EditorState.create({ doc, extensions: exts }),
      parent,
    });
    return {
      cm: true,
      getValue: () => view.state.doc.toString(),
      setValue: (text) =>
        view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: String(text) } }),
      focus: () => view.focus(),
      destroy: () => view.destroy(),
    };
  }

  // --- Fallback: textarea with tab support ---------------------------------
  const ta = document.createElement("textarea");
  ta.className = "gd-code-fallback";
  ta.value = doc;
  ta.spellcheck = false;
  ta.setAttribute("autocomplete", "off");
  ta.setAttribute("autocapitalize", "off");
  ta.setAttribute("autocorrect", "off");
  ta.addEventListener("keydown", (e) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const s = ta.selectionStart;
      const eN = ta.selectionEnd;
      ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(eN);
      ta.selectionStart = ta.selectionEnd = s + 2;
    }
  });
  if (typeof opts.onChange === "function") {
    ta.addEventListener("input", () => opts.onChange(ta.value));
  }
  parent.appendChild(ta);
  return {
    cm: false,
    getValue: () => ta.value,
    setValue: (text) => {
      ta.value = String(text);
    },
    focus: () => ta.focus(),
    destroy: () => ta.remove(),
  };
}
