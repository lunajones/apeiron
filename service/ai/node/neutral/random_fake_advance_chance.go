package neutral

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomFakeAdvanceChanceNode struct {
	Chance float64 // ex: 0.2 para 20%
}

func (n *RandomFakeAdvanceChanceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[FAKE ADVANCE CHANCE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	if rand.Float64() < n.Chance {
		now := time.Now()
		log.Printf("[FAKE ADVANCE CHANCE] [%s] decidiu blefar", c.Handle.String())

		// Para feedback do próprio comportamento
		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "FakeAdvancePerformed",
			Timestamp:    now,
		})

		// Para oponentes detectarem
		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "FakeAdvanceBroadcast",
			Timestamp:    now,
		})

		return core.StatusSuccess
	}

	log.Printf("[FAKE ADVANCE CHANCE] [%s] decidiu não blefar", c.Handle.String())
	return core.StatusFailure
}

func (n *RandomFakeAdvanceChanceNode) Reset() {}
