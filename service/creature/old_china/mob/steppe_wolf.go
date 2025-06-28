package mob

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/ai/node/predator"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func NewSteppeWolf(spawnPoint position.Position, spawnRadius float64) *creature.Creature {
	log.Println("[Creature] Initializing steppe wolf...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:               handle.NewEntityHandle(id, 1),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},
		Generation:           1,
		Creature: model.Creature{
			Name:           "Steppe Wolf",
			MaxHP:          80,
			RespawnTimeSec: 25,
			SpawnPoint:     spawnPoint,
			SpawnRadius:    spawnRadius,
			Faction:        "Beasts",
		},
		MoveCtrl:                movement.NewMovementController(),
		HP:                      40,
		IsAlive:                 true,
		IsCorpse:                false,
		IsHostile:               true,
		PrimaryType:             consts.Wolf,
		Types:                   []consts.CreatureType{consts.Wolf},
		Actions:                 []consts.CreatureAction{consts.ActionIdle, consts.ActionWalk, consts.ActionRun, consts.ActionSkill1, consts.ActionDie},
		CurrentAction:           consts.ActionIdle,
		AIState:                 consts.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[consts.CreatureAction][]consts.CreatureAction),
		FieldOfViewDegrees:      100,
		VisionRange:             12,
		HearingRange:            12,
		SmellRange:              8,
		HitboxRadius:            0.75,
		DesiredBufferDistance:   0.4,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         12.0,
		AttackRange:             1.5,
		Skills:                  []string{"Bite", "Lacerate"},
		SkillCooldowns:          make(map[consts.CreatureAction]time.Time),
		AggroTable:              make(map[handle.EntityHandle]*aggro.AggroEntry),
		WalkSpeed:               2.5,
		RunSpeed:                4.5,
		AttackSpeed:             1.0,
		MaxPosture:              80,
		CurrentPosture:          80,
		PostureRegenRate:        1.2,
		PostureBreakDurationSec: 4,
		Strength:                15,
		Dexterity:               20,
		Intelligence:            3,
		Focus:                   6,
		PhysicalDefense:         0.10,
		MagicDefense:            0.02,
		RangedDefense:           0.05,
		ControlResistance:       0.05,
		StatusResistance:        0.05,
		CriticalResistance:      0.1,
		CriticalChance:          0.03,
		Needs: []consts.Need{
			{Type: consts.NeedHunger, Value: 75, Threshold: 30},
			{Type: consts.NeedSleep, Value: 30, Threshold: 80},
		},
		Tags:            []consts.CreatureTag{consts.TagPredator},
		FacingDirection: position.Vector2D{X: 1, Y: 0},
		Position:        spawnPoint.RandomWithinRadius(spawnRadius),
		LastPosition:    spawnPoint,
		Stagger:         physics.StaggerData{},
		Invincibility:   physics.InvincibilityData{},
		ActiveEffects:   []consts.ActiveEffect{},
		DamageWeakness:  make(map[consts.DamageType]float32),
	}

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
		core.NewCooldownDecorator(&predator.DetectOtherCreatureNode{}, 2*time.Second),
	))

	tree.AddSubtree(consts.AIStateChasing, core.NewParallelNode(
		core.NewCooldownDecorator(&node.UseRandomOffensiveSkillNode{}, 3*time.Second),
		&node.MoveCreatureTowardsTargetNode{},
	))

	tree.AddSubtree(consts.AIStateSearchFood, core.NewParallelNode(
		core.NewCooldownDecorator(&node.FeedOnCorpseNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.MoveTowardsCorpseNode{}, 2*time.Second),
		core.NewCooldownDecorator(&predator.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.RandomWanderNode{}, 5*time.Second),
	))

	tree.AddSubtree(consts.AIStateSearchFood, core.NewSelectorNode(
		core.NewCooldownDecorator(&predator.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewSequenceNode(
			core.NewCooldownDecorator(&node.FeedOnCorpseNode{}, 2*time.Second),
			core.NewCooldownDecorator(&node.MoveTowardsCorpseNode{}, 2*time.Second),
		),
		core.NewCooldownDecorator(&node.RandomWanderNode{}, 5*time.Second),
	))

	tree.AddSubtree(consts.AIStateDrowsy, core.NewSequenceNode(
		core.NewCooldownDecorator(&predator.FindSafePlaceToSleepNode{}, 2*time.Second),
		core.NewCooldownDecorator(&predator.DetectOtherCreatureNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.SleepNode{}, 2*time.Second),
	))

	tree.AddSubtree(consts.AIStateSleeping, core.NewSequenceNode(
		core.NewCooldownDecorator(&node.RegenerateSleepNode{}, 2*time.Second),
		core.NewCooldownDecorator(&node.WakeIfThreatNearbyNode{}, 2*time.Second),
	))

	c.BehaviorTree = tree

	return c
}
