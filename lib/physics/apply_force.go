package physics

import (
	"github.com/lunajones/apeiron/lib/position"
)

// ApplyForce aplica uma força sobre a aceleração atual do Movable.
// Se quiser, pode adicionar massa como divisor no futuro.
func ApplyForce(acceleration *position.Vector3D, force position.Vector3D, maxAccel float64) {
	acceleration.X += force.X
	acceleration.Y += force.Y
	acceleration.Z += force.Z

	// Limita a aceleração ao máximo permitido
	mag := acceleration.Magnitude()
	if mag > maxAccel && maxAccel > 0 {
		*acceleration = acceleration.Normalize().Scale(maxAccel)
	}
}
