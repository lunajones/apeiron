package context

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type AIServiceContext struct {
	Creatures []*creature.Creature
	Players   []*player.Player
}

func (ctx AIServiceContext) GetServiceCreatures() []*creature.Creature {
	return ctx.Creatures
}

func (ctx AIServiceContext) GetServicePlayers() []*player.Player {
	return ctx.Players
}
