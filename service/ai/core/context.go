package core

import (
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type AIContext struct {
	Creatures []*creature.Creature
	Players   []*player.Player
}
