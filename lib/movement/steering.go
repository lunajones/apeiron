package movement

import (
	"math"
	"time"

	"github.com/lunajones/apeiron/lib/position"
)

// SteeringOutput representa a força calculada do steering
type SteeringOutput struct {
	Linear  position.Vector3D
	Angular float64
}

// SteeringBehaviors agrupa os parâmetros de steering para um movimento complexo
type SteeringBehaviors struct {
	MaxAcceleration float64
	MaxSpeed        float64
}

// ApplySteering combina os resultados e aplica limites
func (s *SteeringOutput) ApplyLimits(maxAccel, maxSpeed float64) {
	mag := s.Linear.Magnitude()
	if mag > maxAccel {
		s.Linear = s.Linear.Normalize().Scale(maxAccel)
	}
	if mag > maxSpeed {
		s.Linear = s.Linear.Normalize().Scale(maxSpeed)
	}
}

// Seek gera força para seguir um alvo
func Seek(current, target position.Position, behaviors SteeringBehaviors) SteeringOutput {
	direction := position.NewVector3DFromTo(current, target)
	linear := direction.Normalize().Scale(behaviors.MaxAcceleration)
	return SteeringOutput{Linear: linear}
}

// Flee gera força para fugir de um alvo
func Flee(current, threat position.Position, behaviors SteeringBehaviors) SteeringOutput {
	direction := position.NewVector3DFromTo(threat, current)
	linear := direction.Normalize().Scale(behaviors.MaxAcceleration)
	return SteeringOutput{Linear: linear}
}

// Arrive ajusta para parar suavemente
func Arrive(current, target position.Position, behaviors SteeringBehaviors, slowRadius, stopRadius float64) SteeringOutput {
	direction := position.NewVector3DFromTo(current, target)
	distance := direction.Magnitude()
	if distance < stopRadius {
		return SteeringOutput{}
	}
	goalSpeed := behaviors.MaxSpeed
	if distance < slowRadius {
		goalSpeed = behaviors.MaxSpeed * (distance / slowRadius)
	}
	linear := direction.Normalize().Scale(goalSpeed)
	return SteeringOutput{Linear: linear}
}

// AvoidObstacle tenta gerar uma força de desvio
func AvoidObstacle(current position.Position, obstacles []position.Position, behaviors SteeringBehaviors, avoidRadius float64) SteeringOutput {
	avoid := position.Vector3D{}
	for _, obs := range obstacles {
		direction := position.NewVector3DFromTo(current, obs)
		dist := direction.Magnitude()
		if dist < avoidRadius && dist > 0 {
			force := direction.Normalize().Scale(-1 / dist)
			avoid = avoid.Add(force)
		}
	}
	if avoid.Magnitude() > 0 {
		avoid = avoid.Normalize().Scale(behaviors.MaxAcceleration)
	}
	return SteeringOutput{Linear: avoid}
}

// Wander exemplo simples de variação de direção
func Wander(currentDir position.Vector3D, jitter, radius float64) SteeringOutput {
	randAngle := (randFloat64()*2 - 1) * jitter
	angle := math.Atan2(currentDir.Y, currentDir.X) + randAngle
	linear := position.Vector3D{
		X: math.Cos(angle) * radius,
		Y: math.Sin(angle) * radius,
		Z: 0,
	}
	return SteeringOutput{Linear: linear}
}

// Helper rand float
func randFloat64() float64 {
	return (float64(randInt63() & 0x7fffffffffffffff)) / (1 << 63)
}

func randInt63() int64 {
	return time.Now().UnixNano() & 0x7fffffffffffffff
}
