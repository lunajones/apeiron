package combat

import (	
	"log"
)

func InitSkills() {
	
	log.Println("[Skill Registry] initializing skills...")
	SkillRegistry["SoldierSkill1"] = Skill{
		Name:        "SoldierSkill1",
		Action:      "Skill1",
		SkillType:   "Physical",
		Multiplier:  1.2,
		Range:       2.5,
		CooldownSec: 3,
	}

	SkillRegistry["SoldierGroundSlam"] = Skill{
		Name:            "SoldierGroundSlam",
		Action:          "Skill2",
		SkillType:       "Physical",
		Multiplier:      2.0,
		Range:           4.0,
		AOERadius:       3.0,
		IsGroundTargeted: true,
		CooldownSec:     8,
	}

	SkillRegistry["SoldierThrowSpear"] = Skill{
		Name:            "SoldierThrowSpear",
		Action:          "Skill3",
		SkillType:       "Physical",
		Multiplier:      1.5,
		Range:           10.0,
		HasProjectile:   true,
		ProjectileSpeed: 12.0,
		ProjectileArc:   true,
		CooldownSec:     5,
	}

	SkillRegistry["SoldierCombo1"] = Skill{
		Name:        "SoldierCombo1",
		Action:      "Combo1",
		SkillType:   "Physical",
		Multiplier:  3.5,
		Range:       3.0,
		CooldownSec: 10,
	}
}

var SkillRegistry = map[string]Skill{}
