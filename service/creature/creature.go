package creature

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"math"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/ai_context"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature/aggro"
)

type MemoryEvent struct {
	Description string
	Timestamp   time.Time
}

type Creature struct {
	model.Creature  // Embutindo o modelo base

	PrimaryType     CreatureType
	Types           []CreatureType
	Level           CreatureLevel
	Actions         []CreatureAction
	CurrentAction   CreatureAction
	AIState         AIState
	LastStateChange time.Time
	LastAttackedTime time.Time
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

	// Efeitos ativos
	ActiveEffects []ActiveEffect

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
	AggroTable map[string]*aggro.AggroEntry

	// Movement and Attack speed
	MoveSpeed   float64
	AttackSpeed float64

	// Posture / Stagger system
	MaxPosture              float64
	CurrentPosture          float64
	PostureRegenRate        float64
	IsPostureBroken         bool
	TimePostureBroken       int64
	PostureBreakDurationSec int

	// AI Behavior Tree
	BehaviorTree    BehaviorTree

	// AI Behavior decision fields
	Needs          []Need
	CurrentRole    Role
	Tags           []CreatureTag
	Memory         []MemoryEvent
	MentalState    MentalState
	DamageWeakness map[DamageType]float32

	FacingDirection position.Vector2D
}
type BehaviorTree interface {
	Tick(c *Creature, ctx interface{}) interface{}
}

var creatures []*Creature

func Init() {
	log.Println("Creature service initialized")
	creatures = append(creatures, exampleSpawn())
}

func exampleSpawn() *Creature {
	c := &Creature{
		Creature: model.Creature{
			ID:             lib.NewUUID(),
			HP:             100,
			MaxHP:          100,
			IsAlive:        true,
			RespawnTimeSec: 30,
			SpawnPoint:     position.Position{X: 0, Y: 0, Z: 0},
			SpawnRadius:    5.0,
			Faction:        "Monsters",
			IsHostile:      true,
		},

		PrimaryType:     Soldier,
		Types:           []CreatureType{Soldier, Human},
		Level:           Normal,
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
		Strength:                20,
		Dexterity:               10,
		Intelligence:            5,
		Focus:                   8,
		PhysicalDefense:         0.15,
		MagicDefense:            0.05,
		RangedDefense:           0.10,
		ControlResistance:       0.1,
		StatusResistance:        0.1,
		CriticalResistance:      0.2,
		CriticalChance:          0.05,
		FieldOfViewDegrees:      120,
		VisionRange:             15,
		HearingRange:            10,
		DetectionRadius:         10.0,
		AttackRange:             2.5,
		SkillCooldowns:          make(map[CreatureAction]time.Time),
		AggroTable:              make(map[string]*aggro.AggroEntry),
		MoveSpeed:               3.5,
		AttackSpeed:             1.2,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		Needs: []Need{
			{Type: NeedHunger, Value: 0, Threshold: 50},
		},
		Tags: []CreatureTag{TagHumanoid},
	}
	c.Position = c.GenerateSpawnPosition()
	return c
}



func (c *Creature) GenerateSpawnPosition() position.Position {
	for attempts := 0; attempts < 10; attempts++ {
		offsetX := (rand.Float64()*2 - 1) * c.SpawnRadius
		offsetZ := (rand.Float64()*2 - 1) * c.SpawnRadius
		newPos := position.Position{
			X: c.SpawnPoint.X + offsetX,
			Y: c.SpawnPoint.Y,
			Z: c.SpawnPoint.Z + offsetZ,
		}

		if IsTerrainWalkable(newPos) {
			log.Println("[Creature] spawning in walkable terrain in position %s, %s, %s", newPos.X, newPos.Y, newPos.Z)

			return newPos
		}
	}
	
	log.Println("[Creature] spawning in walkable terrain in poisition %s, %s, %s", c.SpawnPoint.X, c.SpawnPoint.Y, c.SpawnPoint.Z)
	return c.SpawnPoint
}

func IsTerrainWalkable(pos position.Position) bool {
	// TODO: Integrar com TerrainService real
	return true
}

func (c *Creature) Tick(ctx ai_context.AIContext) {
	c.TickEffects()
	c.TickPosture()

	switch c.AIState {
	case AIStateIdle:
		if rand.Float32() < 0.1 {
			c.ChangeAIState(AIStateAlert)
		}

	case AIStateAlert:
		for _, p := range ctx.GetPlayers() {
			singlePlayerSlice := []*model.Player{p}
			if CanSeePlayer(c, singlePlayerSlice) || CanHearPlayer(c, singlePlayerSlice) {
				c.AddThreat(p.ID, 10, "PlayerDetected", "VisionOrSound")
				log.Printf("[AI] %s detectou o player %s e adicionou threat.", c.ID, p.ID)
				c.ChangeAIState(AIStateChasing)
				break
			}
		}

		if time.Since(c.LastStateChange) > 2*time.Second {
			c.ChangeAIState(AIStateIdle)
		}

	case AIStateChasing:
		targetID := c.GetHighestThreatTarget()
		if targetID == "" {
			log.Printf("[AI] %s sem alvo de threat, voltando pra Idle", c.ID)
			c.ChangeAIState(AIStateIdle)
			return
		}

		target := findTargetByID(targetID, ctx.GetCreatures(), ctx.GetPlayers())
		if target == nil {
			log.Printf("[AI] %s: alvo %s não encontrado, limpando aggro", c.ID, targetID)
			c.ClearAggro()
			c.ChangeAIState(AIStateIdle)
			return
		}

		c.MoveTowards(target.GetPosition(), c.MoveSpeed)

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
		c.IsCorpse = true
		c.TimeOfDeath = time.Now()
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
	now := time.Now()
	var remainingEffects []ActiveEffect

	for _, eff := range c.ActiveEffects {
		// Verificar expiração
		if now.Sub(eff.StartTime) >= eff.Duration {
			log.Printf("[Effect] Creature %s: efeito %s expirou.", c.ID, eff.Type)
			continue
		}

		// Processamento de efeito contínuo (DOT, Regen, etc)
		if now.Sub(eff.LastTickTime) >= eff.TickInterval {
			switch eff.Type {
			case EffectPoison, EffectBurn:
				c.HP -= eff.Power
				log.Printf("[Effect] Creature %s sofreu %d de %s. HP atual: %d", c.ID, eff.Power, eff.Type, c.HP)
				if c.HP <= 0 && c.IsAlive {
					c.ChangeAIState(AIStateDead)
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

func TickAll(ctx ai_context.AIContext) {
	for _, c := range creatures {
		c.Tick(ctx)
	}
}

func DebugPrintCreatures() {
	for _, c := range creatures {
		fmt.Printf(
			"Creature: %s, Type: %s, Level: %s, AIState: %s, HP: %d, Action: %s, Pos: (%.2f, %.2f, %.2f), Posture: %.1f/%.1f\n",
			c.ID, c.Types, c.Level, c.AIState, c.HP, c.CurrentAction, c.Position.X, c.Position.Y, c.Position.Z, c.CurrentPosture, c.MaxPosture,
		)
	}
}

func (c *Creature) WasRecentlyAttacked() bool {
	return time.Since(c.LastAttackedTime).Seconds() < 10 // Exemplo básico
}

// --- Função FindByID ---
func FindByID(creatures []*model.Creature, id string) *model.Creature {
	for _, c := range creatures {
		if c.ID == id {
			return c
		}
	}
	return nil
}

// --- Função CanSeeOtherCreatures ---
func CanSeeOtherCreatures(c *Creature, creatures []*Creature) bool {
	// Exemplo simples só para compilar
	return len(creatures) > 0
}

// --- Função CanHearOtherCreatures ---
func CanHearOtherCreatures(c *Creature, creatures []*Creature) bool {
	// Exemplo simples só para compilar
	return len(creatures) > 0
}

func CanSeePlayer(c *Creature, players []*model.Player) bool {
	for _, p := range players {
		toPlayer := position.Vector2D{
			X: p.Position.X - c.Position.X,
			Y: p.Position.Z - c.Position.Z, // Considerando plano XZ (horizontal)
		}

		distance := toPlayer.Magnitude()
		if distance > c.FieldOfViewDegrees {
			continue
		}

		toPlayerNormalized := toPlayer.Normalize()
		facing := c.FacingDirection.Normalize()

		dot := facing.Dot(toPlayerNormalized)

		fovRadians := (c.FieldOfViewDegrees / 2) * (math.Pi / 180)
		cosFov := math.Cos(fovRadians)

		if dot >= cosFov {
			return true
		}
	}
	return false
}


func CanHearPlayer(c *Creature, players []*model.Player) bool {
	for _, p := range players {
		distance := position.CalculateDistance(c.Position, p.Position)
		if distance <= c.HearingRange {
			return true
		}
	}
	return false
}

func (c *Creature) GetNeedValue(needType NeedType) float64 {
	for _, n := range c.Needs {
		if n.Type == needType {
			return n.Value
		}
	}
	return 0
}

func (c *Creature) HasTag(tag CreatureTag) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *Creature) SetNeedValue(needType NeedType, value float64) {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			c.Needs[i].Value = value
			return
		}
	}
}

func (c *Creature) ModifyNeed(needType NeedType, amount float64) {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			c.Needs[i].Value += amount
			if c.Needs[i].Value < 0 {
				c.Needs[i].Value = 0
			}
			if c.Needs[i].Value > 100 {
				c.Needs[i].Value = 100
			}
			break
		}
	}
}

func (c *Creature) HasType(t CreatureType) bool {
	for _, ct := range c.Types {
		if ct == t {
			return true
		}
	}
	return false
}

func (c *Creature) Respawn() {
	c.HP = c.MaxHP
	c.IsAlive = true
	c.Position = c.GenerateSpawnPosition()
	c.TimeOfDeath = time.Time{} // Zera o campo
	c.ClearAggro()
	c.ClearCooldowns()
	// Pode adicionar mais resets aqui depois
}

func (c *Creature) ClearCooldowns() {
	if c.SkillCooldowns != nil {
		for action := range c.SkillCooldowns {
			delete(c.SkillCooldowns, action)
		}
	}
}

func (c *Creature) AddThreat(targetID string, amount float64, source, action string) {
	if c.AggroTable == nil {
		c.AggroTable = make(map[string]*aggro.AggroEntry)
	}

	entry, exists := c.AggroTable[targetID]
	if !exists {
		entry = &aggro.AggroEntry{
			TargetID: targetID,
		}
		c.AggroTable[targetID] = entry
	}

	entry.ThreatValue += amount
	entry.LastDamageTime = time.Now()
	entry.AggroSource = source
	entry.LastAction = action
}

func (c *Creature) GetHighestThreatTarget() string {
	var topTarget string
	var topThreat float64

	for targetID, entry := range c.AggroTable {
		if entry.ThreatValue > topThreat {
			topThreat = entry.ThreatValue
			topTarget = targetID
		}
	}

	return topTarget
}

func (c *Creature) ClearAggro() {
	if c.AggroTable != nil {
		for targetID := range c.AggroTable {
			delete(c.AggroTable, targetID)
		}
	}
}


func (c *Creature) ReduceThreatOverTime(decayRatePerSecond float64) {
	now := time.Now()

	for targetID, entry := range c.AggroTable {
		elapsed := now.Sub(entry.LastDamageTime).Seconds()
		decay := decayRatePerSecond * elapsed

		entry.ThreatValue -= decay
		if entry.ThreatValue <= 0 {
			delete(c.AggroTable, targetID)
			log.Printf("[Aggro] Creature %s perdeu o aggro de %s por decay.", c.ID, targetID)
			// No futuro: Aqui você pode disparar um OnThreatLost
		} else {
			entry.LastDamageTime = now
		}
	}
}

func (c *Creature) MoveTowards(targetPos position.Position, speed float64) {
	dx := targetPos.X - c.Position.X
	dz := targetPos.Z - c.Position.Z
	dist := math.Sqrt(dx*dx + dz*dz)

	if dist < 0.01 {
		return
	}

	// Atualiza a direção que a criatura está olhando
	c.FacingDirection = position.Vector2D{
		X: dx / dist,
		Y: dz / dist,
	}

	moveX := (dx / dist) * speed
	moveZ := (dz / dist) * speed

	c.Position.X += moveX
	c.Position.Z += moveZ

	log.Printf("[AI] %s movendo-se em direção a (%.2f, %.2f)", c.ID, targetPos.X, targetPos.Z)
}


func (c *Creature) GetPosition() position.Position {
	return c.Position
}

func (c *Creature) GetID() string {
	return c.ID
}

func findTargetByID(id string, creatures []*model.Creature, players []*model.Player) model.Targetable {
	for _, c := range creatures {
		if c.ID == id {
			return c
		}
	}
	for _, p := range players {
		if p.ID == id {
			return p
		}
	}
	return nil
}