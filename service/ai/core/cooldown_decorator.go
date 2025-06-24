package core

import (
	"time"

	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type CooldownDecorator struct {
	Child          BehaviorNode
	LastExecution  time.Time
	CooldownPeriod time.Duration
}

func (n *CooldownDecorator) Tick(c *creature.Creature, ctx ai_context.AIContext) interface{} {
	if time.Since(n.LastExecution) < n.CooldownPeriod {
		return StatusFailure
	}

	status := n.Child.Tick(c, ctx).(BehaviorStatus)
	if status == StatusSuccess {
		n.LastExecution = time.Now()
	}
	return status
}

func NewCooldownDecorator(child BehaviorNode, cooldown time.Duration) *CooldownDecorator {
	return &CooldownDecorator{
		Child:          child,
		CooldownPeriod: cooldown,
		LastExecution:  time.Time{},
	}
}