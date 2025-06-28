package node

import (
	"log"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type MoveCreatureTowardsTargetNode struct{}

func (n *MoveCreatureTowardsTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] grid indisponível", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	// Busca o alvo
	var target *creature.Creature
	for _, other := range svcCtx.GetServiceCreatures(c.GetPosition(), c.DetectionRadius) {
		if other.GetHandle().Equals(c.TargetCreatureHandle) {
			target = other
			break
		}
	}

	if target == nil || !target.IsAlive {
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] alvo inválido ou morto (%s)",
			c.Handle.String(), c.PrimaryType, c.TargetCreatureHandle.ID)
		c.ClearTargetHandles()
		c.ChangeAIState(consts.AIStateSearchFood)
		return core.StatusFailure
	}

	stopAt := c.GetHitboxRadius() + target.GetHitboxRadius() + c.GetDesiredBufferDistance() + 0.2
	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())

	if dist <= stopAt {
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] já na distância desejada de [%s] (%.2f ≤ %.2f)",
			c.Handle.String(), c.PrimaryType, target.Handle.String(), dist, stopAt)
		c.IsCrouched = false
		c.SetAction(consts.ActionIdle)
		return core.StatusSuccess
	}

	speed := c.RunSpeed
	if target.AIState == consts.AIStateSleeping {
		speed = c.WalkSpeed * 0.75
		c.IsCrouched = true
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] alvo [%s] dormindo, aproximando em stealth",
			c.Handle.String(), c.PrimaryType, target.Handle.String())
	} else {
		c.IsCrouched = false
	}

	// Seta destino direto no target
	c.MoveCtrl.SetTarget(target.GetPosition(), speed, stopAt)

	if c.MoveCtrl.Update(c, 0.016, grid) {
		log.Printf("[AI-MOVE-TARGET] [%s (%s)] alcançou distância após mover (%.2f ≤ %.2f)",
			c.Handle.String(), c.PrimaryType, dist, stopAt)
		return core.StatusSuccess
	}

	c.SetAction(consts.ActionRun)
	log.Printf("[AI-MOVE-TARGET] [%s (%s)] movendo em direção a [%s] (%.2f > %.2f)",
		c.Handle.String(), c.PrimaryType, target.Handle.String(), dist, stopAt)
	return core.StatusRunning
}

func (n *MoveCreatureTowardsTargetNode) Reset() {
	// Nada a resetar
}
