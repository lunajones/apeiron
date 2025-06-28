package node

import (
	"log"
	"math/rand"

	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type RandomWanderNode struct{}

func (n *RandomWanderNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-WANDER] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[AI-WANDER] [%s (%s)] grid indisponível", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.MoveCtrl.IsMoving {
		if c.MoveCtrl.Update(c, 0.016, grid) {
			return core.StatusSuccess
		}
		c.SetAction(consts.ActionWalk)
		return core.StatusRunning
	}

	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		offset := position.Vector2D{
			X: (rand.Float64() - 0.5) * 2 * 1.0,
			Y: (rand.Float64() - 0.5) * 2 * 1.0,
		}

		destX := c.Position.FastGlobalX() + offset.X
		destY := c.Position.FastGlobalY() + offset.Y
		dest := position.FromGlobal(destX, destY, c.Position.Z)

		if physics.IsWalkable(dest, c.GetHitboxRadius()) {
			c.MoveCtrl.SetTarget(dest, c.WalkSpeed, c.HitboxRadius+c.DesiredBufferDistance)
			_ = c.MoveCtrl.Update(c, 0.016, grid)
			c.SetAction(consts.ActionWalk)

			log.Printf("[AI-WANDER] [%s (%s)] caminhando para destino válido: (%.2f, %.2f, %.2f)",
				c.Handle.ID, c.PrimaryType, dest.FastGlobalX(), dest.FastGlobalY(), dest.Z)
			return core.StatusRunning
		}

		log.Printf("[AI-WANDER] [%s (%s)] tentativa %d bloqueada, recalculando",
			c.Handle.ID, c.PrimaryType, attempt)
	}

	log.Printf("[AI-WANDER] [%s (%s)] falhou em encontrar destino após %d tentativas",
		c.Handle.ID, c.PrimaryType, maxAttempts)
	return core.StatusFailure
}

func (n *RandomWanderNode) Reset() {
	// Nada a resetar
}
