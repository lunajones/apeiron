package core

import "github.com/lunajones/apeiron/service/creature"

type SequenceNode struct {
	Children []BehaviorNode
}

func (s *SequenceNode) Tick(c *creature.Creature, ctx AIContext) BehaviorStatus {
	for _, child := range s.Children {
		status := child.Tick(c, ctx)
		if status != StatusSuccess {
			return status
		}
	}
	return StatusSuccess
}
