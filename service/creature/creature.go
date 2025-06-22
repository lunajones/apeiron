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
	MaxHP           int
	Actions         []CreatureAction
	CurrentAction   CreatureAction
	AIState         AIState
	LastStateChange time.Time
	DynamicCombos   map[CreatureAction][]CreatureAction

	// Atributos base
	Strength     int
	Dexterity    int
	Intelligence int
	Focus        int

	// Defesas e Resistências
	PhysicalDefense    float64
	MagicDefense       float64
	RangedDefense      float64
	ControlResistance  float64
	StatusResistance   float64
	CriticalResistance float64
	CriticalChance     float64

	// Controle de vida e respawn
	IsAlive        bool
	RespawnTimeSec int
	TimeOfDeath    int64
	OwnerPlayerID  string

	// Efeitos ativos
	ActiveEffects []ActiveEffect

	// Posição e spawn
	SpawnPoint  Position
	Position    Position
	SpawnRadius float64

	// AI Targeting
	TargetCreatureID string
	TargetPlayerID   string

	// Perception
	FieldOfViewDegrees float64
	VisionRange        float64
	HearingRange       float64
	IsBlind            bool
	IsDeaf             bool

	// Perception/Combat - Range
	DetectionRadius float64
	AttackRange     float64

	// Skill Cooldowns
	SkillCooldowns map[CreatureAction]time.Time

	// Aggro / Hate
	AggroTable map[string]float64

	// Movement and Attack speed
	MoveSpeed   float64
	AttackSpeed float64

	// Faction / Hostility
	Faction   string
	IsHostile bool

	// Posture / Stagger system
	MaxPosture              float64
	CurrentPosture          float64
	PostureRegenRate        float64
	IsPostureBroken         bool
	TimePostureBroken       int64
	PostureBreakDurationSec int

	// AI Behavior Tree
	BehaviorTree interface {
		Tick(c *Creature)
	}
}

var creatures []*Creature

func Init() {
	log.Println("Creature service initialized")
	creatures = append(creatures, exampleSpawn())
}

func exampleSpawn() *Creature {
	c := &Creature{
		ID:    lib.NewUUID(),
		Type:  Soldier,
		Level: Normal,
		HP:    100,
		MaxHP: 100,
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
		CurrentAction:           ActionIdle,
		AIState:                 AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[CreatureAction][]CreatureAction),
		IsAlive:                 true,
		RespawnTimeSec:          30,
		SpawnPoint:              Position{X: 0, Y: 0, Z: 0},
		SpawnRadius:             5.0,
		FieldOfViewDegrees:      120,
		VisionRange:             15,
		HearingRange:            10,
		IsBlind:                 false,
		IsDeaf:                  false,
		DetectionRadius:         10.0,
		AttackRange:             2.5,
		SkillCooldowns:          make(map[CreatureAction]time.Time),
		AggroTable:              make(map[string]float64),
		MoveSpeed:               3.5,
		AttackSpeed:             1.2,
		Faction:                 "Monsters",
		IsHostile:               true,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		// Atributos de combate
		Strength:          20,
		Dexterity:         10,
		Intelligence:      5,
		Focus:             8,
		PhysicalDefense:   0.15,
		MagicDefense:      0.05,
		RangedDefense:     0.10,
		ControlResistance: 0.1,
		StatusResistance:  0.1,
		CriticalResistance: 0.2,
		CriticalChance:     0.05,
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
	// TODO: Integrar com TerrainService real
	return true
}

func (c *Creature) Tick() {
	c.TickEffects()
	c.TickPosture()

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

func (c *Creature) TickPosture() {
	if c.IsPostureBroken {
		if time.Now().Unix()-c.TimePostureBroken >= int64(c.PostureBreakDurationSec) {
			c.IsPostureBroken = false
			c.CurrentPosture = c.MaxPosture
			c.ChangeAIState(AIStateIdle)
			log.Printf("[Creature %s] Posture recuperada.", c.ID)
		}
	} else if c.CurrentPosture < c.MaxPosture {
		c.CurrentPosture += c.PostureRegenRate
		if c.CurrentPosture > c.MaxPosture {
			c.CurrentPosture = c.MaxPosture
		}
	}
}

func (c *Creature) ApplyPostureDamage(amount float64) {
	if c.IsPostureBroken || !c.IsAlive {
		return
	}

	c.CurrentPosture -= amount
	if c.CurrentPosture <= 0 {
		c.CurrentPosture = 0
		c.IsPostureBroken = true
		c.TimePostureBroken = time.Now().Unix()
		c.ChangeAIState(AIStateStaggered) // Certifique-se de ter esse estado no type.go
		log.Printf("[Creature %s] Posture quebrada! Entrando em stagger.", c.ID)
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
	case AIStateStaggered:
		// Lógica para stagger
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
		// Verificar expiração
		if now-eff.StartTime >= eff.Duration {
			log.Printf("[Effect] Creature %s: efeito %s expirou.", c.ID, eff.Type)
			continue
		}

		// Processamento de efeito contínuo (DOT, Regen, etc)
		if now-eff.LastTickTime >= eff.TickInterval {
			switch eff.Type {
			case EffectPoison, EffectBurn:
				c.HP -= eff.Power
				log.Printf("[Effect] Creature %s sofreu %d de %s. HP atual: %d", c.ID, eff.Power, eff.Type, c.HP)
				if c.HP <= 0 && c.IsAlive {
					c.IsAlive = false
					c.TimeOfDeath = now
					c.CurrentAction = ActionDie
					log.Printf("[Effect] Creature %s morreu por efeito %s", c.ID, eff.Type)
				}

			case EffectRegen:
				c.HP += eff.Power
				log.Printf("[Effect] Creature %s curou %d por %s. HP atual: %d", c.ID, eff.Power, eff.Type, c.HP)
			}

			eff.LastTickTime = now
		}

		// Controle de CC (Stun, Slow etc)
		if eff.Type == EffectStun {
			// Exemplo: podemos impedir ações dentro do ProcessAI
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
		fmt.Printf(
			"Creature: %s, Type: %s, Level: %s, AIState: %s, HP: %d, Action: %s, Pos: (%.2f, %.2f, %.2f), Posture: %.1f/%.1f\n",
			c.ID, c.Type, c.Level, c.AIState, c.HP, c.CurrentAction, c.Position.X, c.Position.Y, c.Position.Z, c.CurrentPosture, c.MaxPosture,
		)
	}
}
