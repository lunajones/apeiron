package node

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib/position"
)

type FeedOnCorpseNode struct{}

func (n *FeedOnCorpseNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx := ctx.(dynamic_context.AIServiceContext)

	log.Printf("[AI] %s executando FeedOnCorpseNode", c.ID)

	// Só predadores devem se alimentar
	if !c.HasTag(creature.TagPredator) {
		log.Printf("[AI] %s não é um predador, não vai se alimentar de cadáveres.", c.ID)
		return core.StatusFailure
	}

	// Verifica corpos de criaturas próximas
	for _, other := range svcCtx.GetServiceCreatures() {
		if other.ID == c.ID || other.IsAlive || !other.IsCorpse {
			continue
		}

		if creature.AreEnemies(c, other) {
			dist := position.CalculateDistance(c.Position, other.Position)
			if dist <= 1.8 {
				log.Printf("[AI] %s encontrou cadáver inimigo próximo, vai se alimentar.", c.ID)
				c.SetAction(creature.ActionSkill1) // Supondo que seja anim de comer
				c.ChangeAIState(creature.AIStateFeeding)

				c.Needs = creature.ReduceNeed(c.Needs, creature.NeedHunger, 25)
				c.Memory = append(c.Memory, creature.MemoryEvent{
					Description: "Alimentou-se de cadáver",
					Timestamp:   time.Now(),
				})
				return core.StatusSuccess
			}
		}
	}

	log.Printf("[AI] %s não encontrou cadáveres adequados.", c.ID)
	return core.StatusFailure
}
