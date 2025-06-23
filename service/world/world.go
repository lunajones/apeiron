package world

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/zone"
	"github.com/lunajones/apeiron/service/ai"
)

func FindCreatureByID(id string) *creature.Creature {
	for _, z := range zone.Zones {
		for _, c := range z.Creatures {
			if c.ID == id {
				return c
			}
		}
	}
	return nil
}

func TickAll() {
	for _, z := range zone.Zones {
		ai.TickZone(z)
	}
}