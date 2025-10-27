<script>
  import { Line } from 'svelte-chartjs';
  import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler
  } from 'chart.js';
  import { getSuccessRateColor } from '../../lib/utils.js';

  // Register Chart.js components
  ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler
  );

  export let selectedJob = null;
  
  $: job = selectedJob;

  // Prepare chart data from all executions within date range
  $: chartData = job?.executions && job.executions.length > 0 ? {
    labels: job.executions
      .slice()
      .reverse()
      .map(exec => {
        const date = new Date(exec.timestamp);
        return date.toLocaleString('en-US', { 
          month: 'short',
          day: 'numeric',
          hour: '2-digit', 
          minute: '2-digit'
        });
      }),
    datasets: [{
      label: 'Response Time',
      data: job.executions
        .slice()
        .reverse()
        .map(exec => exec.response_time),
      borderColor: 'rgba(161, 98, 7, 1)', // Warm brown
      backgroundColor: 'rgba(161, 98, 7, 0.08)',
      tension: 0.3,
      fill: true,
      pointRadius: 5,
      pointHoverRadius: 7,
      pointBackgroundColor: job.executions
        .slice()
        .reverse()
        .map(exec => {
          if (exec.status === 'success') return 'rgba(34, 197, 94, 1)';
          if (exec.status === 'warning') return 'rgba(234, 179, 8, 1)';
          return 'rgba(239, 68, 68, 1)';
        }),
      pointBorderColor: 'rgba(255, 255, 255, 1)',
      pointBorderWidth: 2,
      borderWidth: 2
    }]
  } : null;

  $: chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: 'index'
    },
    plugins: {
      legend: {
        display: false
      },
      tooltip: {
        backgroundColor: 'rgba(17, 24, 39, 0.95)',
        padding: 12,
        titleColor: 'rgba(255, 255, 255, 0.95)',
        bodyColor: 'rgba(255, 255, 255, 0.85)',
        borderColor: 'rgba(161, 98, 7, 0.3)',
        borderWidth: 1,
        cornerRadius: 8,
        displayColors: false,
        callbacks: {
          title: (items) => {
            return items[0].label;
          },
          label: (context) => {
            return `Response Time: ${context.parsed.y}ms`;
          }
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        border: {
          display: false
        },
        grid: {
          color: 'rgba(128, 128, 128, 0.08)',
          lineWidth: 1
        },
        ticks: {
          color: 'rgba(107, 114, 128, 0.8)',
          font: {
            size: 11,
            family: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif"
          },
          padding: 8,
          callback: (value) => `${value}ms`
        }
      },
      x: {
        border: {
          display: false
        },
        grid: {
          display: false
        },
        ticks: {
          color: 'rgba(107, 114, 128, 0.8)',
          font: {
            size: 11,
            family: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif"
          },
          maxRotation: 45,
          minRotation: 45,
          padding: 4
        }
      }
    }
  };

  function timeAgo(dateString) {
    if (!dateString) return 'Never';
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);
    
    if (seconds < 60) return 'Just now';
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    return `${Math.floor(seconds / 86400)}d ago`;
  }
</script>

<div class="detail-panel">
  {#if job}
    <div class="detail-header">
      <div class="detail-title">
        <h2>{job.name}</h2>
        <span class="detail-type">{job.type}</span>
      </div>
      <div class="detail-stats">
        <div class="detail-stat">
          <span class="detail-stat-label">Success Rate</span>
          <span class="detail-stat-value" style="color: {getSuccessRateColor(job.success_rate)}">
            {job.success_rate.toFixed(1)}%
          </span>
        </div>
        <div class="detail-stat">
          <span class="detail-stat-label">Executions</span>
          <span class="detail-stat-value">{job.execution_count?.toLocaleString() || 0}</span>
        </div>
        <div class="detail-stat">
          <span class="detail-stat-label">Last Run</span>
          <span class="detail-stat-value">{timeAgo(job.last_execution)}</span>
        </div>
        {#if job.avg_response_time}
          <div class="detail-stat">
            <span class="detail-stat-label">Avg Response</span>
            <span class="detail-stat-value">{Math.round(job.avg_response_time)}ms</span>
          </div>
        {/if}
      </div>
    </div>

    <!-- Response Time Graph -->
    <div class="graph-container">
      <h3>Response Time History</h3>
      {#if chartData && chartData.datasets[0].data.length > 0}
        <div class="chart-wrapper">
          <Line data={chartData} options={chartOptions} />
        </div>
      {:else}
        <div class="graph-placeholder">
          <p>No execution data available yet</p>
        </div>
      {/if}
    </div>

    <!-- Recent Executions -->
    {#if job.recent_executions && job.recent_executions.length > 0}
      <div class="recent-executions">
        <h3>Recent Executions</h3>
        <div class="execution-list">
          {#each job.recent_executions as execution}
            <div class="execution-row">
              <div class="execution-status">
                <div class="pill pill-{execution.status}"></div>
                <span class="execution-status-text">{execution.status}</span>
              </div>
              <div class="execution-time">{new Date(execution.timestamp).toLocaleString()}</div>
              {#if execution.response_time}
                <div class="execution-response">{execution.response_time}ms</div>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {:else}
    <div class="detail-empty">
      <p>Select a check to view details</p>
    </div>
  {/if}
</div>

<style>
  .detail-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
    min-width: 0;
  }

  .detail-header {
    background: var(--card-bg);
    padding: var(--spacing-lg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
  }

  .detail-title {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);
  }

  .detail-title h2 {
    margin: 0;
    font-size: var(--font-xl);
    color: var(--text-primary);
  }

  .detail-type {
    padding: var(--spacing-xs) var(--spacing-sm);
    background: var(--bg-secondary);
    border-radius: var(--radius-sm);
    font-size: var(--font-xs);
    color: var(--text-secondary);
    text-transform: uppercase;
    font-weight: 500;
    letter-spacing: 0.05em;
  }

  .detail-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: var(--spacing-md);
  }

  .detail-stat {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);
  }

  .detail-stat-label {
    font-size: var(--font-xs);
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-weight: 500;
  }

  .detail-stat-value {
    font-size: var(--font-lg);
    font-weight: 600;
    color: var(--text-primary);
  }

  /* Graph Container */
  .graph-container {
    background: var(--card-bg);
    padding: var(--spacing-lg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
  }

  .graph-container h3 {
    margin: 0 0 var(--spacing-md) 0;
    font-size: var(--font-lg);
    color: var(--text-primary);
  }

  .chart-wrapper {
    height: 300px;
    position: relative;
  }

  .graph-placeholder {
    background: var(--bg-secondary);
    border: 2px dashed var(--border-color);
    border-radius: var(--radius-md);
    padding: var(--spacing-xl);
    text-align: center;
    min-height: 300px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .graph-placeholder p {
    margin: 0;
    color: var(--text-secondary);
    font-size: var(--font-md);
  }

  /* Recent Executions */
  .recent-executions {
    background: var(--card-bg);
    padding: var(--spacing-lg);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
  }

  .recent-executions h3 {
    margin: 0 0 var(--spacing-md) 0;
    font-size: var(--font-lg);
    color: var(--text-primary);
  }

  .execution-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);
  }

  .execution-row {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
    padding: var(--spacing-sm);
    background: var(--bg-secondary);
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
  }

  .execution-status {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    min-width: 100px;
  }

  .execution-status-text {
    text-transform: capitalize;
    font-weight: 500;
  }

  .execution-time {
    flex: 1;
    color: var(--text-secondary);
  }

  .execution-response {
    font-weight: 600;
    color: var(--text-primary);
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

  .detail-empty {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--card-bg);
    border-radius: var(--radius-lg);
    padding: var(--spacing-xl);
    color: var(--text-secondary);
  }

  .detail-empty p {
    font-size: var(--font-lg);
    margin: 0;
  }
</style>
