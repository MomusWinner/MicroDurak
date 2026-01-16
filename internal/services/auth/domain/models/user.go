package models

import "github.com/google/uuid"

type User struct {
	Name string
	Age  int
}

type AuthUser struct {
	Id       uuid.UUID
	PlayerId uuid.UUID
	Email    string
	Password string
}
