package offensive

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node/humanoid"
	"github.com/lunajones/apeiron/service/ai/node/neutral"
	"github.com/lunajones/apeiron/service/ai/node/predator"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type GetApproachNodeForTagNode struct{}

func (n *GetApproachNodeForTagNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	drive := c.GetCombatDrive()
	if drive.Counter < 0.2 {
		return core.StatusFailure
	}

	var node core.BehaviorNode

	switch {
	case c.HasTag(consts.TagPredator):
		node = &predator.PredatorHopApproachNode{}
	case c.HasTag(consts.TagHumanoid):
		node = &humanoid.SneakBehindTargetNode{}
	case c.HasTag(consts.TagUndead):
		node = &neutral.ApproachUntilInRangeNode{}
	default:
		node = &ChaseUntilInRangeNode{}
	}

	// Executa o node retornado
	return node.Tick(c, ctx)
}

func (n *GetApproachNodeForTagNode) Reset(c *creature.Creature) {

}
