package core

import (
	"github.com/lunajones/apeiron/service/creature"
)

type SequenceNode struct {
	Children []BehaviorNode
}

func NewSequenceNode(children ...BehaviorNode) *SequenceNode {
	return &SequenceNode{Children: children}
}

func (n *SequenceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	for _, child := range n.Children {
		status := child.Tick(c, ctx).(BehaviorStatus)
		if status != StatusSuccess {
			return status
		}
	}
	return StatusSuccess
}
