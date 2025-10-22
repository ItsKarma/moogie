import { writable } from "svelte/store";

// Default to last 7 days
function getDefaultDateRange() {
  const to = new Date();
  const from = new Date();
  from.setDate(from.getDate() - 7);

  return {
    from: from.toISOString().split("T")[0], // YYYY-MM-DD format
    to: to.toISOString().split("T")[0],
  };
}

// Parse URL query parameters for date range
function getDateRangeFromURL() {
  if (typeof window === "undefined") return getDefaultDateRange();

  const params = new URLSearchParams(window.location.search);
  const fromParam = params.get("from");
  const toParam = params.get("to");

  // Validate date format (YYYY-MM-DD)
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;

  if (
    fromParam &&
    toParam &&
    dateRegex.test(fromParam) &&
    dateRegex.test(toParam)
  ) {
    // Ensure from date is not after to date
    const fromDate = new Date(fromParam);
    const toDate = new Date(toParam);

    if (fromDate <= toDate && toDate <= new Date()) {
      return { from: fromParam, to: toParam };
    }
  }

  return getDefaultDateRange();
}

// Update URL with date range parameters
function updateURL(from, to) {
  if (typeof window === "undefined") return;

  const url = new URL(window.location);
  url.searchParams.set("from", from);
  url.searchParams.set("to", to);

  window.history.replaceState({}, "", url.toString());
}

// Create the global date range store
function createDateRangeStore() {
  const initialRange = getDateRangeFromURL();
  const { subscribe, set, update } = writable(initialRange);

  return {
    subscribe,
    setRange: (from, to) => {
      const range = { from, to };
      set(range);
      updateURL(from, to);
    },
    setFrom: (from) =>
      update((range) => {
        const newRange = { ...range, from };
        updateURL(newRange.from, newRange.to);
        return newRange;
      }),
    setTo: (to) =>
      update((range) => {
        const newRange = { ...range, to };
        updateURL(newRange.from, newRange.to);
        return newRange;
      }),
    reset: () => {
      const defaultRange = getDefaultDateRange();
      set(defaultRange);
      updateURL(defaultRange.from, defaultRange.to);
    },
    // Initialize from URL (called on app start)
    initFromURL: () => {
      const urlRange = getDateRangeFromURL();
      set(urlRange);
    },
  };
}

export const dateRange = createDateRangeStore();

// Helper function to check if a date falls within the selected range
export function isDateInRange(dateString, range) {
  const date = new Date(dateString);
  const fromDate = new Date(range.from);
  const toDate = new Date(range.to);

  // Set time to start/end of day for accurate comparison
  fromDate.setHours(0, 0, 0, 0);
  toDate.setHours(23, 59, 59, 999);

  return date >= fromDate && date <= toDate;
}

// Helper function to filter execution history by date range
export function filterExecutionsByDateRange(executions, range) {
  return executions.filter((execution) =>
    isDateInRange(execution.timestamp, range)
  );
}

// Helper function to calculate success rate for filtered data
export function calculateSuccessRateForRange(executions, range) {
  const filteredExecutions = filterExecutionsByDateRange(executions, range);
  if (filteredExecutions.length === 0) return 0;

  const successfulExecutions = filteredExecutions.filter(
    (execution) => execution.status === "success"
  ).length;

  return (
    Math.round((successfulExecutions / filteredExecutions.length) * 100 * 10) /
    10
  );
}
