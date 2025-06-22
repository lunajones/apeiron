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
	if !c.IsAlive || c.IsPostureBroken {
		return
	}

	tree, exists := behaviorTrees[c.Type]
	if !exists {
		return
	}

	// Agora os nodes internos podem usar a lista creatures se quiserem
	tree.Tick(c)
}
