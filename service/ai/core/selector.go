package core

import "github.com/lunajones/apeiron/service/creature"

type SelectorNode struct {
	Children []BehaviorNode
}

func (n *SelectorNode) Tick(c *creature.Creature, ctx AIContext) BehaviorStatus {
	for _, child := range n.Children {
		status := child.Tick(c, ctx)
		if status == StatusSuccess || status == StatusRunning {
			return status
		}
	}
	return StatusFailure
}