package creature

import (
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
)


var templateRegistry = map[int]func() *creature.Creature{
	1001: mob.NewChineseSoldier,
	// Adicione aqui outros templates quando criar
}

func CreateFromTemplate(templateID int) *creature.Creature {
	createFunc, exists := templateRegistry[templateID]
	if !exists {
		return nil
	}
	return createFunc()
}
