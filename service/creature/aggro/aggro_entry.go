package aggro

import (
	"time"
)

type AggroEntry struct {
	TargetID       string
	ThreatValue    float64
	LastDamageTime time.Time
	AggroSource    string // Ex: "Damage", "Heal", "Taunt"
	LastAction     string // Ex: Nome da skill, ataque, etc
}
