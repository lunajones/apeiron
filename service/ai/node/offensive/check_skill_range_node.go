package offensive

import (
	"log"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CheckSkillRangeNode struct{}

func (n *CheckSkillRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.NextSkillToUse == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] Sem skill planejada, permite próximo node",
			c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	if c.MoveCtrl.IsMoving {
		return core.StatusRunning
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] nenhum alvo válido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		c.ClearMovementIntent()
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] dentro do range (%.2f <= %.2f), permitindo execução de skill",
			c.Handle.String(), c.PrimaryType, dist, rangeNeeded)
		return core.StatusSuccess
	}

	// Está fora do alcance da skill
	log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] fora do range (%.2f > %.2f), mantendo movimento",
		c.Handle.String(), c.PrimaryType, dist, rangeNeeded)
	return core.StatusRunning
}

func (n *CheckSkillRangeNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
}
