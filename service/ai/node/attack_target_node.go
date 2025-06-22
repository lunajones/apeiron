package node

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/service/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

type AttackTargetNode struct {
	SkillName string
}

func (n *AttackTargetNode) Tick(c *creature.Creature) BehaviorStatus {
	if c.TargetCreatureID == "" {
		return StatusFailure
	}

	target := world.FindCreatureByID(c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		return StatusFailure
	}

	distance := Distance(c.Position, target.Position)
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada no registry.", n.SkillName)
		return StatusFailure
	}

	if distance > skill.Range {
		log.Printf("[AI] Target %s fora do alcance de %s.", target.ID, n.SkillName)
		return StatusFailure
	}

	combat.UseSkill(c, target, target.Position, n.SkillName, nil, nil)
	log.Printf("[AI] Creature %s usou %s contra %s", c.ID, n.SkillName, target.ID)

	return StatusSuccess
}

func Distance(a, b creature.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
