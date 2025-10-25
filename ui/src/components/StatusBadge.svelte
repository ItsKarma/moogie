<script>
  import { getStatusColor } from '../lib/utils.js';
  
  export let status;
  export let size = 'normal'; // 'small', 'normal', 'large'
</script>

<span 
  class="status-badge {size}" 
  style="background-color: {getStatusColor(status)}"
>
  {status.toUpperCase()}
</span>

<style>
  .status-badge {
    padding: var(--spacing-xs) var(--spacing-sm);
    border-radius: var(--radius-lg);
    color: white;
    font-size: var(--font-sm);
    font-weight: bold;
    text-align: center;
    display: inline-block;
  }

  .status-badge.large {
    padding: var(--spacing-sm) var(--spacing-md);
    font-size: var(--font-base);
  }

  .status-badge.small {
    padding: 0.2rem var(--spacing-sm);
    font-size: var(--font-xs);
  }
</style>
