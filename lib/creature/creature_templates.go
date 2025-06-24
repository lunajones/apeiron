package creature

import (
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
)


var templateRegistry = map[int]func() *creature.Creature{
	1001: mob.NewChineseSoldier,
	1002: mob.NewChineseWolf,
	1003: mob.NewChineseArcher,
	1004: mob.NewChineseSpearman,
	1005: mob.NewTerrifiedConcubine,
}

func CreateFromTemplate(templateID int) *creature.Creature {
	createFunc, exists := templateRegistry[templateID]
	if !exists {
		return nil
	}
	return createFunc()
}
