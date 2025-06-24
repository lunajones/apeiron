package creature

import "github.com/lunajones/apeiron/lib/position"

type Creature struct {
	ID               string
	Name             string
	HP               int
	MaxHP            int
	Position         position.Position
	SpawnPoint       position.Position
	SpawnRadius      float64
	IsAlive          bool
	IsCorpse         bool
	RespawnTimeSec   int
	TimeOfDeath      time.Time
	OwnerPlayerID    string
	TargetCreatureID string
	TargetPlayerID   string
	Faction          string
	IsHostile        bool
}