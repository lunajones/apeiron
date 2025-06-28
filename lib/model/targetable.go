package model

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type Targetable interface {
	GetPosition() position.Position
	GetHitboxRadius() float64
	GetDesiredBufferDistance() float64
	GetHandle() handle.EntityHandle

	CheckIsAlive() bool
	TakeDamage(amount int)
	ApplyEffect(effect consts.ActiveEffect)
	GetFacingDirection() position.Vector2D
	IsCreature() bool
}
