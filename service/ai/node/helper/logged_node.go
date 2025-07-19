package helper

import (
	"github.com/fatih/color"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type LoggedNode struct {
	Name string
	Node core.BehaviorNode
}

func NewLoggedNode(name string, node core.BehaviorNode) *LoggedNode {
	return &LoggedNode{
		Name: name,
		Node: node,
	}
}

func (l *LoggedNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	result := l.Node.Tick(c, ctx)

	switch result {
	case core.StatusSuccess:
		color.Green("[SELECTOR-LOG] [%s] %s → SUCCESS", c.Handle.String(), l.Name)
	case core.StatusRunning:
		color.Cyan("[SELECTOR-LOG] [%s] %s → RUNNING", c.Handle.String(), l.Name)
	case core.StatusFailure:
		color.Red("[SELECTOR-LOG] [%s] %s → FAILURE", c.Handle.String(), l.Name)
	default:
		color.Magenta("[SELECTOR-LOG] [%s] %s → UNKNOWN STATUS", c.Handle.String(), l.Name)
	}

	return result
}

func (l *LoggedNode) Reset(c *creature.Creature) {
	l.Node.Reset(c)
}
