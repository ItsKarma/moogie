<script>
  import { dateRange } from '../lib/stores.js';

  let showPicker = false;
  let tempFromDate = '';
  let tempToDate = '';
  let currentFromDate = '';
  let currentToDate = '';

  // Subscribe to store changes
  dateRange.subscribe(range => {
    currentFromDate = range.from;
    currentToDate = range.to;
    // If picker is closed, sync temp values
    if (!showPicker) {
      tempFromDate = range.from;
      tempToDate = range.to;
    }
  });

  function togglePicker() {
    if (!showPicker) {
      // Opening picker - sync temp values with current
      tempFromDate = currentFromDate;
      tempToDate = currentToDate;
    }
    showPicker = !showPicker;
  }

  function handleFromDateChange(event) {
    const dateValue = event.target.value;
    if (dateValue) {
      tempFromDate = new Date(dateValue).toISOString();
    }
  }

  function handleToDateChange(event) {
    const dateValue = event.target.value;
    if (dateValue) {
      tempToDate = new Date(dateValue).toISOString();
    }
  }

  function setQuickRange(days) {
    const to = new Date();
    const from = new Date();
    from.setDate(from.getDate() - days);
    
    // Quick ranges use live "to" date (always "now")
    dateRange.setRange(from.toISOString(), to.toISOString(), true);
    showPicker = false;
  }

  function setQuickRangeHours(hours) {
    const to = new Date();
    const from = new Date();
    from.setHours(from.getHours() - hours);
    
    // Quick ranges use live "to" date (always "now")
    dateRange.setRange(from.toISOString(), to.toISOString(), true);
    showPicker = false;
  }

  function applyDateRange() {
    // Validate dates
    if (!tempFromDate || !tempToDate) {
      alert('Please select both start and end dates.');
      return;
    }

    const fromDateObj = new Date(tempFromDate);
    const toDateObj = new Date(tempToDate);

    if (fromDateObj > toDateObj) {
      alert('Start date cannot be after end date.');
      return;
    }

    if (toDateObj > new Date()) {
      alert('End date cannot be in the future.');
      return;
    }

    // Custom range uses locked "to" date (user-specified)
    dateRange.setRange(tempFromDate, tempToDate, false);
    showPicker = false;
  }

  function cancelDateRange() {
    // Reset temp values to current values
    tempFromDate = currentFromDate;
    tempToDate = currentToDate;
    showPicker = false;
  }

  function handleOutsideClick(event) {
    if (!event.target.closest('.date-picker-container')) {
      // Reset temp values when closing without applying
      tempFromDate = currentFromDate;
      tempToDate = currentToDate;
      showPicker = false;
    }
  }

  function formatDateForDisplay(dateString) {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  // Convert ISO string to datetime-local format (YYYY-MM-DDTHH:MM)
  function toDateTimeLocal(isoString) {
    if (!isoString) return '';
    const date = new Date(isoString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  }

  // Check if temp values are different from current values
  $: hasChanges = tempFromDate !== currentFromDate || tempToDate !== currentToDate;
</script>

<svelte:window on:click={handleOutsideClick} />

<div class="date-picker-container">
  <button class="date-trigger" on:click={togglePicker}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
      <line x1="16" y1="2" x2="16" y2="6"></line>
      <line x1="8" y1="2" x2="8" y2="6"></line>
      <line x1="3" y1="10" x2="21" y2="10"></line>
    </svg>
    <span class="date-display">
      {formatDateForDisplay(currentFromDate)} - {formatDateForDisplay(currentToDate)}
    </span>
    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron" class:open={showPicker}>
      <polyline points="6,9 12,15 18,9"></polyline>
    </svg>
  </button>

  {#if showPicker}
    <div class="date-picker-dropdown">
      <div class="quick-ranges">
        <h4>Quick Ranges</h4>
        <div class="quick-buttons">
          <button type="button" on:click={() => setQuickRangeHours(1)}>Last 1 hour</button>
          <button type="button" on:click={() => setQuickRange(1)}>Last 1 day</button>
          <button type="button" on:click={() => setQuickRangeHours(6)}>Last 6 hours</button>
          <button type="button" on:click={() => setQuickRange(3)}>Last 3 days</button>
          <button type="button" on:click={() => setQuickRangeHours(12)}>Last 12 hours</button>
          <button type="button" on:click={() => setQuickRange(7)}>Last 7 days</button>
        </div>
      </div>
      
      <div class="custom-range">
        <h4>Custom Range</h4>
        <div class="date-inputs">
          <div class="input-group">
            <label for="from-date">From</label>
            <input 
              id="from-date"
              type="datetime-local" 
              value={toDateTimeLocal(tempFromDate)}
              on:change={handleFromDateChange}
              max={toDateTimeLocal(tempToDate)}
            />
          </div>
          <div class="input-group">
            <label for="to-date">To</label>
            <input 
              id="to-date"
              type="datetime-local" 
              value={toDateTimeLocal(tempToDate)}
              on:change={handleToDateChange}
              min={toDateTimeLocal(tempFromDate)}
              max={toDateTimeLocal(new Date().toISOString())}
            />
          </div>
        </div>
        
        <div class="picker-actions">
          <button type="button" class="cancel-btn" on:click={cancelDateRange}>
            Cancel
          </button>
          <button 
            type="button" 
            class="apply-btn" 
            on:click={applyDateRange}
            disabled={!hasChanges}
          >
            Apply
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .date-picker-container {
    position: relative;
  }

  .date-trigger {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm) var(--spacing-md);
    background: var(--card-background);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.9rem;
    color: var(--text-color);
    transition: all 0.2s ease;
  }

  .date-trigger:hover {
    background: var(--table-hover);
    border-color: var(--primary-color);
  }

  .date-display {
    white-space: nowrap;
  }

  .chevron {
    transition: transform 0.2s ease;
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .date-picker-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    z-index: 1000;
    background: var(--card-background);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
    padding: var(--spacing-lg);
    min-width: 300px;
    margin-top: var(--spacing-xs);
  }

  .quick-ranges, .custom-range {
    margin-bottom: var(--spacing-lg);
  }

  .custom-range:last-child {
    margin-bottom: 0;
  }

  h4 {
    margin: 0 0 var(--spacing-sm) 0;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--text-secondary);
    text-transform: uppercase;
  }

  .quick-buttons {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-xs);
  }

  .quick-buttons button {
    padding: var(--spacing-sm) var(--spacing-md);
    background: var(--table-hover);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.85rem;
    color: var(--text-color);
    transition: all 0.2s ease;
    text-align: center;
    white-space: nowrap;
  }

  .quick-buttons button:hover {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .date-inputs {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);
  }

  label {
    font-size: 0.8rem;
    font-weight: 500;
    color: var(--text-secondary);
  }

  input[type="date"],
  input[type="datetime-local"] {
    padding: var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: 0.9rem;
    color: var(--text-color);
    background: var(--card-background);
    transition: border-color 0.2s ease;
  }

  input[type="date"]:focus,
  input[type="datetime-local"]:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  /* Fix for datetime-local in dark mode */
  input[type="datetime-local"]::-webkit-calendar-picker-indicator {
    filter: invert(var(--invert-icon, 0));
  }

  :root.theme-dark input[type="datetime-local"]::-webkit-calendar-picker-indicator,
  :root[data-theme="dark"] input[type="datetime-local"]::-webkit-calendar-picker-indicator {
    filter: invert(1);
  }

  .picker-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-lg);
    padding-top: var(--spacing-md);
    border-top: 1px solid var(--border-color);
  }

  .cancel-btn, .apply-btn {
    padding: var(--spacing-sm) var(--spacing-lg);
    border-radius: var(--radius-sm);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 1px solid;
  }

  .cancel-btn {
    background: transparent;
    color: var(--text-secondary);
    border-color: var(--border-color);
  }

  .cancel-btn:hover {
    background: var(--table-hover);
    color: var(--text-color);
  }

  .apply-btn {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .apply-btn:hover:not(:disabled) {
    background: var(--primary-dark);
    border-color: var(--primary-dark);
  }

  .apply-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: var(--text-muted);
    border-color: var(--text-muted);
  }

  /* Dark mode adjustments */
  @media (prefers-color-scheme: dark) {
    input[type="date"]::-webkit-calendar-picker-indicator {
      filter: invert(1);
    }
  }

  /* Mobile responsiveness */
  @media (max-width: 768px) {
    .date-picker-dropdown {
      right: -1rem;
      left: -1rem;
      min-width: auto;
    }

    .date-inputs {
      grid-template-columns: 1fr;
    }

    .date-display {
      display: none;
    }
  }
</style>
