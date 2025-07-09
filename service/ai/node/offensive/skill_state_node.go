package offensive

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/combat"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type SkillStateNode struct{}

func (n *SkillStateNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SKILL-STATE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if (c.IsBlocking() || c.IsDodging()) && (c.NextSkillToUse == nil || !c.NextSkillToUse.CanCastWhileBlocking) {
		if c.CurrentSkillState() != nil && c.CurrentSkillState().CanBeCancelled() {
			log.Printf("[SKILL-STATE] [%s] cancelando skill %s por bloqueio ou esquiva",
				c.Handle.String(), c.NextSkillToUse.Name)
			c.CancelCurrentSkill()
		}
		return core.StatusFailure
	}

	state := c.SkillStates[c.NextSkillToUse.Action]

	if state == nil || !state.InUse {
		// Inicializa ciclo
		state = c.InitSkillState(c.NextSkillToUse.Action, time.Now())
		log.Printf("[SKILL-STATE] [%s (%s)] iniciando skill %s", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)
		return core.StatusRunning
	}

	if state.WasInterrupted {
		log.Printf("[SKILL-STATE] [%s (%s)] skill %s foi interrompida", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)
		c.ResetSkillState(c.NextSkillToUse.Action)
		c.NextSkillToUse = nil
		return core.StatusSuccess
	}

	if time.Now().Before(state.WindUpUntil) {
		c.SetAnimationState(constslib.AnimationWindup)

		if !state.WindUpFired {
			state.WindUpFired = true

			target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
			now := time.Now()

			windup := state.WindUpUntil.Sub(now)
			if target != nil {
				svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
					SourceHandle: c.Handle,
					TargetHandle: target.GetHandle(),
					BehaviorType: "AggressiveIntention",
					Timestamp:    now,
					WindupTime:   windup,
				})

				log.Printf("%s", color.New(color.FgHiGreen).Sprintf(
					"[EVENTO] Registrado: AggressiveIntention [%s → %s] windup: %dms",
					c.Handle.String(),
					target.GetHandle().String(),
					windup.Milliseconds(),
				))

			}
		}

		return core.StatusRunning
	}

	if time.Now().Before(state.CastUntil) {

		c.SetAnimationState(constslib.AnimationCast)

		target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

		if target != nil && !state.HasCastBeenFired {
			combat.UseSkill(c, target, target.GetPosition(), c.NextSkillToUse, nil, nil, svcCtx.NavMesh, svcCtx)
			state.HasCastBeenFired = true
			log.Printf("[SKILL-STATE] [%s (%s)] skill %s disparada", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)

			// Evento padrão
			svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
				SourceHandle: c.Handle,
				BehaviorType: "SkillExecuted",
				Timestamp:    time.Now(),
			})

			if c.NextSkillToUse.Movement != nil && c.SkillMovementState == nil {
				c.SkillMovementState = combat.ApplySkillMovement(c, target, c.NextSkillToUse)
			}
		}

		return core.StatusRunning
	}

	if c.NextSkillToUse.Movement != nil && c.SkillMovementState != nil && c.SkillMovementState.Active {
		return core.StatusRunning
	}

	if time.Now().Before(state.RecoveryUntil) {
		c.SetAnimationState(constslib.AnimationRecovery)
		return core.StatusRunning
	}

	// Ciclo finalizado
	log.Printf("[SKILL-STATE] [%s (%s)] skill %s ciclo finalizado", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)

	svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
		SourceHandle: c.Handle,
		BehaviorType: "SkillCycleCompleted",
		Timestamp:    time.Now(),
	})

	c.ResetSkillState(c.NextSkillToUse.Action)
	c.NextSkillToUse = nil
	return core.StatusSuccess
}

func (n *SkillStateNode) Reset() {
	// Nada a resetar
}
