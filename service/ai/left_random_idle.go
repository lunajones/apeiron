package ai

import (
	"log"
	"math/rand"

	"github.com/lunajones/apeiron/service/creature"
)

type RandomIdleBehaviorNode struct{}

func (r *RandomIdleBehaviorNode) Tick(c *creature.Creature) BehaviorStatus {
	choice := rand.Float32()
	if choice < 0.5 {
		log.Printf("[AI] Creature %s decidiu andar em Idle.", c.ID)
		c.SetAction(creature.ActionWalk)
	} else {
		log.Printf("[AI] Creature %s continua parado.", c.ID)
		c.SetAction(creature.ActionIdle)
	}
	return StatusSuccess
}
