package helper

import (
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type OnlyIfFarAndNotMovingNode struct {
	Node core.BehaviorNode
}

func (n *OnlyIfFarAndNotMovingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		return core.StatusFailure
	}

	if c.MoveCtrl.IsMoving {
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	if dist <= 1.0 {
		return core.StatusFailure
	}

	return n.Node.Tick(c, ctx)
}

func (n *OnlyIfFarAndNotMovingNode) Reset(c *creature.Creature) {

}
