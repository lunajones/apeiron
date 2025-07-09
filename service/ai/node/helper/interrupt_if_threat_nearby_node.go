package helper

import (
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/ai/sensor"
	"github.com/lunajones/apeiron/service/creature"
)

type InterruptIfThreatNearbyNode struct {
	InterruptAIState   constslib.AIState
	InterruptAnimation constslib.AnimationState
}

func (n *InterruptIfThreatNearbyNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-INTERRUPT] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	for _, t := range svcCtx.GetCachedTargets(c.Handle) {
		other, ok := t.(*creature.Creature)
		if !ok || !other.Alive || other.Handle.Equals(c.Handle) {
			continue
		}

		if !creature.AreEnemies(c, other) {
			continue
		}

		var detectedBy string

		if sensor.CanSee(c, other) {
			detectedBy = "vision"
		} else if sensor.CanHear(c, other) {
			detectedBy = "hearing"
		} else if sensor.CanSmell(c, other) {
			detectedBy = "smell"
		} else {
			continue
		}

		// Stealth chance se agachado (não na visão direta)
		if other.IsCurrentlyCrouched() && detectedBy != "vision" {
			if rand.Float64() >= 0.25 {
				log.Printf("[AI-INTERRUPT] [%s (%s)] ameaça agachada próxima, stealth bem-sucedido.",
					c.Handle.String(), c.PrimaryType)
				continue
			}
			log.Printf("[AI-INTERRUPT] [%s (%s)] ameaça agachada próxima, stealth falhou. Interrompendo.",
				c.Handle.String(), c.PrimaryType)
		} else {
			log.Printf("[AI-INTERRUPT] [%s (%s)] ameaça detectada (%s). Interrompendo.",
				c.Handle.String(), c.PrimaryType, detectedBy)
		}

		// Decide CombatState inicial
		switch detectedBy {
		case "vision":
			c.CombatState = constslib.CombatStateAggressive
		case "hearing", "smell":
			c.CombatState = constslib.CombatStateCautious
		default:
			c.CombatState = constslib.CombatStateIdle
		}
		log.Printf("[AI-INTERRUPT] [%s] combate iniciado em %s", c.Handle.String(), detectedBy)
		log.Printf("%s", color.New(color.FgHiMagenta).Sprintf(
			"[AI-INTERRUPT] [%s (%s)] combate iniciado → CombatState=%s (detectedBy=%s)",
			c.Handle.String(), c.PrimaryType,
			c.CombatState.String(), detectedBy,
		))

		c.TargetCreatureHandle = other.Handle
		c.LastThreatSeen = time.Now()

		if n.InterruptAnimation != "" {
			c.SetAnimationState(n.InterruptAnimation)
		}
		if n.InterruptAIState != "" {
			c.ChangeAIState(n.InterruptAIState)
		}

		// Registra ThreatDetected
		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "ThreatDetected",
			Timestamp:    time.Now(),
		})

		return core.StatusSuccess
	}

	return core.StatusFailure
}

func (n *InterruptIfThreatNearbyNode) Reset() {}
