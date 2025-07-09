package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math"
// )

// type Vertex struct {
// 	X float64 `json:"x"`
// 	Z float64 `json:"z"`
// }

// type Polygon struct {
// 	ID        int      `json:"id"`
// 	Vertices  []Vertex `json:"vertices"`
// 	Neighbors []int    `json:"neighbors"`
// 	AreaType  string   `json:"areaType"`
// 	Slope     float64  `json:"slope"`
// 	Y         float64  `json:"y"`
// }

// type NavMesh struct {
// 	Polygons []Polygon `json:"polygons"`
// }

// func main() {
// 	inputFile := "map.json"
// 	outputFile := "map_subdivided.json"
// 	cellSize := 0.25

// 	data, err := ioutil.ReadFile(inputFile)
// 	if err != nil {
// 		log.Fatalf("Erro ao ler arquivo: %v", err)
// 	}

// 	var mesh NavMesh
// 	if err := json.Unmarshal(data, &mesh); err != nil {
// 		log.Fatalf("Erro ao decodificar JSON: %v", err)
// 	}

// 	var newPolygons []Polygon
// 	idCounter := 1

// 	for _, poly := range mesh.Polygons {
// 		minX, maxX, minZ, maxZ := bounds(poly.Vertices)
// 		for x := minX; x < maxX; x += cellSize {
// 			for z := minZ; z < maxZ; z += cellSize {
// 				cellVerts := []Vertex{
// 					{X: x, Z: z},
// 					{X: x + cellSize, Z: z},
// 					{X: x + cellSize, Z: z + cellSize},
// 					{X: x, Z: z + cellSize},
// 				}
// 				center := Vertex{X: x + cellSize/2, Z: z + cellSize/2}
// 				if pointInPolygon(center, poly.Vertices) {
// 					newPolygons = append(newPolygons, Polygon{
// 						ID:        idCounter,
// 						Vertices:  cellVerts,
// 						AreaType:  poly.AreaType,
// 						Slope:     poly.Slope,
// 						Y:         poly.Y,
// 						Neighbors: []int{},
// 					})
// 					idCounter++
// 				}
// 			}
// 		}
// 	}

// 	// Calcula neighbors
// 	for i := range newPolygons {
// 		for j := range newPolygons {
// 			if i == j {
// 				continue
// 			}
// 			if sharesEdge(newPolygons[i], newPolygons[j]) {
// 				newPolygons[i].Neighbors = append(newPolygons[i].Neighbors, newPolygons[j].ID)
// 			}
// 		}
// 	}

// 	outMesh := NavMesh{Polygons: newPolygons}
// 	outData, _ := json.MarshalIndent(outMesh, "", "  ")
// 	if err := ioutil.WriteFile(outputFile, outData, 0644); err != nil {
// 		log.Fatalf("Erro ao salvar: %v", err)
// 	}

// 	fmt.Printf("NavMesh subdividido salvo em: %s (%d polígonos)\n", outputFile, len(newPolygons))
// }

// func bounds(verts []Vertex) (minX, maxX, minZ, maxZ float64) {
// 	minX, minZ = math.MaxFloat64, math.MaxFloat64
// 	maxX, maxZ = -math.MaxFloat64, -math.MaxFloat64
// 	for _, v := range verts {
// 		if v.X < minX {
// 			minX = v.X
// 		}
// 		if v.X > maxX {
// 			maxX = v.X
// 		}
// 		if v.Z < minZ {
// 			minZ = v.Z
// 		}
// 		if v.Z > maxZ {
// 			maxZ = v.Z
// 		}
// 	}
// 	return
// }

// func pointInPolygon(p Vertex, verts []Vertex) bool {
// 	count := 0
// 	n := len(verts)
// 	for i := 0; i < n; i++ {
// 		v1 := verts[i]
// 		v2 := verts[(i+1)%n]
// 		if ((v1.Z > p.Z) != (v2.Z > p.Z)) &&
// 			(p.X < (v2.X-v1.X)*(p.Z-v1.Z)/(v2.Z-v1.Z)+v1.X) {
// 			count++
// 		}
// 	}
// 	return count%2 == 1
// }

// func sharesEdge(a, b Polygon) bool {
// 	matches := 0
// 	const epsilon = 0.0001
// 	for _, va := range a.Vertices {
// 		for _, vb := range b.Vertices {
// 			if math.Abs(va.X-vb.X) < epsilon && math.Abs(va.Z-vb.Z) < epsilon {
// 				matches++
// 			}
// 		}
// 	}
// 	return matches >= 2
// }
