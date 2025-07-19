package movement

import (
	"log"
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
	// Path + Navegação
	CurrentPath    []position.Position
	PathIndex      int
	TargetPosition position.Position
	StopDistance   float64

	// Movimento ativo
	IsMoving         bool
	Velocity         position.Vector3D
	Acceleration     position.Vector3D
	DesiredDirection position.Vector2D

	// Intenção futura
	Intent            MoveIntent
	CurrentIntentDest position.Position

	// Impulso (esquiva)
	ImpulseState *ImpulseMovementState

	// Controle de bloqueios
	triedSidestep  bool
	WasBlocked     bool
	LastRepath     time.Time
	RepathCooldown time.Duration

	// Planejamento
	MovementPlan *MovementPlan

	// Alvo
	TargetHandle handle.EntityHandle
	Speed        float64
	LastUpdate   time.Time
}

func NewMovementController() *MovementController {
	return &MovementController{
		RepathCooldown: 1 * time.Second,
	}
}

func (m *MovementController) SetMoveTarget(pos position.Position, speed, stopDist float64) {
	m.Intent = MoveIntent{
		TargetPosition: pos,
		Speed:          speed,
		StopDistance:   stopDist,
		HasIntent:      true,
	}

	m.TargetHandle = handle.EntityHandle{}
	m.TargetPosition = pos
	m.Speed = speed
	m.StopDistance = stopDist
	m.CurrentIntentDest = pos

	m.CurrentPath = nil
	m.PathIndex = 0
	m.IsMoving = true
	m.triedSidestep = false
	m.WasBlocked = false
}

func (m *MovementController) UpdateTargetPosition(pos position.Position) {
	m.TargetPosition = pos
}

func (m *MovementController) Update(mov model.Movable, deltaTime float64, ctx *dynamic_context.AIServiceContext) bool {
	handleStr := mov.GetHandle().String()

	// Limpa claim anterior
	ctx.ClearClaims(mov.GetHandle())

	// 🌟 Impulso (dodge, etc)
	if m.updateImpulseMovement(mov) {
		// log.Printf("[MOVE-UPDATE] [%s] movimento via impulso (esquiva ou similar)", handleStr)
		return true
	}

	// 🌟 Intent de movimento normal (setado por AI)
	if m.applyMoveIntent(mov, ctx) {
		log.Printf("[MOVE-UPDATE] [%s] applyMoveIntent aplicou novo destino", handleStr)
		return false
	}

	if !m.IsMoving {
		// log.Printf("[MOVE-UPDATE] [%s] IsMoving = false. Encerrando update", handleStr)
		m.triedSidestep = false
		return false
	}

	dest := m.getCurrentDestination()
	// log.Printf("[MOVE-UPDATE] [%s] destino atual: (%.2f, %.2f, %.2f)", handleStr, dest.X, dest.Y, dest.Z)

	// Claim da posição de destino
	if !ctx.ClaimPosition(dest, mov.GetHandle()) {
		// log.Printf("[MOVE-UPDATE] [%s] célula (%s) já ocupada. Abortando tick.", handleStr, dest.Key())
		m.IsMoving = false
		return false
	}

	if m.checkProximity(mov, dest) {
		// log.Printf("[MOVE-UPDATE] [%s] já está próximo do destino, encerrando", handleStr)
		return true
	}

	dir := m.calculateDirection(mov, dest)

	// ✅ Atualiza a direção do torso baseada no movimento
	dir2D := position.Vector2D{X: dir.X, Z: dir.Z}.Normalize()
	if dir2D.Length() > 0.01 {
		mov.SetTorsoDirection(dir2D)
	}

	m.DesiredDirection = dir2D
	m.Acceleration = dir.Scale(m.Speed)

	// log.Printf("[MOVE-UPDATE] [%s] direção: (%.2f, %.2f), aceleração: (%.2f, %.2f), speed=%.2f",
	// handleStr, dir.X, dir.Z, m.Acceleration.X, m.Acceleration.Z, m.Speed)

	nearby := ctx.SpatialIndex.Query(mov.GetPosition(), m.calculateSearchRadius(mov))
	// log.Printf("[MOVE-UPDATE] [%s] checando colisões com %d entidades próximas", handleStr, len(nearby))

	m.WasBlocked = physics.ApplyPhysics(mov, &m.Velocity, m.Acceleration, deltaTime, true, ctx.NavMesh, nearby)
	// log.Printf("[MOVE-UPDATE] [%s] ApplyPhysics => WasBlocked = %v", handleStr, m.WasBlocked)

	if m.WasBlocked {
		// log.Printf("[MOVE-UPDATE] [%s] movimento bloqueado. Lidando com bloqueio", handleStr)
		m.handleBlockedMovement(mov, ctx)
	}

	m.LastUpdate = time.Now()
	// log.Printf("[MOVE-UPDATE] [%s] movimento aplicado com sucesso", handleStr)
	return false
}

func (m *MovementController) updateImpulseMovement(mov model.Movable) bool {
	if m.ImpulseState == nil || !m.ImpulseState.Active {
		return false
	}
	now := time.Now()
	elapsed := now.Sub(m.ImpulseState.Start)
	t := float64(elapsed) / float64(m.ImpulseState.Duration)

	if t >= 1.0 {
		mov.SetPosition(m.ImpulseState.EndPos)

		// Ajusta a face final com base em último trecho real
		dir := position.CalculateDirection2D(m.ImpulseState.StartPos, m.ImpulseState.EndPos)
		if dir.Length() > 0.01 {
			mov.SetTorsoDirection(dir)
		}

		m.ImpulseState = nil
		return true
	}

	// 🟡 Atualiza posição intermediária
	newPos := m.ImpulseState.StartPos.LerpTo(m.ImpulseState.EndPos, t)
	oldPos := mov.GetPosition()
	mov.SetPosition(newPos)

	// 🟡 Ajusta a face dinâmica: de onde estava → onde chegou agora
	dir := position.CalculateDirection2D(oldPos, newPos)
	if dir.Length() > 0.01 {
		mov.SetTorsoDirection(dir)
	}

	return true
}

func (m *MovementController) applyMoveIntent(mov model.Movable, ctx *dynamic_context.AIServiceContext) bool {
	if !m.Intent.HasIntent {
		return false
	}
	// log.Printf("[MOVE-CTRL] [%s] HasIntent = true. Dest: (%.2f, %.2f)", mov.GetHandle().String(), m.Intent.TargetPosition.X, m.Intent.TargetPosition.Z)

	m.SetMoveTarget(m.Intent.TargetPosition, m.Intent.Speed, m.Intent.StopDistance)
	m.Intent.HasIntent = false

	// log.Printf("[MOVE-CTRL] [%s] SetTarget para (%.2f, %.2f) com speed=%.2f, stop=%.2f",
	// mov.GetHandle().String(), m.TargetPosition.X, m.TargetPosition.Z, m.Speed, m.StopDistance)

	path := ctx.NavMesh.FindPath(mov.GetPosition(), m.TargetPosition)
	m.LastRepath = time.Now()
	if len(path) > 0 {
		// log.Printf("[MOVE-CTRL] [%s] Caminho encontrado com %d pontos", mov.GetHandle().String(), len(path))
		m.SetPath(path, mov)
	} else {
		// log.Printf("[MOVE-CTRL] [%s] Caminho NÃO encontrado para o alvo!", mov.GetHandle().String())
		m.IsMoving = true
	}
	return true
}

func (m *MovementController) getCurrentDestination() position.Position {
	if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath) {
		return m.CurrentPath[m.PathIndex]
	}
	return m.TargetPosition
}

func (m *MovementController) checkProximity(mov model.Movable, dest position.Position) bool {
	currentPos := mov.GetPosition()
	dx := dest.X - currentPos.X
	dy := dest.Y - currentPos.Y
	dz := dest.Z - currentPos.Z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

	if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath)-1 {
		if dist <= m.StopDistance*0.5 {
			m.PathIndex++
			m.triedSidestep = false
			return true
		}
	} else if dist <= m.StopDistance {
		if dist > m.StopDistance*0.5 {
			if position.CalculateDistance2D(m.CurrentIntentDest, m.TargetPosition) > 0.01 {
				m.SetMoveTarget(m.TargetPosition, m.Speed, m.StopDistance)
				return true // novo destino setado, continua movendo
			}
		}

		// Só para se realmente chegou e não tem mais destino pendente
		m.IsMoving = false
		m.triedSidestep = false
		return true
	}

	return false
}

func (m *MovementController) calculateDirection(mov model.Movable, dest position.Position) position.Vector3D {
	currentPos := mov.GetPosition()
	dx := dest.X - currentPos.X
	dy := dest.Y - currentPos.Y
	dz := dest.Z - currentPos.Z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return position.Vector3D{X: dx / dist, Y: dy / dist, Z: dz / dist}
}

func (m *MovementController) calculateSearchRadius(mov model.Movable) float64 {
	ownRadius := mov.GetHitboxRadius()
	maxTargetRadius := 2.0
	buffer := 0.5
	return ownRadius + maxTargetRadius + buffer
}

func (m *MovementController) handleBlockedMovement(mov model.Movable, ctx *dynamic_context.AIServiceContext) {
	// log.Printf("[MOVE-CTRL] [%s] Movimento BLOQUEADO! Tentando sidestep? %v", mov.GetHandle().String(), !m.triedSidestep)

	if !m.triedSidestep {
		angle := rand.Float64() * math.Pi
		sideX := math.Cos(angle)
		sideZ := math.Sin(angle)
		sideStep := position.Vector3D{X: sideX, Y: 0, Z: sideZ}.Scale(1.0)
		newPos := mov.GetPosition().AddOffset(sideStep.X, sideStep.Z)

		if position.CalculateDistance2D(m.CurrentIntentDest, newPos) > 0.01 {
			// log.Printf("[MOVE-CTRL] [%s] Executando sidestep para (%.2f, %.2f)", mov.GetHandle().String(), newPos.X, newPos.Z)
			m.SetMoveTarget(newPos, m.Speed, m.StopDistance)
		}
		m.triedSidestep = true
		return
	}

	if time.Since(m.LastRepath) >= m.RepathCooldown {
		// log.Printf("[MOVE-CTRL] [%s] Fazendo REPATH para alvo (%.2f, %.2f)", mov.GetHandle().String(), m.TargetPosition.X, m.TargetPosition.Z)
		path := ctx.NavMesh.FindPath(mov.GetPosition(), m.TargetPosition)
		m.LastRepath = time.Now()
		if len(path) > 0 {
			// log.Printf("[MOVE-CTRL] [%s] Novo caminho com %d pontos", mov.GetHandle().String(), len(path))
			m.SetPath(path, mov)
		} else {
			// log.Printf("[MOVE-CTRL] [%s] Repath FALHOU. Parando movimento.", mov.GetHandle().String())
			m.IsMoving = false
		}
		m.triedSidestep = false
	}
}

type MovementPlan struct {
	Type               consts.MovementPlanType
	TargetHandle       handle.EntityHandle
	DesiredDistance    float64
	ExpiresAt          time.Time
	LastTargetPosition position.Position
}

func NewMovementPlan(
	planType constslib.MovementPlanType,
	target handle.EntityHandle,
	distance float64,
	duration time.Duration,
	lastTargetPos position.Position,
) *MovementPlan {
	return &MovementPlan{
		Type:               planType,
		TargetHandle:       target,
		DesiredDistance:    distance,
		ExpiresAt:          time.Now().Add(duration),
		LastTargetPosition: lastTargetPos,
	}
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
	m.ImpulseState = nil
	m.MovementPlan = nil
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
