package position

import (
	"math"
	"math/rand"
)

type Position struct {
	X, Y, Z float64
}

func (p Position) RandomWithinRadius(radius float64) Position {
	if radius == 0 {
		return p
	}

	angle := rand.Float64() * 2 * math.Pi
	dist := rand.Float64() * radius

	return Position{
		X: p.X + dist*math.Cos(angle),
		Y: p.Y + dist*math.Sin(angle),
		Z: p.Z,
	}
}

func CalculateDistance(a, b Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}