package types

import "time"

type ItemStatus = int

const (
	MatchPending ItemStatus = iota
	MatchFound
	MatchError
)

type MatchResponse struct {
	Status ItemStatus
	RoomId string
}

type MatchChan struct {
	PlayerId   string
	Rating     int
	SentTime   time.Time
	ReturnChan chan<- MatchResponse
}
