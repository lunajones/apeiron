package ai

import (
	"log"
	"math/rand"

	"github.com/lunajones/apeiron/service/combat"
	"github.com/lunajones/apeiron/service/creature"
)

type UseGroundSkillNode struct {
	SkillName string
	Players   []Player
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature) BehaviorStatus {
	if len(n.Players) == 0 {
		log.Printf("[AI] Nenhum player disponível para %s mirar o skill %s", c.ID, n.SkillName)
		return StatusFailure
	}

	// Exemplo simples: escolher um player aleatório como referência
	targetPlayer := n.Players[rand.Intn(len(n.Players))]
	targetPos := targetPlayer.Position

	// (Opcional futuro) - Você pode randomizar a posição ao redor do player aqui, se quiser:
	// targetPos.X += rand.Float64()*2 - 1
	// targetPos.Z += rand.Float64()*2 - 1

	combat.UseSkill(c, nil, targetPos, n.SkillName, nil, n.Players)

	log.Printf("[AI] Creature %s usou skill ground-targeted %s na posição %.2f, %.2f, %.2f",
		c.ID, n.SkillName, targetPos.X, targetPos.Y, targetPos.Z)

	return StatusSuccess
}
