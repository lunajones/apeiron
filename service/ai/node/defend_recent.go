package node

import (
	"time"
	"github.com/lunajones/apeiron/service/creature"
)

type DefendIfRecentlyDamagedNode struct{}

func (d *DefendIfRecentlyDamagedNode) Tick(c *creature.Creature) BehaviorStatus {
	if time.Now().Unix()-c.TimeOfDeath < 5 {
		c.SetAction(creature.ActionBlock)
		c.ChangeAIState(creature.AIStateIdle)
		return StatusSuccess
	}
	return StatusFailure
}
