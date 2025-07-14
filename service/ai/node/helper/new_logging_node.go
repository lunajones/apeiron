package helper

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type LoggingNode struct {
	label string
	node  core.BehaviorNode
}

func NewLoggingNode(label string, node core.BehaviorNode) *LoggingNode {
	return &LoggingNode{
		label: label,
		node:  node,
	}
}

func (n *LoggingNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	status := n.node.Tick(c, ctx)
	log.Printf("[LOG-NODE] [%s] %s → %v", c.Handle.String(), n.label, status)
	return status
}

func (n *LoggingNode) Reset(c *creature.Creature) {
	n.node.Reset(c)
}
