package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

type AttackIfVulnerableNode struct{}

func (n *AttackIfVulnerableNode) Tick(c *creature.Creature) core.BehaviorStatus {
	if c.TargetCreatureID == "" {
		return core.StatusFailure
	}

	target := world.FindCreatureByID(c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		return core.StatusFailure
	}

	if target.IsPostureBroken {
		c.SetAction(creature.ActionSkill3) // Exemplo: skill mais pesada
		log.Printf("[AI] Creature %s atacando alvo vulnerável %s", c.ID, target.ID)
		return core.StatusSuccess
	}

	return core.StatusFailure
}
