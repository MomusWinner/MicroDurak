package core

import "time"

func (g *Game) checkIsAttacker(userId string) string {
	if userId != g.AttackingId {
		return ERROR_NOT_YOUR_TURN
	}
	return ERROR_EMPTY
}

func (g *Game) checkIsDefender(userId string) string {
	if userId != g.DefendingId {
		return ERROR_NOT_YOUR_TURN
	}
	return ERROR_EMPTY
}

func (g *Game) checkUserHasCard(user *User, card Card) string {
	card, err := getCardBySuitAndRank(user.Cards, card.Suit, card.Rank)
	if err != nil {
		return ERROR_USER_NO_HAS_CARD
	}
	return ERROR_EMPTY
}

func (g *Game) checkAttackTimer() string {
	if !g.AttackTimerIsRunning {
		return ERROR_EMPTY
	}
	now := time.Now()
	if now.Sub(g.AttackTimerStartedAt).Seconds() >= g.Settings.TimeOver {
		return ERROR_ATTACK_TIME_OVER
	}
	return ERROR_EMPTY
}

func (g *Game) checkDefendTimer() string {
	if !g.DefendTimerIsRunning {
		return ERROR_EMPTY
	}
	now := time.Now()
	if now.Sub(g.DefendTimerStartedAt).Seconds() >= g.Settings.TimeOver {
		return ERROR_DEFEND_TIME_OVER
	}
	return ERROR_EMPTY
}

func (g *Game) checkSameRankCard(rank int) string {
	if !tableHasCardRank(g.TableCards, rank) {
		return ERROR_NO_SAME_RANK_CARD_IN_TABLE
	}
	return ERROR_EMPTY
}

func (g *Game) checkCardOnTable(suit int, rank int) string {
	targetCardExist := tableHasCard(g.TableCards, suit, rank)

	if !targetCardExist {
		return ERROR_NOT_FOUND_CART_ON_TABLE
	}

	return ERROR_EMPTY
}

func (g *Game) checkCardGreater(suit int, rank int, tsuit int, trank int) string {
	greaterThanTarget := CardGreater(suit, rank, tsuit, trank, g.Trump.Suit)
	if !greaterThanTarget {
		return ERROR_TARGET_CARD_GREATER_THEN_YOUR
	}

	return ERROR_EMPTY
}

func (g *Game) checkGameStarted() string {
	if !g.IsStarted {
		return ERROR_GAME_SHOULD_BE_STARTED
	}

	return ERROR_EMPTY
}

func (g *Game) checkNotFirstTurn() string {
	if len(g.TableCards) == 0 {
		return ERROR_CANNOT_END_ATTACK_IN_FIRST_TURN
	}

	return ERROR_EMPTY
}

func (g *Game) checkAllCardsBeatOff() string {
	if !allCardBeatOff(g.TableCards) {
		return ERROR_ALL_CARD_SHOULD_BE_BEAT_OFF_BEFORE_END_ATTACK
	}

	return ERROR_EMPTY
}

func (g *Game) checkTableHoldsOnlySixCards() string {
	if len(g.TableCards) >= 6 {
		return ERROR_TABLE_HOLDS_ONLY_SIX_CARDS
	}

	return ERROR_EMPTY
}

func (g *Game) checkDefenderHasCards() string {
	defender, _ := getUserById(g.Users, g.DefendingId)
	if len(defender.Cards) <= 0 {
		return ERROR_DEFENDER_NO_CARDS
	}

	return ERROR_EMPTY
}

func errorChecker(gameErrors []string) string {
	for _, gameError := range gameErrors {
		if gameError != ERROR_EMPTY {
			return gameError
		}
	}

	return ERROR_EMPTY
}
