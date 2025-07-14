package helper

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type FaceTargetNode struct{}

func (n *FaceTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[FACE-TARGET] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[FACE-TARGET] [%s] sem alvo para virar", c.Handle.String())
		return core.StatusFailure
	}

	dirVec := target.GetPosition().Sub2D(c.Position).Normalize()
	c.SetFacingDirection(dirVec)
	log.Printf("[FACE-TARGET] [%s] virou para [%s]", c.Handle.String(), target.GetHandle().String())
	return core.StatusSuccess
}

func (n *FaceTargetNode) Reset(c *creature.Creature) {}
