package model

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
)

var SkillRegistry = map[string]*Skill{}

func InitSkills() {
	log.Println("[Skill Registry] initializing skills...")

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
		ScoreBase: 5.0,
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
		CastTime:          0.3,
		RecoveryTime:      0.2,
		Interruptible:     false,
		Impact: &ImpactEffect{
			PostureDamage:     12,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase: 3.0,
		Movement: &MovementConfig{
			Speed:         4.0,
			DurationSec:   0.3,
			MaxDistance:   2.5,
			DirectionLock: true,
			TargetLock:    false,
			Interruptible: false,
			PushType:      constslib.PushOnImpact,
			Style:         constslib.MovementStyle(constslib.MoveToImpact),
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
		CastTime:          2,
		RecoveryTime:      0.3,
		Interruptible:     false,
		Impact: &ImpactEffect{
			PostureDamage:     12,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.15,
			DefenseStat:       "ControlResistance",
		},
		ScoreBase:     4.5,
		StaminaDamage: 15,
		Movement: &MovementConfig{
			Speed:         8.0,
			DurationSec:   2,
			MaxDistance:   3.0,
			DirectionLock: true,
			MicroHoming:   true,
			TargetLock:    true,
			Interruptible: false,
			ExtraDistance: 0,
			PushType:      constslib.PushOnEnd,
			Style:         constslib.MoveThrough,
		},
	}

	log.Println("[Skill Registry] finishing system...")
}
