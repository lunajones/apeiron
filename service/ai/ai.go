package ai

import (
	"github.com/lunajones/apeiron/service/creature"
)

func Init() {
	InitBehaviorRules()
}

func ProcessAI(c *creature.Creature) {
	if !c.IsAlive || c.IsPostureBroken || c.BehaviorTree == nil {
		return
	}

	c.BehaviorTree.Tick(c)
}
