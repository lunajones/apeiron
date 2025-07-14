package helper

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

// AlwaysSuccessNode executa o node interno mas sempre retorna Success.
type AlwaysSuccessNode struct {
	inner core.BehaviorNode
}

// Tick executa o node interno mas ignora o status e sempre retorna Success.
func (n *AlwaysSuccessNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if n.inner != nil {
		n.inner.Tick(c, ctx)
	}
	return core.StatusSuccess
}

// Reset reseta o node interno, se houver.
func (n *AlwaysSuccessNode) Reset(c *creature.Creature) {
	if n.inner != nil {
		n.inner.Reset(c)
	}
}

// NewAlwaysSuccessNode cria um AlwaysSuccessNode envolvendo o node fornecido.
func NewAlwaysSuccessNode(inner core.BehaviorNode) core.BehaviorNode {
	return &AlwaysSuccessNode{inner: inner}
}
