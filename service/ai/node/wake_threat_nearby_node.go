package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type WakeIfThreatNearbyNode struct{}

func (n *WakeIfThreatNearbyNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI] [%s] contexto inválido em WakeIfThreatNearbyNode", c.Handle.String())
		return core.StatusFailure
	}

	if c.AIState != consts.AIStateSleeping {
		return core.StatusFailure
	}

	hearingRange := c.HearingRange * 0.5
	smellRange := c.SmellRange * 0.5

	for _, other := range svcCtx.GetServiceCreatures(c.Position, hearingRange+smellRange) {
		if other.Handle.ID == c.Handle.ID || !other.IsAlive {
			continue
		}
		if !creature.AreEnemies(c, other) {
			continue
		}

		distance := position.CalculateDistance(c.Position, other.Position)

		if distance <= hearingRange || distance <= smellRange {
			if other.IsCurrentlyCrouched() {
				stealthFailChance := 0.25
				if rand.Float64() >= stealthFailChance {
					log.Printf("[AI] [%s (%s)] ameaça agachada próxima (%.2f), mas não acordou (stealth bem-sucedido).",
						c.Handle.String(), c.PrimaryType, distance)
					continue
				}
				log.Printf("[AI] [%s (%s)] ameaça agachada próxima (%.2f), mas falhou no stealth. Acordando.",
					c.Handle.String(), c.PrimaryType, distance)
			} else {
				log.Printf("[AI] [%s (%s)] sentiu ameaça próxima (%.2f). Acordando.",
					c.Handle.String(), c.PrimaryType, distance)
			}

			c.LastThreatSeen = time.Now()
			c.ChangeAIState(consts.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}

func (n *WakeIfThreatNearbyNode) Reset() {
	// Não há estado interno para resetar neste node
}
