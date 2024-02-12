let ws

export const getWebSocket = () => ws

export const initializeWebSocket = () => {
	ws = new WebSocket(`wss://${location.hostname}:8080/ws`)

	ws.onclose = () => {}

	ws.onerror = (error) => {}
}
