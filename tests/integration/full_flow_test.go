package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/game/core"
)

const (
	HTTP_URL_AUTH     = "http://localhost:8080/api/v1/auth/"
	WS_URL_MATCHMAKER = "ws://localhost:3000/api/v1/matchmaker/find-match"
	WS_URL_GAME       = "ws://localhost:7070/api/v1/game-manager/"
)

type User struct {
	Id       string
	Token    string
	GameId   string
	Name     string
	Age      int
	Email    string
	Password string
}

func TestFullUserFlow(t *testing.T) {
	user1 := User{
		Name:     "Oleg",
		Age:      21,
		Email:    "oleg@gmail.com",
		Password: "password",
	}

	user2 := User{
		Name:     "NeOleg",
		Age:      22,
		Email:    "neoleg@gmail.com",
		Password: "password",
	}

	{ // Authorize
		httpClient := NewHTTPClient()
		httpClient.RegisterAndLogin(t, &user1)
		httpClient.RegisterAndLogin(t, &user2)
		fmt.Println(user1)
		fmt.Println(user2)
	}

	{ // Find Match
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			wsClient := NewWSClient(user1.Token)
			wsClient.Connect(t, WS_URL_MATCHMAKER)
			wsClient.FindMatch(t, &user1)
			wsClient.Close(t)
		}()
		go func() {
			defer wg.Done()
			wsClient := NewWSClient(user2.Token)
			wsClient.Connect(t, WS_URL_MATCHMAKER)
			wsClient.FindMatch(t, &user2)
			wsClient.Close(t)
		}()
		wg.Wait()

		fmt.Println(user1.GameId)
	}

	{ // Play
		var wg sync.WaitGroup

		wg.Add(2)
		ch := make(chan string, 1)
		go func() {
			defer wg.Done()
			wsClient := NewWSClient(user1.Token)
			wsClient.Connect(t, WS_URL_GAME+user1.GameId)
			wsClient.GamePlay(t, user1, ch)
		}()
		go func() {
			defer wg.Done()
			wsClient := NewWSClient(user2.Token)
			wsClient.Connect(t, WS_URL_GAME+user2.GameId)
			wsClient.GamePlay(t, user2, ch)
		}()
		wg.Wait()
	}
}

func (c *HTTPClient) RegisterAndLogin(t *testing.T, user *User) {
	type RegistrationRequest struct {
		Age      int    `json:"age"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type AuthResp struct {
		PlayerId string `json:"player_id"`
		Token    string `json:"token"`
	}

	regData := RegistrationRequest{
		Age:      user.Age,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	var token AuthResp

	regResp, regBoody := c.POST(t, HTTP_URL_AUTH+"register", regData)

	if regResp.StatusCode == http.StatusCreated {
		json.Unmarshal(regBoody, &token)
		user.Id = token.PlayerId
		user.Token = token.Token
		return
	}

	loginData := LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	loginResp, loginBoody := c.POST(t, HTTP_URL_AUTH+"login", loginData)

	if loginResp.StatusCode == http.StatusOK {
		json.Unmarshal(loginBoody, &token)
		user.Id = token.PlayerId
		user.Token = token.Token
		return
	}

	t.Error("Authorization failed")
}

func (c *WSClient) FindMatch(t *testing.T, user *User) {
	type FindMatchStatus struct {
		Status string `json:"status"`
		GameId string `json:"game_id"`
	}

	for {
		message, err := c.Receive(t, 20*time.Second)
		if err != nil {
			t.Error(err)
		}

		var matchStatus FindMatchStatus
		json.Unmarshal(message, &matchStatus)
		fmt.Println(matchStatus)
		if matchStatus.Status == "created" {
			user.GameId = matchStatus.GameId
			break
		}
	}
}

func (c *WSClient) GamePlay(t *testing.T, user User, ch chan string) {
	c.GameSendReady(t, user)
	attacker := false

	time.Sleep(500 * time.Millisecond) // TODO:

	pack := c.GameWaitStart(t, user)

	if pack.GameState.Me.Id == pack.GameState.AttackingId {
		attacker = true
	}

	if attacker {
		fmt.Printf("%s is Attacker", user.Name)

		c.GameSendAttack(t, user, pack.GameState.Me.Cards[0])
		ch <- "attack"

		for {
			msg := <-ch
			if msg == "take_all" {
				time.Sleep(10 * time.Millisecond)
				pack = c.GameReceive(t, user)
				// c.GameEndAttack(t, user)
				// time.Sleep(10 * time.Millisecond)
				// pack = c.GameReceive(t, user)
				time.Sleep(10 * time.Millisecond)
				c.GameSendAttack(t, user, pack.GameState.Me.Cards[0])
				ch <- "attack"
			}
		}
	} else {
		fmt.Printf("%s is not Attacker", user.Name)

		for {
			msg := <-ch
			if msg == "attack" {
				time.Sleep(10 * time.Millisecond)
				pack = c.GameReceive(t, user)
				c.GameSendTakeAllCards(t, user)
				time.Sleep(10 * time.Millisecond)
				pack = c.GameReceive(t, user)
				ch <- "take_all"
			}
		}
	}
}

func (c *WSClient) GameSendCommand(t *testing.T, user User, action string) {
	readyCommand := core.Command{
		Action: action,
		UserId: user.Id,
		GameId: user.GameId,
	}

	command, err := json.Marshal(readyCommand)
	if err != nil {
		t.Error(err)
	}

	c.Send(t, command)
}

func (c *WSClient) GameSendReady(t *testing.T, user User) {
	c.GameSendCommand(t, user, core.ACTION_READY)
	fmt.Printf("%s: Ready\n", user.Name)
}

func (c *WSClient) GameEndAttack(t *testing.T, user User) {
	c.GameSendCommand(t, user, core.ACTION_END_ATTACK)
	fmt.Printf("%s: End attack\n", user.Name)
}

func (c *WSClient) GameSendAttack(t *testing.T, user User, card core.Card) {
	readyCommand := core.AttackCommand{
		Command: core.Command{
			Action: core.ACTION_ATTACK,
			UserId: user.Id,
			GameId: user.GameId,
		},
		Card: card,
	}

	command, err := json.Marshal(readyCommand)
	if err != nil {
		t.Error(err)
	}

	c.Send(t, command)
	fmt.Printf("%s: Attack\n", user.Name)
}

func (c *WSClient) GameSendTakeAllCards(t *testing.T, user User) {
	c.GameSendCommand(t, user, core.ACTION_TAKE_ALL_CARDS)
}

type MessagePack struct {
	Messages  []core.GameEvent       `json:"messages"`
	GameState core.GameStateResponse `json:"game_state"`
}

func (c *WSClient) GameReceive(t *testing.T, user User) MessagePack {
	fmt.Printf("%s: Recive message pack\n", user.Name)
	var pack MessagePack

	data, err := c.Receive(t, 100*time.Second)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(data, &pack)

	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s: Deck size %d\n", user.Name, pack.GameState.DeckLength)
	fmt.Println("--------------------")
	fmt.Println(pack.Messages)

	return pack
}

func (c *WSClient) GameWaitStart(t *testing.T, user User) MessagePack {
	var pack MessagePack
	for {
		pack = c.GameReceive(t, user)
		fmt.Printf("%s: Wait", user.Name)
		for _, event := range pack.Messages {
			if event.Event == core.EVENT_START {
				return pack
			}
		}
	}
}
