package zone

import (
	"fmt"
	"log"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
)

type Player struct {
	ID       string
	Position creature.Position
}

type Zone struct {
	ID        string
	Creatures []*creature.Creature
	Players   []*Player
}

var zones []*Zone

func Init() {
	log.Println("[ZoneService] Inicializando zonas...")

	zone1 := &Zone{ID: "zone_map1"}
	zone1.SpawnCreature(creature.Soldier)
	zones = append(zones, zone1)
}

func (z *Zone) SpawnCreature(cType creature.CreatureType) {
	c := &creature.Creature{
		ID:      generateUniqueCreatureID(),
		Type:    cType,
		IsAlive: true,
		Position: creature.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
	
	playerList := convertToAIPlayers(z.Players)

	// Criar BehaviorTree com a lista atual de players e creatures na zona
	switch cType {
	case creature.Soldier:
		c.BehaviorTree = mob.BuildChineseSoldierBT(playerList, z.Creatures)
	case creature.ChineseSpearman:
		c.BehaviorTree = mob.BuildChineseSpearmanBT(playerList, z.Creatures)
	default:
		log.Printf("[ZoneService] Tipo de criatura %s sem BehaviorTree definida", cType)
	}

	z.Creatures = append(z.Creatures, c)
	log.Printf("[ZoneService] Criada criatura %s do tipo %s na zona %s", c.ID, cType, z.ID)
}

func TickAllZones() {
	for _, z := range zones {
		z.Tick()
	}
}

func (z *Zone) Tick() {
	for _, c := range z.Creatures {
		c.Tick()
		ai.ProcessAI(c)
	}
}

// Helpers

var creatureCounter int

func generateUniqueCreatureID() string {
	creatureCounter++
	return fmt.Sprintf("creature_%d", creatureCounter)
}

func convertToAIPlayers(players []*Player) []ai.Player {
	var aiPlayers []ai.Player
	for _, p := range players {
		aiPlayers = append(aiPlayers, ai.Player{
			ID:       p.ID,
			Position: p.Position,
		})
	}
	return aiPlayers
}
