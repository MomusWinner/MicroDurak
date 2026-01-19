package props

import "github.com/MommusWinner/MicroDurak/internal/services/game-manager/domain"

type HandleMessageReq struct {
	GameId  string
	UserId  string
	Message []byte
}

type HandleMessageResp struct {
	Success bool
}

type ConnectWebSocketReq struct {
	GameId    string
	UserId    string
	WebSocket domain.WebSocket
}
