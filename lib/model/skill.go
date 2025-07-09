package model

import (
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
)

type Skill struct {
	Name              string
	Action            constslib.SkillAction
	SkillType         constslib.SkillType // "Physical", "Magic", "Utility", "Effect"
	Range             float64
	CooldownSec       float64
	InitialMultiplier float64

	HasDOT     bool
	DOT        *DOTConfig
	AOE        *AOEConfig
	Projectile *ProjectileConfig
	Teleport   *TeleportConfig
	Impact     *ImpactEffect
	Conditions *SkillCondition

	TargetLock     bool
	GroundTargeted bool
	RotationLock   bool
	CastTime       float64
	WindUpTime     float64
	RecoveryTime   float64
	Interruptible  bool // NOVO: permite definir se a skill pode ser interrompida durante execução

	ScoreBase float64

	Movement *MovementConfig `json:"movement,omitempty"` // Indica se essa skill causa movimento

	Buff   *BuffConfig   // ✅ Configuração de Buff (nil se não tiver)
	Debuff *DebuffConfig // ✅ Configuração de Debuff (nil se não tiver)

	StaminaDamage        float64 // Dano causado à stamina do bloqueador (ou 0 se não aplicar)
	CanCastWhileBlocking bool
}

type MovementConfig struct {
	Speed         float64                 `json:"speed"`         // Velocidade do avanço
	DurationSec   float64                 `json:"durationSec"`   // Duração total do avanço
	MaxDistance   float64                 `json:"maxDistance"`   // Distância máxima
	ExtraDistance float64                 `json:"extraDistance"` // ← NOVO! opcional, default 0
	DirectionLock bool                    `json:"directionLock"` // Travar direção no início (true = sem homing)
	MicroHoming   bool                    `json:"microHoming"`   // Permite leve ajuste inicial (ex: 10% do avanço)
	TargetLock    bool                    `json:"targetLock"`    // Mira na posição do alvo no início
	Interruptible bool                    `json:"interruptible"` // Se o avanço pode ser interrompido (ex: por parry, block)
	PushType      constslib.SkillPushType `json:"pushType"`      // ← NOVO! Tipo de push associado ao movimento
	Style         constslib.MovementStyle `json:"style"`
}

type DOTConfig struct {
	DurationSec int
	TickSec     int
	TickPower   int
	EffectType  constslib.EffectType
}

type AOEConfig struct {
	Radius float64
	Shape  string
	Angle  float64
}

type ProjectileConfig struct {
	Speed       float64
	HasArc      bool
	LifeTimeSec int
}

type TeleportConfig struct {
	ToBackOfTarget bool
	DistanceOffset float64
}

type BuffConfig struct {
	Name          string             // Nome do buff (ex: "Fortitude", "Haste")
	DurationSec   float64            // Duração em segundos
	StatModifiers map[string]float64 // Modificadores de atributos (ex: "Strength": +10)
	Resistances   map[string]float64 // Modificadores de resistência (ex: "Fire": +15)
	IsStackable   bool               // Permite acumular múltiplas instâncias
	TargetSelf    bool               // ✅ true se o buff só pode ser aplicado no caster

	MaxStacks      int    // Quantidade máxima de stacks permitidos
	VisualEffectID string // ID do efeito visual associado
}

type DebuffConfig struct {
	Name           string             // Nome do debuff (ex: "Poison", "Slow")
	DurationSec    float64            // Duração em segundos
	StatModifiers  map[string]float64 // Reduções de atributos (ex: "Speed": -20)
	Resistances    map[string]float64 // Redução de resistência (ex: "Armor": -10)
	DamagePerSec   float64            // Dano por segundo (caso aplique)
	IsStackable    bool               // Permite acumular múltiplas instâncias
	MaxStacks      int                // Quantidade máxima de stacks permitidos
	VisualEffectID string             // ID do efeito visual associado
}

type ImpactEffect struct {
	PostureDamage     float64 // Agora é o campo único e definitivo para postura
	ScalingStat       string
	ScalingMultiplier float64
	DefenseStat       string
}

type SkillCondition struct {
	RequiredStates     []string
	FacingRequirement  string
	TargetMustBeAlive  bool
	OnlyIfParrySuccess bool
}

type SkillResult struct {
	Success       bool
	WasOnCooldown bool
	DamageDealt   int
	TargetDied    bool
	PostureBroken bool
	EffectApplied bool
	CriticalHit   bool
	Blocked       bool
	Parried       bool
	Interrupted   bool
	WasAOEHit     bool
}

type SkillState struct {
	Skill       *Skill
	StartedAt   time.Time
	WindUpUntil time.Time
	WindUpFired bool // ✅ novo campo

	CastUntil     time.Time
	RecoveryUntil time.Time
	CooldownUntil time.Time
	InUse         bool

	ChargesLeft int
	LastUsedAt  time.Time

	EffectApplied    bool
	HasCastBeenFired bool   // <- novo campo
	NextQueuedSkill  *Skill // NOVO: permite enfileirar próxima skill
	WasInterrupted   bool   // NOVO: indica se esta execução foi interrompida
	AppliedBuffID    string

	HasAggressiveIntentRegistered bool
}

func (s *SkillState) CanBeCancelled() bool {
	now := time.Now()

	if now.Before(s.WindUpUntil) {
		return true // Ainda está no Windup
	}
	if now.After(s.CastUntil) && now.Before(s.RecoveryUntil) {
		return true // Está no Recovery
	}
	return false // Está no Cast ou já terminou
}
