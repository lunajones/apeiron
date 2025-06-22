package creature

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/lunajones/apeiron/lib"
)

type Position struct {
	X float64
	Y float64
	Z float64
}

type Creature struct {
	ID              string
	Type            CreatureType
	Level           CreatureLevel
	HP              int
	Actions         []CreatureAction
	CurrentAction   CreatureAction
	AIState         AIState
	LastStateChange time.Time
	DynamicCombos   map[CreatureAction][]CreatureAction

	// Controle de vida e respawn
	IsAlive        bool
	RespawnTimeSec int
	TimeOfDeath    int64
	OwnerPlayerID  string // Se for de quest, quem é o dono

	// Efeitos ativos (DOTs, Buffs, Debuffs)
	ActiveEffects []ActiveEffect

	// Posição e spawn
	SpawnPoint Position
	Position   Position
	SpawnRadius float64
}

var creatures []*Creature

func Init() {
	log.Println("Creature service initialized")
	creatures = append(creatures, exampleSpawn())
}

func exampleSpawn() *Creature {
	c := &Creature{
		ID:    lib.NewUUID(),
		Type:  Mob,
		Level: Normal,
		HP:    100,
		Actions: []CreatureAction{
			ActionIdle,
			ActionWalk,
			ActionRun,
			ActionParry,
			ActionBlock,
			ActionJump,
			ActionSkill1,
			ActionSkill2,
			ActionSkill3,
			ActionSkill4,
			ActionSkill5,
			ActionCombo1,
			ActionCombo2,
			ActionCombo3,
			ActionDie,
		},
		CurrentAction:   ActionIdle,
		AIState:         AIStateIdle,
		LastStateChange: time.Now(),
		DynamicCombos:   make(map[CreatureAction][]CreatureAction),
		IsAlive:         true,
		RespawnTimeSec:  30,
		SpawnPoint:      Position{X: 0, Y: 0, Z: 0},
		SpawnRadius:     5.0,
	}

	c.Position = c.GenerateSpawnPosition()
	return c
}

func (c *Creature) GenerateSpawnPosition() Position {
	for attempts := 0; attempts < 10; attempts++ {
		offsetX := (rand.Float64()*2 - 1) * c.SpawnRadius
		offsetZ := (rand.Float64()*2 - 1) * c.SpawnRadius
		newPos := Position{
			X: c.SpawnPoint.X + offsetX,
			Y: c.SpawnPoint.Y,
			Z: c.SpawnPoint.Z + offsetZ,
		}

		if IsTerrainWalkable(newPos) {
			return newPos
		}
	}
	return c.SpawnPoint
}

func IsTerrainWalkable(pos Position) bool {
	// TODO: Integrar com o TerrainService para validação real
	return true
}

func (c *Creature) Tick() {
	c.TickEffects()

	switch c.AIState {
	case AIStateIdle:
		if rand.Float32() < 0.1 {
			c.ChangeAIState(AIStateAlert)
		}

	case AIStateAlert:
		if time.Since(c.LastStateChange) > 2*time.Second {
			c.ChangeAIState(AIStateAttack)
		}

	case AIStateAttack:
		log.Printf("[Creature %s] Atacando!", c.ID)
		c.SetAction(ActionAttack)
		c.ChangeAIState(AIStateIdle)

	case AIStateDead:
		// Nada a fazer
	}
}

func (c *Creature) SetAction(action CreatureAction) {
	c.CurrentAction = action
	log.Printf("[Creature %s] Action set to: %s", c.ID, action)
}

func (c *Creature) ChangeAIState(newState AIState) {
	log.Printf("[Creature %s] AI State mudou: %s → %s", c.ID, c.AIState, newState)
	c.AIState = newState
	c.LastStateChange = time.Now()

	switch newState {
	case AIStateIdle:
		c.SetAction(ActionIdle)
	case AIStateAlert:
		c.SetAction(ActionIdle)
	case AIStateAttack:
		c.SetAction(ActionAttack)
	case AIStateDead:
		c.SetAction(ActionDie)
		c.IsAlive = false
		c.TimeOfDeath = time.Now().Unix()
	}
}

func (c *Creature) GenerateRandomCombo(comboAction CreatureAction) {
	possibleSkills := []CreatureAction{
		ActionSkill1,
		ActionSkill2,
		ActionSkill3,
		ActionSkill4,
		ActionSkill5,
	}

	var combo []CreatureAction
	numSkillsInCombo := rand.Intn(4) + 2

	for i := 0; i < numSkillsInCombo; i++ {
		randomSkill := possibleSkills[rand.Intn(len(possibleSkills))]
		combo = append(combo, randomSkill)
	}

	if c.DynamicCombos == nil {
		c.DynamicCombos = make(map[CreatureAction][]CreatureAction)
	}

	c.DynamicCombos[comboAction] = combo
	log.Printf("[Creature %s] Novo combo gerado para %s: %v", c.ID, comboAction, combo)
}

func (c *Creature) IsQuestOnly() bool {
	return strings.TrimSpace(c.OwnerPlayerID) != ""
}

func (c *Creature) ApplyEffect(effect ActiveEffect) {
	c.ActiveEffects = append(c.ActiveEffects, effect)
	log.Printf("[Creature %s] recebeu efeito: %s", c.ID, effect.Type)
}

func (c *Creature) TickEffects() {
	now := time.Now().Unix()
	var remainingEffects []ActiveEffect

	for _, eff := range c.ActiveEffects {
		if now-eff.StartTime >= int64(eff.Duration) {
			log.Printf("[Creature %s] efeito %s expirou", c.ID, eff.Type)
			continue
		}

		if eff.Type.IsDOT() {
			damage := 10
			c.HP -= damage
			log.Printf("[Creature %s] sofreu %d de dano de %s. HP atual: %d", c.ID, damage, eff.Type, c.HP)

			if c.HP <= 0 && c.IsAlive {
				c.IsAlive = false
				c.TimeOfDeath = now
				c.CurrentAction = ActionDie
				log.Printf("[Creature %s] morreu por DOT %s", c.ID, eff.Type)
			}
		}

		remainingEffects = append(remainingEffects, eff)
	}

	c.ActiveEffects = remainingEffects
}

func TickAll() {
	for _, c := range creatures {
		c.Tick()
	}
}

func DebugPrintCreatures() {
	for _, c := range creatures {
		fmt.Printf("Creature: %s, Type: %s, Level: %s, AIState: %s, HP: %d, Action: %s, Pos: (%.2f, %.2f, %.2f)\n",
			c.ID, c.Type, c.Level, c.AIState, c.HP, c.CurrentAction, c.Position.X, c.Position.Y, c.Position.Z)
	}
}
