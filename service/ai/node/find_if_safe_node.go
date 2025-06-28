package node

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type FindIfSafeNode struct {
	SafeDistance float64
}

func (n *FindIfSafeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SAFE CHECK] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	creatures := svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius*2)
	cX := c.Position.FastGlobalX()
	cZ := c.Position.FastGlobalZ()

	isSafe := true

	for _, other := range creatures {
		if other.Handle.ID == c.Handle.ID || !other.IsAlive {
			continue
		}
		if !creature.AreEnemies(c, other) {
			continue
		}

		dist := math.Hypot(other.Position.FastGlobalX()-cX, other.Position.FastGlobalZ()-cZ)
		log.Printf("[SAFE CHECK] [%s (%s)] distância para %s (%s): %.2f",
			c.Handle.String(), c.PrimaryType, other.Handle.String(), other.PrimaryType, dist)

		if dist < n.SafeDistance {
			log.Printf("[SAFE CHECK] [%s (%s)] ainda não seguro (distância %.2f < %.2f)",
				c.Handle.String(), c.PrimaryType, dist, n.SafeDistance)
			isSafe = false
			break
		}
	}

	if isSafe {
		log.Printf("[SAFE CHECK] [%s (%s)] local seguro alcançado", c.Handle.String(), c.PrimaryType)
		c.ChangeAIState(consts.AIStateIdle)
		return core.StatusSuccess
	}

	return core.StatusRunning
}

func (n *FindIfSafeNode) Reset() {
	// Nada a resetar por enquanto
}
