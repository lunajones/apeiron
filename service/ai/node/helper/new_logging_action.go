package helper

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type LoggingActionNode struct {
	message string
	action  func(c *creature.Creature, ctx interface{}) core.BehaviorStatus
}

func NewLoggingAction(message string, action func(c *creature.Creature, ctx interface{}) core.BehaviorStatus) *LoggingActionNode {
	return &LoggingActionNode{
		message: message,
		action:  action,
	}
}

func (n *LoggingActionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	status := n.action(c, ctx)
	log.Printf("[LOG-ACTION] [%s] %s → %v", c.Handle.String(), n.message, status)
	return status
}

func (n *LoggingActionNode) Reset(c *creature.Creature) {}
