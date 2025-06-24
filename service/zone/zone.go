package zone

import (
	"fmt"
	"log"
	"time"
	"path/filepath"

	"github.com/lunajones/apeiron/service/world/spawn"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	creaturelib "github.com/lunajones/apeiron/lib/creature"
)

type Zone struct {
	ID        string
	Creatures []*creature.Creature

}

var Zones []*Zone
var creatureCounter int

func Init() {
	log.Println("[Zone] initializing zones...")

	zone1 := &Zone{ID: "old_china"}

	err := zone1.LoadStaticSpawns()
	if err != nil {
		log.Printf("[Zone] Erro ao carregar spawns da zona: %v", err)
	}

	Zones = append(Zones, zone1)

	log.Println("[Zone] finishing zones...")
}

func (z *Zone) Tick(ctx core.AIContext) {
	for _, c := range z.Creatures {
		if !c.IsAlive {
			if c.TimeOfDeath.IsZero() {
				continue
			}
			if time.Since(c.TimeOfDeath) >= time.Duration(c.RespawnTimeSec)*time.Second {
				log.Printf("[Zone] Respawnando criatura %s", c.ID)
				c.Respawn()
			}
			continue
		}

		if c.BehaviorTree != nil {
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
	filePath := filepath.Join("data", "zone", z.ID, "spawns.json")

	spawnDefs, err := spawn.LoadSpawnsForZone(filePath)
	if err != nil {
		return fmt.Errorf("erro carregando spawns da zona %s: %v", z.ID, err)
	}

	for _, def := range spawnDefs {
		for i := 0; i < def.Count; i++ {
			newCreature := creaturelib.CreateFromTemplate(def.TemplateID)
			if newCreature == nil {
				log.Printf("[Zone] Falha ao criar criatura de templateID %d", def.TemplateID)
				continue
			}

			spawnPos := def.Position.RandomWithinRadius(def.Radius)
			newCreature.Position = spawnPos

			z.AddCreature(newCreature)
			log.Printf("[Zone] Criado NPC %s na zona %s", newCreature.ID, z.ID)
		}
	}

	return nil
}


func (z *Zone) AddCreature(c *creature.Creature) {
	z.Creatures = append(z.Creatures, c)
}