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

func NewChineseSoldier(spawnPoint position.Position, spawnRadius float64, ctx *dynamic_context.AIServiceContext) *creature.Creature {
	log.Println("[Creature] Initializing chinese soldier...")

	id := lib.NewUUID()

	c := &creature.Creature{
		Handle:     handle.NewEntityHandle(id, 1),
		Generation: 1,

		Creature: model.Creature{
			Name:           "Chinese Soldier",
			MaxHP:          250,
			RespawnTimeSec: 30,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    5.0,
			Faction:        "Chinese Human",
		},

		MoveCtrl:             movement.NewMovementController(),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},

		PrimaryType: consts.Human,
		Types:       []consts.CreatureType{consts.Human, consts.Soldier},
		HP:          250,
		Alive:       true,
		IsCorpse:    false,
		Hostile:     true,

		AIState:         constslib.AIStateIdle,
		CombatState:     constslib.CombatStateIdle,
		AnimationState:  constslib.AnimationIdle,
		LastStateChange: time.Now(),

		Strength:     20,
		Dexterity:    10,
		Intelligence: 5,
		Focus:        8,

		HitboxRadius:            0.9,
		DesiredBufferDistance:   0.5,
		MinWanderDistance:       3.0,
		MaxWanderDistance:       8.0,
		WanderStopDistance:      0.2,
		FieldOfViewDegrees:      120,
		VisionRange:             15,
		HearingRange:            10,
		SmellRange:              6,
		DetectionRadius:         12.0,
		AttackRange:             2.5,
		WalkSpeed:               1.3,
		RunSpeed:                2.8,
		OriginalRunSpeed:        3.5,
		AttackSpeed:             1.2,
		MaxPosture:              100,
		Posture:                 100,
		PostureRegenRate:        0.1,
		PostureBreakDurationSec: 5,
		PhysicalDefense:         0.15,
		MagicDefense:            0.05,
		RangedDefense:           0.10,
		ControlResistance:       0.1,
		StatusResistance:        0.1,
		CriticalResistance:      0.2,
		CriticalChance:          0.05,

		BlockableChance: 0.85,
		DodgableChance:  0.15,
		DodgeDistance:   1.8,

		DodgeInvulnerabilityDuration: 2300 * time.Millisecond,

		MaxStamina:         100,
		Stamina:            100,
		StaminaRegenPerSec: 10,   // você pode ajustar por tipo: lobo, humano, coelho...
		DodgeStaminaCost:   40.0, // custo padrão de dodge
		MaxBlockDuration:   3 * time.Second,

		RegisteredSkills: []*model.Skill{
			// model.SkillRegistry["SoldierSlash"],
			// model.SkillRegistry["SoldierShieldBash"],
			// model.SkillRegistry["SoldierGroundSlam"],
			// model.SkillRegistry["SoldierLongStep"],
			model.SkillRegistry["SoldierShieldRush"],
		},
		SkillStates: map[constslib.SkillAction]*model.SkillState{
			constslib.Basic:  &model.SkillState{},
			constslib.Skill1: &model.SkillState{},
			constslib.Skill2: &model.SkillState{},
			constslib.Skill3: &model.SkillState{},
			// constslib.Skill4: &model.SkillState{},
			// constslib.Skill5: &model.SkillState{},
		},
		AggroTable: make(map[handle.EntityHandle]*aggro.AggroEntry),

		Needs: []constslib.Need{
			{Type: constslib.NeedHunger, Value: 0, LowThreshold: 0, Threshold: 50},
			{Type: constslib.NeedProvoke, Value: 5, LowThreshold: 0, Threshold: 70},
			{Type: constslib.NeedRecover, Value: 0, LowThreshold: 0, Threshold: 50},
			{Type: constslib.NeedAdvance, Value: 70, LowThreshold: 45, Threshold: 60},
			{Type: constslib.NeedGuard, Value: 30, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedRetreat, Value: 0, LowThreshold: 0, Threshold: 100},
			{Type: constslib.NeedRage, Value: 0, LowThreshold: 0, Threshold: 99},
			{Type: constslib.NeedFake, Value: 0, LowThreshold: 40, Threshold: 60},
			{Type: constslib.NeedPlan, Value: 0, LowThreshold: 40, Threshold: 60},
		},
		Tags: []consts.CreatureTag{
			consts.TagHumanoid,
		},
		Position:          position.Position{X: 0, Y: 0, Z: 0},
		LastPosition:      position.Position{X: 0, Y: 0, Z: 0},
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
				&helper.OnlyIfCloseAndNotMovingNode{Node: &offensive.GetApproachNodeForTagNode{}},
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
