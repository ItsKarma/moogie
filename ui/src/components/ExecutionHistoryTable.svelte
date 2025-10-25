<script>
  import { getStatusColor, formatDate, formatDuration } from '../lib/utils.js';
  import StatusBadge from './StatusBadge.svelte';
  
  export let executions = [];

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

<div class="execution-history">
  <h2>Execution History</h2>
  {#if executions.length === 0}
    <div class="no-executions">
      <p>No execution history available for the selected date range.</p>
    </div>
  {:else}
    <div class="history-table-container">
      <table class="history-table">
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>Status</th>
            <th>Response Time</th>
            <th>Status Code</th>
            <th>Message</th>
            <th>Details</th>
          </tr>
        </thead>
        <tbody>
          {#each executions as execution, index}
            <tr class="history-row">
              <td>{formatDate(execution.timestamp)}</td>
              <td>
                <StatusBadge status={execution.status} size="small" />
              </td>
              <td class="response-time">{formatDuration(execution.response_time)}</td>
              <td class="status-code">{execution.details?.response_code || 'N/A'}</td>
              <td class="message">{execution.details?.message || execution.details?.error || 'No message'}</td>
              <td class="logs-cell">
                <details>
                  <summary on:click={(e) => toggleLogs(e, index)}>View Details</summary>
                </details>
              </td>
            </tr>
            <!-- Expandable details section (accordion style) -->
            <tr class="logs-row">
              <td colspan="6" class="logs-td">
                <details>
                  <summary style="display: none;"></summary>
                  <div class="log-entries">
                    <div class="log-code-block" style="border-left-color: {getStatusColor(execution.status)}">
                      <pre>{JSON.stringify(execution.details, null, 2)}</pre>
                    </div>
                  </div>
                </details>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .execution-history {
    margin: var(--spacing-xl) 0;
  }

  .execution-history h2 {
    margin: 0 0 var(--spacing-lg) 0;
    font-size: var(--font-xl);
    color: var(--text-primary);
  }

  .no-executions {
    background: var(--card-bg);
    padding: var(--spacing-xl);
    border-radius: var(--radius-md);
    text-align: center;
    color: var(--text-secondary);
  }

  .history-table-container {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    overflow: hidden;
    box-shadow: var(--shadow-md);
  }

  .history-table {
    width: 100%;
    border-collapse: collapse;
  }

  .history-table thead {
    background: var(--bg-secondary);
  }

  .history-table th {
    padding: var(--spacing-md);
    text-align: left;
    font-weight: 600;
    color: var(--text-secondary);
    font-size: var(--font-sm);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .history-row td {
    padding: var(--spacing-md);
    border-bottom: 1px solid var(--border-color);
    color: var(--text-primary);
  }

  .history-row:hover {
    background: var(--hover-bg);
  }

  .response-time {
    font-family: 'Courier New', monospace;
    font-weight: 500;
  }

  .status-code {
    font-family: 'Courier New', monospace;
    font-weight: 500;
  }

  .message {
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .logs-cell details {
    cursor: pointer;
  }

  .logs-cell summary {
    color: var(--primary-color);
    font-size: var(--font-sm);
    user-select: none;
  }

  .logs-cell summary:hover {
    text-decoration: underline;
  }

  .logs-row {
    display: none;
  }

  .logs-row:has(details[open]) {
    display: table-row;
  }

  .logs-td {
    padding: 0 !important;
    background: var(--bg-secondary);
  }

  .log-entries {
    padding: var(--spacing-md);
  }

  .log-code-block {
    border-left: 4px solid;
    padding: var(--spacing-md);
    background: var(--card-bg);
    border-radius: var(--radius-sm);
    overflow-x: auto;
  }

  .log-code-block pre {
    margin: 0;
    font-size: var(--font-sm);
    line-height: 1.5;
    color: var(--text-primary);
  }

  @media (max-width: 768px) {
    .history-table-container {
      overflow-x: auto;
    }
    
    .history-table {
      min-width: 600px;
    }
  }
</style>
