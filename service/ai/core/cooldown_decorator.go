package core

import (
	"time"

	"github.com/lunajones/apeiron/service/creature"
)

type CooldownDecorator struct {
	Child         BehaviorNode
	Cooldown      time.Duration
	lastExecution map[string]time.Time
}

func NewCooldownDecorator(child BehaviorNode, cooldown time.Duration) *CooldownDecorator {
	return &CooldownDecorator{
		Child:         child,
		Cooldown:      cooldown,
		lastExecution: make(map[string]time.Time),
	}
}

func (n *CooldownDecorator) Tick(c *creature.Creature, ctx AIContext) BehaviorStatus {
	last, exists := n.lastExecution[c.ID]
	if exists && time.Since(last) < n.Cooldown {
		return StatusFailure
	}

	n.lastExecution[c.ID] = time.Now()
	return n.Child.Tick(c, ctx)
}
