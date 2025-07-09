package model

import (
	"github.com/lunajones/apeiron/lib/position"
)

type Attacker interface {
	Targetable // herda tudo de Targetable: posição, vida, hitbox, etc.

	// Stats básicos
	GetStrength() int
	GetDexterity() int
	GetIntelligence() int
	GetFocus() int

	// Postura
	ApplyPostureDamage(amount float64)

	// Tipo e direção
	GetPrimaryType() string
	GetFacingDirection() position.Vector2D
	SetFacingDirection(dir position.Vector2D)

	// Movimento especial
	SetSkillMovementState(state *SkillMovementState)
	GetSkillMovementState() *SkillMovementState

	IsHostile() bool
	IsPvPEnabled() bool
	IsAlive() bool
	IsHungry() bool // se precisar
}
