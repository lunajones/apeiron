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
			color.Yellow("[CIRCLE] [%s] impedido por ciclo ativo de skill", c.PrimaryType)
			return core.StatusFailure
		}
	}

	plan := c.MoveCtrl.MovementPlan
	if plan != nil && plan.Type == constslib.MovementPlanCircle {
		if time.Now().Before(plan.ExpiresAt) {
			color.Yellow("[CIRCLE] [%s] plano CIRCLE ativo", c.PrimaryType)
			return core.StatusRunning
		}
		color.HiRed("[CIRCLE] [%s] plano CIRCLE expirou", c.PrimaryType)
		c.GetContext().ClearClaims(c.Handle)
		c.MoveCtrl.MovementPlan = nil
	}

	if c.IsDodging() {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if time.Since(c.GetLastCircleAt()) < 3*time.Second {
		color.Yellow("[CIRCLE] [%s] cooldown ativo", c.PrimaryType)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil && rand.Float64() < 0.4+drive.Counter*0.4 {
		color.White("[CIRCLE] [%s] sem skill ativa, circulando por decisão randômica", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	if c.HasRecentlyMissedSkill() || target.IsBlocking() {
		color.White("[CIRCLE] [%s] falha recente ou bloqueio detectado, reposicionando", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	if rand.Float64() < 0.2+drive.Counter*0.3 {
		color.White("[CIRCLE] [%s] decisão randômica pelo drive", c.PrimaryType)
		return n.executeCircle(c, target)
	}

	return core.StatusFailure
}

func (n *CircleAroundTargetNode) executeCircle(c *creature.Creature, target model.Targetable) interface{} {
	if c.MoveCtrl == nil {
		log.Printf("[CIRCLE] [%s] MoveCtrl nil", c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx := c.GetContext()
	idealRange := 1.5
	if c.NextSkillToUse != nil {
		idealRange = c.NextSkillToUse.Range
	}

	c.AddRecentAction(constslib.CombatActionCircleAround)

	if !finder.IsInFieldOfView(target, c, 100.0) {
		color.Green("[CIRCLE] [%s] pronto para atacar — fora do cone", c.PrimaryType)
		return core.StatusSuccess
	}

	dirVec := position.NewVector2DFromTo(c.GetPosition(), target.GetPosition()).Normalize()
	angle := rand.Float64()*1.5708 - 0.7854 // ±45°
	perp := position.RotateVector2D(dirVec, angle).Normalize()
	moveOffset := perp.Scale(3.0)
	moveTo := c.GetPosition().AddVector2D(moveOffset)

	finalDist := position.CalculateDistance(c.GetPosition(), moveTo)
	if finalDist < 2.9 {
		color.HiRed("[CIRCLE] [%s] destino muito próximo (%.2f)", c.PrimaryType, finalDist)
		return core.StatusFailure
	}

	if !svcCtx.NavMesh.IsWalkable(moveTo) || !svcCtx.ClaimPosition(moveTo, c.Handle) {
		color.HiRed("[CIRCLE] [%s] destino inválido", c.PrimaryType)
		return core.StatusFailure
	}

	c.MoveCtrl.SetMoveTarget(moveTo, c.WalkSpeed*0.7, 0.1)
	c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
		constslib.MovementPlanCircle,
		target.GetHandle(),
		idealRange,
		2*time.Second,
		target.GetPosition(),
	)
	c.SetLastCircleAt(time.Now())

	color.White("[CIRCLE] [%s] circulando o alvo (%.2f)", c.PrimaryType, finalDist)
	return core.StatusRunning
}

func (n *CircleAroundTargetNode) Reset(c *creature.Creature) {
	color.Yellow("[CIRCLE] [RESET] [%s]", c.PrimaryType)
	c.GetContext().ClearClaims(c.Handle)
	c.MoveCtrl.MovementPlan = nil
}
