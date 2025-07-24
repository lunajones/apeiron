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
	// NÃO planeja se estiver bloqueando ou esquivando
	if c.IsBlocking() || c.IsDodging() {
		return nil
	}

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

	// Base boost para ataque básico
	if skill.Action == constslib.Basic {
		score += 1.0
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())

	// 💡 Análise tática baseada em distância e tipo de skill
	if dist <= skill.Range {
		score += 3.0 // Dentro do alcance? ótimo
	} else {
		outOfRange := dist - skill.Range

		if skill.Tags.Has(constslib.SkillTagRush) {
			// ⚡ Rush: ganha incentivo se estiver fora do alcance, especialmente entre 2–5m
			if outOfRange <= 3.5 {
				score += 2.5 - (outOfRange * 0.5) // Diminui com a distância, mas ainda vale a pena
			} else {
				score -= outOfRange // longe demais, penaliza normal
			}
		} else {
			// ❌ Skills normais fora do range: penaliza pesado
			if skill.Action == constslib.Basic {
				score -= outOfRange * 0.3
			} else {
				score -= outOfRange
			}
		}
	}

	// DRIVE INFLUÊNCIA TÁTICA
	drive := c.GetCombatDrive()

	if target.IsCasting() && skill.Tags.Has(constslib.SkillTagInterrupt) {
		score += 3.0
	}

	if drive.Termination > 0.6 {
		if skill.Tags.Has(constslib.SkillTagBurst) {
			score += 2.0
		}
		if skill.Tags.Has(constslib.SkillTagRush) {
			score += 1.5
		}
	}

	if drive.Caution > 0.5 && skill.Range >= 4.0 {
		score += 2.0
	}

	if drive.Rage > 0.5 && skill.Tags.Has(constslib.SkillTagRush) {
		score += 2.5
	}

	// CONTROLE / FINALIZAÇÃO
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

	// 🔀 Pequena variação aleatória pra não parecer robótico
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

func FindBasicAttack(c *creature.Creature) *model.Skill {
	for _, skill := range c.RegisteredSkills {
		if skill != nil && skill.Action == constslib.Basic {
			now := time.Now()
			state := c.SkillStates[skill.Action]
			log.Printf("[DEBUG] [%s] encontrou ataque básico: %s (cooldown até %.2fs)",
				c.Handle.String(),
				skill.Name,
				state.CooldownUntil.Sub(now).Seconds(),
			)
			if state != nil && !now.Before(state.CooldownUntil) {
				return skill
			}
		}
	}
	return nil
}
