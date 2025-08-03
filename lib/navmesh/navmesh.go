package navmesh

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/lunajones/apeiron/lib/position"
)

type Polygon struct {
	ID        int
	GridX     int
	GridZ     int
	OffsetX   float64
	OffsetZ   float64
	Y         float64
	Slope     float64
	AreaType  string
	Neighbors []int
}

type NavMesh struct {
	Polygons    []*Polygon
	PolygonMap  map[string]*Polygon
	QuadTree    *QuadTreeNode // NOVO: estrutura de subdivisão espacial
	BoundsMin   position.Position
	BoundsMax   position.Position
	CellSize    float64
	CellOriginX float64
	CellOriginZ float64
}

func LoadFromTSV(path string) *NavMesh {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("[NAVMESH LOADER] erro ao abrir arquivo: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("[NAVMESH LOADER] erro ao ler cabeçalho: %v", err)
	}

	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	mesh := &NavMesh{
		PolygonMap: make(map[string]*Polygon),
		CellSize:   1.0,
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[NAVMESH LOADER] erro ao ler linha: %v", err)
		}

		id, _ := strconv.Atoi(record[colIdx["id"]])
		gridX, _ := strconv.Atoi(record[colIdx["gridX"]])
		gridZ, _ := strconv.Atoi(record[colIdx["gridZ"]])
		offsetX, _ := strconv.ParseFloat(record[colIdx["offsetX"]], 64)
		offsetZ, _ := strconv.ParseFloat(record[colIdx["offsetZ"]], 64)
		y, _ := strconv.ParseFloat(record[colIdx["y"]], 64)
		slope, _ := strconv.ParseFloat(record[colIdx["slope"]], 64)
		areaType := record[colIdx["areaType"]]

		neighborsStr := strings.Trim(record[colIdx["neighbors"]], "[]")
		neighbors := []int{}
		if neighborsStr != "" {
			for _, s := range strings.Split(neighborsStr, ",") {
				n, _ := strconv.Atoi(strings.TrimSpace(s))
				neighbors = append(neighbors, n)
			}
		}

		p := &Polygon{
			ID:        id,
			GridX:     gridX,
			GridZ:     gridZ,
			OffsetX:   offsetX,
			OffsetZ:   offsetZ,
			Y:         y,
			AreaType:  areaType,
			Neighbors: neighbors,
			Slope:     slope,
		}

		mesh.Polygons = append(mesh.Polygons, p)

		key := fmt.Sprintf("%d,%d", gridX, gridZ)
		mesh.PolygonMap[key] = p
	}

	mesh.calculateBounds()

	// NOVO: constrói o Quadtree após calcular bounds
	bounds := position.AABB{
		Min: mesh.BoundsMin,
		Max: mesh.BoundsMax,
	}
	mesh.QuadTree = BuildQuadTree(mesh.Polygons, bounds, 0, 6, 16)

	log.Printf("[NAVMESH LOADER] NavMesh carregado com %d polígonos", len(mesh.Polygons))
	return mesh
}

func (m *NavMesh) calculateBounds() {
	if len(m.Polygons) == 0 {
		return
	}

	minX := m.Polygons[0].OffsetX
	minZ := m.Polygons[0].OffsetZ
	maxX := m.Polygons[0].OffsetX
	maxZ := m.Polygons[0].OffsetZ

	for _, p := range m.Polygons {
		x := float64(p.GridX) + p.OffsetX
		z := float64(p.GridZ) + p.OffsetZ

		if x < minX {
			minX = x
		}
		if z < minZ {
			minZ = z
		}
		if x > maxX {
			maxX = x
		}
		if z > maxZ {
			maxZ = z
		}
	}

	m.BoundsMin = position.Position{X: minX, Z: minZ}
	m.BoundsMax = position.Position{X: maxX, Z: maxZ}
	m.CellOriginX = minX
	m.CellOriginZ = minZ
}

func (m *NavMesh) IsWalkable(pos position.Position) bool {
	gridX := int(math.Floor(pos.X))
	gridZ := int(math.Floor(pos.Z))
	key := fmt.Sprintf("%d,%d", gridX, gridZ)
	// for k := range m.PolygonMap {
	// 	log.Printf("[POLYGONMAP DEBUG] chave carregada: %s", k)
	// }

	log.Printf("[WALKABLE DEBUG] buscando chave: %s", key)

	poly, exists := m.PolygonMap[key]
	if !exists {
		log.Printf("[WALKABLE DEBUG] (%d, %d) fora do grid", gridX, gridZ)
		return false
	}

	log.Printf("[WALKABLE DEBUG] (%d, %d) ok — área: %s", gridX, gridZ, poly.AreaType)
	return true
}

func (p *Polygon) GetVertices() []position.Position {
	baseX := float64(p.GridX)
	baseZ := float64(p.GridZ)
	return []position.Position{
		{X: baseX, Y: p.Y, Z: baseZ},
		{X: baseX + 1, Y: p.Y, Z: baseZ},
		{X: baseX + 1, Y: p.Y, Z: baseZ + 1},
		{X: baseX, Y: p.Y, Z: baseZ + 1},
	}
}

func (p *Polygon) GetCenterPosition() position.Position {
	return position.Position{
		X: float64(p.GridX) + p.OffsetX,
		Y: p.Y,
		Z: float64(p.GridZ) + p.OffsetZ,
	}
}

func (mesh *NavMesh) GetEscapePoint(current position.Position, threats []position.Position, maxDist float64) position.Position {
	bestPos := current
	maxScore := -1.0

	for _, poly := range mesh.Polygons {
		center := poly.GetCenterPosition()

		if position.CalculateDistance(current, center) > maxDist {
			continue
		}

		score := 0.0
		for _, threat := range threats {
			score += position.CalculateDistance(center, threat)
		}

		if score > maxScore {
			maxScore = score
			bestPos = center
		}
	}

	return bestPos
}

func (mesh *NavMesh) GetRandomWalkablePoint(origin position.Position, minDist, maxDist float64) position.Position {
	candidates := []position.Position{}

	for _, poly := range mesh.Polygons {
		center := poly.GetCenterPosition()
		dist := position.CalculateDistance2D(origin, center)
		if dist >= minDist && dist <= maxDist {
			candidates = append(candidates, center)
		}
	}

	if len(candidates) == 0 {
		return origin
	}

	return candidates[rand.Intn(len(candidates))]
}

func (m *NavMesh) FindClosestPolygonByGrid(pos position.Position) *Polygon {
	gridX := int(math.Floor(pos.X))
	gridZ := int(math.Floor(pos.Z))
	key := fmt.Sprintf("%d,%d", gridX, gridZ)

	if poly, exists := m.PolygonMap[key]; exists {
		return poly
	}

	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			neighborKey := fmt.Sprintf("%d_%d", gridX+dx, gridZ+dz)
			if poly, exists := m.PolygonMap[neighborKey]; exists {
				return poly
			}
		}
	}

	return nil
}
