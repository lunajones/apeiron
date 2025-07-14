package movement

import (
	"math"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/consts"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
)

type MoveIntent struct {
	TargetPosition position.Position
	Speed          float64
	StopDistance   float64
	HasIntent      bool
}

type ImpulseMovementState struct {
	Active   bool
	Start    time.Time
	Duration time.Duration
	StartPos position.Position
	EndPos   position.Position
}

type MovementController struct {
	TargetHandle      handle.EntityHandle
	TargetPosition    position.Position
	CurrentPath       []position.Position
	PathIndex         int
	Speed             float64
	StopDistance      float64
	IsMoving          bool
	Velocity          position.Vector3D
	Acceleration      position.Vector3D
	DesiredDirection  position.Vector2D
	RepathCooldown    time.Duration
	LastRepath        time.Time
	LastUpdate        time.Time
	Intent            MoveIntent
	triedSidestep     bool
	WasBlocked        bool
	CurrentIntentDest position.Position // 🌟 novo campo
	ImpulseState      *ImpulseMovementState
	MovementPlan      *MovementPlan
}

func NewMovementController() *MovementController {
	return &MovementController{
		RepathCooldown: 1 * time.Second,
	}
}

func (m *MovementController) SetMoveIntent(pos position.Position, speed, stopDist float64) {
	m.Intent = MoveIntent{
		TargetPosition: pos,
		Speed:          speed,
		StopDistance:   stopDist,
		HasIntent:      true,
	}
	m.CurrentIntentDest = pos

}

func (m *MovementController) UpdateTargetPosition(pos position.Position) {
	m.TargetPosition = pos
}

func (m *MovementController) Update(mov model.Movable, deltaTime float64, ctx *dynamic_context.AIServiceContext) bool {

	// 🌟 Execução de impulso lateral (dodge lateral)
	if m.ImpulseState != nil && m.ImpulseState.Active {
		now := time.Now()
		elapsed := now.Sub(m.ImpulseState.Start)
		t := float64(elapsed) / float64(m.ImpulseState.Duration)

		if t >= 1.0 {
			mov.SetPosition(m.ImpulseState.EndPos)
			m.ImpulseState = nil
			return true
		}

		newPos := m.ImpulseState.StartPos.LerpTo(m.ImpulseState.EndPos, t)
		mov.SetPosition(newPos)
		return true
	}

	if m.Intent.HasIntent {
		m.SetTarget(m.Intent.TargetPosition, m.Intent.Speed, m.Intent.StopDistance)
		m.Intent.HasIntent = false

		path := ctx.NavMesh.FindPath(
			mov.GetPosition(),
			m.TargetPosition,
		)
		m.LastRepath = time.Now()
		if len(path) > 0 {
			m.SetPath(path, mov)
		} else {
			m.IsMoving = true
		}
	}

	if !m.IsMoving {
		m.triedSidestep = false
		return false
	}

	var dest position.Position
	if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath) {
		dest = m.CurrentPath[m.PathIndex]
	} else {
		dest = m.TargetPosition
	}

	currentPos := mov.GetPosition()
	dx := dest.X - currentPos.X
	dz := dest.Z - currentPos.Z
	dy := dest.Y - currentPos.Y
	dist := math.Sqrt(dx*dx + dz*dz + dy*dy)

	if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath)-1 {
		if dist <= m.StopDistance*0.5 {
			m.PathIndex++
			m.triedSidestep = false
			return false
		}
	} else if dist <= m.StopDistance {
		m.IsMoving = false
		m.triedSidestep = false

		if dist > m.StopDistance*0.5 {
			// Evita intents redundantes
			if position.CalculateDistance2D(m.CurrentIntentDest, m.TargetPosition) > 0.01 {
				m.SetMoveIntent(m.TargetPosition, m.Speed, m.StopDistance)
			}
		}

		return true
	}

	dir := position.Vector3D{
		X: dx / dist,
		Y: dy / dist,
		Z: dz / dist,
	}
	m.DesiredDirection = position.Vector2D{X: dir.X, Z: dir.Z}
	m.Acceleration = dir.Scale(m.Speed)

	ownRadius := mov.GetHitboxRadius()
	maxTargetRadius := 2.0
	buffer := 0.5
	searchRadius := ownRadius + maxTargetRadius + buffer

	nearby := ctx.SpatialIndex.Query(mov.GetPosition(), searchRadius)

	blocked := physics.ApplyPhysics(mov, &m.Velocity, m.Acceleration, deltaTime, true, ctx.NavMesh, nearby)
	m.WasBlocked = blocked

	if blocked {
		if !m.triedSidestep {
			angle := rand.Float64() * math.Pi
			sideX := math.Cos(angle)
			sideZ := math.Sin(angle)
			sideStep := position.Vector3D{X: sideX, Y: 0, Z: sideZ}.Scale(1.0)

			newPos := currentPos.AddOffset(sideStep.X, sideStep.Z)

			// Evita sidestep redundante
			if position.CalculateDistance2D(m.CurrentIntentDest, newPos) > 0.01 {
				m.SetMoveIntent(newPos, m.Speed, m.StopDistance)
			}

			m.triedSidestep = true
			return false
		}

		if time.Since(m.LastRepath) >= m.RepathCooldown {
			path := ctx.NavMesh.FindPath(
				mov.GetPosition(),
				m.TargetPosition,
			)
			m.LastRepath = time.Now()
			if len(path) > 0 {
				m.SetPath(path, mov)
			} else {
				m.IsMoving = false
			}
			m.triedSidestep = false
		}
	}

	m.LastUpdate = time.Now()
	return false
}

type MovementPlan struct {
	Type            consts.MovementPlanType
	TargetHandle    handle.EntityHandle
	DesiredDistance float64
	ExpiresAt       time.Time
}

func NewMovementPlan(planType constslib.MovementPlanType, target handle.EntityHandle, distance float64, duration time.Duration) *MovementPlan {
	return &MovementPlan{
		Type:            planType,
		TargetHandle:    target,
		DesiredDistance: distance,
		ExpiresAt:       time.Now().Add(duration),
	}
}

func (m *MovementController) SetTarget(pos position.Position, speed, stopDist float64) {
	m.TargetHandle = handle.EntityHandle{}
	m.TargetPosition = pos
	m.Speed = speed
	m.StopDistance = stopDist
	m.CurrentPath = nil
	m.PathIndex = 0
	m.IsMoving = true
	m.triedSidestep = false
	m.WasBlocked = false
}

func (m *MovementController) SetPath(path []position.Position, mov model.Movable) {
	for len(path) > 0 {
		dist := position.CalculateDistance(mov.GetPosition(), path[0])
		if dist < 0.01 {
			// log.Printf("[MOVE CTRL] [%s] Ignorando ponto redundante no path (dist=%.4f)", mov.GetHandle().ID, dist)
			path = path[1:]
		} else {
			break
		}
	}

	m.CurrentPath = path
	m.PathIndex = 0
	m.IsMoving = len(path) > 0
	m.triedSidestep = false
	m.WasBlocked = false
}

func (m *MovementController) Stop() {
	m.IsMoving = false
	m.CurrentPath = nil
	m.PathIndex = 0
	m.Intent.HasIntent = false
}

func (m *MovementController) SetImpulseMovement(current position.Position, dest position.Position, duration time.Duration) {
	m.ImpulseState = &ImpulseMovementState{
		Active:   true,
		Start:    time.Now(),
		Duration: duration,
		StartPos: current,
		EndPos:   dest,
	}
	m.IsMoving = false
	m.CurrentPath = nil
	m.PathIndex = 0
}

func (p *MovementPlan) IsActive() bool {
	return p != nil && time.Now().Before(p.ExpiresAt)
}

func (p *MovementPlan) Is(planType consts.MovementPlanType) bool {
	return p != nil && p.Type == planType && p.IsActive()
}
