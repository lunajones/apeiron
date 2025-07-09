package main

import (
	"log"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/world"
	"github.com/lunajones/apeiron/service/zone"
)

func main() {
	log.Println("[Main] initializing system...")

	// Inicializa habilidades
	model.InitSkills()

	// Inicializa jogadores (pode vir vazio por enquanto)
	// Inicializa zonas e spawns
	zone.Init()

	// Inicia loop de AI
	world.TickAll()
}
