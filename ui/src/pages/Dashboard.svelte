<script>
  import { onMount } from 'svelte';
  import { getSuccessRateColor } from '../lib/utils.js';
  import { dateRange, dashboardStore } from '../lib/stores.js';
  import StatusBadge from '../components/StatusBadge.svelte';

  let currentDateRange = {};
  let dashboardData = { data: null, loading: true, error: null };

  // Subscribe to date range changes
  dateRange.subscribe(range => {
    currentDateRange = range;
    // Fetch new data when date range changes
    dashboardStore.fetchSummary(range.from, range.to);
  });

  // Subscribe to dashboard store
  dashboardStore.subscribe(value => {
    dashboardData = value;
  });

  // Fetch initial data
  onMount(() => {
    dashboardStore.fetchSummary(currentDateRange.from, currentDateRange.to);
  });

  // Extract data with defaults
  $: summary = dashboardData.data || {};
  $: totalJobs = summary.total_jobs || 0;
  $: activeJobs = summary.active_jobs || 0;
  $: totalExecutions = summary.total_executions || 0;
  $: overallSuccessRate = summary.overall_success_rate || 0;
  $: jobSummaries = summary.job_summaries || [];
  $: recentActivity = summary.recent_activity || [];

  // Map job summaries to display format
  $: jobs = jobSummaries.map(job => ({
    id: job.id,
    name: job.name,
    type: job.type,
    successRate: job.success_rate,
    lastExecution: job.last_execution,
    avgResponseTime: job.avg_response_time,
    executionCount: job.execution_count,
    enabled: job.enabled,
    // Determine status based on success rate and last execution
    status: !job.enabled ? 'disabled' : 
            job.success_rate >= 95 ? 'success' : 
            job.success_rate >= 80 ? 'warning' : 'failed'
  }));
</script>

<div class="dashboard">
  {#if dashboardData.loading}
    <div class="loading">
      <p>Loading dashboard data...</p>
    </div>
  {:else if dashboardData.error}
    <div class="error">
      <p>Error loading dashboard: {dashboardData.error}</p>
      <button on:click={() => dashboardStore.fetchSummary(currentDateRange.from, currentDateRange.to)}>
        Try Again
      </button>
    </div>
  {:else}
    <div class="stats-grid">
      <div class="stat-card">
        <h3>Total Jobs</h3>
        <div class="stat-number">{totalJobs}</div>
      </div>
      <div class="stat-card">
        <h3>Active Jobs</h3>
        <div class="stat-number">{activeJobs}</div>
      </div>
      <div class="stat-card">
        <h3>Total Executions</h3>
        <div class="stat-number">{totalExecutions.toLocaleString()}</div>
      </div>
      <div class="stat-card">
        <h3>Success Rate</h3>
        <div class="stat-number" style="color: {getSuccessRateColor(overallSuccessRate)}">{overallSuccessRate.toFixed(1)}%</div>
      </div>
    </div>

    <div class="jobs-overview">
      <h2>Recent Jobs</h2>
      {#if jobs.length === 0}
        <p>No jobs found for the selected date range.</p>
      {:else}
        <div class="jobs-list">
          {#each jobs as job}
            <div class="job-card">
              <div class="job-header">
                <h3>{job.name}</h3>
                <StatusBadge status={job.status} size="small" />
              </div>
              <div class="job-details">
                <p><strong>Type:</strong> {job.type}</p>
                <p><strong>Success Rate:</strong> <span style="color: {getSuccessRateColor(job.successRate)}">{job.successRate.toFixed(1)}%</span></p>
                <p><strong>Executions:</strong> {job.executionCount.toLocaleString()}</p>
                {#if job.lastExecution}
                  <p><strong>Last Run:</strong> {new Date(job.lastExecution).toLocaleString()}</p>
                {/if}
                {#if job.avgResponseTime}
                  <p><strong>Avg Response:</strong> {Math.round(job.avgResponseTime)}ms</p>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    {#if recentActivity.length > 0}
      <div class="recent-activity">
        <h2>Recent Activity</h2>
        <div class="activity-list">
          {#each recentActivity.slice(0, 10) as execution}
            <div class="activity-item">
              <div class="activity-status">
                <StatusBadge status={execution.status} size="small" />
              </div>
              <div class="activity-details">
                <p><strong>{execution.job?.name || 'Unknown Job'}</strong></p>
                <p>{new Date(execution.timestamp).toLocaleString()}</p>
                {#if execution.response_time}
                  <p>{execution.response_time}ms</p>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .dashboard {
    padding: var(--spacing-xl);
    max-width: 1200px;
    margin: 0 auto;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-xl);
  }

  .stat-card {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    padding: var(--spacing-lg);
    box-shadow: var(--shadow-lg);
    text-align: center;
  }

  .stat-card h3 {
    margin: 0 0 var(--spacing-sm) 0;
    color: var(--text-primary);
    font-size: var(--font-md);
    text-transform: uppercase;
    font-weight: 600;
  }

  .stat-number {
    font-size: var(--font-xxl);
    font-weight: bold;
    color: var(--text-primary);
  }

  .job-card {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    padding: var(--spacing-lg);
    box-shadow: var(--shadow-lg);
  }

  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-md);
  }

  .job-header h3 {
    margin: 0;
    color: var(--text-primary);
    font-size: var(--font-lg);
  }

  .job-details p {
    margin: 0.25rem 0;
    color: var(--text-secondary);
    font-size: var(--font-md);
  }

  .loading {
    text-align: center;
    padding: var(--spacing-xl);
    color: var(--text-secondary);
  }

  .error {
    text-align: center;
    padding: var(--spacing-xl);
    color: var(--status-failed);
  }

  .error button {
    background: var(--primary-color);
    color: white;
    border: none;
    padding: var(--spacing-sm) var(--spacing-md);
    border-radius: var(--radius-sm);
    cursor: pointer;
    margin-top: var(--spacing-md);
  }

  .error button:hover {
    background: var(--primary-color-light);
  }

  .recent-activity {
    margin-top: var(--spacing-xl);
  }

  .recent-activity h2 {
    margin-bottom: var(--spacing-md);
  }

  .activity-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .activity-item {
    display: flex;
    align-items: center;
    background: var(--card-bg);
    border-radius: var(--radius-md);
    padding: var(--spacing-md);
    box-shadow: var(--shadow-lg);
  }

  .activity-status {
    margin-right: var(--spacing-md);
  }

  .activity-details p {
    margin: 0.125rem 0;
    font-size: var(--font-sm);
  }

  .activity-details p:first-child {
    font-weight: 600;
    color: var(--text-primary);
  }

  .activity-details p:not(:first-child) {
    color: var(--text-secondary);
  }
</style>
