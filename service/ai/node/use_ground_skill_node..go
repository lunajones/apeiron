package node

import (
	"log"
	"math/rand"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type UseGroundSkillNode struct {
	SkillName string
	Players   []player.Player
}

func (n *UseGroundSkillNode) Tick(c *creature.Creature) core.BehaviorStatus {
	skill, exists := combat.SkillRegistry[n.SkillName]
	if !exists {
		log.Printf("[AI] Ground skill %s não encontrada.", n.SkillName)
		return core.StatusFailure
	}

	if len(n.Players) == 0 {
		log.Printf("[AI] Nenhum player disponível como target de skill ground.")
		return core.StatusFailure
	}

	targetPlayer := n.Players[rand.Intn(len(n.Players))]
	targetPos := targetPlayer.Position

	combat.UseSkill(c, nil, targetPos, n.SkillName, nil, nil)
	log.Printf("[AI] Creature %s usou %s em posição (%f, %f, %f)", c.ID, n.SkillName, targetPos.X, targetPos.Y, targetPos.Z)

	return core.StatusSuccess
}
