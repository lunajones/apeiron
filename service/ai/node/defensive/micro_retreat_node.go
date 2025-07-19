package defensive

import (
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type MicroRetreatNode struct {
	actionRegistered bool
}

func (n *MicroRetreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[MICRO-RETREAT] [%s] contexto inválido.", c.Handle.String())
		return core.StatusFailure
	}

	plan := c.MoveCtrl.MovementPlan
	if plan != nil && plan.Type == constslib.MovementPlanMicroRetreat {
		if time.Now().Before(plan.ExpiresAt) {
			color.Yellow("[MICRO-RETREAT] [%s] plano de movimento MICRO-RETREAT ativo, ignorando novo plano", c.Handle.String())
			return core.StatusRunning
		}
		color.HiRed("[MICRO-RETREAT] [%s] plano MICRO-RETREAT expirou, limpando movimento", c.Handle.String())
		c.ClearMovementIntent()
		c.SetAnimationState(constslib.AnimationIdle)
		c.MoveCtrl.MovementPlan = nil
		svcCtx.ClearClaims(c.Handle)
	}

	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if drive.Caution < 0.8 {
		log.Printf("[MICRO-RETREAT] [%s] Caution insuficiente (%.2f), recuo ignorado", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[MICRO-RETREAT] [%s] alvo não encontrado.", c.Handle.String())
		return core.StatusFailure
	}

	if !n.actionRegistered {
		c.RecentActions = append(c.RecentActions, constslib.CombatActionMicroRetreat)
		n.actionRegistered = true
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	if dist > 2.0 {
		return core.StatusFailure
	}

	// Decide direção de recuo
	mode := rand.Intn(2)
	forward := position.NewVector2DFromTo(target.GetPosition(), c.GetPosition()).Normalize()
	var retreatDir position.Vector2D

	if mode == 0 {
		retreatDir = forward.Scale(0.4)
	} else {
		lateral := forward.Perpendicular()
		side := 1.0
		if rand.Float64() < 0.5 {
			side = -1.0
		}
		retreatDir = forward.Scale(0.3).Add(lateral.Scale(0.3 * side))
	}

	dest := c.GetPosition().AddVector2D(retreatDir)
	if !svcCtx.NavMesh.IsWalkable(dest) || svcCtx.IsClaimedByOther(dest, c.Handle) {
		color.HiRed("[MICRO-RETREAT] [%s] posição bloqueada (walkable=%v, claimed)", c.Handle.String(), svcCtx.NavMesh.IsWalkable(dest))
		return core.StatusFailure
	}

	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 {
		svcCtx.ClearClaims(c.Handle)
		if !svcCtx.ClaimPosition(dest, c.Handle) {
			color.Red("[MICRO-RETREAT] [%s] falha ao claimar célula de destino", c.Handle.String())
			return core.StatusFailure
		}

		c.MoveCtrl.SetMoveTarget(dest, c.WalkSpeed*0.7, 0.1)
		c.SetAnimationState(constslib.AnimationWalk)

		c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
			constslib.MovementPlanMicroRetreat,
			target.GetHandle(),
			2.0,
			1*time.Second,
			target.GetPosition(),
		)

		log.Printf("[MICRO-RETREAT] [%s] recuando (modo=%d, dist=%.2f)", c.Handle.String(), mode, dist)
	}

	return core.StatusRunning
}

func (n *MicroRetreatNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
	c.MoveCtrl.MovementPlan = nil

	svcCtx := c.GetContext()
	if svcCtx != nil {
		svcCtx.ClearClaims(c.Handle)
	}

	log.Printf("[MICRO-RETREAT] [RESET] [%s (%s)]", c.Handle.String(), c.PrimaryType)
}
