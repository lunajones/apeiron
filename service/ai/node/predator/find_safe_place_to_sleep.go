package predator

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

type FindSafePlaceToSleepNode struct{}

func (n *FindSafePlaceToSleepNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SAFE SLEEP - PREDATOR] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid() // 🚀 você vai implementar esse método no seu context
	if grid == nil {
		log.Printf("[SAFE SLEEP - PREDATOR] [%s (%s)] grid não disponível no contexto", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	creatures := svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius)
	centerX, centerY := 0.0, 0.0
	count := 0.0

	for _, other := range creatures {
		if other.Handle.ID == c.Handle.ID || !other.IsAlive {
			continue
		}
		if other.HasTag("Predator") {
			centerX += other.Position.FastGlobalX()
			centerY += other.Position.FastGlobalY()
			count++
		}
	}

	if count == 0 {
		log.Printf("[SAFE SLEEP - PREDATOR] [%s (%s)] nenhum outro predador detectado, local seguro", c.Handle.String(), c.PrimaryType)
		c.ChangeAIState(consts.AIStateSleeping)
		return core.StatusSuccess
	}

	centerX /= count
	centerY /= count

	cX := c.Position.FastGlobalX()
	cY := c.Position.FastGlobalY()

	dirX := cX - centerX
	dirY := cY - centerY
	mag := math.Hypot(dirX, dirY)
	if mag == 0 {
		angle := rand.Float64() * 2 * math.Pi
		dirX = math.Cos(angle)
		dirY = math.Sin(angle)
		mag = 1
	}

	dirX /= mag
	dirY /= mag

	destX := cX + dirX*6.0
	destY := cY + dirY*6.0
	dest := position.FromGlobal(destX, destY, c.Position.Z)

	if c.MoveCtrl.IsMoving {
		c.MoveCtrl.Update(c, 0.016, grid)
		return core.StatusRunning
	}

	slowWalkSpeed := c.WalkSpeed * 0.5
	c.MoveCtrl.SetTarget(dest, slowWalkSpeed, 1.5)
	log.Printf("[SAFE SLEEP - PREDATOR] [%s (%s)] buscando local seguro (destino: %.2f, %.2f, %.2f)",
		c.Handle.String(), c.PrimaryType, dest.FastGlobalX(), dest.FastGlobalY(), dest.Z)
	return core.StatusRunning
}

func (n *FindSafePlaceToSleepNode) Reset() {}
