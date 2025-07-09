package defensive

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type IsTargetDefendingNode struct{}

func (n *IsTargetDefendingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[DEF-DEFEND-CHECK] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[DEF-DEFEND-CHECK] [%s] sem alvo para checar defesa", c.Handle.String())
		return core.StatusFailure
	}

	if target.IsBlocking() {
		log.Printf("[DEF-DEFEND-CHECK] [%s] alvo [%s] está defendendo", c.Handle.String(), target.GetHandle().String())

		// Registra o evento
		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "TargetDefendingDetected",
			Timestamp:    time.Now(),
		})

		return core.StatusSuccess
	}

	log.Printf("[DEF-DEFEND-CHECK] [%s] alvo [%s] não está defendendo", c.Handle.String(), target.GetHandle().String())
	return core.StatusFailure
}

func (n *IsTargetDefendingNode) Reset() {
	// Nada a resetar
}
