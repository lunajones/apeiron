package creature

import (
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
)

// Agora cada função recebe posição e raio de spawn
var templateRegistry = map[int]func(position.Position, float64) *creature.Creature{
	1001: func(pos position.Position, radius float64) *creature.Creature {
		return mob.NewChineseSoldier(pos, radius)
	},
	1002: func(pos position.Position, radius float64) *creature.Creature {
		return mob.NewSteppeWolf(pos, radius)
	},
	// 1003: func(pos position.Position, radius float64) *creature.Creature { return mob.NewChineseArcher() },
	// 1004: func(pos position.Position, radius float64) *creature.Creature { return mob.NewChineseSpearman() },
	// 1005: func(pos position.Position, radius float64) *creature.Creature { return mob.NewTerrifiedConcubine() },
	1006: func(pos position.Position, radius float64) *creature.Creature {
		return mob.NewMountainRabbit(pos, radius)
	},
}

func CreateFromTemplate(templateID int, spawnPoint position.Position, spawnRadius float64) *creature.Creature {
	createFunc, exists := templateRegistry[templateID]
	if !exists {
		return nil
	}
	return createFunc(spawnPoint, spawnRadius)
}
