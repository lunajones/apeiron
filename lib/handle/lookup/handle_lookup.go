package lookup

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

// FindByHandle procura uma criatura ou jogador pelo EntityHandle
func FindByHandle(h handle.EntityHandle, creatures []*creature.Creature, players []*player.Player) model.Targetable {
	for _, c := range creatures {
		if c.Handle.Equals(h) {
			return c
		}
	}
	for _, p := range players {
		if p.Handle.Equals(h) {
			return p
		}
	}
	return nil
}

// ValidateHandle verifica se um targetable corresponde ao handle informado
func ValidateHandle(h handle.EntityHandle, tgt model.Targetable) bool {
	if tgt == nil {
		return false
	}
	return tgt.GetHandle().Equals(h)
}
