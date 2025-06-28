package world

import (
	"time"

	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/world/spatial"
	"github.com/lunajones/apeiron/service/zone"
)

var Players []*player.Player

func TickAll() {
	for {
		// Exporte o grid de pathfinding
		grid := spatial.GlobalGrid.ExportGridForPathfinding(50, 50, 1.0) // ajuste os valores conforme seu mapa

		svcCtx := &dynamic_context.AIServiceContext{
			PathfindingGrid: grid,
		}

		for _, z := range zone.Zones {
			z.Tick(svcCtx)

			for _, c := range z.Creatures {
				if !c.Position.Equals(c.LastPosition) {
					spatial.GlobalGrid.UpdateEntity(c)
				}
			}
		}

		var allCreatures []*creature.Creature
		for _, z := range zone.Zones {
			allCreatures = append(allCreatures, z.Creatures...)
		}
		PrintWorldGridAAA(allCreatures)

		time.Sleep(1 * time.Second)
	}
}
