package world

import (
	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/zone"
)

var Players []*player.Player

func TickAll() {
	for _, z := range zone.Zones {
		for _, c := range z.Creatures {
			if c.IsAlive {
				ai.ProcessAI(c, z.Creatures, Players)
				c.TickEffects()
				c.TickPosture()
			}
		}
	}
}
