package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type ProcessAIStateNode struct{}

func (n *ProcessAIStateNode) Tick(c *creature.Creature, ctx interface{}) core.BehaviorStatus {
	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI] Contexto inválido para ProcessAIStateNode")
		return core.StatusFailure
	}

	switch c.AIState {
	case creature.AIStateIdle:
		if rand.Float32() < 0.1 {
			c.ChangeAIState(creature.AIStateAlert)
		}

	case creature.AIStateAlert:
		for _, p := range svcCtx.GetServicePlayers() {
			if creature.CanSeePlayer(c, []*player.Player{p}) || creature.CanHearPlayer(c, []*player.Player{p}) {
				c.AddThreat(p.ID, 10, "PlayerDetected", "VisionOrSound")
				log.Printf("[AI] %s detectou o player %s e adicionou threat.", c.ID, p.ID)
				c.ChangeAIState(creature.AIStateChasing)
				break
			}
		}
		if time.Since(c.LastStateChange) > 2*time.Second {
			c.ChangeAIState(creature.AIStateIdle)
		}

	case creature.AIStateChasing:
		targetID := c.GetHighestThreatTarget()
		if targetID == "" {
			log.Printf("[AI] %s sem alvo de threat, voltando pra Idle", c.ID)
			c.ChangeAIState(creature.AIStateIdle)
			return core.StatusSuccess
		}
		target := creature.FindTargetByID(targetID, svcCtx.GetServiceCreatures(), svcCtx.GetServicePlayers())
		if target == nil {
			log.Printf("[AI] %s: alvo %s não encontrado, limpando aggro", c.ID, targetID)
			c.ClearAggro()
			c.ChangeAIState(creature.AIStateIdle)
			return core.StatusSuccess
		}
		c.MoveTowards(target.GetPosition(), c.MoveSpeed)

	case creature.AIStateAttack:
		log.Printf("[Creature %s] Atacando!", c.ID)
		c.SetAction(creature.ActionAttack)
		c.ChangeAIState(creature.AIStateIdle)

	case creature.AIStateDead:
		// Nada a fazer
	}

	return core.StatusSuccess
}
