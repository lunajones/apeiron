package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type DetectOtherCreatureNode struct{}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature, ctx dynamic_context.AIServiceContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando DetectOtherCreatureNode", c.ID)
	for _, other := range ctx.GetServiceCreatures() {
		if other.ID == c.ID || !other.IsAlive {
			continue
		}

		if c.MentalState == creature.MentalStateAfraid {
			log.Printf("[AI] %s está com medo e ignorando %s.", c.ID, other.ID)
			continue
		}

		if c.HasTag(creature.TagAnimal) && !other.HasTag(creature.TagPrey) {
			continue
		}

		hunger := c.GetNeedValue(creature.NeedHunger)
		if hunger > 80 && c.HasTag(creature.TagPredator) && other.HasTag(creature.TagPrey) {
			log.Printf("[AI] %s está faminto e detectou %s como presa.", c.ID, other.ID)
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		if creature.CanSeeOtherCreatures(c, []*creature.Creature{other}) || creature.CanHearOtherCreatures(c, []*creature.Creature{other}) {
			log.Printf("[AI] %s detectou %s por visão ou audição.", c.ID, other.ID)
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}
