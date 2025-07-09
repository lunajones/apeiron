package model

import (
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type Targetable interface {
	GetPosition() position.Position
	GetLastPosition() position.Position
	SetPosition(position.Position) // NOVO: Permite o push e outras alterações de posição

	GetHitboxRadius() float64
	GetDesiredBufferDistance() float64
	GetHandle() handle.EntityHandle

	IsAlive() bool
	TakeDamage(amount int)
	ApplyEffect(effect constslib.ActiveEffect)
	GetFacingDirection() position.Vector2D
	IsCreature() bool

	IsHostile() bool
	IsPvPEnabled() bool
	IsHungry() bool // se precisar
	IsBlocking() bool
	HasTag(tag consts.CreatureTag) bool
	IsInvulnerableNow() bool
	ApplyPostureDamage(amount float64)
	IsInParryWindow() bool
}
