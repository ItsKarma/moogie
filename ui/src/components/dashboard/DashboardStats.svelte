<script>
  import { getSuccessRateColor } from '../../lib/utils.js';

  export let totalJobs = 0;
  export let activeJobs = 0;
  export let totalExecutions = 0;
  export let overallSuccessRate = 0;
</script>

<div class="stats-bar">
  <div class="stat-item">
    <span class="stat-label">Total Jobs</span>
    <span class="stat-value">{totalJobs}</span>
  </div>
  <div class="stat-divider"></div>
  <div class="stat-item">
    <span class="stat-label">Active</span>
    <span class="stat-value">{activeJobs}</span>
  </div>
  <div class="stat-divider"></div>
  <div class="stat-item">
    <span class="stat-label">Executions</span>
    <span class="stat-value">{totalExecutions.toLocaleString()}</span>
  </div>
  <div class="stat-divider"></div>
  <div class="stat-item">
    <span class="stat-label">Success Rate</span>
    <span class="stat-value" style="color: {getSuccessRateColor(overallSuccessRate)}">{overallSuccessRate.toFixed(1)}%</span>
  </div>
</div>

<style>
  .stats-bar {
    display: flex;
    align-items: center;
    justify-content: space-around;
    background: var(--card-bg);
    padding: var(--spacing-lg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-md);
    margin-bottom: var(--spacing-xl);
  }

  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--spacing-xs);
  }

  .stat-label {
    font-size: var(--font-sm);
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-weight: 500;
  }

  .stat-value {
    font-size: var(--font-xl);
    font-weight: 700;
    color: var(--text-primary);
  }

  .stat-divider {
    width: 1px;
    height: 40px;
    background: var(--border-color);
  }

  @media (max-width: 768px) {
    .stats-bar {
      flex-wrap: wrap;
      gap: var(--spacing-md);
    }

    .stat-divider {
      display: none;
    }
  }
</style>
