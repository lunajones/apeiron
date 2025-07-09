package helper

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type ConditionFunc func(c *creature.Creature, ctx interface{}) bool

type ConditionNode struct {
	condition ConditionFunc
}

func NewConditionNode(cond ConditionFunc) *ConditionNode {
	return &ConditionNode{
		condition: cond,
	}
}

func (n *ConditionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if n.condition(c, ctx) {
		return core.StatusSuccess
	}
	return core.StatusFailure
}

func (n *ConditionNode) Reset() {
	// Nada a resetar
}
