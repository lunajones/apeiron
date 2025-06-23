package mob

import (
	"github.com/lunajones/apeiron/service/creature"
)

func NewChineseArcher(id string) *creature.Creature {
	return &creature.Creature{
		ID:    id,
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
		DamageWeakness: map[creature.DamageType]float32{
			creature.Piercing: 1.2,
			creature.Magic:    0.8,
		},
		IsAlive:        true,
		RespawnTimeSec: 60,
	}
}
