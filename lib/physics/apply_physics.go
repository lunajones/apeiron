package physics

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
)

func ApplyPhysics(mov model.Movable, velocity *position.Vector3D, acceleration position.Vector3D, deltaTime float64, checkCollision bool, mesh *navmesh.NavMesh, all []model.Targetable) bool {
	if deltaTime <= 0 {
		return false
	}

	velocity.X = acceleration.X
	velocity.Y = acceleration.Y
	velocity.Z = acceleration.Z

	currentPos := mov.GetPosition()
	newX := currentPos.X + velocity.X*deltaTime
	newZ := currentPos.Z + velocity.Z*deltaTime
	newY := currentPos.Y + velocity.Y*deltaTime
	newPos := position.Position{X: newX, Y: newY, Z: newZ}

	if checkCollision {
		if CheckCollision(newPos, mov.GetHitboxRadius(), mov.GetHandle(), mesh, all) {
			velocity.Zero()
			return true
		}
	}

	mov.SetPosition(newPos)

	if velocity.X != 0 || velocity.Z != 0 {
		mag := math.Sqrt(velocity.X*velocity.X + velocity.Z*velocity.Z)
		if mag > 0 {
			mov.SetFacingDirection(position.Vector2D{
				X: velocity.X / mag,
				Z: velocity.Z / mag,
			})
		}
	}

	// log.Printf("[PHYSICS] [%s] ApplyPhysics aplicou: newPos=(%.2f, %.2f, %.2f)", mov.GetHandle().ID, newX, newY, newZ)
	return false
}

func CheckCollision(newPos position.Position, radius float64, selfHandle handle.EntityHandle, mesh *navmesh.NavMesh, all []model.Targetable) bool {
	if mesh != nil && !mesh.IsWalkable(newPos) {
		log.Printf("[PHYSICS] [%s] colisão NavMesh em (%.2f, %.2f)", selfHandle.ID, newPos.X, newPos.Z)
		return true
	}

	for _, t := range all {
		if t.GetHandle().Equals(selfHandle) {
			continue
		}
		if !t.IsAlive() {
			continue
		}

		dx := newPos.X - t.GetLastPosition().X
		dz := newPos.Z - t.GetLastPosition().Z
		dist := math.Sqrt(dx*dx + dz*dz)

		if dist < radius+t.GetHitboxRadius() {
			log.Printf("[PHYSICS] [%s] colisão com entidade %s", selfHandle.ID, t.GetHandle().ID)
			return true
		}
	}
	return false
}
