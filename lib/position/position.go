package position

import (
	"fmt"
	"math"
	"math/rand"
)

type Position struct {
	X float64
	Y float64
	Z float64
}

// Offset cria uma nova posição deslocada
func (p Position) Offset(dx, dz float64) Position {
	return Position{
		X: p.X + dx,
		Y: p.Y,
		Z: p.Z + dz,
	}
}

// Equals compara duas posições com margem mínima de erro
func (p Position) Equals(other Position) bool {
	const epsilon = 0.001
	return math.Abs(p.X-other.X) < epsilon &&
		math.Abs(p.Z-other.Z) < epsilon &&
		math.Abs(p.Y-other.Y) < epsilon
}

// CalculateDistance calcula distância 3D entre duas posições

func CalculateDistance(a, b Position) float64 {
	dx := a.X - b.X
	dz := a.Z - b.Z
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dz*dz + dy*dy)
}

// CalculateDistance2D calcula distância 2D no plano X-Z
func CalculateDistance2D(a, b Position) float64 {
	dx := a.X - b.X
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dz*dz)
}

// RandomWithinRadius gera nova posição aleatória dentro de um raio no plano
func (p Position) RandomWithinRadius(radius float64) Position {
	minDist := math.Min(radius*0.6, radius-0.1)
	dist := minDist + rand.Float64()*(radius-minDist)
	angle := rand.Float64() * 2 * math.Pi
	newX := p.X + dist*math.Cos(angle)
	newZ := p.Z + dist*math.Sin(angle)
	return Position{
		X: newX,
		Y: p.Y,
		Z: newZ,
	}
}

func (p Position) AddOffset(dx, dz float64) Position {
	return Position{
		X: p.X + dx,
		Y: p.Y,
		Z: p.Z + dz,
	}
}

// AddVector3D soma um Vector3D a uma Position
func (p Position) AddVector3D(v Vector3D) Position {
	return Position{
		X: p.X + v.X,
		Y: p.Y + v.Y,
		Z: p.Z + v.Z,
	}
}

func (p Position) ToVector2D() Vector2D {
	return Vector2D{X: p.X, Z: p.Z}
}

// Sub calcula o vetor diferença entre duas posições (3D)
func (p Position) Sub(other Position) Vector3D {
	return Vector3D{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p Position) Sub2D(other Position) Vector2D {
	return Vector2D{
		X: p.X - other.X,
		Z: p.Z - other.Z,
	}
}

// AddVector2D soma um Vector2D (X,Z) a uma Position, mantendo o Y original
func (p Position) AddVector2D(v Vector2D) Position {
	return Position{
		X: p.X + v.X,
		Y: p.Y, // mantemos altura original
		Z: p.Z + v.Z,
	}
}

func LerpVector2D(v1, v2 Vector2D, t float64) Vector2D {
	return Vector2D{
		X: v1.X + (v2.X-v1.X)*t,
		Z: v1.Z + (v2.Z-v1.Z)*t,
	}
}

func (p Position) LerpTo(dest Position, t float64) Position {
	return Position{
		X: p.X + (dest.X-p.X)*t,
		Y: p.Y + (dest.Y-p.Y)*t,
		Z: p.Z + (dest.Z-p.Z)*t,
	}
}

func (p Position) Key() string {
	gridX := int(math.Floor(p.X))
	gridZ := int(math.Floor(p.Z))
	return fmt.Sprintf("%d:%d", gridX, gridZ)
}

func CalculateDirection2D(start, end Position) Vector2D {
	dx := end.X - start.X
	dz := end.Z - start.Z
	length := math.Sqrt(dx*dx + dz*dz)
	if length == 0 {
		return Vector2D{X: 0, Z: 0}
	}
	return Vector2D{X: dx / length, Z: dz / length}
}
