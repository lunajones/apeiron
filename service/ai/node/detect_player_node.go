package node

import (
	"math"
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type DetectPlayerNode struct {
	Players []player.Player
}

func (n *DetectPlayerNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando DetectPlayerNode", c.ID)

	for _, p := range ctx.Players {
		dx := p.Position.X - c.Position.X
		dy := p.Position.Y - c.Position.Y
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dy*dy + dz*dz)

		visionRange := 10.0 // Exemplo de alcance de detecção
		if distance <= visionRange {
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}
	return core.StatusFailure
}
