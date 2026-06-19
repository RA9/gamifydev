/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
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
    "text-grass-600",
    "text-amber-500",
    "text-rose-500",
    "from-brand-500",
    "to-brand-700",
  ],
  theme: {
    extend: {
      fontFamily: {
        // Friendly, rounded display/body face with a robust system fallback so
        // the app still looks intentional offline (PWA / low-bandwidth).
        sans: [
          "Nunito",
          "ui-rounded",
          "Quicksand",
          '"Segoe UI"',
          "system-ui",
          "-apple-system",
          "sans-serif",
        ],
      },
      colors: {
        // Primary brand: a friendly indigo/violet evolved from the old #757195.
        brand: {
          50: "#f2f1ff",
          100: "#e7e5ff",
          200: "#d2cdff",
          300: "#b4aaff",
          400: "#937dff",
          500: "#7857f7",
          600: "#6740e8",
          700: "#5731c8",
          800: "#482aa2",
          900: "#3d2880",
        },
        // Success / gamification green (Duolingo-ish "grass").
        grass: {
          50: "#f1fcebff",
          100: "#dff7d2",
          200: "#c0ee9f",
          300: "#9ce26a",
          400: "#7ad23e",
          500: "#58cc02",
          600: "#46a302",
          700: "#377e06",
          800: "#2d630c",
          900: "#27530f",
        },
      },
      boxShadow: {
        // Soft, diffuse card shadow for the friendly look.
        soft: "0 10px 30px -12px rgba(40, 28, 90, 0.18)",
        card: "0 4px 16px -6px rgba(40, 28, 90, 0.12)",
        pop: "0 12px 0 0 rgba(0,0,0,0.06)",
      },
      borderRadius: {
        "2xl": "1.25rem",
        "3xl": "1.75rem",
      },
      keyframes: {
        "pop-in": {
          "0%": { transform: "scale(0.85)", opacity: "0" },
          "70%": { transform: "scale(1.04)", opacity: "1" },
          "100%": { transform: "scale(1)" },
        },
        "fade-up": {
          "0%": { transform: "translateY(12px)", opacity: "0" },
          "100%": { transform: "translateY(0)", opacity: "1" },
        },
        float: {
          "0%,100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-8px)" },
        },
        shimmer: {
          "0%": { backgroundPosition: "-200% 0" },
          "100%": { backgroundPosition: "200% 0" },
        },
        wiggle: {
          "0%,100%": { transform: "rotate(-6deg)" },
          "50%": { transform: "rotate(6deg)" },
        },
        "bar-fill": {
          "0%": { width: "0%" },
        },
      },
      animation: {
        "pop-in": "pop-in 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both",
        "fade-up": "fade-up 0.5s ease-out both",
        float: "float 4s ease-in-out infinite",
        shimmer: "shimmer 2.5s linear infinite",
        wiggle: "wiggle 0.5s ease-in-out",
        "bar-fill": "bar-fill 1s ease-out",
      },
    },
  },
  plugins: [],
};
