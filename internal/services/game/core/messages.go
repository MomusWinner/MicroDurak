package core

const (
	ACTION_READY              = "ACTION_READY"
	ACTION_ATTACK             = "ACTION_ATTACK"
	ACTION_DEFEND             = "ACTION_DEFEND"
	ACTION_END_ATTACK         = "ACTION_END_ATTACK"
	ACTION_TAKE_ALL_CARDS     = "ACTION_TAKE_ALL_CARDS"
	ACTION_CHECK_ATTACK_TIMER = "ACTION_CHECK_ATTACK_TIMER" // TODO:
	ACTION_CHECK_DEFEND_TIMER = "ACTION_CHECK_DEFEND_TIMER" // TODO:
)

const (
	ERROR_EMPTY                                         = ""
	ERROR_SERVER                                        = "SERVER_ERROR"
	ERROR_BAD_REQUEST                                   = "BAD_REQUEST"
	ERROR_USER_ALREADY_READY                            = "USER_ALREADY_READY"
	ERROR_NOT_YOUR_TURN                                 = "NOT_YOUR_TURN"
	ERROR_INCORRECT_CARD                                = "INCORRECT_CARD"
	ERROR_USER_NO_HAS_CARD                              = "USER_NO_HAS_CARD"
	ERROR_ATTACK_TIME_OVER                              = "ATTACK_TIME_OVER"
	ERROR_DEFEND_TIME_OVER                              = "DEFEND_TIME_OVER"
	ERROR_NO_SAME_RANK_CARD_IN_TABLE                    = "NO_SAME_RANK_CARD_IN_TABLE"
	ERROR_NOT_FOUND_CART_ON_TABLE                       = "NOT_FOUND_CART_ON_TABLE"
	ERROR_TARGET_CARD_GREATER_THEN_YOUR                 = "TARGET_CARD_GREATER_THEN_YOUR"
	ERROR_GAME_SHOULD_BE_STARTED                        = "GAME_SHOULD_BE_STARTED"
	ERROR_CANNOT_END_ATTACK_IN_FIRST_TURN               = "CANNOT_END_ATTACK_IN_FIRST_TURN"
	ERROR_ALL_CARD_SHOULD_BE_BEAT_OFF_BEFORE_END_ATTACK = "ALL_CARD_SHOULD_BE_BEAT_OFF_BEFORE_END_ATTACK"
	ERROR_TABLE_HOLDS_ONLY_SIX_CARDS                    = "TABLE_HOLDS_ONLY_SIX_CARDS"
	ERROR_DEFENDER_NO_CARDS                             = "DEFENDER_NO_CARDS"
	ERROR_ALREADY_END_ATTACK                            = "ALREADY_END_ATTACK" // TODO: implement
	ERROR_UNREGISTERED_ACTION                           = "UNREGISTERED_ACTION"
)

type MessagePack struct {
	Messages  []interface{}     `json:"messages"`
	GameState GameStateResponse `json:"game_state"`
}

func gameToGameStateResponse(game *Game, targetUser *User) GameStateResponse {
	return GameStateResponse{
		Me:          *targetUser,
		Users:       usersToUserResponses(game.Users),
		AttackingId: game.AttackingId,
		DefendingId: game.DefendingId,
		DeckLength:  len(game.Deck),
		TrumpSuit:   game.TrumpSuit,
		TableCards:  game.TableCards,
	}
}

func usersToUserResponses(users []*User) []UserResponse {
	userResponses := make([]UserResponse, len(users))
	for i := 0; i < len(users); i++ {
		userResponses[i] = userToUserResponse(*users[i])
	}

	return userResponses
}

func userToUserResponse(user User) UserResponse {
	return UserResponse{
		Id:               user.Id,
		Name:             user.Name,
		CardLength:       len(user.Cards),
		TakenCardsLength: len(user.TakenCards),
	}
}

type GameStateResponse struct {
	Me          User           `json:"me"`
	Users       []UserResponse `json:"users"`
	AttackingId string         `json:"attacking_id"`
	DefendingId string         `json:"defending_id"`
	DeckLength  int            `json:"deck_length"`
	TrumpSuit   int            `json:"trump_suit"`
	TableCards  []TableCard    `json:"table_cards"`
}

// Requeste messages
type Command struct {
	GameId string `json:"game_id"`
	Action string `json:"action"`
	UserId string `json:"user_id"`
}

type AttackCommand struct {
	Card Card `json:"card"`
	Command
}

type DefendCommand struct {
	TargetCard Card `json:"target_card"`
	UserCard   Card `json:"user_card"`
	Command
}

// Response messages
type CommandResponse struct {
	Error   string            `json:"error"`
	Command any               `json:"command"`
	State   GameStateResponse `json:"state"`
}
