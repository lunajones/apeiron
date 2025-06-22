package ai

import (
	"github.com/lunajones/apeiron/service/creature"
)

var behaviorTrees map[creature.CreatureType]BehaviorNode

func Init() {
	behaviorTrees = map[creature.CreatureType]BehaviorNode{
		creature.Soldier: BuildChineseSoldierBT(nil, nil), // Modifique se precisar
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
