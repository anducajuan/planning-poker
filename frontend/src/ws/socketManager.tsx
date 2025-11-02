// ws/socketManager.ts
type EventCallback = (event: string, data: unknown) => void;

class SocketManager {
  private ws: WebSocket | null = null;
  private sessionId: string | null = null;
  private eventQueue: Array<{ event: string; data: unknown }> = [];
  private onEventCallback: EventCallback | null = null;
  private reconnecting = false;

  connect(sessionId: string, onEvent: EventCallback) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    this.sessionId = sessionId;
    this.onEventCallback = onEvent;

    this.ws = new WebSocket(`ws://localhost:8080/ws?session_id=${sessionId}`);

    this.ws.addEventListener("open", () => {
      this.flushQueue();
    });

    this.ws.addEventListener("message", (msg) => {
      try {
        const { event, data } = JSON.parse(msg.data as string);
        if (this.onEventCallback) this.onEventCallback(event, data);
      } catch (e) {
        console.error("[WS] Parse error:", e);
      }
    });

    this.ws.addEventListener("close", () => {
      console.warn("[WS] Disconnected");
      this.scheduleReconnect();
    });

    this.ws.addEventListener("error", (err) => {
      console.error("[WS] Error:", err);
      this.scheduleReconnect();
    });
  }

  private scheduleReconnect() {
    if (this.reconnecting || !this.sessionId) return;
    this.reconnecting = true;
    console.log("[WS] Attempting reconnect in 2s...");
    setTimeout(() => {
      this.reconnecting = false;
      this.connect(this.sessionId!, this.onEventCallback!);
    }, 2000);
  }

  private flushQueue() {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    this.eventQueue.forEach(({ event, data }) => this.send(event, data));
    this.eventQueue = [];
  }

  send(event: string, data: unknown) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.eventQueue.push({ event, data });
      return;
    }
    try {
      this.ws.send(JSON.stringify({ event, data }));
    } catch {
      this.eventQueue.push({ event, data });
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.eventQueue = [];
  }

  waitForOpen(sessionId: string, timeout = 5000): Promise<void> {
    return new Promise((resolve, reject) => {
      const start = Date.now();

      const check = () => {
        if (
          this.sessionId === sessionId &&
          this.ws &&
          this.ws.readyState === WebSocket.OPEN
        ) {
          cleanup();
          return resolve();
        }
        if (Date.now() - start >= timeout) {
          cleanup();
          return reject(new Error("WS open timeout"));
        }
      };

      const interval = setInterval(check, 50);

      let onOpen: (() => void) | null = null;
      const wsRef: WebSocket | null = this.ws;
      if (this.ws) {
        onOpen = () => {
          if (this.sessionId === sessionId) {
            cleanup();
            resolve();
          }
        };
        this.ws.addEventListener("open", onOpen);
      }

      function cleanup() {
        clearInterval(interval);
        if (onOpen && wsRef) {
          try {
            wsRef.removeEventListener("open", onOpen);
          } catch {
            // ignore
          }
        }
      }
    });
  }
}

export const socketManager = new SocketManager();
