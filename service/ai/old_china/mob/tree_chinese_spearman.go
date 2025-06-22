package mob

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type MaintainMediumDistanceNode struct {
	Players []player.Player
}


func BuildChineseSpearmanBT(players []player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SequenceNode{
		Children: []core.BehaviorNode{
			&node.FleeIfLowHPNode{},
			&node.DetectPlayerNode{Players: players},
			&MaintainMediumDistanceNode{Players: players},
			&node.DetectOtherCreatureNode{Creatures: creatures},
			&node.RandomIdleBehaviorNode{},
		},
	}
}

