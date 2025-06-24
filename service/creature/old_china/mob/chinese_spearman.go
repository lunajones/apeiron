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

func NewChineseSpearman() *creature.Creature {
	log.Println("[Creature] Initializing chinese spearman...")

	c := &creature.Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			Name:           "Chinese Spearman",
			HP:             120,
			MaxHP:          120,
			IsAlive:        true,
			IsCorpse:       false,
			RespawnTimeSec: 30,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    5.0,
			Faction:        "Monsters",
			IsHostile:      true,
		},
		PrimaryType: creature.Human,
		Types:       []creature.CreatureType{creature.Human, creature.Soldier},
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionRun,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionDie,
		},
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
		AttackSpeed:             1.1,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		Strength:                18,
		Dexterity:               12,
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
		FacingDirection: position.Vector2D{X: 1, Y: 0},
	}

	c.Position = c.GenerateSpawnPosition()
	return c
}
