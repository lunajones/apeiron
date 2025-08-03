package position

type AABB struct {
	Min Position
	Max Position
}

func (a AABB) Contains(p Position) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Z >= a.Min.Z && p.Z <= a.Max.Z
}
