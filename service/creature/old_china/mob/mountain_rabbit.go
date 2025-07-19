package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func NewMountainRabbit(spawnPoint position.Position, spawnRadius float64, ctx *dynamic_context.AIServiceContext) *creature.Creature {
	log.Println("[Creature] Initializing mountain rabbit...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:               handle.NewEntityHandle(id, 1),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},
		Generation:           1,
		Creature: model.Creature{
			Name:           "Mountain Rabbit",
			MaxHP:          200,
			RespawnTimeSec: 15,
			SpawnPoint:     spawnPoint,
			SpawnRadius:    spawnRadius,
			Faction:        "Rodents",
		},
		MoveCtrl:                movement.NewMovementController(),
		HP:                      200,
		Alive:                   true,
		IsCorpse:                false,
		Hostile:                 false,
		PrimaryType:             consts.Rabbit,
		Types:                   []consts.CreatureType{consts.Rabbit},
		AIState:                 constslib.AIStateIdle,
		CombatState:             constslib.CombatStateIdle,
		AnimationState:          constslib.AnimationIdle,
		LastStateChange:         time.Now(),
		Strength:                1,
		Dexterity:               18,
		Intelligence:            2,
		Focus:                   6,
		HitboxRadius:            0.25,
		DesiredBufferDistance:   0.2,
		MinWanderDistance:       2.0,
		MaxWanderDistance:       10.0,
		WanderStopDistance:      0.2,
		FieldOfViewDegrees:      140,
		VisionRange:             10,
		HearingRange:            14,
		SmellRange:              6,
		DetectionRadius:         10.0,
		AttackRange:             0.0,
		WalkSpeed:               1.0,
		RunSpeed:                3.2,
		AttackSpeed:             0.0,
		MaxPosture:              20,
		Posture:                 20,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 2,
		PhysicalDefense:         0.01,
		MagicDefense:            0.0,
		RangedDefense:           0.02,
		ControlResistance:       0.1,
		StatusResistance:        0.0,
		CriticalResistance:      0.05,
		CriticalChance:          0.0,

		// NOVOS CAMPOS DE STAMINA
		MaxStamina:         20,
		Stamina:            20,
		StaminaRegenPerSec: 15,

		SkillStates: make(map[constslib.SkillAction]*model.SkillState),
		AggroTable:  make(map[handle.EntityHandle]*aggro.AggroEntry),
		Needs: []constslib.Need{
			{Type: constslib.NeedHunger, Value: 10, Threshold: 50},
			{Type: constslib.NeedSleep, Value: 30, Threshold: 50},
		},
		Tags:              []consts.CreatureTag{consts.TagAnimal, consts.TagPrey, consts.TagCoward},
		Position:          spawnPoint.RandomWithinRadius(spawnRadius),
		LastPosition:      spawnPoint,
		ActiveEffects:     []constslib.ActiveEffect{},
		DamageWeakness:    make(map[constslib.DamageType]float32),
		LastKnownDistance: 0,
	}

	tree := core.NewStateSelectorNode()

	// tree.AddSubtree(constslib.AIStateIdle,
	// 	core.NewSelectorNode(
	// 		&node.InterruptIfAttackedRecentlyNode{
	// 			InterruptAIState:   constslib.AIStateFleeing,
	// 			InterruptAnimation: constslib.AnimationRun,
	// 		},
	// 		core.NewCooldownDecorator(
	// 			&node.EvaluateNeedsNode{
	// 				PriorityOrder: []constslib.NeedType{
	// 					constslib.NeedSleep,
	// 					constslib.NeedHunger,
	// 				},
	// 			},
	// 			2*time.Second,
	// 		),
	// 		core.NewCooldownDecorator(&node.WanderNode{
	// 			MaxDistance:      1.5,
	// 			SniffChance:      0.2,
	// 			LookAroundChance: 0.1,
	// 			IdleChance:       0.1,
	// 			ScratchChance:    0.05,
	// 			VocalizeChance:   0.05,
	// 			PlayChance:       0.05,
	// 			ThreatChance:     0.05,
	// 			CuriousChance:    0.05,
	// 		}, 3*time.Second),
	// 	),
	// )

	// tree.AddSubtree(constslib.AIStateDrowsy,
	// 	core.NewSelectorNode(
	// 		&node.InterruptIfAttackedRecentlyNode{
	// 			InterruptAIState:   constslib.AIStateFleeing,
	// 			InterruptAnimation: constslib.AnimationRun,
	// 		},
	// 		core.NewCooldownDecorator(
	// 			&node.FindSafePlaceToSleepNode{},
	// 			2*time.Second,
	// 		),
	// 	),
	// )

	// tree.AddSubtree(constslib.AIStateSeekingSafePlace,
	// 	core.NewSelectorNode(
	// 		&node.InterruptIfAttackedRecentlyNode{
	// 			InterruptAIState:   constslib.AIStateFleeing,
	// 			InterruptAnimation: constslib.AnimationRun,
	// 		},
	// 		core.NewCooldownDecorator(
	// 			&node.FindSafePlaceToSleepNode{},
	// 			2*time.Second,
	// 		),
	// 	),
	// )

	// tree.AddSubtree(constslib.AIStateSleeping,
	// 	core.NewSelectorNode(
	// 		&node.InterruptIfAttackedRecentlyNode{
	// 			InterruptAIState:   constslib.AIStateFleeing,
	// 			InterruptAnimation: constslib.AnimationRun,
	// 		},
	// 		&node.InterruptIfThreatNearbyNode{
	// 			InterruptAIState:   constslib.AIStateAlert,
	// 			InterruptAnimation: constslib.AnimationWake,
	// 		},
	// 		&node.RegenerateNeedNode{
	// 			NeedType:            constslib.NeedSleep,
	// 			CompletionThreshold: 0.0,
	// 			RegenAmount:         -0.0004,
	// 			OnCompleteAI:        constslib.AIStateIdle,
	// 			OnCompleteAnim:      constslib.AnimationWake,
	// 			RunningAnim:         constslib.AnimationSleep,
	// 		},
	// 	),
	// )

	// tree.AddSubtree(constslib.AIStateAlert,
	// 	core.NewSelectorNode(
	// 		&node.InterruptIfAttackedRecentlyNode{
	// 			InterruptAIState:   constslib.AIStateFleeing,
	// 			InterruptAnimation: constslib.AnimationRun,
	// 		},
	// 		core.NewCooldownDecorator(
	// 			&node.FleeFromThreatNode{
	// 				SafeDistance: 6.0,
	// 			},
	// 			500*time.Millisecond,
	// 		),
	// 	),
	// )

	c.BehaviorTree = tree

	c.UpdateFacingDirection(ctx)

	return c
}
