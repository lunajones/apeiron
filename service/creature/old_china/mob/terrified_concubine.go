package mob

import "github.com/lunajones/apeiron/lib/combat"

func NewTerrifiedConcubine(id string, ctype CreatureType) *Creature {
	return &Creature{
		ID:    id,
		Name:  "Soldier",
		TypLevele:  clevel,
		Types: []CreatureType{
			Human,
			Soldier,
		},
		HP:    250,
		MaxHP: 250,
		Actions: []CreatureAction{
			ActionIdle,
			ActionWalk,
			ActionSkill1,
			ActionSkill2,
			ActionCombo1,
			ActionDie,
		},
		DamageWeakness: map[combat.DamageType]float32{
			combat.Piercing: 1.2,
			combat.Magic:    0.8,
		},
		IsAlive:        true,
		RespawnTimeSec: 60,
	}
}