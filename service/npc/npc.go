package npc

import (
	"time"
	"github.com/lunajones/apeiron/service/npc/consts"
)

type NPC struct {
	ID         string
	Name       string
	Role       consts.Role
	City       string
	IsAlive    bool
	IsRelative bool     // true se for parente de jogador
	RelatedTo  string   // ID do jogador, se for parente
	LastSeen   time.Time
}
