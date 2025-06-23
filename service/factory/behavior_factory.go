package behaviorfactory

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

func BuildChineseSoldierBT(players []player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return mob.BuildChineseSoldierBT(players, creatures)
}

func BuildChineseSpearmanBT(players []player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return mob.BuildChineseSpearmanBT(players, creatures)
}

func CreateBehaviorTree(cType creature.CreatureType, players []player.Player, creatures []*creature.Creature) core.BehaviorNode {
	switch cType {
	case creature.Soldier:
		return BuildChineseSoldierBT(players, creatures)
	case creature.ChineseSpearman:
		return BuildChineseSpearmanBT(players, creatures)
	default:
		log.Printf("[BehaviorFactory] Tipo de criatura %s sem BehaviorTree definida", cType)
		return nil
	}
}
