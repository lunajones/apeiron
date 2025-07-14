package core

import (
	"github.com/lunajones/apeiron/service/creature"
)

type ParallelNode struct {
	children []BehaviorNode
}

func NewParallelNode(children ...BehaviorNode) *ParallelNode {
	return &ParallelNode{
		children: children,
	}
}

func (n *ParallelNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	allSuccess := true
	anyRunning := false

	for _, child := range n.children {
		result := child.Tick(c, ctx)
		switch result {
		case StatusFailure:
			return StatusFailure // Se algum falhar, falha tudo
		case StatusRunning:
			anyRunning = true
			allSuccess = false
		case StatusSuccess:
			// Continua verificando os outros
		default:
			allSuccess = false
		}
	}

	if anyRunning {
		return StatusRunning
	}
	if allSuccess {
		return StatusSuccess
	}

	return StatusRunning
}

func (n *ParallelNode) Reset(c *creature.Creature) {
	for _, child := range n.children {
		child.Reset(c)
	}
}
