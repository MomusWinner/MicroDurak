package core

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func nextPlace(place int, usersLength int) int {
	place += 1
	if place >= usersLength {
		place = 0
	}

	return place
}

func (g *Game) ImproveUserFirstCard(user *User) error {
	if len(user.Cards) == 0 {
		return errors.New("The user hasn't cards left.")
	}

	user.Cards[0] = Card{g.TrumpSuit, 15}

	return nil
}

func checkCommandResponse(response CommandResponse) error {
	if response.Error != ERROR_EMPTY {
		return errors.New(response.Error)
	}

	return nil
}

func (g *Game) SendAttackCommand(user *User, card Card) CommandResponse {
	attackC := AttackCommand{
		Command: Command{
			Action: ACTION_ATTACK,
			UserId: user.Id,
		},
		Card: card,
	}

	return g.AttackHandler(attackC, user)
}

func (g *Game) SendAttackCommandSafe(user *User, card Card) error {
	return checkCommandResponse(g.SendAttackCommand(user, card))
}

func (g *Game) SendDefendCommand(user *User, target_card Card, defend_card Card) CommandResponse {
	defendC := DefendCommand{
		Command: Command{
			Action: ACTION_DEFEND,
			UserId: user.Id,
		},
		UserCard:   defend_card,
		TargetCard: target_card,
	}

	return g.DefendHandler(defendC, user)
}

func (g *Game) SendDefendCommandSafe(user *User, target_card Card, defend_card Card) error {
	return checkCommandResponse(g.SendDefendCommand(user, target_card, defend_card))
}

func (g *Game) SendEndAttackCommand(user *User) CommandResponse {
	return g.EndAttackHandler(
		Command{
			Action: ACTION_END_ATTACK,
			UserId: user.Id,
		}, user)
}

func (g *Game) SendEndAttackCommandSafe(user *User) error {
	return checkCommandResponse(g.SendEndAttackCommand(user))
}

func (g *Game) SendTakeAllCardsCommand(user *User) CommandResponse {
	return g.TakeAllCardHandler(
		Command{
			Action: ACTION_TAKE_ALL_CARDS,
			UserId: user.Id,
		}, user)
}

func (g *Game) SendTakeAllCardsCommandSafe(user *User) error {
	return checkCommandResponse(g.SendTakeAllCardsCommand(user))
}

func checkGameEvents(game *Game, events ...string) error {
	var errMsg string = ""

	if len(events) != len(game.GameEventBuffer) {
		errMsg += fmt.Sprintf(
			"\nTarget events size (%d) does not match game events size (%d)\n",
			len(events),
			len(game.GameEventBuffer),
		)
	}

	if errMsg == "" {
		for i, v := range game.GameEventBuffer {
			eventType := GameEventToType(v)

			if events[i] != eventType {
				errMsg += fmt.Sprintf(
					"Bad event: expected type %s, got %s at index %d\n",
					eventType,
					events[i],
					i,
				)
				break
			}
		}
	}

	if errMsg != "" {
		errMsg += "\nTarget events:\n"
		for _, e := range events {
			errMsg += " - " + e + "\n"
		}

		errMsg += "\nGame events:\n"

		for _, e := range game.GameEventBuffer {
			errMsg += " - " + GameEventToType(e) + "\n"
		}

		return errors.New(errMsg)
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

	err := checkGameEvents(game, EVENT_READY, EVENT_READY, EVENT_START)

	return game, attackUser, defendUser, err
}

func createGameWithReadyUsers3() (*Game, *User, *User, *User, error) {
	game, _ := CreateNewGame([]string{"user1", "user2", "user3"})

	attacking, _ := game.getUserById(game.AttackingId)
	defending, _ := game.getUserById(game.DefendingId)
	observing := game.getObservingUsers()[0]

	game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: attacking.Id,
	}, attacking)

	game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: defending.Id,
	}, defending)

	game.ReadyHandler(Command{
		Action: ACTION_READY,
		UserId: observing.Id,
	}, observing)

	err := checkGameEvents(game, EVENT_READY, EVENT_READY, EVENT_READY, EVENT_START)

	return game, attacking, defending, observing, err
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
	game, attacker, defender, err := createGameWithReadyUsers()

	if err != nil {
		t.Error(err)
	}

	// ATTACK
	attackCard := attacker.Cards[0]

	err = game.SendAttackCommandSafe(attacker, attackCard)
	if err != nil {
		t.Error(err)
	}

	if len(attacker.Cards) != 5 {
		t.Error("The card was not removed after the attack.")
	}

	if game.TableCards[0].Suit != attackCard.Suit || game.TableCards[0].Rank != attackCard.Rank {
		t.Error("Table shouldn't be empty after attack")
	}

	// DEFEND
	err = game.ImproveUserFirstCard(defender)
	if err != nil {
		t.Error(err)
	}

	err = game.SendDefendCommandSafe(defender, attackCard, defender.Cards[0])
	if err != nil {
		t.Error(err)
	}

	// END ATTACK
	err = game.SendEndAttackCommandSafe(attacker)
	if err != nil {
		t.Error(err)
	}

	if len(attacker.Cards) != 6 || len(defender.Cards) != 6 {
		t.Error("Atfter AttackEndCommand game server should hand out the cards.")
	}

	if game.DefendingId != attacker.Id || game.AttackingId != defender.Id {
		t.Error("After AttackEndCommand game should change user roles.")
	}

	if len(game.TableCards) != 0 {
		t.Error("After AttackEndCommand should be emtpy")
	}

	err = checkGameEvents(game, EVENT_ATTACK, EVENT_DEFEND, EVENT_END_ATTACK)
	if err != nil {
		t.Error(err)
	}
}

func TestTakeAllCards(t *testing.T) {
	game, attacker, defender, observer, err := createGameWithReadyUsers3()

	if err != nil {
		t.Error(err)
	}

	for i := range attacker.Cards {
		attacker.Cards[i] = Card{
			Suit: 1,
			Rank: 6,
		}
	}

	for range 3 {
		err = game.SendAttackCommandSafe(attacker, attacker.Cards[0])
		if err != nil {
			t.Error(err)
		}
	}

	// Beat off first card on table
	game.TableCards[0].BeatOff = &Card{
		Suit: 1,
		Rank: 14,
	}

	err = game.SendTakeAllCardsCommandSafe(defender)
	if err != nil {
		t.Error(err)
	}

	if len(defender.Cards) != 10 {
		t.Error("User should take all cards in table")
	}

	if game.AttackingId == attacker.Id {
		t.Errorf("Hmm %s", observer.Id)
	}
	if game.AttackingId != observer.Id {
		t.Errorf("Expected next attacker to be user %s", observer.Id)
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

func TestAttackCardsOnTableHaveSameRank(t *testing.T) {
	game, attackUser, _, err := createGameWithReadyUsers()

	attackUser.Cards[0] = Card{Suit: 1, Rank: 6}
	attackUser.Cards[1] = Card{Suit: 1, Rank: 7}

	// ATTACK
	attackCard := attackUser.Cards[0]

	err = game.SendAttackCommandSafe(attackUser, attackCard)
	if err != nil {
		t.Error(err)
	}

	attackCard = attackUser.Cards[0]

	response := game.SendAttackCommand(attackUser, attackCard)
	if response.Error != ERROR_NO_SAME_RANK_CARD_IN_TABLE {
		t.Errorf("Expected CommandResponse error (%s) when user attacks with an invalid card rank", ERROR_NO_SAME_RANK_CARD_IN_TABLE)
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

func Test3UserGame(t *testing.T) {
	game, a, d, o, err := createGameWithReadyUsers3()
	_, _, _, _ = game, a, d, o

	fmt.Printf("Attacker: %s\nDefender: %s\nObserver: %s\n ", a.Id, d.Id, o.Id)

	if err != nil {
		t.Error(err)
	}

	// ATTACK
	attackCard := a.Cards[0]

	err = game.SendAttackCommandSafe(a, attackCard)
	if err != nil {
		t.Error(err)
	}

	// DEFEND
	err = game.ImproveUserFirstCard(d)
	if err != nil {
		t.Error(err)
	}

	game.SendDefendCommand(d, attackCard, d.Cards[0])
	if err != nil {
		t.Error(err)
	}

	startAttackerPlace := a.Place
	startDeffenderPlace := d.Place
	// startObserverPlace := d.Place

	// ATTACKER END ATTCK
	err = game.SendEndAttackCommandSafe(a)
	if err != nil {
		t.Error(err)
	}

	// OBSERVER END ATTCK
	err = game.SendEndAttackCommandSafe(o)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("\nNew Attacker: %s\nNew Defender: %s\n", game.AttackingId, game.DefendingId)

	next_attacker, err := game.getUserByPlace(nextPlace(startAttackerPlace, len(game.Users)))
	if err != nil {
		t.Error(err)
	}

	if next_attacker.Id != d.Id {
		t.Error("Error")
	}

	next_defender, err := game.getUserByPlace(nextPlace(startDeffenderPlace, len(game.Users)))
	if err != nil {
		t.Error(err)
	}

	if next_defender.Id != o.Id {
		t.Error("Error")
	}

	if a.Id == game.AttackingId || a.Id == game.DefendingId {
		t.Error("Last attacker should be observer")
	}
}
