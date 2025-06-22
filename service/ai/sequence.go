package ai

import "github.com/lunajones/apeiron/service/creature"

type SequenceNode struct {
	Children []BehaviorNode
}

func (s *SequenceNode) Tick(c *creature.Creature) BehaviorStatus {
	for _, child := range s.Children {
		if child.Tick(c) != StatusSuccess {
			return StatusFailure
		}
	}
	return StatusSuccess
}
