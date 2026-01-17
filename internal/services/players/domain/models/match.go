package models

import "github.com/google/uuid"

type GameResult int32

const (
	GameResult_WIN         GameResult = 0
	GameResult_DRAW        GameResult = 1
	GameResult_INTERRUPTED GameResult = 2
)

var (
	GameResult_name = map[int32]string{
		0: "win",
		1: "draw",
		2: "interrupted",
	}
	GameResult_value = map[string]int32{
		"win":         0,
		"draw":        1,
		"interrupted": 2,
	}
)

type Match struct {
	Id          uuid.UUID
	PlayerCount int
	GameResult  GameResult
}

type PlayerPlacement struct {
	PlayerId    string // TODO: rename
	PlayerPlace int
}

type PlayerMatchResult struct {
	Id           uuid.UUID
	Rating       int32
	RatingChange int32
}
