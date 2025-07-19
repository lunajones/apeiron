package defensive

import (
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CircleAroundTargetNode struct{}

func (n *CircleAroundTargetNode) Tick(c *creature.Creature, _ interface{}) interface{} {
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		color.HiRed("[CIRCLE] [%s] alvo não encontrado", c.PrimaryType)
		return core.StatusFailure
	}

	if c.NextSkillToUse != nil {
		state := c.SkillStates[c.NextSkillToUse.Action]
		if state != nil && state.InUse && (state.WindUpFired || state.CastFired || state.RecoveryFired) {
			color.Yellow("[CIRCLE] [%s] impedido por ciclo ativo de skill: windup/cast/recovery", c.PrimaryType)
			return core.StatusFailure
		}
	}

	plan := c.MoveCtrl.MovementPlan
	if plan != nil && plan.Type == constslib.MovementPlanCircle {
		if time.Now().Before(plan.ExpiresAt) {
			color.Yellow("[CIRCLE] [%s] plano CIRCLE ativo, ignorando novo plano", c.PrimaryType)
			return core.StatusRunning
		}
		color.HiRed("[CIRCLE] [%s] plano CIRCLE expirou, limpando movimento", c.PrimaryType)
		c.GetContext().ClearClaims(c.Handle)
		c.ClearMovementIntent()
		c.SetAnimationState(constslib.AnimationIdle)
		c.MoveCtrl.MovementPlan = nil
	}

	if c.IsDodging() {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if time.Since(c.GetLastCircleAt()) < 3*time.Second {
		color.Yellow("[CIRCLE] [%s] cooldown interno ativo", c.PrimaryType)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil && rand.Float64() < 0.4+drive.Counter*0.4 {
		color.White("[CIRCLE] [%s] sem skill ativa, circulando por decisão randômica", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	if c.HasRecentlyMissedSkill() || target.IsBlocking() {
		color.White("[CIRCLE] [%s] falha recente ou alvo bloqueando, reposicionando", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	if rand.Float64() < 0.2+drive.Counter*0.3 {
		color.White("[CIRCLE] [%s] decisão randômica influenciada pelo drive", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	return core.StatusFailure
}

func (n *CircleAroundTargetNode) executeCircle(c *creature.Creature, target model.Targetable) interface{} {
	if c.MoveCtrl == nil {
		log.Printf("[CIRCLE] [%s] MoveCtrl está nil — abortando", c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx := c.GetContext()
	idealRange := c.NextSkillToUse.Range
	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	inRange := dist >= idealRange*0.9 && dist <= idealRange*1.1

	c.RecentActions = append(c.RecentActions, constslib.CombatActionCircleAround)

	if !finder.IsInFieldOfView(target, c, 100.0) && inRange {
		color.Green("[CIRCLE] [%s] pronto para atacar — fora do cone e no range", c.PrimaryType)
		return core.StatusSuccess
	}

	if dist > idealRange*0.9 {
		moveTo := target.GetPosition()

		if !svcCtx.NavMesh.IsWalkable(moveTo) || !svcCtx.ClaimPosition(moveTo, c.Handle) {
			color.HiRed("[CIRCLE] [%s] posição inválida ou já ocupada para aproximação", c.PrimaryType)
			return core.StatusFailure
		}

		c.ClearMovementIntent()
		c.MoveCtrl.IsMoving = false
		c.MoveCtrl.SetMoveTarget(moveTo, c.WalkSpeed, 0.0)
		c.SetAnimationState(constslib.AnimationWalk)
		c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
			constslib.MovementPlanCircle,
			target.GetHandle(),
			idealRange,
			3*time.Second,
			target.GetPosition(),
		)
		c.SetLastCircleAt(time.Now())
		color.HiCyan("[CIRCLE] [%s] aproximando para distância ideal (dist=%.2f)", c.PrimaryType, dist)
		return core.StatusRunning
	}

	minMoveDist := 3.0
	dirVec := position.NewVector2DFromTo(target.GetPosition(), c.GetPosition()).Normalize()
	angle := rand.Float64()*1.5708 - 0.7854
	perp := position.RotateVector2D(dirVec, angle)

	moveOffset := perp.Scale(dist)
	if moveOffset.Length() < minMoveDist {
		moveOffset = perp.Normalize().Scale(minMoveDist)
	}

	moveTo := target.GetPosition().AddVector2D(moveOffset)
	finalDist := position.CalculateDistance(c.GetPosition(), moveTo)

	if finalDist < 2.9 {
		color.HiRed("[CIRCLE] [%s] ponto de destino muito próximo (%.2f), abortando", c.PrimaryType, finalDist)
		return core.StatusFailure
	}

	if !svcCtx.NavMesh.IsWalkable(moveTo) || !svcCtx.ClaimPosition(moveTo, c.Handle) {
		color.HiRed("[CIRCLE] [%s] destino de movimento inválido ou já ocupado", c.PrimaryType)
		return core.StatusFailure
	}

	c.ClearMovementIntent()
	c.MoveCtrl.IsMoving = false
	c.MoveCtrl.SetMoveTarget(moveTo, c.WalkSpeed*0.7, 0.1)
	c.SetAnimationState(constslib.AnimationWalk)
	c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
		constslib.MovementPlanCircle,
		target.GetHandle(),
		idealRange,
		2*time.Second,
		target.GetPosition(),
	)
	c.SetLastCircleAt(time.Now())

	color.White("[CIRCLE] [%s] circulando o alvo (finalDist=%.2f)", c.PrimaryType, finalDist)
	return core.StatusRunning
}

func (n *CircleAroundTargetNode) Reset(c *creature.Creature) {
	color.Yellow("[CIRCLE] [RESET] [%s] limpando movimento", c.PrimaryType)
	c.GetContext().ClearClaims(c.Handle)
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
	c.MoveCtrl.MovementPlan = nil
}
