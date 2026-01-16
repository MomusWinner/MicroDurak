package props

type LoginReq struct {
	Email    string
	Password string
}

type LoginResp struct {
	PlayerId string `json:"player_id"`
	Token    string `json:"token"`
}

// func DtoToProps(dto dto.CreateUserDto) CreateUserReq { // TODO:
// 	return CreateUserReq{
// 		Email:      dto.Email,
// 		FirstName:  dto.FirstName,
// 		SecondName: dto.SecondName,
// 	}
// }
