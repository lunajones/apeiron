package core

import (
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type SequenceNode struct {
	Children []BehaviorNode
}

func (n *SequenceNode) Tick(c *creature.Creature, ctx ai_context.AIContext) interface{} {
	for _, child := range n.Children {
		status := child.Tick(c, ctx).(BehaviorStatus)
		if status != StatusSuccess {
			return status
		}
	}
	return StatusSuccess
}
