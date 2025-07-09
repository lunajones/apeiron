package neutral

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type CheckOrEnterExitCombatNode struct{}

func (n *CheckOrEnterExitCombatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[COMBAT] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := c.GetBestTargetFromTargets(svcCtx.GetCachedTargets(c.Handle))
	if target == nil {
		// Caso não tenha alvo, saímos do estado de combate
		log.Printf("[COMBAT] [%s] sem alvo válido — saindo do AIStateCombat", c.Handle.String())
		if c.AIState == constslib.AIStateCombat {
			c.ChangeAIState(constslib.AIStateIdle) // Saindo do combate
		}
		return core.StatusSuccess
	}

	// Caso tenha alvo válido, garantimos que a criatura está no estado de combate
	if c.AIState != constslib.AIStateCombat {
		log.Printf("[COMBAT] [%s] entrando no estado de combate", c.Handle.String())
		c.ChangeAIState(constslib.AIStateCombat) // Entrando no combate
	}
	return core.StatusSuccess
}

func (n *CheckOrEnterExitCombatNode) Reset() {
	// Nada necessário
}
