package offensive

import (
	"log"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper"
)

type PlanOffensiveSkillNode struct{}

func (n *PlanOffensiveSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[PLAN-OFFENSIVE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	// Já tem skill pronta
	if c.NextSkillToUse != nil {
		color.Red("[PLAN-OFFENSIVE] Já tem skill planejada")

		return core.StatusSuccess
	}

	// color.Red("[PLAN-OFFENSIVE] [%s] acionou planejamento: %s", c.GetPrimaryType())

	c.RecentActions = append(c.RecentActions, constslib.CombatActionAttackPrepared)

	// Busca melhor skill usando helper (já validando bloqueio, dodge, distância, cast, drive, etc)
	bestSkill := helper.FindBestOffensiveSkill(c, svcCtx, time.Now())
	if bestSkill == nil {
		return core.StatusFailure
	}

	c.NextSkillToUse = bestSkill
	c.LastSkillPlannedAt = time.Now()

	c.CombatState = constslib.CombatStateMoving

	color.Red("[PLAN-OFFENSIVE] [%s] Planejada skill: %s", c.Handle.String(), bestSkill.Name)
	return core.StatusSuccess
}

func (n *PlanOffensiveSkillNode) Reset(c *creature.Creature) {
	log.Printf("[PLAN-SKILL] [RESET] [%s (%s)]", c.Handle.String(), c.PrimaryType)
}
