package core

import (
	"time"

	"github.com/lunajones/apeiron/service/creature"
)

type CooldownDecorator struct {
	Node          BehaviorNode
	Cooldown      time.Duration
	LastExecution time.Time
}

func NewCooldownDecorator(node BehaviorNode, cooldown time.Duration) BehaviorNode {
	return &CooldownDecorator{
		Node:          node,
		Cooldown:      cooldown,
		LastExecution: time.Unix(0, 0),
	}
}

func (d *CooldownDecorator) Tick(c *creature.Creature, ctx interface{}) interface{} {
	now := time.Now()
	if now.Sub(d.LastExecution) < d.Cooldown {
		return StatusFailure
	}

	d.LastExecution = now
	return d.Node.Tick(c, ctx)
}

func (d *CooldownDecorator) Reset(c *creature.Creature) {
	if d.Node != nil {
		d.Node.Reset(c)
	}
}
