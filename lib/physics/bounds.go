package physics

import (
	"github.com/lunajones/apeiron/lib/position"
)

type Bounds struct {
	Min position.Position
	Max position.Position
}

func NewBounds(min, max position.Position) Bounds {
	return Bounds{Min: min, Max: max}
}

func (b Bounds) Contains(pos position.Position) bool {
	return IsWithinBounds(pos, b)
}
