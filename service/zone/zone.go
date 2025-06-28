package zone

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	libcreature "github.com/lunajones/apeiron/lib/creature"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/world/spatial"
	"github.com/lunajones/apeiron/service/world/spawn"
)

type Zone struct {
	ID        string
	Creatures []*creature.Creature
}

var Zones []*Zone
var creatureCounter int

// Init inicializa todas as zonas do jogo.
func Init() {
	zone := &Zone{ID: "old_china"}

	if err := zone.LoadStaticSpawns(); err != nil {
		log.Printf("[ZONE] erro ao carregar spawns da zona %s: %v", zone.ID, err)
	}

	Zones = append(Zones, zone)
}

// Tick processa AI, respawn e ações de criaturas vivas.
func (z *Zone) Tick(ctx interface{}) {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[ZONE] contexto inválido fornecido à zona %s (recebido: %T)", z.ID, ctx)
		return
	}

	for _, c := range z.Creatures {
		if !c.IsAlive {
			z.processRespawn(c)
			continue
		}

		if c.BehaviorTree != nil {
			c.BehaviorTree.Tick(c, svcCtx)
		}
	}
}

func (z *Zone) processRespawn(c *creature.Creature) {
	if c.TimeOfDeath.IsZero() {
		return
	}
	if time.Since(c.TimeOfDeath) < time.Duration(c.RespawnTimeSec)*time.Second {
		return
	}

	log.Printf("[ZONE] Respawnando %s (%s)", c.Name, c.PrimaryType)
	c.Respawn()
	spatial.GlobalGrid.UpdateEntity(c)
}

func generateUniqueCreatureID() string {
	creatureCounter++
	return fmt.Sprintf("creature_%d", creatureCounter)
}

// LoadStaticSpawns carrega todas as criaturas estáticas da zona a partir do JSON.
func (z *Zone) LoadStaticSpawns() error {
	filePath := filepath.Join("data", "zone", z.ID, "spawns.json")

	spawnDefs, err := spawn.LoadSpawnsForZone(filePath)
	if err != nil {
		return fmt.Errorf("erro carregando spawns da zona %s: %v", z.ID, err)
	}

	for _, def := range spawnDefs {
		for i := 0; i < def.Count; i++ {
			newCreature := libcreature.CreateFromTemplate(def.TemplateID, def.Position, def.Radius)
			if newCreature == nil {
				log.Printf("[ZONE] falha ao criar criatura (templateID %d)", def.TemplateID)
				continue
			}

			spawnPos := def.Position.RandomWithinRadius(def.Radius)
			newCreature.SetPosition(spawnPos)

			log.Printf("[ZONE] Criado %s (%s) em (%.2f, %.2f, %.2f)",
				newCreature.Name, newCreature.PrimaryType,
				spawnPos.FastGlobalX(), spawnPos.FastGlobalY(), spawnPos.FastGlobalZ(),
			)

			z.AddCreature(newCreature)
			spatial.GlobalGrid.Add(newCreature)
		}
	}

	return nil
}

// AddCreature adiciona a criatura à zona, se ainda não estiver presente.
func (z *Zone) AddCreature(c *creature.Creature) {
	for _, existing := range z.Creatures {
		if existing.Handle.ID == c.Handle.ID {
			log.Printf("[ZONE] Ignorando criatura duplicada: %s (%s)", c.Handle.String(), c.Name)
			return
		}
	}
	z.Creatures = append(z.Creatures, c)
	log.Printf("[ZONE] Adicionado à zona %s: %s (%s)", z.ID, c.Handle.String(), c.Name)
}
