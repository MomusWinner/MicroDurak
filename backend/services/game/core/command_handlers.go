package core

import "time"

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
			Error:     gameError,
			Command:   command,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	g.EndAttackUserId = append(g.EndAttackUserId, command.UserId)

	defendUser, _ := getUserById(g.Users, g.DefendingId)
	attackUser, _ := getUserById(g.Users, g.AttackingId)

	if len(g.EndAttackUserId) == len(g.Users)-1 {
		g.EndAttack()
	}

	if len(g.Deck) == 0 && len(defendUser.Cards) == 0 && len(attackUser.Cards) == 0 {
		g.AddEventToBuffer(NewEndGameEvent(GameResultDraw))
	} else if len(g.Deck) == 0 && len(defendUser.Cards) == 0 {
		g.AddEventToBuffer(NewEndGameEvent(GameResultWin))
	} else if len(g.Deck) == 0 && len(attackUser.Cards) == 0 {
		g.AddEventToBuffer(NewEndGameEvent(GameResultWin))
	}

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   command,
		GameState: gameToGameStateResponse(g, *user),
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
		g.EndAttack()
	}

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:     gameError,
			Command:   attackCommand,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	tableCard := TableCard{}
	tableCard.Suit = attackCommand.Card.Suit
	tableCard.Rank = attackCommand.Card.Rank

	g.TableCards = append(g.TableCards, tableCard)
	err := g.removeUserCard(user.Id, tableCard.Suit, tableCard.Rank)
	if err != nil {
		return CommandResponse{
			Error:     ERROR_SERVER,
			Command:   attackCommand,
			GameState: gameToGameStateResponse(g, *user),
		}
	}
	g.EndAttackUserId = make([]string, 0)
	g.StartDefendTimer()
	g.AddEventToBuffer(NewAttackEvent(attackCommand.Card, attackCommand.UserId))

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   attackCommand,
		GameState: gameToGameStateResponse(g, *user),
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

	if gameError == ERROR_DEFEND_TIME_OVER {
		g.AddEventToBuffer(NewEndAttackEvent())
	}

	if gameError != ERROR_EMPTY {
		return CommandResponse{
			Error:     gameError,
			Command:   defendCommand,
			GameState: gameToGameStateResponse(g, *user),
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
		),
	)

	if allCardBeatOff(g.TableCards) && len(g.TableCards) < 6 {
		g.StartAttackTimer()
	}

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   defendCommand,
		GameState: gameToGameStateResponse(g, *user),
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
			Error:     gameError,
			Command:   command,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	tableCards := tableCardsToCards(g.TableCards)
	user.Cards = append(user.Cards, tableCards...)
	g.TableCards = []TableCard{}
	g.StopDefendTimer()

	g.AddEventToBuffer(NewTakeAllCardsEvent(user.Id))
	g.AddEventToBuffer(NewEndAttackEvent())

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   command,
		GameState: gameToGameStateResponse(g, *user),
	}
}

func (g *Game) ReadyHandler(command Command, user *User) CommandResponse {
	if contains(g.ReadyUsers, user.Id) {
		return CommandResponse{
			Error:     ERROR_USER_ALREADY_READY,
			Command:   command,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	g.ReadyUsers = append(g.ReadyUsers, user.Id)

	if len(g.ReadyUsers) >= len(g.Users) {
		g.IsStarted = true
		g.StartAttackTimer()
		g.AddEventToBuffer(NewStartGameEvent())
	} else {
		g.AddEventToBuffer(NewReadyEvent(user.Id))
	}

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   command,
		GameState: gameToGameStateResponse(g, *user),
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
			Error:     gameError,
			Command:   command,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	if g.AttackTimerIsRunning {
		timeEndAt := g.AttackTimerStartedAt.Add(time.Duration(g.Settings.TimeOver) * time.Second)
		g.AddEventToBuffer(NewAttackTimerStateEvent(false, &timeEndAt))
	} else {
		g.AddEventToBuffer(NewAttackTimerStateEvent(true, nil))
	}

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   command,
		GameState: gameToGameStateResponse(g, *user),
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
			Error:     gameError,
			Command:   command,
			GameState: gameToGameStateResponse(g, *user),
		}
	}

	if g.DefendTimerIsRunning {
		timeEndAt := g.DefendTimerStartedAt.Add(time.Duration(g.Settings.TimeOver) * time.Second)
		g.AddEventToBuffer(NewDefendTimerStateEvent(false, &timeEndAt))
	} else {
		g.AddEventToBuffer(NewAttackTimerStateEvent(true, nil))
	}

	return CommandResponse{
		Error:     ERROR_EMPTY,
		Command:   command,
		GameState: gameToGameStateResponse(g, *user),
	}
}
