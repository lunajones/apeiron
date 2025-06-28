package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/ai/node/prey"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func NewTerrifiedConcubine() *creature.Creature {
	log.Println("[Creature] Initializing terrified concubine...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:     handle.NewEntityHandle(id, 1),
		Generation: 1,

		Creature: model.Creature{
			Name:           "Terrified Concubine",
			MaxHP:          50,
			RespawnTimeSec: 120,
			SpawnPoint: position.Position{
				GridX: 0, GridY: 0, OffsetX: 0, OffsetY: 0, Z: 0,
			},
			SpawnRadius: 3.0,
		},

		HP:                      50,
		IsAlive:                 true,
		IsCorpse:                false,
		IsHostile:               false,
		PrimaryType:             consts.Human,
		Types:                   []consts.CreatureType{consts.Human},
		Level:                   consts.Normal,
		Actions:                 []consts.CreatureAction{consts.ActionIdle, consts.ActionWalk, consts.ActionRun, consts.ActionDie},
		CurrentAction:           consts.ActionIdle,
		AIState:                 consts.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[consts.CreatureAction][]consts.CreatureAction),
		FieldOfViewDegrees:      90,
		VisionRange:             10,
		HearingRange:            8,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         8.0,
		AttackRange:             0,
		SkillCooldowns:          make(map[consts.CreatureAction]time.Time),
		AggroTable:              make(map[handle.EntityHandle]*aggro.AggroEntry),
		WalkSpeed:               1.0,
		RunSpeed:                2.5,
		AttackSpeed:             0,
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
		Needs: []consts.Need{
			{Type: consts.NeedHunger, Value: 0, Threshold: 30},
			{Type: consts.NeedThirst, Value: 0, Threshold: 30},
		},
		Tags: []consts.CreatureTag{
			consts.TagHumanoid,
		},
		DamageWeakness: map[consts.DamageType]float32{
			consts.Piercing: 1.0,
			consts.Magic:    1.0,
		},
		FacingDirection: position.Vector2D{X: 1, Y: 0},
	}

	c.Position = c.GenerateSpawnPosition()

	tree := core.NewStateSelectorNode()

	tree.AddSubtree(consts.AIStateIdle, core.NewSelectorNode(
		// core.NewCooldownDecorator(&prey.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 6*time.Second),
		core.NewCooldownDecorator(&node.RandomWanderNode{}, 8*time.Second),
	))

	tree.AddSubtree(consts.AIStateFleeing, core.NewSelectorNode(
		core.NewCooldownDecorator(&prey.FleeFromThreatNode{}, 3*time.Second),
	))

	// tree.AddSubtree(consts.AIStateAlert, core.NewSelectorNode(
	// 	core.NewCooldownDecorator(&node.DetectOtherCreatureNode{}, 2*time.Second),
	// ))

	c.BehaviorTree = tree

	return c
}
