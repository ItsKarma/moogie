<script>
  import { push } from 'svelte-spa-router';
  import { getStatusColor, getCheckTypeColor, formatDate, formatDuration, formatLabels } from '../lib/utils.js';
  import { dateRange, filterExecutionsByDateRange, calculateSuccessRateForRange } from '../lib/stores.js';
  import StatusBadge from '../components/StatusBadge.svelte';
  
  export let selectedJob;
  export let executionHistory;

  let currentDateRange = {};
  
  // Subscribe to date range changes
  dateRange.subscribe(range => {
    currentDateRange = range;
  });

  // Filter execution history and calculate metrics based on date range
  $: filteredExecutions = filterExecutionsByDateRange(executionHistory, currentDateRange);
  $: filteredSuccessRate = calculateSuccessRateForRange(executionHistory, currentDateRange);
  $: filteredExecutionCount = filteredExecutions.length;

  function goBack() {
    push('/jobs');
  }

  function toggleLogs(event, index) {
    event.preventDefault();
    
    // Find both details elements for this row
    const table = event.target.closest('table');
    const rows = table.querySelectorAll('tbody tr');
    const logsTrigger = rows[index * 2].querySelector('.logs-cell details');
    const logsContent = rows[index * 2 + 1].querySelector('.logs-td details');
    
    // Toggle both elements
    if (logsTrigger.open) {
      logsTrigger.open = false;
      logsContent.open = false;
    } else {
      logsTrigger.open = true;
      logsContent.open = true;
    }
  }
</script>

<!-- Job Detail View -->
<div class="job-detail">
  <div class="detail-header">
    <button class="btn-back" on:click={goBack}>← Back to Jobs</button>
    <h1>{selectedJob.config.metadata.displayName}</h1>
    <p class="job-description">{selectedJob.config.metadata.description}</p>
  </div>

  <div class="job-summary">
    <div class="summary-card">
      <h3>Current Status</h3>
      <StatusBadge status={selectedJob.status} size="large" />
    </div>
    <div class="summary-card">
      <h3>Total Executions</h3>
      <div class="summary-number">{filteredExecutionCount}</div>
    </div>
    <div class="summary-card">
      <h3>Success Rate</h3>
      <div class="summary-number">{filteredSuccessRate}%</div>
    </div>
    <div class="summary-card">
      <h3>Last Run</h3>
      <div class="summary-text">{formatDate(selectedJob.lastRun)}</div>
    </div>
  </div>

  <div class="job-config">
    <details class="config-details">
      <summary class="config-summary">
        <h2>Job Configuration</h2>
        <div class="config-preview">
          <div class="check-type" style="background-color: {getCheckTypeColor(selectedJob.config.kind)}">
            {selectedJob.config.kind}
          </div>
          <span class="config-preview-text">{selectedJob.config.metadata.id}</span>
        </div>
      </summary>
      
      <div class="config-content">
        <div class="config-sections">
          <div class="config-section">
            <h3>Metadata</h3>
            <div class="config-grid">
              <div class="config-item">
                <strong>API Version:</strong> {selectedJob.config.apiVersion}
              </div>
              <div class="config-item">
                <strong>Kind:</strong> {selectedJob.config.kind}
              </div>
              <div class="config-item">
                <strong>ID:</strong> {selectedJob.config.metadata.id}
              </div>
              <div class="config-item">
                <strong>Display Name:</strong> {selectedJob.config.metadata.displayName}
              </div>
              <div class="config-item">
                <strong>Description:</strong> {selectedJob.config.metadata.description}
              </div>
              <div class="config-item">
                <strong>Labels:</strong> {formatLabels(selectedJob.config.metadata.labels)}
              </div>
            </div>
          </div>

          <div class="config-section">
            <h3>Spec</h3>
            <div class="config-grid">
              {#each Object.entries(selectedJob.config.spec) as [key, value]}
                {#if key !== 'alerts'}
                  <div class="config-item">
                    <strong>{key}:</strong> 
                    {#if typeof value === 'object'}
                      <pre class="config-object">{JSON.stringify(value, null, 2)}</pre>
                    {:else}
                      {value}
                    {/if}
                  </div>
                {/if}
              {/each}
            </div>
          </div>

          <div class="config-section">
            <h3>Alerts</h3>
            <div class="config-grid">
              {#if selectedJob.config.spec.alerts}
                {#each Object.entries(selectedJob.config.spec.alerts) as [key, value]}
                  <div class="config-item">
                    <strong>{key}:</strong> 
                    {#if typeof value === 'object'}
                      <pre class="config-object">{JSON.stringify(value, null, 2)}</pre>
                    {:else}
                      {value}
                    {/if}
                  </div>
                {/each}
              {:else}
                <div class="config-item">
                  <strong>Status:</strong> No alerts configured
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    </details>
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
            <th>Logs</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredExecutions as execution, index}
            <tr class="history-row">
              <td>{formatDate(execution.timestamp)}</td>
              <td>
                <StatusBadge status={execution.status} size="small" />
              </td>
              <td class="response-time">{formatDuration(execution.responseTime)}</td>
              <td class="status-code">{execution.statusCode}</td>
              <td class="message">{execution.message}</td>
              <td class="logs-cell">
                <details>
                  <summary on:click={(e) => toggleLogs(e, index)}>View Logs</summary>
                </details>
              </td>
            </tr>
            <!-- Expandable logs section (accordion style) -->
            <tr class="logs-row">
              <td colspan="6" class="logs-td">
                <details>
                  <summary style="display: none;"></summary>
                  <div class="log-entries">
                    <div class="log-code-block" style="border-left-color: {getStatusColor(execution.status)}">{execution.logs.join('\n')}</div>
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

<style>
  /* Buttons */
  .btn-back {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s;
    text-decoration: none;
    display: inline-block;
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
    margin: 0 0 0.5rem 0;
  }

  .job-description {
    color: #666;
    font-size: 1rem;
    line-height: 1.5;
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
  .logs-cell {
    width: 120px;
    text-align: center;
  }

  .logs-cell summary {
    cursor: pointer;
    color: #ff3e00;
    font-size: 0.9rem;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    list-style: none;
    text-align: center;
  }

  .logs-cell summary:hover {
    background: #e9ecef;
  }

  .logs-row {
    background-color: #f8f9fa;
  }

  .logs-row .logs-td {
    padding: 0;
    border-top: none;
  }

  .logs-row details {
    margin: 0;
  }

  .logs-row details summary {
    display: none;
  }

  .log-entries {
    padding: 1rem;
    background: var(--card-background);
    border-top: 1px solid var(--border-color);
  }

  .log-code-block {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 0.8rem;
    color: var(--text-color);
    background: var(--code-background);
    padding: 0.1rem 1rem;
    border-radius: 4px;
    white-space: pre-wrap;
    word-break: break-word;
    overflow-x: auto;
    line-height: 1.4;
    border-left-width: 3px;
    border-left-style: solid;
  }

  /* Job Configuration */
  .job-config {
    margin-bottom: 2rem;
  }

  .config-details {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    overflow: hidden;
  }

  .config-summary {
    padding: 1rem;
    background: #f8f9fa;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    list-style: none;
    transition: background-color 0.2s;
  }

  .config-summary:hover {
    background: #e9ecef;
  }

  .config-summary h2 {
    color: #333;
    margin: 0;
    font-size: 1.2rem;
  }

  .config-preview {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .check-type {
    padding: 0.25rem 0.75rem;
    border-radius: 15px;
    color: white;
    font-size: 0.7rem;
    font-weight: bold;
  }

  .config-preview-text {
    font-size: 0.9rem;
    color: #666;
    font-weight: 500;
  }

  .config-content {
    border-top: 1px solid #e9ecef;
    padding: 1.5rem;
  }

  .config-sections {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 2rem;
  }

  .config-section {
    border: none;
    padding: 0;
  }

  .config-section h3 {
    margin: 0 0 1rem 0;
    color: #666;
    font-size: 1rem;
    text-transform: uppercase;
    font-weight: 600;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e9ecef;
  }

  .config-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0;
  }

  .config-item {
    padding: 0.25rem 0;
  }

  .config-item strong {
    color: #333;
    min-width: 120px;
    display: inline-block;
  }

  .config-object {
    background: #f8f9fa;
    border-radius: 4px;
    padding: 0.5rem;
    margin-top: 0.5rem;
    font-size: 0.8rem;
    overflow-x: auto;
    white-space: pre-wrap;
    border: 1px solid #e9ecef;
  }

  /* Dark Mode */
  @media (prefers-color-scheme: dark) {
    .detail-header h1, .execution-history h2, .job-config h2 {
      color: #ffffff;
    }
    
    .job-description {
      color: #cccccc;
    }
    
    .message {
      color: #cccccc;
    }
    
    .summary-card, .history-table-container, .config-details {
      background: #2a2a2a;
    }
    
    .history-table th, .config-summary {
      background: #333;
      color: #ffffff;
      border-bottom-color: #555;
    }
    
    .config-summary:hover {
      background: #404040;
    }
    
    .history-table td, .config-content {
      border-bottom-color: #555;
      color: #ffffff;
      background: #2a2a2a;
    }
    
    .config-section h3 {
      color: #cccccc;
      border-bottom-color: #555;
    }
    
    .summary-number, .summary-text {
      color: #ffffff;
    }
    
    .summary-card h3 {
      color: #cccccc;
    }

    .config-preview-text, .config-item strong, .config-item {
      color: #ffffff;
    }
    
    .config-summary h2 {
      color: #ffffff;
    }
    
    .btn-back {
      background: #444;
      color: #ffffff;
    }
    
    .btn-back:hover {
      background: #555;
    }
    
    .log-entries, .config-object {
      background: #333;
      border-color: #555;
    }

    .logs-cell summary {
      background: #444;
      border-color: #555;
      color: #ff6b6b;
    }

    .logs-cell summary:hover {
      background: #555;
    }
  }

  /* Responsive Design */
  @media (max-width: 768px) {
    .history-table-container {
      overflow-x: auto;
    }
    
    .history-table {
      min-width: 600px;
    }
    
    .job-summary {
      grid-template-columns: 1fr;
    }
    
    .config-sections {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
  }
</style>
