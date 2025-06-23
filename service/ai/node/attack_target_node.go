package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
)

type AttackTargetNode struct {
	SkillName string
}

func (n *AttackTargetNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando AttackTargetNode", c.ID)

	if c.TargetCreatureID == "" {
		log.Printf("[AI] %s não tem alvo para atacar.", c.ID)
		return core.StatusFailure
	}

	target := creature.FindByID(ctx.Creatures, c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		log.Printf("[AI] %s: Target inválido ou morto.", c.ID)
		return core.StatusFailure
	}

	// Regra 1: Estado mental
	if c.MentalState == creature.MentalStateAfraid {
		log.Printf("[AI] %s está com medo, recusando-se a atacar.", c.ID)
		return core.StatusFailure
	}

	// Regra 2: Analisar fome extrema (Exemplo: se a fome está acima de 90, ignora o medo)
	hunger := c.GetNeedValue(creature.NeedHunger)
	if hunger > 90 {
		log.Printf("[AI] %s está faminto demais, vai atacar de qualquer jeito!", c.ID)
	} else {
	// Regra 3: Se for animal, só atacar se o alvo for "prey"
	if c.HasTag(creature.TagAnimal) && !target.HasTag(creature.TagPrey) {
		log.Printf("[AI] %s é um animal e o alvo %s não é presa. Abortando ataque.", c.ID, target.ID)
		return core.StatusFailure
	}
}

	distance := calculateDistance(c.Position, target.Position)
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada para %s.", n.SkillName, c.ID)
		return core.StatusFailure
	}

	if distance > skill.Range {
		log.Printf("[AI] %s: Alvo %s fora de alcance da skill %s.", c.ID, target.ID, n.SkillName)
		return core.StatusFailure
	}

	combat.UseSkill(c, target, target.Position, n.SkillName, ctx.Creatures, ctx.Players)
	log.Printf("[AI] %s atacou %s com %s.", c.ID, target.ID, n.SkillName)
	return core.StatusSuccess
}
