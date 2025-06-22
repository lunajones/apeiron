package node

import (
	"math/rand"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomIdleNode struct{}

func (n *RandomIdleNode) Tick(c *creature.Creature) core.BehaviorStatus {
	actions := []creature.CreatureAction{
		creature.ActionIdle,
		creature.ActionWalk,
		creature.ActionJump,
	}

	chosen := actions[rand.Intn(len(actions))]
	c.SetAction(chosen)

	return core.StatusSuccess
}
