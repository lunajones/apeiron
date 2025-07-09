package aggro

import (
	"time"

	"github.com/lunajones/apeiron/lib/handle"
)

type AggroEntry struct {
	TargetHandle   handle.EntityHandle
	ThreatValue    float64
	LastDamageTime time.Time
	AggroSource    string
	LastAction     string
}
