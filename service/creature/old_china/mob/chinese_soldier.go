package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
)

func NewChineseSoldier() *creature.Creature {
	log.Println("[Creature] Initializing chinese soldier...")

	c := &creature.Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			Name:           "Chinese Soldier",
			MaxHP:          100,
			RespawnTimeSec: 30,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    5.0,
			Faction:        "Monsters",
		},
		PrimaryType: creature.Human,
		Types:       []creature.CreatureType{creature.Human, creature.Soldier},
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionRun,
			creature.ActionParry,
			creature.ActionBlock,
			creature.ActionJump,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionSkill3,
			creature.ActionSkill4,
			creature.ActionSkill5,
			creature.ActionCombo1,
			creature.ActionCombo2,
			creature.ActionCombo3,
			creature.ActionDie,
		},
		HP:             100,
		IsAlive:        true,
		IsCorpse:       false,
		IsHostile:      true,
		CurrentAction:           creature.ActionIdle,
		AIState:                 creature.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[creature.CreatureAction][]creature.CreatureAction),
		FieldOfViewDegrees:      120,
		VisionRange:             15,
		HearingRange:            10,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         10.0,
		AttackRange:             2.5,
		SkillCooldowns:          make(map[creature.CreatureAction]time.Time),
		AggroTable:              make(map[string]*aggro.AggroEntry),
		MoveSpeed:               3.5,
		AttackSpeed:             1.2,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		Strength:                20,
		Dexterity:               10,
		Intelligence:            5,
		Focus:                   8,
		PhysicalDefense:         0.15,
		MagicDefense:            0.05,
		RangedDefense:           0.10,
		ControlResistance:       0.1,
		StatusResistance:        0.1,
		CriticalResistance:      0.2,
		CriticalChance:          0.05,
		Needs: []creature.Need{
			{Type: creature.NeedHunger, Value: 0, Threshold: 50},
		},
		Tags: []creature.CreatureTag{
			creature.TagHumanoid,
		},
		FacingDirection: position.Vector2D{X: 1, Y: 0}, // Exemplo: olhando pro eixo X positivo
	}

	c.Position = c.GenerateSpawnPosition()
	return c
}
