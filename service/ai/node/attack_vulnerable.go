package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
)

type AttackIfEnemyVulnerableNode struct {
	SkillName string
}

func (n *AttackIfEnemyVulnerableNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	if c.TargetCreatureID == "" {
		return core.StatusFailure
	}

	target := creature.FindByID(ctx.Creatures, c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		return core.StatusFailure
	}

	if target.HP > 30 {
		log.Printf("[AI] Target %s com HP alto demais para vulnerável.", target.ID)
		return core.StatusFailure
	}

	distance := calculateDistance(c.Position, target.Position)
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada.", n.SkillName)
		return core.StatusFailure
	}

	if distance > skill.Range {
		log.Printf("[AI] Target %s fora de alcance de %s.", target.ID, n.SkillName)
		return core.StatusFailure
	}

	combat.UseSkill(c, target, target.Position, n.SkillName, ctx.Creatures, ctx.Players)
	log.Printf("[AI] %s executou ataque vulnerável em %s com %s", c.ID, target.ID, n.SkillName)
	return core.StatusSuccess
}