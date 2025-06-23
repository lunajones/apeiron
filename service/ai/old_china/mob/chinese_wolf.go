package mob

import (
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type MaintainDistanceForWolfNode struct {
	Players []*player.Player
}

func (n *MaintainDistanceForWolfNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	for _, p := range ctx.Players {
		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		idealMin := 6.0
		idealMax := 12.0

		if distance < idealMin {
			c.SetAction(creature.ActionRun) // Recuar
			return core.StatusSuccess
		}

		if distance > idealMax {
			c.SetAction(creature.ActionRun) // Aproximar
			return core.StatusSuccess
		}

		c.SetAction(creature.ActionSkill2) // Exemplo: mordida
		c.ChangeAIState(creature.AIStateAttack)
		return core.StatusSuccess
	}

	return core.StatusFailure
}

func BuildChineseWolfBT(players []*player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node.DetectPlayerNode{Players: players},
			&MaintainDistanceForWolfNode{Players: players},
			&node.DetectOtherCreatureNode{},
			&node.AttackIfVulnerableNode{SkillName: "WolfBite"},
			&node.RandomIdleNode{},
		},
	}
}
