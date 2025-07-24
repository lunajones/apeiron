package offensive

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type ChaseUntilInRangeNode struct{}

func (n *ChaseUntilInRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	drive := c.GetCombatDrive()

	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		return core.StatusFailure
	}

	if drive.Caution >= 0.2 {
		color.Yellow("[CHASE-IN-RANGE] [%s] Caution insuficiente (%.2f), ignorando perseguição", c.GetPrimaryType(), drive.Caution)
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		color.HiRed("[CHASE-IN-RANGE] [%s] alvo não encontrado", c.GetPrimaryType())
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())

	var rangeNeeded float64
	if c.NextSkillToUse != nil {
		rangeNeeded = c.NextSkillToUse.Range
	} else {
		rangeNeeded = 1.5
		color.HiMagenta("[CHASE-IN-RANGE] [%s] Sem skill planejada, usando fallback range %.2f", c.GetPrimaryType(), rangeNeeded)
	}

	color.White("[CHASE-IN-RANGE] [%s] distância até alvo: %.2f (range considerado: %.2f)", c.GetPrimaryType(), dist, rangeNeeded)

	c.AddRecentAction(constslib.CombatActionChase)

	if dist <= rangeNeeded {
		color.Green("[CHASE-IN-RANGE] [%s] no range (dist=%.2f)", c.GetPrimaryType(), dist)
		return core.StatusSuccess
	}

	plan := c.MoveCtrl.MovementPlan
	if plan != nil && plan.Type == constslib.MovementPlanChase {
		const minRecalcDist = 1.0
		lastPos := plan.LastTargetPosition
		currentPos := target.GetPosition()
		moved := position.CalculateDistance(lastPos, currentPos)

		if moved > minRecalcDist {
			color.Red("[CHASE-IN-RANGE] [%s] alvo se moveu %.2f m, refazendo plano", c.GetPrimaryType(), moved)
			// FSM cuidará do intent, não precisamos limpar aqui
		} else if time.Now().Before(plan.ExpiresAt) {
			color.Yellow("[CHASE-IN-RANGE] [%s] plano CHASE ainda válido, mantendo (%.2fs restantes)", c.GetPrimaryType(), time.Until(plan.ExpiresAt).Seconds())
			return core.StatusRunning
		}
	}

	color.HiCyan("[CHASE-IN-RANGE] [%s] iniciando perseguição (dist=%.2f, range=%.2f)", c.GetPrimaryType(), dist, rangeNeeded)

	baseTargetPos := target.GetPosition()
	if !c.GetContext().ClaimPosition(baseTargetPos, c.Handle) {
		offset := rand.Float64()*2 - 1
		forward := position.NewVector2DFromTo(c.GetPosition(), baseTargetPos).Normalize()
		side := forward.Perpendicular().Scale(offset * 1.5)
		alt := baseTargetPos.AddVector2D(side)

		if c.GetContext().ClaimPosition(alt, c.Handle) {
			baseTargetPos = alt
			color.Red("[CHASE-IN-RANGE] [%s] célula ocupada, usando posição lateral %.2f", c.GetPrimaryType(), offset)
		} else {
			color.Red("[CHASE-IN-RANGE] [%s] nenhuma posição viável encontrada", c.GetPrimaryType())
			return core.StatusFailure
		}
	}

	c.MoveCtrl.SetMoveTarget(baseTargetPos, c.RunSpeed, rangeNeeded)

	c.MoveCtrl.MovementPlan = &movement.MovementPlan{
		Type:               constslib.MovementPlanChase,
		TargetHandle:       target.GetHandle(),
		DesiredDistance:    rangeNeeded,
		ExpiresAt:          time.Now().Add(3 * time.Second),
		LastTargetPosition: target.GetPosition(),
	}

	color.White("[CHASE-IN-RANGE] [%s] posição atual: %s | destino: %s", c.GetPrimaryType(), c.GetPosition(), c.MoveCtrl.TargetPosition)
	return core.StatusRunning
}

func (n *ChaseUntilInRangeNode) Reset(c *creature.Creature) {
	color.Yellow("[CHASE-IN-RANGE] [%s] Reset chamado", c.GetPrimaryType())
}
