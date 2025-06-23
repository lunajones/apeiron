package node

import (
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib/combat"
)

type UseGroundSkillNode struct {
	SkillName string
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	if len(ctx.Players) == 0 {
		return core.StatusFailure
	}
	targetPlayer := ctx.Players[0] // Exemplo: sempre escolhe o primeiro jogador da lista
	distanceToTarget := calculateDistance(c.Position, targetPlayer.Position)
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		return core.StatusFailure
	}
	if distanceToTarget > skill.Range {
		return core.StatusFailure
	}
	combat.UseSkill(c, nil, targetPlayer.Position, n.SkillName, ctx.Creatures, ctx.Players)
	return core.StatusSuccess
}