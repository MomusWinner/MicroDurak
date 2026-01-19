package http

import (
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"
	"github.com/gorilla/websocket"
)

type webSocketAdapter struct {
	conn *websocket.Conn
}

func NewWebSocketAdapter(conn *websocket.Conn) domain.WebSocket {
	return &webSocketAdapter{conn: conn}
}

func (w *webSocketAdapter) ReadMessage() (message []byte, err error) {
	_, msg, err := w.conn.ReadMessage()
	return msg, err
}

func (w *webSocketAdapter) WriteMessage(message []byte) error {
	return w.conn.WriteMessage(1, message)
}

func (w *webSocketAdapter) Close() error {
	return w.conn.Close()
}
