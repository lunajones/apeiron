package core

type BehaviorStatus int

const (
	StatusSuccess BehaviorStatus = iota
	StatusFailure
	StatusRunning
)
