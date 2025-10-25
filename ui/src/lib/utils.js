// Shared utility functions for the Moogie UI

/**
 * Get CSS variable value from root
 * @param {string} variable - CSS variable name (without --)
 * @returns {string} CSS variable reference
 */
function getCSSVar(variable) {
  return `var(--${variable})`;
}

/**
 * Get color for job status
 * @param {string} status - The status string
 * @returns {string} CSS variable reference or hex color
 */
export function getStatusColor(status) {
  const statusLower = status.toLowerCase();
  switch (statusLower) {
    case "running":
      return getCSSVar("status-running");
    case "success":
      return getCSSVar("status-success");
    case "failed":
    case "failure":
    case "error":
      return getCSSVar("status-failed");
    case "warning":
      return getCSSVar("status-warning");
    case "timeout":
      return getCSSVar("status-timeout");
    case "disabled":
      return getCSSVar("status-default");
    default:
      return getCSSVar("status-default");
  }
}

/**
 * Get color based on success rate percentage
 * @param {number} rate - Success rate percentage
 * @returns {string} CSS variable reference
 */
export function getSuccessRateColor(rate) {
  if (rate >= 100) return getCSSVar("success-perfect");
  if (rate >= 95) return getCSSVar("success-good");
  if (rate >= 90) return getCSSVar("success-fair");
  return getCSSVar("success-poor");
}

/**
 * Get color for check type
 * @param {string} kind - The check type/kind
 * @returns {string} CSS variable reference or hex color
 */
export function getCheckTypeColor(kind) {
  switch (kind) {
    case "HttpCheck":
    case "http":
      return getCSSVar("check-http");
    case "TcpCheck":
    case "tcp":
      return getCSSVar("check-tcp");
    case "DnsCheck":
    case "dns":
      return getCSSVar("check-dns");
    case "SslCheck":
    case "ssl":
      return getCSSVar("check-ssl");
    case "PingCheck":
    case "ping":
      return getCSSVar("check-ping");
    default:
      return getCSSVar("status-default");
  }
}

/**
 * Format date string to locale string
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted date string
 */
export function formatDate(dateString) {
  return new Date(dateString).toLocaleString();
}

/**
 * Format duration in milliseconds
 * @param {number} ms - Duration in milliseconds
 * @returns {string} Formatted duration string
 */
export function formatDuration(ms) {
  return `${ms}ms`;
}

/**
 * Format labels object to comma-separated string
 * @param {Object} labels - Key-value pairs of labels
 * @returns {string} Formatted labels string
 */
export function formatLabels(labels) {
  return Object.entries(labels)
    .map(([key, value]) => `${key}: ${value}`)
    .join(", ");
}

/**
 * Calculate overall success rate from array of jobs
 * @param {Array} jobs - Array of job objects with successRate property
 * @returns {number} Average success rate rounded to 1 decimal place
 */
export function calculateOverallSuccessRate(jobs) {
  if (jobs.length === 0) return 0;
  const totalRate = jobs.reduce((sum, job) => sum + job.successRate, 0);
  return Math.round((totalRate / jobs.length) * 10) / 10;
}
