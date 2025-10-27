<script>
  import { formatDate, getSuccessRateColor } from '../../lib/utils.js';
  import StatusBadge from '../StatusBadge.svelte';
  
  export let job;
  export let executionCount;
  export let successRate;
</script>

<div class="job-summary">
  <div class="summary-card">
    <h3>Status</h3>
    <StatusBadge status={job.enabled ? 'success' : 'disabled'} size="large" />
    <div class="summary-text">{job.enabled ? 'Enabled' : 'Disabled'}</div>
  </div>
  <div class="summary-card">
    <h3>Total Executions</h3>
    <div class="summary-number">{executionCount}</div>
  </div>
  <div class="summary-card">
    <h3>Success Rate</h3>
    <div class="summary-number" style="color: {getSuccessRateColor(successRate)}">{successRate.toFixed(1)}%</div>
  </div>
  <div class="summary-card">
    <h3>Last Run</h3>
    <div class="summary-text">{job.last_execution ? formatDate(job.last_execution) : 'Never'}</div>
  </div>
</div>

<style>
  .job-summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--spacing-lg);
    margin: var(--spacing-xl) 0;
  }

  .summary-card {
    background: var(--card-bg);
    padding: var(--spacing-lg);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
    text-align: center;
  }

  .summary-card h3 {
    margin: 0 0 var(--spacing-md) 0;
    color: var(--text-secondary);
    font-size: var(--font-sm);
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .summary-number {
    font-size: var(--font-xxl);
    font-weight: bold;
    color: var(--primary-color);
  }

  .summary-text {
    font-size: var(--font-md);
    color: var(--text-primary);
    margin-top: var(--spacing-sm);
  }

  @media (max-width: 768px) {
    .job-summary {
      grid-template-columns: 1fr;
    }
  }
</style>
