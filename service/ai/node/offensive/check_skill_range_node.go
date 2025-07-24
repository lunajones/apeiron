package offensive

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CheckSkillRangeNode struct{}

func (n *CheckSkillRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.NextSkillToUse == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s] Sem skill planejada, permite próximo node",
			c.PrimaryType)
		return core.StatusSuccess
	}

	if c.GetCombatState() == constslib.CombatStateCasting {
		return core.StatusRunning
	}

	if c.MoveCtrl.IsMoving || c.IsInTacticalMovement() {
		return core.StatusRunning
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s] nenhum alvo válido", c.PrimaryType)
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		c.ClearMovementIntent()
		c.SetCombatState(constslib.CombatStateCasting)
		log.Printf("[CHECK-SKILL-RANGE] [%s] dentro do range (%.2f <= %.2f), iniciando cast",
			c.PrimaryType, dist, rangeNeeded)
		return core.StatusSuccess
	}

	log.Printf("[CHECK-SKILL-RANGE] [%s] fora do range (%.2f > %.2f), mantendo movimento",
		c.PrimaryType, dist, rangeNeeded)

	c.CombatState = constslib.CombatStateMoving
	return core.StatusRunning
}

func (n *CheckSkillRangeNode) Reset(c *creature.Creature) {
	// c.ClearMovementIntent()
}
