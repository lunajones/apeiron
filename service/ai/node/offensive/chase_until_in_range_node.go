package offensive

import (
	"log"
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

type ChaseUntilInRangeNode struct{}

func (n *ChaseUntilInRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	// ⚠️ VALIDAÇÃO DO PLANO DE MOVIMENTO ATUAL
	plan := c.MoveCtrl.MovementPlan
	if plan != nil {
		if plan.Type == constslib.MovementPlanChase {
			if time.Now().Before(plan.ExpiresAt) {
				color.Yellow("[CHASE-IN-RANGE] [%s] plano de movimento CHASE ativo, ignorando novo plano", c.Handle.String())
				return core.StatusRunning
			} else {
				color.HiRed("[CHASE-IN-RANGE] [%s] plano CHASE expirou, limpando movimento", c.Handle.String())
				c.ClearMovementIntent()
				c.SetAnimationState(constslib.AnimationIdle)
				c.MoveCtrl.MovementPlan = nil
			}
		}
	}

	drive := c.GetCombatDrive()
	// color.Cyan("[CHASE-IN-RANGE] [%s] PERSEGUINDO (%.2f)", c.GetPrimaryType(), drive.Caution)

	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		// color.HiRed("[CHASE-IN-RANGE] [%s] FALHOU EM PERSEGUIR (%.2f)", c.GetPrimaryType(), drive.Caution)
		return core.StatusFailure
	}

	if drive.Caution >= 0.4 {
		color.Yellow("[CHASE-IN-RANGE] [%s] Caution insuficiente (%.2f), ignorando perseguição", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil {
		color.HiMagenta("[CHASE-IN-RANGE] [%s (%s)] Sem skill planejada", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		color.HiRed("[CHASE-IN-RANGE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.IsMovementLocked() {
		color.Blue("[CHASE-IN-RANGE] [%s] movimento travado até %v", c.Handle.String(), c.GetMovementLockUntil())
		return core.StatusRunning
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		color.HiRed("[CHASE-IN-RANGE] [%s (%s)] alvo não encontrado", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range
	log.Println(color.WhiteString("[CHASE-IN-RANGE] [%s] distância até alvo: %.2f (range da skill: %.2f)", c.Handle.String(), dist, rangeNeeded))

	c.RecentActions = append(c.RecentActions, constslib.CombatActionChase)

	if dist <= rangeNeeded {
		color.Green("[CHASE-IN-RANGE] [%s] no range para %s (dist=%.2f)", c.Handle.String(), c.NextSkillToUse.Name, dist)
		c.ClearMovementIntent()
		c.SetAnimationState(constslib.AnimationIdle)
		c.MoveCtrl.MovementPlan = nil
		return core.StatusSuccess
	}
	log.Println(color.WhiteString("[CHASE-IN-RANGE] [%s] distância até alvo: %.2f, is moving %s", c.GetPrimaryType(), dist, rangeNeeded, c.MoveCtrl.IsMoving))

	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 || (c.MoveCtrl.MovementPlan != nil && time.Now().After(c.MoveCtrl.MovementPlan.ExpiresAt)) {

		color.HiCyan("[CHASE-IN-RANGE] [%s] iniciando perseguição (dist=%.2f, range=%.2f)", c.Handle.String(), dist, rangeNeeded)
		c.MoveCtrl.SetMoveIntent(target.GetPosition(), c.RunSpeed, rangeNeeded)
		c.SetMovementLock(5 * time.Second)

		c.MoveCtrl.MovementPlan = &movement.MovementPlan{
			Type:            constslib.MovementPlanChase,
			TargetHandle:    target.GetHandle(), // <- usa handle do alvo
			DesiredDistance: rangeNeeded,
			ExpiresAt:       time.Now().Add(3 * time.Second), // <- calcula vencimento
		}

	}

	color.White("[CHASE-IN-RANGE] [%s] posição atual: %s | destino: %s", c.Handle.String(), c.GetPosition(), c.MoveCtrl.TargetPosition)
	c.SetAnimationState(constslib.AnimationRun)
	return core.StatusRunning
}

func (n *ChaseUntilInRangeNode) Reset(c *creature.Creature) {
	color.Yellow("[CHASE-IN-RANGE] [%s] Reset chamado, limpando movimento", c.Handle.String())
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
}
