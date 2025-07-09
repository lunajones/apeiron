package creature

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/lunajones/apeiron/lib"
	"github.com/lunajones/apeiron/lib/combat"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature/aggro"
	"github.com/lunajones/apeiron/service/creature/consts"
	"github.com/lunajones/apeiron/service/helper/finder"
	"github.com/lunajones/apeiron/service/player"
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
	AIState          constslib.AIState
	CombatState      constslib.CombatState
	AnimationState   constslib.AnimationState
	LastStateChange  time.Time
	LastAttackedTime time.Time

	RegisteredSkills []*model.Skill
	NextSkillToUse   *model.Skill

	SkillStates        map[constslib.SkillAction]*model.SkillState
	SkillMovementState *model.SkillMovementState

	Strength              int
	Dexterity             int
	Intelligence          int
	Focus                 int
	HitboxRadius          float64
	DesiredBufferDistance float64
	MinWanderDistance     float64
	MaxWanderDistance     float64
	WanderStopDistance    float64
	Position              position.Position
	LastPosition          position.Position
	HP                    int
	IsCrouched            bool
	Hostile               bool
	Alive                 bool
	IsCorpse              bool
	TimeOfDeath           time.Time

	PhysicalDefense    float64
	MagicDefense       float64
	RangedDefense      float64
	ControlResistance  float64
	StatusResistance   float64
	CriticalResistance float64
	CriticalChance     float64

	ActiveEffects []constslib.ActiveEffect

	FieldOfViewDegrees float64
	VisionRange        float64
	HearingRange       float64
	SmellRange         float64
	IsBlind            bool
	IsDeaf             bool

	DetectionRadius float64
	AttackRange     float64

	WalkSpeed        float64
	RunSpeed         float64
	OriginalRunSpeed float64
	AttackSpeed      float64

	// POSTURA
	MaxPosture           float64
	Posture              float64
	PostureRegenRate     float64
	PostureBroken        bool
	TimePostureBroken    int64
	PostureBrokenElapsed float64 // em segundos

	PostureBreakDurationSec float64
	ReceivedDamageRecently  bool

	blocking         bool
	BlockStartedAt   time.Time
	BlockDuration    time.Duration
	MaxBlockDuration time.Duration
	// STAMINA
	dodging            bool
	Stamina            float64
	BlockSpentStamina  float64
	MaxStamina         float64
	StaminaRegenPerSec float64

	DodgeStartedAt               time.Time
	DodgeDistance                float64
	DodgeStaminaCost             float64
	DodgeDisabledUntil           time.Time
	DodgeInvulnerabilityDuration time.Duration

	invulnerableUntil time.Time

	BehaviorTree BehaviorTree

	NextAggressiveDecisionAllowed time.Time
	NextDefensiveDecisionAllowed  time.Time
	NextStrategicDecisionAllowed  time.Time

	Needs          []constslib.Need
	CurrentRole    consts.Role
	Tags           []consts.CreatureTag
	Memory         []MemoryEvent
	MentalState    consts.MentalState
	DamageWeakness map[constslib.DamageType]float32

	FacingDirection position.Vector2D
	LastThreatSeen  time.Time

	AggroTable        map[handle.EntityHandle]*aggro.AggroEntry
	LastKnownDistance float64

	ParryWindowStart      time.Time
	ParryWindowEnd        time.Time
	BlockStaminaTolerance float64

	BlockableChance float64 // ex: 0.7 (70%)
	DodgableChance  float64 // ex: 0.3 (30%)
}

func (c *Creature) GetHandle() handle.EntityHandle {
	return c.Handle
}

type BehaviorTree interface {
	Tick(c *Creature, ctx interface{}) interface{}
}

// Ajuste no método GenerateSpawnPosition esperado no Creature
func (c *Creature) GenerateSpawnPosition(mesh *navmesh.NavMesh) position.Position {
	for attempts := 0; attempts < 10; attempts++ {
		offsetXVal := (rand.Float64()*2 - 1) * c.Creature.SpawnRadius
		offsetZVal := (rand.Float64()*2 - 1) * c.Creature.SpawnRadius

		x := c.Creature.SpawnPoint.X + offsetXVal
		z := c.Creature.SpawnPoint.Z + offsetZVal
		y := c.Creature.SpawnPoint.Y

		newPos := position.Position{
			X: x,
			Y: y,
			Z: z,
		}

		if mesh.IsWalkable(newPos) {
			log.Printf("[SPAWN DEBUG] Gerado para Position: X=%.2f Z=%.2f Y=%.2f", x, z, y)
			return newPos
		}
	}

	log.Printf("[SPAWN DEBUG] Falha ao gerar posição válida, retornando ponto original")
	return c.Creature.SpawnPoint
}

func (c *Creature) SetPosition(newPos position.Position) {
	// log.Printf("[Creature] [%s (%s)] SetPosition: nova posição = %.2f, %.2f, %.2f",
	// 	c.Handle.String(), c.PrimaryType, newPos.X, newPos.Y, newPos.Z)

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

func (c *Creature) Tick(ctx *dynamic_context.AIServiceContext, deltaTime float64) {

	if !c.Alive {
		return
	}

	c.PerformDefensiveAction(ctx, deltaTime)
	// 3️⃣ Atualiza movimento baseado em habilidades (ex: Leap)
	if c.SkillMovementState != nil && c.SkillMovementState.Active {
		target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, ctx)

		if combat.UpdateSkillMovement(c, c.SkillMovementState, target, ctx.NavMesh, ctx, deltaTime) {
			log.Printf("[LEAP] [%s] SkillMovement concluído", c.Handle.String())
			c.SkillMovementState = nil
			c.CombatState = constslib.CombatStateIdle
		}
	}

	// 4️⃣ Se não está em movimento de habilidade, atualiza movimento normal
	if c.SkillMovementState == nil || !c.SkillMovementState.Active {
		c.MoveCtrl.Update(c, deltaTime, ctx)
	}

	// 5️⃣ Ticks de status
	c.TickNeeds(deltaTime)
	c.TickEffects(deltaTime)
	c.TickPosture(deltaTime)
	c.TickStamina(deltaTime)

	if c.IsDodging() && time.Since(c.DodgeStartedAt) > c.DodgeInvulnerabilityDuration {
		c.SetDodging(false)
		log.Printf("[DODGE] [%s] fim da invulnerabilidade (%.1fs)", c.Handle.String(), c.DodgeInvulnerabilityDuration.Seconds())
	}

	// 6️⃣ AI Behavior Tree
	if c.BehaviorTree != nil {
		c.BehaviorTree.Tick(c, ctx)
	}

}

func (c *Creature) PerformDefensiveAction(ctx *dynamic_context.AIServiceContext, deltaTime float64) {
	if !c.Alive {
		return
	}

	total := c.BlockableChance + c.DodgableChance
	if total == 0 {
		return // criatura não defende
	}

	events := ctx.GetRecentAggressorsAgainst(c.Handle, time.Now().Add(-1*time.Second))
	if len(events) > 0 {
		// 🔎 Seleciona o evento mais recente
		latest := events[len(events)-1]

		log.Printf("\033[38;5;154m[DEFENSE] [%s] reagindo ao evento mais recente (type=%s, posture=%s)\033[0m",
			c.Handle.String(), latest.BehaviorType, c.Posture)

		r := rand.Float64() * total
		if r < c.DodgableChance {
			c.TryDodgeReaction(latest)
		} else {
			c.TryBlockReaction(latest)
		}
	}

	// Sempre atualiza estados atuais, mesmo sem eventos
	c.PerformDodge(ctx)
	c.PerformBlock(deltaTime)
}

func (c *Creature) TryBlockReaction(e dynamic_context.CombatBehaviorEvent) {
	if !c.Alive || c.PostureBroken || c.IsBlocking() || c.IsDodging() {
		return
	}
	// Ajusta margem de tempo conforme o estado de combate
	baseMargin := 400 * time.Millisecond
	randFactor := 0.25
	stateReactionFactor := 2.0
	switch c.CombatState {
	case constslib.CombatStateDefensive:
		stateReactionFactor = 1.4
	case constslib.CombatStateAggressive:
		stateReactionFactor = 0.6
	}
	adjustedMargin := time.Duration(stateReactionFactor * float64(baseMargin))
	randomMargin := time.Duration(rand.Float64() * randFactor * float64(e.WindupTime))
	start := e.Timestamp
	end := start.Add(e.WindupTime).Add(adjustedMargin).Add(randomMargin)
	now := time.Now()

	if now.After(start) && now.Before(end) {
		c.SetBlocking(true)
		c.SetDodging(false)
		c.BlockStartedAt = now // usado opcionalmente para logging ou cálculo
		log.Printf("[REACT] [%s] iniciou bloqueio contra %s", c.Handle.String(), e.SourceHandle.ID)
	}
}

func (c *Creature) PerformBlock(deltaTime float64) {
	// Se não está bloqueando, nada a fazer (reseta stamina gasta)
	if !c.IsBlocking() {
		c.BlockSpentStamina = 0
		return
	}
	now := time.Now()

	// Início do bloqueio: aplica custo inicial e define janela de parry
	if c.BlockDuration == 0 {
		initialCost := 10.0
		c.ReduceStamina(initialCost)
		var parryDuration time.Duration
		switch c.CombatState {
		case constslib.CombatStateStrategic:
			parryDuration = 450 * time.Millisecond
		case constslib.CombatStateAggressive:
			parryDuration = 200 * time.Millisecond
		default:
			parryDuration = 300 * time.Millisecond
		}
		c.ParryWindowStart = now
		c.ParryWindowEnd = now.Add(parryDuration)
		log.Printf("[PARRY] [%s] parry window ativa por %v", c.Handle.String(), parryDuration)

		c.BlockStaminaTolerance = 5.0 + rand.Float64()*5.0 // rand.Float64 em [0,1):contentReference[oaicite:6]{index=6}
		c.MaxBlockDuration = generateBlockDuration(c.CombatState)
	}

	// Atualiza duração e consumo de stamina do bloqueio
	c.BlockDuration += time.Duration(deltaTime * float64(time.Second))
	staminaPerSecond := 1.0
	staminaThisTick := staminaPerSecond * deltaTime * 10
	c.BlockSpentStamina += staminaThisTick
	c.ReduceStamina(staminaThisTick)
	log.Printf("[BLOCK] [%s] mantendo bloqueio, consumiu %.3f stamina neste tick", c.Handle.String(), staminaThisTick)
	c.SetDodging(false)
	log.Printf("[BLOCK-CHECK] [%s] stamina gasta=%.2f, tolerância=%.2f, duração=%.2fs / limite=%.2fs",
		c.Handle.String(),
		c.BlockSpentStamina,
		c.BlockStaminaTolerance,
		c.BlockDuration.Seconds(),
		c.MaxBlockDuration.Seconds(),
	)

	// Finaliza o bloqueio quando exceder tolerância ou duração máxima
	if c.BlockSpentStamina > c.BlockStaminaTolerance || c.BlockDuration >= c.MaxBlockDuration {
		c.SetBlocking(false)
		c.ParryWindowStart = time.Time{}
		c.ParryWindowEnd = time.Time{}
		c.BlockStaminaTolerance = 0
		c.BlockDuration = 0
		c.MaxBlockDuration = 0
		c.BlockSpentStamina = 0
		log.Printf("[BLOCK] [%s] soltou o bloqueio — motivo: %s",
			c.Handle.String(),
			func() string {
				if c.Stamina <= c.BlockStaminaTolerance {
					return "stamina abaixo da tolerância"
				}
				return "duração máxima atingida"
			}(),
		)
	}
}

func (c *Creature) TryDodgeReaction(e dynamic_context.CombatBehaviorEvent) {
	if !c.Alive || c.IsDodging() || c.IsBlocking() {
		return
	}
	// Verifica se há stamina suficiente para esquivar
	if c.Stamina < c.DodgeStaminaCost+5.0 {
		log.Printf("[DODGE] [%s] recusou esquiva — stamina insuficiente (%.2f)", c.Handle.String(), c.Stamina)
		return
	}
	// Chance base de esquiva por estado de combate
	baseChance := 0.75
	switch c.CombatState {
	case constslib.CombatStateDefensive:
		baseChance = 1.0
	case constslib.CombatStateStrategic:
		baseChance = 0.85
	case constslib.CombatStateAggressive:
		baseChance = 0.75
	}
	staminaRatio := c.Stamina / c.MaxStamina // de 0.0 a 1.0
	finalChance := baseChance * (0.25 + 0.75*staminaRatio)
	if rand.Float64() >= finalChance {
		log.Printf("[DODGE] [%s] não esquivou — chance final %.2f (stamina ratio: %.2f)", c.Handle.String(), finalChance, staminaRatio)
		return
	}
	start := e.Timestamp
	end := start.Add(e.WindupTime)
	now := time.Now()
	if now.After(start) && now.Before(end) {
		c.SetDodging(true)
		c.SetBlocking(false)
		c.invulnerableUntil = now.Add(c.DodgeInvulnerabilityDuration)
		c.DodgeStartedAt = now
		log.Printf("[REACT] [%s] iniciou esquiva contra %s", c.Handle.String(), e.SourceHandle.ID)
	}
}

func (c *Creature) PerformDodge(svcCtx *dynamic_context.AIServiceContext) {
	if !c.IsDodging() {
		return
	}
	if time.Now().Before(c.DodgeDisabledUntil) {
		log.Printf("[DODGE] [%s] exausto, não pode esquivar ainda", c.Handle.String())
		return
	}
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		return
	}
	dirVec := target.GetPosition().Sub2D(c.Position).Normalize()
	perp := position.RotateVector2D(dirVec, math.Pi/2)
	if rand.Float64() < 0.5 {
		perp = perp.Multiply(-1)
	}
	newPos := c.Position.AddVector3D(position.Vector3D{X: perp.X, Y: 0, Z: perp.Z}.Multiply(c.DodgeDistance))

	if svcCtx.NavMesh.IsWalkable(newPos) {
		c.SetBlocking(false)
		c.MoveCtrl.SetMoveIntent(newPos, c.RunSpeed*1.5, 0.0)
		c.MoveCtrl.SetTarget(newPos, c.RunSpeed*1.5, 0.0)
		c.ReduceStamina(c.DodgeStaminaCost)
		c.invulnerableUntil = time.Now().Add(c.DodgeInvulnerabilityDuration)

		if c.Stamina <= 0 {
			c.DodgeDisabledUntil = time.Now().Add(2 * time.Second)
			c.RunSpeed *= 0.5
			log.Printf("[DODGE] [%s] exausto — esquiva desativada e velocidade reduzida", c.Handle.String())
		}
		log.Printf("[DODGE] [%s] esquiva bem-sucedida para (%.2f, %.2f) — stamina restante: %.2f",
			c.Handle.String(), newPos.X, newPos.Z, c.Stamina)
	} else {
		log.Printf("[DODGE] [%s] tentativa de esquiva falhou — destino não andável", c.Handle.String())
	}
}

func (c *Creature) CancelCurrentSkill() {
	state := c.SkillStates[c.NextSkillToUse.Action]
	if state != nil {
		state.InUse = false
		state.WasInterrupted = true
	}
	c.NextSkillToUse = nil // limpa a referência
}

func (c *Creature) ApplyPostureDamage(amount float64) {
	if c.PostureBroken || !c.Alive {
		return
	}

	c.Posture -= amount
	if c.Posture <= 0 {
		c.Posture = 0
		c.PostureBroken = true
		c.TimePostureBroken = time.Now().Unix()
		c.ChangeAIState(constslib.AIStateStaggered)
		c.CombatState = constslib.CombatStateStaggered
		// Exemplo de possível animação: c.SetAnimationState(constslib.AnimationIdle) ou custom para stagger
		// physics.StartStagger(&c.Stagger, c.TimePostureBroken, float64(c.PostureBreakDurationSec))
		log.Printf("[Creature %s] Posture quebrada! Entrando em stagger.", c.Handle.ID)
	}
}

func (c *Creature) TickPosture(deltaTime float64) {
	if c.PostureBroken {
		c.PostureBrokenElapsed += deltaTime
		if c.PostureBrokenElapsed >= c.PostureBreakDurationSec {
			c.PostureBroken = false
			c.Posture = c.MaxPosture
			c.PostureBrokenElapsed = 0
			c.ChangeAIState(constslib.AIStateIdle)
			c.CombatState = constslib.CombatStateRecovering
			log.Printf("[Creature %s] Posture recuperada. CombatState ajustado para Recovering.", c.Handle.String())
		}
	} else if c.Posture < c.MaxPosture {
		c.Posture += c.PostureRegenRate * deltaTime
		if c.Posture > c.MaxPosture {
			c.Posture = c.MaxPosture
		}
	}
}

func (c *Creature) TickStamina(deltaTime float64) {
	if time.Now().After(c.DodgeDisabledUntil) {
		if c.RunSpeed < c.OriginalRunSpeed {
			c.RunSpeed = c.OriginalRunSpeed
			log.Printf("[STAMINA] [%s] Velocidade restaurada após exaustão", c.Handle.String())
		}

		if c.Stamina < c.MaxStamina {
			regen := c.StaminaRegenPerSec * deltaTime
			c.Stamina += regen
			if c.Stamina > c.MaxStamina {
				c.Stamina = c.MaxStamina
			}
			// log.Printf("[STAMINA] [%s] Recuperação acelerada: %.2f / %.2f", c.Handle.String(), c.Stamina, c.MaxStamina)
		}
	} else {
		log.Printf("[STAMINA] [%s] Penalidade ativa, sem regeneração", c.Handle.String())
	}
}

func (c *Creature) TickEffects(deltaTime float64) {
	if !c.Alive {
		return
	}

	var remainingEffects []constslib.ActiveEffect
	ccActive := false

	for _, eff := range c.ActiveEffects {
		eff.Elapsed += deltaTime

		// Verifica expiração
		if eff.Elapsed >= eff.Duration.Seconds() {
			log.Printf("[Effect] Creature %s: efeito %s expirou.", c.Handle.String(), eff.Type)
			continue
		}

		// DOTs e cura
		if eff.Elapsed-eff.LastTickElapsed >= eff.TickInterval.Seconds() {
			switch eff.Type {
			case constslib.EffectPoison, constslib.EffectBurn, constslib.EffectBleed:
				c.HP -= eff.Power
				log.Printf("[Effect] Creature %s sofreu %d de %s. HP atual: %d", c.Handle.String(), eff.Power, eff.Type, c.HP)
				if c.HP <= 0 && c.Alive {
					c.ChangeAIState(constslib.AIStateDead)
					c.CombatState = constslib.CombatStateDead
					c.SetAnimationState(constslib.AnimationDie)
					log.Printf("[Effect] Creature %s morreu por efeito %s", c.Handle.String(), eff.Type)
				}
			case constslib.EffectRegen:
				c.HP += eff.Power
				if c.HP > c.Creature.MaxHP {
					c.HP = c.Creature.MaxHP
				}
				log.Printf("[Effect] Creature %s curou %d por %s. HP atual: %d", c.Handle.String(), eff.Power, eff.Type, c.HP)
			}

			eff.LastTickElapsed = eff.Elapsed
		}

		// Controle de CC
		if eff.IsCC {
			ccActive = true
		}

		remainingEffects = append(remainingEffects, eff)
	}

	// Atualiza estado de CC
	if ccActive {
		if c.CombatState != constslib.CombatStateStaggered {
			c.CombatState = constslib.CombatStateStaggered
			log.Printf("[Effect] Creature %s entrou em estado de CC (Staggered).", c.Handle.String())
		}
	} else if c.CombatState == constslib.CombatStateStaggered {
		c.CombatState = constslib.CombatStateRecovering
		log.Printf("[Effect] Creature %s saiu de CC e está se recuperando.", c.Handle.String())
	}

	c.ActiveEffects = remainingEffects
}

func (c *Creature) TickNeeds(deltaTime float64) {
	// Necessidades fisiológicas (proporcional ao tempo)
	ModifyNeed(c, constslib.NeedHunger, 0.007*deltaTime)
	ModifyNeed(c, constslib.NeedThirst, 0.008*deltaTime)
	ModifyNeed(c, constslib.NeedSleep, 0.004*deltaTime)

	// Tendência de estabilidade: puxa valores para o ponto médio entre Min e Threshold
	for i := range c.Needs {
		n := &c.Needs[i]
		middle := (n.LowThreshold + n.Threshold) / 2
		var delta float64

		if n.Value < middle {
			// Tendência de crescer
			if n.Type == constslib.NeedAdvance || n.Type == constslib.NeedGuard {
				delta = rand.Float64() * 0.1
			} else {
				delta = rand.Float64() * 0.05
			}
		} else {
			// Tendência de reduzir
			if n.Type == constslib.NeedAdvance || n.Type == constslib.NeedGuard {
				delta = -(rand.Float64() * 0.1)
			} else {
				delta = -(rand.Float64() * 0.05)
			}
		}

		ModifyNeed(c, n.Type, delta*deltaTime) // << aqui também
	}

	// Clamp em 0-100
	for _, t := range []constslib.NeedType{
		constslib.NeedAdvance, constslib.NeedGuard, constslib.NeedRetreat,
		constslib.NeedProvoke, constslib.NeedRecover, constslib.NeedPlan,
		constslib.NeedFake, constslib.NeedRage,
	} {
		val := c.GetNeedValue(t)
		if val < 0 {
			c.SetNeedValue(t, 0)
		}
		if val > 100 {
			c.SetNeedValue(t, 100)
		}
	}
}

func (c *Creature) SetAnimationState(state constslib.AnimationState) {
	if c.AnimationState != state {
		c.AnimationState = state
		// log.Printf("[Creature] %s (%s) animação definida para %s", c.Handle.String(), c.PrimaryType, state)
	}
}

func (c *Creature) ChangeAIState(newState constslib.AIState) {
	if c.AIState == newState {
		return
	}

	// log.Printf("[Creature] %s (%s) AI State mudou: %s → %s", c.Handle.String(), c.PrimaryType, c.AIState, newState)
	c.AIState = newState
	c.LastStateChange = time.Now()

	if newState == constslib.AIStateDead {
		c.Alive = false
		c.IsCorpse = true
		c.TimeOfDeath = time.Now()
		c.SetAnimationState(constslib.AnimationDie)
	} else if newState == constslib.AIStateIdle {
		c.SetAnimationState(constslib.AnimationIdle)
	}
}

func (c *Creature) IsQuestOnly() bool {
	return strings.TrimSpace(c.Creature.OwnerPlayerID) != ""
}

func (c *Creature) ApplyEffect(effect constslib.ActiveEffect) {
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

func (c *Creature) Respawn(navMesh *navmesh.NavMesh) {
	c.HP = c.Creature.MaxHP
	c.Alive = true
	c.IsCorpse = false
	c.AIState = constslib.AIStateIdle
	c.CombatState = constslib.CombatStateIdle
	c.AnimationState = constslib.AnimationIdle
	c.TimeOfDeath = time.Time{}
	c.LastThreatSeen = time.Time{}
	c.ClearAggro()
	c.ClearCooldowns()

	pos := c.GenerateSpawnPosition(navMesh)
	c.SetPosition(pos)

	c.Generation++
	c.Handle = handle.NewEntityHandle(lib.NewUUID(), c.Generation)
}

func (c *Creature) ClearCooldowns() {
	for _, state := range c.SkillStates {
		if state != nil {
			state.CooldownUntil = time.Time{}
			state.InUse = false
			// Opcional: zera outras fases
			state.StartedAt = time.Time{}
			state.WindUpUntil = time.Time{}
			state.CastUntil = time.Time{}
			state.RecoveryUntil = time.Time{}
			state.WindUpFired = false
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

	// Exemplo AAA: ajuste AIState ao receber aggro relevante
	if c.AIState == constslib.AIStateIdle && amount > 0 {
		c.ChangeAIState(constslib.AIStateAlert)
	}
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

func (c *Creature) IsAlive() bool {
	return c.Alive
}

func (c *Creature) IsHostile() bool {
	return c.Hostile
}

func (c *Creature) GetCurrentSpeed() float64 {
	switch c.AnimationState {
	case constslib.AnimationRun:
		return c.RunSpeed
	case constslib.AnimationWalk:
		return c.WalkSpeed
	default:
		return 0
	}
}

func (c *Creature) TakeDamage(amount int) {
	if !c.Alive {
		return
	}

	finalDamage := int(math.Round(float64(amount) * (1.0 - c.PhysicalDefense)))
	if finalDamage <= 0 {
		finalDamage = 1
	}

	c.HP -= finalDamage
	c.LastAttackedTime = time.Now()
	c.ReceivedDamageRecently = true

	log.Printf("[Creature %s] sofreu %d de dano. HP restante: %d", c.Handle.String(), finalDamage, c.HP)

	// ⚔️ Bloqueio reflexo caso o dano seja brutal
	if !c.IsBlocking() && !c.PostureBroken {
		maxHP := float64(c.Creature.MaxHP)
		if maxHP <= 0 {
			maxHP = 1
		}
		percent := float64(finalDamage) / maxHP

		if percent >= 0.35 {
			c.SetBlocking(true)
			c.BlockStartedAt = time.Now()
			c.CombatState = constslib.CombatStateBlocking
			c.ReduceStamina(10.0)
			log.Printf("[BLOCK-REFLEX] [%s] ativou bloqueio reflexo após dano brutal (%.1f%% HP)",
				c.Handle.String(), percent*100)
		}
	}

	if c.HP <= 0 {
		c.ChangeAIState(constslib.AIStateDead)
		c.CombatState = constslib.CombatStateDead
		c.SetAnimationState(constslib.AnimationDie)
		log.Printf("[Creature %s] morreu após receber dano.", c.Handle.String())
	} else {
		c.CombatState = constslib.CombatStateRecovering
	}
}

func (c *Creature) GetFacingDirection() position.Vector2D {
	return c.FacingDirection
}

// creature/creature.go

func (c *Creature) GetBestTargetFromTargets(targets []model.Targetable) model.Targetable {
	var best model.Targetable
	var minDist float64 = math.MaxFloat64

	for _, t := range targets {
		if t == nil || !t.IsAlive() {
			continue
		}
		if t.GetHandle().Equals(c.Handle) {
			continue
		}
		dist := position.CalculateDistance(c.Position, t.GetPosition())
		if dist < minDist {
			minDist = dist
			best = t
		}
	}

	c.LastKnownDistance = minDist // Atualiza a distância conhecida
	return best
}

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
		if !p.IsAlive() {
			continue
		}
		dist := position.CalculateDistance(c.Position, p.Position)
		if dist < minDist {
			minDist = dist
			closest = p
		}
	}

	for _, c2 := range creatures {
		if c2.Handle.Equals(c.Handle) || !c2.Alive {
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
		if n.Type == constslib.NeedHunger && n.Value >= n.Threshold {
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

func (c *Creature) ConsumeCorpse(index navmesh.SpatialIndex) {
	c.IsCorpse = false
	index.Remove(c)
	log.Printf("[Creature] %s (%s) foi consumido e removido do SpatialIndex.", c.Handle.String(), c.PrimaryType)
}

func (c *Creature) IsCreature() bool { return true }

func (c *Creature) IsStaticObstacle() bool {
	return false // criaturas nunca são obstáculo absoluto
}

func (c *Creature) InitSkillState(action constslib.SkillAction, now time.Time) *model.SkillState {
	windupUntil := now.Add(time.Duration(c.NextSkillToUse.WindUpTime * float64(time.Second)))
	castUntil := windupUntil.Add(time.Duration(c.NextSkillToUse.CastTime * float64(time.Second)))
	recoveryUntil := castUntil.Add(time.Duration(c.NextSkillToUse.RecoveryTime * float64(time.Second)))

	state := &model.SkillState{
		InUse:            true,
		StartedAt:        now,
		WindUpUntil:      windupUntil,
		CastUntil:        castUntil,
		RecoveryUntil:    recoveryUntil,
		CooldownUntil:    now.Add(time.Duration(c.NextSkillToUse.CooldownSec * float64(time.Second))), // AQUI
		HasCastBeenFired: false,
		WindUpFired:      false,
	}

	c.SkillStates[action] = state
	return state
}

func (c *Creature) ResetSkillState(action constslib.SkillAction) {
	state, exists := c.SkillStates[action]
	if !exists || state == nil {
		// Nada a resetar
		return
	}

	// Reseta flags do ciclo
	state.InUse = false
	state.HasCastBeenFired = false
	state.StartedAt = time.Time{}
	state.WindUpUntil = time.Time{}
	state.CastUntil = time.Time{}
	state.RecoveryUntil = time.Time{}
	state.WindUpFired = false

	// Limpa movimento associado
	c.SkillMovementState = nil
}

func (c *Creature) ClearMovementIntent() {
	if c.MoveCtrl != nil {
		c.MoveCtrl.Intent = movement.MoveIntent{} // Zera o intent
		c.MoveCtrl.IsMoving = false
		c.MoveCtrl.CurrentPath = nil
		c.MoveCtrl.PathIndex = 0
		c.MoveCtrl.WasBlocked = false

		// Log opcional:
		// log.Printf("[MOVEMENT] [%s (%s)] Movimento pendente e path limpos", c.Handle.String(), c.PrimaryType)
	}
}

func (c *Creature) GetStrength() int {
	return c.Strength
}

func (c *Creature) GetDexterity() int {
	return c.Dexterity
}

func (c *Creature) GetIntelligence() int {
	return c.Intelligence
}

func (c *Creature) GetFocus() int {
	return c.Focus
}

func (c *Creature) GetPrimaryType() string {
	return string(c.PrimaryType)
}

func (c *Creature) GetSkillMovementState() *model.SkillMovementState {
	return c.SkillMovementState
}

func (c *Creature) SetSkillMovementState(state *model.SkillMovementState) {
	c.SkillMovementState = state
}

func (c *Creature) IsPvPEnabled() bool {
	return false
}

func (c *Creature) IsDodging() bool {
	return c.dodging
}

func (c *Creature) SetDodging(dodging bool) {
	c.dodging = dodging
}

func (c *Creature) IsBlocking() bool {
	return c.blocking
}

func (c *Creature) SetBlocking(blocking bool) {
	c.blocking = blocking
}

func (c *Creature) GetCurrentTarget(svcCtx *dynamic_context.AIServiceContext) model.Targetable {
	if c.TargetCreatureHandle.IsValid() {
		return svcCtx.FindByHandle(c.TargetCreatureHandle)
	}
	if c.TargetPlayerHandle.IsValid() {
		return svcCtx.FindByHandle(c.TargetPlayerHandle)
	}
	return nil
}

func (c *Creature) ApplyDodgeExhaustionPenalty() {
	c.DodgeDisabledUntil = time.Now().Add(2 * time.Second)
	c.RunSpeed *= 0.5
	log.Printf("[DODGE] [%s] exausto após esquiva — dodge desativado e velocidade reduzida", c.Handle.String())
}

func (c *Creature) IsInvulnerableNow() bool {
	return time.Now().Before(c.invulnerableUntil)
}

func (c *Creature) ReduceStamina(amount float64) {
	c.Stamina -= amount
	if c.Stamina < 0 {
		c.Stamina = 0
	}
}

func generateBlockDuration(combatState constslib.CombatState) time.Duration {
	var base float64

	switch combatState {
	case constslib.CombatStateDefensive:
		// Tendência mais longa: 1.0 a 3.3s
		base = 1.0 + rand.Float64()*2.3
	case constslib.CombatStateStrategic:
		// Tendência média: 0.5 a 2.3s
		base = 0.5 + rand.Float64()*1.8
	case constslib.CombatStateAggressive:
		// Tendência curta: 0.5 a 1.5s
		base = 0.5 + rand.Float64()*1.0
	default:
		// Igual ao Strategic
		base = 0.5 + rand.Float64()*1.8
	}

	return time.Duration(base * float64(time.Second))
}

// OnBlockHit aplica dano de postura ao bloquear um golpe
func OnBlockHit(c *Creature, postureDamage float64) {
	finalPostureDamage := postureDamage * 2 // Dano dobrado na postura ao bloquear
	c.ApplyPostureDamage(finalPostureDamage)
	log.Printf("[BLOCK] [%s] bloqueou ataque, aplicou %.1f posture damage (dobrado)", c.Handle.String(), finalPostureDamage)
}

func (c *Creature) IsInParryWindow() bool {
	now := time.Now()
	return now.After(c.ParryWindowStart) && now.Before(c.ParryWindowEnd)
}

func (c *Creature) CurrentSkillState() *model.SkillState {
	if c.NextSkillToUse == nil {
		return nil
	}
	return c.SkillStates[c.NextSkillToUse.Action]
}
