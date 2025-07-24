package finder

import (
	"math"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
)

func FindTargetByHandles(
	searcherHandle, creatureHandle, playerHandle handle.EntityHandle,
	svcCtx *dynamic_context.AIServiceContext,
) model.Targetable {
	targets := svcCtx.GetCachedTargets(searcherHandle)

	for _, t := range targets {
		if t.GetHandle().Equals(creatureHandle) && t.IsAlive() {
			return t
		}
		if t.GetHandle().Equals(playerHandle) && t.IsAlive() {
			return t
		}
	}
	return nil
}

// Retorna true se 'target' estiver dentro do cone de visão de 'source'
func IsInFieldOfView(source, target model.Targetable, fovDegrees float64) bool {
	dirVec3D := source.GetFacingDirection() // deve retornar um Vector3D
	forward := position.Vector2D{X: dirVec3D.X, Z: dirVec3D.Z}.Normalize()

	toTarget3D := position.NewVector3DFromTo(source.GetPosition(), target.GetPosition())
	toTarget := position.Vector2D{X: toTarget3D.X, Z: toTarget3D.Z}.Normalize()

	dot := forward.Dot(toTarget)
	angle := math.Acos(dot) * (180.0 / math.Pi)

	return angle <= (fovDegrees / 2.0)
}

func FindNearbyAllies(ctx *dynamic_context.AIServiceContext, self model.Targetable, faction string, maxDist float64) []model.Targetable {
	var result []model.Targetable
	candidates := ctx.SpatialIndex.Query(self.GetPosition(), maxDist)

	for _, t := range candidates {
		if t == nil || t.GetHandle().Equals(self.GetHandle()) {
			continue
		}

		if !t.IsAlive() {
			continue
		}

		if t.GetFaction() != faction {
			continue
		}

		dist := position.CalculateDistance(self.GetPosition(), t.GetPosition())
		if dist <= maxDist {
			result = append(result, t)
		}
	}

	return result
}

func FindNearbyTargets(ctx *dynamic_context.AIServiceContext, self model.Targetable, maxDist float64) []model.Targetable {
	var result []model.Targetable
	candidates := ctx.SpatialIndex.Query(self.GetPosition(), maxDist)

	for _, t := range candidates {
		if t == nil || t.GetHandle().Equals(self.GetHandle()) {
			continue
		}
		if !t.IsAlive() {
			continue
		}
		dist := position.CalculateDistance(self.GetPosition(), t.GetPosition())
		if dist <= maxDist {
			result = append(result, t)
		}
	}

	return result
}
