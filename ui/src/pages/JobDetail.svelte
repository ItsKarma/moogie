<script>
  import { push } from 'svelte-spa-router';
  import { dateRange, filterExecutionsByDateRange, calculateSuccessRateForRange } from '../lib/stores.js';
  import JobSummaryCards from '../components/JobSummaryCards.svelte';
  import JobConfigSection from '../components/JobConfigSection.svelte';
  import ExecutionHistoryTable from '../components/ExecutionHistoryTable.svelte';
  
  export let selectedJob;
  export let executionHistory = [];

  let currentDateRange = {};
  
  // Subscribe to date range changes
  dateRange.subscribe(range => {
    currentDateRange = range;
  });

  // Filter execution history and calculate metrics based on date range
  $: filteredExecutions = filterExecutionsByDateRange(executionHistory || [], currentDateRange);
  $: filteredSuccessRate = calculateSuccessRateForRange(executionHistory || [], currentDateRange);
  $: filteredExecutionCount = filteredExecutions.length;

  function goBack() {
    push('/jobs');
  }
</script>

{#if selectedJob}
<!-- Job Detail View -->
<div class="job-detail">
  <div class="detail-header">
    <button class="btn-back" on:click={goBack}>← Back to Jobs</button>
    <h1>{selectedJob.name}</h1>
    <p class="job-description">{selectedJob.config?.metadata?.description || selectedJob.type + ' check for ' + selectedJob.name}</p>
  </div>

  <JobSummaryCards 
    job={selectedJob} 
    executionCount={filteredExecutionCount} 
    successRate={filteredSuccessRate} 
  />

  <JobConfigSection job={selectedJob} />

  <ExecutionHistoryTable executions={filteredExecutions} />
</div>

<style>
  .job-detail {
    max-width: 1400px;
    margin: 0 auto;
    padding: var(--spacing-xl);
  }

  .detail-header {
    margin-bottom: var(--spacing-xl);
  }

  .detail-header h1 {
    margin: var(--spacing-md) 0 var(--spacing-sm) 0;
    font-size: var(--font-xxl);
    color: var(--text-primary);
  }

  .job-description {
    color: var(--text-secondary);
    font-size: var(--font-base);
    margin: 0;
  }

  .btn-back {
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: var(--font-md);
    transition: background-color 0.2s;
    text-decoration: none;
    background: var(--bg-secondary);
    color: var(--text-primary);
  }

  .btn-back:hover {
    background: var(--hover-bg);
  }

  .loading {
    text-align: center;
    padding: var(--spacing-xl);
    color: var(--text-secondary);
  }
</style>
{:else}
<div class="loading">
  <p>No job selected</p>
</div>
{/if}
