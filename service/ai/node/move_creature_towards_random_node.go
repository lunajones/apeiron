package node

import (
	"log"
	"math"
	"math/rand"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

// MoveCreatureTowardsRandomNode move a criatura aleatoriamente dentro de um raio controlado
func MoveCreatureTowardsRandomNode(c *creature.Creature, maxDist float64, ctx interface{}) {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-MOVE] [%s (%s)] contexto inválido no MoveCreatureTowardsRandomNode", c.Handle.String(), c.PrimaryType)
		return
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[AI-MOVE] [%s (%s)] grid indisponível no MoveCreatureTowardsRandomNode", c.Handle.String(), c.PrimaryType)
		return
	}

	current := c.GetPosition()

	// Gera deslocamento aleatório dentro do range
	distance := rand.Float64() * maxDist
	angle := rand.Float64() * 2 * math.Pi

	dx := distance * math.Cos(angle)
	dy := distance * math.Sin(angle)
	dz := (rand.Float64() - 0.5) * 2 * (maxDist * 0.2)

	newX := current.FastGlobalX() + dx
	newY := current.FastGlobalY() + dy
	newZ := current.Z + dz

	target := position.FromGlobal(newX, newY, newZ)
	stopAt := c.GetHitboxRadius() + c.GetDesiredBufferDistance()

	// Configura o target
	c.MoveCtrl.SetTarget(target, c.WalkSpeed, stopAt)

	// Aplica o movimento
	if c.MoveCtrl.Update(c, 0.016, grid) { // deltaTime 16ms
		log.Printf(
			"[AI-MOVE] [%s (%s)] chegou no ponto aleatório: (%.2f, %.2f, %.2f)",
			c.Handle.ID, c.PrimaryType,
			target.FastGlobalX(), target.FastGlobalY(), target.Z,
		)
	} else {
		c.SetAction(consts.ActionWalk)
		log.Printf(
			"[AI-MOVE] [%s (%s)] andando aleatoriamente de (%.2f, %.2f, %.2f) para (%.2f, %.2f, %.2f)",
			c.Handle.ID, c.PrimaryType,
			current.FastGlobalX(), current.FastGlobalY(), current.Z,
			target.FastGlobalX(), target.FastGlobalY(), target.Z,
		)
	}
}
