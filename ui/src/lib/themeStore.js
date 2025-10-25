import { writable } from "svelte/store";

const THEME_KEY = "moogie-theme";
const THEMES = {
  LIGHT: "light",
  DARK: "dark",
  SYSTEM: "system",
};

/**
 * Get the system theme preference
 * @returns {string} 'light' or 'dark'
 */
function getSystemTheme() {
  if (typeof window === "undefined") return THEMES.DARK;
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? THEMES.DARK
    : THEMES.LIGHT;
}

/**
 * Get the stored theme preference from localStorage
 * @returns {string} 'light', 'dark', or 'system'
 */
function getStoredTheme() {
  if (typeof window === "undefined") return THEMES.SYSTEM;
  const stored = localStorage.getItem(THEME_KEY);
  return stored && Object.values(THEMES).includes(stored)
    ? stored
    : THEMES.SYSTEM;
}

/**
 * Calculate the effective theme (resolving 'system' to actual theme)
 * @param {string} theme - 'light', 'dark', or 'system'
 * @returns {string} 'light' or 'dark'
 */
function resolveTheme(theme) {
  return theme === THEMES.SYSTEM ? getSystemTheme() : theme;
}

/**
 * Apply the theme to the document
 * @param {string} effectiveTheme - 'light' or 'dark'
 */
function applyTheme(effectiveTheme) {
  if (typeof document === "undefined") return;

  const root = document.documentElement;

  // Remove both classes first
  root.classList.remove("theme-light", "theme-dark");

  // Add the appropriate class
  root.classList.add(`theme-${effectiveTheme}`);

  // Also set data attribute for CSS selectors
  root.setAttribute("data-theme", effectiveTheme);
}

/**
 * Create the theme store
 */
function createThemeStore() {
  const initialTheme = getStoredTheme();
  const initialEffectiveTheme = resolveTheme(initialTheme);

  const { subscribe, set } = writable({
    theme: initialTheme,
    effectiveTheme: initialEffectiveTheme,
  });

  // Theme is already applied by inline script in index.html
  // No need to apply it again here to avoid forced layout warning

  // Listen for system theme changes
  if (typeof window !== "undefined") {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    mediaQuery.addEventListener("change", (e) => {
      // Only update if current theme is 'system'
      const currentTheme = getStoredTheme();
      if (currentTheme === THEMES.SYSTEM) {
        const newEffectiveTheme = e.matches ? THEMES.DARK : THEMES.LIGHT;
        applyTheme(newEffectiveTheme);
        set({
          theme: THEMES.SYSTEM,
          effectiveTheme: newEffectiveTheme,
        });
      }
    });
  }

  return {
    subscribe,
    /**
     * Set the theme preference
     * @param {string} newTheme - 'light', 'dark', or 'system'
     */
    setTheme: (newTheme) => {
      if (!Object.values(THEMES).includes(newTheme)) {
        console.error(`Invalid theme: ${newTheme}`);
        return;
      }

      // Store in localStorage
      localStorage.setItem(THEME_KEY, newTheme);

      // Resolve and apply
      const effectiveTheme = resolveTheme(newTheme);
      applyTheme(effectiveTheme);

      // Update store
      set({
        theme: newTheme,
        effectiveTheme,
      });
    },
    /**
     * Toggle between light and dark (ignores system)
     */
    toggle: () => {
      const currentTheme = getStoredTheme();
      const currentEffective = resolveTheme(currentTheme);
      const newTheme =
        currentEffective === THEMES.LIGHT ? THEMES.DARK : THEMES.LIGHT;

      localStorage.setItem(THEME_KEY, newTheme);
      applyTheme(newTheme);

      set({
        theme: newTheme,
        effectiveTheme: newTheme,
      });
    },
  };
}

export const themeStore = createThemeStore();
export const THEME_OPTIONS = THEMES;
