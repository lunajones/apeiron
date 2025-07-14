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
	decorator "github.com/lunajones/apeiron/service/ai/core/decorator"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/ai/node/defensive"
	"github.com/lunajones/apeiron/service/ai/node/helper"
	"github.com/lunajones/apeiron/service/ai/node/neutral"
	"github.com/lunajones/apeiron/service/ai/node/offensive"
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
			MaxHP:          150,
			RespawnTimeSec: 25,
			SpawnPoint:     spawnPoint,
			SpawnRadius:    spawnRadius,
			Faction:        "Beasts",
		},
		MoveCtrl:              movement.NewMovementController(),
		HP:                    150,
		Alive:                 true,
		IsCorpse:              false,
		Hostile:               true,
		PrimaryType:           consts.Wolf,
		Types:                 []consts.CreatureType{consts.Wolf},
		AIState:               constslib.AIStateIdle,
		CombatState:           constslib.CombatStateIdle,
		AnimationState:        constslib.AnimationIdle,
		LastStateChange:       time.Now(),
		Strength:              15,
		Dexterity:             20,
		Intelligence:          3,
		Focus:                 6,
		HitboxRadius:          0.75,
		DesiredBufferDistance: 0.4,
		MinWanderDistance:     4.0,
		MaxWanderDistance:     10.0,
		WanderStopDistance:    0.2,
		FieldOfViewDegrees:    100,
		VisionRange:           12,
		HearingRange:          12,
		SmellRange:            14,
		DetectionRadius:       24.0,
		AttackRange:           1.5,
		WalkSpeed:             2.5,
		RunSpeed:              4.5,
		BlockableChance:       0.0,
		DodgableChance:        1.0,
		DodgeDistance:         3.5,

		DodgeInvulnerabilityDuration: 2300 * time.Millisecond,

		OriginalRunSpeed:        4.5,
		AttackSpeed:             1.0,
		MaxPosture:              100,
		Posture:                 100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		PhysicalDefense:         0.10,
		MagicDefense:            0.02,
		RangedDefense:           0.05,
		ControlResistance:       0.05,
		StatusResistance:        0.05,
		CriticalResistance:      0.1,
		CriticalChance:          0.03,

		// NOVOS CAMPOS DE STAMINA
		MaxStamina:         100,
		Stamina:            100,
		StaminaRegenPerSec: 10,
		DodgeStaminaCost:   10.0, // custo padrão de dodge
		RegisteredSkills: []*model.Skill{
			model.SkillRegistry["Bite"],
			model.SkillRegistry["Lacerate"],
		},
		SkillStates: map[constslib.SkillAction]*model.SkillState{
			constslib.Basic:  &model.SkillState{},
			constslib.Skill1: &model.SkillState{},
		},
		AggroTable: make(map[handle.EntityHandle]*aggro.AggroEntry),
		Needs: []constslib.Need{
			{Type: constslib.NeedHunger, Value: 90, LowThreshold: 0, Threshold: 80},
			{Type: constslib.NeedThirst, Value: 10, LowThreshold: 0, Threshold: 80},
			{Type: constslib.NeedSleep, Value: 30, LowThreshold: 0, Threshold: 80},
			{Type: constslib.NeedProvoke, Value: 5, LowThreshold: 0, Threshold: 50},
			{Type: constslib.NeedRecover, Value: 0, LowThreshold: 0, Threshold: 60},
			{Type: constslib.NeedAdvance, Value: 30, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedGuard, Value: 60, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedRetreat, Value: 5, LowThreshold: 0, Threshold: 30},
			{Type: constslib.NeedFake, Value: 10, LowThreshold: 0, Threshold: 20},
			{Type: constslib.NeedPlan, Value: 0, LowThreshold: 0, Threshold: 30},
			{Type: constslib.NeedRage, Value: 0, LowThreshold: 20, Threshold: 40},
		},
		Tags:              []consts.CreatureTag{consts.TagPredator},
		FacingDirection:   position.Vector2D{X: 1, Z: 0},
		Position:          spawnPoint.RandomWithinRadius(spawnRadius),
		LastPosition:      spawnPoint,
		ActiveEffects:     []constslib.ActiveEffect{},
		DamageWeakness:    make(map[constslib.DamageType]float32),
		LastKnownDistance: 0,
	}

	for _, skill := range c.RegisteredSkills {
		if skill == nil {
			continue
		}

		state, exists := c.SkillStates[skill.Action]

		if !exists {
			state = &model.SkillState{}
			c.SkillStates[skill.Action] = state
		}
		state.Skill = skill
		state.ChargesLeft = 1
	}

	tree := core.NewStateSelectorNode()
	tree.AddSubtree(constslib.AIStateIdle,
		decorator.NewInterruptOnThreatDecorator(
			core.NewSelectorNode(
				core.NewCooldownDecorator(
					&node.EvaluateNeedsNode{
						PriorityOrder: []constslib.NeedType{
							constslib.NeedHunger,
							constslib.NeedSleep,
						},
					},
					2*time.Second,
				),
				core.NewCooldownDecorator(
					core.NewCooldownDecorator(&node.WanderNode{
						MaxDistance:      1.5,
						SniffChance:      0.2,
						LookAroundChance: 0.1,
						IdleChance:       0.1,
						ScratchChance:    0.05,
						VocalizeChance:   0.05,
						PlayChance:       0.05,
						ThreatChance:     0.05,
						CuriousChance:    0.05,
					}, 3*time.Second),
					3*time.Second,
				),
			),
			constslib.AIStateCombat,
			constslib.AnimationCombatReady,
		),
	)

	tree.AddSubtree(constslib.AIStateSearchFood,
		decorator.NewInterruptOnThreatDecorator(
			core.NewSelectorNode(
				core.NewSequenceNode(
					&predator.SearchPreyNode{
						TargetTags: []consts.CreatureTag{
							consts.TagPrey,
							consts.TagCoward,
						},
					},
				),
				core.NewCooldownDecorator(
					&node.WanderNode{
						MaxDistance:      3.5,
						SniffChance:      0.3,
						LookAroundChance: 0.2,
						IdleChance:       0.1,
						ScratchChance:    0.05,
						VocalizeChance:   0.05,
						PlayChance:       0.05,
						ThreatChance:     0.05,
						CuriousChance:    0.05,
					},
					3*time.Second,
				),
			),
			constslib.AIStateCombat,
			constslib.AnimationCombatReady,
		),
	)

	// tree.AddSubtree(constslib.AIStateChasing,
	// 	decorator.NewInterruptOnThreatDecorator(
	// 		core.NewSequenceNode(
	// 			core.NewSelectorNode(
	// 				helper.NewConditionNode(func(c *creature.Creature, ctx interface{}) bool {
	// 					return c.NextSkillToUse != nil
	// 				}),
	// 				&offensive.PlanOffensiveSkillNode{},
	// 			),
	// 			core.NewSelectorNode(
	// 				&offensive.CheckSkillRangeNode{},
	// 				&offensive.ChaseTargetNode{},
	// 			),
	// 			&offensive.SkillStateNode{},
	// 		),
	// 		constslib.AIStateCombat,
	// 		constslib.AnimationCombatReady,
	// 	),
	// )
	tree.AddSubtree(constslib.AIStateCombat,
		core.NewSelectorNode(
			// OFENSIVO
			core.NewSequenceNode(
				core.NewSelectorNode(
					helper.NewConditionNode(func(c *creature.Creature, ctx interface{}) bool {
						return c.NextSkillToUse != nil
					}),
					&offensive.PlanOffensiveSkillNode{},
				),
				&offensive.CheckSkillRangeNode{},
				&offensive.SkillStateNode{},
			),

			// DEFENSIVO — só roda se distância < 3.0 e não estiver se movendo
			core.NewSequenceNode(
				&helper.OnlyIfCloseAndNotMovingNode{Node: &defensive.CounterMoveNode{}},
				&helper.OnlyIfCloseAndNotMovingNode{Node: &defensive.MicroRetreatNode{}},
				&helper.OnlyIfCloseAndNotMovingNode{Node: &defensive.CircleAroundTargetNode{}},
			),

			// POSICIONAMENTO — só roda se distância > 3.0 e não estiver se movendo
			core.NewSelectorNode(
				&helper.OnlyIfFarAndNotMovingNode{Node: &neutral.ApproachUntilInRangeNode{}},
				&helper.OnlyIfFarAndNotMovingNode{Node: &offensive.ChaseUntilInRangeNode{}},
				&helper.OnlyIfFarAndNotMovingNode{Node: &defensive.CircleAroundTargetNode{}},
			),
		),
	)

	c.BehaviorTree = tree

	return c
}
