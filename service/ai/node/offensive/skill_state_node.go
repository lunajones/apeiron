package offensive

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/combat"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type SkillStateNode struct{}

func (n *SkillStateNode) Tick(c *creature.Creature, ctx interface{}) interface{} {

	if c.CombatState != constslib.CombatStateCasting {
		return core.StatusSuccess

	}

	drive := c.GetCombatDrive()
	if drive.Rage > 0 && drive.Caution > drive.Rage {
		log.Printf("[SKILL-STATE] [%s] evitou iniciar estado de skill por Caution > Rage (%.2f > %.2f)", c.Handle.String(), drive.Caution, drive.Rage)
		c.CombatState = constslib.CombatStateMoving

		// ⚠️ Sinaliza que hesitou atacar

		return core.StatusSuccess
	}

	if (c.IsBlocking() || c.IsDodging()) && c.CurrentSkillState() != nil && c.CurrentSkillState().CanBeCancelled() {
		log.Printf("[SKILL-STATE] [%s] cancelando skill %s por bloqueio ou esquiva",
			c.Handle.String(), c.NextSkillToUse.Name)
		c.CancelCurrentSkill()
		n.Reset(c)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SKILL-STATE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		n.Reset(c)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil {
		log.Printf("[SKILL-STATE] [%s] não possui skill para executar", c.PrimaryType)
		c.CombatState = constslib.CombatStateMoving
		return core.StatusRunning
	}

	log.Printf("[SKILL-STATE] [%s] iniciando estado de skill: %s", c.PrimaryType, c.NextSkillToUse.Name)

	state := c.SkillStates[c.NextSkillToUse.Action]

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[CIRCLE] [%s (%s)] alvo não encontrado", c.Handle.String(), c.PrimaryType)
		return core.StatusRunning
	}

	dist := position.CalculateDistance(c.Position, target.GetPosition())
	if dist > c.NextSkillToUse.Range {
		log.Printf("[SKILL-STATE] [%s (%s)] cancelada skill %s: fora do alcance (%.2f > %.2f)",
			c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, dist, c.NextSkillToUse.Range)
		// n.Reset(c)
		c.CombatState = constslib.CombatStateMoving
		return core.StatusSuccess
	}
	if state == nil || !state.InUse {
		if target == nil {
			log.Printf("[SKILL-STATE] [%s (%s)] alvo inválido", c.Handle.String(), c.PrimaryType)
			n.Reset(c)
			return core.StatusFailure
		}

		if c.NextSkillToUse.Movement == nil {
			dist := position.CalculateDistance(c.Position, target.GetPosition())
			if dist > c.NextSkillToUse.Range {
				log.Printf("[SKILL-STATE] [%s (%s)] cancelada skill %s: fora do alcance (%.2f > %.2f)",
					c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, dist, c.NextSkillToUse.Range)
				n.Reset(c)
				return core.StatusFailure
			}
		}

		state = c.InitSkillState(c.NextSkillToUse.Action, time.Now())
		log.Printf("[SKILL-STATE] [%s (%s)] iniciando skill %s", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)
		return core.StatusRunning
	}

	if state.WasInterrupted {
		log.Printf("[SKILL-STATE] [%s (%s)] skill %s foi interrompida", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)
		n.Reset(c)
		return core.StatusSuccess
	}

	now := time.Now()

	// WINDUP
	if !state.WindUpFired {
		state.WindUpUntil = now.Add(time.Duration(c.NextSkillToUse.WindUpTime * float64(time.Second)))
		state.WindUpFired = true
		c.SetAnimationState(constslib.AnimationWindup)

		log.Printf("[SKILL-STATE] [%s (%s)] skill %s entrando em WindUp até %v", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, state.WindUpUntil)
		return core.StatusRunning
	}
	if now.Before(state.WindUpUntil) {
		return core.StatusRunning
	}

	// CAST
	if !state.CastFired {

		if targetCreature, ok := target.(*creature.Creature); ok {
			if targetCreature.IsBlocking() && !c.NextSkillToUse.CanCastWhileBlocking {
				log.Printf("[SKILL-STATE] [%s] evitou cast pois alvo está bloqueando", c.Handle.String())
				return core.StatusFailure
			}

			for _, evt := range targetCreature.GetCombatEvents() {
				if evt.BehaviorType == "ParryFailed" || evt.BehaviorType == "BlockBroken" {
					if time.Since(evt.Timestamp) < 1*time.Second {
						log.Printf("[SKILL-STATE] [%s] antecipou cast por erro de defesa do alvo", c.Handle.String())
						break
					}
				}
			}
		}

		allies := finder.FindNearbyAllies(svcCtx, c, c.GetFaction(), 8.0)
		for _, ally := range allies {
			if ally != nil && ally.GetHandle() != c.GetHandle() && ally.GetCombatDrive().Termination > 0.6 {
				log.Printf("[SKILL-STATE] [%s] detectou aliado [%s] com Termination alto → sincronizando cast", c.Handle.String(), ally.GetHandle().String())
				break
			}
		}

		state.CastUntil = now.Add(time.Duration(c.NextSkillToUse.CastTime * float64(time.Second)))
		state.CastFired = true

		c.SetAnimationState(constslib.AnimationCast)

		log.Printf("[SKILL-STATE] [%s (%s)] skill %s entrando em Cast até %v", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, state.CastUntil)

		if targetCreature, ok := target.(*creature.Creature); ok {
			event := model.CombatEvent{
				SourceHandle:   c.Handle,
				BehaviorType:   "AggressiveIntention",
				Timestamp:      now,
				ExpectedImpact: state.CastUntil,
			}
			targetCreature.RegisterCombatEvent(event)
			log.Printf("%s", color.New(color.FgHiGreen).Sprintf("[AGGRO] Registrado: AggressiveIntention [%s → %s]", c.Handle.String(), targetCreature.Handle.String()))
		}

		return core.StatusRunning
	}
	if now.Before(state.CastUntil) {
		return core.StatusRunning
	}

	// RECOVERY
	if !state.RecoveryFired {
		combat.UseSkill(c, target, target.GetPosition(), c.NextSkillToUse, nil, nil, svcCtx.NavMesh, svcCtx)

		if c.NextSkillToUse.Movement != nil && c.SkillMovementState == nil {
			c.SkillMovementState = combat.ApplySkillMovement(c, target, c.NextSkillToUse)
		}

		state.CooldownUntil = now.Add(time.Duration(c.NextSkillToUse.CooldownSec * float64(time.Second)))
		state.RecoveryUntil = now.Add(time.Duration(c.NextSkillToUse.RecoveryTime * float64(time.Second)))
		state.RecoveryFired = true

		c.SetAnimationState(constslib.AnimationRecovery)
		log.Printf("[SKILL-STATE] [%s (%s)] skill %s disparada e entrando em Recovery até %v", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name, state.RecoveryUntil)

		return core.StatusRunning
	}
	if now.Before(state.RecoveryUntil) {
		return core.StatusRunning
	}
	c.CombatState = constslib.CombatStateMoving
	log.Printf("[SKILL-STATE] [%s (%s)] skill %s ciclo finalizado", c.Handle.String(), c.PrimaryType, c.NextSkillToUse.Name)
	n.Reset(c)
	return core.StatusSuccess
}

func (n *SkillStateNode) Reset(c *creature.Creature) {
	log.Printf("[SKILL-STATE] [RESET] [%s (%s)]", c.Handle.String(), c.PrimaryType)

	if c.NextSkillToUse != nil {
		state := c.SkillStates[c.NextSkillToUse.Action]
		if state != nil {
			state.InUse = false
			state.WindUpFired = false
			state.CastFired = false
			state.RecoveryFired = false
		}
		c.NextSkillToUse = nil
	}
	if c.SkillMovementState != nil && c.SkillMovementState.Active {
		c.SkillMovementState = nil
	}
}
