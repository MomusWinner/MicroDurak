package core

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func checkGameEvents(game *Game, events ...string) error { // TODO
	for _, v := range game.GameEventBuffer {
		eventType := GameEventToType(v)

		if eventType == EVENT_NONE {
			return errors.New("Event shouldn't be NONE")
		}

		if events[0] != eventType {
			return errors.New(fmt.Sprintf("Bad event (%s) target (%s) \n", eventType, events[0]))
		}
		events = events[1:]
	}
	game.GameEventBuffer = []GameEventContainer{}

	return nil
}

func checkPack(pack map[string][]byte) error {
	all := ""

	for userId, msg := range pack {
		all += string(msg) + "\n"
		if msg == nil {
			return errors.New("Message pack has nil message. UserId: " + userId)
		}
	}

	return errors.New(all)
}

func createGameWithReadyUsers() (*Game, *User, *User, error) {
	game, _ := CreateNewGame([]string{"user1", "user2"})

	attackUser, _ := game.getUserById(game.AttackingId)
	defendUser, _ := game.getUserById(game.DefendingId)

	response := game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: attackUser.Id,
	}, attackUser)
	_ = response

	game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: defendUser.Id,
	}, defendUser)

	err := checkGameEvents(game, EVENT_READY, EVENT_START)

	return game, attackUser, defendUser, err
}

func TestGanaratePack(t *testing.T) {
	game, _ := CreateNewGame([]string{"user1", "user2"})

	attackUser, _ := game.getUserById(game.AttackingId)
	defendUser, _ := game.getUserById(game.DefendingId)

	game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: attackUser.Id,
	}, attackUser)

	response := game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: defendUser.Id,
	}, defendUser)

	// ATTACK
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackUser.Cards[0],
	}

	game.AttackHandler(attackC, attackUser)

	packByUser := game.GeneratePack(response, defendUser)
	for user, pack := range packByUser {
		fmt.Println(user)
		fmt.Println("----------------------------------------")
		fmt.Println(string(pack))
	}
}

// func TestGameToGameStateResponse(t *testing.T) {
// 	game, _ := CreateNewGame([]string{"user1", "user2"})
// 	attackUser, _ := game.getUserById(game.AttackingId)
// 	gameState := gameToGameStateResponse(game, attackUser)
//
// 	if gameState.Me.Id != attackUser.Id {
// 		t.Error("Encorrect GameStateResponse")
// 	}
// }

//	func TestCreateNewGame(t *testing.T) {
//		game, err := CreateNewGame([]string{"test1", "test2"})
//		if err != nil {
//			t.Error(err)
//		}
//		for _, user := range game.Users {
//			if len(user.Cards) != 6 {
//				t.Error("User must have 6 cards.")
//			}
//		}
//
//		_ = game
//	}
func TestAttackCycle(t *testing.T) {
	game, attackUser, defendUser, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	// ATTACK
	attackCard := attackUser.Cards[0]
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackUser.Cards[0],
	}

	r := game.AttackHandler(attackC, attackUser)
	if r.Error != ERROR_EMPTY {
		t.Error("Attack error CODE: " + r.Error)
	}

	if len(attackUser.Cards) != 5 {
		t.Error("The card was not removed after the attack.")
	}

	if game.TableCards[0].Suit != attackCard.Suit || game.TableCards[0].Rank != attackCard.Rank {
		t.Error("Table shouldn't be empty after attack")
	}

	// DEFEND
	defendUser.Cards[0] = Card{Suit: game.TrumpSuit, Rank: 15}
	defendC := DefendCommand{
		Command: Command{
			Action: ACTION_DEFEND,
			UserId: defendUser.Id,
		},
		UserCard:   defendUser.Cards[0],
		TargetCard: attackCard,
	}

	r = game.DefendHandler(defendC, defendUser)

	if r.Error != ERROR_EMPTY {
		t.Error("Defend error CODE: " + r.Error)
	}

	// END ATTACK
	game.EndAttackHandler(Command{
		Action: ACTION_END_ATTACK,
		UserId: attackUser.Id,
	}, attackUser)

	if len(attackUser.Cards) != 6 || len(defendUser.Cards) != 6 {
		t.Error("Atfter attack end game server should hand out the cards.")
	}

	if game.DefendingId != attackUser.Id || game.AttackingId != defendUser.Id {
		t.Error("After end attack game should change user roles.")
	}

	if len(game.TableCards) != 0 {
		t.Error("After end attack table should be emtpy")
	}

	err = checkGameEvents(game, EVENT_ATTACK, EVENT_DEFEND, EVENT_END_ATTACK)
	if err != nil {
		t.Error(err)
	}
}

func TestTakeAllCards(t *testing.T) {
	game, attackUser, defendUser, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	for i := range attackUser.Cards {

		attackUser.Cards[i] = Card{
			Suit: 1,
			Rank: 6,
		}
	}

	for range 3 {
		attackC := AttackCommand{
			Command: Command{
				Action: ACTION_ATTACK,
				UserId: attackUser.Id,
			},
			Card: attackUser.Cards[0],
		}

		r := game.AttackHandler(attackC, attackUser)
		if r.Error != ERROR_EMPTY {
			t.Error("Attack error CODE: " + r.Error)
		}
	}

	// Beat off first card on table
	game.TableCards[0].BeatOff = &Card{
		Suit: 1,
		Rank: 14,
	}

	rsp := game.TakeAllCardHandler(Command{
		Action: ACTION_TAKE_ALL_CARDS,
		UserId: defendUser.Id,
	}, defendUser)

	if rsp.Error != ERROR_EMPTY {
		t.Error("TakeAllCard Error: " + rsp.Error)
	}

	if len(defendUser.Cards) != 10 {
		t.Error("User should take all cards in table")
	}

	err = checkGameEvents(
		game,
		EVENT_ATTACK,
		EVENT_ATTACK,
		EVENT_ATTACK,
		EVENT_TAKE_ALL_CARDS,
		EVENT_END_ATTACK,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestDefendTimeOver(t *testing.T) {
	game, attackUser, defendUser, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	game.Settings.TimeOver = 0.0001

	attackUser.Cards[0] = Card{
		Suit: 1,
		Rank: 6,
	}

	attackCard := attackUser.Cards[0]
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackCard,
	}

	game.AttackHandler(attackC, attackUser)

	if !game.DefendTimerIsRunning {
		t.Error("Defend timer should be started")
	}

	time.Sleep(2 * time.Millisecond)
	defendUser.Cards[0] = Card{
		Suit: 1,
		Rank: 14,
	}

	defendC := DefendCommand{
		Command: Command{
			Action: ACTION_DEFEND,
			UserId: defendUser.Id,
		},
		UserCard:   defendUser.Cards[0],
		TargetCard: attackCard,
	}

	r := game.DefendHandler(defendC, defendUser)
	if r.Error != ERROR_DEFEND_TIME_OVER {
		t.Error("Defend timer should be time over. ERROR: " + r.Error)
	}

	err = checkGameEvents(
		game,
		EVENT_ATTACK,
		EVENT_END_ATTACK,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestAttackTimeOver(t *testing.T) {
	game, attackUser, _, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}
	game.Settings.TimeOver = 0.0001

	for i := range 2 {
		attackUser.Cards[i] = Card{
			Suit: 1,
			Rank: 6,
		}
	}

	time.Sleep(2 * time.Millisecond)

	attackCard := attackUser.Cards[0]
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackCard,
	}

	r := game.AttackHandler(attackC, attackUser)

	if r.Error != ERROR_ATTACK_TIME_OVER {
		t.Error("Attack timer should be time over. ERROR: " + r.Error)
	}

	err = checkGameEvents(game, EVENT_END_ATTACK)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckAttackTimerStatus(t *testing.T) {
	game, attackUser, _, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	game.CheckAttackTimerHandler(Command{Action: ACTION_CHECK_ATTACK_TIMER}, attackUser)
	attackTimerStateEvent := game.GameEventBuffer[0].(AttackTimerStateEvent)
	if attackTimerStateEvent.Completed == true {
		t.Error("Timer shouldn't be completed")
	}
}

func TestCheckDefendTimerStatus(t *testing.T) {
	game, attackUser, defendUser, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	attackCard := attackUser.Cards[0]
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackCard,
	}

	game.AttackHandler(attackC, attackUser)
	checkGameEvents(game, EVENT_ATTACK)

	game.CheckDefendTimerHandler(Command{Action: ACTION_CHECK_DEFEND_TIMER}, defendUser)
	defendTimerStateEvent := game.GameEventBuffer[0].(DefendTimerStateEvent)
	if defendTimerStateEvent.Completed == true {
		t.Error("Timer shouldn't be completed")
	}
}

func TestGameEnd(t *testing.T) {
	game, attackUser, defendUser, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	game.Deck = []Card{}
	attackUser.Cards = []Card{attackUser.Cards[0]}

	// ATTACK
	attackCard := attackUser.Cards[0]
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: attackUser.Id,
		},
		Card: attackCard,
	}

	game.AttackHandler(attackC, attackUser)

	// DEFEND
	defendUser.Cards[0] = Card{Suit: game.TrumpSuit, Rank: 15}
	defendC := DefendCommand{
		Command: Command{
			Action: ACTION_DEFEND,
			UserId: defendUser.Id,
		},
		UserCard:   defendUser.Cards[0],
		TargetCard: attackCard,
	}

	r := game.DefendHandler(defendC, defendUser)

	if r.Error != ERROR_EMPTY {
		t.Error("Defend error CODE: " + r.Error)
	}
	rsp := game.EndAttackHandler(Command{
		Action: ACTION_END_ATTACK,
		UserId: attackUser.Id,
	}, attackUser)

	if rsp.Error != ERROR_EMPTY {
		t.Error("End Attack Action end with error: " + rsp.Error)
	}

	err = checkGameEvents(game, EVENT_ATTACK, EVENT_DEFEND, EVENT_END_ATTACK, EVENT_END_GAME)
	if err != nil {
		t.Error(err)
	}
}
