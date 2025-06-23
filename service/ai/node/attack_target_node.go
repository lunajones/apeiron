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
		return core.StatusFailure
	}

	target := creature.FindByID(ctx.Creatures, c.TargetCreatureID)
	if target == nil || !target.IsAlive {
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
	return core.StatusSuccess
}