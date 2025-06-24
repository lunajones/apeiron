package model

import (
	"time"
	"github.com/lunajones/apeiron/lib/position"
)

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

func (c *Creature) GetID() string {
	return c.ID
}

func (c *Creature) GetPosition() position.Position {
	return c.Position
}