package defensive

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type RetreatNode struct{}

func (n *RetreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[RETREAT] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[RETREAT] [%s] sem alvo para recuar", c.Handle.String())
		return core.StatusFailure
	}

	// Calcula direção de recuo
	dir := position.NewVector3DFromTo(target.GetPosition(), c.GetPosition()).Normalize()
	dest := c.GetPosition().AddVector3D(dir.Scale(2.0)) // Passo de recuo

	// Move com half walk speed
	c.MoveCtrl.SetMoveIntent(dest, c.WalkSpeed*0.5, 0.0)

	c.SetAnimationState(constslib.AnimationWalk)
	log.Printf("[RETREAT] [%s] recuando cauteloso para %v", c.Handle.String(), dest)

	svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
		SourceHandle: c.Handle,
		BehaviorType: "RetreatPerformed",
		Timestamp:    time.Now(),
	})

	return core.StatusRunning
}

func (n *RetreatNode) Reset() {
	// Nada a resetar
}
