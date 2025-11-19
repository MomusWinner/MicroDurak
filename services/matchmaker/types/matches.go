package types

import "time"

type ItemStatus int

const (
	MatchPending ItemStatus = iota
	MatchFoundGroup
	MatchCreated
	MatchError
)

func (s ItemStatus) String() string {
	switch s {
	case MatchPending:
		return "pending"
	case MatchCreated:
		return "created"
	case MatchFoundGroup:
		return "found_group"
	case MatchError:
		return "error"
	}
	return "unknown"
}

type MatchResponse struct {
	Status    ItemStatus
	RoomId    string
	GroupSize int
	Error     error
}

type MatchCancel struct {
	PlayerId string
}

type MatchChan struct {
	PlayerId   string
	Rating     int
	SentTime   time.Time
	ReturnChan chan<- MatchResponse
}
