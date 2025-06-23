package node

import (
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type DetectOtherCreatureNode struct{}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando DetectOtherCreatureNode", c.ID)
	for _, other := range ctx.Creatures {
		if other.ID == c.ID || !other.IsAlive {
			continue
		}

		if creature.CanSeeOtherCreatures(c, []*creature.Creature{other}) || creature.CanHearOtherCreatures(c, []*creature.Creature{other}) {
			c.TargetCreatureID = other.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}
	return core.StatusFailure
}