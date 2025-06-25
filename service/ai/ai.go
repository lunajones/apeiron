package ai

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

// InitBehaviorTrees agora não faz nada, pois cada criatura já tem sua árvore setada
func InitBehaviorTrees(players []*player.Player, creatures []*creature.Creature) {
	log.Println("[AI] InitBehaviorTrees chamado, mas não faz nada.")
}

// ProcessAI usa dynamic_context.AIServiceContext com as criaturas vivas
func ProcessAI(c *creature.Creature, creatures []*creature.Creature, players []*player.Player) {
	if c.BehaviorTree == nil {
		log.Printf("[AI] Nenhuma árvore de comportamento encontrada para %s (%s)", c.ID, c.PrimaryType)
		return
	}

	ctx := dynamic_context.AIServiceContext{
		Creatures: creatures,
		Players:   players,
	}

	c.BehaviorTree.Tick(c, ctx)
}
