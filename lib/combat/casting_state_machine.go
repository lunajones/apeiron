package combat

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/lunajones/apeiron/lib/consts"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CastingStateMachine struct {
	creature model.Attacker
	target   model.Targetable
	state    *model.SkillState
	now      time.Time
}

type CastingConfig struct {
	Creature model.Attacker
	Target   model.Targetable
	Skill    *model.Skill
	State    *model.SkillState
	Now      time.Time
	Allies   []model.Targetable
}

func ProcessCastingFSM(cfg *CastingConfig, svcCtx *dynamic_context.AIServiceContext) {
	// Inicializa state se estiver nulo
	if cfg.State == nil {
		cfg.State = cfg.Creature.InitSkillState(cfg.Skill.Action, cfg.Now)
	}

	// Marca como InUse se ainda não foi
	if !cfg.State.InUse {
		fmt.Printf("\033[38;5;81m[CAST-FSM] [%s] FSM iniciando skill %s\033[0m\n", cfg.Creature.GetPrimaryType(), cfg.Skill.Name)
		cfg.State.InUse = true
	}

	fsm := &CastingStateMachine{
		creature: cfg.Creature,
		target:   cfg.Target,
		state:    cfg.State,
		now:      cfg.Now,
	}
	fsm.Process(cfg.Skill, svcCtx)
}

func (fsm *CastingStateMachine) Process(skill *model.Skill, svcCtx *dynamic_context.AIServiceContext) {
	c := fsm.creature
	state := fsm.state
	now := fsm.now

	// Valida alvo
	if fsm.target == nil || !fsm.target.IsAlive() {
		fmt.Printf("\033[91m[CAST-FSM] [%s] cancelando cast: alvo inválido ou morto\033[0m\n", c.GetPrimaryType())
		c.CancelCurrentSkill()
		c.SetCombatState(constslib.CombatStateMoving)
		return
	}

	// WINDUP
	if !state.WindUpFired {
		fmt.Printf("\033[38;5;45m[CAST-FSM] [%s] Iniciando WINDUP para skill %s\033[0m\n", c.GetPrimaryType(), skill.Name)
		state.WindUpUntil = now.Add(time.Duration(skill.WindUpTime * float64(time.Second)))
		state.WindUpFired = true
		c.SetAnimationState(consts.AnimationWindup)
		return
	}
	if now.Before(state.WindUpUntil) {
		return
	}

	// CAST
	if !state.CastFired {
		fmt.Printf("\033[38;5;51m[CAST-FSM] [%s] Iniciando CAST para skill %s\033[0m\n", c.GetPrimaryType(), skill.Name)

		if rand.Float64() < 0.5 {
			c.AddRecentAction(consts.CombatActionHesitatedAttack)
		}

		if fsm.target.IsBlocking() && !skill.CanCastWhileBlocking {
			c.CancelCurrentSkill()
			c.SetCombatState(constslib.CombatStateMoving)
			return
		}

		for _, evt := range fsm.target.GetCombatEvents() {
			if evt.BehaviorType == "ParryFailed" || evt.BehaviorType == "BlockBroken" {
				if time.Since(evt.Timestamp) < 1*time.Second {
					break
				}
			}
		}

		allies := finder.FindNearbyAllies(svcCtx, c, c.GetFaction(), 8.0)
		for _, ally := range allies {
			if ally != nil && ally.GetHandle() != c.GetHandle() && ally.GetCombatDrive().Termination > 0.6 {
				break
			}
		}

		state.CastUntil = now.Add(time.Duration(skill.CastTime * float64(time.Second)))
		state.CastFired = true
		c.SetAnimationState(consts.AnimationCast)

		fsm.target.RegisterCombatEvent(model.CombatEvent{
			SourceHandle:   c.GetHandle(),
			BehaviorType:   "AggressiveIntention",
			Timestamp:      now,
			ExpectedImpact: state.CastUntil,
		})

		return
	}
	if now.Before(state.CastUntil) {
		return
	}

	// RECOVERY
	if !state.RecoveryFired {
		fmt.Printf("\033[38;5;87m[CAST-FSM] [%s] Iniciando RECOVERY para skill %s\033[0m\n", c.GetPrimaryType(), skill.Name)

		UseSkill(c, fsm.target, fsm.target.GetPosition(), skill, nil, nil, svcCtx.NavMesh, svcCtx)

		if skill.Movement != nil && c.GetSkillMovementState() == nil {
			c.SetSkillMovementState(ApplySkillMovement(c, fsm.target, skill))
		}

		state.RecoveryUntil = now.Add(time.Duration(skill.RecoveryTime * float64(time.Second)))
		state.CooldownUntil = now.Add(time.Duration(skill.CooldownSec * float64(time.Second)))
		state.RecoveryFired = true

		c.SetAnimationState(consts.AnimationRecovery)
		return
	}
	if now.Before(state.RecoveryUntil) {
		return
	}

	fmt.Printf("\033[36m[CAST-FSM] [%s] Finalizou cast da skill %s — limpando estado e saindo do FSM\033[0m\n", c.GetPrimaryType(), skill.Name)
	c.SetCombatState(consts.CombatStateMoving)
	c.CancelCurrentSkill()
}
