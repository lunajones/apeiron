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

	plan := c.MoveCtrl.MovementPlan
	if plan != nil {
		if plan.Type == constslib.MovementPlanMicroRetreat {
			if time.Now().Before(plan.ExpiresAt) {
				color.Yellow("[MICRO-RETREAT] [%s] plano de movimento MICRO-RETREAT ativo, ignorando novo plano", c.Handle.String())
				return core.StatusRunning
			} else {
				color.HiRed("[MICRO-RETREAT] [%s] plano MICRO-RETREAT expirou, limpando movimento", c.Handle.String())
				c.ClearMovementIntent()
				c.SetAnimationState(constslib.AnimationIdle)
				c.MoveCtrl.MovementPlan = nil
			}
		}
	}
	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if drive.Caution < 0.9 {
		log.Printf("[MICRO-RETREAT] [%s] Caution insuficiente (%.2f), recuo ignorado", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[MICRO-RETREAT] [%s] contexto inválido.", c.Handle.String())
		return core.StatusFailure
	}

	if c.IsMovementLocked() {
		return core.StatusRunning
	}

	// Validação de plano já existente
	if plan := c.MoveCtrl.MovementPlan; plan != nil {
		if plan.Type == constslib.MovementPlanMicroRetreat && time.Now().Before(plan.ExpiresAt) {
			log.Printf("[MICRO-RETREAT] [%s] plano de recuo ativo, ignorando novo plano", c.Handle.String())
			return core.StatusRunning
		}
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

	// Decide o tipo de recuo
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

	destination := c.GetPosition().AddVector2D(retreatDir)

	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 {
		c.MoveCtrl.SetMoveIntent(destination, c.WalkSpeed*0.7, 0.1)
		c.SetAnimationState(constslib.AnimationWalk)
		c.SetMovementLock(1 * time.Second)

		c.MoveCtrl.MovementPlan = &movement.MovementPlan{
			Type:            constslib.MovementPlanMicroRetreat,
			TargetHandle:    target.GetHandle(),
			DesiredDistance: 2.0,
			ExpiresAt:       time.Now().Add(1 * time.Second),
		}
	}

	log.Printf("[MICRO-RETREAT] [%s] recuando (modo=%d, dist=%.2f)", c.Handle.String(), mode, dist)

	return core.StatusRunning
}

func (n *MicroRetreatNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
	c.MoveCtrl.MovementPlan = nil
	log.Printf("[MICRO-RETREAT] [RESET] [%s (%s)]", c.Handle.String(), c.PrimaryType)
}
