package ai

import "github.com/lunajones/apeiron/service/creature"

type SelectorNode struct {
	Children []BehaviorNode
}

func (s *SelectorNode) Tick(c *creature.Creature) BehaviorStatus {
	for _, child := range s.Children {
		if child.Tick(c) == StatusSuccess {
			return StatusSuccess
		}
	}
	return StatusFailure
}
