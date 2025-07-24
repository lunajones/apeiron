package model

import (
	"time"

	"github.com/lunajones/apeiron/lib/position"
)

// SkillMovementState representa o estado de avanço de uma skill em execução (como pulo ou investida)
type SkillMovementState struct {
	Active        bool              // Se o movimento está em andamento
	StartTime     time.Time         // Momento em que o movimento começou
	Duration      time.Duration     // Duração total do movimento
	Speed         float64           // Velocidade do avanço
	Direction     position.Vector3D // Direção do movimento
	TargetPos     position.Position // Posição inicial do alvo no momento do cast (se TargetLock)
	Config        *MovementConfig   // Configuração da skill que gerou o movimento
	DamageApplied bool
	Skill         *Skill
	HitTargets    map[string]bool
	EngagedTarget Targetable // Usado para empurrar o alvo enquanto se move

}

func (s *SkillMovementState) IsComplete(now time.Time, currentPos position.Position) bool {
	if !s.Active {
		return true
	}
	// Chegou ao destino?
	if position.CalculateDistance2D(currentPos, s.TargetPos) < 0.1 {
		return true
	}
	// Duração expirou?
	if now.Sub(s.StartTime) >= s.Duration {
		return true
	}
	return false
}

func (s *SkillMovementState) HasAlreadyHit(t Targetable) bool {
	if s.HitTargets == nil {
		s.HitTargets = make(map[string]bool)
	}
	return s.HitTargets[t.GetHandle().ID]
}

func (s *SkillMovementState) MarkAsHit(t Targetable) {
	if s.HitTargets == nil {
		s.HitTargets = make(map[string]bool)
	}
	s.HitTargets[t.GetHandle().ID] = true
}
