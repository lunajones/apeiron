package physics

import (
	"github.com/lunajones/apeiron/lib/position"
)

func IsWithinBounds(pos position.Position, bounds Bounds) bool {
	return pos.GridX >= bounds.Min.GridX && pos.GridX <= bounds.Max.GridX &&
		pos.GridY >= bounds.Min.GridY && pos.GridY <= bounds.Max.GridY &&
		pos.Z >= bounds.Min.Z && pos.Z <= bounds.Max.Z
}
