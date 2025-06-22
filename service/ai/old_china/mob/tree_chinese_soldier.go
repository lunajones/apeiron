package mob

import (
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
)

func BuildChineseSoldierBT(players []node.Player, creatures []*creature.Creature) node.BehaviorNode {
	return &node.SelectorNode{
		Children: []node.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node.DetectOtherCreatureNode{},
			&node.AttackTargetNode{SkillName: "SoldierSkill1"},
			&node.UseGroundSkillNode{
				SkillName: "SoldierGroundSlam",
				Players:   players,
			},
			&node.AttackTargetNode{
				SkillName: "SoldierSkill1",
			},
			&node.AttackIfEnemyVulnerableNode{},
			&node.RandomIdleBehaviorNode{},
		},
	}
}
