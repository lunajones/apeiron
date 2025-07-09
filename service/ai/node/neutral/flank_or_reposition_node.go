package neutral

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type FlankOrRepositionNode struct{}

func (n *FlankOrRepositionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[FLANK] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[FLANK] [%s] sem alvo para flanquear", c.Handle.String())
		return core.StatusFailure
	}

	dirVec3D := target.GetPosition().Sub(c.Position)
	dirVec2D := position.Vector2D{X: dirVec3D.X, Z: dirVec3D.Z}.Normalize()

	baseAngle := (math.Pi / 4) + rand.Float64()*(math.Pi/2) // 45° a 135°
	if rand.Float64() < 0.5 {
		baseAngle = -baseAngle
	}
	flankVec := position.RotateVector2D(dirVec2D, baseAngle)

	dist := 1.0 + rand.Float64()*1.5
	flankVec = flankVec.Multiply(dist)

	dest := c.Position.AddVector3D(position.Vector3D{X: flankVec.X, Y: 0, Z: flankVec.Z})

	if position.CalculateDistance(c.MoveCtrl.CurrentIntentDest, dest) < 0.5 {
		// Já está indo para um destino parecido, deixa seguir
		return core.StatusRunning
	}

	if svcCtx.NavMesh.IsWalkable(dest) {
		c.MoveCtrl.SetMoveIntent(dest, c.WalkSpeed, 0.0) // stopAt = 0.0 → chega exatamente
		log.Printf("[FLANK] [%s] flanqueando para (%.2f, %.2f)", c.Handle.String(), dest.X, dest.Z)

		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "FlankExecuted",
			Timestamp:    time.Now(),
		})

		return core.StatusSuccess
	}

	log.Printf("[FLANK] [%s] destino inválido (%.2f, %.2f)", c.Handle.String(), dest.X, dest.Z)
	return core.StatusFailure
}

func (n *FlankOrRepositionNode) Reset() {}
