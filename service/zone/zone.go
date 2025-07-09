package zone

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	libcreature "github.com/lunajones/apeiron/lib/creature"
	"github.com/lunajones/apeiron/lib/handle/lookup"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/world/spawn"
)

type Zone struct {
	ID           string
	Creatures    []*creature.Creature
	NavMesh      *navmesh.NavMesh
	SpatialIndex navmesh.SpatialIndex
}

var Zones []*Zone
var Players []*player.Player
var creatureCounter int

func Init() {
	zone := &Zone{ID: "old_china"}

	if err := zone.LoadNavMesh(); err != nil {
		log.Printf("[ZONE] erro ao carregar NavMesh da zona %s: %v", zone.ID, err)
	}

	zone.SpatialIndex = navmesh.NewSimpleSpatialIndex()

	if err := zone.LoadStaticSpawns(); err != nil {
		log.Printf("[ZONE] erro ao carregar spawns da zona %s: %v", zone.ID, err)
	}

	Zones = append(Zones, zone)
}

func (z *Zone) LoadNavMesh() error {
	filePath := filepath.Join("data", "zone", z.ID, "map.json")
	mesh := navmesh.LoadNavMesh(filePath)
	if mesh == nil {
		return fmt.Errorf("falha ao carregar NavMesh em %s", filePath)
	}
	z.NavMesh = mesh
	log.Printf("[ZONE] NavMesh carregado para zona %s", z.ID)
	return nil
}

func (z *Zone) Tick(elapsed float64) {
	// Cria o contexto uma vez só
	svcCtx := dynamic_context.NewAIServiceContext(z.NavMesh, z.SpatialIndex)

	for _, c := range z.Creatures {
		if !c.Alive {
			z.processRespawn(c)
			continue
		}

		if c.MoveCtrl.TargetHandle.IsValid() {
			if tgt := lookup.FindByHandle(c.MoveCtrl.TargetHandle, z.Creatures, Players); tgt != nil {
				c.MoveCtrl.UpdateTargetPosition(tgt.GetPosition())
			}
		}

		// Atualiza os alvos visíveis no contexto antes de usar
		svcCtx.CacheFor(c.Handle, c.Position, c.DetectionRadius)
		c.Tick(svcCtx, elapsed)
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
	c.Respawn(z.NavMesh)
	z.SpatialIndex.Insert(c)
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
		// Se UseNavMeshCenter for true, calcula o centro do navmesh como base
		basePos := def.Position
		if def.UseNavMeshCenter {
			bMinX, bMaxX, bMinZ, bMaxZ := z.NavMesh.BoundingBox()
			basePos = position.Position{
				X: (bMinX + bMaxX) / 2,
				Y: 0, // A altura pode ser ajustada conforme necessário
				Z: (bMinZ + bMaxZ) / 2,
			}
		}

		for i := 0; i < def.Count; i++ {
			var spawnPos position.Position
			valid := false
			for attempt := 0; attempt < 10; attempt++ {
				spawnPos = basePos.RandomWithinRadius(def.Radius)
				if z.NavMesh.IsWalkable(spawnPos) {
					valid = true
					break
				}
			}
			if !valid {
				log.Printf("[ZONE] Falha ao encontrar posição válida no NavMesh para templateID %d após 10 tentativas", def.TemplateID)
				continue
			}

			newCreature := libcreature.CreateFromTemplate(def.TemplateID, def.Position, def.Radius)
			if newCreature == nil {
				log.Printf("[ZONE] falha ao criar criatura (templateID %d)", def.TemplateID)
				continue
			}

			newCreature.SetPosition(spawnPos)
			z.AddCreature(newCreature)
			z.SpatialIndex.Insert(newCreature)

			log.Printf("[ZONE] Criado %s (%s) em (%.2f, %.2f, %.2f)",
				newCreature.Name, newCreature.PrimaryType,
				spawnPos.X, spawnPos.Y, spawnPos.Z,
			)
		}
	}

	return nil
}

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

type CreatureSnapshot struct {
	X, Z float64
	Type string
}
