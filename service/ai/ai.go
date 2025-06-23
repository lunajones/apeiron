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
		if tree != nil {
			primaryType := string(c.PrimaryType)
			behaviorTrees[primaryType] = tree
			log.Printf("[AI] BehaviorTree carregada para %s", primaryType)
		} else {
			log.Printf("[AI] Nenhuma BehaviorTree encontrada para %s", c.ID)
		}
	}
}

func ProcessAI(c *creature.Creature, creatures []*creature.Creature, players []*player.Player) {
	tree, exists := behaviorTrees[string(c.PrimaryType)]
	if !exists {
		log.Printf("[AI] Nenhuma árvore de comportamento encontrada para %s (%s)", c.ID, c.PrimaryType)
		return
	}

	ctx := core.AIContext{
		Creatures: creatures,
		Players:   players,
	}

	tree.Tick(c, ctx)
}
