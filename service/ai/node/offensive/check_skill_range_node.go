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

	if dist <= c.NextSkillToUse.Range {
		// Já está no range → limpa movimento e permite seguir no sequence
		c.ClearMovementIntent()
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] no range para %s (dist=%.2f, range=%.2f)",
			c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, dist, c.NextSkillToUse.Range)
		return core.StatusSuccess
	}

	// Fora do range → indica que precisa seguir perseguindo
	log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] fora do range (dist=%.2f, range=%.2f), seguir chase",
		c.Handle.String(), c.PrimaryType, dist, c.NextSkillToUse.Range)
	return core.StatusRunning
}

func (n *CheckSkillRangeNode) Reset() {
	// Nada a resetar
}
