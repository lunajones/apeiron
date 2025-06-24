package player

import (
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
)

// Wrapper Player (composição)
type Player struct {
	model.Player
}

// Interface Targetable
func (p *Player) GetPosition() position.Position {
	return p.Position
}

func (p *Player) GetID() string {
	return p.ID
}