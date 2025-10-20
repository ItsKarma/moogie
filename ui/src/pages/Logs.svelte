<script>
  import { onMount } from 'svelte';
  
  // Mock job data - will be replaced with API calls later
  let jobs = [
    { 
      id: 1, 
      name: "Website Health Check", 
      status: "running",
      lastRun: "2024-10-20 14:30:00",
      executions: 1247,
      successRate: 98.2
    },
    { 
      id: 2, 
      name: "API Response Time", 
      status: "success",
      lastRun: "2024-10-20 14:15:00",
      executions: 892,
      successRate: 99.1
    },
    { 
      id: 3, 
      name: "Database Connection", 
      status: "failed",
      lastRun: "2024-10-20 13:30:00",
      executions: 456,
      successRate: 95.8
    },
    { 
      id: 4, 
      name: "SSL Certificate Check", 
      status: "warning",
      lastRun: "2024-10-20 11:00:00",
      executions: 124,
      successRate: 100.0
    },
  ];

  // Mock execution history for selected job
  let selectedJob = null;
  let executionHistory = [];

  // Generate mock execution history
  function generateExecutionHistory(jobId) {
    const statuses = ['success', 'failed', 'timeout', 'warning'];
    const history = [];
    
    for (let i = 0; i < 20; i++) {
      const date = new Date();
      date.setMinutes(date.getMinutes() - (i * 30)); // Every 30 minutes back
      
      const status = Math.random() > 0.8 ? 'failed' : 'success'; // 80% success rate
      const responseTime = Math.floor(Math.random() * 2000) + 100; // 100-2100ms
      
      history.push({
        id: i + 1,
        timestamp: date.toISOString(),
        status: status,
        responseTime: responseTime,
        statusCode: status === 'success' ? 200 : (Math.random() > 0.5 ? 404 : 500),
        message: status === 'success' 
          ? 'Request completed successfully' 
          : 'Connection timeout or server error',
        logs: generateLogEntries(status, responseTime)
      });
    }
    
    return history;
  }

  function generateLogEntries(status, responseTime) {
    const timestamp = new Date().toISOString();
    const logs = [
      `[${timestamp}] Starting job execution`,
      `[${timestamp}] Sending HTTP request...`,
      `[${timestamp}] Response received in ${responseTime}ms`,
    ];
    
    if (status === 'success') {
      logs.push(`[${timestamp}] ✓ Job completed successfully`);
    } else {
      logs.push(`[${timestamp}] ✗ Job failed: Connection timeout`);
      logs.push(`[${timestamp}] Retrying in 30 seconds...`);
    }
    
    return logs;
  }

  function selectJob(job) {
    selectedJob = job;
    executionHistory = generateExecutionHistory(job.id);
  }

  function goBack() {
    selectedJob = null;
    executionHistory = [];
  }

  function getStatusColor(status) {
    switch(status) {
      case 'running': return '#2196F3';
      case 'success': return '#4CAF50';
      case 'failed': return '#F44336';
      case 'warning': return '#FF9800';
      case 'timeout': return '#FF5722';
      default: return '#9E9E9E';
    }
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleString();
  }

  function formatDuration(ms) {
    return `${ms}ms`;
  }
</script>

<div class="logs-container">
  {#if !selectedJob}
    <!-- Job List View -->
    <div class="logs-header">
      <h1>Job Logs</h1>
      <p>Click on any job to view its execution history and logs</p>
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
            <tr class="job-row" on:click={() => selectJob(job)}>
              <td class="job-name">{job.name}</td>
              <td>
                <span class="status-badge" style="background-color: {getStatusColor(job.status)}">
                  {job.status.toUpperCase()}
                </span>
              </td>
              <td>{formatDate(job.lastRun)}</td>
              <td>{job.executions}</td>
              <td class="success-rate">{job.successRate}%</td>
              <td>
                <button class="btn-view" on:click|stopPropagation={() => selectJob(job)}>
                  View Logs
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <!-- Job Detail View -->
    <div class="job-detail">
      <div class="detail-header">
        <button class="btn-back" on:click={goBack}>← Back to Jobs</button>
        <h1>{selectedJob.name}</h1>
      </div>

      <div class="job-summary">
        <div class="summary-card">
          <h3>Current Status</h3>
          <span class="status-badge large" style="background-color: {getStatusColor(selectedJob.status)}">
            {selectedJob.status.toUpperCase()}
          </span>
        </div>
        <div class="summary-card">
          <h3>Total Executions</h3>
          <div class="summary-number">{selectedJob.executions}</div>
        </div>
        <div class="summary-card">
          <h3>Success Rate</h3>
          <div class="summary-number">{selectedJob.successRate}%</div>
        </div>
        <div class="summary-card">
          <h3>Last Run</h3>
          <div class="summary-text">{formatDate(selectedJob.lastRun)}</div>
        </div>
      </div>

      <div class="execution-history">
        <h2>Execution History</h2>
        <div class="history-table-container">
          <table class="history-table">
            <thead>
              <tr>
                <th>Timestamp</th>
                <th>Status</th>
                <th>Response Time</th>
                <th>Status Code</th>
                <th>Message</th>
              </tr>
            </thead>
            <tbody>
              {#each executionHistory as execution}
                <tr class="history-row">
                  <td>{formatDate(execution.timestamp)}</td>
                  <td>
                    <span class="status-badge small" style="background-color: {getStatusColor(execution.status)}">
                      {execution.status.toUpperCase()}
                    </span>
                  </td>
                  <td class="response-time">{formatDuration(execution.responseTime)}</td>
                  <td class="status-code">{execution.statusCode}</td>
                  <td class="message">{execution.message}</td>
                </tr>
                <!-- Expandable logs section -->
                <tr class="logs-row">
                  <td colspan="5">
                    <details>
                      <summary>View Logs</summary>
                      <div class="log-entries">
                        {#each execution.logs as logEntry}
                          <div class="log-entry">{logEntry}</div>
                        {/each}
                      </div>
                    </details>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
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

  .logs-header h1 {
    color: #333;
    margin-bottom: 0.5rem;
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
    cursor: pointer;
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

  .success-rate {
    font-weight: 500;
  }

  /* Status Badges */
  .status-badge {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    color: white;
    font-size: 0.8rem;
    font-weight: bold;
    text-align: center;
    display: inline-block;
  }

  .status-badge.large {
    padding: 0.5rem 1rem;
    font-size: 1rem;
  }

  .status-badge.small {
    padding: 0.2rem 0.5rem;
    font-size: 0.7rem;
  }

  /* Buttons */
  .btn-view, .btn-back {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s;
  }

  .btn-view {
    background: #7D471F;
    color: white;
  }

  .btn-view:hover {
    background: #e63200;
  }

  .btn-back {
    background: #f5f5f5;
    color: #333;
    margin-bottom: 1rem;
  }

  .btn-back:hover {
    background: #e0e0e0;
  }

  /* Job Detail View */
  .detail-header h1 {
    color: #333;
    margin: 0;
  }

  .job-summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin: 2rem 0;
  }

  .summary-card {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    text-align: center;
  }

  .summary-card h3 {
    margin: 0 0 1rem 0;
    color: #666;
    font-size: 0.9rem;
    text-transform: uppercase;
  }

  .summary-number {
    font-size: 2rem;
    font-weight: bold;
    color: #333;
  }

  .summary-text {
    font-size: 1rem;
    color: #333;
  }

  /* Execution History */
  .execution-history h2 {
    color: #333;
    margin-bottom: 1rem;
  }

  .history-table-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    overflow: hidden;
  }

  .history-table {
    width: 100%;
    border-collapse: collapse;
  }

  .history-table th {
    background: #f8f9fa;
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #333;
    border-bottom: 2px solid #e9ecef;
  }

  .history-table td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #e9ecef;
  }

  .response-time {
    font-family: 'Monaco', 'Menlo', monospace;
    color: #2196F3;
  }

  .status-code {
    font-family: 'Monaco', 'Menlo', monospace;
    font-weight: bold;
  }

  .message {
    color: #666;
  }

  /* Log Entries */
  .logs-row td {
    padding: 0;
  }

  .logs-row details {
    margin: 0.5rem 1rem;
  }

  .logs-row summary {
    cursor: pointer;
    color: #ff3e00;
    font-size: 0.9rem;
    padding: 0.5rem 0;
  }

  .log-entries {
    background: #f8f9fa;
    border-radius: 4px;
    padding: 1rem;
    margin: 0.5rem 0;
    max-height: 200px;
    overflow-y: auto;
  }

  .log-entry {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 0.8rem;
    color: #333;
    margin-bottom: 0.25rem;
    word-break: break-all;
  }

  /* Dark Mode */
  @media (prefers-color-scheme: dark) {
    .logs-container {
      color: #ffffff;
    }
    
    .logs-header h1, .detail-header h1, .execution-history h2 {
      color: #ffffff;
    }
    
    .logs-header p, .message {
      color: #cccccc;
    }
    
    .jobs-table-container, .summary-card, .history-table-container {
      background: #2a2a2a;
    }
    
    .jobs-table th, .history-table th {
      background: #333;
      color: #ffffff;
      border-bottom-color: #555;
    }
    
    .jobs-table td, .history-table td {
      border-bottom-color: #555;
      color: #ffffff;
    }
    
    .job-name, .summary-number, .summary-text {
      color: #ffffff;
    }
    
    .summary-card h3 {
      color: #cccccc;
    }
    
    .job-row:hover {
      background-color: #333;
    }
    
    .btn-back {
      background: #444;
      color: #ffffff;
    }
    
    .btn-back:hover {
      background: #555;
    }
    
    .log-entries {
      background: #333;
    }
    
    .log-entry {
      color: #ffffff;
    }
    
    .logs-row summary {
      color: #ff6b6b;
    }
  }

  /* Responsive Design */
  @media (max-width: 768px) {
    .jobs-table-container, .history-table-container {
      overflow-x: auto;
    }
    
    .jobs-table, .history-table {
      min-width: 600px;
    }
    
    .job-summary {
      grid-template-columns: 1fr;
    }
  }
</style>
