package cases

import "errors"

var (
	ErrInternal                  = errors.New("Server internal error")
	ErrNoPlayers                 = errors.New("No players error")
	ErrUnprocessableId           = errors.New("Unprocessable id")
	ErrPlayerNotFound            = errors.New("Player not found")
	ErrIncorrectPlayersPlacement = errors.New("Incorrect players placement")
)
