package enum

type GameStatus int

const (
	GameStatusWaiting GameStatus = iota
	GameStatusInProgress
	GameStatusFinished
	GameStatusAbandoned
)
