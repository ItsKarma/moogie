<script>
  import { getCheckTypeColor } from '../lib/utils.js';
  
  export let job;
</script>

<div class="job-config">
  <details class="config-details">
    <summary class="config-summary">
      <h2>Job Configuration</h2>
      <div class="config-preview">
        <div class="check-type" style="background-color: {getCheckTypeColor(job.type)}">
          {job.type}
        </div>
        <span class="config-preview-text">{job.name}</span>
      </div>
    </summary>
    
    <div class="config-content">
      <div class="config-sections">
        <div class="config-section">
          <h3>Metadata</h3>
          <div class="config-grid">
            <div class="config-item">
              <strong>API Version:</strong> {job.config.apiVersion || 'N/A'}
            </div>
            <div class="config-item">
              <strong>Kind:</strong> {job.config.kind || job.type}
            </div>
            <div class="config-item">
              <strong>Job ID:</strong> {job.id}
            </div>
            <div class="config-item">
              <strong>Name:</strong> {job.name}
            </div>
            <div class="config-item">
              <strong>Type:</strong> {job.type}
            </div>
            {#if job.config.metadata}
              {#each Object.entries(job.config.metadata) as [key, value]}
                <div class="config-item">
                  <strong>{key}:</strong> {typeof value === 'object' ? JSON.stringify(value) : value}
                </div>
              {/each}
            {/if}
          </div>
        </div>

        <div class="config-section">
          <h3>Spec</h3>
          <div class="config-grid">
            {#each Object.entries(job.config.spec) as [key, value]}
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
            {#if job.config.spec.alerts}
              {#each Object.entries(job.config.spec.alerts) as [key, value]}
                <div class="config-item">
                  <strong>{key}:</strong> {value}
                </div>
              {/each}
            {:else}
              <div class="config-item">
                <em>No alerts configured</em>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </details>
</div>

<style>
  .job-config {
    margin: var(--spacing-xl) 0;
  }

  .config-details {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
  }

  .config-summary {
    padding: var(--spacing-lg);
    cursor: pointer;
    list-style: none;
    display: flex;
    justify-content: space-between;
    align-items: center;
    user-select: none;
  }

  .config-summary::-webkit-details-marker {
    display: none;
  }

  .config-summary h2 {
    margin: 0;
    font-size: var(--font-lg);
    color: var(--text-primary);
  }

  .config-preview {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
  }

  .check-type {
    padding: var(--spacing-sm) var(--spacing-md);
    border-radius: var(--radius-sm);
    color: white;
    font-weight: 500;
    font-size: var(--font-sm);
    text-transform: uppercase;
  }

  .config-preview-text {
    color: var(--text-secondary);
    font-size: var(--font-md);
  }

  .config-content {
    padding: 0 var(--spacing-lg) var(--spacing-lg);
    border-top: 1px solid var(--border-color);
  }

  .config-sections {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: var(--spacing-xl);
    margin-top: var(--spacing-lg);
  }

  .config-section h3 {
    margin: 0 0 var(--spacing-md) 0;
    font-size: var(--font-base);
    color: var(--text-primary);
    border-bottom: 2px solid var(--primary-color);
    padding-bottom: var(--spacing-sm);
  }

  .config-grid {
    display: grid;
    gap: var(--spacing-md);
  }

  .config-item {
    font-size: var(--font-md);
    line-height: 1.5;
    color: var(--text-primary);
  }

  .config-item strong {
    color: var(--text-secondary);
    display: inline-block;
    min-width: 120px;
  }

  .config-object {
    margin-top: var(--spacing-sm);
    padding: var(--spacing-md);
    background: var(--bg-secondary);
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    overflow-x: auto;
  }

  @media (max-width: 768px) {
    .config-sections {
      grid-template-columns: 1fr;
      gap: var(--spacing-lg);
    }
  }
</style>
