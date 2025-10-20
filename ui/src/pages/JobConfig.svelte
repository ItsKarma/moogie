<script>
  // Form data
  let jobName = '';
  let jobType = 'http';
  let url = '';
  let method = 'GET';
  let headers = '';
  let schedule = '*/30 * * * *'; // Every 30 minutes
  let timeout = 30;
  let retries = 3;
  let alertOnFailure = true;
  let alertEmail = '';
  
  // Job types
  const jobTypes = [
    { value: 'http', label: 'HTTP/HTTPS Request' },
    { value: 'tcp', label: 'TCP Port Check' },
    { value: 'dns', label: 'DNS Resolution' },
    { value: 'ssl', label: 'SSL Certificate' },
    { value: 'ping', label: 'Ping Test' }
  ];

  function handleSubmit() {
    // This will be connected to an API endpoint later
    console.log('Creating job:', {
      jobName,
      jobType,
      url,
      method,
      headers,
      schedule,
      timeout,
      retries,
      alertOnFailure,
      alertEmail
    });
    
    // Show success message (placeholder)
    alert('Job created successfully! (This is a placeholder - will be replaced with proper API integration)');
    
    // Reset form
    resetForm();
  }

  function resetForm() {
    jobName = '';
    jobType = 'http';
    url = '';
    method = 'GET';
    headers = '';
    schedule = '*/30 * * * *';
    timeout = 30;
    retries = 3;
    alertOnFailure = true;
    alertEmail = '';
  }
</script>

<div class="job-config">
  <h1>Create New Job</h1>
  
  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-section">
      <h2>Basic Configuration</h2>
      
      <div class="form-group">
        <label for="jobName">Job Name</label>
        <input
          type="text"
          id="jobName"
          bind:value={jobName}
          placeholder="e.g., Website Health Check"
          required
        />
      </div>

      <div class="form-group">
        <label for="jobType">Job Type</label>
        <select id="jobType" bind:value={jobType}>
          {#each jobTypes as type}
            <option value={type.value}>{type.label}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <label for="url">Target URL/Address</label>
        <input
          type="text"
          id="url"
          bind:value={url}
          placeholder="https://example.com/api/health"
          required
        />
      </div>

      {#if jobType === 'http'}
        <div class="form-group">
          <label for="method">HTTP Method</label>
          <select id="method" bind:value={method}>
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
            <option value="HEAD">HEAD</option>
          </select>
        </div>

        <div class="form-group">
          <label for="headers">Custom Headers (JSON format)</label>
          <textarea
            id="headers"
            bind:value={headers}
            placeholder={'{"Authorization": "Bearer token", "Content-Type": "application/json"}'}
            rows="3"
          ></textarea>
        </div>
      {/if}
    </div>

    <div class="form-section">
      <h2>Schedule & Timing</h2>
      
      <div class="form-group">
        <label for="schedule">Cron Schedule</label>
        <input
          type="text"
          id="schedule"
          bind:value={schedule}
          placeholder="*/30 * * * * (every 30 minutes)"
        />
        <div class="help-text">
          Use cron syntax: minute hour day month weekday
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="timeout">Timeout (seconds)</label>
          <input
            type="number"
            id="timeout"
            bind:value={timeout}
            min="1"
            max="300"
          />
        </div>

        <div class="form-group">
          <label for="retries">Retry Attempts</label>
          <input
            type="number"
            id="retries"
            bind:value={retries}
            min="0"
            max="10"
          />
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2>Alerting</h2>
      
      <div class="form-group">
        <label class="checkbox-label">
          <input
            type="checkbox"
            bind:checked={alertOnFailure}
          />
          Send alerts on failure
        </label>
      </div>

      {#if alertOnFailure}
        <div class="form-group">
          <label for="alertEmail">Alert Email</label>
          <input
            type="email"
            id="alertEmail"
            bind:value={alertEmail}
            placeholder="alerts@example.com"
          />
        </div>
      {/if}
    </div>

    <div class="form-actions">
      <button type="button" on:click={resetForm} class="btn-secondary">Reset</button>
      <button type="submit" class="btn-primary">Create Job</button>
    </div>
  </form>
</div>

<style>
  .job-config {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }

  h1 {
    color: #333;
    text-align: center;
    margin-bottom: 2rem;
  }

  .form-section {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  }

  .form-section h2 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1.2rem;
    border-bottom: 2px solid #f0f0f0;
    padding-bottom: 0.5rem;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #555;
    font-weight: 500;
  }

  input[type="text"],
  input[type="email"],
  input[type="number"],
  select,
  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.2s;
  }

  input:focus,
  select:focus,
  textarea:focus {
    outline: none;
    border-color: #ff3e00;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    width: auto;
    margin-right: 0.5rem;
  }

  .help-text {
    font-size: 0.8rem;
    color: #666;
    margin-top: 0.25rem;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .btn-primary {
    background: #ff3e00;
    color: white;
  }

  .btn-primary:hover {
    background: #e63200;
  }

  .btn-secondary {
    background: #f5f5f5;
    color: #333;
  }

  .btn-secondary:hover {
    background: #e0e0e0;
  }

  @media (prefers-color-scheme: dark) {
    .job-config {
      color: #ffffff;
    }
    
    h1, .form-section h2 {
      color: #ffffff;
    }
    
    .form-section {
      background: #2a2a2a;
      color: #ffffff;
    }
    
    .form-section h2 {
      border-bottom-color: #444;
    }
    
    label, .help-text {
      color: #cccccc;
    }
    
    input, select, textarea {
      background: #333;
      color: #ffffff;
      border-color: #555;
    }
    
    input:focus, select:focus, textarea:focus {
      border-color: #ff3e00;
    }
    
    .btn-secondary {
      background: #444;
      color: #ffffff;
    }
    
    .btn-secondary:hover {
      background: #555;
    }
  }

  @media (max-width: 768px) {
    .form-row {
      grid-template-columns: 1fr;
    }
    
    .form-actions {
      flex-direction: column;
    }
  }
</style>
