package main

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/combat"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/zone"
)

func main() {
	log.Println("[Server] Iniciando servidor de zona...")

	creature.Init()
	combat.InitSkills()
	ai.Init()
	zone.Init()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		zone.TickAll()
	}
}
