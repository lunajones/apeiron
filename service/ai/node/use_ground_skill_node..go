package node

import (
	"log"
	"math"

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

	// Verificar se tem players dentro do alcance da skill
	var bestTargetPos position.Position
	var playersInRange int

	for _, p := range ctx.Players {
		dist := distance(c.Position, p.Position)
		if dist <= skill.Range {
			playersInRange++
			bestTargetPos = p.Position // Por simplicidade, usa o primeiro encontrado
		}
	}

	// Se nenhum player no range, não usa
	if playersInRange == 0 {
		log.Printf("[AI] %s não encontrou players próximos pra usar %s.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	// Se estiver faminto e for predador, pode usar mesmo em alvo único
	hunger := c.GetNeedValue(creature.NeedHunger)
	if playersInRange == 1 && !(hunger > 80 && c.HasTag(creature.TagPredator)) {
		log.Printf("[AI] %s preferiu guardar a skill %s para mais inimigos.", c.ID, n.SkillName)
		return core.StatusFailure
	}

	// Executa skill
	log.Printf("[AI] %s usando %s em posição (%f, %f, %f)", c.ID, n.SkillName, bestTargetPos.X, bestTargetPos.Y, bestTargetPos.Z)
	combat.UseSkill(c, nil, bestTargetPos, n.SkillName, ctx.Creatures, ctx.Players)

	return core.StatusSuccess
}

func distance(a, b position.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
