package types

import (
	"errors"
	"fmt"
)

var ErrGroupNotFound = errors.New("matchmaker: group not found")

type ErrGroupTooSmall struct {
	Gid int
}

func (e ErrGroupTooSmall) Error() string {
	return fmt.Sprintf("matchmaker: group %d too small found", e.Gid)
}

func NewGroupTooSmall(gid int) error {
	return &ErrGroupTooSmall{gid}
}
