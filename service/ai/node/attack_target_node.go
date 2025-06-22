package node

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

type AttackTargetNode struct {
	SkillName string
}

func (n *AttackTargetNode) Tick(c *creature.Creature) core.BehaviorStatus {
	if c.TargetCreatureID == "" {
		return core.StatusFailure
	}

	target := world.FindCreatureByID(c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		return core.StatusFailure
	}

	distance := distance(c.Position, target.Position)
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada.", n.SkillName)
		return core.StatusFailure
	}

	if distance > skill.Range {
		log.Printf("[AI] Target %s fora de alcance de %s.", target.ID, n.SkillName)
		return core.StatusFailure
	}

	combat.UseSkill(c, target, target.Position, n.SkillName, nil, nil)
	return core.StatusSuccess
}

func distance(a, b position.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
