package zone

import (
	"fmt"
	"log"
	"github.com/lunajones/apeiron/service/ai/core"
	
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)


type Zone struct {
	ID        string
	Creatures []*creature.Creature
}

var Zones []*Zone
var creatureCounter int

func Init() {
	log.Println("[Zone] initializing zones...")

	zone1 := &Zone{ID: "zone_map1"}

	c1 := creature.NewChineseSoldier()
	c1.Position = position.Position{X: 0, Y: 0, Z: 0}

	c2 := creature.NewChineseSoldier()
	c2.Position = position.Position{X: 2, Y: 0, Z: 0}

	c3 := creature.NewChineseWolf()
	c3.Position = position.Position{X: 0, Y: 0, Z: 0}

	c4 := creature.NewChineseWolf()
	c4.Position = position.Position{X: 2, Y: 0, Z: 0}

	zone1.Creatures = append(zone1.Creatures, c1, c2, c3, c4)

	Zones = append(Zones, zone1)
	
	log.Println("[Zone] finishing zones...")

}

func (z *Zone) Tick(ctx core.AIContext) {
	for _, c := range z.Creatures {
		if c.IsAlive && c.BehaviorTree != nil {
			c.BehaviorTree.Tick(c, ctx)
		}
	}
}

type BehaviorNode interface {
	Tick(c *creature.Creature) interface{}
}

func generateUniqueCreatureID() string {
	creatureCounter++
	return fmt.Sprintf("creature_%d", creatureCounter)
}

func (z *Zone) SpawnCreature(cType creature.CreatureType, players []*player.Player, tree creature.BehaviorTree) {
	c := &creature.Creature{
		ID:      generateUniqueCreatureID(),
		Type:    cType,
		IsAlive: true,
		Position: position.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}

	c.BehaviorTree = tree

	z.Creatures = append(z.Creatures, c)
	log.Printf("[ZoneService] Criada criatura %s do tipo %s na zona %s", c.ID, cType, z.ID)
}



func convertToAIPlayers(players []*player.Player) []player.Player {
	var aiPlayers []player.Player
	for _, p := range players {
		aiPlayers = append(aiPlayers, player.Player{
			ID:       p.ID,
			Position: p.Position,
		})
	}
	return aiPlayers
}