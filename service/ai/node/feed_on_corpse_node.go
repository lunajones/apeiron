package node

import (
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FeedOnCorpseNode struct{}

func (n *FeedOnCorpseNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	if !c.HasTag(creature.TagPredator) {
		return core.StatusFailure
	}

	for _, other := range ctx.Creatures {
		if other.ID == c.ID || !other.IsCorpse {
			continue
		}

		distance := CalculateDistance(c.Position, other.Position)
		if distance > 2.0 {
			continue
		}

		log.Printf("[AI] %s está se alimentando do corpo de %s", c.ID, other.ID)
		c.ModifyNeed(creature.NeedHunger, -50) // Reduz a fome em 50 unidades (ajuste como quiser)
		c.SetAction(creature.ActionIdle)        // Após comer, volta pro idle

		// Opcional: marcar o corpo como já consumido
		other.IsCorpse = false

		return core.StatusSuccess
	}

	return core.StatusFailure
}