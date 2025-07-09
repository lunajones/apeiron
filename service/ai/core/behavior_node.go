package core

import (
	"github.com/lunajones/apeiron/service/creature"
)

type BehaviorNode interface {
	Tick(c *creature.Creature, ctx interface{}) interface{}
	Reset()
}
