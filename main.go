package main

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/world"
	"github.com/lunajones/apeiron/service/zone"
	"github.com/lunajones/apeiron/service/creature"
)

func main() {
	log.Println("[Main] initializing system...")

	combat.InitSkills()
	zone.Init()

	var allCreatures []*creature.Creature
	for _, z := range zone.Zones {
		allCreatures = append(allCreatures, z.Creatures...)
	}

	ai.InitBehaviorTrees(world.Players, allCreatures)

	// Loop de Tick a cada 50ms (~20 ticks por segundo)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		world.TickAll()
	}
}
