<script>
  import { onMount } from 'svelte';
  import { dateRange, dashboardStore, jobDetailStore } from '../lib/stores.js';
  import ChecksListSidebar from '../components/dashboard/ChecksListSidebar.svelte';
  import JobDetailPanel from '../components/dashboard/JobDetailPanel.svelte';

  let currentDateRange = {};
  let dashboardData = { data: null, loading: true, error: null };
  let jobDetailData = { data: null, loading: false, error: null };
  let selectedJobId = null;
  let searchQuery = '';
  let filterStatus = 'all';

  dateRange.subscribe(range => {
    currentDateRange = range;
    dashboardStore.fetchSummary(range.from, range.to);
  });

  dashboardStore.subscribe(value => {
    dashboardData = value;
  });

  jobDetailStore.subscribe(value => {
    jobDetailData = value;
  });

  onMount(() => {
    dashboardStore.fetchSummary(currentDateRange.from, currentDateRange.to);
  });

  $: summary = dashboardData.data || {};
  $: jobSummaries = summary.job_summaries || [];

  $: allJobs = jobSummaries.map(job => {
    const recentExecs = job.recent_executions || [];
    let status = 'success';
    let statusPriority = 3;
    
    if (!job.enabled) {
      status = 'disabled';
      statusPriority = 4;
    } else if (recentExecs.length > 0) {
      const mostRecentStatus = recentExecs[0].status;
      if (mostRecentStatus === 'failed' || mostRecentStatus === 'error') {
        status = 'failing';
        statusPriority = 0;
      } else if (mostRecentStatus === 'warning' || mostRecentStatus === 'timeout') {
        status = 'warning';
        statusPriority = 1;
      } else if (mostRecentStatus === 'success') {
        status = 'success';
        statusPriority = 2;
      }
    }
    
    return { ...job, recent_executions: recentExecs, status, statusPriority };
  });

  $: filteredJobs = allJobs
    .filter(job => {
      const matchesSearch = !searchQuery || 
        job.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        job.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
        job.labels?.service?.toLowerCase().includes(searchQuery.toLowerCase());
      
      const matchesStatus = filterStatus === 'all' || 
        (filterStatus === 'failing' && job.status === 'failing') ||
        (filterStatus === 'warning' && job.status === 'warning') ||
        (filterStatus === 'success' && job.status === 'success');
      
      return matchesSearch && matchesStatus;
    })
    .sort((a, b) => {
      if (a.statusPriority !== b.statusPriority) return a.statusPriority - b.statusPriority;
      if (a.success_rate !== b.success_rate) return a.success_rate - b.success_rate;
      return a.name.localeCompare(b.name);
    });

  $: selectedJob = selectedJobId ? allJobs.find(j => j.id === selectedJobId) : null;
  
  // Fetch detailed job data with full execution history when job is selected
  $: if (selectedJobId && currentDateRange.from && currentDateRange.to) {
    jobDetailStore.fetchJob(selectedJobId, currentDateRange.from, currentDateRange.to, 1000);
  }
  
  // Merge detailed execution data with summary data
  $: selectedJobWithExecutions = selectedJob && jobDetailData.data 
    ? { ...selectedJob, executions: jobDetailData.data.executions || [] }
    : selectedJob;
  
  $: if (!selectedJobId && filteredJobs.length > 0) {
    selectedJobId = filteredJobs[0].id;
  }

  function selectJob(jobId) {
    selectedJobId = jobId;
  }
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
    <div class="split-view">
      <ChecksListSidebar 
        jobs={filteredJobs}
        {selectedJobId}
        bind:searchQuery
        bind:filterStatus
        onSelectJob={selectJob}
      />
      
      <JobDetailPanel selectedJob={selectedJobWithExecutions} />
    </div>
  {/if}
</div>

<style>
  .dashboard {
    padding: var(--spacing-lg);
    max-width: calc(100vw - 2 * var(--spacing-lg));
    display: flex;
    flex-direction: column;
    margin: 0 auto;
    box-sizing: border-box;
  }

  .split-view {
    display: grid;
    grid-template-columns: 30% 70%;
    gap: var(--spacing-lg);
    min-height: 600px;
  }

  .split-view > :global(*) {
    min-width: 0;
  }

  .loading,
  .error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--spacing-xl);
    text-align: center;
  }

  .loading p,
  .error p {
    color: var(--text-secondary);
    font-size: var(--font-lg);
    margin: 0 0 var(--spacing-md) 0;
  }

  .error button {
    padding: var(--spacing-sm) var(--spacing-lg);
    background: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--radius-md);
    cursor: pointer;
    font-size: var(--font-md);
    transition: opacity 0.2s;
  }

  .error button:hover {
    opacity: 0.9;
  }

  @media (max-width: 1024px) {
    .split-view {
      grid-template-columns: 1fr;
    }
  }
</style>
