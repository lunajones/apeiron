package mob

// import (
// 	"log"
// 	"time"

// 	"github.com/lunajones/apeiron/lib"
// 	"github.com/lunajones/apeiron/lib/handle"
// 	"github.com/lunajones/apeiron/lib/model"
// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/node"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/aggro"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// func NewChineseSpearman() *creature.Creature {
// 	log.Println("[Creature] Initializing chinese spearman...")

// 	id := lib.NewUUID()

// 	c := &creature.Creature{
// 		Handle:     handle.NewEntityHandle(id, 1),
// 		Generation: 1,

// 		Creature: model.Creature{
// 			Name:           "Chinese Spearman",
// 			MaxHP:          100,
// 			RespawnTimeSec: 30,
// 			SpawnPoint: position.Position{
// 				GridX: 0, GridY: 0, OffsetX: 0, OffsetY: 0, Z: 0,
// 			},
// 			SpawnRadius: 4.0,
// 			Faction:     "Han Dynasty",
// 		},

// 		HP:                      100,
// 		IsAlive:                 true,
// 		IsCorpse:                false,
// 		IsHostile:               true,
// 		PrimaryType:             consts.Soldier,
// 		Types:                   []consts.CreatureType{consts.Soldier},
// 		Actions:                 []consts.CreatureAction{consts.ActionIdle, consts.ActionWalk, consts.ActionRun, consts.ActionSkill1, consts.ActionDie},
// 		CurrentAction:           consts.ActionIdle,
// 		AIState:                 consts.AIStateIdle,
// 		LastStateChange:         time.Now(),
// 		DynamicCombos:           make(map[consts.CreatureAction][]consts.CreatureAction),
// 		FieldOfViewDegrees:      100,
// 		VisionRange:             10,
// 		HearingRange:            10,
// 		IsBlind:                 false,
// 		IsDeaf:                  false,
// 		DetectionRadius:         7.5,
// 		AttackRange:             2.5,
// 		SkillCooldowns:          make(map[consts.CreatureAction]time.Time),
// 		AggroTable:              make(map[handle.EntityHandle]*aggro.AggroEntry),
// 		WalkSpeed:               2.5,
// 		RunSpeed:                3.8,
// 		AttackSpeed:             1.2,
// 		MaxPosture:              85,
// 		CurrentPosture:          85,
// 		PostureRegenRate:        1.1,
// 		PostureBreakDurationSec: 5,
// 		Strength:                18,
// 		Dexterity:               16,
// 		Intelligence:            5,
// 		Focus:                   7,
// 		PhysicalDefense:         0.12,
// 		MagicDefense:            0.04,
// 		RangedDefense:           0.08,
// 		ControlResistance:       0.05,
// 		StatusResistance:        0.05,
// 		CriticalResistance:      0.1,
// 		CriticalChance:          0.04,
// 		FacingDirection:         position.Vector2D{X: 1, Y: 0},
// 		Needs: []consts.Need{
// 			{Type: consts.NeedHunger, Value: 0, Threshold: 50},
// 		},
// 		Tags: []consts.CreatureTag{
// 			consts.TagHumanoid,
// 		},
// 	}

// 	c.Position = c.GenerateSpawnPosition()

// 	c.BehaviorTree = core.NewSelectorNode(
// 		core.NewCooldownDecorator(&node.FleeIfLowHPNode{}, 5*time.Second),
// 		// core.NewCooldownDecorator(&node.DetectOtherCreatureNode{}, 2*time.Second),
// 		core.NewCooldownDecorator(&node.DetectPlayerNode{}, 2*time.Second),
// 		core.NewCooldownDecorator(&node.MaintainMediumDistanceNode{}, 3*time.Second),
// 		// core.NewCooldownDecorator(&node.UseGroundSkillNode{SkillName: "SpearStorm"}, 4*time.Second),
// 		core.NewCooldownDecorator(&node.RandomIdleNode{}, 5*time.Second),
// 	)

// 	return c
// }
