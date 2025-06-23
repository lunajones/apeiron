package core

import "github.com/lunajones/apeiron/service/creature"

type BehaviorStatus int

const (
	StatusSuccess BehaviorStatus = iota
	StatusFailure
	StatusRunning
)

type BehaviorNode interface {
	Tick(c *creature.Creature, ctx AIContext) BehaviorStatus
}