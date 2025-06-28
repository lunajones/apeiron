package pathfinding

import (
	"container/heap"
	"fmt"
	"log"
	"math"

	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
)

// Node representa um ponto no grid para o A*
type Node struct {
	Pos    position.Position
	GCost  float64
	HCost  float64
	FCost  float64
	Parent *Node
	Index  int // Para o heap
}

// NodeHeap é o heap do open set
type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].FCost < h[j].FCost }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i]; h[i].Index = i; h[j].Index = j }

func (h *NodeHeap) Push(x interface{}) {
	n := x.(*Node)
	n.Index = len(*h)
	*h = append(*h, n)
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := old[len(old)-1]
	*h = old[0 : len(old)-1]
	return n
}

// FindPath executa A* e retorna o caminho
func FindPath(start, end position.Position, grid [][]int) []position.Position {
	openSet := &NodeHeap{}
	openSetMap := make(map[string]*Node)
	heap.Init(openSet)

	startNode := &Node{Pos: start, GCost: 0, HCost: heuristic(start, end)}
	startNode.FCost = startNode.HCost
	heap.Push(openSet, startNode)
	openSetMap[key(start)] = startNode

	costSoFar := map[string]float64{
		key(start): 0,
	}

	dirs := []position.Vector2D{
		{X: 1, Y: 0}, {X: -1, Y: 0},
		{X: 0, Y: 1}, {X: 0, Y: -1},
		{X: 1, Y: 1}, {X: -1, Y: -1},
		{X: -1, Y: 1}, {X: 1, Y: -1},
	}

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)
		delete(openSetMap, key(current.Pos))

		if current.Pos.Equals(end) {
			return reconstructPath(current)
		}

		for _, dir := range dirs {
			nextPos := current.Pos.Offset(dir.X, dir.Y)
			if !isWalkable(nextPos, grid) {
				continue
			}

			newCost := costSoFar[key(current.Pos)] + distance(current.Pos, nextPos)
			nextKey := key(nextPos)
			if c, ok := costSoFar[nextKey]; !ok || newCost < c {
				costSoFar[nextKey] = newCost
				hCost := heuristic(nextPos, end)

				if existing, ok := openSetMap[nextKey]; ok {
					if newCost < existing.GCost {
						existing.GCost = newCost
						existing.FCost = newCost + hCost
						existing.Parent = current
						heap.Fix(openSet, existing.Index)
					}
				} else {
					nextNode := &Node{
						Pos:    nextPos,
						GCost:  newCost,
						HCost:  hCost,
						Parent: current,
					}
					nextNode.FCost = nextNode.GCost + nextNode.HCost
					heap.Push(openSet, nextNode)
					openSetMap[nextKey] = nextNode
				}
			}
		}
	}

	log.Println("[PATHFINDER] Caminho não encontrado")
	return nil
}

// Gera chave única para map
func key(p position.Position) string {
	return fmt.Sprintf("%d_%d_%.2f_%.2f", p.GridX, p.GridY, p.OffsetX, p.OffsetY)
}

func heuristic(a, b position.Position) float64 {
	dx := math.Abs(a.FastGlobalX() - b.FastGlobalX())
	dy := math.Abs(a.FastGlobalY() - b.FastGlobalY())
	return dx + dy
}

func distance(a, b position.Position) float64 {
	dx := a.FastGlobalX() - b.FastGlobalX()
	dy := a.FastGlobalY() - b.FastGlobalY()
	return math.Sqrt(dx*dx + dy*dy)
}

func isWalkable(pos position.Position, grid [][]int) bool {
	x := int(pos.FastGlobalX())
	y := int(pos.FastGlobalY())
	if x < 0 || y < 0 || y >= len(grid) || x >= len(grid[0]) {
		return false
	}
	if grid[y][x] == 1 {
		return false
	}
	if physics.CheckCollision(pos, 0.5) {
		return false
	}
	return true
}

func reconstructPath(endNode *Node) []position.Position {
	var path []position.Position
	for n := endNode; n != nil; n = n.Parent {
		path = append([]position.Position{n.Pos}, path...)
	}
	return path
}
