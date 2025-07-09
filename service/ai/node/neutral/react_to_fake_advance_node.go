package neutral

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type ReactToFakeAdvanceNode struct{}

func (n *ReactToFakeAdvanceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[REACT FAKE ADVANCE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	events := svcCtx.GetRecentCombatBehaviors(c.TargetCreatureHandle, time.Now().Add(-2*time.Second))
	for _, e := range events {
		if e.BehaviorType == "FakeAdvanceBroadcast" {
			log.Printf("[REACT FAKE ADVANCE] [%s] detectou FakeAdvance do alvo %s", c.Handle.String(), e.SourceHandle.String())

			c.SetBlocking(true)
			c.CombatState = consts.CombatStateBlocking

			// Registra a reação
			svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
				SourceHandle: c.Handle,
				BehaviorType: "ReactedToFakeAdvance",
				Timestamp:    time.Now(),
			})

			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}

func (n *ReactToFakeAdvanceNode) Reset() {}
