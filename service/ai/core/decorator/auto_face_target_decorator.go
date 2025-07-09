package helper

import (
	"log"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type AutoFaceTargetDecorator struct {
	Child core.BehaviorNode
}

func NewAutoFaceTargetDecorator(child core.BehaviorNode) *AutoFaceTargetDecorator {
	return &AutoFaceTargetDecorator{
		Child: child,
	}
}

func (d *AutoFaceTargetDecorator) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AUTO-FACE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target != nil {
		dir := position.Vector2D{
			X: target.GetPosition().X - c.Position.X,
			Z: target.GetPosition().Z - c.Position.Z,
		}.Normalize()
		c.SetFacingDirection(dir)
		// log.Printf("[AUTO-FACE] [%s (%s)] ajustou facing para alvo",
		// 	c.Handle.String(), c.PrimaryType)

	}

	return d.Child.Tick(c, ctx)
}

func (d *AutoFaceTargetDecorator) Reset() {
	if d.Child != nil {
		d.Child.Reset()
	}
}
