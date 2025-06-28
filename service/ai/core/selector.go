package core

import (
	"github.com/lunajones/apeiron/service/creature"
)

type SelectorNode struct {
	Children []BehaviorNode
}

func NewSelectorNode(children ...BehaviorNode) *SelectorNode {
	return &SelectorNode{Children: children}
}

func (n *SelectorNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	for _, child := range n.Children {
		status := child.Tick(c, ctx).(BehaviorStatus)
		if status == StatusSuccess {
			return StatusSuccess
		}
	}
	return StatusFailure
}

func (n *SelectorNode) Reset() {
	// Se quiser, pode iterar e resetar os filhos também
	for _, child := range n.Children {
		child.Reset()
	}
}
