package neutral

import (
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

type ApproachUntilInRangeNode struct{}

func (n *ApproachUntilInRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {

	plan := c.MoveCtrl.MovementPlan
	if plan != nil {
		if plan.Type == constslib.MovementPlanApproach {
			if time.Now().Before(plan.ExpiresAt) {
				color.Yellow("[APPROACH-IN-RANGE] [%s] plano de movimento APPROACH-IN-RANGE ativo, ignorando novo plano", c.Handle.String())
				return core.StatusRunning
			} else {
				color.HiRed("[APPROACH-IN-RANGE] [%s] plano APPROACH-IN-RANGE expirou, limpando movimento", c.Handle.String())
				c.ClearMovementIntent()
				c.SetAnimationState(constslib.AnimationIdle)
				c.MoveCtrl.MovementPlan = nil
			}
		}
	}

	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		color.HiRed("[APPROACH-IN-RANGE] [%s] não está em estado de movimento ou está esquivando", c.Handle.String())
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if drive.Caution < 0.4 || drive.Caution >= 0.7 {
		color.Yellow("[APPROACH-IN-RANGE] [%s] Caution %.2f fora do intervalo (0.4 ~ 0.7)", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil {
		color.HiMagenta("[APPROACH-IN-RANGE] [%s (%s)] Sem skill planejada", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		color.HiRed("[APPROACH-IN-RANGE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.IsMovementLocked() {
		color.Blue("[APPROACH-IN-RANGE] [%s] movimento travado até %v", c.Handle.String(), c.GetMovementLockUntil())
		return core.StatusRunning
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		color.HiRed("[APPROACH-IN-RANGE] [%s (%s)] alvo não encontrado", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	c.RecentActions = append(c.RecentActions, constslib.CombatActionApproach)

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		c.ClearMovementIntent()
		c.SetAnimationState(constslib.AnimationIdle)
		c.MoveCtrl.MovementPlan = nil
		color.Green("[APPROACH-IN-RANGE] [%s] no range para %s (dist=%.2f)", c.Handle.String(), c.NextSkillToUse.Name, dist)
		return core.StatusSuccess
	}

	// 🔸 Caminho com leve desvio lateral
	offsetDir := rand.Float64()*2 - 1
	forwardVec := position.NewVector2DFromTo(c.GetPosition(), target.GetPosition()).Normalize()
	lateralVec := forwardVec.Perpendicular().Scale(offsetDir * 1.0)
	destination := target.GetPosition().AddVector2D(lateralVec)

	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 {
		walkSpeed := c.WalkSpeed * 0.5
		c.MoveCtrl.SetMoveIntent(destination, walkSpeed, rangeNeeded)
		c.SetAnimationState(constslib.AnimationWalk)
		c.SetMovementLock(3 * time.Second)

		// ✅ Alimenta MovementPlan
		c.MoveCtrl.MovementPlan = &movement.MovementPlan{
			Type:            constslib.MovementPlanApproach,
			TargetHandle:    target.GetHandle(),
			DesiredDistance: rangeNeeded,
			ExpiresAt:       time.Now().Add(5 * time.Second),
		}

		color.HiCyan("[APPROACH-IN-RANGE] [%s] avançando com cautela (dist=%.2f offset=%.2f)", c.Handle.String(), dist, offsetDir)
	}

	return core.StatusRunning
}

func (n *ApproachUntilInRangeNode) Reset(c *creature.Creature) {
	color.Yellow("[APPROACH-IN-RANGE] [RESET] [%s (%s)] limpando movimento", c.Handle.String(), c.PrimaryType)
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
	c.MoveCtrl.MovementPlan = nil
}
