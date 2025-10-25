// API service for communicating with the Moogie backend
class ApiService {
  constructor() {
    this.baseURL = "http://localhost:8080/api/v1";
  }

  // Helper method to make API requests
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;

    const config = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error);
      throw error;
    }
  }

  // Jobs API
  async getJobs(from, to) {
    const params = new URLSearchParams();
    if (from) params.append("from", from);
    if (to) params.append("to", to);

    const endpoint = `/jobs${params.toString() ? `?${params.toString()}` : ""}`;
    return this.request(endpoint);
  }

  async getJob(id, from, to, limit = 100) {
    const params = new URLSearchParams();
    if (from) params.append("from", from);
    if (to) params.append("to", to);
    if (limit) params.append("limit", limit.toString());

    const endpoint = `/jobs/${id}${
      params.toString() ? `?${params.toString()}` : ""
    }`;
    return this.request(endpoint);
  }

  // Dashboard API
  async getDashboardSummary(from, to) {
    const params = new URLSearchParams();
    if (from) params.append("from", from);
    if (to) params.append("to", to);

    const endpoint = `/dashboard/summary${
      params.toString() ? `?${params.toString()}` : ""
    }`;
    return this.request(endpoint);
  }

  // Executions API
  async createExecution(executionData) {
    return this.request("/executions", {
      method: "POST",
      body: JSON.stringify(executionData),
    });
  }

  // Health check
  async healthCheck() {
    return this.request("/health");
  }
}

// Create and export a singleton instance
export const apiService = new ApiService();

// Export the class for testing or custom instances
export default ApiService;
