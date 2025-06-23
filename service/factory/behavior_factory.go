package factory

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

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

func CreateBehaviorTree(types []creature.CreatureType, players []*player.Player, creatures []*creature.Creature) creature.BehaviorTree {
	var tree core.BehaviorNode

	for _, t := range types {
		switch t {
		case creature.Soldier:
			tree = mob.BuildChineseSoldierBT(players, creatures)
			break
		case creature.Wolf:
			tree = mob.BuildChineseWolfBT(players, creatures)
			break
		case creature.Human:
			tree = mob.BuildChineseSoldierBT(players, creatures)
			break
		}
		if tree != nil {
			break
		}
	}

	if tree == nil {
		log.Printf("[BehaviorFactory] Nenhuma BehaviorTree encontrada para tipos %v", types)
		return nil
	}

	return &behaviorTreeAdapter{tree: tree}
}
