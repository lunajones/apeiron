package node

import (
	"log"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type DetectOtherCreatureNode struct{}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx := ctx.(dynamic_context.AIServiceContext)
	log.Printf("[AI] %s executando DetectOtherCreatureNode", c.ID)

	if c.IsBlind {
		log.Printf("[AI] %s está cego e não pode detectar criaturas", c.ID)
		return core.StatusFailure
	}

	for _, other := range svcCtx.GetServiceCreatures() {
		if other.ID == c.ID || !other.IsAlive {
			continue
		}

		dist := position.CalculateDistance(c.Position, other.Position)
		if dist > c.DetectionRadius {
			continue
		}

		// 1. Fome extrema prioriza caça
		hunger := c.GetNeedValue(creature.NeedHunger)
		if hunger > 80 && c.HasTag(creature.TagPredator) && other.HasTag(creature.TagPrey) {
			log.Printf("[AI] %s está faminto e detectou %s como presa.", c.ID, other.ID)
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		// 2. Predadores ignoram quem não for presa
		if c.HasTag(creature.TagPredator) && !other.HasTag(creature.TagPrey) {
			continue
		}

		// 3. Checa inimizade
		if !creature.AreEnemies(c, other) {
			continue
		}

		log.Printf("[AI] %s detectou %s como inimigo a %.2f de distância", c.ID, other.ID, dist)
		c.TargetCreatureID = other.ID
		c.ChangeAIState(creature.AIStateCombat)
		return core.StatusSuccess
	}

	log.Printf("[AI] %s não detectou nenhuma criatura relevante", c.ID)
	return core.StatusFailure
}
