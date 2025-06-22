package world

import (
	"github.com/lunajones/apeiron/service/creature"
)

var Zones []*Zone

type Zone struct {
	ID        string
	Creatures []*creature.Creature
}

func FindCreatureByID(id string) *creature.Creature {
	for _, z := range Zones {
		for _, c := range z.Creatures {
			if c.ID == id {
				return c
			}
		}
	}
	return nil
}
