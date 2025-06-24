package core

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type AIContext struct {
	Creatures []*creature.Creature
	Players   []*player.Player
}

func (ctx AIContext) GetCreatures() []interface{} {
	creatures := make([]interface{}, len(ctx.Creatures))
	for i, c := range ctx.Creatures {
		creatures[i] = c
	}
	return creatures
}

func (ctx AIContext) GetPlayers() []interface{} {
	players := make([]interface{}, len(ctx.Players))
	for i, p := range ctx.Players {
		players[i] = p
	}
	return players
}
