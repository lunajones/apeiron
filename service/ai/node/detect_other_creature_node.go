package node

import (
	"github.com/lunajones/apeiron/service/creature"
)

type DetectOtherCreatureNode struct {
	Creatures []*creature.Creature
}

func (n *DetectOtherCreatureNode) Tick(c *creature.Creature) BehaviorStatus {
	for _, target := range n.Creatures {
		if target.ID == c.ID || !target.IsAlive {
			continue
		}

		// Verifica visão
		if CanSeeOtherCreatures(c, []*creature.Creature{target}) != nil {
			EvaluateBehavior(c, target)
			return StatusSuccess
		}

		// Verifica audição
		if CanHearOtherCreatures(c, []*creature.Creature{target}) != nil {
			EvaluateBehavior(c, target)
			return StatusSuccess
		}
	}

	return StatusFailure
}
