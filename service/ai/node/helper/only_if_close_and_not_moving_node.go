package helper

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type OnlyIfNotMovingNode struct {
	Node core.BehaviorNode
}

func (n *OnlyIfNotMovingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.MoveCtrl.IsMoving {
		return core.StatusFailure
	}
	return n.Node.Tick(c, ctx)
}

func (n *OnlyIfNotMovingNode) Reset(c *creature.Creature) {

}
