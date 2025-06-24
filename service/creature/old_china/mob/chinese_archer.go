package mob

import (
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
)

func NewChineseArcher() *creature.Creature {
	return &creature.Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			Name:           "Chinese Archer",
			HP:             250,
			MaxHP:          250,
			RespawnTimeSec: 60,
			Position:       position.Position{},
			SpawnPoint:     position.Position{},
			SpawnRadius:    5.0,
			IsAlive:        true,
		},
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionCombo1,
			creature.ActionDie,
		},
		LastStateChange: time.Now(),
	}
}
