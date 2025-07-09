package navmesh

import (
	"container/heap"
	"log"
	"math"

	"github.com/lunajones/apeiron/lib/position"
)

type PathOption func(*PathSettings)

type PathSettings struct {
	UseFunnel      bool
	MaxSlope       float64
	AvoidAreas     []string
	PrefAreas      []string
	TargetHeight   float64
	MaxHeightDiff  float64
	CostModifiers  map[string]float64
	ConsiderHeight bool
}

func defaultPathSettings() *PathSettings {
	return &PathSettings{
		UseFunnel:      false,
		MaxSlope:       45,
		AvoidAreas:     nil,
		PrefAreas:      nil,
		TargetHeight:   0,
		MaxHeightDiff:  3,
		CostModifiers:  make(map[string]float64),
		ConsiderHeight: false,
	}
}

type PathNode struct {
	PolygonID int
	Pos       position.Position
	GCost     float64
	HCost     float64
	FCost     float64
	Parent    *PathNode
	Index     int
}

type NodeHeap []*PathNode

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].FCost < h[j].FCost }
func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}
func (h *NodeHeap) Push(x interface{}) {
	n := x.(*PathNode)
	n.Index = len(*h)
	*h = append(*h, n)
}
func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := old[len(old)-1]
	*h = old[:len(old)-1]
	return n
}

func (mesh *NavMesh) FindPath(start, end position.Position, opts ...PathOption) []position.Position {
	settings := defaultPathSettings()
	for _, opt := range opts {
		opt(settings)
	}

	startPoly := mesh.findContainingPolygon(start)
	endPoly := mesh.findContainingPolygon(end)

	if startPoly == nil || endPoly == nil {
		log.Printf("[NAVMESH PATHFINDER] Start ou End fora da NavMesh")
		return nil
	}

	openSet := &NodeHeap{}
	heap.Init(openSet)

	startNode := &PathNode{
		PolygonID: startPoly.ID,
		Pos:       start,
		GCost:     0,
		HCost:     heuristicCost(start, end, settings),
	}
	startNode.FCost = startNode.HCost

	heap.Push(openSet, startNode)
	cameFrom := make(map[int]*PathNode)
	cameFrom[startPoly.ID] = startNode

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*PathNode)

		if current.PolygonID == endPoly.ID {
			path := reconstructPath(current, end)
			if settings.UseFunnel {
				log.Printf("[NAVMESH PATHFINDER] Funnel aplicado (placeholder)")
			}
			return path
		}

		poly := mesh.getPolygonByID(current.PolygonID)
		for _, neighborID := range poly.Neighbors {
			if _, ok := cameFrom[neighborID]; ok {
				continue
			}
			neighborPoly := mesh.getPolygonByID(neighborID)

			if !isPolyAllowed(current.Pos, neighborPoly, settings) {
				continue
			}

			center := neighborPoly.centerPosition()
			gCost := current.GCost + heuristicCost(current.Pos, center, settings)

			if areaCost, ok := settings.CostModifiers[neighborPoly.AreaType]; ok {
				gCost *= areaCost
			}

			hCost := heuristicCost(center, end, settings)
			fCost := gCost + hCost

			node := &PathNode{
				PolygonID: neighborID,
				Pos:       center,
				GCost:     gCost,
				HCost:     hCost,
				FCost:     fCost,
				Parent:    current,
			}
			heap.Push(openSet, node)
			cameFrom[neighborID] = node
		}
	}

	log.Printf("[NAVMESH PATHFINDER] Nenhum caminho encontrado")
	return nil
}

func isPolyAllowed(fromPos position.Position, poly *Polygon, settings *PathSettings) bool {
	if settings.AvoidAreas != nil {
		for _, area := range settings.AvoidAreas {
			if poly.AreaType == area {
				return false
			}
		}
	}
	if settings.MaxSlope > 0 && poly.Slope > settings.MaxSlope {
		return false
	}
	if settings.ConsiderHeight {
		if math.Abs(fromPos.Y-poly.centerPosition().Y) > settings.MaxHeightDiff {
			return false
		}
	}
	return true
}

func heuristicCost(a, b position.Position, settings *PathSettings) float64 {
	dx := a.X - b.X
	dz := a.Z - b.Z
	dist := math.Sqrt(dx*dx + dz*dz)
	if settings.ConsiderHeight {
		dy := a.Y - b.Y
		dist = math.Sqrt(dist*dist + dy*dy)
	}
	return dist
}

func (mesh *NavMesh) getPolygonByID(id int) *Polygon {
	for i := range mesh.Polygons {
		if mesh.Polygons[i].ID == id {
			return &mesh.Polygons[i]
		}
	}
	return nil
}

func (mesh *NavMesh) findContainingPolygon(pos position.Position) *Polygon {
	for i := range mesh.Polygons {
		if pointInPolygonXZ(pos, mesh.Polygons[i]) {
			return &mesh.Polygons[i]
		}
	}
	return nil
}

func pointInPolygonXZ(pos position.Position, poly Polygon) bool {
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

func (poly *Polygon) centerPosition() position.Position {
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

func reconstructPath(node *PathNode, end position.Position) []position.Position {
	var path []position.Position
	for n := node; n != nil; n = n.Parent {
		path = append([]position.Position{n.Pos}, path...)
	}
	path = append(path, end)
	return path
}

func interpolatePath(path []position.Position, maxSegment float64) []position.Position {
	var refined []position.Position
	for i := 0; i < len(path)-1; i++ {
		refined = append(refined, path[i])
		a := path[i]
		b := path[i+1]
		dist := position.CalculateDistance(a, b)
		steps := int(dist / maxSegment)
		for s := 1; s < steps; s++ {
			t := float64(s) / float64(steps)
			refined = append(refined, position.Position{
				X: a.X + t*(b.X-a.X),
				Y: a.Y + t*(b.Y-a.Y),
				Z: a.Z + t*(b.Z-a.Z),
			})
		}
	}
	refined = append(refined, path[len(path)-1])
	return refined
}
