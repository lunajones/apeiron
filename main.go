package main

import (
	"log"

	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/world"
	"github.com/lunajones/apeiron/service/zone"
)

func main() {
	log.Println("[Main] initializing system...")

	// Inicializa habilidades
	combat.InitSkills()

	// Inicializa jogadores (pode vir vazio por enquanto)
	world.Players = []*player.Player{}

	// Inicializa zonas e spawns
	zone.Init()

	// Inicia loop de AI
	world.TickAll()
}
