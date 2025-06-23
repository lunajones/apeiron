package main

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

func main() {
	log.Println("[Server] Iniciando servidor de zona...")

	creature.Init()
	combat.InitSkills()
	ai.InitBehaviorTrees()
	zone.Init()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		world.TickAll()
	}
}
