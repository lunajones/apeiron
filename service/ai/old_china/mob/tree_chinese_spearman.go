package mob

import (
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type MaintainMediumDistanceNode struct {
	Players []ai.Player
}


func BuildChineseSpearmanBT(players []node.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node1.DetectPlayerNode{Players: players},
			&MaintainMediumDistanceNode{Players: players},
			&node.DetectOtherCreatureNode{Creatures: creatures},
			&node.RandomIdleBehaviorNode{},
		},
	}
}

