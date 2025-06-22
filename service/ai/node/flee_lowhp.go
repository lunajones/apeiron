package node

import "github.com/lunajones/apeiron/service/creature"

type FleeIfLowHPNode struct{}

func (f *FleeIfLowHPNode) Tick(c *creature.Creature) BehaviorStatus {
	if c.HP < 30 {
		c.SetAction(creature.ActionRun)
		c.ChangeAIState(creature.AIStateIdle)
		return StatusSuccess
	}
	return StatusFailure
}
