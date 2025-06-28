package combat

import (
	"log"

	"github.com/lunajones/apeiron/service/creature/consts"
)

var SkillRegistry = map[string]Skill{}

func InitSkills() {
	log.Println("[Skill Registry] initializing skills...")

	SkillRegistry["SoldierSkill1"] = Skill{
		Name:              "SoldierSkill1",
		Action:            consts.ActionSkill1,
		SkillType:         "Physical",
		InitialMultiplier: 1.2,
		Range:             2.5,
		CooldownSec:       3,
		Impact: &ImpactEffect{
			PostureDamageBase: 10,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
	}

	SkillRegistry["SoldierGroundSlam"] = Skill{
		Name:              "SoldierGroundSlam",
		Action:            consts.ActionSkill2,
		SkillType:         "Physical",
		InitialMultiplier: 2.0,
		Range:             4.0,
		CooldownSec:       8,
		GroundTargeted:    true,
		AOE: &AOEConfig{
			Radius: 3.0,
			Shape:  "Circle",
		},
		Impact: &ImpactEffect{
			PostureDamageBase: 20,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.2,
			DefenseStat:       "ControlResistance",
		},
	}

	SkillRegistry["SoldierThrowSpear"] = Skill{
		Name:              "SoldierThrowSpear",
		Action:            consts.ActionSkill3,
		SkillType:         "Physical",
		InitialMultiplier: 1.5,
		Range:             10.0,
		CooldownSec:       5,
		Projectile: &ProjectileConfig{
			Speed:       12.0,
			HasArc:      true,
			LifeTimeSec: 2,
		},
		Impact: &ImpactEffect{
			PostureDamageBase: 12,
			ScalingStat:       "Dexterity",
			ScalingMultiplier: 0.15,
			DefenseStat:       "ControlResistance",
		},
	}

	SkillRegistry["SoldierCombo1"] = Skill{
		Name:              "SoldierCombo1",
		Action:            consts.ActionCombo1,
		SkillType:         "Physical",
		InitialMultiplier: 3.5,
		Range:             3.0,
		CooldownSec:       10,
		Impact: &ImpactEffect{
			PostureDamageBase: 25,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.3,
			DefenseStat:       "ControlResistance",
		},
	}

	SkillRegistry["Bite"] = Skill{
		Name:              "Bite",
		Action:            consts.ActionSkill1,
		SkillType:         "Physical",
		InitialMultiplier: 0.5,
		Range:             2.0,
		CooldownSec:       1,
		Impact: &ImpactEffect{
			PostureDamageBase: 4,
			ScalingStat:       "Strength",
			ScalingMultiplier: 0.05,
			DefenseStat:       "ControlResistance",
		},
	}

	SkillRegistry["Lacerate"] = Skill{
		Name:              "Lacerate",
		Action:            consts.ActionSkill2,
		SkillType:         "Physical",
		InitialMultiplier: 1.0,
		Range:             3.0,
		CooldownSec:       4,
		HasDOT:            true,
		DOT: &DOTConfig{
			DurationSec: 6,
			TickSec:     2,
			TickPower:   3,
			EffectType:  consts.EffectPoison,
		},
		Impact: &ImpactEffect{
			PostureDamageBase: 8,
			ScalingStat:       "Dexterity",
			ScalingMultiplier: 0.1,
			DefenseStat:       "ControlResistance",
		},
	}

	log.Println("[Skill Registry] finishing system...")
}
