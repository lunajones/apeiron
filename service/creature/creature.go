package creature

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/lunajones/apeiron/lib"
	combatlib "github.com/lunajones/apeiron/lib/combat"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/movement/fsm"
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

	facingDirection position.Vector2D
	torsoDirection  position.Vector2D

	LastThreatSeen time.Time

	AggroTable        map[handle.EntityHandle]*aggro.AggroEntry
	LastKnownDistance float64

	ParryWindowStart      time.Time
	ParryWindowEnd        time.Time
	BlockStaminaTolerance float64

	BlockableChance float64 // ex: 0.7 (70%)
	DodgableChance  float64 // ex: 0.3 (30%)

	lastDodgeEvent      model.CombatEvent
	cachedDodgePosition position.Position
	combatDrive         model.CombatDrive
	combatEvents        []model.CombatEvent
	lastAggressionEvent model.CombatEvent

	movementLockedUntil time.Time
	recentActions       []constslib.CombatAction

	casting            bool
	LastSkillPlannedAt time.Time
	LastDodgeDirection position.Vector2D

	lastMissedSkillAt time.Time
	lastCircleAt      time.Time
	lastRetreatAt     time.Time

	context *dynamic_context.AIServiceContext

	movementFSM fsm.MovementFSM
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

func (c *Creature) GetLastDodgeEvent() model.CombatEvent {
	return c.lastDodgeEvent
}

func (c *Creature) SetLastDodgeEvent(event model.CombatEvent) {
	c.lastDodgeEvent = event
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
	c.facingDirection = dir
}

func (c *Creature) GetFacingDirection() position.Vector2D {
	return c.facingDirection
}

func (c *Creature) SetTorsoDirection(dir position.Vector2D) {
	c.torsoDirection = dir
}

func (c *Creature) GetTorsoDirection() position.Vector2D {
	return c.torsoDirection
}

var creatures []*Creature

func (c *Creature) Tick(ctx *dynamic_context.AIServiceContext, deltaTime float64) {
	if !c.Alive {
		return
	}

	c.SetContext(ctx)

	// FSM de casting
	if c.CombatState == constslib.CombatStateCasting {
		skill := c.NextSkillToUse

		combatlib.ProcessCastingFSM(&combatlib.CastingConfig{
			Creature: c,
			Target:   finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, ctx),
			Skill:    skill,
			State:    c.SkillStates[skill.Action],
			Now:      time.Now(),
			Allies:   finder.FindNearbyAllies(ctx, c, c.GetFaction(), 8.0),
		}, ctx)
	}

	c.UpdateFacingDirection(ctx)

	// Defesa (block/dodge automático)
	c.PerformDefensiveAction(deltaTime)

	// 1️⃣ Movimento por habilidade (leap, charge, etc)
	if c.SkillMovementState != nil && c.SkillMovementState.Active {
		target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, ctx)
		if combatlib.UpdateSkillMovement(c, c.SkillMovementState, target, ctx.NavMesh, ctx, deltaTime) {
			log.Printf("[LEAP] [%s] SkillMovement concluído", c.Handle.String())
			c.SkillMovementState = nil
			c.CombatState = constslib.CombatStateIdle
		}
	}

	// 2️⃣ Movimento normal
	if c.SkillMovementState == nil || !c.SkillMovementState.Active {
		c.MoveCtrl.Update(c, deltaTime, ctx)
	}

	if c.IsDodging() && !c.MoveCtrl.IsImpulsing() {
		c.SetDodging(false)
	}

	// 3️⃣ FSM de movimento (decisão automática de estado)
	if c.GetMovementFSM() == nil {
		fsm := fsm.ProcessMovementFSM(fsm.FSMHooks{
			OnClearIntent:      c.ClearMovementIntent,
			OnSetAnimation:     c.SetAnimationState,
			OnKnockbackImpulse: func() { c.MoveCtrl.ForceImpulseAwayFromTarget(2.0) },
			OnHasArrived:       func() bool { return c.MoveCtrl.HasArrived(c) },
			OnIsDodging:        c.IsDodging,
			OnSetImpulse: func(state *movement.ImpulseMovementState) {
				c.MoveCtrl.ImpulseState = state
			},
			OnShouldWalk: func() bool {
				return c.MoveCtrl.IsMoving // ou qualquer flag que você tenha
			},
			OnIsCasting: func() bool {
				return c.IsCasting() // ou outro método seu que detecta WindUp/Cast/Recovery
			},
		})
		c.SetMovementFSM(fsm)
	}

	c.GetMovementFSM().Tick(deltaTime)

	// 4️⃣ Ticks de status
	c.TickNeeds(deltaTime)
	c.TickEffects(deltaTime)
	c.TickPosture(deltaTime)
	c.TickStamina(deltaTime)

	// 5️⃣ Feedback de combate
	c.ProcessCombatFeedback()

	// 6️⃣ AI
	if c.BehaviorTree != nil {
		c.BehaviorTree.Tick(c, ctx)
	}
}

func (c *Creature) PerformDefensiveAction(deltaTime float64) {
	if !c.Alive {
		return
	}

	total := c.BlockableChance + c.DodgableChance
	if total == 0 {
		return
	}

	now := time.Now()
	cutoff := now.Add(-1 * time.Second)

	baseMargin := 120 * time.Millisecond
	randFactor := 0.8
	stateReactionFactor := 1.0
	adjustedMargin := time.Duration(stateReactionFactor * float64(baseMargin))

	var best model.CombatEvent
	var found bool
	var bestStart time.Time

	for _, e := range c.GetRecentCombatEvents(cutoff) {
		if e.BehaviorType != "AggressiveIntention" {
			continue
		}

		castDuration := e.ExpectedImpact.Sub(e.Timestamp)
		randomMargin := time.Duration(rand.Float64() * randFactor * float64(castDuration))

		start := e.ExpectedImpact.Add(-adjustedMargin).Add(-randomMargin)
		end := e.ExpectedImpact

		if now.After(start) && now.Before(end) {
			if !found || start.Before(bestStart) {
				best = e
				bestStart = start
				found = true
			}
		} else if now.After(end) {
			log.Printf("[DEFENSE] [%s] janela passou sem defesa (%s → %s), agora: %s",
				c.Handle.String(), start.Format("15:04:05.000"), end.Format("15:04:05.000"), now.Format("15:04:05.000"))
			c.ConsumeCombatEvent(e)
		}
	}

	if found {
		r := rand.Float64() * total

		log.Printf("[DEFENSE] [%s] total=%.2f | dodge=%.2f | block=%.2f | roll=%.2f",
			c.GetPrimaryType(), total, c.DodgableChance, c.BlockableChance, r)
		if r < c.DodgableChance {
			if c.TryDodgeReaction(best) {
				c.ConsumeCombatEvent(best)
			}
		} else {
			if c.TryBlockReaction(best) {
				c.ConsumeCombatEvent(best)
			}
		}
	}

	c.PerformBlock(deltaTime)
}

func (c *Creature) TryBlockReaction(e model.CombatEvent) bool {
	if !c.Alive || c.PostureBroken || c.IsBlocking() || c.IsDodging() {
		return false
	}

	if !c.CurrentSkillState().CanBeCancelled() {
		return false
	} else {
		c.CancelCurrentSkill()
	}

	c.SetBlocking(true)
	c.SetDodging(false)
	c.BlockStartedAt = time.Now()
	log.Printf("[REACT] [%s] iniciou bloqueio contra %s", c.Handle.String(), e.SourceHandle.ID)
	return true
}

func (c *Creature) PerformBlock(deltaTime float64) {
	if !c.IsBlocking() {
		c.BlockSpentStamina = 0
		return
	}

	now := time.Now()

	if c.BlockDuration == 0 {
		c.ReduceStamina(10.0)

		drive := c.GetCombatDrive()
		parryDuration := 300 * time.Millisecond
		switch {
		case drive.Caution > 0.8:
			parryDuration = 450 * time.Millisecond
		case drive.Caution > 0.6:
			parryDuration = 380 * time.Millisecond
		case drive.Caution > 0.4:
			parryDuration = 320 * time.Millisecond
		case drive.Caution > 0.2:
			parryDuration = 260 * time.Millisecond
		}

		c.ParryWindowStart = now
		c.ParryWindowEnd = now.Add(parryDuration)
		log.Printf("[PARRY] [%s] parry window ativa por %v", c.Handle.String(), parryDuration)

		c.BlockStaminaTolerance = 5.0 + rand.Float64()*5.0

		base := 2.5
		variation := rand.Float64() * 2.5
		c.MaxBlockDuration = time.Duration((base + variation) * float64(time.Second))
	}

	c.BlockDuration += time.Duration(deltaTime * float64(time.Second))
	staminaPerSecond := 1.0
	staminaThisTick := staminaPerSecond * deltaTime * 10
	c.BlockSpentStamina += staminaThisTick
	c.ReduceStamina(staminaThisTick)
	c.SetDodging(false)

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

func (c *Creature) TryDodgeReaction(e model.CombatEvent) bool {
	if !c.Alive || c.IsDodging() || c.IsBlocking() {
		return false
	}

	if last := c.GetLastDodgeEvent(); last.BehaviorType != "" {
		return false
	}

	if c.Stamina < c.DodgeStaminaCost+5.0 {
		log.Printf("[DODGE] [%s] recusou esquiva — stamina insuficiente (%.2f)", c.GetPrimaryType(), c.Stamina)
		return false
	}

	drive := c.GetCombatDrive()
	staminaRatio := c.Stamina / c.MaxStamina
	baseChance := 0.75
	switch {
	case drive.Caution > 0.6 || drive.Counter > 0.6:
		baseChance = 1.0
	case drive.Caution > 0.4 || drive.Counter > 0.4:
		baseChance = 0.85
	case drive.Caution > 0.2 || drive.Counter > 0.2:
		baseChance = 0.75
	default:
		baseChance = 0.5
	}
	finalChance := baseChance * (0.3 + 0.7*staminaRatio)
	if rand.Float64() >= finalChance {
		log.Printf("[DODGE] [%s] falhou — finalChance %.2f | Caution=%.2f | Counter=%.2f | StaminaRatio=%.2f",
			c.GetPrimaryType(), finalChance, drive.Caution, drive.Counter, staminaRatio)
		c.ConsumeCombatEvent(e)
		return false
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		return false
	}

	dirVec := target.GetPosition().Sub2D(c.Position).Normalize()
	perp := position.RotateVector2D(dirVec, math.Pi/2).Normalize()
	back := dirVec.Multiply(-1).Normalize()

	applyJitter := func(v position.Vector2D) position.Vector2D {
		angle := (rand.Float64() - 0.5) * 0.3
		return position.RotateVector2D(v, angle).Normalize().Multiply(c.DodgeDistance)
	}

	candidateOffsets := []position.Vector2D{
		applyJitter(perp),
		applyJitter(perp.Multiply(-1)),
		applyJitter(back),
	}
	rand.Shuffle(len(candidateOffsets), func(i, j int) {
		candidateOffsets[i], candidateOffsets[j] = candidateOffsets[j], candidateOffsets[i]
	})

	var chosen position.Position
	var chosenVec position.Vector2D
	found := false
	for _, offset := range candidateOffsets {
		newPos := c.Position.AddVector3D(position.Vector3D{X: offset.X, Y: 0, Z: offset.Z})
		if c.GetContext().NavMesh.IsWalkable(newPos) && c.GetContext().ClaimPosition(newPos, c.Handle) {
			chosen = newPos
			chosenVec = offset.Normalize()
			found = true
			break
		}
	}
	if !found {
		log.Printf("[DODGE] [%s] negado — sem destino andável", c.GetPrimaryType())
		return false
	}

	if !c.CurrentSkillState().CanBeCancelled() {
		return false
	}
	c.CancelCurrentSkill()

	c.ConsumeCombatEvent(e)

	// ✅ Ativa FSM de movimento
	c.GetMovementFSM().EnterDodgingState(c, chosen, chosenVec, e)

	log.Printf("[REACT] [%s] esquiva iniciada para (%.2f, %.2f) contra %s — stamina: %.2f",
		c.GetPrimaryType(), chosen.X, chosen.Z, e.SourceHandle.ID, c.Stamina)

	return true
}

func (c *Creature) EndCurrentSkill() {
	if c.NextSkillToUse == nil {
		return
	}

	state := c.SkillStates[c.NextSkillToUse.Action]
	if state != nil {
		state.InUse = false

		// Reset flags de execução
		state.WindUpFired = false
		state.CastFired = false
		state.RecoveryFired = false

		// Reset tempos transitórios
		state.WindUpUntil = time.Time{}
		state.CastUntil = time.Time{}
		state.RecoveryUntil = time.Time{}
		// ⚠️ Mantém CooldownUntil!
	}

	c.SetCombatState(constslib.CombatStateMoving)
	c.NextSkillToUse = nil
}

func (c *Creature) CancelCurrentSkill() {
	if c.NextSkillToUse == nil {
		return
	}

	state := c.SkillStates[c.NextSkillToUse.Action]
	if state != nil {
		state.InUse = false
		state.WasInterrupted = true

		state.WindUpFired = false
		state.CastFired = false
		state.RecoveryFired = false
		state.WindUpUntil = time.Time{}
		state.CastUntil = time.Time{}
		state.RecoveryUntil = time.Time{}
		state.CooldownUntil = time.Time{}
	}
	c.SetCombatState(constslib.CombatStateMoving)

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

	c.AddRecentAction(constslib.CombatActionTookDamage)

	log.Printf("[Creature %s] sofreu %d de dano. HP restante: %d", c.Handle.String(), finalDamage, c.HP)

	if c.HP <= 0 {
		c.ChangeAIState(constslib.AIStateDead)
		c.CombatState = constslib.CombatStateDead
		c.SetAnimationState(constslib.AnimationDie)
		log.Printf("[Creature %s] morreu após receber dano.", c.Handle.String())
	}
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

	state := &model.SkillState{
		StartedAt:     now,
		InUse:         true,
		WindUpFired:   false,
		CastFired:     false,
		RecoveryFired: false,
	}

	c.SkillStates[action] = state
	return state
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

func (c *Creature) GetCachedDodgePosition() position.Position {
	return c.cachedDodgePosition
}

func (c *Creature) SetCachedDodgePosition(pos position.Position) {
	c.cachedDodgePosition = pos
}

func (c *Creature) GetCombatDrive() *model.CombatDrive {
	return &c.combatDrive
}

func (c *Creature) SetCombatDrive(drive model.CombatDrive) {
	c.combatDrive = drive
}

func (c *Creature) RegisterCombatEvent(event model.CombatEvent) {
	c.combatEvents = append(c.combatEvents, event)
}

func (c *Creature) GetRecentCombatEvents(since time.Time) []model.CombatEvent {
	var result []model.CombatEvent
	for _, e := range c.combatEvents {
		if e.Timestamp.After(since) {
			result = append(result, e)
		}
	}
	return result
}

func (c *Creature) ClearOldCombatEvents(before time.Time) {
	var filtered []model.CombatEvent
	for _, e := range c.combatEvents {
		if e.Timestamp.After(before) {
			filtered = append(filtered, e)
		}
	}
	c.combatEvents = filtered
}

func (c *Creature) AddCaution(delta float64) {
	drive := c.GetCombatDrive()
	drive.Caution += delta
	if drive.Caution > 1.0 {
		drive.Caution = 1.0
	}
}

func (c *Creature) RegisterAggressionFrom(source handle.EntityHandle, now time.Time) {
	c.lastAggressionEvent = model.CombatEvent{
		SourceHandle: source,
		BehaviorType: "AggressiveIntention",
		Timestamp:    now,
	}
}

func (c *Creature) RecalculateDrive() {
	drive := c.GetCombatDrive()
	// log.Printf("[SKILL-STATE] [%s] CombatDrive: Rage=%.2f Termination=%.2f Value=%.2f",
	// 	c.Handle.String(), drive.Rage, drive.Termination, drive.Value)
	drive.Value = RecalculateCombatDrive(drive)
}

// RecalculateCombatDrive aplica regras para consolidar os componentes em um único valor entre 0.0 e 1.0
func RecalculateCombatDrive(d *model.CombatDrive) float64 {
	// Normaliza os componentes
	rage := clamp01(d.Rage)               // motivação por dor pessoal
	caution := clamp01(d.Caution)         // medo ou precaução
	vengeance := clamp01(d.Vengeance)     // perdas de aliados
	termination := clamp01(d.Termination) // falta de combate recente

	// Ponderação comportamental:
	// Rage empurra o combate → impulsivo
	// Caution segura → estratégico/defensivo
	// Vengeance aumenta persistência em combate
	// Termination busca estímulo → hostilidade aumentada com tédio

	// Ajuste mais comportamental: Caution reduz o valor final
	raw := 0.0
	raw += rage * 0.4
	raw += vengeance * 0.25
	raw += termination * 0.25
	raw -= caution * 0.3 // o medo pode segurar a criatura, até paralisá-la

	return clamp01(raw)
}

func clamp01(value float64) float64 {
	if value < 0.0 {
		return 0.0
	}
	if value > 1.0 {
		return 1.0
	}
	return value
}

func (c *Creature) RemoveCombatEventAt(index int) {
	if index < 0 || index >= len(c.combatEvents) {
		log.Printf("[REMOVE-COMBAT-EVENT] Índice inválido: %d", index)
		return
	}
	c.combatEvents = append(c.combatEvents[:index], c.combatEvents[index+1:]...)
}

func (c *Creature) GetCombatEvents() []model.CombatEvent {
	return c.combatEvents
}

// IsMovementLocked retorna true se a criatura ainda estiver com o movimento travado
func (c *Creature) IsMovementLocked() bool {
	return time.Now().Before(c.movementLockedUntil)
}

// SetMovementLock define um tempo de travamento de movimentação (lock)
func (c *Creature) SetMovementLock(duration time.Duration) {
	c.movementLockedUntil = time.Now().Add(duration)
}

func (c *Creature) GetMovementLockUntil() time.Time {
	return c.movementLockedUntil
}

func (c *Creature) UpdateFacingDirection(ctx *dynamic_context.AIServiceContext) {
	// 1. Se está se movendo (usa MoveCtrl e Velocity), orienta na direção da velocidade
	if c.MoveCtrl != nil && c.MoveCtrl.IsMoving {
		dir := c.MoveCtrl.Velocity.ToVector2D().Normalize()
		if dir.X != 0 || dir.Z != 0 {
			c.facingDirection = dir
			return
		}
	}

	// 2. Se tem um alvo válido, tenta encontrar e mirar
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, ctx)
	if target != nil {
		toTarget := position.NewVector2DFromTo(c.Position, target.GetPosition())
		if toTarget.X != 0 || toTarget.Z != 0 {
			c.facingDirection = toTarget
			return
		}
	}

	// 3. Caso ainda esteja zerado (no início, por exemplo), seta um default
	if c.GetFacingDirection().X == 0 && c.GetFacingDirection().Z == 0 {
		c.facingDirection = position.Vector2D{X: 0, Z: 1}
	}
}

func (c *Creature) UpdateTorsoDirection() {
	if c.MoveCtrl != nil && c.MoveCtrl.IsMoving {
		dir := c.MoveCtrl.Velocity.ToVector2D().Normalize()
		if dir.Length() > 0.01 {
			c.SetTorsoDirection(dir)
			return
		}
	}

	// Se não estiver se movendo, mantém o último valor
}

func (c *Creature) ConsumeCombatEvent(target model.CombatEvent) {
	newEvents := make([]model.CombatEvent, 0, len(c.combatEvents))
	for _, e := range c.combatEvents {
		if e != target {
			newEvents = append(newEvents, e)
		}
	}
	c.combatEvents = newEvents
}

func (c *Creature) IsSkillAvailable(skillID string) bool {
	action := constslib.SkillAction(skillID)
	state, ok := c.SkillStates[action]
	if !ok || state == nil {
		return true // nunca usada = disponível
	}
	return time.Now().After(state.CooldownUntil)
}

func (c *Creature) IsCasting() bool {
	if c.NextSkillToUse == nil {
		return false
	}
	state, ok := c.SkillStates[c.NextSkillToUse.Action]
	if !ok || state == nil {
		return false
	}
	return state.WindUpFired && !state.RecoveryFired
}

func (c *Creature) SetCasting(casting bool) {
	c.casting = casting
}

func (c *Creature) GetFaction() string {
	return c.Faction
}

func (c *Creature) ProcessCombatFeedback() {
	drive := c.GetCombatDrive()
	seen := make(map[constslib.CombatAction]bool)

	for _, action := range c.GetRecentActions() {
		if seen[action] {
			continue
		}
		seen[action] = true

		switch action {
		case constslib.CombatActionBlockSuccess:
			drive.Rage += 0.1
			drive.Caution -= 0.03
			drive.Termination += 0.03
			drive.Counter += 0.2

		case constslib.CombatActionParrySuccess:
			drive.Rage += 0.1
			drive.Caution -= 0.06
			drive.Termination += 0.45
			drive.Counter += 0.45

		case constslib.CombatActionDodgeSuccess:
			drive.Rage += 0.1
			drive.Caution -= 0.1
			drive.Termination += 0.03
			drive.Counter += 0.2

		case constslib.CombatActionMicroRetreat:
			drive.Rage -= 0.02
			drive.Caution += 0.04
			drive.Termination += 0.01
			drive.Counter += 0.2

		case constslib.CombatActionCircleAround:
			drive.Rage += 0.005
			drive.Caution -= 0.005
			drive.Termination += 0.01
			drive.Counter += 0.005

		case constslib.CombatActionAttackWindow:
			drive.Rage += 0.05
			drive.Caution -= 0.05
			drive.Termination += 0.05
			drive.Counter += 0.05

		case constslib.CombatActionApproach:
			drive.Rage += 0.02
			drive.Caution -= 0.02
			drive.Termination += 0.01
			drive.Counter += 0.05

		case constslib.CombatActionChase:
			drive.Rage += 0.03
			drive.Caution -= 0.02
			drive.Termination += 0.03
			drive.Counter += 0.05

		case constslib.CombatActionAttackPrepared:
			drive.Rage -= 0.005
			drive.Caution += 0.005
			drive.Termination += 0.005
			drive.Counter -= 0.001

		case constslib.CombatActionAttackSuccess:
			drive.Rage += 0.08
			drive.Caution -= 0.02
			drive.Termination += 0.06
			drive.Counter += 0.05

		case constslib.CombatActionAttackMissed:
			drive.Rage += 0.02
			drive.Caution += 0.04
			drive.Termination += 0.07
			drive.Counter += 0.07

		case constslib.CombatActionSkillInterrupted:
			drive.Rage -= 0.05
			drive.Caution += 0.06
			drive.Termination -= 0.01
			drive.Counter += 0.12

		case constslib.CombatActionCounter:
			drive.Rage += 0.01
			drive.Caution -= 0.01
			drive.Termination -= 0.01
			drive.Counter = 0.0

		case constslib.CombatActionTookDamage:
			drive.Rage += 0.07
			drive.Caution += 0.12
			drive.Termination += 0.01
			drive.Counter += 0.24

		case constslib.CombatActionHesitatedAttack:
			drive.Caution += 0.12
			drive.Rage -= 0.04
			drive.Termination -= 0.01
			drive.Counter -= 0.04

		}
	}

	// Decaimento
	drive.Rage *= 0.995
	drive.Caution *= 0.995
	drive.Termination *= 0.995
	drive.Counter *= 0.97

	// Clamp entre 0.0 e 1.0
	drive.Rage = math.Max(0, math.Min(1, drive.Rage))
	drive.Caution = math.Max(0, math.Min(1, drive.Caution))
	drive.Termination = math.Max(0, math.Min(1, drive.Termination))
	drive.Counter = math.Max(0, math.Min(1, drive.Counter))

	// Log visual
	// color.New(color.FgHiMagenta, color.Bold).Printf(
	// 	"[COMBAT-FEEDBACK] [%s] drive atualizado: Rage=%.2f | Caution=%.2f | Termination=%.2f | Counter=%.2f\n",
	// 	c.PrimaryType, drive.Rage, drive.Caution, drive.Termination, drive.Counter,
	// )

	// Limpa ações
	c.ClearRecentActions()
}

func (c *Creature) GetYaw() float64 {
	dir := c.GetFacingDirection() // position.Vector2D

	if dir.X == 0 && dir.Z == 0 {
		return 0 // parado, usa default
	}

	angle := math.Atan2(dir.X, dir.Z) * (180 / math.Pi)
	if angle < 0 {
		angle += 360
	}
	return angle
}

func (c *Creature) SetLastMissedSkillAt(t time.Time) {
	c.lastMissedSkillAt = t
}

func (c *Creature) HasRecentlyMissedSkill() bool {
	return time.Since(c.lastMissedSkillAt) < 2*time.Second
}

func (c *Creature) CanCircleAgain(cooldown time.Duration) bool {
	return time.Since(c.lastCircleAt) >= cooldown
}

func (c *Creature) GetLastSkillMissedAt() time.Time {
	return c.lastMissedSkillAt
}

func (c *Creature) GetLastCircleAt() time.Time {
	return c.lastCircleAt
}

func (c *Creature) SetLastCircleAt(t time.Time) {
	c.lastCircleAt = t
}

func (c *Creature) CanRetreatAgain(cooldown time.Duration) bool {
	return time.Since(c.lastRetreatAt) >= cooldown
}

func (c *Creature) SetLastRetreatAt(t time.Time) {
	c.lastRetreatAt = t
}

func (c *Creature) GetLastRetreatAt() time.Time {
	return c.lastRetreatAt
}

func (c *Creature) SetContext(ctx *dynamic_context.AIServiceContext) {
	c.context = ctx
}

func (c *Creature) GetContext() *dynamic_context.AIServiceContext {
	return c.context
}

func (c *Creature) IsInTacticalMovement() bool {
	return c.MoveCtrl.MovementPlan != nil &&
		c.MoveCtrl.MovementPlan.Type == constslib.MovementPlanCircle
}

func (c *Creature) GetCombatState() constslib.CombatState {
	return c.CombatState
}

func (c *Creature) SetCombatState(state constslib.CombatState) {
	c.CombatState = state
}

func (c *Creature) GetRecentActions() []constslib.CombatAction {
	return c.recentActions
}

func (c *Creature) ClearRecentActions() {
	c.recentActions = nil
}

func (c *Creature) AddRecentAction(action constslib.CombatAction) {
	c.recentActions = append(c.recentActions, action)
}

func (c *Creature) GetMovementFSM() fsm.MovementFSM {
	return c.movementFSM
}

func (c *Creature) SetMovementFSM(fsm fsm.MovementFSM) {
	c.movementFSM = fsm
}

func (c *Creature) GetDodgeInvulnerabilityDuration() time.Duration {
	return c.DodgeInvulnerabilityDuration
}

func (c *Creature) SetDodgeInvulnerabilityDuration(d time.Duration) {
	c.DodgeInvulnerabilityDuration = d
}

func (c *Creature) GetDodgeStaminaCost() float64 {
	return c.DodgeStaminaCost
}

func (c *Creature) SetDodgeStaminaCost(cost float64) {
	c.DodgeStaminaCost = cost
}

func (c *Creature) GetStamina() float64 {
	return c.Stamina
}

func (c *Creature) SetStamina(value float64) {
	c.Stamina = value
}

func (c *Creature) GetDodgeStartedAt() time.Time {
	return c.DodgeStartedAt
}

func (c *Creature) SetDodgeStartedAt(t time.Time) {
	c.DodgeStartedAt = t
}

func (c *Creature) GetInvulnerableUntil() time.Time {
	return c.invulnerableUntil
}

func (c *Creature) SetInvulnerableUntil(t time.Time) {
	c.invulnerableUntil = t
}

func (c *Creature) GetMoveCtrl() *movement.MovementController {
	return c.MoveCtrl
}

func (c *Creature) SetMoveCtrl(ctrl *movement.MovementController) {
	c.MoveCtrl = ctrl
}

func (c *Creature) ApplyImpulseFrom(from position.Position, duration time.Duration) {
	dir := position.CalculateDirection2D(from, c.Position)
	if dir.Length() == 0 {
		return
	}
	dir = dir.Normalize()
	dist := position.CalculateDistance2D(from, c.Position)
	dest := from.AddOffset(dir.X*dist, dir.Z*dist)

	c.MoveCtrl.SetImpulseMovement(from, dest, duration)
}
