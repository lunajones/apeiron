package offensive

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper"
)

type PlanOffensiveSkillNode struct{}

func (n *PlanOffensiveSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	now := time.Now()

	// Verifica se é hora de planejar novamente
	// if now.Before(c.NextAggressiveDecisionAllowed) {
	// 	log.Printf("[PLAN-SKILL] [%s (%s)] aguardando próxima decisão agressiva até %s",
	// 		c.Handle.String(), c.PrimaryType, c.NextAggressiveDecisionAllowed.Format("15:04:05.000"))
	// 	return core.StatusRunning
	// }

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[PLAN-SKILL] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	best := helper.FindBestOffensiveSkill(c, svcCtx, now)
	if best != nil {
		c.NextSkillToUse = best

		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "OffensiveSkillPlanned",
			Timestamp:    now,
		})

		log.Printf("%s", color.New(color.FgRed).Sprintf(
			"[PLAN-SKILL] [%s (%s)] planejou skill ofensiva: %s",
			c.Handle.String(), c.PrimaryType, best.Name,
		))

		// Atualiza próximo momento permitido para planejar
		// c.NextAggressiveDecisionAllowed = now.Add(500 * time.Millisecond) // ajuste o valor conforme o ritmo desejado
		return core.StatusSuccess
	}

	basic := helper.FindBasicAttack(c, now)
	if basic != nil {
		c.NextSkillToUse = basic

		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "OffensiveSkillPlanned",
			Timestamp:    now,
		})

		log.Printf("%s", color.New(color.FgYellow).Sprintf(
			"[PLAN-SKILL] [%s (%s)] fallback para ataque básico: %s",
			c.Handle.String(), c.PrimaryType, basic.Name,
		))

		// c.NextAggressiveDecisionAllowed = now.Add(500 * time.Millisecond) // mesmo delay no fallback
		return core.StatusSuccess
	}

	log.Printf("%s", color.New(color.FgHiBlack).Sprintf(
		"[PLAN-SKILL] [%s (%s)] nenhuma skill disponível para planejar",
		c.Handle.String(), c.PrimaryType,
	))
	c.NextSkillToUse = nil

	// Mesmo sem skill, aplica um pequeno delay para evitar spam
	// c.NextAggressiveDecisionAllowed = now.Add(300 * time.Millisecond)
	return core.StatusRunning
}

func (n *PlanOffensiveSkillNode) Reset() {}
