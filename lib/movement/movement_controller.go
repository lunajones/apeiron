package movement

import (
	"log"
	"math"
	"time"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement/pathfinding"
	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
)

type MovementController struct {
	TargetHandle     handle.EntityHandle
	TargetPosition   position.Position
	CurrentPath      []position.Position
	PathIndex        int
	Speed            float64
	StopDistance     float64
	IsMoving         bool
	WasBlocked       bool
	Velocity         position.Vector3D
	Acceleration     position.Vector3D
	DesiredDirection position.Vector2D
	RepathCooldown   time.Duration
	LastRepath       time.Time
	LastUpdate       time.Time
}

func NewMovementController() *MovementController {
	return &MovementController{
		RepathCooldown: 1 * time.Second,
	}
}

func (m *MovementController) SetTarget(pos position.Position, speed, stopDist float64) {
	m.TargetPosition = pos
	m.Speed = speed
	m.StopDistance = stopDist
	m.CurrentPath = nil
	m.PathIndex = 0
	m.IsMoving = true
	m.WasBlocked = false
}

func (m *MovementController) SetPath(path []position.Position, speed, stopDist float64) {
	m.CurrentPath = path
	m.PathIndex = 0
	m.Speed = speed
	m.StopDistance = stopDist
	m.IsMoving = true
	m.WasBlocked = false
}

func (m *MovementController) Update(mov model.Movable, deltaTime float64, grid [][]int) bool {
	if !m.IsMoving {
		return false
	}

	var dest position.Position
	if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath) {
		dest = m.CurrentPath[m.PathIndex]
	} else {
		dest = m.TargetPosition
	}

	currentPos := mov.GetPosition()
	dx := dest.FastGlobalX() - currentPos.FastGlobalX()
	dy := dest.FastGlobalY() - currentPos.FastGlobalY()
	dz := dest.Z - currentPos.Z

	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dist <= m.StopDistance {
		if len(m.CurrentPath) > 0 && m.PathIndex < len(m.CurrentPath)-1 {
			m.PathIndex++
			return false
		}
		m.IsMoving = false
		log.Printf("[MOVE CTRL] [%s] chegou ao destino (%.2f ≤ %.2f)", mov.GetHandle().ID, dist, m.StopDistance)
		return true
	}

	dir := position.Vector3D{
		X: dx / dist,
		Y: dy / dist,
		Z: dz / dist,
	}
	m.DesiredDirection = position.Vector2D{X: dir.X, Y: dir.Z}
	m.Acceleration = dir.Scale(m.Speed)

	physics.ApplyPhysics(mov, &m.Velocity, m.Acceleration, deltaTime, true)

	if m.WasBlocked {
		if grid != nil {
			log.Printf("[MOVE CTRL] [%s] bloqueado, tentando pathfinding...", mov.GetHandle().ID)
			path := pathfinding.FindPath(currentPos, m.TargetPosition, grid)
			if len(path) > 0 {
				m.SetPath(path, m.Speed, m.StopDistance)
				log.Printf("[MOVE CTRL] [%s] novo path definido com %d pontos", mov.GetHandle().ID, len(path))
			} else {
				log.Printf("[MOVE CTRL] [%s] pathfinding falhou ao recalcular", mov.GetHandle().ID)
			}
		} else {
			log.Printf("[MOVE CTRL] [%s] bloqueado mas grid não fornecido, não pode recalcular path", mov.GetHandle().ID)
		}
		m.WasBlocked = false
	}

	m.LastUpdate = time.Now()
	return false
}

func (m *MovementController) WasBlockedNow() bool {
	return m.WasBlocked
}
