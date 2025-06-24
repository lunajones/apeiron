package core

import (
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type BehaviorNode interface {
	Tick(c *creature.Creature, ctx ai_context.AIContext) interface{}
}
