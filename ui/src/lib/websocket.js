// WebSocket service for real-time updates
import { writable } from "svelte/store";

/**
 * WebSocket connection states
 */
export const ConnectionState = {
  CONNECTING: "connecting",
  CONNECTED: "connected",
  DISCONNECTED: "disconnected",
  ERROR: "error",
};

/**
 * WebSocket message types from backend
 */
export const MessageType = {
  EXECUTION_CREATED: "execution_created",
  JOB_UPDATED: "job_updated",
  DASHBOARD_UPDATED: "dashboard_updated",
};

/**
 * WebSocket service for managing real-time connections
 */
class WebSocketService {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000; // Start with 1 second
    this.maxReconnectDelay = 30000; // Max 30 seconds
    this.reconnectTimer = null;
    this.messageHandlers = new Map();
    this.isIntentionallyClosed = false;

    // Create store for connection status
    this.connectionStatus = writable(ConnectionState.DISCONNECTED);
  }

  /**
   * Connect to WebSocket server
   * @param {string} url - WebSocket server URL
   */
  connect(url = "ws://localhost:8080/ws") {
    if (
      this.ws &&
      (this.ws.readyState === WebSocket.CONNECTING ||
        this.ws.readyState === WebSocket.OPEN)
    ) {
      console.log("WebSocket already connected or connecting");
      return;
    }

    this.isIntentionallyClosed = false;
    this.connectionStatus.set(ConnectionState.CONNECTING);

    try {
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log("WebSocket connected");
        this.connectionStatus.set(ConnectionState.CONNECTED);
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          console.error("Error parsing WebSocket message:", error);
        }
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        this.connectionStatus.set(ConnectionState.ERROR);
      };

      this.ws.onclose = (event) => {
        console.log("WebSocket disconnected:", event.code, event.reason);
        this.connectionStatus.set(ConnectionState.DISCONNECTED);
        this.ws = null;

        // Attempt to reconnect if not intentionally closed
        if (!this.isIntentionallyClosed) {
          this.scheduleReconnect(url);
        }
      };
    } catch (error) {
      console.error("Error creating WebSocket connection:", error);
      this.connectionStatus.set(ConnectionState.ERROR);
      this.scheduleReconnect(url);
    }
  }

  /**
   * Schedule a reconnection attempt with exponential backoff
   * @param {string} url - WebSocket server URL
   */
  scheduleReconnect(url) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error("Max reconnection attempts reached");
      this.connectionStatus.set(ConnectionState.ERROR);
      return;
    }

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }

    const delay = Math.min(
      this.reconnectDelay * Math.pow(2, this.reconnectAttempts),
      this.maxReconnectDelay
    );

    console.log(
      `Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts + 1}/${
        this.maxReconnectAttempts
      })`
    );

    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++;
      this.connect(url);
    }, delay);
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect() {
    this.isIntentionallyClosed = true;

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.connectionStatus.set(ConnectionState.DISCONNECTED);
  }

  /**
   * Handle incoming WebSocket message
   * @param {Object} message - WebSocket message with type and data
   */
  handleMessage(message) {
    const { type, data } = message;

    if (!type) {
      console.warn("Received message without type:", message);
      return;
    }

    // Call all registered handlers for this message type
    const handlers = this.messageHandlers.get(type) || [];
    handlers.forEach((handler) => {
      try {
        handler(data);
      } catch (error) {
        console.error(`Error in message handler for type ${type}:`, error);
      }
    });
  }

  /**
   * Register a handler for a specific message type
   * @param {string} type - Message type (e.g., 'execution_created')
   * @param {Function} handler - Handler function that receives message data
   * @returns {Function} Unsubscribe function
   */
  on(type, handler) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, []);
    }

    this.messageHandlers.get(type).push(handler);

    // Return unsubscribe function
    return () => {
      const handlers = this.messageHandlers.get(type);
      if (handlers) {
        const index = handlers.indexOf(handler);
        if (index > -1) {
          handlers.splice(index, 1);
        }
      }
    };
  }

  /**
   * Send a message to the WebSocket server
   * @param {Object} message - Message to send
   */
  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn("Cannot send message: WebSocket not connected");
    }
  }

  /**
   * Get current connection state
   * @returns {string} Current connection state
   */
  getConnectionState() {
    if (!this.ws) return ConnectionState.DISCONNECTED;

    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return ConnectionState.CONNECTING;
      case WebSocket.OPEN:
        return ConnectionState.CONNECTED;
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
        return ConnectionState.DISCONNECTED;
      default:
        return ConnectionState.DISCONNECTED;
    }
  }
}

// Create and export singleton instance
export const websocketService = new WebSocketService();

// Export the class for testing
export default WebSocketService;
