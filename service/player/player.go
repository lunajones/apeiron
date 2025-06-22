package player

import "github.com/lunajones/apeiron/service/creature"

type Player struct {
	ID       string
	Position creature.Position
	// No futuro: HP, MP, Atributos, Inventory, etc
}
