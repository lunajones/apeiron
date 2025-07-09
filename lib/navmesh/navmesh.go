package navmesh

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"math/rand/v2"

	"github.com/lunajones/apeiron/lib/position"
)

type Vertex struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
}

type Polygon struct {
	ID        int      `json:"id"`
	Vertices  []Vertex `json:"vertices"`
	Neighbors []int    `json:"neighbors"`
	AreaType  string   `json:"areaType"`
	Slope     float64  `json:"slope"`
	Y         float64  `json:"y"` // altura média ou referência do polígono
}

type NavMesh struct {
	Polygons []Polygon `json:"polygons"`
}

// LoadNavMesh carrega o NavMesh a partir de um arquivo JSON
func LoadNavMesh(path string) *NavMesh {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[NAVMESH LOADER] Erro ao ler arquivo: %v", err)
	}

	var mesh NavMesh
	if err := json.Unmarshal(data, &mesh); err != nil {
		log.Fatalf("[NAVMESH LOADER] Erro ao decodificar JSON: %v", err)
	}

	log.Printf("[NAVMESH LOADER] NavMesh carregado com %d polígonos", len(mesh.Polygons))
	return &mesh
}

// Verifica se um ponto está dentro do polígono no plano XZ
func PointInPolygonXZ(pos position.Position, poly Polygon) bool {
	count := 0
	n := len(poly.Vertices)
	for i := 0; i < n; i++ {
		v1 := poly.Vertices[i]
		v2 := poly.Vertices[(i+1)%n]
		if ((v1.Z > pos.Z) != (v2.Z > pos.Z)) &&
			(pos.X < (v2.X-v1.X)*(pos.Z-v1.Z)/(v2.Z-v1.Z)+v1.X) {
			count++
		}
	}
	return count%2 == 1
}

func (mesh *NavMesh) IsWalkable(pos position.Position) bool {
	for _, poly := range mesh.Polygons {
		if PointInPolygonXZ(pos, poly) {
			return true
		}
	}
	return false
}

func (mesh *NavMesh) BoundingBox() (minX, maxX, minZ, maxZ float64) {
	if len(mesh.Polygons) == 0 {
		return 0, 0, 0, 0
	}
	minX, maxX = math.MaxFloat64, -math.MaxFloat64
	minZ, maxZ = math.MaxFloat64, -math.MaxFloat64

	for _, poly := range mesh.Polygons {
		for _, v := range poly.Vertices {
			if v.X < minX {
				minX = v.X
			}
			if v.X > maxX {
				maxX = v.X
			}
			if v.Z < minZ {
				minZ = v.Z
			}
			if v.Z > maxZ {
				maxZ = v.Z
			}
		}
	}
	return
}

func (mesh *NavMesh) FindClosestPolygon(pos position.Position) *Polygon {
	var closest *Polygon
	minDist := math.MaxFloat64
	for i := range mesh.Polygons {
		center := mesh.Polygons[i].CenterPosition()
		dist := position.CalculateDistance2D(pos, center)
		if dist < minDist {
			minDist = dist
			closest = &mesh.Polygons[i]
		}
	}
	return closest
}

func (poly *Polygon) CenterPosition() position.Position {
	var sumX, sumZ float64
	for _, v := range poly.Vertices {
		sumX += v.X
		sumZ += v.Z
	}
	n := float64(len(poly.Vertices))
	return position.Position{
		X: sumX / n,
		Y: poly.Y,
		Z: sumZ / n,
	}
}

func (mesh *NavMesh) GetEscapePoint(current position.Position, threats []position.Position, distance float64) position.Position {
	var centerX, centerZ float64
	for _, t := range threats {
		centerX += t.X
		centerZ += t.Z
	}
	count := float64(len(threats))
	if count == 0 {
		return current
	}
	centerX /= count
	centerZ /= count

	dirX := current.X - centerX
	dirZ := current.Z - centerZ
	mag := math.Hypot(dirX, dirZ)
	if mag == 0 {
		angle := rand.Float64() * 2 * math.Pi
		dirX = math.Cos(angle)
		dirZ = math.Sin(angle)
		mag = 1
	}
	dirX /= mag
	dirZ /= mag

	return position.Position{
		X: current.X + dirX*distance,
		Y: current.Y,
		Z: current.Z + dirZ*distance,
	}
}

func (mesh *NavMesh) GetRandomWalkablePoint(origin position.Position, minDist, maxDist float64) position.Position {
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		distance := rand.Float64()*(maxDist-minDist) + minDist
		angle := rand.Float64() * 2 * math.Pi

		dest := position.Position{
			X: origin.X + math.Cos(angle)*distance,
			Y: origin.Y,
			Z: origin.Z + math.Sin(angle)*distance,
		}

		if mesh.IsWalkable(dest) {
			return dest
		}
	}

	// fallback: retorna o próprio ponto de origem
	return origin
}
