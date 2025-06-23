package ai

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/factory"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

var behaviorTrees map[creature.CreatureType]creature.BehaviorTree

func InitBehaviorTrees(players []*player.Player, creatures []*creature.Creature) {
	behaviorTrees = make(map[creature.CreatureType]creature.BehaviorTree)

	for _, c := range creatures {
		tree := factory.CreateBehaviorTree(c.Type, players, creatures)
		c.BehaviorTree = tree
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

func GetBehaviorTreeForType(cType creature.CreatureType) (creature.BehaviorTree, bool) {
	tree, exists := behaviorTrees[cType]
	return tree, exists
}
