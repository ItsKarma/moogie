<script>
  import { push, link } from 'svelte-spa-router';
  import { formatDate } from '../lib/utils.js';
  import { mockJobs, mockExecutionHistory } from '../lib/mockData.js';
  import StatusBadge from '../components/StatusBadge.svelte';
  import JobDetail from './JobDetail.svelte';
  
  // Route parameters
  export let params = {};
  
  // Use imported mock data
  let jobs = mockJobs;

  // Mock execution history for selected job
  let selectedJob = null;
  let executionHistory = [];

  // Watch for route parameter changes
  $: if (params.jobId) {
    const jobId = params.jobId;
    const job = jobs.find(j => j.id === jobId || j.config.metadata.id === jobId);
    if (job) {
      selectJob(job);
    } else {
      // Job not found, redirect to jobs list
      push('/jobs');
    }
  } else {
    selectedJob = null;
    executionHistory = [];
  }

  function selectJob(job) {
    selectedJob = job;
    executionHistory = mockExecutionHistory; // Use imported mock data
    // Update URL if not already there
    if (!params.jobId || params.jobId !== job.id) {
      push(`/job/${job.id}`);
    }
  }
</script>

<div class="logs-container">
  {#if !selectedJob}
    <!-- Job List View -->
    <div class="logs-header">
      <p>Click on any job to view its configuration and execution history</p>
    </div>

    <div class="jobs-table-container">
      <table class="jobs-table">
        <thead>
          <tr>
            <th>Job Name</th>
            <th>Status</th>
            <th>Last Run</th>
            <th>Executions</th>
            <th>Success Rate</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each jobs as job}
            <tr class="job-row">
              <td class="job-name">
                <a href="/job/{job.id}" use:link class="job-link">
                  {job.config.metadata.displayName}
                </a>
              </td>
              <td>
                <StatusBadge status={job.status} />
              </td>
              <td>{formatDate(job.lastRun)}</td>
              <td>{job.executions}</td>
              <td class="success-rate">{job.successRate}%</td>
              <td>
                <a href="/job/{job.id}" use:link class="btn-view">
                  View Job
                </a>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <!-- Job Detail View -->
    <JobDetail {selectedJob} {executionHistory} />
  {/if}
</div>

<style>
  .logs-container {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .logs-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .logs-header p {
    color: #666;
  }

  /* Job List Table */
  .jobs-table-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    overflow: hidden;
  }

  .jobs-table {
    width: 100%;
    border-collapse: collapse;
  }

  .jobs-table th {
    background: #f8f9fa;
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #333;
    border-bottom: 2px solid #e9ecef;
  }

  .job-row {
    transition: background-color 0.2s;
  }

  .job-row:hover {
    background-color: #f8f9fa;
  }

  .jobs-table td {
    padding: 1rem;
    border-bottom: 1px solid #e9ecef;
  }

  .job-name {
    font-weight: 500;
    color: #333;
  }

  .job-link {
    color: inherit;
    text-decoration: none;
    display: block;
    width: 100%;
  }

  .job-link:hover {
    color: #7D471F;
    text-decoration: underline;
  }

  .success-rate {
    font-weight: 500;
  }

  /* Buttons */
  .btn-view {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s;
    text-decoration: none;
    display: inline-block;
    background: #7D471F;
    color: white;
  }

  .btn-view:hover {
    background: #e63200;
  }

  /* Dark Mode */
  @media (prefers-color-scheme: dark) {
    .logs-container {
      color: #ffffff;
    }
    
    .logs-header p {
      color: #cccccc;
    }
    
    .jobs-table-container {
      background: #2a2a2a;
    }
    
    .jobs-table th {
      background: #333;
      color: #ffffff;
      border-bottom-color: #555;
    }
    
    .jobs-table td {
      border-bottom-color: #555;
      color: #ffffff;
    }
    
    .job-name {
      color: #ffffff;
    }

    .job-link {
      color: #ffffff;
    }

    .job-link:hover {
      color: #ff6b6b;
    }
    
    .job-row:hover {
      background-color: #333;
    }
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
