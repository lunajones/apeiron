package dynamic_context

import (
	"time"

	"github.com/lunajones/apeiron/lib/handle"
)

// CombatBehaviorEvent representa uma ação ofensiva ou defensiva registrada no contexto de combate
type CombatBehaviorEvent struct {
	SourceHandle handle.EntityHandle
	TargetHandle handle.EntityHandle
	BehaviorType string        // Exemplo: "FakeAdvance", "Provoke", "GuardRaise", "AggressiveIntention"
	Timestamp    time.Time     // Momento do evento (ex: início do windup)
	Damage       float64       // Dano causado (se aplicável)
	WindupTime   time.Duration // ⬅️ Adicionado
}
