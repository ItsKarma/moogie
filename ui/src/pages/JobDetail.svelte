<script>
  import { push } from 'svelte-spa-router';
  import JobSummaryCards from '../components/job-detail/JobSummaryCards.svelte';
  import JobConfigSection from '../components/job-detail/JobConfigSection.svelte';
  import ExecutionHistoryTable from '../components/job-detail/ExecutionHistoryTable.svelte';
  
  export let selectedJob;
  export let executionHistory = [];

  // Parse config if it's a string
  $: parsedJob = selectedJob ? (() => {
    try {
      return {
        ...selectedJob,
        config: typeof selectedJob.config === 'string' 
          ? JSON.parse(selectedJob.config) 
          : selectedJob.config
      };
    } catch (e) {
      console.error('Error parsing job config:', e);
      return {
        ...selectedJob,
        config: {}
      };
    }
  })() : null;

  // The API already filters executions by date range, so we don't need to filter again
  // Just use the data as-is from the API
  $: allExecutions = executionHistory || [];
  $: executionCount = allExecutions.length;
  $: successRate = parsedJob?.success_rate || 0;

  function goBack() {
    push('/jobs');
  }
</script>

{#if parsedJob}
<!-- Job Detail View -->
<div class="job-detail">
  <div class="detail-header">
    <button class="btn-back" on:click={goBack}>← Back to Jobs</button>
    <h1>{parsedJob.name}</h1>
    <p class="job-description">{parsedJob.config?.metadata?.description || parsedJob.type + ' check for ' + parsedJob.name}</p>
  </div>

  <JobSummaryCards 
    job={parsedJob} 
    executionCount={executionCount} 
    successRate={successRate} 
  />

  <JobConfigSection job={parsedJob} />

  <ExecutionHistoryTable executions={allExecutions} />
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
