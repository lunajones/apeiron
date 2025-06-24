package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
)

func NewChineseSpearman() *creature.Creature {
	log.Println("[Creature] Initializing chinese spearman...")

	c := &creature.Creature{
		ID:    lib.NewUUID(),
		Name:  "Chinese Spearman",
		Types: []creature.CreatureType{creature.Soldier, creature.Human},
		Level: creature.Normal,
		HP:    120,
		MaxHP: 120,
		Actions: []creature.CreatureAction{
			creature.ActionIdle,
			creature.ActionWalk,
			creature.ActionRun,
			creature.ActionParry,
			creature.ActionBlock,
			creature.ActionJump,
			creature.ActionSkill1,
			creature.ActionSkill2,
			creature.ActionCombo1,
			creature.ActionDie,
		},
		CurrentAction:           creature.ActionIdle,
		AIState:                 creature.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[creature.CreatureAction][]creature.CreatureAction),
		IsAlive:                 true,
		IsCorpse:                false,
		RespawnTimeSec:          30,
		SpawnPoint:              position.Position{X: 0, Y: 0, Z: 0},
		SpawnRadius:             5.0,
		FieldOfViewDegrees:      110,
		VisionRange:             12,
		HearingRange:            10,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         10.0,
		AttackRange:             3.0,
		SkillCooldowns:          make(map[creature.CreatureAction]time.Time),
		AggroTable:              make(map[string]*aggro.AggroEntry),
		MoveSpeed:               3.2,
		AttackSpeed:             1.0,
		Faction:                 "Monsters",
		IsHostile:               true,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		Strength:                18,
		Dexterity:               12,
		Intelligence:            5,
		Focus:                   7,
		PhysicalDefense:         0.12,
		MagicDefense:            0.05,
		RangedDefense:           0.10,
		ControlResistance:       0.1,
		StatusResistance:        0.1,
		CriticalResistance:      0.2,
		CriticalChance:          0.05,
		Needs: []creature.Need{
			{Type: creature.NeedHunger, Value: 0, Threshold: 50},
		},
		Tags: []creature.CreatureTag{creature.TagHumanoid},
		DamageWeakness: map[creature.DamageType]float32{
			creature.Piercing: 1.1,
			creature.Magic:    0.9,
		},
		FacingDirection: position.Vector2D{X: 1, Y: 0}, // Exemplo: olhando pro eixo X positivo

	}

	c.Position = c.GenerateSpawnPosition()
	return c
}
