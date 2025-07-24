package model

import (
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/position"
)

// Movable define o contrato para qualquer entidade que possa se mover no mundo,
// incluindo criaturas, jogadores, NPCs e afins.
type Movable interface {
	// Posição e orientação
	GetPosition() position.Position
	SetPosition(position.Position)
	SetFacingDirection(position.Vector2D)
	SetTorsoDirection(position.Vector2D)
	GetHitboxRadius() float64
	GetHandle() handle.EntityHandle
	GetPrimaryType() string // Para logs/debug

	ClearMovementIntent()
	IsDodging() bool

	// Estado de animação
	SetAnimationState(state constslib.AnimationState)

	// Skill movement (pulo, leap, etc)
	GetSkillMovementState() *SkillMovementState

	// FSM de movimento
	SetDodging(active bool)
	SetBlocking(active bool)
	SetCachedDodgePosition(pos position.Position)
	SetLastDodgeEvent(e CombatEvent)
	ReduceStamina(amount float64)
	GetStamina() float64
	GetDodgeStaminaCost() float64
	GetDodgeInvulnerabilityDuration() time.Duration
	SetInvulnerableUntil(t time.Time)
	SetDodgeStartedAt(t time.Time)
}
