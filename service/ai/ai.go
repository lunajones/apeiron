package ai

import (
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

var behaviorTrees map[creature.CreatureType]core.BehaviorNode

func Init() {
	InitBehaviorRules()
	dummyPlayers := []player.Player{}              // Por enquanto vazio
	dummyCreatures := []*creature.Creature{}       // Por enquanto vazio

	behaviorTrees = map[creature.CreatureType]core.BehaviorNode{
		creature.Soldier: mob.BuildChineseSoldierBT(dummyPlayers, dummyCreatures),
		// No futuro: adicionar outros tipos
	}
}

func ProcessAI(c *creature.Creature, creatures []*creature.Creature) {
	tree, exists := behaviorTrees[c.Type]
	if !exists {
		return
	}

	// Se a BehaviorTree do tipo aceitar o creatures[], ótimo. Caso contrário, ajuste os nodes que precisam.
	tree.Tick(c)
}
