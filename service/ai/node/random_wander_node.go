package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomWanderNode struct{}

func (n *RandomWanderNode) Tick(c *creature.Creature, ctx ai_context.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando RandomWanderNode", c.ID)

	rand.Seed(time.Now().UnixNano())

	// Gera deslocamento aleatório pequeno (entre -1.0 e +1.0)
	dx := (rand.Float64() - 0.5) * 2
	dz := (rand.Float64() - 0.5) * 2

	newPos := position.Position{
		X: c.Position.X + dx,
		Y: c.Position.Y,
		Z: c.Position.Z + dz,
	}

	// Atualiza posição
	log.Printf("[AI] %s andando de (%.2f, %.2f, %.2f) para (%.2f, %.2f, %.2f)",
		c.ID, c.Position.X, c.Position.Y, c.Position.Z, newPos.X, newPos.Y, newPos.Z)

	c.Position = newPos
	c.SetAction(creature.ActionWalk)

	return core.StatusSuccess
}
