package world

import (
	"time"

	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/zone"
)

var Players []*player.Player

func TickAll() {
	for {
		for _, z := range zone.Zones {
			ctx := dynamic_context.AIServiceContext{
				Creatures: z.Creatures,
				Players:   Players,
			}
			z.Tick(ctx)
		}
		time.Sleep(1 * time.Second)
	}
}
