package defensive

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

type CircleAroundTargetNode struct {
	actionRegistered bool
}

func (n *CircleAroundTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	plan := c.MoveCtrl.MovementPlan
	if plan != nil {
		if plan.Type == constslib.MovementPlanCircle {
			if time.Now().Before(plan.ExpiresAt) {
				color.Yellow("[CIRCLE] [%s] plano de movimento CIRCLE ativo, ignorando novo plano", c.Handle.String())
				return core.StatusRunning
			} else {
				color.HiRed("[CIRCLE] [%s] plano CIRCLE expirou, limpando movimento", c.Handle.String())
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
	if drive.Caution < 0.7 || drive.Caution >= 0.9 {
		color.Yellow("[CIRCLE] [%s] Caution %.2f fora do intervalo (0.7 ~ 0.9)", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil {
		color.HiMagenta("[CIRCLE] [%s (%s)] Sem skill planejada", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		color.HiRed("[CIRCLE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.IsMovementLocked() {
		color.Blue("[CIRCLE] [%s] movimento travado até %v", c.Handle.String(), c.GetMovementLockUntil())
		return core.StatusRunning
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		color.HiRed("[CIRCLE] [%s (%s)] alvo não encontrado", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if !n.actionRegistered {
		c.RecentActions = append(c.RecentActions, constslib.CombatActionCircleAround)
		n.actionRegistered = true
	}

	idealRange := c.NextSkillToUse.Range
	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	inRange := dist >= idealRange*0.9 && dist <= idealRange*1.1

	if !finder.IsInFieldOfView(target, c, 100.0) && inRange {
		color.Green("[CIRCLE] [%s] pronto para atacar — fora do cone e no range", c.Handle.String())
		return core.StatusSuccess
	}

	if dist > idealRange*0.9 {
		moveTo := target.GetPosition()
		c.MoveCtrl.SetMoveIntent(moveTo, c.WalkSpeed, 0.0)
		c.SetAnimationState(constslib.AnimationWalk)
		c.SetMovementLock(1 * time.Second)

		// ✅ Alimenta MovementPlan
		c.MoveCtrl.MovementPlan = &movement.MovementPlan{
			Type:            constslib.MovementPlanCircle,
			TargetHandle:    target.GetHandle(),
			DesiredDistance: idealRange,
			ExpiresAt:       time.Now().Add(3 * time.Second),
		}

		color.HiCyan("[CIRCLE] [%s] aproximando para distância ideal (dist=%.2f)", c.Handle.String(), dist)
		return core.StatusRunning
	}

	// Reposiciona lateralmente ao redor do alvo
	dirVec := position.NewVector2DFromTo(target.GetPosition(), c.GetPosition()).Normalize()
	angle := rand.Float64()*1.5708 - 0.7854 // ±45°
	perp := position.RotateVector2D(dirVec, angle)
	moveTo := target.GetPosition().AddVector2D(perp.Scale(dist))

	if !svcCtx.NavMesh.IsWalkable(moveTo) {
		color.HiRed("[CIRCLE] [%s] posição de movimento não caminhável", c.Handle.String())
		return core.StatusFailure
	}

	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 {
		c.MoveCtrl.SetMoveIntent(moveTo, c.WalkSpeed, 0.0)
		c.SetAnimationState(constslib.AnimationWalk)
		c.SetMovementLock(1 * time.Second)

		// ✅ Alimenta MovementPlan
		c.MoveCtrl.MovementPlan = &movement.MovementPlan{
			Type:            constslib.MovementPlanCircle,
			TargetHandle:    target.GetHandle(),
			DesiredDistance: dist,
			ExpiresAt:       time.Now().Add(3 * time.Second),
		}

		color.White("[CIRCLE] [%s] circulando o alvo (dist=%.2f)", c.Handle.String(), dist)
	}

	return core.StatusRunning
}

func (n *CircleAroundTargetNode) Reset(c *creature.Creature) {
	color.Yellow("[CIRCLE] [RESET] [%s (%s)] limpando movimento", c.Handle.String(), c.PrimaryType)
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
	c.MoveCtrl.MovementPlan = nil
	n.actionRegistered = false
}
