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
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func NewChineseSoldier() *creature.Creature {
	log.Println("[Creature] Initializing chinese soldier...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:     handle.NewEntityHandle(id, 1),
		Generation: 1,

		Creature: model.Creature{
			Name:           "Chinese Soldier",
			MaxHP:          100,
			RespawnTimeSec: 30,
			SpawnPoint: position.Position{
				GridX: 0, GridY: 0, OffsetX: 0, OffsetY: 0, Z: 0,
			},
			SpawnRadius: 5.0,
			Faction:     "Monsters",
		},

		PrimaryType: consts.Human,
		Types:       []consts.CreatureType{consts.Human, consts.Soldier},
		Actions: []consts.CreatureAction{
			consts.ActionIdle,
			consts.ActionWalk,
			consts.ActionRun,
			consts.ActionParry,
			consts.ActionBlock,
			consts.ActionJump,
			consts.ActionSkill1,
			consts.ActionSkill2,
			consts.ActionSkill3,
			consts.ActionSkill4,
			consts.ActionSkill5,
			consts.ActionCombo1,
			consts.ActionCombo2,
			consts.ActionCombo3,
			consts.ActionDie,
		},
		HP:                 100,
		IsAlive:            true,
		IsCorpse:           false,
		IsHostile:          true,
		CurrentAction:      consts.ActionIdle,
		AIState:            consts.AIStateIdle,
		LastStateChange:    time.Now(),
		DynamicCombos:      make(map[consts.CreatureAction][]consts.CreatureAction),
		FieldOfViewDegrees: 120,
		VisionRange:        15,
		HearingRange:       10,
		IsBlind:            false,
		IsDeaf:             false,
		DetectionRadius:    10.0,
		AttackRange:        2.5,
		SkillCooldowns:     make(map[consts.CreatureAction]time.Time),
		AggroTable:         make(map[handle.EntityHandle]*aggro.AggroEntry),

		WalkSpeed:               2.5,
		RunSpeed:                3.5,
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
		Needs: []consts.Need{
			{Type: consts.NeedHunger, Value: 0, Threshold: 50},
		},
		Tags: []consts.CreatureTag{
			consts.TagHumanoid,
		},
		FacingDirection: position.Vector2D{X: 1, Y: 0},
	}

	c.Position = c.GenerateSpawnPosition()

	c.BehaviorTree = core.NewSelectorNode(
		core.NewCooldownDecorator(&node.FleeIfLowHPNode{}, 5*time.Second),
		//core.NewCooldownDecorator(&predator.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.DetectPlayerNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.MaintainMediumDistanceNode{}, 3*time.Second),
		// core.NewCooldownDecorator(&node.UseGroundSkillNode{SkillName: "SoldierGroundSlam"}, 4*time.Second),
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 5*time.Second),
	)

	return c
}
