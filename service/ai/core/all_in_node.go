package core

import "github.com/lunajones/apeiron/service/creature"

type AllInNode struct {
	children []BehaviorNode
}

func NewAllInNode(children ...BehaviorNode) *AllInNode {
	return &AllInNode{children: children}
}

func (n *AllInNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	finalStatus := StatusSuccess
	for _, child := range n.children {
		status := child.Tick(c, ctx)
		if status == StatusRunning {
			finalStatus = StatusRunning
		}
	}
	return finalStatus
}

func (n *AllInNode) Reset(c *creature.Creature) {
	for _, child := range n.children {
		child.Reset(c)
	}
}
