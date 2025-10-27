<script>
  import { onMount } from 'svelte';
  import { push, link } from 'svelte-spa-router';
  import { formatDate } from '../lib/utils.js';
  import { dateRange, jobsStore, jobDetailStore } from '../lib/stores.js';
  import StatusBadge from '../components/StatusBadge.svelte';
  import DashboardStats from '../components/dashboard/DashboardStats.svelte';
  import JobDetail from './JobDetail.svelte';
  
  // Route parameters
  export let params = {};
  
  let currentDateRange = {};
  let jobsData = { data: [], loading: true, error: null };
  let jobDetailData = { data: null, loading: false, error: null };
  let selectedJob = null;

  // Subscribe to stores
  dateRange.subscribe(range => {
    currentDateRange = range;
    // Fetch jobs when date range changes
    jobsStore.fetchJobs(range.from, range.to);
  });

  jobsStore.subscribe(value => {
    jobsData = value;
  });

  jobDetailStore.subscribe(value => {
    jobDetailData = value;
    selectedJob = value.data;
  });

  // Watch for route parameter changes
  $: {
    if (params && params.jobId) {
      const jobId = parseInt(params.jobId, 10);
      if (!isNaN(jobId)) {
        jobDetailStore.fetchJob(jobId, currentDateRange.from, currentDateRange.to, 100);
      } else {
        // Invalid job ID, redirect to jobs list
        push('/jobs');
      }
    } else {
      selectedJob = null;
      jobDetailStore.reset();
    }
  }

  // Fetch initial data
  onMount(() => {
    jobsStore.fetchJobs(currentDateRange.from, currentDateRange.to);
  });

  function selectJob(job) {
    push(`/job/${job.id}`);
  }

  // Map API job data to display format for backward compatibility
  $: jobs = jobsData.data.map(job => ({
    id: job.id,
    name: job.name,
    type: job.type,
    enabled: job.enabled,
    successRate: job.success_rate || 0,
    lastExecution: job.last_execution,
    avgResponseTime: job.avg_response_time || 0,
    executions: job.executions || [],
    // Determine status based on enabled state and success rate
    status: !job.enabled ? 'disabled' : 
            job.success_rate >= 95 ? 'success' : 
            job.success_rate >= 80 ? 'warning' : 'failed',
    config: {
      metadata: {
        displayName: job.name,
        id: job.id.toString()
      }
    }
  }));

  $: filteredJobs = jobs; // API already handles date filtering
  
  // Calculate stats for DashboardStats component
  $: totalJobs = jobs.length;
  $: activeJobs = jobs.filter(job => job.enabled).length;
  $: totalExecutions = jobs.reduce((sum, job) => sum + (job.executions?.length || 0), 0);
  $: overallSuccessRate = jobs.length > 0
    ? jobs.reduce((sum, job) => sum + job.successRate, 0) / jobs.length
    : 0;
</script>

<div class="logs-container">
  {#if !selectedJob}
    <!-- Job List View -->
    <DashboardStats 
      {totalJobs}
      {activeJobs}
      {totalExecutions}
      {overallSuccessRate}
    />
    
    <div class="logs-header">
      <p>Click on any job to view its configuration and execution history</p>
    </div>

    {#if jobsData.loading}
      <div class="loading">
        <p>Loading jobs...</p>
      </div>
    {:else if jobsData.error}
      <div class="error">
        <p>Error loading jobs: {jobsData.error}</p>
        <button on:click={() => jobsStore.fetchJobs(currentDateRange.from, currentDateRange.to)}>
          Try Again
        </button>
      </div>
    {:else if jobs.length === 0}
      <div class="no-data">
        <p>No jobs found for the selected date range.</p>
      </div>
    {:else}
      <div class="jobs-table-container">
        <table class="jobs-table">
          <thead>
            <tr>
              <th>Job Name</th>
              <th>Status</th>
              <th>Last Run</th>
              <th>Avg Response Time</th>
              <th>Success Rate</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredJobs as job}
              <tr class="job-row">
                <td>
                  <a href="/job/{job.id}" use:link class="job-link">
                    {job.config.metadata.displayName}
                  </a>
                </td>
                <td>
                  <StatusBadge status={job.status} />
                </td>
                <td>
                  {#if job.lastExecution}
                    {formatDate(job.lastExecution)}
                  {:else}
                    Never
                  {/if}
                </td>
                <td>
                  {#if job.avgResponseTime}
                    {Math.round(job.avgResponseTime)}ms
                  {:else}
                    -
                  {/if}
                </td>
                <td class="success-rate">{Math.round(job.successRate * 10) / 10}%</td>
                <td>
                  <a href="/job/{job.id}" use:link class="btn-view">
                    View Details
                  </a>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  {:else}
    <!-- Job Detail View -->
    {#if jobDetailData.loading}
      <div class="loading">
        <p>Loading job details...</p>
      </div>
    {:else if jobDetailData.error}
      <div class="error">
        <p>Error loading job details: {jobDetailData.error}</p>
        <button on:click={() => push('/jobs')}>
          Back to Jobs
        </button>
      </div>
    {:else if selectedJob}
      <JobDetail selectedJob={selectedJob} executionHistory={selectedJob.executions || []} />
    {/if}
  {/if}
</div>

<style>
  .logs-container {
    padding: var(--spacing-xl);
    max-width: 1400px;
    margin: 0 auto;
  }

  .logs-header {
    text-align: center;
    margin-bottom: var(--spacing-xl);
  }

  .logs-header p {
    color: var(--text-secondary);
  }

  /* Job List Table */
  .jobs-table-container {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-lg);
    overflow: hidden;
  }

  .jobs-table {
    width: 100%;
    border-collapse: collapse;
  }

  .jobs-table th {
    background: var(--table-header);
    padding: var(--spacing-md);
    text-align: left;
    font-weight: 600;
    color: var(--text-primary);
    border-bottom: 2px solid var(--border-color);
  }

  .job-row {
    transition: background-color 0.2s;
  }

  .job-row:hover {
    background-color: var(--table-hover);
  }

  .jobs-table td {
    padding: var(--spacing-md);
    border-bottom: 1px solid var(--border-color);
  }

  .job-link {
    color: inherit;
    text-decoration: none;
    display: block;
    width: 100%;
  }

  .job-link:hover {
    color: var(--primary-color);
    text-decoration: underline;
  }

  .success-rate {
    font-weight: 500;
  }

  /* Buttons */
  .btn-view {
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: var(--font-md);
    transition: background-color 0.2s;
    text-decoration: none;
    display: inline-block;
    background: var(--primary-color);
    color: white;
  }

  .btn-view:hover {
    background: var(--primary-color-light);
  }

  .loading, .error, .no-data {
    text-align: center;
    padding: var(--spacing-xl);
    background: var(--card-bg);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-lg);
  }

  .loading p {
    color: var(--text-secondary);
  }

  .error p {
    color: var(--status-failed);
    margin-bottom: var(--spacing-md);
  }

  .error button {
    background: var(--primary-color);
    color: white;
    border: none;
    padding: var(--spacing-sm) var(--spacing-md);
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .error button:hover {
    background: var(--primary-color-light);
  }

  .no-data p {
    color: var(--text-secondary);
  }

  /* Responsive Design */
  @media (max-width: 768px) {
    .jobs-table-container {
      overflow-x: auto;
    }
    
    .jobs-table {
      min-width: 600px;
    }
  }
</style>
