package behaviorfactory

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

func CreateBehaviorTree(cType creature.CreatureType, players []player.Player, creatures []*creature.Creature) creature.BehaviorNode {
	switch cType {
	case creature.Soldier:
		return mob.BuildChineseSoldierBT(players, creatures)
	case creature.ChineseSpearman:
		return mob.BuildChineseSpearmanBT(players, creatures)
	default:
		log.Printf("[BehaviorFactory] Tipo de criatura %s sem BehaviorTree definida", cType)
		return nil
	}
}
