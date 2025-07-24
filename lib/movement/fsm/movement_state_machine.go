package fsm

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
)

// MovementState representa os possíveis estados de movimento
type MovementState int

const (
	MovementStateIdle MovementState = iota
	MovementStateWalking
	MovementStateDodging
	MovementStateKnockback
)

// MovementFSM define a interface pública para qualquer FSM de movimento
type MovementFSM interface {
	GetState() MovementState
	Enter(state MovementState)
	Tick(deltaTime float64)
	EnterDodgingState(mov model.Movable, dest position.Position, dir position.Vector2D, e model.CombatEvent)
}

// Callbacks para ações externas
type FSMHooks struct {
	OnClearIntent      func()
	OnSetAnimation     func(state constslib.AnimationState)
	OnKnockbackImpulse func()
	OnIsCasting        func() bool // 🔹 AQUI
	OnHasArrived       func() bool
	OnIsDodging        func() bool
	OnSetImpulse       func(state *movement.ImpulseMovementState)
	OnShouldWalk       func() bool
}

// MovementStateMachine é a implementação concreta
type MovementStateMachine struct {
	state        MovementState
	enteredAt    time.Time
	stateChanged bool
	hooks        FSMHooks
}

// NewMovementStateMachine inicializa com os hooks necessários
func ProcessMovementFSM(hooks FSMHooks) *MovementStateMachine {
	return &MovementStateMachine{
		state:        MovementStateIdle,
		enteredAt:    time.Now(),
		stateChanged: true,
		hooks:        hooks,
	}
}

func (fsm *MovementStateMachine) GetState() MovementState {
	return fsm.state
}

func (fsm *MovementStateMachine) Enter(state MovementState) {
	if fsm.state == state {
		fsm.stateChanged = false
		return
	}

	fsm.state = state
	fsm.enteredAt = time.Now()
	fsm.stateChanged = true

	log.Printf("[FSM-MOVEMENT] entrou em estado %s", fsm.state.String())

	switch state {
	case MovementStateIdle:
		if fsm.hooks.OnClearIntent != nil {
			fsm.hooks.OnClearIntent()
		}
		if fsm.hooks.OnSetAnimation != nil {
			fsm.hooks.OnSetAnimation(constslib.AnimationIdle)
		}

	case MovementStateWalking:
		if fsm.hooks.OnIsCasting != nil && fsm.hooks.OnIsCasting() {
			fsm.Enter(MovementStateIdle)
			return
		}
		if fsm.hooks.OnHasArrived != nil && fsm.hooks.OnHasArrived() {
			fsm.Enter(MovementStateIdle)
		}

	case MovementStateKnockback:
		if fsm.hooks.OnSetAnimation != nil {
			fsm.hooks.OnSetAnimation(constslib.AnimationKnockback)
		}
		if fsm.hooks.OnKnockbackImpulse != nil {
			fsm.hooks.OnKnockbackImpulse()
		}
	}
}

func (fsm *MovementStateMachine) Tick(deltaTime float64) {
	switch fsm.state {
	case MovementStateDodging:
		if fsm.hooks.OnIsDodging != nil && !fsm.hooks.OnIsDodging() {
			if fsm.hooks.OnShouldWalk != nil && fsm.hooks.OnShouldWalk() {
				fsm.Enter(MovementStateWalking)
			} else {
				fsm.Enter(MovementStateIdle)
			}
		}
	case MovementStateKnockback:
		if time.Since(fsm.enteredAt) > 1*time.Second {
			fsm.Enter(MovementStateIdle)
		}

	case MovementStateWalking:
		if fsm.hooks.OnHasArrived != nil && fsm.hooks.OnHasArrived() {
			fsm.Enter(MovementStateIdle)
		}
	}
}

// String retorna o nome do estado para debug
func (s MovementState) String() string {
	switch s {
	case MovementStateIdle:
		return "IDLE"
	case MovementStateWalking:
		return "WALKING"
	case MovementStateDodging:
		return "DODGING"
	case MovementStateKnockback:
		return "KNOCKBACK"
	default:
		return "UNKNOWN"
	}
}

func (fsm *MovementStateMachine) EnterDodgingState(mov model.Movable, dest position.Position, dir position.Vector2D, e model.CombatEvent) {
	now := time.Now()

	// Ativa impulso
	if fsm.hooks.OnSetImpulse != nil {
		fsm.hooks.OnSetImpulse(&movement.ImpulseMovementState{
			Active:   true,
			Start:    now,
			Duration: 300 * time.Millisecond,
			StartPos: mov.GetPosition(),
			EndPos:   dest,
		})
	}

	mov.SetDodging(true)
	mov.SetBlocking(false)
	mov.SetAnimationState(constslib.AnimationDodge)
	mov.SetCachedDodgePosition(dest)
	mov.ReduceStamina(mov.GetDodgeStaminaCost())
	mov.SetLastDodgeEvent(e)
	mov.SetInvulnerableUntil(now.Add(mov.GetDodgeInvulnerabilityDuration()))
	mov.SetDodgeStartedAt(now)

	fsm.Enter(MovementStateDodging)
}
