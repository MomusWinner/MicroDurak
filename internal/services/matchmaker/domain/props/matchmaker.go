package props

import (
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/types"
)

type FindMatchReq struct {
	PlayerId string
	Rating   int
	SentTime time.Time
}

type FindMatchResp struct {
	Status    types.ItemStatus
	RoomId    string
	GroupSize int
	Error     error
}

type CancelMatchReq struct {
	PlayerId string
}
