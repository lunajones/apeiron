package model

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
)

var SkillRegistry = map[string]*Skill{}

func InitSkills() {
	log.Println("[Skill Registry] initializing skills...")

	// ===== SOLDIER SKILLS =====

	SkillRegistry["SoldierSlash"] = &Skill{
		ID:                "SoldierSlash",
		Name:              "SoldierSlash",
		Tags:              NewSkillTags("Burst"),
		Action:            constslib.Basic,
		SkillType:         "Physical",
		InitialMultiplier: 1.0,
		Range:             2.2,
		CooldownSec:       0.4,
		WindUpTime:        0.1,
		CastTime:          0.4,
		RecoveryTime:      0.4,
		Interruptible:     true,
		Impact: &ImpactEffect{
			PostureDamage:     5,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.05,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 1.0,
		Hitbox: &HitboxConfig{
			Shape:     HitboxCone,
			MinRadius: 0.4,
			MaxRadius: 2.2,
			Angle:     90,
		},
	}

	SkillRegistry["SoldierShieldBash"] = &Skill{
		ID:                "SoldierShieldBash",
		Name:              "SoldierShieldBash",
		Tags:              NewSkillTags("Interrupt", "Burst"),
		Action:            constslib.Skill1,
		SkillType:         "Physical",
		InitialMultiplier: 0.8,
		Range:             2.0,
		CooldownSec:       3.0,
		WindUpTime:        0.2,
		CastTime:          0.2,
		RecoveryTime:      0.2,
		Interruptible:     true,
		Impact: &ImpactEffect{
			PostureDamage:     10,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 2.0,
		Hitbox: &HitboxConfig{
			Shape:     HitboxCone,
			MinRadius: 0.2,
			MaxRadius: 2.0,
			Angle:     60,
		},
	}

	SkillRegistry["SoldierGroundSlam"] = &Skill{
		ID:                "SoldierGroundSlam",
		Name:              "SoldierGroundSlam",
		Tags:              NewSkillTags("AOE", "Burst"),
		Action:            constslib.Skill2,
		SkillType:         "Physical",
		InitialMultiplier: 1.5,
		Range:             3.0,
		CooldownSec:       6.0,
		WindUpTime:        0.4,
		CastTime:          0.4,
		RecoveryTime:      0.3,
		Interruptible:     false,
		AOE: &AOEConfig{
			Radius: 3.0,
			Shape:  "Circle",
		},
		Impact: &ImpactEffect{
			PostureDamage:     15,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.15,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 3.5,
		Hitbox: &HitboxConfig{
			Shape:     HitboxCircle,
			MaxRadius: 3.0,
		},
	}

	SkillRegistry["SoldierLongStep"] = &Skill{
		ID:                "SoldierLongStep",
		Name:              "SoldierLongStep",
		Tags:              NewSkillTags("Rush", "Burst"),
		Action:            constslib.Skill3,
		SkillType:         "Physical",
		InitialMultiplier: 1.3,
		Range:             3.0,
		CooldownSec:       4.0,
		WindUpTime:        0.2,
		CastTime:          0.3,
		RecoveryTime:      0.2,
		Interruptible:     true,
		Impact: &ImpactEffect{
			PostureDamage:     6,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 4.0,
		Hitbox: &HitboxConfig{
			Shape:  HitboxBox,
			Length: 3.0,
			Width:  1.0,
		},
	}

	SkillRegistry["SoldierShieldRush"] = &Skill{
		ID:                "SoldierShieldRush",
		Name:              "SoldierShieldRush",
		Tags:              NewSkillTags("Rush", "Burst"),
		Action:            constslib.Skill4,
		SkillType:         "Physical",
		InitialMultiplier: 1.2,
		Range:             2.5,
		CooldownSec:       5.0,
		WindUpTime:        0.3,
		CastTime:          2.4,
		RecoveryTime:      0.2,
		Interruptible:     false,
		ScoreBase:         4.0,

		Impact: &ImpactEffect{
			PostureDamage:     12,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},

		Movement: &MovementConfig{
			Speed:                    4.0,
			DurationSec:              0.3,
			MaxDistance:              4.0,
			DirectionLock:            true,
			TargetLock:               true,
			Interruptible:            false,
			PushType:                 constslib.PushOnImpact,
			Style:                    constslib.MoveToFront,
			SeparationRadius:         0.8,
			SeparationForce:          0.4,
			BlockDuringMovement:      true,  // ✅ escudo levantado
			StopOnFirstHit:           false, // ✅ continua empurrando
			PushTargetDuringMovement: true,  // ✅ engata e arrasta
		},

		Hitbox: &HitboxConfig{
			Shape:  HitboxLine,
			Length: 2.5,
			Width:  1.2,
		},
	}

	SkillRegistry["SoldierRiposteStance"] = &Skill{
		ID:                "SoldierRiposteStance",
		Name:              "SoldierRiposteStance",
		Tags:              NewSkillTags("Utility"),
		Action:            constslib.Combo1,
		SkillType:         "Physical",
		InitialMultiplier: 1.0,
		Range:             1.5,
		CooldownSec:       8.0,
		WindUpTime:        0.0,
		CastTime:          3.0,
		RecoveryTime:      0.5,
		Interruptible:     false,
		Impact: &ImpactEffect{
			PostureDamage:     15,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.2,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 4.0,
	}

	// ===== WOLF SKILLS =====

	SkillRegistry["Bite"] = &Skill{
		ID:                "Bite",
		Name:              "Bite",
		Tags:              NewSkillTags("Interrupt"),
		Action:            constslib.Basic,
		SkillType:         "Physical",
		InitialMultiplier: 0.8,
		Range:             2.2,
		CooldownSec:       1.1,
		WindUpTime:        1.0,
		CastTime:          0.4,
		RecoveryTime:      1.1,
		Interruptible:     true,
		Impact: &ImpactEffect{
			PostureDamage:     4,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.05,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase:     1.0,
		StaminaDamage: 5,
		Hitbox: &HitboxConfig{
			Shape:  HitboxBox,
			Length: 2.2,
			Width:  1.0,
		},
	}

	SkillRegistry["Lacerate"] = &Skill{
		ID:                "Lacerate",
		Name:              "Lacerate",
		Tags:              NewSkillTags("DOT", "Burst"),
		Action:            constslib.Skill1,
		SkillType:         "Physical",
		InitialMultiplier: 1.0,
		Range:             2.5,
		CooldownSec:       6.0,
		WindUpTime:        0.6,
		CastTime:          0.5,
		RecoveryTime:      1.4,
		Interruptible:     true,
		HasDOT:            true,
		DOT: &DOTConfig{
			DurationSec: 6,
			TickSec:     2,
			TickPower:   3,
			EffectType:  constslib.EffectPoison,
		},
		Impact: &ImpactEffect{
			PostureDamage:     8,
			ScalingStat:       "Dexterity",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase:     5.0,
		StaminaDamage: 10,
		Hitbox: &HitboxConfig{
			Shape:     HitboxCone,
			MinRadius: 0.3,
			MaxRadius: 2.5,
			Angle:     100,
		},
	}

	SkillRegistry["Leap"] = &Skill{
		ID:                "Leap",
		Name:              "Leap",
		Tags:              NewSkillTags("Rush", "Burst"),
		Action:            constslib.Skill2,
		SkillType:         "Physical",
		InitialMultiplier: 1.8,
		Range:             3.0,
		CooldownSec:       6.0,
		WindUpTime:        0.3,
		CastTime:          1.2,
		RecoveryTime:      0.4,
		Interruptible:     false,
		ScoreBase:         4.5,
		StaminaDamage:     15,

		Impact: &ImpactEffect{
			PostureDamage:     12,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.15,
			DefenseStat:       "ControlResistance",
		},

		Movement: &MovementConfig{
			Speed:                    7.5,
			DurationSec:              0.66,
			MaxDistance:              5.0,
			DirectionLock:            true,
			MicroHoming:              true,
			TargetLock:               true,
			Interruptible:            false,
			ExtraDistance:            0,
			PushType:                 constslib.PushOnImpact,
			Style:                    constslib.MoveToFront,
			SeparationRadius:         0.9,
			SeparationForce:          0.5,
			BlockDuringMovement:      false, // ❌ Leap não bloqueia
			StopOnFirstHit:           false, // ❌ Leap atravessa
			PushTargetDuringMovement: false, // ❌ Leap não engata
		},

		Hitbox: &HitboxConfig{
			Shape:  HitboxLine,
			Length: 5.0,
			Width:  1.4,
		},
	}

	log.Println("[Skill Registry] finishing system...")
}
