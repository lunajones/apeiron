package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeLowHPNode struct{}

func (n *FleeLowHPNode) Tick(c *creature.Creature) core.BehaviorStatus {
	hpPercent := float64(c.HP) / float64(c.MaxHP)

	if hpPercent < 0.2 {
		c.SetAction(creature.ActionRun)
		c.ChangeAIState(creature.AIStateFlee)
		return core.StatusSuccess
	}

	return core.StatusFailure
}
