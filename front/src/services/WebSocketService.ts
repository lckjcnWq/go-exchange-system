import { message } from 'antd';
import { CONFIG } from '../config';

export class WebSocketService{
    private static instance: WebSocketService;
    private ws: WebSocket | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout = 3000;
    private messageHandlers: Map<string, Function[]> = new Map();

    private constructor() {
        this.connect();
    }

    public static getInstance(): WebSocketService {
        if (!WebSocketService.instance) {
            WebSocketService.instance = new WebSocketService();
        }
        return WebSocketService.instance;
    }

    public connect() {
        try {
            this.ws = new WebSocket(CONFIG.WS_URL);
            this.ws.onopen = this.onOpen.bind(this);
            this.ws.onmessage = this.onMessage.bind(this);
            this.ws.onclose = this.onClose.bind(this);
            this.ws.onerror = this.onError.bind(this);
        } catch (error) {
            console.error('WebSocket connection error:', error);
        }
    }

    private onOpen() {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
    }

    private onClose() {
        console.log('WebSocket closed');
        this.reconnect();
    }

    private onMessage(event: MessageEvent) {
        try {
            const data = JSON.parse(event.data);
            const handlers = this.messageHandlers.get(data.type);
            if (handlers) {
                handlers.forEach(handler => handler(data));
            }
        } catch (e) {
            console.error('Error parsing WebSocket message:', e);
        }
    }

    private onError(error: Event) {
        console.error('WebSocket error:', error);
    }

    private reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            setTimeout(() => {
                console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);
                this.connect();
            }, this.reconnectTimeout);
        } else {
            message.error('WebSocket connection failed');
        }
    }

    public subscribe(type: string, handler: Function) {
        if (!this.messageHandlers.has(type)) {
            this.messageHandlers.set(type, []);
        }
        this.messageHandlers.get(type)?.push(handler);
    }

    public unsubscribe(type: string, handler: Function) {
        const handlers = this.messageHandlers.get(type);
        if (handlers) {
            this.messageHandlers.set(
                type,
                handlers.filter(h => h !== handler)
            );
        }
    }

    public send(message: any) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not connected');
        }
    }
}