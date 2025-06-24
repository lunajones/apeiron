package mob

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib"
)

func NewChineseArcher() *creature.Creature {
	return &creature.Creature{
		ID:    lib.NewUUID(),
		Name:  "Chinese Archer",
		Types: []creature.CreatureType{
			creature.Human,
			creature.Soldier,
		},
		HP:    250,
		MaxHP: 250,
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionCombo1,
			creature.ActionDie,
		},
		IsAlive:        true,
		RespawnTimeSec: 60,
	}
}
