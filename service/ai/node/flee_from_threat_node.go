package node

import (
	"log"
	"math"
	"math/rand"

	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeFromThreatNode struct {
	SafeDistance float64
}

func (n *FleeFromThreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-FLEE] [%s (%s)] contexto inválido", c.Handle.ID, c.PrimaryType)
		return core.StatusFailure
	}

	threat := resolveThreat(c, svcCtx.GetCachedTargets(c.Handle))
	if threat == nil {
		log.Printf("[AI-FLEE] [%s (%s)] nenhuma ameaça encontrada", c.Handle.ID, c.PrimaryType)
		c.ClearTargetHandles()
		return core.StatusFailure
	}

	c.TargetCreatureHandle = threat.GetHandle()

	dist := position.CalculateDistance(c.Position, threat.GetPosition())
	if dist >= n.SafeDistance {
		log.Printf("[AI-FLEE] [%s (%s)] local seguro alcançado (%.2f u) longe de [%s]", c.Handle.ID, c.PrimaryType, dist, threat.GetHandle().ID)
		c.ClearTargetHandles()
		c.ChangeAIState(consts.AIStateIdle)
		c.SetAnimationState(consts.AnimationIdle)
		return core.StatusSuccess
	}

	dirX := c.Position.X - threat.GetPosition().X
	dirZ := c.Position.Z - threat.GetPosition().Z
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
	targetX := c.Position.X + dirX*fleeDistance
	targetZ := c.Position.Z + dirZ*fleeDistance
	newPos := position.Position{
		X: targetX,
		Y: c.Position.Y,
		Z: targetZ,
	}

	path := svcCtx.NavMesh.FindPath(c.Position, newPos)
	if len(path) == 0 {
		closestPoly := svcCtx.NavMesh.FindClosestPolygon(c.Position)
		if closestPoly != nil {
			center := closestPoly.CenterPosition()
			path = svcCtx.NavMesh.FindPath(c.Position, center)
		}
	}

	if len(path) == 0 {
		log.Printf("[AI-FLEE] [%s (%s)] fuga impossível, sem caminho no NavMesh", c.Handle.ID, c.PrimaryType)
		return core.StatusFailure
	}

	c.MoveCtrl.SetPath(path, c)
	c.SetAnimationState(consts.AnimationRun)

	log.Printf("[AI-FLEE] [%s (%s)] fugindo de [%s]: distância atual %.2f u, path com %d pontos",
		c.Handle.ID, c.PrimaryType, threat.GetHandle().ID, dist, len(path))

	return core.StatusRunning
}

func (n *FleeFromThreatNode) Reset(c *creature.Creature) {}

func resolveThreat(c *creature.Creature, targets []model.Targetable) model.Targetable {
	for _, t := range targets {
		other, ok := t.(*creature.Creature)
		if !ok || !other.Alive || other.Handle.Equals(c.Handle) {
			continue
		}
		if creature.AreEnemies(c, other) {
			return other
		}
	}
	return nil
}
