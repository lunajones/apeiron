package creature

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/world/spatial"
)

type MemoryEvent struct {
	Description string
	Timestamp   time.Time
}

type Creature struct {
	Handle               handle.EntityHandle
	TargetCreatureHandle handle.EntityHandle
	TargetPlayerHandle   handle.EntityHandle

	Generation int

	model.Creature

	MoveCtrl         *movement.MovementController
	PrimaryType      consts.CreatureType
	Types            []consts.CreatureType
	Level            consts.CreatureLevel
	Actions          []consts.CreatureAction
	CurrentAction    consts.CreatureAction
	AIState          consts.AIState
	LastStateChange  time.Time
	LastAttackedTime time.Time

	Skills []string

	DynamicCombos map[consts.CreatureAction][]consts.CreatureAction

	Stagger       physics.StaggerData
	Invincibility physics.InvincibilityData

	Strength              int
	Dexterity             int
	Intelligence          int
	Focus                 int
	HitboxRadius          float64
	DesiredBufferDistance float64
	Position              position.Position
	LastPosition          position.Position
	HP                    int
	IsCrouched            bool
	IsHostile             bool
	IsAlive               bool
	IsCorpse              bool
	TimeOfDeath           time.Time

	PhysicalDefense    float64
	MagicDefense       float64
	RangedDefense      float64
	ControlResistance  float64
	StatusResistance   float64
	CriticalResistance float64
	CriticalChance     float64

	ActiveEffects []consts.ActiveEffect

	FieldOfViewDegrees float64
	VisionRange        float64
	HearingRange       float64
	SmellRange         float64
	IsBlind            bool
	IsDeaf             bool

	DetectionRadius float64
	AttackRange     float64

	SkillCooldowns map[consts.CreatureAction]time.Time
	AggroTable     map[handle.EntityHandle]*aggro.AggroEntry

	WalkSpeed   float64
	RunSpeed    float64
	AttackSpeed float64

	MaxPosture              float64
	CurrentPosture          float64
	PostureRegenRate        float64
	IsPostureBroken         bool
	TimePostureBroken       int64
	PostureBreakDurationSec int

	BehaviorTree BehaviorTree

	Needs          []consts.Need
	CurrentRole    consts.Role
	Tags           []consts.CreatureTag
	Memory         []MemoryEvent
	MentalState    consts.MentalState
	DamageWeakness map[consts.DamageType]float32

	FacingDirection position.Vector2D
	LastThreatSeen  time.Time
}

func (c *Creature) GetHandle() handle.EntityHandle {
	return c.Handle
}

type BehaviorTree interface {
	Tick(c *Creature, ctx interface{}) interface{}
}

func (c *Creature) GenerateSpawnPosition() position.Position {
	for attempts := 0; attempts < 10; attempts++ {
		offsetX := (rand.Float64()*2 - 1) * c.SpawnRadius
		offsetZ := (rand.Float64()*2 - 1) * c.SpawnRadius

		x := c.SpawnPoint.FastGlobalX() + offsetX
		y := c.SpawnPoint.FastGlobalY()
		z := c.SpawnPoint.Z + offsetZ

		newPos := position.FromGlobal(x, y, z)
		if IsTerrainWalkable(newPos) {
			log.Printf("[Creature] spawn em terreno válido: x=%.2f y=%.2f z=%.2f", x, y, z)
			return newPos
		}
	}

	x := c.SpawnPoint.FastGlobalX()
	y := c.SpawnPoint.FastGlobalY()
	z := c.SpawnPoint.Z
	log.Printf("[Creature] fallback no ponto de origem: x=%.2f y=%.2f z=%.2f", x, y, z)
	return c.SpawnPoint
}

func (c *Creature) SetPosition(newPos position.Position) {
	c.LastPosition = c.Position
	c.Position = newPos
}

func (c *Creature) GetPosition() position.Position {
	return c.Position
}

func (c *Creature) GetLastPosition() position.Position {
	return c.LastPosition
}

func (c *Creature) SetFacingDirection(dir position.Vector2D) {
	c.FacingDirection = dir
}

var creatures []*Creature

func Init() {
	log.Println("Creature service initialized")
	creatures = append(creatures, exampleSpawn())
}

func exampleSpawn() *Creature {
	id := lib.NewUUID()

	c := &Creature{
		Handle:               handle.NewEntityHandle(id, 1),
		TargetCreatureHandle: handle.EntityHandle{},
		TargetPlayerHandle:   handle.EntityHandle{},
		Generation:           1,
		Creature: model.Creature{
			Name:           "Example Creature",
			MaxHP:          100,
			SpawnPoint:     position.FromGlobal(0, 0, 0),
			SpawnRadius:    5.0,
			Faction:        "Monsters",
			RespawnTimeSec: 30,
		},
		HP:          100,
		IsAlive:     true,
		IsHostile:   true,
		IsCrouched:  false,
		PrimaryType: consts.Soldier,
		Types:       []consts.CreatureType{consts.Soldier, consts.Human},
		Level:       consts.Normal,
		Actions: []consts.CreatureAction{
			consts.ActionIdle, consts.ActionWalk, consts.ActionRun, consts.ActionParry,
			consts.ActionBlock, consts.ActionJump, consts.ActionSkill1, consts.ActionSkill2,
			consts.ActionSkill3, consts.ActionSkill4, consts.ActionSkill5,
			consts.ActionCombo1, consts.ActionCombo2, consts.ActionCombo3, consts.ActionDie,
		},
		Skills:                  []string{"Bite", "Lacerate"},
		CurrentAction:           consts.ActionIdle,
		AIState:                 consts.AIStateIdle,
		LastStateChange:         time.Now(),
		DynamicCombos:           make(map[consts.CreatureAction][]consts.CreatureAction),
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
		SmellRange:              10,
		DetectionRadius:         10.0,
		AttackRange:             2.5,
		SkillCooldowns:          make(map[consts.CreatureAction]time.Time),
		WalkSpeed:               1.0,
		RunSpeed:                3.5,
		AttackSpeed:             1.2,
		MaxPosture:              100,
		CurrentPosture:          100,
		PostureRegenRate:        1.5,
		PostureBreakDurationSec: 5,
		Needs:                   []consts.Need{{Type: consts.NeedHunger, Value: 0, Threshold: 50}},
		Tags:                    []consts.CreatureTag{consts.TagHumanoid},
	}

	c.Position = c.GenerateSpawnPosition()
	c.LastPosition = c.SpawnPoint

	return c
}

func (c *Creature) Tick(ctx interface{}) {
	if !c.IsAlive {
		return
	}

	c.TickNeeds()
	c.TickEffects()
	c.TickPosture()
	c.TickPhysicsStates()

	if c.BehaviorTree != nil {
		c.BehaviorTree.Tick(c, ctx)
	}
}

func (c *Creature) TickPhysicsStates() {
	now := time.Now().Unix()
	physics.TickStagger(&c.Stagger, now)
	physics.TickInvincibility(&c.Invincibility, now)
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
		c.ChangeAIState(consts.AIStateStaggered)
		now := time.Now().Unix()
		c.TimePostureBroken = now
		physics.StartStagger(&c.Stagger, now, float64(c.PostureBreakDurationSec))
		log.Printf("[Creature %s] Posture quebrada! Entrando em stagger.", c.Handle.ID)
	}
}

func (c *Creature) TickPosture() {
	if c.IsPostureBroken {
		if time.Now().Unix()-c.TimePostureBroken >= int64(c.PostureBreakDurationSec) {
			c.IsPostureBroken = false
			c.CurrentPosture = c.MaxPosture
			c.ChangeAIState(consts.AIStateIdle)
			log.Printf("[Creature %s] Posture recuperada.", c.Handle.String())
		}
	} else if c.CurrentPosture < c.MaxPosture {
		c.CurrentPosture += c.PostureRegenRate
		if c.CurrentPosture > c.MaxPosture {
			c.CurrentPosture = c.MaxPosture
		}
	}
}

func (c *Creature) TickEffects() {
	if !c.IsAlive {
		return
	}
	now := time.Now()
	var remainingEffects []consts.ActiveEffect

	for _, eff := range c.ActiveEffects {
		// Verificar expiração
		if now.Sub(eff.StartTime) >= eff.Duration {
			log.Printf("[Effect] Creature %s: efeito %s expirou.", c.Handle.String(), eff.Type)
			continue
		}

		// Processamento de efeito contínuo (DOT, Regen, etc)
		if now.Sub(eff.LastTickTime) >= eff.TickInterval {
			switch eff.Type {
			case consts.EffectPoison, consts.EffectBurn:
				c.HP -= eff.Power
				log.Printf("[Effect] Creature %s sofreu %d de %s. HP atual: %d", c.Handle.String(), eff.Power, eff.Type, c.HP)
				if c.HP <= 0 && c.IsAlive {
					c.ChangeAIState(consts.AIStateDead)
					log.Printf("[Effect] Creature %s morreu por efeito %s", c.Handle.String(), eff.Type)
				}

			case consts.EffectRegen:
				c.HP += eff.Power
				log.Printf("[Effect] Creature %s curou %d por %s. HP atual: %d", c.Handle.String(), eff.Power, eff.Type, c.HP)
			}

			eff.LastTickTime = now
		}

		// Controle de CC (Stun, Slow etc)
		if eff.Type == consts.EffectStun {
			// Exemplo: podemos impedir ações dentro do ProcessAI
		}

		remainingEffects = append(remainingEffects, eff)
	}

	c.ActiveEffects = remainingEffects
}

func (c *Creature) TickNeeds() {
	ModifyNeed(c, consts.NeedHunger, 0.007) // 0 → 50 em ~2 horas reais
	ModifyNeed(c, consts.NeedThirst, 0.008) // ligeiramente mais rápida
	ModifyNeed(c, consts.NeedSleep, 0.004)  // mais lenta
}

func IsTerrainWalkable(pos position.Position) bool {
	// TODO: Integrar com TerrainService real
	return true
}

func CanSeePlayer(c *Creature, players []*player.Player) bool {
	for _, p := range players {
		toPlayer := position.Vector2D{
			X: p.Position.FastGlobalX() - c.Position.FastGlobalX(),
			Y: p.Position.Z - c.Position.Z,
		}
		distance := toPlayer.Magnitude()
		if distance > c.VisionRange {
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

func CanHearPlayer(c *Creature, players []*player.Player) bool {
	for _, p := range players {
		distance := position.CalculateDistance(c.Position, p.Position)
		if distance <= c.HearingRange {
			return true
		}
	}
	return false
}

var aiStateToAction = map[consts.AIState]consts.CreatureAction{
	consts.AIStateIdle:   consts.ActionIdle,
	consts.AIStateAlert:  consts.ActionIdle,
	consts.AIStateAttack: consts.ActionAttack,
	consts.AIStateDead:   consts.ActionDie,
	// consts.AIStateStaggered: pode ser adicionado se necessário
}

func (c *Creature) SetAction(action consts.CreatureAction) {
	if c.CurrentAction != action {
		c.CurrentAction = action
		log.Printf("[Creature] %s (%s) ação definida para %s", c.Handle.String(), c.PrimaryType, action)
	}
}

func (c *Creature) ChangeAIState(newState consts.AIState) {
	if c.AIState == newState {
		return
	}

	log.Printf("[Creature] %s (%s) AI State mudou: %s → %s", c.Handle.String(), c.PrimaryType, c.AIState, newState)
	c.AIState = newState
	c.LastStateChange = time.Now()

	if act, ok := aiStateToAction[newState]; ok {
		c.SetAction(act)
	}

	if newState == consts.AIStateDead {
		c.IsAlive = false
		c.IsCorpse = true
		c.TimeOfDeath = time.Now()
	}

	// lógica de stagger pode ser inserida aqui
}

func (c *Creature) GenerateRandomCombo(comboAction consts.CreatureAction) {
	possibleSkills := []consts.CreatureAction{
		consts.ActionSkill1,
		consts.ActionSkill2,
		consts.ActionSkill3,
		consts.ActionSkill4,
		consts.ActionSkill5,
	}

	var combo []consts.CreatureAction
	numSkillsInCombo := rand.Intn(4) + 2

	for i := 0; i < numSkillsInCombo; i++ {
		randomSkill := possibleSkills[rand.Intn(len(possibleSkills))]
		combo = append(combo, randomSkill)
	}

	if c.DynamicCombos == nil {
		c.DynamicCombos = make(map[consts.CreatureAction][]consts.CreatureAction)
	}

	c.DynamicCombos[comboAction] = combo
	log.Printf("[Creature %s] Novo combo gerado para %s: %v", c.Handle.String(), comboAction, combo)
}

func (c *Creature) IsQuestOnly() bool {
	return strings.TrimSpace(c.OwnerPlayerID) != ""
}

func (c *Creature) ApplyEffect(effect consts.ActiveEffect) {
	c.ActiveEffects = append(c.ActiveEffects, effect)
	log.Printf("[Creature %s] recebeu efeito: %s", c.Handle.String(), effect.Type)
}

func (c *Creature) WasRecentlyAttacked() bool {
	return time.Since(c.LastAttackedTime).Seconds() < 10 // Exemplo básico
}

// --- Função FindByID ---
func FindByHandleID(creatures []*Creature, id string) *Creature {
	for _, c := range creatures {
		if c.Handle.ID == id {
			return c
		}
	}
	return nil
}

func (c *Creature) HasTag(tag consts.CreatureTag) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *Creature) HasType(t consts.CreatureType) bool {
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
	c.LastThreatSeen = time.Time{}
	c.ClearAggro()
	c.ClearCooldowns()

	c.Generation++
	c.Handle = handle.NewEntityHandle(lib.NewUUID(), c.Generation)
	// Pode adicionar mais resets aqui depois
}

func (c *Creature) ClearCooldowns() {
	if c.SkillCooldowns != nil {
		for action := range c.SkillCooldowns {
			delete(c.SkillCooldowns, action)
		}
	}
}

// --- THREAT SYSTEM ---

// ----------------------
// AGGRO SYSTEM (HANDLE-BASED)
// ----------------------

func (c *Creature) AddThreat(targetHandle handle.EntityHandle, amount float64, source, action string) {
	if c.AggroTable == nil {
		c.AggroTable = make(map[handle.EntityHandle]*aggro.AggroEntry)
	}

	entry, exists := c.AggroTable[targetHandle]
	if !exists {
		entry = &aggro.AggroEntry{
			TargetHandle: targetHandle,
		}
		c.AggroTable[targetHandle] = entry
	}

	entry.ThreatValue += amount
	entry.LastDamageTime = time.Now()
	entry.AggroSource = source
	entry.LastAction = action

	log.Printf("[Aggro] %s recebeu %.2f de ameaça de %s (source: %s, action: %s). Ameaça total: %.2f",
		c.Handle.ID, amount, targetHandle.ID, source, action, entry.ThreatValue)
}

func (c *Creature) GetHighestThreatTarget() handle.EntityHandle {
	var topHandle handle.EntityHandle
	var topThreat float64

	for h, entry := range c.AggroTable {
		if entry.ThreatValue > topThreat {
			topThreat = entry.ThreatValue
			topHandle = h
		}
	}

	return topHandle
}

func (c *Creature) ClearAggro() {
	if c.AggroTable != nil {
		for h := range c.AggroTable {
			delete(c.AggroTable, h)
		}
		log.Printf("[Aggro] %s limpou toda tabela de ameaça", c.Handle.ID)
	}
}

func (c *Creature) ReduceThreatOverTime(decayRatePerSecond float64) {
	now := time.Now()

	for h, entry := range c.AggroTable {
		elapsed := now.Sub(entry.LastDamageTime).Seconds()
		decay := decayRatePerSecond * elapsed

		entry.ThreatValue -= decay
		if entry.ThreatValue <= 0 {
			delete(c.AggroTable, h)
			log.Printf("[Aggro] %s perdeu o aggro de %s por decay.", c.Handle.ID, h.ID)
			// Futuro: disparar OnThreatLost
		} else {
			entry.LastDamageTime = now
		}
	}
}

func (c *Creature) CheckIsAlive() bool {
	return c.IsAlive
}

func (c *Creature) GetCurrentSpeed() float64 {
	switch c.CurrentAction {
	case consts.ActionRun:
		return c.WalkSpeed
	default:
		return c.RunSpeed
	}
}

func (c *Creature) TakeDamage(amount int) {
	if !c.IsAlive {
		return
	}

	finalDamage := int(math.Round(float64(amount) * (1.0 - c.PhysicalDefense)))
	if finalDamage <= 0 {
		finalDamage = 1
	}

	c.HP -= finalDamage
	c.LastAttackedTime = time.Now()

	log.Printf("[Creature %s] sofreu %d de dano. HP restante: %d", c.Handle.String(), finalDamage, c.HP)

	if c.HP <= 0 {
		c.ChangeAIState(consts.AIStateDead)
		log.Printf("[Creature %s] morreu após receber dano.", c.Handle.String())
	}

	if c.HP > 0 && !c.IsHostile && c.HP < c.MaxHP && c.AIState != consts.AIStateFleeing && c.AIState != consts.AIStateDead {
		log.Printf("[Creature %s] recebeu dano e vai fugir!", c.Handle.String())
		c.ChangeAIState(consts.AIStateFleeing)
	}
}

func (c *Creature) GetFacingDirection() position.Vector2D {
	return c.FacingDirection
}

// creature/creature.go

func (c *Creature) GetBestTarget(creatures []*Creature, players []*player.Player) model.Targetable {
	bestHandle := c.GetHighestThreatTarget()
	if !bestHandle.Equals(handle.EntityHandle{}) {
		// Procura entre criaturas
		for _, c2 := range creatures {
			if c2.Handle.Equals(bestHandle) {
				return c2
			}
		}
		// Procura entre jogadores
		for _, p := range players {
			if p.Handle.Equals(bestHandle) {
				return p
			}
		}
	}

	// Se não houver aggro válido, busca o mais próximo
	var closest model.Targetable
	var minDist float64 = math.MaxFloat64

	for _, p := range players {
		if !p.CheckIsAlive() {
			continue
		}
		dist := position.CalculateDistance(c.Position, p.Position)
		if dist < minDist {
			minDist = dist
			closest = p
		}
	}

	for _, c2 := range creatures {
		if c2.Handle.Equals(c.Handle) || !c2.IsAlive {
			continue
		}
		dist := position.CalculateDistance(c.Position, c2.Position)
		if dist < minDist {
			minDist = dist
			closest = c2
		}
	}

	return closest
}

func (c *Creature) IsHungry() bool {
	for _, n := range c.Needs {
		if n.Type == consts.NeedHunger && n.Value >= n.Threshold {
			return true
		}
	}
	return false
}

func (c *Creature) ClearTargetHandles() {
	if !c.TargetCreatureHandle.Equals(handle.EntityHandle{}) || !c.TargetPlayerHandle.Equals(handle.EntityHandle{}) {
		log.Printf("[Creature] [%s (%s)] limpando alvos: criatura=%s, jogador=%s",
			c.Handle.ID, c.PrimaryType, c.TargetCreatureHandle.ID, c.TargetPlayerHandle.ID)
	}
	c.TargetCreatureHandle = handle.EntityHandle{}
	c.TargetPlayerHandle = handle.EntityHandle{}
}

func (c *Creature) IsCurrentlyCrouched() bool {
	return c.IsCrouched
}

func (c *Creature) GetHitboxRadius() float64 {
	return c.HitboxRadius
}

func (c *Creature) GetDesiredBufferDistance() float64 {
	return c.DesiredBufferDistance
}

func (c *Creature) ConsumeCorpse() {
	c.IsCorpse = false // já indica que não é mais válido como corpo
	// Aqui você pode acionar a remoção do spatial grid se tiver essa função
	spatial.GlobalGrid.Remove(c)
	log.Printf("[Creature] %s (%s) foi consumido e removido do mundo.", c.Handle.String(), c.PrimaryType)
}

func (c *Creature) IsCreature() bool { return true }
