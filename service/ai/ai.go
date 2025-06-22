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
	if !c.IsAlive || c.IsPostureBroken {
		return
	}

	tree, exists := behaviorTrees[c.Type]
	if !exists {
		return
	}

	tree.Tick(c)
}
