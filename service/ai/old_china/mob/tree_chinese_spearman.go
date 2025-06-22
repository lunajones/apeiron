package mob

import (
	"math"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/creature"
)

type MaintainMediumDistanceNode struct {
	Players []ai.Player
}

func (n *MaintainMediumDistanceNode) Tick(c *creature.Creature) ai.BehaviorStatus {
	for _, p := range n.Players {
		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		idealMin := 4.0  // Exemplo: distância mínima segura
		idealMax := 8.0  // Exemplo: distância máxima para atacar

		if distance < idealMin {
			c.SetAction(creature.ActionRun) // Recuar
			return ai.StatusSuccess
		}

		if distance > idealMax {
			c.SetAction(creature.ActionRun) // Avançar
			return ai.StatusSuccess
		}

		// Se estiver em distância média, pode atacar
		c.SetAction(creature.ActionSkill2) // Exemplo: ataque de investida
		c.ChangeAIState(creature.AIStateAttack)
		return ai.StatusSuccess
	}

	return ai.StatusFailure
}

func BuildChineseSpearmanBT(players []ai.Player, creatures []*creature.Creature) ai.BehaviorNode {
	return &ai.SelectorNode{
		Children: []ai.BehaviorNode{
			&ai.FleeIfLowHPNode{},
			&ai.DetectPlayerNode{Players: players},
			&MaintainMediumDistanceNode{Players: players},
			&ai.DetectOtherCreatureNode{Creatures: creatures},
			&ai.RandomIdleBehaviorNode{},
		},
	}
}
