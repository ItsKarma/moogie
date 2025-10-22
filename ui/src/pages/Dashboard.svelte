<script>
  import { getSuccessRateColor, calculateOverallSuccessRate } from '../lib/utils.js';
  import { mockDashboardJobs } from '../lib/mockData.js';
  import StatusBadge from '../components/StatusBadge.svelte';

  // Use imported mock data
  let jobs = mockDashboardJobs;

  $: overallSuccessRate = calculateOverallSuccessRate(jobs);
</script>

<div class="dashboard">
  <div class="stats-grid">
    <div class="stat-card">
      <h3>Active Jobs</h3>
      <div class="stat-number">4</div>
    </div>
    <div class="stat-card">
      <h3>Running</h3>
      <div class="stat-number running">1</div>
    </div>
    <div class="stat-card">
      <h3>Failed</h3>
      <div class="stat-number failed">1</div>
    </div>
    <div class="stat-card">
      <h3>Success Rate</h3>
      <div class="stat-number" style="color: {getSuccessRateColor(overallSuccessRate)}">{overallSuccessRate}%</div>
    </div>
  </div>

  <div class="jobs-section">
    <h2>Job Status</h2>
    <div class="jobs-grid">
      {#each jobs as job}
        <div class="job-card">
          <div class="job-header">
            <h3>{job.name}</h3>
            <StatusBadge status={job.status} size="small" />
          </div>
          <div class="job-details">
            <p><strong>Last Run:</strong> {job.lastRun}</p>
            <p><strong>Next Run:</strong> {job.nextRun}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .dashboard {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    text-align: center;
  }

  .stat-card h3 {
    margin: 0 0 0.5rem 0;
    color: #333;
    font-size: 0.9rem;
    text-transform: uppercase;
    font-weight: 600;
  }

  .stat-number {
    font-size: 2rem;
    font-weight: bold;
    color: #333;
  }

  .stat-number.running {
    color: #2196F3;
  }

  .stat-number.failed {
    color: #F44336;
  }

  .jobs-section h2 {
    color: #333;
    margin-bottom: 1rem;
  }

  .jobs-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
  }

  .job-card {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  }

  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .job-header h3 {
    margin: 0;
    color: #333;
    font-size: 1.1rem;
  }

  .job-details p {
    margin: 0.5rem 0;
    color: #666;
    font-size: 0.9rem;
  }

  @media (prefers-color-scheme: dark) {
    .dashboard {
      color: #ffffff;
    }
    
    .jobs-section h2, .job-header h3 {
      color: #ffffff;
    }
    
    .stat-card, .job-card {
      background: #2a2a2a;
      color: #ffffff;
    }
    
    .stat-card h3 {
      color: #ffffff;
      font-weight: 700;
    }
    
    .stat-number {
      color: #ffffff;
    }
    
    .job-details p {
      color: #e0e0e0;
    }
  }
</style>
