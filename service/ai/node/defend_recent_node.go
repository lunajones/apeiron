package node

import (
	"log"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type DefendRecentNode struct{}

func (n *DefendRecentNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI] [%s (%s)] contexto inválido em DefendRecentNode", c.Handle.ID, c.PrimaryType)
		return core.StatusFailure
	}

	if len(c.AggroTable) == 0 {
		return core.StatusFailure
	}

	var (
		bestTarget    *creature.Creature
		highestThreat float64
	)

	nearby := svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius)

	for _, entry := range c.AggroTable {
		target := svcCtx.FindCreatureByHandle(entry.TargetHandle)
		if target == nil || !target.IsAlive {
			continue
		}

		// Verifica se está realmente por perto
		foundNearby := false
		for _, n := range nearby {
			if n.GetHandle().Equals(target.GetHandle()) {
				foundNearby = true
				break
			}
		}
		if !foundNearby {
			continue
		}

		if entry.ThreatValue > highestThreat {
			bestTarget = target
			highestThreat = entry.ThreatValue
		}
	}

	if bestTarget == nil {
		c.TargetCreatureHandle = handle.EntityHandle{}
		return core.StatusFailure
	}

	c.TargetCreatureHandle = bestTarget.GetHandle()
	c.ChangeAIState(consts.AIStateCombat)

	log.Printf("[AI] [%s (%s)] defendendo-se de [%s (%s)] (ameaça: %.2f)",
		c.Handle.ID, c.PrimaryType,
		bestTarget.Handle.ID, bestTarget.PrimaryType,
		highestThreat,
	)

	return core.StatusSuccess
}

func (n *DefendRecentNode) Reset() {
	// Esse node não tem estado interno, então o Reset não faz nada
}
