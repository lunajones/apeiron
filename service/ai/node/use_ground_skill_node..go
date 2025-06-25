package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib/position"
)

type UseGroundSkillNode struct {
	SkillName string
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx := ctx.(dynamic_context.AIServiceContext)

	log.Printf("[AI] %s executando UseGroundSkillNode", c.ID)

	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada.", n.SkillName)
		return core.StatusFailure
	}

	if c.MentalState == creature.MentalStateAfraid {
		log.Printf("[AI] %s está com medo, recusando-se a usar %s.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	var bestTargetPos position.Position
	var targetsInRange int

	for _, p := range svcCtx.GetServicePlayers() {
		dist := CalculateDistance(c.Position, p.Position)
		if dist <= skill.Range {
			targetsInRange++
			bestTargetPos = p.Position
		}
	}

	for _, other := range svcCtx.GetServiceCreatures() {
		if other.ID == c.ID || !other.IsAlive {
			continue
		}
		dist := CalculateDistance(c.Position, other.Position)
		if dist <= skill.Range {
			targetsInRange++
			bestTargetPos = other.Position
		}
	}

	if targetsInRange == 0 {
		log.Printf("[AI] %s não encontrou alvos próximos pra usar %s.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	hunger := c.GetNeedValue(creature.NeedHunger)
	if targetsInRange == 1 && !(hunger > 80 && c.HasTag(creature.TagPredator)) {
		log.Printf("[AI] %s preferiu guardar a skill %s para mais inimigos.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	log.Printf("[AI] %s usando %s em posição (%f, %f, %f)", c.ID, n.SkillName, bestTargetPos.X, bestTargetPos.Y, bestTargetPos.Z)
	combat.UseSkill(c, nil, bestTargetPos, n.SkillName, svcCtx.GetServiceCreatures(), svcCtx.GetServicePlayers())

	return core.StatusSuccess
}
