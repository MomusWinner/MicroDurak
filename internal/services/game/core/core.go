package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"slices"
)

var (
	DefaultGameSettings = GameSettings{
		TimeOver: 3000000000,
	}
)

type GameSettings struct {
	TimeOver float64
}

type Card struct {
	Suit int `json:"suit"`
	Rank int `json:"rank"`
}

type UserStatus int

const (
	Wait UserStatus = iota
	Attack
	Deffend
)

type User struct {
	Id         string     `json:"id"`
	Place      int        `json:"place"`
	Status     UserStatus `json:"status"`
	Name       string     `json:"name"`
	Cards      []Card     `json:"cards"`
	TakenCards []Card     `json:"taken_cards"`
}

type UserResponse struct {
	Id               string     `json:"id"`
	Status           UserStatus `json:"status"`
	Name             string     `json:"name"`
	CardLength       int        `json:"card_length"`
	TakenCardsLength int        `json:"taken_cards_length"`
}

type TableCard struct {
	Card
	BeatOff *Card `json:"beat_off"`
}

type Game struct {
	Id              string        `json:"id"`
	Settings        *GameSettings `json:"settings"`
	Users           []*User       `json:"users"`
	AttackingId     string        `json:"attacking_id"`
	DefendingId     string        `json:"defending_id"`
	Deck            []Card        `json:"deck"`
	TrumpSuit       int           `json:"trump_suit"` // TODO: remove
	TableCards      []TableCard   `json:"table_cards"`
	EndAttackUserId []string      `json:"end_attack_user_id"`
	ReadyUsers      []string      `json:"ready_users"`
	IsStarted       bool          `json:"is_Started"`

	AttackTimerIsRunning bool      `json:"attack_timer_is_running"`
	AttackTimerStartedAt time.Time `json:"attack_timer_started_at"`
	AttackTimerEndedAt   time.Time `json:"attack_timer_ended_at"`

	DefendTimerIsRunning bool      `json:"defend_timer_is_running"`
	DefendTimerStartedAt time.Time `json:"defend_timer_started_at"`
	DefendTimerEndedAt   time.Time `json:"defend_timer_ended_at"`

	GameEventBuffer []GameEventContainer `json:"game_event_buffer"`
}

func CreateNewGameAndSaveInRedis(redis *redis.Client, userIds []string) (*Game, error) { // TODO: move to handler layer
	game, err := CreateNewGame(userIds)

	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(game)
	log.Print(string(result))

	ctx := context.Background()
	// TODO: Change expiration time
	status := redis.Set(ctx, "game:"+game.Id, string(result), 0)
	if status.Err() != nil {
		log.Fatal("Couldn't create game room: " + game.Id)
	} else {
		log.Print("Create game room: " + game.Id)
	}

	return game, nil
}

func SaveGame(game *Game, redis *redis.Client) { // TODO: move to handler layer
	result, _ := json.Marshal(game)
	log.Print(string(result))

	ctx := context.Background()
	// TODO: Change expiration time
	status := redis.Set(ctx, "game:"+game.Id, string(result), 0)
	if status.Err() != nil {
		log.Fatal("Couldn't save game room: " + game.Id)
	} else {
		log.Print("Save game room: " + game.Id)
	}
}

func CreateNewGame(userIds []string) (*Game, error) {
	id := uuid.New()
	deck := generateDeck()
	trum_suit := deck[0].Suit

	shackeCards(deck)
	users := make([]*User, len(userIds))
	for i := range users {
		userCards := deck[len(deck)-6:]
		deck = deck[:len(deck)-6]

		users[i] = &User{
			Id:     userIds[i],
			Place:  i,
			Status: Wait,
			Name:   "",
			Cards:  userCards,
		}
	}

	game := Game{
		Id:         id.String(),
		Settings:   &DefaultGameSettings,
		Users:      users,
		Deck:       deck,
		TrumpSuit:  trum_suit,
		TableCards: []TableCard{},
	}

	// select first attacking and defending user
	attacking_i := rand.IntN(len(users))
	game.AttackingId = users[attacking_i].Id
	defending, err := game.nextUser(game.AttackingId)
	if err != nil {
		// TODO: log
	}
	game.DefendingId = defending.Id

	return &game, nil
}

func LoadGame(redis *redis.Client, gameId string) (*Game, error) {
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

func (g *Game) HandleMessage(msg []byte) (map[string][]byte, error) {
	var command Command
	err := json.Unmarshal(msg, &command)
	if err != nil {
		return nil, err
	}

	user, err := g.getUserById(command.UserId)
	if err != nil {
		return nil, err
	}

	var response CommandResponse

	switch command.Action {
	case ACTION_READY:
		response = g.ReadyHandler(command, user)
	case ACTION_ATTACK:
		var attackCommand AttackCommand
		json.Unmarshal(msg, &attackCommand)
		response = g.AttackHandler(attackCommand, user)
	case ACTION_DEFEND:
		var defendCommand DefendCommand
		json.Unmarshal(msg, &defendCommand)
		response = g.DefendHandler(defendCommand, user)
	case ACTION_END_ATTACK:
		response = g.EndAttackHandler(command, user)
	case ACTION_TAKE_ALL_CARDS:
		response = g.TakeAllCardHandler(command, user)
	case ACTION_CHECK_ATTACK_TIMER:
		response = g.CheckAttackTimerHandler(command, user)
	case ACTION_CHECK_DEFEND_TIMER:
		response = g.CheckDefendTimerHandler(command, user)
	default:
		response = CommandResponse{
			Error:   ERROR_UNREGISTERED_ACTION,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	return g.GeneratePack(response, user), nil
}

func (g *Game) StartAttackTimer() {
	g.AttackTimerIsRunning = true
	g.AttackTimerStartedAt = time.Now()
}

func (g *Game) StartDefendTimer() {
	g.DefendTimerIsRunning = true
	g.DefendTimerStartedAt = time.Now()
}

func (g *Game) StopAttackTimer() {
	g.AttackTimerIsRunning = false
}

func (g *Game) StopDefendTimer() {
	g.DefendTimerIsRunning = false
}

func (g *Game) removeUserCard(userId string, suit int, rank int) error {
	user, err := g.getUserById(userId)
	if err != nil {
		return err
	}

	for i := range user.Cards {
		if user.Cards[i].Suit == suit && user.Cards[i].Rank == rank {
			user.Cards = append(user.Cards[i+1:], user.Cards[:i]...)
			return nil
		}
	}

	return errors.New("User not found")
}

func (g *Game) beatOffCard(suit int, rank int, targetCard Card) bool {
	for i := range g.TableCards {
		if g.TableCards[i].Suit == targetCard.Suit && g.TableCards[i].Rank == targetCard.Rank {
			if CardGreater(suit, rank, targetCard.Suit, targetCard.Rank, g.TrumpSuit) {
				g.TableCards[i].BeatOff = &Card{Suit: suit, Rank: rank}
				return true
			} else {
				return false
			}
		}
	}

	return false
}

func (g *Game) GeneratePack(response CommandResponse, user *User) map[string][]byte {
	responseByUser := make(map[string][]byte, 0)

	messagePackByUser := g.CreateMessangePackByUserFromEventBuffer()

	for userId, messagePack := range messagePackByUser {
		if userId == user.Id {
			r := []any{response}
			messagePack.Messages = append(r, messagePack.Messages...)
		}
		messageString, err := json.Marshal(messagePack)
		if messageString == nil || err != nil {
			panic(err)
		}
		responseByUser[userId] = messageString
	}

	return responseByUser
}

func (g *Game) CreateMessangePackByUserFromEventBuffer() map[string]MessagePack {
	result := make(map[string]MessagePack)
	for _, user := range g.Users {
		events := g.GameEventBuffer

		eventPackEvents := []any{}
		for _, event := range events {

			estring, _ := json.Marshal(event)
			fmt.Println("---------------------" + user.Id)
			fmt.Println(string(estring))
			fmt.Println("---------------------")
			eventPackEvents = append(eventPackEvents, event)
		}
		result[user.Id] = MessagePack{
			Messages:  eventPackEvents,
			GameState: gameToGameStateResponse(g, user),
		}
	}

	g.GameEventBuffer = []GameEventContainer{}

	return result
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

func (g *Game) nextUser(userId string) (*User, error) {
	user, err := g.getUserById(userId)
	if err != nil {
		return nil, err
	}

	nextPlace := user.Place + 1
	if nextPlace >= len(g.Users) {
		nextPlace = 0
	}

	return g.Users[nextPlace], nil
}

func (g *Game) getUserById(userId string) (*User, error) {
	for _, u := range g.Users {
		if u.Id == userId {
			return u, nil
		}
	}

	return nil, errors.New("Not found")
}

func (g *Game) getUserByPlace(place int) (*User, error) {
	for _, u := range g.Users {
		if u.Place == place {
			return u, nil
		}
	}

	return nil, errors.New("Not found")
}

func (g *Game) getObservingUsers() []*User {
	users := make([]*User, 0, len(g.Users)-2)
	if len(g.Users) == 2 {
		return users
	}

	for _, u := range g.Users {
		if u.Id != g.AttackingId && u.Id != g.DefendingId {
			users = append(users, u)
		}
	}

	return users
}

func getCardBySuitAndRank(cards []Card, suit int, rank int) (Card, error) {
	for i := range cards {
		if cards[i].Suit == suit && cards[i].Rank == rank {
			return cards[i], nil
		}
	}

	return Card{}, errors.New("Not found")
}

func tableHasCardRank(tableCards []TableCard, rank int) bool {
	for i := range tableCards {
		if tableCards[i].Rank == rank {
			return true
		}

		beatOffCard := tableCards[i].BeatOff
		if beatOffCard != nil {
			if beatOffCard.Rank == rank {
				return true
			}
		}
	}

	return false
}

func tableHasCard(tableCards []TableCard, suit int, rank int) bool {
	for i := range tableCards {
		if tableCards[i].Suit == suit && tableCards[i].Rank == rank {
			return true
		}

		beatOffCard := tableCards[i].BeatOff
		if beatOffCard != nil {
			if beatOffCard.Suit == suit && beatOffCard.Rank == rank {
				return true
			}
		}
	}

	return false
}

func CardGreater(fsuit int, frank int, ssuit int, srank int, trump int) bool {
	if fsuit == trump && ssuit != trump {
		return true
	} else if fsuit != trump && ssuit == trump {
		return false
	} else {
		return frank > srank
	}
}

func allCardBeatOff(cards []TableCard) bool {
	for i := range cards {
		if cards[i].BeatOff == nil {
			return false
		}
	}
	return true
}

func tableCardsToCards(tableCards []TableCard) []Card {
	cards := []Card{}

	for i := range tableCards {
		cards = append(cards, Card{Suit: tableCards[i].Suit, Rank: tableCards[i].Rank})
		if tableCards[i].BeatOff != nil {
			cards = append(cards, *tableCards[i].BeatOff)
		}
	}

	return cards
}

func contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func cardInfo(cards []Card) string {
	result := ""
	for _, card := range cards {
		result += fmt.Sprintf("Card |suit:%d |rank:%d\n", card.Suit, card.Rank)
	}
	return result
}

func (g *Game) AddEventToBuffer(event GameEventContainer) {
	g.GameEventBuffer = append(g.GameEventBuffer, event)
}

func (g *Game) EndAttack(switchUsers bool) {
	attacker, _ := g.getUserById(g.AttackingId)
	defender, _ := g.getUserById(g.DefendingId)

	otherUsers := []*User{}
	copy(otherUsers, g.Users)
	removeUser(otherUsers, attacker.Id, defender.Id)

	g.AddCardsToUser(attacker)
	for _, user := range otherUsers {
		g.AddCardsToUser(user)
	}
	g.AddCardsToUser(defender)

	if switchUsers {
		newDefending, _ := g.nextUser(g.DefendingId)

		g.AttackingId = g.DefendingId
		g.DefendingId = newDefending.Id
	}

	g.TableCards = []TableCard{}
	g.StopAttackTimer()
	g.StopDefendTimer()
	endEvent := NewEndAttackEvent()
	g.AddEventToBuffer(endEvent)
}

func (g *Game) AddCardsToUser(user *User) {
	if len(user.Cards) >= 6 {
		return
	}

	user.TakenCards = []Card{}

	addCardAmount := 6 - len(user.Cards)
	for range addCardAmount {
		card, err := g.TakeCardFromDeck()
		if err != nil {
			break
		}

		user.Cards = append(user.Cards, card)
		user.TakenCards = append(user.TakenCards, card)
	}
}

func (g *Game) TakeCardFromDeck() (Card, error) {
	if len(g.Deck) <= 0 {
		return Card{}, errors.New("Deck is empty")
	}

	card := g.Deck[len(g.Deck)-1]
	g.Deck = g.Deck[:len(g.Deck)-1]

	return card, nil
}

func removeUser(users []*User, userIds ...string) {
	for i, user := range users {
		for _, removeUser := range userIds {
			if user.Id == removeUser {
				users = slices.Delete(users, i, i+1)
			}
		}
	}
}
