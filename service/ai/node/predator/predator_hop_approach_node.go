package predator

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

type PredatorHopApproachNode struct{}

func (n *PredatorHopApproachNode) Tick(c *creature.Creature, _ interface{}) interface{} {
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		return core.StatusFailure
	}

	svcCtx := c.GetContext()

	// Não reexecuta se já tiver um plano de movimento
	if c.MoveCtrl.MovementPlan != nil {
		return core.StatusRunning
	}

	dir := position.NewVector2DFromTo(c.GetPosition(), target.GetPosition()).Normalize()
	hop := dir.Scale(1.8) // salto curto
	dest := c.GetPosition().AddVector2D(hop)

	if !svcCtx.NavMesh.IsWalkable(dest) || !svcCtx.ClaimPosition(dest, c.Handle) {
		color.HiRed("[HOP] [%s] destino inválido ou ocupado", c.PrimaryType)
		return core.StatusFailure
	}

	c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
		consts.MovementPlanHop,
		target.GetHandle(),
		1.0,
		500*time.Millisecond,
		target.GetPosition(),
	)

	color.White("[HOP] [%s] saltando em direção ao alvo", c.PrimaryType)
	return core.StatusRunning
}

func (n *PredatorHopApproachNode) Reset(c *creature.Creature) {
	color.Yellow("[HOP] [RESET] [%s] limpando plano", c.GetPrimaryType())
	c.MoveCtrl.MovementPlan = nil
	// FSM limpará o resto
}
