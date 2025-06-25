package model

import (
	"github.com/lunajones/apeiron/lib/position"
)

type Creature struct {
	ID               string
	Name             string
	MaxHP            int
	SpawnPoint       position.Position
	SpawnRadius      float64
	RespawnTimeSec   int
	OwnerPlayerID    string
	Faction          string
}

func (c *Creature) GetID() string {
	return c.ID
}
