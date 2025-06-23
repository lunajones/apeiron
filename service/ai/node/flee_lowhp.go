package node

import (
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeIfLowHPNode struct{}

func (n *FleeIfLowHPNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando FleeIfLowHPNode", c.ID)

	hpThreshold := 30
	if c.HP < hpThreshold {
		c.ChangeAIState(creature.AIStateFleeing)
		return core.StatusSuccess
	}
	return core.StatusFailure
}