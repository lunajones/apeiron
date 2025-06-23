package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomIdleNode struct{}

func (n *RandomIdleNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	c.SetAction(creature.ActionIdle)
	return core.StatusSuccess
}
