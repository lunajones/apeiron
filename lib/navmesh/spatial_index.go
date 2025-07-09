package navmesh

import (
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
)

type SpatialIndex interface {
	Insert(entity model.Targetable)
	Remove(entity model.Targetable)
	Update(entity model.Targetable)
	Query(center position.Position, radius float64) []model.Targetable
}

type SimpleSpatialIndex struct {
	entities []model.Targetable
}

func NewSimpleSpatialIndex() *SimpleSpatialIndex {
	return &SimpleSpatialIndex{
		entities: []model.Targetable{},
	}
}

func (s *SimpleSpatialIndex) Insert(entity model.Targetable) {
	s.entities = append(s.entities, entity)
}

func (s *SimpleSpatialIndex) Remove(entity model.Targetable) {
	newList := s.entities[:0]
	for _, e := range s.entities {
		if !e.GetHandle().Equals(entity.GetHandle()) {
			newList = append(newList, e)
		}
	}
	s.entities = newList
}

func (s *SimpleSpatialIndex) Update(entity model.Targetable) {
	// No-op
}

func (s *SimpleSpatialIndex) Query(center position.Position, radius float64) []model.Targetable {
	var result []model.Targetable
	r2 := radius * radius
	for _, e := range s.entities {
		if !e.IsAlive() {
			continue
		}
		last := e.GetLastPosition()
		dx := center.X - last.X
		dz := center.Z - last.Z
		if dx*dx+dz*dz <= r2 {
			result = append(result, e)
		}
	}
	return result
}
