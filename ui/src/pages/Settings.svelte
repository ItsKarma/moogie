<script>
  import { themeStore, THEME_OPTIONS } from '../lib/themeStore.js';
  
  let currentTheme = 'system';
  let effectiveTheme = 'dark';
  
  themeStore.subscribe(value => {
    currentTheme = value.theme;
    effectiveTheme = value.effectiveTheme;
  });
  
  function handleThemeChange(theme) {
    themeStore.setTheme(theme);
  }
</script>

<div class="settings-container">
  <div class="settings-header">
    <h1>Settings</h1>
    <p class="subtitle">Configure your Moogie monitoring dashboard preferences</p>
  </div>

  <div class="settings-content">
    <!-- Appearance Section -->
    <section class="settings-section">
      <h2>Appearance</h2>
      
      <div class="setting-item">
        <div class="setting-info">
          <h3>Theme</h3>
          <p>Choose how Moogie looks. Select a single theme, or sync with your system and automatically switch between light and dark themes.</p>
          {#if currentTheme === THEME_OPTIONS.SYSTEM}
            <p class="theme-hint">
              Currently using: <strong>{effectiveTheme === 'light' ? 'Light' : 'Dark'}</strong> (from system)
            </p>
          {/if}
        </div>
        
        <div class="theme-toggle-group">
          <button 
            class="theme-toggle-option {currentTheme === THEME_OPTIONS.LIGHT ? 'active' : ''}"
            on:click={() => handleThemeChange(THEME_OPTIONS.LIGHT)}
          >
            Light
          </button>
          
          <button 
            class="theme-toggle-option {currentTheme === THEME_OPTIONS.DARK ? 'active' : ''}"
            on:click={() => handleThemeChange(THEME_OPTIONS.DARK)}
          >
            Dark
          </button>
          
          <button 
            class="theme-toggle-option {currentTheme === THEME_OPTIONS.SYSTEM ? 'active' : ''}"
            on:click={() => handleThemeChange(THEME_OPTIONS.SYSTEM)}
          >
            System
          </button>
        </div>
      </div>
    </section>

    <!-- Future sections can go here -->
    <section class="settings-section coming-soon">
      <h2>Notifications</h2>
      <p class="coming-soon-text">Coming soon: Configure alert preferences and notification settings</p>
    </section>
    
    <section class="settings-section coming-soon">
      <h2>Data & Refresh</h2>
      <p class="coming-soon-text">Coming soon: Set data refresh intervals and API preferences</p>
    </section>
  </div>
</div>

<style>
  .settings-container {
    max-width: 900px;
    margin: 0 auto;
    padding: var(--spacing-xl);
  }

  .settings-header {
    margin-bottom: var(--spacing-xl);
  }

  .settings-header h1 {
    margin: 0 0 var(--spacing-sm) 0;
    font-size: var(--font-xxl);
    color: var(--text-primary);
  }

  .subtitle {
    color: var(--text-secondary);
    font-size: var(--font-base);
    margin: 0;
  }

  .settings-content {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xl);
  }

  .settings-section {
    background: var(--card-bg);
    border-radius: var(--radius-md);
    padding: var(--spacing-xl);
    box-shadow: var(--shadow-lg);
  }

  .settings-section h2 {
    margin: 0 0 var(--spacing-lg) 0;
    font-size: var(--font-xl);
    color: var(--text-primary);
    border-bottom: 2px solid var(--border-color);
    padding-bottom: var(--spacing-sm);
  }

  .setting-item {
    margin-bottom: var(--spacing-lg);
  }

  .setting-item:last-child {
    margin-bottom: 0;
  }

  .setting-info h3 {
    margin: 0 0 var(--spacing-sm) 0;
    font-size: var(--font-lg);
    color: var(--text-primary);
  }

  .setting-info p {
    margin: 0 0 var(--spacing-lg) 0;
    color: var(--text-secondary);
    font-size: var(--font-md);
    line-height: 1.6;
  }

  .theme-toggle-group {
    display: inline-flex;
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    padding: 4px;
    gap: 4px;
  }

  .theme-toggle-option {
    padding: var(--spacing-sm) var(--spacing-xl);
    border: none;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--text-primary);
    font-size: var(--font-md);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .theme-toggle-option:hover {
    background: var(--hover-bg);
  }

  .theme-toggle-option.active {
    background: var(--primary-color);
    color: white;
    box-shadow: var(--shadow-sm);
  }

  .theme-hint {
    color: var(--text-secondary);
    font-size: var(--font-sm);
    margin: var(--spacing-md) 0 0 0;
    font-style: italic;
  }

  .theme-hint strong {
    color: var(--primary-color);
  }

  .coming-soon {
    opacity: 0.6;
  }

  .coming-soon-text {
    color: var(--text-secondary);
    font-size: var(--font-md);
    font-style: italic;
  }
</style>
