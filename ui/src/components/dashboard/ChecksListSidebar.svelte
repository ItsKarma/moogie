<script>
  import { getSuccessRateColor } from '../../lib/utils.js';

  export let jobs = [];
  export let selectedJobId = null;
  export let searchQuery = '';
  export let filterStatus = 'all';
  export let onSelectJob;

  function selectJob(jobId) {
    if (onSelectJob) {
      onSelectJob(jobId);
    }
  }
</script>

<div class="sidebar">
  <!-- Search and Filter -->
  <div class="search-card">
    <input 
      type="text" 
      class="search-input" 
      placeholder="Search checks..."
      bind:value={searchQuery}
    />
    <div class="filter-buttons">
      <button 
        class="filter-btn" 
        class:active={filterStatus === 'all'}
        on:click={() => filterStatus = 'all'}
      >
        All
      </button>
      <button 
        class="filter-btn failing" 
        class:active={filterStatus === 'failing'}
        on:click={() => filterStatus = 'failing'}
      >
        Failing
      </button>
      <button 
        class="filter-btn warning" 
        class:active={filterStatus === 'warning'}
        on:click={() => filterStatus = 'warning'}
      >
        Warning
      </button>
      <button 
        class="filter-btn success" 
        class:active={filterStatus === 'success'}
        on:click={() => filterStatus = 'success'}
      >
        Success
      </button>
    </div>
  </div>

  <!-- Checks List -->
  <div class="checks-list">
    {#if jobs.length === 0}
      <div class="empty-state">
        <p>No checks found.</p>
      </div>
    {:else}
      {#each jobs as job}
        <button 
          class="check-item" 
          class:selected={selectedJobId === job.id}
          class:status-failing={job.status === 'failing'}
          class:status-warning={job.status === 'warning'}
          class:status-success={job.status === 'success'}
          on:click={() => selectJob(job.id)}
        >
          <div class="check-info">
            <span class="check-name">{job.name}</span>
            <span class="check-meta">
              <span class="check-type">{job.type}</span>
              <span class="check-separator">•</span>
              <span class="check-rate" style="color: {getSuccessRateColor(job.success_rate)}">
                {job.success_rate.toFixed(1)}%
              </span>
            </span>
          </div>
          
          <div class="check-pills">
            {#if job.recent_executions && job.recent_executions.length > 0}
              {#each job.recent_executions.slice(0, 5).reverse() as execution}
                <div class="pill pill-{execution.status}"></div>
              {/each}
            {:else}
              <span class="no-data-small">—</span>
            {/if}
          </div>
        </button>
      {/each}
    {/if}
  </div>
</div>

<style>
  .sidebar {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
    height: 100%;
  }

  /* Search Card */
  .search-card {
    background: var(--card-bg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    padding: var(--spacing-md);
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  /* Filters */
  .filters {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .search-input {
    width: 100%;
    padding: var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    background: var(--bg-secondary);
    color: var(--text-primary);
    font-size: var(--font-sm);
  }

  .search-input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .filter-buttons {
    display: flex;
    gap: var(--spacing-xs);
  }

  .filter-btn {
    flex: 1;
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-size: var(--font-xs);
    cursor: pointer;
    transition: all 0.2s;
    text-transform: uppercase;
    font-weight: 500;
    letter-spacing: 0.05em;
  }

  .filter-btn:hover {
    background: var(--hover-bg);
  }

  .filter-btn.active {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .filter-btn.failing.active {
    background: var(--status-failed);
    border-color: var(--status-failed);
  }

  .filter-btn.warning.active {
    background: var(--status-warning);
    border-color: var(--status-warning);
  }

  .filter-btn.success.active {
    background: var(--status-success);
    border-color: var(--status-success);
  }

  /* Checks List */
  .checks-list {
    background: var(--card-bg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    overflow: hidden;
  }

  /* Check Items */
  .check-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-bottom: 1px solid var(--border-color);
    background: transparent;
    width: 100%;
    text-align: left;
    cursor: pointer;
    transition: all 0.2s;
    color: inherit;
  }

  .check-item:last-child {
    border-bottom: none;
  }

  .check-item:hover {
    background: var(--hover-bg);
  }

  .check-item.selected {
    background: var(--primary-color);
    color: white;
  }

  .check-item.selected .check-type,
  .check-item.selected .check-rate,
  .check-item.selected .check-separator {
    color: rgba(255, 255, 255, 0.8);
  }

  .check-item.status-failing {
    border-left: 3px solid var(--status-failed);
  }

  .check-item.status-warning {
    border-left: 3px solid var(--status-warning);
  }

  .check-item.status-success {
    border-left: 3px solid var(--status-success);
  }

  .check-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    flex: 1;
    min-width: 0;
  }

  .check-name {
    font-size: var(--font-sm);
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .check-meta {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    font-size: var(--font-xs);
    color: var(--text-secondary);
  }

  .check-type {
    text-transform: capitalize;
  }

  .check-separator {
    opacity: 0.5;
  }

  .check-rate {
    font-weight: 600;
  }

  /* Status pills */
  .check-pills {
    display: flex;
    gap: 3px;
    align-items: center;
    padding-left: var(--spacing-sm);
  }

  .pill {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .pill-success {
    background: var(--status-success);
  }

  .pill-failed,
  .pill-error {
    background: var(--status-failed);
  }

  .pill-warning,
  .pill-timeout {
    background: var(--status-warning);
  }

  .no-data-small {
    font-size: var(--font-xs);
    color: var(--text-secondary);
  }

  .empty-state {
    padding: var(--spacing-xl);
    text-align: center;
    color: var(--text-secondary);
  }
</style>
