package core

import (
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/creature"
)

type StateSelectorNode struct {
	Subtrees map[constslib.AIState]BehaviorNode
}

func NewStateSelectorNode() *StateSelectorNode {
	return &StateSelectorNode{
		Subtrees: make(map[constslib.AIState]BehaviorNode),
	}
}

func (n *StateSelectorNode) AddSubtree(state constslib.AIState, subtree BehaviorNode) {
	n.Subtrees[state] = subtree
}

func (n *StateSelectorNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	subtree, ok := n.Subtrees[c.AIState]
	if !ok {
		return StatusFailure
	}
	return subtree.Tick(c, ctx)
}
