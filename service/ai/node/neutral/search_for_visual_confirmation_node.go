package neutral

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/ai/sensor"
	"github.com/lunajones/apeiron/service/creature"
)

type SearchForVisualConfirmationNode struct{}

func (n *SearchForVisualConfirmationNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[VISUAL-CONFIRM] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	target := svcCtx.FindByHandle(c.TargetCreatureHandle)
	if target == nil {
		return core.StatusFailure
	}

	creatureTarget, ok := target.(*creature.Creature)
	if !ok {
		log.Printf("[VISUAL-CONFIRM] [%s (%s)] alvo não é uma criatura válida", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if sensor.CanSee(c, creatureTarget) {
		c.CombatState = consts.CombatStateAggressive
		log.Printf("%s", color.New(color.FgHiRed).Sprintf(
			"[VISUAL-CONFIRM] [%s (%s)] visão confirmada → CombatState=Aggressive",
			c.Handle.String(), c.PrimaryType))
		return core.StatusSuccess
	}

	// Se passou mais de 5 segundos desde a detecção → voltar para strategic
	if time.Since(c.LastThreatSeen) > 5*time.Second {
		c.CombatState = consts.CombatStateStrategic
		log.Printf("%s", color.New(color.FgHiCyan).Sprintf(
			"[VISUAL-CONFIRM] [%s (%s)] ameaça não confirmada, retornando para Strategic",
			c.Handle.String(), c.PrimaryType))
		return core.StatusSuccess
	}

	return core.StatusRunning
}

func (n *SearchForVisualConfirmationNode) Reset() {}
