package zone

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/lunajones/apeiron/service/world/spawn"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/old_china/mob"
	"github.com/lunajones/apeiron/lib/creature"
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

func (z *Zone) LoadStaticSpawns() error {
	filePath := filepath.Join("data", "zones", fmt.Sprintf("%s_spawns.json", z.Name))

	spawnDefs, err := spawn.LoadSpawnsForZone(filePath)
	if err != nil {
		return fmt.Errorf("erro carregando spawns da zona %s: %v", z.Name, err)
	}

	for _, def := range spawnDefs {
		for i := 0; i < def.Count; i++ {
			newCreature := creature.CreateFromTemplate(def.TemplateID)
			if newCreature == nil {
				log.Printf("[Zone] Falha ao criar criatura de templateID %d", def.TemplateID)
				continue
			}

			spawnPos := def.Position.RandomWithinRadius(def.Radius)
			newCreature.Position = spawnPos

			z.AddCreature(newCreature)
			log.Printf("[Zone] Criado NPC %s na zona %s", newCreature.ID, z.Name)
		}
	}

	return nil
}
