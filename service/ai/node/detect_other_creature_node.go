package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type DetectOtherCreatureNode struct{}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature, ctx ai_context.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando DetectOtherCreatureNode", c.ID)
	for _, other := range ctx.Creatures {
		// Ignorar a si mesmo
		if other.ID == c.ID || !other.IsAlive {
			continue
		}

		// Se a criatura estiver com medo, ela evita olhar para outras criaturas
		if c.MentalState == creature.MentalStateAfraid {
			log.Printf("[AI] %s está com medo e ignorando %s.", c.ID, other.ID)
			continue
		}

		// Se for animal e o outro não for prey, ignora
		if c.HasTag(creature.TagAnimal) && !other.HasTag(creature.TagPrey) {
			continue
		}

		// Se for predador com fome extrema, detecta qualquer prey
		hunger := c.GetNeedValue(creature.NeedHunger)
		if hunger > 80 && c.HasTag(creature.TagPredator) && other.HasTag(creature.TagPrey) {
			log.Printf("[AI] %s está faminto e detectou %s como presa.", c.ID, other.ID)
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		// Comportamento normal: perceber se pode ver ou ouvir o outro
		if creature.CanSeeOtherCreatures(c, []*creature.Creature{other}) || creature.CanHearOtherCreatures(c, []*creature.Creature{other}) {
			log.Printf("[AI] %s detectou %s por visão ou audição.", c.ID, other.ID)
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}
