package core

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand/v2"

	"github.com/MommusWinner/MicroDurak/services/game/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Command struct {
	Action string `json:"action"`
}

type AttackCommand struct {
	Command
}

const (
	ATTACKING_STATE = "attacking"
	DEFENDING_STATE = "defending"
)

type Card struct {
	Suit int `json:"suit"`
	Rank int `json:"rank"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type UserResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CardLength int    `json:"card_length"`
}

type TableCard struct {
	BeatOff Card `json:"rank"`
	Card
}

type Game struct {
	Redis       *redis.Client
	Conf        *config.Config
	Id          string      `json:"id"`
	State       string      `json:"state"` // attacking defending
	Users       []User      `json:"users"`
	AttackingId string      `json:"attacking_id"`
	DefendingId string      `json:"defending_id"`
	Deck        []Card      `json:"deck"`
	TableCards  []TableCard `json:"table_cards"`
}

type GameStateResponse struct {
	State       string         `json:"state"` // attacking defending
	Me          User           `json:"me"`
	Users       []UserResponse `json:"users"`
	AttackingId string         `json:"attacking_id"`
	DefendingId string         `json:"defending_id"`
	DeckLength  int            `json:"deck_length"`
	TableCards  []TableCard    `json:"table_cards"`
}

func CreateNewGame(redis *redis.Client, userIds []string) (*Game, error) {
	id := uuid.New()
	deck := generateDeck()
	// TODO: The last cards will not be an ace.
	shackeCards(deck)
	users := make([]User, len(userIds))
	for i := 0; i < len(users); i++ {
		userCards := deck[:6]
		deck = deck[6:]

		users[i] = User{
			Id:    userIds[i],
			Name:  "",
			Cards: userCards,
		}
	}

	// select first attacking and defending user
	attacking_i := rand.IntN(len(users))
	attacking_id := users[attacking_i].Id
	defending, _ := nextUser(users, attacking_id)

	game := Game{
		Redis:       redis,
		Id:          id.String(),
		State:       "attacking",
		Users:       users,
		AttackingId: attacking_id,
		DefendingId: defending.Id,
		Deck:        deck,
		TableCards:  make([]TableCard, 0),
	}

	result, _ := json.Marshal(game)
	log.Print(string(result))

	ctx := context.Background()
	// TODO: Change expiration time
	status := redis.Set(ctx, "game:"+id.String(), string(result), 0)
	if status.Err() != nil {
		log.Fatal("Couldn't create game room: " + id.String())
	} else {
		log.Print("Create game room: " + id.String())
	}

	return &game, nil
}

func LoadGame(redis *redis.Client, gameId string) (*Game, error) {
	// TODO: implement
	ctx := context.Background()
	value, err := redis.Get(ctx, "game:"+gameId).Result()
	if err != nil {
		return nil, err
	}

	var game Game

	log.Print(value)
	err = json.Unmarshal([]byte(value), &game)

	if err != nil {
		return nil, err
	}

	log.Printf("Success load game (%s)\n%s", gameId, value)

	return &game, err
}

func (g *Game) HandleMessage(msg []byte) error {
	var c Command
	err := json.Unmarshal(msg, &c)
	if err != nil {
		return err
	}

	return err
}

func generateDeck() []Card {
	deck := make([]Card, 36)
	i := 0

	for suit := 1; suit <= 4; suit++ {
		for rank := 6; rank <= 14; rank++ {
			deck[i] = Card{
				Suit: suit,
				Rank: rank,
			}
			i++
		}
	}

	return deck
}

func shackeCards(cards []Card) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func nextUser(users []User, userId string) (User, error) {
	length := len(users)
	for i := 0; i < length; i++ {
		if users[i].Id == userId {
			if i == length-1 {
				return users[0], nil
			} else {
				return users[i+1], nil
			}
		}
	}

	return User{}, errors.New("User not found")
}
