package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type DefendRecentNode struct{}

func (n *DefendRecentNode) Tick(c *creature.Creature) core.BehaviorStatus {
	if c.WasRecentlyAttacked() {
		c.SetAction(creature.ActionBlock)
		c.ChangeAIState(creature.AIStateDefend)
		return core.StatusSuccess
	}
	return core.StatusFailure
}
