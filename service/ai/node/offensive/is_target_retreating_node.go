package offensive

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

// IsTargetRetreatingNode verifica se o alvo está recuando (distanciando-se de maneira ativa)
type IsTargetRetreatingNode struct{}

func (n *IsTargetRetreatingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[RETREAT-CHECK] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[RETREAT-CHECK] [%s] sem alvo válido para verificar recuo", c.Handle.String())
		return core.StatusFailure
	}

	if c.LastKnownDistance == 0 {
		c.LastKnownDistance = position.CalculateDistance(c.Position, target.GetPosition())
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.Position, target.GetPosition())

	if dist > c.LastKnownDistance {
		c.LastKnownDistance = dist
		log.Printf("[RETREAT-CHECK] [%s] Alvo está recuando", c.Handle.String())

		// Registra evento
		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "TargetRetreatingDetected",
			Timestamp:    time.Now(),
		})

		return core.StatusSuccess
	}

	c.LastKnownDistance = dist
	return core.StatusFailure
}

func (n *IsTargetRetreatingNode) Reset() {
	// Nada a resetar
}
