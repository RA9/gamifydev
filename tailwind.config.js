/** @type {import('tailwindcss').Config} */
module.exports = {
  // Scan the HTML entry and every JS file, since most markup (and therefore
  // most utility classes) lives inside JS template strings.
  content: ["./index.html", "./index.js", "./js/**/*.js"],
  // Classes assembled dynamically at runtime (e.g. `text-${color}-500`) can't be
  // discovered by the content scanner, so list them explicitly.
  safelist: [
    "text-green-500",
    "text-yellow-500",
    "text-red-500",
    "text-green-400",
    "text-red-400",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
