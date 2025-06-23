package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib/position"
)

type UseGroundSkillNode struct {
	SkillName string
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando UseGroundSkillNode", c.ID)

	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Skill %s não encontrada.", n.SkillName)
		return core.StatusFailure
	}

	// Regra: Não gastar skill se estiver com medo (salvo skills defensivas, se quiser tratar isso depois)
	if c.MentalState == creature.MentalStateAfraid {
		log.Printf("[AI] %s está com medo, recusando-se a usar %s.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	// Verificar se tem players OU criaturas dentro do alcance da skill
	var bestTargetPos position.Position
	var targetsInRange int

	// Verificar players
	for _, p := range ctx.Players {
		dist := CalculateDistance(c.Position, p.Position)
		if dist <= skill.Range {
			targetsInRange++
			bestTargetPos = p.Position
		}
	}

	// Verificar criaturas (exceto ele mesmo e mortos)
	for _, other := range ctx.Creatures {
		if other.ID == c.ID || !other.IsAlive {
			continue
		}
		dist := CalculateDistance(c.Position, other.Position)
		if dist <= skill.Range {
			targetsInRange++
			bestTargetPos = other.Position
		}
	}

	// Se nenhum alvo no range, não usa
	if targetsInRange == 0 {
		log.Printf("[AI] %s não encontrou alvos próximos pra usar %s.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	// Se estiver faminto e for predador, pode usar mesmo em alvo único
	hunger := c.GetNeedValue(creature.NeedHunger)
	if targetsInRange == 1 && !(hunger > 80 && c.HasTag(creature.TagPredator)) {
		log.Printf("[AI] %s preferiu guardar a skill %s para mais inimigos.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	// Executa skill
	log.Printf("[AI] %s usando %s em posição (%f, %f, %f)", c.ID, n.SkillName, bestTargetPos.X, bestTargetPos.Y, bestTargetPos.Z)
	combat.UseSkill(c, nil, bestTargetPos, n.SkillName, ctx.Creatures, ctx.Players)

	return core.StatusSuccess
}