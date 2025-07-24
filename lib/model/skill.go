package model

import (
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
)

// -- Core Skill Struct --

type Skill struct {
	ID                string
	Name              string
	Tags              *SkillTags
	Action            constslib.SkillAction
	SkillType         constslib.SkillType
	Range             float64
	CooldownSec       float64
	InitialMultiplier float64

	WindUpTime           float64
	CastTime             float64
	RecoveryTime         float64
	Interruptible        bool
	CanCastWhileBlocking bool

	TargetLock     bool
	GroundTargeted bool
	RotationLock   bool

	Impact     *ImpactEffect
	Movement   *MovementConfig
	AOE        *AOEConfig
	Projectile *ProjectileConfig
	Teleport   *TeleportConfig
	DOT        *DOTConfig
	HasDOT     bool
	Buff       *BuffConfig
	Debuff     *DebuffConfig
	Conditions *SkillCondition
	Hitbox     *HitboxConfig

	StaminaDamage float64
	ScoreBase     float64
}

// -- Runtime Skill State --

type SkillState struct {
	Skill                         *Skill
	StartedAt                     time.Time
	WindUpUntil                   time.Time
	CastUntil                     time.Time
	RecoveryUntil                 time.Time
	CooldownUntil                 time.Time
	InUse                         bool
	ChargesLeft                   int
	LastUsedAt                    time.Time
	EffectApplied                 bool
	WindUpFired                   bool
	CastFired                     bool
	RecoveryFired                 bool
	NextQueuedSkill               *Skill
	WasInterrupted                bool
	AppliedBuffID                 string
	HasAggressiveIntentRegistered bool
}

func (s *SkillState) CanBeCancelled() bool {
	if s == nil || s.Skill == nil {
		return true
	}

	if s.InUse && !s.Skill.Interruptible {
		return false
	}

	now := time.Now()

	if now.Before(s.WindUpUntil) {
		return true
	}

	if now.After(s.WindUpUntil) && now.Before(s.CastUntil) {
		return s.Skill.Interruptible
	}

	if now.After(s.CastUntil) && now.Before(s.RecoveryUntil) {
		return s.RecoveryFired
	}

	return false
}

// -- Config Structs --

type ImpactEffect struct {
	PostureDamage     float64
	ScalingStat       string
	ScalingMultiplier float64
	DefenseStat       string
}

type MovementConfig struct {
	Speed                    float64
	DurationSec              float64
	MaxDistance              float64
	ExtraDistance            float64
	DirectionLock            bool
	MicroHoming              bool
	TargetLock               bool
	Interruptible            bool
	PushType                 constslib.SkillPushType
	Style                    constslib.MovementStyle
	BlockDuringMovement      bool // Ativa bloqueio frontal enquanto a skill estiver em movimento
	StopOnFirstHit           bool // Finaliza o movimento ao colidir com o primeiro inimigo
	PushTargetDuringMovement bool
	SeparationRadius         float64 // Distância para empurrar outros ao final do movimento
	SeparationForce          float64
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

type DOTConfig struct {
	DurationSec int
	TickSec     int
	TickPower   int
	EffectType  constslib.EffectType
}

type BuffConfig struct {
	Name           string
	DurationSec    float64
	StatModifiers  map[string]float64
	Resistances    map[string]float64
	IsStackable    bool
	TargetSelf     bool
	MaxStacks      int
	VisualEffectID string
}

type DebuffConfig struct {
	Name           string
	DurationSec    float64
	StatModifiers  map[string]float64
	Resistances    map[string]float64
	DamagePerSec   float64
	IsStackable    bool
	MaxStacks      int
	VisualEffectID string
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

type HitboxShape int

const (
	HitboxBox HitboxShape = iota
	HitboxCone
	HitboxCircle
	HitboxLine
)

type HitboxConfig struct {
	Shape       HitboxShape
	Length      float64
	Width       float64
	Angle       float64
	MinRadius   float64
	MaxRadius   float64
	OffsetFront float64
}

// -- Tag Handling --

type SkillTags struct {
	values map[string]bool
}

func NewSkillTags(tags ...string) *SkillTags {
	t := &SkillTags{values: make(map[string]bool)}
	for _, tag := range tags {
		t.values[string(tag)] = true
	}
	return t
}

func (t *SkillTags) Has(tag constslib.SkillTag) bool {
	if t == nil {
		return false
	}
	return t.values[string(tag)]
}
