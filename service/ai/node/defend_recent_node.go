package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type DefendRecentNode struct{}

func (n *DefendRecentNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI] Contexto inválido em DefendRecentNode para %s", c.ID)
		return core.StatusFailure
	}

	log.Printf("[AI] %s executando DefendRecentNode", c.ID)

	if len(c.AggroTable) == 0 {
		log.Printf("[AI] %s não possui aggro registrado", c.ID)
		return core.StatusFailure
	}

	var (
		bestTarget *creature.Creature
		highestAggro float64
	)

	for targetID, entry := range c.AggroTable {
		target := svcCtx.FindCreatureByID(targetID)
		if target == nil || !target.IsAlive {
			continue
		}

		if entry.Aggro > highestAggro {
			bestTarget = target
			highestAggro = entry.Aggro
		}
	}

	if bestTarget == nil {
		log.Printf("[AI] %s não encontrou inimigos válidos para se defender", c.ID)
		return core.StatusFailure
	}

	c.TargetCreatureID = bestTarget.ID
	c.ChangeAIState(creature.AIStateCombat)
	log.Printf("[AI] %s decidiu se defender contra %s com aggro %.2f", c.ID, bestTarget.ID, highestAggro)
	return core.StatusSuccess
}
