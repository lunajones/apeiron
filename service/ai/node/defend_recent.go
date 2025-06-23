package node

import (
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type DefendRecentNode struct{}

func (n *DefendRecentNode) Tick(c *creature.Creature) core.BehaviorStatus {
	log.Printf("[AI] %s executando DefendRecentNode", c.ID)

	if c.WasRecentlyAttacked() {
		c.SetAction(creature.ActionBlock)
		c.ChangeAIState(creature.AIStateDefending)
		return core.StatusSuccess
	}
	return core.StatusFailure
}
