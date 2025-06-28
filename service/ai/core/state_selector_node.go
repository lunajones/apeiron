package core

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type StateSelectorNode struct {
	Subtrees map[consts.AIState]BehaviorNode
}

func NewStateSelectorNode() *StateSelectorNode {
	return &StateSelectorNode{
		Subtrees: make(map[consts.AIState]BehaviorNode),
	}
}

func (n *StateSelectorNode) AddSubtree(state consts.AIState, subtree BehaviorNode) {
	n.Subtrees[state] = subtree
}

func (n *StateSelectorNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	subtree, ok := n.Subtrees[c.AIState]
	if !ok {
		return StatusFailure
	}
	return subtree.Tick(c, ctx)
}
