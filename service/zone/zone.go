package zone

import (
	"fmt"
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
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

	// Exemplo de criação de soldados e lobos
	zone1.Creatures = append(zone1.Creatures, mob.NewChineseSoldier())
	//zone1.Creatures = append(zone1.Creatures, mob.NewChineseSoldier())
	zone1.Creatures = append(zone1.Creatures, mob.NewChineseWolf())
	//zone1.Creatures = append(zone1.Creatures, mob.NewChineseWolf())

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
