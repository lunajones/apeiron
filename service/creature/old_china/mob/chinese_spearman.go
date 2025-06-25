package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
)

func NewChineseSpearman() *creature.Creature {
	log.Println("[Creature] Initializing chinese spearman...")

	c := &creature.Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			Name:           "Chinese Spearman",
			MaxHP:          100,
			RespawnTimeSec: 30,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    4.0,
			Faction:        "Han Dynasty",
		},
		HP:                     100,
		IsAlive:                true,
		IsCorpse:               false,
		IsHostile:              true,
		PrimaryType:            creature.Soldier,
		Types:                  []creature.CreatureType{creature.Soldier},
		Actions:                []creature.CreatureAction{creature.ActionIdle, creature.ActionWalk, creature.ActionRun, creature.ActionSkill1, creature.ActionDie},
		CurrentAction:          creature.ActionIdle,
		AIState:                creature.AIStateIdle,
		LastStateChange:        time.Now(),
		DynamicCombos:          make(map[creature.CreatureAction][]creature.CreatureAction),
		FieldOfViewDegrees:     100,
		VisionRange:            10,
		HearingRange:           10,
		IsBlind:                false,
		IsDeaf:                 false,
		DetectionRadius:        7.5,
		AttackRange:            2.5,
		SkillCooldowns:         make(map[creature.CreatureAction]time.Time),
		AggroTable:             make(map[string]*aggro.AggroEntry),
		MoveSpeed:              4.2,
		AttackSpeed:            1.2,
		MaxPosture:             85,
		CurrentPosture:         85,
		PostureRegenRate:       1.1,
		PostureBreakDurationSec: 5,
		Strength:               18,
		Dexterity:              16,
		Intelligence:           5,
		Focus:                  7,
		PhysicalDefense:        0.12,
		MagicDefense:           0.04,
		RangedDefense:          0.08,
		ControlResistance:      0.05,
		StatusResistance:       0.05,
		CriticalResistance:     0.1,
		CriticalChance:         0.04,
		FacingDirection:        position.Vector2D{X: 1, Y: 0},
	}

	c.Position = c.GenerateSpawnPosition()

	c.BehaviorTree = core.NewSelectorNode(
		core.NewCooldownDecorator(&node.FleeIfLowHPNode{}, 5*time.Second),
		core.NewCooldownDecorator(&node.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.DetectPlayerNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.MaintainMediumDistanceNode{}, 3*time.Second),
		core.NewCooldownDecorator(&node.UseGroundSkillNode{SkillName: "SpearStorm"}, 4*time.Second),
		core.NewCooldownDecorator(&node.AttackTargetNode{SkillName: "SpearThrust"}, 3*time.Second),
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 5*time.Second),
	)

	return c
}
