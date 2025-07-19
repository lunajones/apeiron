package creature

import (
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
)

var templateRegistry = map[int]func(position.Position, float64, *dynamic_context.AIServiceContext) *creature.Creature{
	1001: mob.NewChineseSoldier,
	1002: mob.NewSteppeWolf,
	1006: mob.NewMountainRabbit,
}

func CreateFromTemplate(templateID int, spawnPoint position.Position, spawnRadius float64, ctx *dynamic_context.AIServiceContext) *creature.Creature {
	createFunc, exists := templateRegistry[templateID]
	if !exists {
		return nil
	}
	return createFunc(spawnPoint, spawnRadius, ctx)
}
