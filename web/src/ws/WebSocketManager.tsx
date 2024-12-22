type WebSocketManagerOptions = {
  heartbeatInterval?: number; 
  reconnectInterval?: number; 
  maxReconnectAttempts?: number; 
  onMessage?: (message: string) => void; 
  onOpen?: () => void; 
  onClose?: (event: CloseEvent) => void; 
  onError?: (error: ErrorEvent) => void;
};

export class WebSocketManager {
  private url: string;
  private heartbeatInterval: number;
  private reconnectInterval: number;
  private maxReconnectAttempts: number;
  private currentReconnectAttempts: number;
  private ws: WebSocket | null;
  private heartbeatTimer: number | null;
  private reconnectTimer: number | null;
  private onMessage: (message: string) => void;
  private onOpen: () => void;
  private onClose: (event: CloseEvent) => void;
  private onError: (error: ErrorEvent) => void;

  constructor(url: string, options: WebSocketManagerOptions = {}) {
    this.url = url;
    this.heartbeatInterval = options.heartbeatInterval || 5000;
    this.reconnectInterval = options.reconnectInterval || 3000;
    this.maxReconnectAttempts = options.maxReconnectAttempts || 10;
    this.currentReconnectAttempts = 0;
    this.ws = null;
    this.heartbeatTimer = null;
    this.reconnectTimer = null;
    this.onMessage = options.onMessage || (() => {});
    this.onOpen = options.onOpen || (() => {});
    this.onClose = options.onClose || (() => {});
    this.onError = options.onError || (() => {});

    this.connect();
  }

  private connect(): void {
    if (this.ws) {
      this.ws.close();
    }

    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      this.currentReconnectAttempts = 0;
      this.startHeartbeat();
      this.onOpen();
    };

    this.ws.onmessage = (event: MessageEvent) => {
      this.onMessage(event.data);
    };

    this.ws.onclose = (event: CloseEvent) => {
      this.stopHeartbeat();
      this.onClose(event);
      this.reconnect();
    };

    this.ws.onerror = (error: Event) => {
      if (this.ws) {
        this.ws.close();  // Safely close WebSocket if error occurs
      }
      this.onError(error as ErrorEvent);
    };
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();
    this.heartbeatTimer = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: "heartbeat", timestamp: Date.now() }));
      }
    }, this.heartbeatInterval);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private reconnect(): void {
    if (this.currentReconnectAttempts >= this.maxReconnectAttempts) {
      return;
    }

    this.currentReconnectAttempts++;
    this.reconnectTimer = window.setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  public send(data: any): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      
    }
  }

  public close(): void {
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close();
    }
    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer);
    }
  }
}
