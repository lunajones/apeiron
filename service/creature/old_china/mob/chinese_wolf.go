package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
)

func NewChineseWolf() *creature.Creature {
	log.Println("[Creature] Initializing chinese wolf...")

	c := &creature.Creature{
		ID:    lib.NewUUID(),
		Name:  "Chinese Wolf",
		PrimaryType: creature.Wolf,
		Types: []creature.CreatureType{creature.Wolf},
		Level: creature.Normal,
		HP:    80,
		MaxHP: 80,
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionRun,
			creature.ActionSkill1,
			creature.ActionDie,
		},
		CurrentAction:           creature.ActionIdle,
		AIState:                 creature.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[creature.CreatureAction][]creature.CreatureAction),
		IsAlive:                 true,
		IsCorpse:                false,
		RespawnTimeSec:          45,
		SpawnPoint:              position.Position{X: 0, Y: 0, Z: 0},
		SpawnRadius:             8.0,
		FieldOfViewDegrees:      150,
		VisionRange:             20,
		HearingRange:            15,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         12.0,
		AttackRange:             1.5,
		SkillCooldowns:          make(map[creature.CreatureAction]time.Time),
		AggroTable:              make(map[string]float64),
		MoveSpeed:               4.2,
		AttackSpeed:             1.5,
		Faction:                 "Wild",
		IsHostile:               true,
		MaxPosture:              50,
		CurrentPosture:          50,
		PostureRegenRate:        1.2,
		PostureBreakDurationSec: 3,
		Strength:                15,
		Dexterity:               20,
		Intelligence:            3,
		Focus:                   5,
		PhysicalDefense:         0.08,
		MagicDefense:            0.02,
		RangedDefense:           0.05,
		ControlResistance:       0.05,
		StatusResistance:        0.05,
		CriticalResistance:      0.1,
		CriticalChance:          0.1,
		Needs: []creature.Need{
			{Type: creature.NeedHunger, Value: 0, Threshold: 50},
		},
		Tags: []creature.CreatureTag{
			creature.TagAnimal,
			creature.TagPredator,
		},
		DamageWeakness: map[creature.DamageType]float32{
			creature.Piercing: 1.0,
			creature.Magic:    1.1,
		},
	}

	c.Position = c.GenerateSpawnPosition()
	return c
}
