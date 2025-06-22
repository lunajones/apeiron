package core

import "github.com/lunajones/apeiron/service/creature"

type SequenceNode struct {
	Children []BehaviorNode
}

func (n *SequenceNode) Tick(c *creature.Creature) BehaviorStatus {
	for _, child := range n.Children {
		status := child.Tick(c)
		if status != StatusSuccess {
			return status
		}
	}
	return StatusSuccess
}
