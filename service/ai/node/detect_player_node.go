package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type DetectPlayerNode struct {
	Players []player.Player
}

func (n *DetectPlayerNode) Tick(c *creature.Creature) core.BehaviorStatus {
	for _, p := range n.Players {
		// Verifica visão
		if creature.CanSeePlayer(c, []player.Player{p}) != nil {
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		// Verifica audição
		if creature.CanHearPlayer(c, []player.Player{p}) != nil {
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}
