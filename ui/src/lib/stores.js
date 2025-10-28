import { writable, derived } from "svelte/store";
import { apiService } from "./api.js";
import { websocketService, MessageType } from "./websocket.js";

// Default to last 7 days
function getDefaultDateRange() {
  const to = new Date();
  const from = new Date();
  from.setDate(from.getDate() - 7);

  return {
    from: from.toISOString(), // ISO 8601 format with time
    to: to.toISOString(),
  };
}

// Parse URL query parameters for date range
function getDateRangeFromURL() {
  if (typeof window === "undefined") return getDefaultDateRange();

  const params = new URLSearchParams(window.location.search);
  const fromParam = params.get("from");
  const toParam = params.get("to");

  // Validate ISO 8601 format
  const iso8601Regex = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/;

  if (
    fromParam &&
    toParam &&
    iso8601Regex.test(fromParam) &&
    iso8601Regex.test(toParam)
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
  const { subscribe, set, update } = writable({
    ...initialRange,
    isToDateLive: true, // Track if "to" should auto-update to "now"
  });

  return {
    subscribe,
    setRange: (from, to, isToDateLive = true) => {
      const range = { from, to, isToDateLive };
      set(range);
      updateURL(from, to);
    },
    setFrom: (from) =>
      update((range) => {
        const newRange = { ...range, from };
        updateURL(newRange.from, newRange.to);
        return newRange;
      }),
    setTo: (to, isToDateLive = false) =>
      update((range) => {
        const newRange = { ...range, to, isToDateLive };
        updateURL(newRange.from, newRange.to);
        return newRange;
      }),
    // Update "to" to current time if it's in live mode
    updateToIfLive: () =>
      update((range) => {
        if (range.isToDateLive) {
          const newTo = new Date().toISOString();
          updateURL(range.from, newTo);
          return { ...range, to: newTo };
        }
        return range;
      }),
    reset: () => {
      const defaultRange = getDefaultDateRange();
      set({ ...defaultRange, isToDateLive: true });
      updateURL(defaultRange.from, defaultRange.to);
    },
    // Initialize from URL (called on app start)
    initFromURL: () => {
      const urlRange = getDateRangeFromURL();
      set({ ...urlRange, isToDateLive: true });
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

// Create API-driven stores

// Jobs store with loading and error states
function createJobsStore() {
  const { subscribe, set, update } = writable({
    data: [],
    loading: false,
    error: null,
  });

  return {
    subscribe,
    async fetchJobs(from, to) {
      update((state) => ({ ...state, loading: true, error: null }));

      try {
        const jobs = await apiService.getJobs(from, to);
        set({ data: jobs, loading: false, error: null });
      } catch (error) {
        set({ data: [], loading: false, error: error.message });
      }
    },
    // Update a specific job in the store
    updateJob: (updatedJob) => {
      update((state) => ({
        ...state,
        data: state.data.map((job) =>
          job.id === updatedJob.id ? { ...job, ...updatedJob } : job
        ),
      }));
    },
    // Add a new execution to a job's recent_executions
    addExecution: (execution) => {
      update((state) => ({
        ...state,
        data: state.data.map((job) => {
          if (job.id === execution.job_id) {
            // Add to beginning of recent_executions array
            const recentExecutions = [
              execution,
              ...(job.recent_executions || []),
            ].slice(0, 10); // Keep only last 10

            // Recalculate success rate
            const successCount = recentExecutions.filter(
              (e) => e.status === "success"
            ).length;
            const successRate =
              recentExecutions.length > 0
                ? (successCount / recentExecutions.length) * 100
                : 0;

            return {
              ...job,
              recent_executions: recentExecutions,
              success_rate: successRate,
              last_execution: execution.timestamp,
              avg_response_time: execution.response_time, // Simplified - would need full recalc
            };
          }
          return job;
        }),
      }));
    },
    reset: () => set({ data: [], loading: false, error: null }),
  };
}

// Dashboard store
function createDashboardStore() {
  const { subscribe, set, update } = writable({
    data: null,
    loading: false,
    error: null,
  });

  return {
    subscribe,
    async fetchSummary(from, to) {
      update((state) => ({ ...state, loading: true, error: null }));

      try {
        const summary = await apiService.getDashboardSummary(from, to);
        set({ data: summary, loading: false, error: null });
      } catch (error) {
        set({ data: null, loading: false, error: error.message });
      }
    },
    // Update dashboard with new data
    updateSummary: (summary) => {
      update((state) => ({
        ...state,
        data: { ...state.data, ...summary },
      }));
    },
    // Add a new execution to a job in the dashboard's job_summaries
    addExecutionToJob: (execution) => {
      update((state) => {
        if (!state.data || !state.data.job_summaries) {
          return state;
        }

        const updatedJobSummaries = state.data.job_summaries.map((job) => {
          if (job.id === execution.job_id) {
            // Add to beginning of recent_executions array
            const recentExecutions = [
              execution,
              ...(job.recent_executions || []),
            ].slice(0, 10); // Keep only last 10

            // Recalculate success rate based on recent executions
            const successCount = recentExecutions.filter(
              (e) => e.status === "success"
            ).length;
            const successRate =
              recentExecutions.length > 0
                ? (successCount / recentExecutions.length) * 100
                : 0;

            return {
              ...job,
              recent_executions: recentExecutions,
              success_rate: successRate,
              last_execution: execution.timestamp,
              avg_response_time: execution.response_time,
            };
          }
          return job;
        });

        return {
          ...state,
          data: {
            ...state.data,
            job_summaries: updatedJobSummaries,
          },
        };
      });
    },
    reset: () => set({ data: null, loading: false, error: null }),
  };
}

// Individual job store for job details page
function createJobDetailStore() {
  const { subscribe, set, update } = writable({
    data: null,
    loading: false,
    error: null,
  });

  return {
    subscribe,
    async fetchJob(id, from, to, limit = 100) {
      update((state) => ({ ...state, loading: true, error: null }));

      try {
        const job = await apiService.getJob(id, from, to, limit);
        set({ data: job, loading: false, error: null });
      } catch (error) {
        set({ data: null, loading: false, error: error.message });
      }
    },
    // Add a new execution to the current job's executions
    addExecution: (execution) => {
      update((state) => {
        if (!state.data || state.data.id !== execution.job_id) {
          return state;
        }

        const executions = [execution, ...(state.data.executions || [])];

        return {
          ...state,
          data: {
            ...state.data,
            executions,
          },
        };
      });
    },
    reset: () => set({ data: null, loading: false, error: null }),
  };
}

// Export store instances
export const jobsStore = createJobsStore();
export const dashboardStore = createDashboardStore();
export const jobDetailStore = createJobDetailStore();

// Derived store that automatically fetches data when date range changes
export const autoJobsStore = derived([dateRange, jobsStore], ([range], set) => {
  // Fetch jobs when date range changes
  jobsStore.fetchJobs(range.from, range.to);
});

export const autoDashboardStore = derived(
  [dateRange, dashboardStore],
  ([range], set) => {
    // Fetch dashboard data when date range changes
    dashboardStore.fetchSummary(range.from, range.to);
  }
);

// Initialize WebSocket and set up message handlers
export function initializeWebSocket() {
  // Connect to WebSocket server
  const wsUrl = `ws://${window.location.hostname}:8080/ws`;
  websocketService.connect(wsUrl);

  // Helper to get current date range value
  let currentRangeValue;
  const unsubscribeDateRange = dateRange.subscribe((value) => {
    currentRangeValue = value;
  });

  const getCurrentDateRange = () => currentRangeValue;

  // Handle execution_created messages
  websocketService.on(MessageType.EXECUTION_CREATED, (execution) => {
    const currentRange = getCurrentDateRange();
    const executionTime = new Date(execution.timestamp);
    const fromTime = new Date(currentRange.from);

    // Check if execution is within the "from" date
    if (executionTime < fromTime) {
      return;
    }

    // If "to" date is locked, also check the upper bound
    if (!currentRange.isToDateLive) {
      const toTime = new Date(currentRange.to);
      if (executionTime > toTime) {
        return;
      }
    }
    // If in live mode, we don't check upper bound (always accept newer executions)

    // Update dashboard store (for Dashboard page)
    dashboardStore.addExecutionToJob(execution);

    // Update jobs store (for Jobs page)
    jobsStore.addExecution(execution);

    // Update job detail store if viewing this job
    jobDetailStore.addExecution(execution);
  });

  // Handle job_updated messages
  websocketService.on(MessageType.JOB_UPDATED, (job) => {
    // Update the job in the jobs store
    jobsStore.updateJob(job);
  });

  // Handle dashboard_updated messages
  websocketService.on(MessageType.DASHBOARD_UPDATED, (summary) => {
    // Update dashboard store with new summary
    dashboardStore.updateSummary(summary);
  });

  // Return cleanup function
  return () => {
    unsubscribeDateRange();
    websocketService.disconnect();
  };
}

// Export WebSocket service and connection status for UI components
export { websocketService };
