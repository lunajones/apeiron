package position

import (
	"math"
	"math/rand"
)

const GridSize = 100.0 // tamanho de cada grid em metros

// Position representa uma coordenada no mundo segmentado em grids
type Position struct {
	GridX, GridY     int
	OffsetX, OffsetY float64
	Z                float64
}

// ToGlobal retorna as coordenadas globais (absolutas) X,Y
func (p Position) ToGlobal() (float64, float64) {
	return float64(p.GridX)*GridSize + p.OffsetX,
		float64(p.GridY)*GridSize + p.OffsetY
}

// FastGlobalX retorna X global sem alocação
func (p Position) FastGlobalX() float64 {
	return float64(p.GridX)*GridSize + p.OffsetX
}

// FastGlobalY retorna Y global sem alocação
func (p Position) FastGlobalY() float64 {
	return float64(p.GridY)*GridSize + p.OffsetY
}

// FastGlobalZ retorna Z global
func (p Position) FastGlobalZ() float64 {
	return p.Z
}

// FromGlobal cria uma Position a partir de coordenadas absolutas
func FromGlobal(x, y, z float64) Position {
	gx := int(x / GridSize)
	gy := int(y / GridSize)
	ox := x - float64(gx)*GridSize
	oy := y - float64(gy)*GridSize
	return Position{
		GridX:   gx,
		GridY:   gy,
		OffsetX: ox,
		OffsetY: oy,
		Z:       z,
	}
}

// Offset cria uma nova posição deslocada
func (p Position) Offset(dx, dy float64) Position {
	return FromGlobal(
		p.FastGlobalX()+dx,
		p.FastGlobalY()+dy,
		p.Z,
	)
}

// WithOffset cria uma nova posição mantendo o grid atual e alterando o offset
func (p Position) WithOffset(ox, oy float64) Position {
	return Position{
		GridX:   p.GridX,
		GridY:   p.GridY,
		OffsetX: ox,
		OffsetY: oy,
		Z:       p.Z,
	}
}

// AddOffset soma ao offset atual (pode mudar de grid se ultrapassar os limites)
func (p Position) AddOffset(dx, dy float64) Position {
	return FromGlobal(
		p.FastGlobalX()+dx,
		p.FastGlobalY()+dy,
		p.Z,
	)
}

// Equals compara duas posições com margem mínima de erro
func (p Position) Equals(other Position) bool {
	const epsilon = 0.001
	return math.Abs(p.FastGlobalX()-other.FastGlobalX()) < epsilon &&
		math.Abs(p.FastGlobalY()-other.FastGlobalY()) < epsilon &&
		math.Abs(p.Z-other.Z) < epsilon
}

// CalculateDistance calcula distância 3D entre duas posições
func CalculateDistance(a, b Position) float64 {
	dx := a.FastGlobalX() - b.FastGlobalX()
	dy := a.FastGlobalY() - b.FastGlobalY()
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// CalculateDistance2D calcula distância 2D entre duas posições
func CalculateDistance2D(a, b Position) float64 {
	dx := a.FastGlobalX() - b.FastGlobalX()
	dy := a.FastGlobalY() - b.FastGlobalY()
	return math.Sqrt(dx*dx + dy*dy)
}

// RandomWithinRadius gera nova posição aleatória dentro de um raio
func (p Position) RandomWithinRadius(radius float64) Position {
	minDist := math.Min(radius*0.6, radius-0.1)
	dist := minDist + rand.Float64()*(radius-minDist)
	angle := rand.Float64() * 2 * math.Pi
	newX := p.FastGlobalX() + dist*math.Cos(angle)
	newY := p.FastGlobalY() + dist*math.Sin(angle)
	return FromGlobal(newX, newY, p.Z)
}
