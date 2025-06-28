package physics

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/lib/position"
)

func MoveTowards(current position.Position, target position.Position, maxStep float64) position.Position {
	dx := float64(target.GridX-current.GridX) + (target.OffsetX - current.OffsetX)
	dy := float64(target.GridY-current.GridY) + (target.OffsetY - current.OffsetY)
	dz := float64(target.Z - current.Z)

	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dist < 0.001 {
		log.Printf("[Physics] MoveTowards ignorado: destino muito próximo.")
		return current
	}

	step := math.Min(maxStep, dist)

	newPos := position.Position{
		GridX:   current.GridX,
		GridY:   current.GridY,
		OffsetX: current.OffsetX + (dx/dist)*step,
		OffsetY: current.OffsetY + (dy/dist)*step,
		Z:       current.Z + (dz/dist)*step,
	}

	log.Printf("[Physics] MoveTowards: step=%.2f ΔX=%.2f ΔY=%.2f ΔZ=%.2f", step, dx, dy, dz)
	return newPos
}
