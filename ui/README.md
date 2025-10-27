# Moogie UI

Modern, real-time monitoring dashboard built with Svelte and Vite.

## Features

- 📊 **Real-time Dashboard** - Split-view interface with live job monitoring
- 📈 **Response Time Graphs** - Chart.js powered time series visualization
- 🎨 **Theme Support** - Light, Dark, and System preference modes
- 📅 **Flexible Date Ranges** - Quick ranges (1h-7d) and custom date/time selection
- 🔍 **Smart Filtering** - Filter jobs by status (failed, warning, success)
- 📱 **Responsive Design** - Optimized for desktop and mobile
- 🚀 **Fast & Lightweight** - Vite-powered dev server and optimized builds

## Architecture

### Component Structure

```
ui/src/
├── components/
│   ├── dashboard/           # Dashboard-specific components
│   │   ├── ChecksListSidebar.svelte   # Job list with search/filter
│   │   ├── DashboardStats.svelte      # Summary statistics bar
│   │   └── JobDetailPanel.svelte      # Job details with graph
│   ├── job-detail/          # Job detail page components
│   │   ├── ExecutionHistoryTable.svelte
│   │   ├── JobConfigSection.svelte
│   │   └── JobSummaryCards.svelte
│   ├── DateRangePicker.svelte   # Shared date range selector
│   └── StatusBadge.svelte       # Shared status indicator
├── pages/
│   ├── Dashboard.svelte     # Main split-view dashboard
│   ├── Jobs.svelte          # Job list view
│   ├── JobDetail.svelte     # Individual job details
│   └── Settings.svelte      # Theme and preferences
├── lib/
│   ├── api.js              # API service layer
│   ├── stores.js           # Svelte stores (jobs, dashboard, date range)
│   ├── themeStore.js       # Theme management
│   ├── utils.js            # Shared utility functions
│   └── mockData.js         # Sample data for development
├── app.css                 # Global styles and CSS variables
├── App.svelte             # Root component with routing
└── main.js                # App entry point
```

### State Management

**Stores:**

- `dateRange` - Global date range filter with URL sync
- `jobsStore` - Job list with loading/error states
- `dashboardStore` - Dashboard summary data
- `jobDetailStore` - Individual job details with executions
- `themeStore` - Theme preference (light/dark/system)

**API Service:**

- Centralized API client (`apiService`)
- Automatic error handling and retries
- ISO 8601 timestamp support

### Styling

**CSS Architecture:**

- Custom properties (CSS variables) for theming
- Dark mode with class-based and system preference support
- Responsive design with mobile-first approach
- Consistent spacing, typography, and color scales

**Design Tokens:**

```css
--primary-color, --status-*, --success-*
--spacing-xs through --spacing-xxl
--font-xs through --font-xxl
--radius-sm, --radius-md, --radius-lg
```

## Development

### Prerequisites

- Node.js 18+ and npm
- API server running on `localhost:8080`

### Setup

1. **Install dependencies:**

   ```bash
   cd ui
   npm install
   ```

2. **Configure API endpoint:**
   Update `VITE_API_URL` in `src/lib/api.js` if needed (defaults to `http://localhost:8080`)

3. **Start development server:**

   ```bash
   npm run dev
   ```

   Dashboard available at: http://localhost:3000

### Build for Production

```bash
npm run build
```

Production files output to `dist/` directory.

Preview production build:

```bash
npm run preview
```

## Docker

### Development Mode

```bash
# From project root
docker-compose up ui
```

Hot-reload enabled with volume mounts.

### Production Build

```bash
# Build production image
docker build -t moogie-ui .

# Run with nginx
docker run -p 3000:80 moogie-ui
```

## Key Features Explained

### Split-View Dashboard

The main dashboard (`pages/Dashboard.svelte`) uses a split layout:

- **Left (30%)**: Job list with search, filters, and status pills
- **Right (70%)**: Selected job details with response time graph

Jobs with failures automatically appear at the top. Status pills show the last 5 executions with color coding.

### Response Time Graph

Powered by Chart.js (`JobDetailPanel.svelte`):

- Adaptive time intervals based on date range (5min, 15min, 1h, 6h)
- Shows full time range with gaps for missing data
- Color-coded status points (green/yellow/red)
- Limited axis labels for clean visualization (max 7 ticks)

### Date Range Filtering

`DateRangePicker.svelte` provides:

- Quick ranges: 1h, 6h, 12h, 1d, 3d, 7d (3x2 grid)
- Custom date/time selection with ISO 8601 timestamps
- URL persistence (shareable filtered views)
- Real-time updates across all components

### Theme System

Three-way theme toggle (`Settings.svelte` + `themeStore.js`):

- **Light**: Explicit light theme
- **Dark**: Explicit dark theme
- **System**: Follows OS preference

Persisted in localStorage, applied via CSS classes on `<html>` root.

## API Integration

### Endpoints Used

```javascript
GET  /api/v1/jobs?from={ISO8601}&to={ISO8601}
GET  /api/v1/jobs/:id?from={ISO8601}&to={ISO8601}&limit={number}
GET  /api/v1/dashboard/summary?from={ISO8601}&to={ISO8601}
WS   /ws  (WebSocket for real-time updates)
```

### Data Flow

1. `stores.js` fetches data via `api.js`
2. Components subscribe to stores reactively
3. Date range changes trigger automatic refetches
4. WebSocket updates push new executions in real-time (planned)

## Testing

### Manual Testing

Use the seed data script to populate realistic test data:

```bash
# From project root
docker-compose run --rm seeder
```

This creates 15 jobs with 180 days of execution history (14,000+ records).

## Common Tasks

### Adding a New Component

1. Create `.svelte` file in appropriate directory
2. Import shared utilities from `lib/utils.js`
3. Use CSS variables from `app.css`
4. Subscribe to stores if needed for reactive data

### Adding a New Page

1. Create page component in `src/pages/`
2. Add route in `App.svelte` (using svelte-spa-router)
3. Add navigation link in header if needed

### Modifying the Graph

Chart configuration in `JobDetailPanel.svelte`:

- Data generation: `chartData` reactive statement
- Chart options: `chartOptions` object
- Chart.js docs: https://www.chartjs.org/

### Updating Styles

1. **Global changes**: Edit `app.css` CSS variables
2. **Component-specific**: Add `<style>` block in `.svelte` file
3. **Theme-aware**: Use CSS variables, avoid hardcoded colors

## Performance

- **Code splitting**: Automatic with Vite
- **Lazy loading**: Router-based page loading
- **Optimized builds**: Minification, tree-shaking
- **Efficient rendering**: Svelte compiler optimizations
- **Limited data**: API limits execution history (default 100-1000)

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## Troubleshooting

### API Connection Issues

Check:

1. API server is running (`docker-compose logs api`)
2. CORS is configured for UI origin
3. `VITE_API_URL` matches API server address

### Build Errors

```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Hot Reload Not Working

- Ensure volume mounts are correct in `docker-compose.yaml`
- Check file permissions on mounted volumes
- Restart development server

## Resources

- [Svelte Documentation](https://svelte.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [Chart.js Documentation](https://www.chartjs.org/)
- [Svelte SPA Router](https://github.com/ItalyPaleAle/svelte-spa-router)
