package models

import "github.com/google/uuid"

type User struct {
	Id     uuid.UUID
	Name   string
	Age    int
	Rating int
}

type PlayerScore struct {
	Id           uuid.UUID
	Place        int
	NewRating    int
	RatingChange int
}

type PlayerStats struct {
	Id     uuid.UUID
	Rating int
	Place  int
}
