package combat

import (
	"time"

	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/lib/position"
)

type Skill struct {
	Name            string
	Action          creature.CreatureAction
	SkillType       string  // "Physical", "Magic", etc
	Multiplier      float64
	Range           float64
	CooldownSec     int
	IsDOT           bool
	DOTDurationSec  int
	DOTTickSec      int
	AOERadius       float64
	IsGroundTargeted bool
	HasProjectile   bool
	ProjectileSpeed float64
	ProjectileArc   bool
}

type SkillExecution struct {
	SkillName  string
	CasterID   string
	TargetID   string
	TargetPos  position.Position
	ExecuteAt  time.Time
}
