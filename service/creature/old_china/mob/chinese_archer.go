package mob

import (
	"time"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature/aggro"
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
		SkillCooldowns:          make(map[creature.CreatureAction]time.Time),
		AggroTable:              make(map[string]*aggro.AggroEntry),
		FacingDirection: position.Vector2D{X: 1, Y: 0}, // Exemplo: olhando pro eixo X positivo
	}
}
