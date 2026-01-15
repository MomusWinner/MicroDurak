package core

import (
	"time"
)

// Attacker end attack
func (g *Game) EndAttackHandler(command Command, user *User) CommandResponse {
	gameErrors := []string{
		g.checkNotFirstTurn(),
		g.checkGameStarted(),
		g.checkIsAttacker(command.UserId),
		g.checkAttackTimer(),
		g.checkAllCardsBeatOff(),
	}

	gameError := errorChecker(gameErrors)

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	g.EndAttackUserId = append(g.EndAttackUserId, command.UserId)

	defendUser, _ := g.getUserById(g.DefendingId)
	attackUser, _ := g.getUserById(g.AttackingId)

	if len(g.EndAttackUserId) == len(g.Users)-1 {
		g.EndAttack(true)

		if len(g.Deck) == 0 && len(defendUser.Cards) == 0 && len(attackUser.Cards) == 0 {
			g.AddEventToBuffer(NewEndGameEvent(GameResultDraw))
		} else if len(g.Deck) == 0 && len(defendUser.Cards) == 0 {
			g.AddEventToBuffer(NewEndGameEvent(GameResultWin))
		} else if len(g.Deck) == 0 && len(attackUser.Cards) == 0 {
			g.AddEventToBuffer(NewEndGameEvent(GameResultWin))
		}
	}

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: command,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) AttackHandler(attackCommand AttackCommand, user *User) CommandResponse {
	gameErrors := []string{
		g.checkGameStarted(),
		g.checkIsAttacker(attackCommand.UserId),
		g.checkUserHasCard(user, attackCommand.Card),
		g.checkAttackTimer(),
		g.checkDefenderHasCards(),
		g.checkTableHoldsOnlySixCards(),
	}

	if len(g.TableCards) != 0 {
		// It is correct if this is the first card in the table or
		// if there is a card of the same rank in the table
		gameErrors = append(gameErrors, g.checkSameRankCard(attackCommand.Card.Rank))
	}

	gameError := errorChecker(gameErrors)
	if gameError == ERROR_ATTACK_TIME_OVER {
		g.EndAttack(true)
	}

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: attackCommand,
			State:   gameToGameStateResponse(g, user),
		}
	}

	tableCard := TableCard{}
	tableCard.Suit = attackCommand.Card.Suit
	tableCard.Rank = attackCommand.Card.Rank

	g.TableCards = append(g.TableCards, tableCard)
	err := g.removeUserCard(user.Id, tableCard.Suit, tableCard.Rank)
	if err != nil {
		return CommandResponse{
			Error:   ERROR_SERVER,
			Command: attackCommand,
			State:   gameToGameStateResponse(g, user),
		}
	}
	g.EndAttackUserId = make([]string, 0)
	g.StartDefendTimer()
	g.AddEventToBuffer(
		NewAttackEvent(attackCommand.Card, attackCommand.UserId),
	)

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: attackCommand,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) DefendHandler(defendCommand DefendCommand, user *User) CommandResponse {
	gameError := errorChecker(
		[]string{
			g.checkGameStarted(),
			g.checkDefendTimer(),
			g.checkIsDefender(defendCommand.UserId),
			g.checkUserHasCard(user, defendCommand.UserCard),
			g.checkCardOnTable(defendCommand.TargetCard.Suit, defendCommand.TargetCard.Rank),
			g.checkCardGreater(
				defendCommand.UserCard.Suit,
				defendCommand.UserCard.Rank,
				defendCommand.TargetCard.Suit,
				defendCommand.TargetCard.Rank,
			),
		},
	)

	if gameError == ERROR_DEFEND_TIME_OVER { // TODO: remove this stuff
		g.AddEventToBuffer(NewEndAttackEvent())
	}

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: defendCommand,
			State:   gameToGameStateResponse(g, user),
		}
	}

	g.removeUserCard(user.Id, defendCommand.UserCard.Suit, defendCommand.UserCard.Rank)
	g.beatOffCard(
		defendCommand.UserCard.Suit,
		defendCommand.UserCard.Rank,
		defendCommand.TargetCard,
	)

	g.AddEventToBuffer(
		NewDefendEvent(
			defendCommand.UserCard,
			defendCommand.TargetCard,
			user.Id,
			gameToGameStateResponse(g, user),
		),
	)

	if allCardBeatOff(g.TableCards) && len(g.TableCards) < 6 {
		g.StartAttackTimer()
	}

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: defendCommand,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) TakeAllCardHandler(command Command, user *User) CommandResponse {
	gameError := errorChecker(
		[]string{
			g.checkGameStarted(),
			g.checkIsDefender(command.UserId),
		},
	)

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	tableCards := tableCardsToCards(g.TableCards)
	user.Cards = append(user.Cards, tableCards...)
	g.TableCards = []TableCard{}

	g.AddEventToBuffer(NewTakeAllCardsEvent(user.Id))

	g.EndAttack(false)

	newAttacker, _ := g.nextUser(g.DefendingId)
	newDefender, _ := g.nextUser(newAttacker.Id)

	g.AttackingId = newAttacker.Id
	g.DefendingId = newDefender.Id

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: command,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) ReadyHandler(command Command, user *User) CommandResponse {
	if contains(g.ReadyUsers, user.Id) {
		return CommandResponse{
			Error:   ERROR_USER_ALREADY_READY,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	g.ReadyUsers = append(g.ReadyUsers, user.Id)

	if len(g.ReadyUsers) >= len(g.Users) {
		g.IsStarted = true
		g.StartAttackTimer()
		g.AddEventToBuffer(NewReadyEvent(user.Id))
		g.AddEventToBuffer(NewStartGameEvent(gameToGameStateResponse(g, user)))
	} else {
		g.AddEventToBuffer(NewReadyEvent(user.Id))
	}

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: command,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) CheckAttackTimerHandler(command Command, user *User) CommandResponse {
	gameError := errorChecker(
		[]string{
			g.checkGameStarted(),
		},
	)

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	if g.checkAttackTimer() == ERROR_ATTACK_TIME_OVER {
		g.EndAttack(true)
	}

	if g.AttackTimerIsRunning {
		timeEndAt := g.AttackTimerStartedAt.Add(time.Duration(g.Settings.TimeOver) * time.Second)
		g.AddEventToBuffer(
			NewAttackTimerStateEvent(false, &timeEndAt),
		)
	} else {
		g.AddEventToBuffer(NewAttackTimerStateEvent(true, nil))
	}

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: command,
		State:   gameToGameStateResponse(g, user),
	}
}

func (g *Game) CheckDefendTimerHandler(command Command, user *User) CommandResponse {
	gameError := errorChecker(
		[]string{
			g.checkGameStarted(),
		},
	)

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:   gameError,
			Command: command,
			State:   gameToGameStateResponse(g, user),
		}
	}

	if g.checkDefendTimer() == ERROR_DEFEND_TIME_OVER {
		g.EndAttack(true)
	}

	if g.DefendTimerIsRunning {
		timeEndAt := g.DefendTimerStartedAt.Add(time.Duration(g.Settings.TimeOver) * time.Second)
		g.AddEventToBuffer(
			NewDefendTimerStateEvent(false, &timeEndAt),
		)
	} else {
		g.AddEventToBuffer(NewDefendTimerStateEvent(true, nil))
	}

	return CommandResponse{
		Error:   ERROR_EMPTY,
		Command: command,
		State:   gameToGameStateResponse(g, user),
	}
}
