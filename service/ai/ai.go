package ai

import (
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
)

var behaviorTrees map[creature.CreatureType]BehaviorNode

func Init() {
	InitBehaviorRules()

	behaviorTrees = map[creature.CreatureType]BehaviorNode{
		creature.Soldier: mob.BuildChineseSoldierBT(),
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
