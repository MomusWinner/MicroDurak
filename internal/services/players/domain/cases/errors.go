package cases

import "errors"

var (
	ErrInternal  = errors.New("Server internal error")
	ErrNoPlayers = errors.New("No players error")
)
