package helper

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node/helper"
	"github.com/lunajones/apeiron/service/creature"
)

type InterruptOnThreatDecorator struct {
	ChildNode          core.BehaviorNode
	InterruptAIState   constslib.AIState
	InterruptAnimation constslib.AnimationState
}

func NewInterruptOnThreatDecorator(
	child core.BehaviorNode,
	interruptAIState constslib.AIState,
	interruptAnimation constslib.AnimationState,
) core.BehaviorNode {
	return &InterruptOnThreatDecorator{
		ChildNode:          child,
		InterruptAIState:   interruptAIState,
		InterruptAnimation: interruptAnimation,
	}
}

func (d *InterruptOnThreatDecorator) Tick(c *creature.Creature, ctx interface{}) interface{} {
	interruptNode := &helper.InterruptIfThreatNearbyNode{
		InterruptAIState:   d.InterruptAIState,
		InterruptAnimation: d.InterruptAnimation,
	}

	status := interruptNode.Tick(c, ctx)
	if status == core.StatusSuccess {
		log.Printf("[AI-DECORATOR] [%s (%s)] Interrompido por ameaça próxima.",
			c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	if d.ChildNode != nil {
		return d.ChildNode.Tick(c, ctx)
	}

	return core.StatusFailure
}

func (d *InterruptOnThreatDecorator) Reset(c *creature.Creature) {
	if d.ChildNode != nil {
		d.ChildNode.Reset(c)
	}
}
