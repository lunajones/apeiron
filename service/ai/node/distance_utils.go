package node

import (
	"math"

	"github.com/lunajones/apeiron/lib/position"
)

func calculateDistance(a, b position.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}