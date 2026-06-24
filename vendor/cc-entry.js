// Bundled (by esbuild) into js/external/jscpp.bundle.js. JSCPP is a C/C++
// interpreter in pure JS — no WASM toolchain — so beginner C challenges run
// fully offline. Attach to globalThis so it works in both window and Web Worker
// contexts (the C runner imports this inside a worker).
import JSCPP from "JSCPP";
globalThis.JSCPP = JSCPP && JSCPP.run ? JSCPP : (JSCPP && JSCPP.default) || JSCPP;
