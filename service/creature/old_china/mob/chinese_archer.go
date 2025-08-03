package mob

import (
	"time"

	"github.com/lunajones/apeiron/lib"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	decorator "github.com/lunajones/apeiron/service/ai/core/decorator"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
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

func NewChineseArcher(spawnPoint position.Position, spawnRadius float64, ctx *dynamic_context.AIServiceContext) *creature.Creature {
	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:     handle.NewEntityHandle(id, 1),
		Generation: 1,

		Creature: model.Creature{
			Name:           "Chinese Archer",
			MaxHP:          180,
			RespawnTimeSec: 30,
			SpawnPoint:     spawnPoint,
			SpawnRadius:    spawnRadius,
			Faction:        "Chinese Human",
		},

		MoveCtrl:             movement.NewMovementController(),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},

		PrimaryType: consts.ChineseArcher,
		Types:       []consts.CreatureType{consts.Human, consts.Soldier},
		HP:          180,
		Alive:       true,
		IsCorpse:    false,
		Hostile:     true,

		AIState:         constslib.AIStateIdle,
		CombatState:     constslib.CombatStateIdle,
		AnimationState:  constslib.AnimationIdle,
		LastStateChange: time.Now(),

		Strength:     8,
		Dexterity:    20,
		Intelligence: 6,
		Focus:        10,

		HitboxRadius:            0.7,
		DesiredBufferDistance:   2.0,
		MinWanderDistance:       4.0,
		MaxWanderDistance:       9.0,
		WanderStopDistance:      0.3,
		FieldOfViewDegrees:      120,
		VisionRange:             20,
		HearingRange:            12,
		SmellRange:              4,
		DetectionRadius:         15.0,
		AttackRange:             8.0,
		WalkSpeed:               1.5,
		RunSpeed:                3.0,
		OriginalRunSpeed:        3.5,
		AttackSpeed:             1.0,
		MaxPosture:              80,
		Posture:                 80,
		PostureRegenRate:        0.08,
		PostureBreakDurationSec: 4,
		PhysicalDefense:         0.10,
		MagicDefense:            0.05,
		RangedDefense:           0.25,
		ControlResistance:       0.05,
		StatusResistance:        0.1,
		CriticalResistance:      0.1,
		CriticalChance:          0.15,

		BlockableChance: 0.0,
		DodgableChance:  1.0,
		DodgeDistance:   2.0,

		DodgeInvulnerabilityDuration: 2500 * time.Millisecond,

		MaxStamina:         90,
		Stamina:            90,
		StaminaRegenPerSec: 12,
		DodgeStaminaCost:   35.0,
		MaxBlockDuration:   2 * time.Second,

		RegisteredSkills: []*model.Skill{
			model.SkillRegistry["ArcherBasicShot"],
			model.SkillRegistry["ArcherPowerShot"],
		},

		SkillStates: map[constslib.SkillAction]*model.SkillState{
			constslib.Basic:  &model.SkillState{},
			constslib.Skill1: &model.SkillState{},
		},

		AggroTable: make(map[handle.EntityHandle]*aggro.AggroEntry),

		Needs: []constslib.Need{
			{Type: constslib.NeedAdvance, Value: 80, LowThreshold: 45, Threshold: 60},
			{Type: constslib.NeedFake, Value: 0, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedPlan, Value: 0, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedRetreat, Value: 0, LowThreshold: 0, Threshold: 100},
		},

		Tags: []consts.CreatureTag{
			consts.TagHumanoid,
		},

		Position:          spawnPoint,
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
				&helper.ValidateCombatStateNode{},

				core.NewSelectorNode(
					helper.NewConditionNode(func(c *creature.Creature, ctx interface{}) bool {
						return c.NextSkillToUse != nil
					}),
					&offensive.PlanOffensiveSkillNode{},
				),
				&offensive.CheckSkillRangeNode{},
			),

			// DEFENSIVO — só roda se distância < 3.0 e não estiver se movendo
			core.NewSelectorNode(
				// &helper.OnlyIfCloseAndNotMovingNode{Node: &offensive.GetApproachNodeForTagNode{}},
				&helper.OnlyIfCloseAndNotMovingNode{Node: &defensive.MicroRetreatNode{}},
				&helper.OnlyIfCloseAndNotMovingNode{Node: &defensive.CircleAroundTargetNode{}},
			),

			// POSICIONAMENTO — só roda se distância > 3.0 e não estiver se movendo
			core.NewSelectorNode(
				&helper.OnlyIfFarAndNotMovingNode{Node: &neutral.ApproachUntilInRangeNode{}},
				&helper.OnlyIfFarAndNotMovingNode{Node: &offensive.ChaseUntilInRangeNode{}},
			),

			// 🛑 SAI DO COMBATE SE NÃO HÁ ALVOS VÁLIDOS
			&neutral.ExitCombatIfNoValidTargetsNode{},
		),
	)

	c.BehaviorTree = tree

	c.UpdateFacingDirection(ctx)

	return c
}
