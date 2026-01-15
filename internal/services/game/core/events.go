package core

import "time"

const (
	EVENT_NONE                       = "NONE"
	EVENT_START                      = "START"
	EVENT_READY                      = "READY"
	EVENT_ATTACK                     = "ATTACK"
	EVENT_DEFEND                     = "DEFEND"
	EVENT_END_ATTACK                 = "END_ATTACK"
	EVENT_TAKE_ALL_CARDS             = "TAKE_ALL_CARDS"
	EVENT_ATTACK_TIMER_NOT_COMPLETED = "ATTACK_TIMER_NOT_COMPLETED"
	EVENT_DEFEND_TIMER_NOT_COMPLETED = "DEFEND_TIMER_NOT_COMPLETED"
	EVENT_ATTACK_TIMER_COMPLETED     = "ATTACK_TIMER_COMPLETED"
	EVENT_DEFEND_TIMER_COMPLETED     = "DEFEND_TIMER_COMPLETED"
	EVENT_USER_EXIT                  = "USER_EXIT" // connection loss
	EVENT_USER_HAS_FINISHED          = "USER_HAS_FINISHED"
	EVENT_END_GAME                   = "END_GAME"
)

type GameResult string

const (
	GameResultWin         GameResult = "win"
	GameResultDraw        GameResult = "draw"
	GameResultInterrupted GameResult = "interrupted"
)

type GameEventContainer any

// type GameEven

type GameEvent struct {
	Event string `json:"event"`
}

type StartGameEvent struct {
	GameEvent
}

type ReadyEvent struct {
	GameEvent
	UserId string `json:"user_id"`
}

type AttackEvent struct {
	GameEvent
	Card       Card   `json:"card"`
	AttackerId string `json:"attacker_id"`
}

type DefendEvent struct {
	GameEvent
	TargetCard Card   `json:"target_card"`
	UserCard   Card   `json:"user_card"`
	DefenderId string `json:"defender_id"`
}

type TakeAllCardsEvent struct {
	GameEvent
	UserId string `json:"user_id"`
}

type EndAttackEvent struct {
	GameEvent
}

type AttackTimerStateEvent struct {
	GameEvent
	Completed  bool       `json:"completed"`
	TimerEndAt *time.Time `json:"timer_end_at"`
}

type DefendTimerStateEvent struct {
	GameEvent
	Completed  bool       `json:"completed"`
	TimerEndAt *time.Time `json:"timer_end_at"`
}

type EndGameEvent struct {
	GameEvent
	GameResult GameResult `json:"game_result"`
}

func NewReadyEvent(userId string) ReadyEvent {
	return ReadyEvent{
		GameEvent: GameEvent{
			Event: EVENT_READY,
		},
		UserId: userId,
	}
}

func NewStartGameEvent(gameState GameStateResponse) StartGameEvent {
	return StartGameEvent{
		GameEvent: GameEvent{
			Event: EVENT_START,
		},
	}
}

func NewAttackEvent(card Card, attackerId string) AttackEvent {
	return AttackEvent{
		GameEvent: GameEvent{
			Event: EVENT_ATTACK,
		},
		Card:       card,
		AttackerId: attackerId,
	}
}

func NewDefendEvent(
	userCard Card,
	targetCard Card,
	defenderId string,
	gameState GameStateResponse,
) DefendEvent {
	return DefendEvent{
		GameEvent: GameEvent{
			Event: EVENT_DEFEND,
		},
		UserCard:   userCard,
		TargetCard: targetCard,
		DefenderId: defenderId,
	}
}

func NewTakeAllCardsEvent(userId string) TakeAllCardsEvent {
	return TakeAllCardsEvent{
		GameEvent: GameEvent{
			Event: EVENT_TAKE_ALL_CARDS,
		},
		UserId: userId,
	}
}

func NewEndAttackEvent() EndAttackEvent {
	return EndAttackEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_ATTACK,
		},
	}
}

func NewAttackTimerStateEvent(
	completed bool,
	timerEndAt *time.Time,
) AttackTimerStateEvent {
	return AttackTimerStateEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_ATTACK,
		},
		Completed:  completed,
		TimerEndAt: timerEndAt,
	}
}

func NewDefendTimerStateEvent(
	completed bool,
	timerEndAt *time.Time,
) DefendTimerStateEvent {
	return DefendTimerStateEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_ATTACK,
		},
		Completed:  completed,
		TimerEndAt: timerEndAt,
	}
}

func NewEndGameEvent(
	result GameResult,
) EndGameEvent { // TODO: handle many users. Now handle only 2
	return EndGameEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_GAME,
		},
		GameResult: result,
	}
}

func GameEventToType(e GameEventContainer) string {
	switch event := e.(type) {
	case ReadyEvent:
		return event.Event
	case StartGameEvent:
		return event.Event
	case AttackEvent:
		return event.Event
	case DefendEvent:
		return event.Event
	case EndAttackEvent:
		return event.Event
	case TakeAllCardsEvent:
		return event.Event
	case EndGameEvent:
		return event.Event
	case AttackTimerStateEvent:
		return event.Event
	case DefendTimerStateEvent:
		return event.Event
	}

	return EVENT_NONE
}
