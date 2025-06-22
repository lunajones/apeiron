package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

type DetectOtherCreatureNode struct{}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature) core.BehaviorStatus {
	for _, zone := range world.Zones {
		for _, other := range zone.Creatures {
			if other.ID == c.ID || !other.IsAlive {
				continue
			}

			if creature.CanSeeOtherCreatures(c, other) || creature.CanHearOtherCreatures(c, other) {
				c.TargetCreatureID = other.ID
				c.ChangeAIState(creature.AIStateAlert)
				return core.StatusSuccess
			}
		}
	}
	return core.StatusFailure
}
