package game

type State int

const (
	StateMenu State = iota
	StatePlaying
	StatePaused
	StateLevelComplete
	StateGameOver
)
