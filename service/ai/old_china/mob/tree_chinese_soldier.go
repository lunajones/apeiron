package mob

import (
	
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

func BuildChineseSoldierBT(players []player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node.DetectOtherCreatureNode{},
			&node.AttackTargetNode{SkillName: "SoldierSkill1"},
			&node.UseGroundSkillNode{
				SkillName: "SoldierGroundSlam",
			},
			&node.AttackTargetNode{
				SkillName: "SoldierSkill1",
			},
			&node.AttackIfVulnerableNode{},
			&node.RandomIdleNode{},
		},
	}
}
