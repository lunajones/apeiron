package ai_context

import (
	"github.com/lunajones/apeiron/lib/model"
)

type AIContext struct {
	Creatures []*model.Creature
	Players   []*model.Player
}

func (ctx AIContext) GetCreatures() []*model.Creature {
	return ctx.Creatures
}

func (ctx AIContext) GetPlayers() []*model.Player {
	return ctx.Players
}
