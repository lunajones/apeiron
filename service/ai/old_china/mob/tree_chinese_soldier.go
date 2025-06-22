package mob

import (
	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/creature"
)

func BuildChineseSoldierBT(creatures []*creature.Creature) ai.BehaviorNode {
	return &ai.SelectorNode{
		Children: []ai.BehaviorNode{
			&ai.FleeIfLowHPNode{},
			&ai.DetectOtherCreatureNode{Creatures: creatures},
			&ai.AttackIfEnemyVulnerableNode{},
			&ai.RandomIdleBehaviorNode{},
		},
	}
}
