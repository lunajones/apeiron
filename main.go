package main

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/world"
	"github.com/lunajones/apeiron/service/zone"
)

func main() {
	log.Println("[Main] Inicializando sistema...")

	// Inicializa as skills
	combat.InitSkills()

	// Inicializa as zonas e criaturas
	zone.Init()

	// Inicializa as árvores de comportamento para as criaturas da primeira zona
	if len(zone.Zones) > 0 {
		ai.InitBehaviorTrees(world.Players, zone.Zones[0].Creatures)
	}

	// Loop de Tick a cada 50ms (~20 ticks por segundo)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		world.TickAll()
	}
}
