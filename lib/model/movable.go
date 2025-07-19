package model

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
)

type Movable interface {
	GetPosition() position.Position
	SetPosition(position.Position)
	SetFacingDirection(position.Vector2D)
	SetTorsoDirection(position.Vector2D)
	GetHitboxRadius() float64
	GetHandle() handle.EntityHandle // ⚡ Necessário para logs, grid, physics
}
