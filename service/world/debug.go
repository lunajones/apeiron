package world

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/service/creature"
)

// PrintWorldGridAAA imprime um grid 100x100 com cores e eixos no estilo AAA
func PrintWorldGridAAA(creatures []*creature.Creature) {
	const gridSize = 33
	const halfGrid = gridSize / 2

	grid := make([][]rune, gridSize)
	for i := range grid {
		grid[i] = make([]rune, gridSize)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	var outOfBounds []*creature.Creature
	rabbitCounter := 0

	for _, c := range creatures {
		x := c.Position.FastGlobalX()
		y := c.Position.FastGlobalY()

		gridX := int(math.Round(x)) + halfGrid
		gridY := int(math.Round(y)) + halfGrid

		if gridX >= 0 && gridX < gridSize && gridY >= 0 && gridY < gridSize {
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
			if grid[gridY][gridX] != '.' {
				grid[gridY][gridX] = '*'
			} else {
				grid[gridY][gridX] = symbol
			}
		} else {
			outOfBounds = append(outOfBounds, c)
		}
	}

	fmt.Println()
	color.New(color.FgHiCyan).Printf("MAPA 100x100 VISUAL DEBUG AAA\n")
	color.New(color.FgHiCyan).Printf("Y / X: -50 .. 0 .. +50\n\n")

	for i := gridSize - 1; i >= 0; i-- {
		yVal := i - halfGrid
		color.New(color.FgHiCyan).Printf("%3d| ", yVal)
		for j := 0; j < gridSize; j++ {
			val := grid[i][j]
			switch {
			case val == 'W':
				color.New(color.FgRed).Printf("%c ", val)
			case val >= '0' && val <= '9':
				color.New(color.FgGreen).Printf("%c ", val)
			case val == '*':
				color.New(color.FgMagenta).Printf("%c ", val)
			case val == '.':
				color.New(color.FgHiBlack).Printf("%c ", val)
			default:
				color.New(color.FgYellow).Printf("%c ", val)
			}
		}
		fmt.Println()
	}

	color.New(color.FgHiCyan).Printf("     ")
	for j := 0; j < gridSize; j++ {
		if (j-halfGrid)%10 == 0 {
			xVal := j - halfGrid
			color.New(color.FgHiCyan).Printf("%3d", xVal)
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
	color.New(color.FgHiBlack).Print(".")
	fmt.Print(" = Empty  ")
	color.New(color.FgYellow).Print("?")
	fmt.Print(" = Unknown\n")

	if len(outOfBounds) > 0 {
		color.New(color.FgHiYellow).Println("\n⚠ CRIATURAS FORA DO GRID VISUAL:")
		for _, c := range outOfBounds {
			fmt.Printf(" - %s (%s) em (%.2f, %.2f, %.2f)\n",
				c.Handle.String(), c.PrimaryType,
				c.Position.FastGlobalX(),
				c.Position.FastGlobalY(),
				c.Position.Z)
		}
		fmt.Println()
	}
}
