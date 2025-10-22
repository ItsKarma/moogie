<script>
  import Router, { link, location } from 'svelte-spa-router';
  import Dashboard from './pages/Dashboard.svelte';
  import Jobs from './pages/Jobs.svelte';

  // Define routes
  const routes = {
    '/': Dashboard,
    '/dashboard': Dashboard,
    '/jobs': Jobs,
    '/job/:jobId': Jobs,
  };

  // Use the reactive location store from svelte-spa-router
  $: currentPath = $location;
</script>

<div class="app">
  <nav class="navbar">
    <div class="nav-container">
      <div class="nav-brand">
        <h1>Moogie</h1>
      </div>
      
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
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .nav-brand h1 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: bold;
  }

  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 0;
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

    .nav-links {
      flex-wrap: wrap;
      justify-content: center;
      gap: 0.5rem;
    }

    .nav-links a {
      padding: 0.5rem 1rem;
      font-size: 0.9rem;
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
