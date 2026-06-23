// Entry point bundled (by esbuild) into js/external/codemirror.bundle.js.
// Exposes a single global `window.CM` so the no-bundler app can use CodeMirror 6
// offline, the same way dexie.js is vendored.
import { EditorView, basicSetup } from "codemirror";
import { EditorState, Compartment } from "@codemirror/state";
import { keymap, lineNumbers } from "@codemirror/view";
import { indentWithTab, defaultKeymap } from "@codemirror/commands";
import { javascript } from "@codemirror/lang-javascript";
import { python } from "@codemirror/lang-python";
import { cpp } from "@codemirror/lang-cpp";
import { java } from "@codemirror/lang-java";

window.CM = {
  EditorView,
  EditorState,
  Compartment,
  basicSetup,
  keymap,
  lineNumbers,
  indentWithTab,
  defaultKeymap,
  langs: { javascript, python, cpp, java },
};
