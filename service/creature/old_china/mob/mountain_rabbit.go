package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	prey "github.com/lunajones/apeiron/service/ai/node/prey"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func NewMountainRabbit(spawnPoint position.Position, spawnRadius float64) *creature.Creature {
	log.Println("[Creature] Initializing mountain rabbit...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:               handle.NewEntityHandle(id, 1),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},
		Creature: model.Creature{
			Name:           "Mountain Rabbit",
			MaxHP:          25,
			RespawnTimeSec: 15,
			SpawnPoint:     spawnPoint,
			SpawnRadius:    spawnRadius,
			Faction:        "Rodents",
		},
		HP:       25,
		MoveCtrl: movement.NewMovementController(),

		HitboxRadius:            0.25,
		DesiredBufferDistance:   0.2,
		IsAlive:                 true,
		IsCorpse:                false,
		IsHostile:               false,
		PrimaryType:             consts.Rabbit,
		Types:                   []consts.CreatureType{consts.Rabbit},
		Actions:                 []consts.CreatureAction{consts.ActionIdle, consts.ActionWalk, consts.ActionRun, consts.ActionDie},
		CurrentAction:           consts.ActionIdle,
		AIState:                 consts.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[consts.CreatureAction][]consts.CreatureAction),
		FieldOfViewDegrees:      140,
		VisionRange:             10,
		HearingRange:            14,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         12.0,
		AttackRange:             0.0,
		SkillCooldowns:          make(map[consts.CreatureAction]time.Time),
		AggroTable:              make(map[handle.EntityHandle]*aggro.AggroEntry),
		WalkSpeed:               1.0,
		RunSpeed:                3.2,
		AttackSpeed:             0.0,
		MaxPosture:              20,
		CurrentPosture:          20,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 2,
		Strength:                1,
		Dexterity:               18,
		Intelligence:            2,
		Focus:                   6,
		PhysicalDefense:         0.01,
		MagicDefense:            0.0,
		RangedDefense:           0.02,
		ControlResistance:       0.1,
		StatusResistance:        0.0,
		CriticalResistance:      0.05,
		CriticalChance:          0.0,
		Needs: []consts.Need{
			{Type: consts.NeedHunger, Value: 44, Threshold: 50},
			{Type: consts.NeedSleep, Value: 10, Threshold: 50},
		},
		Tags:            []consts.CreatureTag{consts.TagAnimal, consts.TagPrey, consts.TagCoward},
		FacingDirection: position.Vector2D{X: 0, Y: 1},
	}

	c.Position = c.GenerateSpawnPosition()

	tree := core.NewStateSelectorNode()

	tree.AddSubtree(consts.AIStateIdle, core.NewSelectorNode(
		core.NewCooldownDecorator(&node.EvaluateNeedsNode{
			PriorityOrder: []consts.NeedType{
				consts.NeedHunger,
				consts.NeedThirst,
				consts.NeedSleep,
			},
		}, 3*time.Second),
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 5*time.Second),
		core.NewCooldownDecorator(&node.RandomWanderNode{}, 5*time.Second),
	))

	tree.AddSubtree(consts.AIStateAlert, core.NewSequenceNode(
		core.NewCooldownDecorator(&prey.DetectOtherCreatureNode{}, 2*time.Second),
	))

	tree.AddSubtree(consts.AIStateSearchFood, core.NewSelectorNode(
		core.NewCooldownDecorator(&node.MoveTowardsCorpseNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.FeedOnCorpseNode{}, 3*time.Second),
		core.NewCooldownDecorator(&prey.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.RandomWanderNode{}, 5*time.Second),
	))

	tree.AddSubtree(consts.AIStateFeeding, core.NewSequenceNode(
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 3*time.Second),
		core.NewCooldownDecorator(&node.EvaluateNeedsNode{
			PriorityOrder:  []consts.NeedType{consts.NeedHunger, consts.NeedSleep},
			CheckOnlyThese: []consts.NeedType{consts.NeedHunger, consts.NeedSleep},
		}, 3*time.Second),
	))

	tree.AddSubtree(consts.AIStateDrowsy, core.NewSequenceNode(
		core.NewCooldownDecorator(&prey.FindSafePlaceToSleepNode{}, 2*time.Second),
		core.NewCooldownDecorator(&prey.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.RandomIdleNode{}, 3*time.Second),
		core.NewCooldownDecorator(&node.SleepNode{}, 3*time.Second),
	))

	tree.AddSubtree(consts.AIStateSleeping, core.NewSequenceNode(
		core.NewCooldownDecorator(&node.WakeIfThreatNearbyNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.RegenerateSleepNode{}, 2*time.Second),
	))

	tree.AddSubtree(consts.AIStateFleeing, core.NewParallelNode(
		core.NewCooldownDecorator(&prey.FleeFromThreatNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.FindIfSafeNode{SafeDistance: 10.0}, 2*time.Second),
		// Checa necessidades, mas só age se não houver ameaça válida
		core.NewCooldownDecorator(&node.EvaluateNeedsNode{
			PriorityOrder: []consts.NeedType{
				consts.NeedSleep,
			},
			CheckOnlyThese: []consts.NeedType{
				consts.NeedSleep,
			},
		}, 3*time.Second),
	))

	c.BehaviorTree = tree

	return c
}
