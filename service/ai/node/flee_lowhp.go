package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeIfLowHPNode struct{}

func (n *FleeIfLowHPNode) Tick(c *creature.Creature) core.BehaviorStatus {
	hpPercent := float64(c.HP) / float64(c.MaxHP)

	if hpPercent < 0.2 {
		c.SetAction(creature.ActionRun)
		c.ChangeAIState(creature.AIStateFleeing)
		return core.StatusSuccess
	}

	return core.StatusFailure
}
