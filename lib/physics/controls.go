package physics

import (
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/world/spatial"
)

// StaggerData controla o estado de stagger (atordoamento temporário)
type StaggerData struct {
	IsStaggered bool
	DurationSec float64
	EndTimeUnix int64
}

// InvincibilityData controla o estado de invencibilidade temporária
type InvincibilityData struct {
	IsInvincible bool
	DurationSec  float64
	EndTimeUnix  int64
}

// StartStagger inicia o estado de stagger
func StartStagger(data *StaggerData, now int64, durationSec float64) {
	data.IsStaggered = true
	data.DurationSec = durationSec
	data.EndTimeUnix = now + int64(durationSec)
}

// TickStagger atualiza o estado de stagger e retorna true se acabou
func TickStagger(data *StaggerData, now int64) bool {
	if data.IsStaggered && now >= data.EndTimeUnix {
		data.IsStaggered = false
		return true
	}
	return false
}

// StartInvincibility inicia o estado de invencibilidade
func StartInvincibility(data *InvincibilityData, durationSec float64) {
	now := time.Now().Unix()
	data.IsInvincible = true
	data.DurationSec = durationSec
	data.EndTimeUnix = now + int64(durationSec)
}

// TickInvincibility atualiza o estado de invencibilidade e retorna true se acabou
func TickInvincibility(data *InvincibilityData, now int64) bool {
	if data.IsInvincible && now >= data.EndTimeUnix {
		data.IsInvincible = false
		return true
	}
	return false
}

// CalculateKnockback calcula o deslocamento inicial do knockback
// No novo padrão, preferimos gerar uma força via ApplyForce em vez de teleport
func CalculateKnockback(pos, from position.Position, force float64) position.Vector3D {
	dir := position.NewVector3DFromTo(from, pos).Normalize()
	return dir.Scale(force)
}

// CheckCollision verifica colisão com o grid
func CheckCollision(newPos position.Position, radius float64) bool {
	nearby := spatial.GlobalGrid.GetNearby(newPos, radius*2)
	for _, e := range nearby {
		if e.GetHandle().ID == "" || !e.CheckIsAlive() {
			continue
		}
		dist := position.CalculateDistance(newPos, e.GetPosition())
		if dist < radius+e.GetHitboxRadius() {
			return true
		}
	}
	// TODO: Integrar colisão com cenário, layers e tipos de objeto no futuro (padrão AAA)
	return false
}
