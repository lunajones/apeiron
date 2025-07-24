package helper

import (
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type OnlyIfCloseAndNotMovingNode struct {
	Node core.BehaviorNode
}

func (n *OnlyIfCloseAndNotMovingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {

	if c.MoveCtrl.IsMoving {
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	if dist > 4.5 {
		return core.StatusFailure
	}

	return n.Node.Tick(c, ctx)
}

func (n *OnlyIfCloseAndNotMovingNode) Reset(c *creature.Creature) {

}
