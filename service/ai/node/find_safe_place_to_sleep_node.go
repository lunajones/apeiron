package node

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type FindSafePlaceToSleepNode struct{}

func (n *FindSafePlaceToSleepNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SAFE SLEEP] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	var threats []position.Position

	for _, t := range svcCtx.GetCachedTargets(c.Handle) {
		other, ok := t.(*creature.Creature)
		if !ok || !other.Alive || other.Handle.Equals(c.Handle) {
			continue
		}
		if other.PrimaryType != c.PrimaryType {
			threats = append(threats, other.Position)
		}
	}

	if len(threats) == 0 {
		log.Printf("[SAFE SLEEP] [%s (%s)] local seguro detectado, entrando em estado de sono", c.Handle.String(), c.PrimaryType)
		c.MoveCtrl.CurrentPath = nil
		c.MoveCtrl.IsMoving = false
		c.MoveCtrl.Intent.HasIntent = false
		c.SetAnimationState(constslib.AnimationSleep)
		c.ChangeAIState(constslib.AIStateSleeping)
		return core.StatusSuccess
	}

	dest := svcCtx.NavMesh.GetEscapePoint(c.Position, threats, 6.0)
	c.MoveCtrl.SetMoveIntent(dest, c.WalkSpeed*0.5, 1.5)

	if c.AIState != constslib.AIStateSeekingSafePlace {
		c.ChangeAIState(constslib.AIStateSeekingSafePlace)
		c.SetAnimationState(constslib.AnimationWalk)
	}

	log.Printf("[SAFE SLEEP] [%s (%s)] buscando local seguro em (%.2f, %.2f, %.2f)",
		c.Handle.String(), c.PrimaryType, dest.X, dest.Z, dest.Y)

	return core.StatusRunning
}

func (n *FindSafePlaceToSleepNode) Reset() {}
