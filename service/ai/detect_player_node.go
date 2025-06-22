package ai

import (
	"github.com/lunajones/apeiron/service/creature"
)

type DetectPlayerNode struct {
	Players []Player
}

func (n *DetectPlayerNode) Tick(c *creature.Creature) BehaviorStatus {
	for _, p := range n.Players {
		// Verifica visão
		if CanSeePlayer(c, []Player{p}) != nil {
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return StatusSuccess
		}

		// Verifica audição
		if CanHearPlayer(c, []Player{p}) != nil {
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return StatusSuccess
		}
	}

	return StatusFailure
}
