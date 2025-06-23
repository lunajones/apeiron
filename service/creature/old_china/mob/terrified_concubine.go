package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
)

func NewTerrifiedConcubine() *creature.Creature {
	log.Println("[Creature] Initializing terrified concubine...")

	c := &creature.Creature{
		ID:    lib.NewUUID(),
		Name:  "Terrified Concubine",
		Types: []creature.CreatureType{creature.Human},
		Level: creature.Normal,
		HP:    50,
		MaxHP: 50,
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionRun,
			creature.ActionDie,
		},
		CurrentAction:           creature.ActionIdle,
		AIState:                 creature.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[creature.CreatureAction][]creature.CreatureAction),
		IsAlive:                 true,
		IsCorpse:                false,
		RespawnTimeSec:          120,
		SpawnPoint:              position.Position{X: 0, Y: 0, Z: 0},
		SpawnRadius:             3.0,
		FieldOfViewDegrees:      90,
		VisionRange:             10,
		HearingRange:            8,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         8.0,
		AttackRange:             0,
		SkillCooldowns:          make(map[creature.CreatureAction]time.Time),
		AggroTable:              make(map[string]float64),
		MoveSpeed:               2.5,
		AttackSpeed:             0,
		Faction:                 "Civilians",
		IsHostile:               false,
		MaxPosture:              50,
		CurrentPosture:          50,
		PostureRegenRate:        1.0,
		PostureBreakDurationSec: 3,
		Strength:                2,
		Dexterity:               3,
		Intelligence:            4,
		Focus:                   2,
		PhysicalDefense:         0.01,
		MagicDefense:            0.01,
		RangedDefense:           0.01,
		ControlResistance:       0.05,
		StatusResistance:        0.05,
		CriticalResistance:      0.1,
		CriticalChance:          0,
		Needs: []creature.Need{
			{Type: creature.NeedHunger, Value: 0, Threshold: 30},
			{Type: creature.NeedThirst, Value: 0, Threshold: 30},
		},
		Tags: []creature.CreatureTag{
			creature.TagHumanoid,
		},
		DamageWeakness: map[creature.DamageType]float32{
			creature.Piercing: 1.0,
			creature.Magic:    1.0,
		},
	}

	c.Position = c.GenerateSpawnPosition()
	return c
}
