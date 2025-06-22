package mob

import (
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
)

func BuildChineseSpearmanBT(players []node.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
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
