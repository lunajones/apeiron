package humanoid

import (
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type SneakBehindTargetNode struct{}

func (n *SneakBehindTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		return core.StatusFailure
	}

	svcCtx := c.GetContext()

	if c.MoveCtrl.MovementPlan != nil {
		return core.StatusRunning
	}

	dir := position.NewVector2DFromTo(target.GetPosition(), c.GetPosition()).Normalize() // vai pras costas
	offset := dir.Scale(1.5)
	dest := target.GetPosition().AddVector2D(offset)

	if !svcCtx.NavMesh.IsWalkable(dest) || !svcCtx.ClaimPosition(dest, c.Handle) {
		color.HiRed("[SNEAK] [%s] destino inválido ou ocupado", c.PrimaryType)
		return core.StatusFailure
	}

	c.MoveCtrl.SetMoveTarget(dest, c.RunSpeed*1.2, 0.1)
	c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
		consts.MovementPlanSneak,
		target.GetHandle(),
		1.5,
		1*time.Second,
		target.GetPosition(),
	)

	color.White("[SNEAK] [%s] correndo pra flanquear alvo", c.PrimaryType)
	return core.StatusRunning
}

func (n *SneakBehindTargetNode) Reset(c *creature.Creature) {
	color.Yellow("[SNEAK] [RESET] [%s] limpando plano", c.GetPrimaryType())
	c.MoveCtrl.MovementPlan = nil
	// FSM limpará o resto
}
