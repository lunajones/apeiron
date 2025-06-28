package physics

import (
	"log"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
)

// ApplyPhysics aplica física básica: atualiza velocidade e posição, lida com colisão.
func ApplyPhysics(mov model.Movable, velocity *position.Vector3D, acceleration position.Vector3D, deltaTime float64, checkCollision bool) {
	if deltaTime <= 0 {
		return
	}

	// Atualiza velocidade: v = v0 + a * t
	velocity.X += acceleration.X * deltaTime
	velocity.Y += acceleration.Y * deltaTime
	velocity.Z += acceleration.Z * deltaTime

	// Calcula novo passo: s = s0 + v * t
	currentPos := mov.GetPosition()
	newX := currentPos.FastGlobalX() + velocity.X*deltaTime
	newY := currentPos.FastGlobalY() + velocity.Y*deltaTime
	newZ := currentPos.Z + velocity.Z*deltaTime
	newPos := position.FromGlobal(newX, newY, newZ)

	// Checa colisão se necessário
	if checkCollision && CheckCollision(newPos, mov.GetHitboxRadius()) {
		log.Printf("[PHYSICS] [%s] colisão detectada em (%.2f, %.2f, %.2f), movimento bloqueado", mov.GetHandle().ID, newPos.FastGlobalX(), newPos.FastGlobalY(), newPos.Z)
		velocity.Zero()
		return
	}

	// Aplica nova posição
	mov.SetPosition(newPos)

	// Atualiza direção se estiver se movendo
	mag := velocity.Magnitude()
	if mag > 0 {
		mov.SetFacingDirection(position.Vector2D{
			X: velocity.X / mag,
			Y: velocity.Y / mag,
		})
	}
}
