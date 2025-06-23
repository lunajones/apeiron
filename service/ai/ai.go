package ai

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/factory"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

var behaviorTrees map[creature.CreatureType]core.BehaviorNode

func InitBehaviorTrees(players []player.Player, creatures []*creature.Creature) {
	behaviorTrees = map[creature.CreatureType]core.BehaviorNode{
		creature.Soldier:        factory.BuildChineseSoldierBT(players, creatures),
	}
}

func ProcessAI(c *creature.Creature, creatures []*creature.Creature, players []*player.Player) {
	tree, exists := behaviorTrees[c.Type]
	if !exists {
		return
	}

	ctx := core.AIContext{
		Creatures: creatures,
		Players:   players,
	}

	tree.Tick(c, ctx)
}
