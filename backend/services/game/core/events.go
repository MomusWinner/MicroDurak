package core

import "time"

const (
	EVENT_START                      = "EVENT_START"
	EVENT_READY                      = "EVENT_READY"
	EVENT_ATTACK                     = "EVENT_ATTACK"
	EVENT_DEFEND                     = "EVENT_DEFEND"
	EVENT_END_ATTACK                 = "EVENT_END_ATTACK"
	EVENT_TAKE_ALL_CARDS             = "EVENT_TAKE_ALL_CARDS"
	EVENT_ATTACK_TIMER_NOT_COMPLETED = "ATTACK_TIMER_NOT_COMPLETED"
	EVENT_DEFEND_TIMER_NOT_COMPLETED = "DEFEND_TIMER_NOT_COMPLETED"
	EVENT_ATTACK_TIMER_COMPLETED     = "ATTACK_TIMER_COMPLETED"
	EVENT_DEFEND_TIMER_COMPLETED     = "DEFEND_TIMER_COMPLETED"
	EVENT_END_GAME                   = "END_GAME"
)

type GameResult string

const (
	GameResultWin         GameResult = "win"
	GameResultDraw        GameResult = "draw"
	GameResultInterrupted GameResult = "interrupted"
)

type GameEventContainer interface {
	GetGameEvent() *GameEvent
}

type EventPack struct {
	Events []GameEventContainer `json:"events"`
}

type GameEvent struct {
	Event string            `json:"event"`
	State GameStateResponse `json:"state"`
}

type StartGameEvent struct {
	GameEvent
}

type ReadyEvent struct {
	UserId string `json:"user_id"`
	GameEvent
}

type AttackEvent struct {
	Card       Card   `json:"card"`
	AttackerId string `json:"attacker_id"`
	GameEvent
}

type DefendEvent struct {
	TargetCard Card   `json:"target_card"`
	UserCard   Card   `json:"user_card"`
	DefenderId string `json:"defender_id"`
	GameEvent
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

func NewStartGameEvent() StartGameEvent {
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

func NewDefendEvent(userCard Card, targetCard Card, defenderId string) DefendEvent {
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

func NewAttackTimerStateEvent(completed bool, timerEndAt *time.Time) AttackTimerStateEvent {
	return AttackTimerStateEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_ATTACK,
		},
		Completed:  completed,
		TimerEndAt: timerEndAt,
	}
}

func NewDefendTimerStateEvent(completed bool, timerEndAt *time.Time) DefendTimerStateEvent {
	return DefendTimerStateEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_ATTACK,
		},
		Completed:  completed,
		TimerEndAt: timerEndAt,
	}
}

func NewEndGameEvent(result GameResult) EndGameEvent { // TODO: handle many users. Now handle only 2
	return EndGameEvent{
		GameEvent: GameEvent{
			Event: EVENT_END_GAME,
		},
		GameResult: result,
	}
}

func (e GameEvent) GetGameEvent() *GameEvent             { return &e }
func (e StartGameEvent) GetGameEvent() *GameEvent        { return &e.GameEvent }
func (e ReadyEvent) GetGameEvent() *GameEvent            { return &e.GameEvent }
func (e AttackEvent) GetGameEvent() *GameEvent           { return &e.GameEvent }
func (e DefendEvent) GetGameEvent() *GameEvent           { return &e.GameEvent }
func (e TakeAllCardsEvent) GetGameEvent() *GameEvent     { return &e.GameEvent }
func (e EndAttackEvent) GetGameEvent() *GameEvent        { return &e.GameEvent }
func (e AttackTimerStateEvent) GetGameEvent() *GameEvent { return &e.GameEvent }
func (e DefendTimerStateEvent) GetGameEvent() *GameEvent { return &e.GameEvent }
