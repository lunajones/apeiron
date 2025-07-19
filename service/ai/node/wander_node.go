package node

import (
	"log"
	"math"
	"math/rand"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type WanderNode struct {
	MaxDistance      float64
	SniffChance      float64
	LookAroundChance float64
	IdleChance       float64
	ScratchChance    float64
	VocalizeChance   float64
	PlayChance       float64
	ThreatChance     float64
	CuriousChance    float64
}

func (n *WanderNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI-WANDER] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.MoveCtrl.IsMoving {
		c.SetAnimationState(constslib.AnimationWalk)
		return core.StatusRunning
	}

	if tryWanderDestination(c, svcCtx) {
		c.SetAnimationState(constslib.AnimationWalk)
		return core.StatusRunning
	}

	roll := rand.Float64()
	acc := 0.0

	acc += n.SniffChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationSniff)
		log.Printf("[AI-WANDER] [%s (%s)] farejando o chão.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.LookAroundChance
	if roll < acc {
		angle := (rand.Float64()*2 - 1) * 30
		rotateFacing(c, angle)
		c.SetAnimationState(constslib.AnimationLookAround)
		log.Printf("[AI-WANDER] [%s (%s)] olhando em volta (%.1f graus).", c.Handle.String(), c.PrimaryType, angle)
		return core.StatusSuccess
	}

	acc += n.IdleChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationIdle)
		log.Printf("[AI-WANDER] [%s (%s)] permanece parado.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.ScratchChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationScratch)
		log.Printf("[AI-WANDER] [%s (%s)] coçando-se ou sacudindo o corpo.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.VocalizeChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationVocalize)
		log.Printf("[AI-WANDER] [%s (%s)] vocalizando.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.PlayChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationPlay)
		log.Printf("[AI-WANDER] [%s (%s)] exibindo comportamento lúdico.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.ThreatChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationThreat)
		log.Printf("[AI-WANDER] [%s (%s)] fazendo gesto de ameaça.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	acc += n.CuriousChance
	if roll < acc {
		c.SetAnimationState(constslib.AnimationCurious)
		log.Printf("[AI-WANDER] [%s (%s)] comportamento curioso com o ambiente.", c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	log.Printf("[AI-WANDER] [%s (%s)] nada a fazer neste tick.", c.Handle.String(), c.PrimaryType)
	return core.StatusFailure
}

func tryWanderDestination(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext) bool {
	var dest position.Position
	var distance float64
	found := false

	for i := 0; i < 5; i++ {
		dest = svcCtx.NavMesh.GetRandomWalkablePoint(c.Position, c.MinWanderDistance, c.MaxWanderDistance)
		distance = position.CalculateDistance2D(c.Position, dest)
		if distance >= 3.0 {
			found = true
			break
		}
	}

	if !found {
		log.Printf("[AI-WANDER] [%s] nenhuma posição válida (dist < 3.0)", c.PrimaryType)
		return false
	}

	wanderStop := 0.4 * distance
	if wanderStop < 1.0 {
		wanderStop = 1.0
	}
	if wanderStop > 2.5 {
		wanderStop = 2.5
	}

	c.MoveCtrl.SetMoveTarget(dest, c.WalkSpeed, wanderStop)
	log.Printf("[AI-WANDER] [%s] caminhando para (%.2f, %.2f, %.2f) com stop=%.2f (dist=%.2f)",
		c.PrimaryType, dest.X, dest.Z, dest.Y, wanderStop, distance)
	log.Printf("[AI-WANDER] [%s] está em movimento, aguardando chegada", c.PrimaryType)
	return true
}

func rotateFacing(c *creature.Creature, angleDegrees float64) {
	angleRad := angleDegrees * (math.Pi / 180)
	c.SetFacingDirection(position.RotateVector2D(c.GetFacingDirection(), angleRad))
}

func (n *WanderNode) Reset(c *creature.Creature) {}
