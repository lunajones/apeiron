package core

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type AIContext struct {
	Creatures []*creature.Creature
	Players   []*player.Player
}

func (ctx AIContext) GetCreatures() []*creature.Creature {
	return ctx.Creatures
}

func (ctx AIContext) GetPlayers() []*player.Player {
	return ctx.Players
}