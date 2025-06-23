package node

import (
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/lib/position"
)

type UseGroundSkillNode struct {
	SkillName string
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	if len(ctx.Players) == 0 {
		return core.StatusFailure
	}
	targetPlayer := ctx.Players[0] // Exemplo: sempre escolhe o primeiro jogador da lista
	distanceToTarget := distance(c.Position, targetPlayer.Position)
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

func distance(a, b position.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}