package zone

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
)

var ShowTrail = true
var trailGrid = make(map[string]bool)

func PrintWorldGridAAA(creatures []*creature.Creature, mesh *navmesh.NavMesh) {
	const cellStep = 1

	minX, maxX, minZ, maxZ := mesh.BoundingBox()
	margin := 2.0
	minX -= margin
	maxX += margin
	minZ -= margin
	maxZ += margin

	gridWidth := int(math.Ceil((maxX - minX) / cellStep))
	gridHeight := int(math.Ceil((maxZ - minZ) / cellStep))

	grid := make([][]rune, gridHeight)
	for i := range grid {
		grid[i] = make([]rune, gridWidth)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	rabbitCounter := 0
	var outOfBounds []*creature.Creature

	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			x := minX + float64(j)*cellStep
			z := minZ + float64(i)*cellStep
			pos := position.Position{X: x, Y: 0, Z: z}

			for _, poly := range mesh.Polygons {
				if navmesh.PointInPolygonXZ(pos, poly) {
					grid[i][j] = '#'
					break
				}
			}
		}
	}

	for _, c := range creatures {
		x := c.Position.X
		z := c.Position.Z

		gridX := int(math.Round((x - minX) / cellStep))
		gridZ := int(math.Round((z - minZ) / cellStep))

		key := fmt.Sprintf("%d_%d", gridX, gridZ)
		trailGrid[key] = true

		if gridX >= 0 && gridX < gridWidth && gridZ >= 0 && gridZ < gridHeight {
			var symbol rune
			switch c.PrimaryType {
			case "Wolf":
				symbol = 'W'
			case "Rabbit":
				symbol = rune('0' + (rabbitCounter % 10))
				rabbitCounter++
			default:
				symbol = '?'
			}
			if grid[gridZ][gridX] != '.' && grid[gridZ][gridX] != '#' {
				grid[gridZ][gridX] = '*'
			} else {
				grid[gridZ][gridX] = symbol
			}
		} else {
			outOfBounds = append(outOfBounds, c)
		}
	}

	for _, c := range creatures {
		if c.MoveCtrl != nil && len(c.MoveCtrl.CurrentPath) > 0 {
			for _, pt := range c.MoveCtrl.CurrentPath {
				gridX := int(math.Round((pt.X - minX) / cellStep))
				gridZ := int(math.Round((pt.Z - minZ) / cellStep))
				if gridX >= 0 && gridX < gridWidth && gridZ >= 0 && gridZ < gridHeight {
					if grid[gridZ][gridX] == '.' || grid[gridZ][gridX] == '#' {
						grid[gridZ][gridX] = '+'
					}
				}
			}
		}
	}

	fmt.Println()
	color.New(color.FgHiCyan).Printf("MAPA VISUAL DYNAMIC NAVMESH\n")
	color.New(color.FgHiCyan).Printf("Z/X %.1f to %.1f / %.1f to %.1f\n\n", minZ, maxZ, minX, maxX)

	for i := gridHeight - 1; i >= 0; i-- {
		zVal := minZ + float64(i)*cellStep
		color.New(color.FgHiCyan).Printf("%6.1f| ", zVal)
		for j := 0; j < gridWidth; j++ {
			val := grid[i][j]
			key := fmt.Sprintf("%d_%d", j, i)
			if val == '.' && ShowTrail && trailGrid[key] {
				val = 'o'
			}
			switch {
			case val == 'W':
				color.New(color.FgRed).Printf("%c ", val)
			case val >= '0' && val <= '9':
				color.New(color.FgGreen).Printf("%c ", val)
			case val == '*':
				color.New(color.FgMagenta).Printf("%c ", val)
			case val == 'o':
				color.New(color.FgYellow).Printf("%c ", val)
			case val == '#':
				color.New(color.FgBlue).Printf("%c ", val)
			case val == '+':
				color.New(color.FgHiYellow).Printf("%c ", val)
			case val == '.':
				color.New(color.FgHiBlack).Printf("%c ", val)
			default:
				color.New(color.FgYellow).Printf("%c ", val)
			}
		}
		fmt.Println()
	}

	color.New(color.FgHiCyan).Printf("        ")
	for j := 0; j < gridWidth; j++ {
		xVal := minX + float64(j)*cellStep
		if int(xVal)%10 == 0 {
			color.New(color.FgHiCyan).Printf("%3.0f", xVal)
		} else {
			fmt.Print("   ")
		}
	}
	fmt.Println()

	fmt.Println()
	color.New(color.FgRed).Print("W")
	fmt.Print(" = Wolf  ")
	color.New(color.FgGreen).Print("0-9")
	fmt.Print(" = Rabbits  ")
	color.New(color.FgMagenta).Print("*")
	fmt.Print(" = Overlap  ")
	color.New(color.FgYellow).Print("o")
	fmt.Print(" = Trilha  ")
	color.New(color.FgBlue).Print("#")
	fmt.Print(" = NavMesh  ")
	color.New(color.FgHiYellow).Print("+")
	fmt.Print(" = Path  ")
	color.New(color.FgHiBlack).Print(".")
	fmt.Print(" = Empty  ")
	color.New(color.FgYellow).Print("?")
	fmt.Print(" = Unknown\n")

	if len(outOfBounds) > 0 {
		color.New(color.FgHiYellow).Println("\n⚠ CRIATURAS FORA DO GRID VISUAL:")
		for _, c := range outOfBounds {
			fmt.Printf(" - %s (%s) em (%.2f, %.2f, %.2f)\n",
				c.Handle.String(), c.PrimaryType,
				c.Position.X,
				c.Position.Z,
				c.Position.Y)
		}
		fmt.Println()
	}
}
