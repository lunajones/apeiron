package combat

import (
	"github.com/lunajones/apeiron/service/creature/consts"
)

type Skill struct {
	Name              string
	Action            consts.CreatureAction
	SkillType         string // "Physical", "Magic", "Utility", "Effect"
	Range             float64
	CooldownSec       int
	InitialMultiplier float64 // Dano base direto

	HasDOT     bool
	DOT        *DOTConfig
	AOE        *AOEConfig
	Projectile *ProjectileConfig
	Teleport   *TeleportConfig
	Buff       *BuffConfig
	Impact     *ImpactEffect
	Conditions *SkillCondition

	TargetLock     bool
	GroundTargeted bool
	RotationLock   bool
	CastTime       float64
	WindUpTime     float64
	RecoveryTime   float64
}

type DOTConfig struct {
	DurationSec int
	TickSec     int
	TickPower   int
	EffectType  consts.EffectType // Poison, Burn, Bleed
}

type AOEConfig struct {
	Radius float64
	Shape  string  // "Circle", "Cone", "Rectangle"
	Angle  float64 // Para "Cone"
}

type ProjectileConfig struct {
	Speed       float64
	HasArc      bool
	LifeTimeSec int
}

type TeleportConfig struct {
	ToBackOfTarget bool    // Se verdadeiro, move para trás do alvo
	DistanceOffset float64 // Distância exata atrás do alvo
}

type BuffConfig struct {
	StatAffected string
	Modifier     float64
	DurationSec  int
	TargetSelf   bool
}

type ImpactEffect struct {
	PostureDamageBase float64
	ScalingStat       string  // Ex: "Strength"
	ScalingMultiplier float64 // Ex: 0.2 → 20% do stat
	DefenseStat       string  // Ex: "ControlResistance"
}

type SkillCondition struct {
	RequiredStates     []string // Ex: "Invisible"
	FacingRequirement  string   // "Behind", "Front", "Any"
	TargetMustBeAlive  bool
	OnlyIfParrySuccess bool // Útil para skills como "Shadow Reversal"
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
