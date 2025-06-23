package ai

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/factory"
	"github.com/lunajones/apeiron/service/player"
)

var behaviorTrees map[string]creature.BehaviorTree

func InitBehaviorTrees(players []*player.Player, creatures []*creature.Creature) {
	behaviorTrees = make(map[string]creature.BehaviorTree)

	for _, c := range creatures {
		tree := factory.CreateBehaviorTree(c.Types, players, creatures)
		if tree == nil {
			log.Printf("[AI] Nenhuma BehaviorTree atribuída para criatura %s (tipos: %v)", c.ID, c.Types)
			continue
		}
		behaviorTrees[c.ID] = tree
	}
}

func ProcessAI(c *creature.Creature, creatures []*creature.Creature, players []*player.Player) {
	tree, exists := behaviorTrees[c.ID]
	if !exists || tree == nil {
		return
	}

	ctx := core.AIContext{
		Creatures: creatures,
		Players:   players,
	}

	tree.Tick(c, ctx)
}
