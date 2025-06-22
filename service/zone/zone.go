package zone

import (
	"fmt"
	"log"

	"github.com/lunajones/apeiron/service/ai"
	"github.com/lunajones/apeiron/service/ai/old_china/mob"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world"
)

type Player struct {
	ID       string
	Position creature.Position
}

func Init() {
	log.Println("[ZoneService] Inicializando zonas...")

	zone1 := &world.Zone{ID: "zone_map1"}

	c1 := creature.NewChineseSoldier()
	c1.Position = creature.Position{X: 0, Y: 0, Z: 0}

	c2 := creature.NewChineseSoldier()
	c2.Position = creature.Position{X: 2, Y: 0, Z: 0}

	zone1.Creatures = append(zone1.Creatures, c1, c2)

	world.Zones = append(world.Zones, zone1)
}

func TickAll() {
	for _, z := range world.Zones {
		z.Tick()
	}
}

func (z *world.Zone) Tick() {
	for _, c := range z.Creatures {
		if c.IsAlive {
			ai.ProcessAI(c, z.Creatures)
			c.TickEffects()
			c.TickPosture()
		}
	}
}

var creatureCounter int

func generateUniqueCreatureID() string {
	creatureCounter++
	return fmt.Sprintf("creature_%d", creatureCounter)
}

func (z *world.Zone) SpawnCreature(cType creature.CreatureType, players []*Player) {
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

	playerList := convertToAIPlayers(players)

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