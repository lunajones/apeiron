package helper

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type LoggingConditionNode struct {
	message   string
	condition func(c *creature.Creature, ctx interface{}) bool
}

func NewLoggingCondition(message string, condition func(c *creature.Creature, ctx interface{}) bool) *LoggingConditionNode {
	return &LoggingConditionNode{
		message:   message,
		condition: condition,
	}
}

func (n *LoggingConditionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	result := n.condition(c, ctx)
	log.Printf("[LOG-COND] [%s] %s → %v", c.Handle.String(), n.message, result)
	if result {
		return core.StatusSuccess
	}
	return core.StatusFailure
}

func (n *LoggingConditionNode) Reset(c *creature.Creature) {}
