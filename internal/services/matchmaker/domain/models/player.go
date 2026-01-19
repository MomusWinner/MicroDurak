package models

type PlayerStatus = int

const (
	StatusEmpty PlayerStatus = iota
	StatusSearch
	StatusMoved
)

type RedisPlayer struct {
	Status PlayerStatus
	Id     string
	Gid    int
}
