package domain

type Messaging interface {
	ProcessQueue(gameId string, userId string, processMessage func([]byte)) error
	SendMessageToGame(message []byte) error
}
