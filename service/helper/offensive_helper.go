// combo_planner.go
package helper

import (
	"log"
	"math/rand/v2"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

func FindBestOffensiveSkill(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext, now time.Time) *model.Skill {
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		return nil
	}

	var bestSkill *model.Skill
	bestScore := -9999.0

	for _, skill := range c.RegisteredSkills {
		if skill == nil {
			continue
		}
		state := c.SkillStates[skill.Action]
		if state == nil || now.Before(state.CooldownUntil) {
			continue
		}
		if skill.Conditions != nil && !ValidateSkillConditions(c, skill) {
			continue
		}
		score := CalculateSkillScore(c, target, skill)
		if score > bestScore {
			bestScore = score
			bestSkill = skill
		}
		log.Printf("[PLAN-BEST-SKILL] %s score: %.2f", skill.Name, score)
	}

	return bestSkill
}

func CalculateSkillScore(c *creature.Creature, target model.Targetable, skill *model.Skill) float64 {
	score := skill.ScoreBase

	// Dê um pequeno boost base pro ataque básico (caso queira garantir que figure em combos)
	if skill.Action == constslib.Basic {
		score += 1.0
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	if dist <= skill.Range {
		score += 3.0
	} else {
		// Penalização ajustada para o ataque básico parecer mais natural (estilo Deathmarch TW3)
		if skill.Action == constslib.Basic {
			score -= (dist - skill.Range) * 0.3
		} else {
			score -= (dist - skill.Range)
		}
	}

	if tgtCreature, ok := target.(*creature.Creature); ok {
		if tgtCreature.Posture < 20 && skill.Impact != nil && skill.Impact.PostureDamage > 0 {
			score += 2.0
		}
		if skill.HasDOT && !HasDOTEffectOfType(tgtCreature, skill.DOT.EffectType) {
			score += 1.5
		}
		if skill.HasDOT && HasDOTEffectOfType(tgtCreature, skill.DOT.EffectType) {
			score -= 1.0
		}
		if tgtCreature.HP < 10 {
			score += 2.0 / (skill.WindUpTime + skill.CastTime + 0.1)
		}
	}

	score += rand.Float64() * 0.5
	return score
}

func ValidateSkillConditions(c *creature.Creature, skill *model.Skill) bool {
	cond := skill.Conditions
	if cond == nil {
		return true
	}
	if cond.FacingRequirement != "" && cond.FacingRequirement == "Behind" {
		// Aqui você pode expandir a lógica real de verificação de facing
		return false
	}
	// Outros requisitos podem ser adicionados aqui
	return true
}

func FindBasicAttack(c *creature.Creature, now time.Time) *model.Skill {
	for _, skill := range c.RegisteredSkills {
		if skill != nil && skill.Action == constslib.Basic {

			state := c.SkillStates[skill.Action]
			log.Printf("[DEBUG] [%s] encontrou ataque básico: %s (cooldown até %.2fs)",
				c.Handle.String(),
				skill.Name,
				state.CooldownUntil.Sub(time.Now()).Seconds(),
			)
			if state != nil && !now.Before(state.CooldownUntil) {
				return skill
			}
		}
	}
	return nil
}
