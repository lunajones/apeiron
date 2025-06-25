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
			MaxHP:          250,
			RespawnTimeSec: 60,
			SpawnPoint:     position.Position{},
			SpawnRadius:    5.0,
		},
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionCombo1,
			creature.ActionDie,
		},
		HP:             250,
		LastStateChange: time.Now(),
		IsAlive:        true,
		Position:       position.Position{},
	}
}
