package neutral

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type ExitCombatIfNoValidTargetsNode struct{}

func (n *ExitCombatIfNoValidTargetsNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[EXIT-COMBAT] [%s] contexto inválido", c.GetPrimaryType())
		return core.StatusFailure
	}

	targets := svcCtx.GetCachedTargets(c.Handle)
	best := c.GetBestTargetFromTargets(targets)
	if best == nil {
		if c.AIState == constslib.AIStateCombat {
			log.Printf("[EXIT-COMBAT] [%s] nenhum alvo restante — voltando para Idle", c.GetPrimaryType())
			c.ChangeAIState(constslib.AIStateIdle)
		}
		return core.StatusSuccess
	}

	// Ainda tem alvo, segue no combate
	return core.StatusSuccess
}

func (n *ExitCombatIfNoValidTargetsNode) Reset(c *creature.Creature) {}
