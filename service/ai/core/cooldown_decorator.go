package core

import (
	"time"
	"github.com/lunajones/apeiron/service/creature"
)

type CooldownDecorator struct {
	Child      BehaviorNode
	LastUsed   map[string]int64
	CooldownMS int64
}

func NewCooldownDecorator(child BehaviorNode, cooldownMS int64) *CooldownDecorator {
	return &CooldownDecorator{
		Child:      child,
		LastUsed:   make(map[string]int64),
		CooldownMS: cooldownMS,
	}
}

func (d *CooldownDecorator) Tick(c *creature.Creature, ctx interface{}) interface{} {
	now := time.Now().UnixMilli()
	last := d.LastUsed[c.ID]
	if now-last < d.CooldownMS {
		return StatusFailure
	}

	status := d.Child.Tick(c, ctx).(BehaviorStatus)
	if status == StatusSuccess {
		d.LastUsed[c.ID] = now
	}
	return status
}