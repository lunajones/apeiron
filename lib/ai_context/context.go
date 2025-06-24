package ai_context

import (
	"github.com/lunajones/apeiron/lib/model"
)

type CreatureContext interface {
	GetCreatures() []*creature.Creature
	GetPlayers() []*player.Player
}

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
