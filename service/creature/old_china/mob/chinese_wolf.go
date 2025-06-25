package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
)

func NewChineseWolf() *creature.Creature {
	log.Println("[Creature] Initializing chinese wolf...")

	c := &creature.Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			Name:           "Chinese Wolf",
			MaxHP:          80,
			RespawnTimeSec: 25,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    4.0,
			Faction:        "Beasts",
		},
		HP:             80,
		IsAlive:        true,
		IsCorpse:       false,
		IsHostile:      true,
		PrimaryType:            creature.Wolf,
		Types:                  []creature.CreatureType{creature.Wolf},
		Actions:                []creature.CreatureAction{creature.ActionIdle, creature.ActionWalk, creature.ActionRun, creature.ActionSkill1, creature.ActionDie},
		CurrentAction:          creature.ActionIdle,
		AIState:                creature.AIStateIdle,
		LastStateChange:        time.Now(),
		DynamicCombos:          make(map[creature.CreatureAction][]creature.CreatureAction),
		FieldOfViewDegrees:     100,
		VisionRange:            12,
		HearingRange:           12,
		IsBlind:                false,
		IsDeaf:                 false,
		DetectionRadius:        8.0,
		AttackRange:            1.5,
		SkillCooldowns:         make(map[creature.CreatureAction]time.Time),
		AggroTable:             make(map[string]*aggro.AggroEntry),
		MoveSpeed:              4.5,
		AttackSpeed:            1.0,
		MaxPosture:             80,
		CurrentPosture:         80,
		PostureRegenRate:       1.2,
		PostureBreakDurationSec: 4,
		Strength:               15,
		Dexterity:              20,
		Intelligence:           3,
		Focus:                  6,
		PhysicalDefense:        0.10,
		MagicDefense:           0.02,
		RangedDefense:          0.05,
		ControlResistance:      0.05,
		StatusResistance:       0.05,
		CriticalResistance:     0.1,
		CriticalChance:         0.03,
		Needs:                  []creature.Need{{Type: creature.NeedHunger, Value: 0, Threshold: 40}},
		Tags:                   []creature.CreatureTag{creature.TagPredator},
		FacingDirection:        position.Vector2D{X: 1, Y: 0},
	}

	c.Position = c.GenerateSpawnPosition()

	c.SetBehavior(core.NewSequenceNode(
	core.NewCooldownDecorator(&node.FleeIfLowHPNode{}, 5*time.Second),
	core.NewCooldownDecorator(&node.FeedOnCorpseNode{}, 3*time.Second),
	core.NewCooldownDecorator(&node.DetectOtherCreatureNode{}, 2*time.Second),
	core.NewCooldownDecorator(&node.DetectPlayerNode{}, 2*time.Second),
	core.NewCooldownDecorator(&node.AttackIfVulnerableNode{SkillName: "Bite"}, 4*time.Second),
	core.NewCooldownDecorator(&node.AttackTargetNode{SkillName: "Bite"}, 3*time.Second),
	core.NewCooldownDecorator(&node.RandomIdleNode{}, 5*time.Second),
))

	return c
}
