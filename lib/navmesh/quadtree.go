package navmesh

import (
	"math"

	"github.com/lunajones/apeiron/lib/position"
)

type QuadTreeNode struct {
	Bounds     position.AABB
	Polygons   []*Polygon
	Children   [4]*QuadTreeNode
	Depth      int
	MaxDepth   int
	MaxPerNode int
}

func BuildQuadTree(polygons []*Polygon, bounds position.AABB, depth, maxDepth, maxPerNode int) *QuadTreeNode {
	node := &QuadTreeNode{
		Bounds:     bounds,
		Depth:      depth,
		MaxDepth:   maxDepth,
		MaxPerNode: maxPerNode,
	}

	for _, poly := range polygons {
		p := poly.GetCenterPosition()
		if bounds.Contains(p) {
			node.Polygons = append(node.Polygons, poly)
		}
	}

	if depth >= maxDepth || len(node.Polygons) <= maxPerNode {
		return node // folha
	}

	// subdividir em 4
	cx := (bounds.Min.X + bounds.Max.X) / 2
	cz := (bounds.Min.Z + bounds.Max.Z) / 2

	quadrants := [4]position.AABB{
		{Min: position.Position{X: bounds.Min.X, Z: bounds.Min.Z}, Max: position.Position{X: cx, Z: cz}},
		{Min: position.Position{X: cx, Z: bounds.Min.Z}, Max: position.Position{X: bounds.Max.X, Z: cz}},
		{Min: position.Position{X: bounds.Min.X, Z: cz}, Max: position.Position{X: cx, Z: bounds.Max.Z}},
		{Min: position.Position{X: cx, Z: cz}, Max: position.Position{X: bounds.Max.X, Z: bounds.Max.Z}},
	}

	for i := 0; i < 4; i++ {
		node.Children[i] = BuildQuadTree(node.Polygons, quadrants[i], depth+1, maxDepth, maxPerNode)
	}

	node.Polygons = nil // otimização: folhas apenas armazenam polígonos
	return node
}

func (n *QuadTreeNode) FindClosestPolygon(p position.Position) *Polygon {
	if n == nil {
		return nil
	}
	if len(n.Children[0].Children) == 0 { // folha
		var best *Polygon
		bestDist := math.MaxFloat64
		for _, poly := range n.Children[0].Polygons {
			d := position.CalculateDistance2D(p, poly.GetCenterPosition())
			if d < bestDist {
				bestDist = d
				best = poly
			}
		}
		return best
	}
	for _, child := range n.Children {
		if child.Bounds.Contains(p) {
			return child.FindClosestPolygon(p)
		}
	}
	// fallback: procura em todas
	for _, child := range n.Children {
		if poly := child.FindClosestPolygon(p); poly != nil {
			return poly
		}
	}
	return nil
}
