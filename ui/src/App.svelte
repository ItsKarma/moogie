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
          <li>
            <a 
              href="/settings" 
              use:link 
              class:active={currentPath === '/settings'}
            >
              Settings
            </a>
          </li>
        </ul>
        
        <div class="nav-actions">
          <DateRangePicker />
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
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 1rem;
  }

  .nav-brand h1 {
    margin: 0;
    color: white;
    font-size: 1.5rem;
  }

  .nav-right {
    display: flex;
    align-items: center;
    gap: 2rem;
  }

  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 1rem;
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

  .main-content {
    flex: 1;
    background: #f5f5f5;
    min-height: calc(100vh - 80px);
  }

  /* Dark mode support */
  @media (prefers-color-scheme: dark) {
    .main-content {
      background: #1a1a1a;
    }
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
