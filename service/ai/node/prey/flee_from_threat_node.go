package prey

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

type FleeFromThreatNode struct{}

func (n *FleeFromThreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-FLEE] [%s] contexto inválido em FleeFromThreatNode", c.Handle.ID)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[AI-FLEE] [%s] grid indisponível no contexto", c.Handle.ID)
		return core.StatusFailure
	}

	threat := resolveThreat(c, svcCtx)
	if threat == nil {
		log.Printf("[AI-FLEE] [%s] nenhuma ameaça encontrada", c.Handle.ID)
		c.ClearTargetHandles()
		return core.StatusFailure
	}

	c.TargetCreatureHandle = threat.Handle

	dist := position.CalculateDistance(c.Position, threat.Position)
	if dist >= 6.0 {
		log.Printf("[AI-FLEE] [%s] local seguro alcançado (%.2f u) longe de [%s]", c.Handle.ID, dist, threat.Handle.ID)
		c.ClearTargetHandles()
		c.ChangeAIState(consts.AIStateIdle)
		c.SetAction(consts.ActionIdle)
		return core.StatusSuccess
	}

	dirX := c.Position.FastGlobalX() - threat.Position.FastGlobalX()
	dirZ := c.Position.Z - threat.Position.Z
	mag := math.Hypot(dirX, dirZ)

	if mag == 0 {
		angle := rand.Float64() * 2 * math.Pi
		dirX = math.Cos(angle)
		dirZ = math.Sin(angle)
		mag = 1
	}

	dirX /= mag
	dirZ /= mag

	fleeDistance := 4.0
	targetX := c.Position.FastGlobalX() + dirX*fleeDistance
	targetZ := c.Position.Z + dirZ*fleeDistance
	newPos := position.FromGlobal(targetX, c.Position.FastGlobalY(), targetZ)

	stopAt := c.HitboxRadius + threat.HitboxRadius + c.DesiredBufferDistance
	speed := c.GetCurrentSpeed()

	if c.MoveCtrl.IsMoving {
		c.MoveCtrl.Update(c, 0.016, grid)
		return core.StatusRunning
	}

	c.MoveCtrl.SetTarget(newPos, speed, stopAt)
	log.Printf("[AI-FLEE] [%s] fugindo de [%s] (%.2f u) para (%.2f, %.2f, %.2f)",
		c.Handle.ID, threat.Handle.ID, dist,
		newPos.FastGlobalX(), newPos.FastGlobalY(), newPos.Z)

	return core.StatusRunning
}

func (n *FleeFromThreatNode) Reset() {}

func resolveThreat(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext) *creature.Creature {
	creatures := svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius)
	for _, other := range creatures {
		if other.Handle.Equals(c.Handle) || !other.IsAlive || !other.IsHostile {
			continue
		}
		return other
	}
	if !c.TargetCreatureHandle.IsEmpty() {
		for _, other := range creatures {
			if other.Handle.Equals(c.TargetCreatureHandle) && other.IsAlive {
				return other
			}
		}
	}
	return nil
}
