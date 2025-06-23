package node

import (
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeIfLowHPNode struct{}

func (n *FleeIfLowHPNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	hpThreshold := 30
	if c.HP < hpThreshold {
		c.ChangeAIState(creature.AIStateFlee)
		return core.StatusSuccess
	}
	return core.StatusFailure
}