package node

import (
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type MaintainMediumDistanceNode struct {
	Players []player.Player
}

func (n *MaintainMediumDistanceNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	for _, p := range ctx.Players {

		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		idealMin := 4.0
		idealMax := 8.0

		if distance < idealMin {
			c.SetAction(creature.ActionRun)
			return core.StatusSuccess
		}

		if distance > idealMax {
			c.SetAction(creature.ActionRun)
			return core.StatusSuccess
		}

		c.SetAction(creature.ActionSkill2)
		c.ChangeAIState(creature.AIStateAttack)
		return core.StatusSuccess
	}

	return core.StatusFailure
}
