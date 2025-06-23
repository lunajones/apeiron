package mob

import (
	"math"
	
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type MaintainMediumDistanceNode struct {
	Players []player.Player
}


func BuildChineseSpearmanBT(players []*player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node.DetectPlayerNode{Players: players},
			&MaintainMediumDistanceNode{},
			&node.DetectOtherCreatureNode{},
			&node.RandomIdleNode{},
		},
	}
}


func (n *MaintainMediumDistanceNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	for _, p := range ctx.Players {
		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		idealMin := 4.0  // Exemplo: distância mínima segura
		idealMax := 8.0  // Exemplo: distância máxima para atacar

		if distance < idealMin {
			c.SetAction(creature.ActionRun) // Recuar
			return core.StatusSuccess
		}

		if distance > idealMax {
			c.SetAction(creature.ActionRun) // Avançar
			return core.StatusSuccess
		}

		// Se estiver em distância média, pode atacar
		c.SetAction(creature.ActionSkill2) // Exemplo: ataque de investida
		c.ChangeAIState(creature.AIStateAttack)
		return core.StatusSuccess
	}

	return core.StatusFailure
}

