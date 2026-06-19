// Dark-mode toggle. The initial theme is applied by a tiny inline script in
// <head> (before paint, to avoid a flash); this file just handles toggling and
// keeping the button icons in sync. Theme is persisted in localStorage so it
// can be read synchronously on load.
const THEME_KEY = "gamifydev-theme";

function isDarkMode() {
  return document.documentElement.classList.contains("dark");
}

function updateThemeToggleIcons() {
  // Show a sun in dark mode (tap to go light) and a moon in light mode.
  const svg =
    typeof icon === "function"
      ? icon(isDarkMode() ? "sun" : "moon", "w-5 h-5")
      : "";
  document
    .querySelectorAll(".theme-toggle-icon")
    .forEach((el) => (el.innerHTML = svg));
}

function toggleTheme() {
  const dark = document.documentElement.classList.toggle("dark");
  try {
    localStorage.setItem(THEME_KEY, dark ? "dark" : "light");
  } catch (e) {
    console.log("theme persist failed:", e);
  }
  updateThemeToggleIcons();
}

document.addEventListener("DOMContentLoaded", updateThemeToggleIcons);

// --- Mobile navigation menu -------------------------------------------------
function toggleMobileMenu() {
  const menu = document.getElementById("gd-mobile-menu");
  const btn = document.getElementById("gd-menu-btn");
  if (!menu) return;
  const open = menu.classList.toggle("hidden") === false;
  if (btn) btn.setAttribute("aria-expanded", String(open));
}

function closeMobileMenu() {
  const menu = document.getElementById("gd-mobile-menu");
  const btn = document.getElementById("gd-menu-btn");
  if (menu && !menu.classList.contains("hidden")) {
    menu.classList.add("hidden");
    if (btn) btn.setAttribute("aria-expanded", "false");
  }
}
