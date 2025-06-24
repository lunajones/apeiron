package core

import (
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type SelectorNode struct {
	Children []BehaviorNode
}

func (n *SelectorNode) Tick(c *creature.Creature, ctx ai_context.AIContext) interface{} {
	for _, child := range n.Children {
		status := child.Tick(c, ctx).(BehaviorStatus)
		if status == StatusSuccess {
			return StatusSuccess
		}
	}
	return StatusFailure
}
