package dynamic_context

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/world/spatial"
)

type AIServiceContext struct {
	PathfindingGrid [][]int
}

func (ctx *AIServiceContext) GetServiceCreatures(center position.Position, radius float64) []*creature.Creature {
	entities := spatial.GlobalGrid.GetNearby(center, radius)
	result := make([]*creature.Creature, 0, len(entities))
	for _, e := range entities {
		if c, ok := e.(*creature.Creature); ok && c.IsAlive {
			result = append(result, c)
		}
	}
	return result
}

func (ctx *AIServiceContext) GetServicePlayers(center position.Position, radius float64) []*player.Player {
	entities := spatial.GlobalGrid.GetNearby(center, radius)
	result := make([]*player.Player, 0, len(entities))
	for _, e := range entities {
		if p, ok := e.(*player.Player); ok {
			result = append(result, p)
		}
	}
	return result
}

func (ctx *AIServiceContext) GetServiceCorpses(center position.Position, radius float64) []*creature.Creature {
	entities := spatial.GlobalGrid.GetNearbyIncludingDead(center, radius)
	result := make([]*creature.Creature, 0, len(entities))
	for _, e := range entities {
		if c, ok := e.(*creature.Creature); ok && !c.IsAlive {
			result = append(result, c)
		}
	}
	return result
}

func (ctx *AIServiceContext) FindCreatureByHandle(h handle.EntityHandle) *creature.Creature {
	entities := spatial.GlobalGrid.GetNearbyIncludingDead(position.Position{}, 99999)
	for _, e := range entities {
		if c, ok := e.(*creature.Creature); ok {
			if c.GetHandle().Equals(h) {
				return c
			}
		}
	}
	return nil
}

func (ctx *AIServiceContext) GetPathfindingGrid() [][]int {
	return ctx.PathfindingGrid
}
