package mob

import (
	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/creature"
)

func BuildChineseSoldierBT(players []ai.Player, creatures []*creature.Creature) ai.BehaviorNode {
	return &ai.SelectorNode{
		Children: []ai.BehaviorNode{
			&ai.FleeIfLowHPNode{},
			&ai.DetectPlayerNode{Players: players},
			&ai.DetectOtherCreatureNode{Creatures: creatures},
			&ai.UseGroundSkillNode{
				SkillName: "SoldierGroundSlam",
				Players:   players,
			},
			&ai.AttackTargetNode{
				SkillName: "SoldierSkill1",
			},
			&ai.AttackIfEnemyVulnerableNode{},
			&ai.RandomIdleBehaviorNode{},
		},
	}
}
