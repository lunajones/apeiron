package factory

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

// Adapter para transformar um core.BehaviorNode em creature.BehaviorTree
type behaviorTreeAdapter struct {
	tree core.BehaviorNode
}

func (a *behaviorTreeAdapter) Tick(c *creature.Creature, ctx interface{}) interface{} {
	realCtx, ok := ctx.(core.AIContext)
	if !ok {
		return nil
	}
	return a.tree.Tick(c, realCtx)
}

func CreateBehaviorTree(cType creature.CreatureType, players []*player.Player, creatures []*creature.Creature) creature.BehaviorTree {
	var tree core.BehaviorNode

	switch cType {
	case creature.Soldier:
		tree = mob.BuildChineseSoldierBT(players, creatures)
	case creature.Wolf:
		tree = mob.BuildChineseWolfBT(players, creatures)
	default:
		log.Printf("[BehaviorFactory] Tipo de criatura %s sem BehaviorTree definida", cType)
		return nil
	}

	// Retorna o adapter, já no formato aceito por creature.BehaviorTree
	return &behaviorTreeAdapter{tree: tree}
}
