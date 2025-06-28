package spatial

import (
	"math"
	"sync"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
)

var GlobalGrid = NewSpatialGrid(10.0) // tamanho da célula pode ser ajustado

type SpatialEntity interface {
	GetHitboxRadius() float64
	GetHandle() handle.EntityHandle
	GetPosition() position.Position
	GetLastPosition() position.Position
	CheckIsAlive() bool
}

type cellKey struct {
	X, Z int
}

type SpatialGrid struct {
	cellSize float64
	entities map[cellKey][]SpatialEntity
	locks    map[cellKey]*sync.RWMutex
	mu       sync.RWMutex
}

func NewSpatialGrid(cellSize float64) *SpatialGrid {
	return &SpatialGrid{
		cellSize: cellSize,
		entities: make(map[cellKey][]SpatialEntity),
		locks:    make(map[cellKey]*sync.RWMutex),
	}
}

func (g *SpatialGrid) getCellKey(pos position.Position) cellKey {
	return cellKey{
		X: int(pos.FastGlobalX() / g.cellSize),
		Z: int(pos.Z / g.cellSize),
	}
}

func (g *SpatialGrid) Add(e SpatialEntity) {
	key := g.getCellKey(e.GetPosition())

	g.mu.Lock()
	lock, exists := g.locks[key]
	if !exists {
		lock = &sync.RWMutex{}
		g.locks[key] = lock
	}
	g.mu.Unlock()

	lock.Lock()
	defer lock.Unlock()

	for _, ent := range g.entities[key] {
		if ent.GetHandle().Equals(e.GetHandle()) {
			return
		}
	}

	g.entities[key] = append(g.entities[key], e)
}

func (g *SpatialGrid) Remove(e SpatialEntity) {
	key := g.getCellKey(e.GetPosition())

	g.mu.RLock()
	lock, exists := g.locks[key]
	g.mu.RUnlock()
	if !exists {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	entities := g.entities[key]
	for i, ent := range entities {
		if ent.GetHandle().Equals(e.GetHandle()) {
			g.entities[key] = append(entities[:i], entities[i+1:]...)
			break
		}
	}
}

func (g *SpatialGrid) Update(e SpatialEntity, oldPos, newPos position.Position) {
	oldKey := g.getCellKey(oldPos)
	newKey := g.getCellKey(newPos)
	if oldKey == newKey {
		return
	}
	g.Remove(e)
	g.Add(e)
}

func (g *SpatialGrid) UpdateEntity(e SpatialEntity) {
	oldKey := g.getCellKey(e.GetLastPosition())
	newKey := g.getCellKey(e.GetPosition())

	if oldKey == newKey {
		return
	}

	g.mu.RLock()
	oldLock, exists := g.locks[oldKey]
	g.mu.RUnlock()
	if exists {
		oldLock.Lock()
		entities := g.entities[oldKey]
		for i, ent := range entities {
			if ent.GetHandle().Equals(e.GetHandle()) {
				g.entities[oldKey] = append(entities[:i], entities[i+1:]...)
				break
			}
		}
		oldLock.Unlock()
	}

	g.Add(e)
}

func (g *SpatialGrid) GetNearby(pos position.Position, radius float64) []SpatialEntity {
	nearby := []SpatialEntity{}
	centerX := pos.FastGlobalX()
	centerZ := pos.Z
	minX := int((centerX - radius) / g.cellSize)
	maxX := int((centerX + radius) / g.cellSize)
	minZ := int((centerZ - radius) / g.cellSize)
	maxZ := int((centerZ + radius) / g.cellSize)

	for x := minX; x <= maxX; x++ {
		for z := minZ; z <= maxZ; z++ {
			key := cellKey{X: x, Z: z}

			g.mu.RLock()
			lock, exists := g.locks[key]
			g.mu.RUnlock()
			if !exists {
				continue
			}

			lock.RLock()
			entities := g.entities[key]
			for _, e := range entities {
				if !e.CheckIsAlive() {
					continue
				}
				entX := e.GetPosition().FastGlobalX()
				entZ := e.GetPosition().Z
				dist := math.Hypot(entX-centerX, entZ-centerZ)
				if dist <= radius {
					nearby = append(nearby, e)
				}
			}
			lock.RUnlock()
		}
	}

	return nearby
}

func (g *SpatialGrid) GetNearbyIncludingDead(pos position.Position, radius float64) []SpatialEntity {
	nearby := []SpatialEntity{}
	centerX := pos.FastGlobalX()
	centerZ := pos.Z
	minX := int((centerX - radius) / g.cellSize)
	maxX := int((centerX + radius) / g.cellSize)
	minZ := int((centerZ - radius) / g.cellSize)
	maxZ := int((centerZ + radius) / g.cellSize)

	for x := minX; x <= maxX; x++ {
		for z := minZ; z <= maxZ; z++ {
			key := cellKey{X: x, Z: z}

			g.mu.RLock()
			lock, exists := g.locks[key]
			g.mu.RUnlock()
			if !exists {
				continue
			}

			lock.RLock()
			for _, e := range g.entities[key] {
				entX := e.GetPosition().FastGlobalX()
				entZ := e.GetPosition().Z
				dist := math.Hypot(entX-centerX, entZ-centerZ)
				if dist <= radius {
					nearby = append(nearby, e)
				}
			}
			lock.RUnlock()
		}
	}

	return nearby
}

func (g *SpatialGrid) ExportGridForPathfinding(gridWidth, gridHeight int, cellSize float64) [][]int {
	grid := make([][]int, gridHeight)
	for y := 0; y < gridHeight; y++ {
		grid[y] = make([]int, gridWidth)
	}

	// Marca bloqueios
	for key, entities := range g.entities {
		// Calcula centro da célula
		centerX := float64(key.X)*g.cellSize + g.cellSize/2
		centerZ := float64(key.Z)*g.cellSize + g.cellSize/2

		// Converte para índice no grid exportado
		gridX := int(centerX / cellSize)
		gridY := int(centerZ / cellSize)

		// Verifica bounds
		if gridY < 0 || gridY >= gridHeight || gridX < 0 || gridX >= gridWidth {
			continue
		}

		for _, e := range entities {
			if !e.CheckIsAlive() {
				continue
			}
			// Marca como bloqueado
			grid[gridY][gridX] = 1
		}
	}

	// TODO: Adicione aqui marcações de obstáculos fixos (rochas, paredes) se tiver

	return grid
}
