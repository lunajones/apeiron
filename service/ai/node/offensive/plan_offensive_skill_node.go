package offensive

import (
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type PlanOffensiveSkillNode struct{}

func (n *PlanOffensiveSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[PLAN-OFFENSIVE] [%s] contexto inválido", c.GetPrimaryType())
		return core.StatusFailure
	}

	if c.NextSkillToUse != nil {
		color.Red("[PLAN-OFFENSIVE] [%s] Já tem skill planejada", c.GetPrimaryType())
		return core.StatusSuccess
	}

	drive := c.GetCombatDrive()
	stamina := c.Stamina
	hpRatio := float64(c.HP) / float64(c.MaxHP)
	recentMiss := time.Since(c.GetLastSkillMissedAt()) < 3*time.Second

	// 🔍 Obtém o alvo atual (criatura ou player)
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	blocking := false
	if target != nil && target.IsBlocking() {
		blocking = true
	}

	// 💡 Hesita por erro recente + cautela/counter
	if recentMiss && (drive.Caution > drive.Rage || drive.Counter < 0.3) {
		log.Printf("[PLAN-OFFENSIVE] [%s] hesitando por erro recente + cautela/counter", c.GetPrimaryType())
		return core.StatusFailure
	}

	// 💡 Hesita por alvo bloqueando + cautela
	if blocking && stamina > 10.0 && drive.Caution > drive.Rage {
		log.Printf("[PLAN-OFFENSIVE] [%s] hesitando por bloqueio + cautela alta (stamina %.1f)", c.GetPrimaryType(), stamina)
		return core.StatusFailure
	}

	// 💡 HP alto + counter baixo → chance de hesitar
	if hpRatio > 0.7 && drive.Counter < 0.2 {
		if rand.Float64() < 0.4 {
			log.Printf("[PLAN-OFFENSIVE] [%s] hesitando por HP alto (%.2f) + counter baixo (%.2f)", c.GetPrimaryType(), hpRatio, drive.Counter)
			return core.StatusFailure
		}
	}

	// 📌 Marca tentativa de planejar
	c.AddRecentAction(constslib.CombatActionAttackPrepared)

	// 🔍 Busca melhor skill
	bestSkill := helper.FindBestOffensiveSkill(c, svcCtx, time.Now())
	if bestSkill == nil {
		log.Printf("[PLAN-OFFENSIVE] [%s] nenhuma skill disponível", c.GetPrimaryType())
		return core.StatusFailure
	}

	c.NextSkillToUse = bestSkill
	c.LastSkillPlannedAt = time.Now()

	// TODO: migrar controle de CombatState para CombatFSM no futuro
	c.SetCombatState(constslib.CombatStateMoving)

	color.Red("[PLAN-OFFENSIVE] [%s] Planejada skill: %s", c.GetPrimaryType(), bestSkill.Name)
	return core.StatusSuccess
}

func (n *PlanOffensiveSkillNode) Reset(c *creature.Creature) {
	log.Printf("[PLAN-SKILL] [RESET] [%s]", c.GetPrimaryType())
}
