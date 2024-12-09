const WebSocket = await import('websocket');

const ws = new WebSocket('ws://localhost:5117/ws');

export default ws;