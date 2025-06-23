package world

import (
	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/zone"
)

var Players []*player.Player

func TickAll() {
	for _, z := range zone.Zones {
		ctx := core.AIContext{
			Creatures: z.Creatures,
			Players:   Players,
		}
		z.Tick(ctx)
	}
}