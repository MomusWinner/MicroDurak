package domain

type WebSocket interface {
	ReadMessage() (message []byte, err error)
	WriteMessage(message []byte) error
	Close() error
}
