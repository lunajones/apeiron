package helper

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type ActionFunc func(c *creature.Creature, ctx interface{}) core.BehaviorStatus

type ActionNode struct {
	action ActionFunc
}

func NewActionNode(action ActionFunc) *ActionNode {
	return &ActionNode{action: action}
}

func (n *ActionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	return n.action(c, ctx)
}

func (n *ActionNode) Reset() {
	// Nada a resetar
}
