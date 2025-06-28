package node

import (
	"log"
	"math"
	"math/rand"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type RandomIdleNode struct{}

func (n *RandomIdleNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-IDLE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[AI-IDLE] [%s (%s)] grid indisponível", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	nearbyCreatures := svcCtx.GetServiceCreatures(c.GetPosition(), c.DetectionRadius)

	nearThreat := false
	for _, other := range nearbyCreatures {
		if other == nil || other.Handle.Equals(c.Handle) || !other.IsAlive {
			continue
		}
		if other.IsHostile && position.CalculateDistance(c.GetPosition(), other.GetPosition()) < 6.0 {
			nearThreat = true
			break
		}
	}

	walkChance := 0.4
	if nearThreat {
		walkChance = 0.7
	}

	roll := rand.Float64()

	switch {
	case roll < walkChance:
		angle := rand.Float64() * 2 * math.Pi
		dist := 1.0 + rand.Float64()*1.0
		newX := c.GetPosition().FastGlobalX() + math.Cos(angle)*dist
		newY := c.GetPosition().FastGlobalY() + math.Sin(angle)*dist
		dest := position.FromGlobal(newX, newY, c.GetPosition().Z)

		stopAt := c.GetHitboxRadius() + c.GetDesiredBufferDistance()
		c.MoveCtrl.SetTarget(dest, c.WalkSpeed, stopAt)
		c.SetAction(consts.ActionWalk)

		c.MoveCtrl.Update(c, 0.016, grid)

		log.Printf("[AI-IDLE] [%s (%s)] caminha ociosamente (roll=%.2f, chance=%.2f)",
			c.Handle.String(), c.PrimaryType, roll, walkChance)

	case rand.Float64() < 0.2:
		c.SetAction(consts.ActionSniff)
		log.Printf("[AI-IDLE] [%s (%s)] fareja o chão", c.Handle.String(), c.PrimaryType)

	default:
		c.SetAction(consts.ActionIdle)
		log.Printf("[AI-IDLE] [%s (%s)] permanece parado", c.Handle.String(), c.PrimaryType)
	}

	return core.StatusSuccess
}

func (n *RandomIdleNode) Reset() {
	// Nada a resetar
}
