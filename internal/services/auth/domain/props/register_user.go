package props

import "github.com/google/uuid"

type RegisterReq struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResp struct {
	PlayerId uuid.UUID `json:"player_id"`
	Token    string    `json:"token"`
}
