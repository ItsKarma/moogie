<script>
  // Mock job data - will be replaced with API calls later
  let jobs = [
    { id: 1, name: "Website Health Check", status: "running", lastRun: "2 minutes ago", nextRun: "in 28 minutes" },
    { id: 2, name: "API Response Time", status: "success", lastRun: "15 minutes ago", nextRun: "in 15 minutes" },
    { id: 3, name: "Database Connection", status: "failed", lastRun: "1 hour ago", nextRun: "in 2 hours" },
    { id: 4, name: "SSL Certificate Check", status: "warning", lastRun: "3 hours ago", nextRun: "in 21 hours" },
  ];

  function getStatusColor(status) {
    switch(status) {
      case 'running': return '#2196F3';
      case 'success': return '#4CAF50';
      case 'failed': return '#F44336';
      case 'warning': return '#FF9800';
      default: return '#9E9E9E';
    }
  }
</script>

<div class="dashboard">
  <h1>Job Dashboard</h1>
  
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
      <div class="stat-number">75%</div>
    </div>
  </div>

  <div class="jobs-section">
    <h2>Job Status</h2>
    <div class="jobs-grid">
      {#each jobs as job}
        <div class="job-card">
          <div class="job-header">
            <h3>{job.name}</h3>
            <span class="status-badge" style="background-color: {getStatusColor(job.status)}">
              {job.status.toUpperCase()}
            </span>
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

  h1 {
    color: #333;
    text-align: center;
    margin-bottom: 2rem;
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
    color: #666;
    font-size: 0.9rem;
    text-transform: uppercase;
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

  .status-badge {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    color: white;
    font-size: 0.8rem;
    font-weight: bold;
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
    
    h1, .jobs-section h2, .job-header h3 {
      color: #ffffff;
    }
    
    .stat-card, .job-card {
      background: #2a2a2a;
      color: #ffffff;
    }
    
    .stat-card h3, .job-details p {
      color: #cccccc;
    }
  }
</style>
