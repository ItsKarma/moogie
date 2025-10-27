<script>
  import { onMount } from 'svelte';
  import Router, { link, location } from 'svelte-spa-router';
  import Dashboard from './pages/Dashboard.svelte';
  import Jobs from './pages/Jobs.svelte';
  import Settings from './pages/Settings.svelte';
  import DateRangePicker from './components/DateRangePicker.svelte';
  import { dateRange } from './lib/stores.js';

  // Define routes
  const routes = {
    '/': Dashboard,
    '/dashboard': Dashboard,
    '/jobs': Jobs,
    '/job/:jobId': Jobs,
    '/settings': Settings,
  };

  // Use the reactive location store from svelte-spa-router
  $: currentPath = $location;

  // Initialize date range from URL parameters when app starts
  onMount(() => {
    dateRange.initFromURL();
  });
</script>

<div class="app">
  <nav class="navbar">
    <div class="nav-container">
      <div class="nav-brand">
        <img src="/moogie.svg" alt="Moogie" class="nav-logo" />
        <h1>Moogie</h1>
      </div>
      
      <div class="nav-right">
        <ul class="nav-links">
          <li>
            <a 
              href="/" 
              use:link 
              class:active={currentPath === '/' || currentPath === '/dashboard'}
            >
              Dashboard
            </a>
          </li>
          <li>
            <a 
              href="/jobs" 
              use:link 
              class:active={currentPath === '/jobs' || currentPath.startsWith('/job/')}
            >
              Jobs
            </a>
          </li>
        </ul>
        
        <div class="nav-actions">
          <DateRangePicker />
          <a 
            href="/settings" 
            use:link 
            class="settings-icon"
            class:active={currentPath === '/settings'}
            title="Settings"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path>
              <circle cx="12" cy="12" r="3"></circle>
            </svg>
          </a>
        </div>
      </div>
    </div>
  </nav>

  <main class="main-content">
    <Router {routes} />
  </main>
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  .navbar {
    background: #7D471F;
    color: white;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .nav-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 var(--spacing-lg);
    width: 100%;
  }

  .nav-brand {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .nav-logo {
    width: 32px;
    height: 32px;
  }

  .nav-brand h1 {
    margin: 0;
    color: white;
    font-size: 1.5rem;
  }

  .nav-right {
    display: flex;
    align-items: center;
    gap: var(--spacing-lg);
  }

  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: var(--spacing-sm);
  }

  .nav-actions {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
  }

  .nav-links li {
    margin: 0;
  }

  .nav-links a {
    display: block;
    padding: 0.75rem 1.5rem;
    color: white;
    text-decoration: none;
    border-radius: 4px;
    transition: background-color 0.2s;
    font-weight: 500;
  }

  .nav-links a:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .nav-links a.active {
    background-color: rgba(255, 255, 255, 0.2);
    font-weight: 600;
  }

  .settings-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--spacing-sm);
    color: white;
    text-decoration: none;
    border-radius: 4px;
    transition: background-color 0.2s;
    cursor: pointer;
  }

  .settings-icon:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .settings-icon.active {
    background-color: rgba(255, 255, 255, 0.2);
  }

  .main-content {
    flex: 1;
    background: var(--background-color);
    min-height: calc(100vh - 80px);
  }

  /* Mobile responsiveness */
  @media (max-width: 768px) {
    .nav-container {
      flex-direction: column;
      gap: 1rem;
      padding: 1rem;
    }

    .nav-brand {
      text-align: center;
    }

    .nav-right {
      flex-direction: column;
      gap: 1rem;
      width: 100%;
    }

    .nav-links {
      flex-wrap: wrap;
      justify-content: center;
      gap: 0.5rem;
    }

    .nav-links a {
      padding: 0.5rem 1rem;
      font-size: 0.9rem;
    }

    .nav-actions {
      width: 100%;
      justify-content: center;
    }

    .main-content {
      min-height: calc(100vh - 120px);
    }
  }

  @media (max-width: 480px) {
    .nav-links {
      flex-direction: column;
      width: 100%;
    }

    .nav-links a {
      text-align: center;
    }
  }
</style>
