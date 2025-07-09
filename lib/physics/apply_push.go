package physics

import (
	"log"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
)

func PushTargetFromWall(
	target model.Targetable,
	attacker model.Movable,
	navMesh *navmesh.NavMesh,
	pushDistance float64,
) {
	dirVec := position.NewVector3DFromTo(target.GetPosition(), attacker.GetPosition()).Normalize()

	// Calcula nova posição proposta para o empurrão direto
	newPos := target.GetPosition().AddVector3D(dirVec.Scale(pushDistance))

	if navMesh.IsWalkable(newPos) {
		target.SetPosition(newPos)
		log.Printf("[PUSH] [%s] alvo empurrado diretamente para (%.2f, %.2f)",
			target.GetHandle().ID, newPos.X, newPos.Z)
		return
	}

	// Tenta ajuste lateral (primeiro lado)
	perpVec := position.Vector3D{
		X: -dirVec.Z,
		Y: 0,
		Z: dirVec.X,
	}.Normalize()

	adjustedPos := target.GetPosition().AddVector3D(perpVec.Scale(pushDistance * 0.5))
	if navMesh.IsWalkable(adjustedPos) {
		target.SetPosition(adjustedPos)
		log.Printf("[PUSH] [%s] avanço bloqueado, alvo empurrado lateralmente para (%.2f, %.2f)",
			target.GetHandle().ID, adjustedPos.X, adjustedPos.Z)
		return
	}

	// Tenta ajuste lateral no lado oposto
	oppositePos := target.GetPosition().AddVector3D(perpVec.Scale(-pushDistance * 0.5))
	if navMesh.IsWalkable(oppositePos) {
		target.SetPosition(oppositePos)
		log.Printf("[PUSH] [%s] avanço bloqueado, alvo empurrado para lado oposto em (%.2f, %.2f)",
			target.GetHandle().ID, oppositePos.X, oppositePos.Z)
		return
	}

	log.Printf("[PUSH] [%s] empurrão bloqueado, nenhuma posição válida encontrada", target.GetHandle().ID)
}
