package creature

type CreatureType string
type CreatureSubtype string
type DamageType string
type CreatureAction string

type CreatureLevel string

const (
    Mob              CreatureType = "Mob"
    Elite            CreatureType = "Elite"
    Boss             CreatureType = "Boss"
    RegionBoss       CreatureType = "RegionBoss"
    WorldBoss        CreatureType = "WorldBoss"
)

type CreatureType string

const (
    Zombie           CreatureSubtype = "Zombie"
    Spider           CreatureSubtype = "Spider"
    Wolf             CreatureSubtype = "Wolf"
    Acolyte          CreatureSubtype = "Acolyte"
    Rabbit           CreatureSubtype = "Rabbit"
    Human            CreatureSubtype = "Human"
    Soldier          CreatureSubtype = "Soldier"
	Victim           CreatureSubtype = "Victim"
	Concubine        CreatureSubtype = "Concubine"
	Slave            CreatureSubtype = "Slave"
    Soldier          CreatureSubtype = "Soldier"
)

type CreatureAction string

const (
    ActionIdle       CreatureAction = "Idle"
    ActionWalk       CreatureAction = "Walk"
    ActionParry      CreatureAction = "Parry"
    ActionBlock      CreatureAction = "Block"
    ActionRun        CreatureAction = "Run"
    ActionJump       CreatureAction = "Jump"
    ActionSkill1     CreatureAction = "Skill1"
    ActionSkill2     CreatureAction = "Skill2"
    ActionSkill3     CreatureAction = "Skill3"
    ActionCombo1     CreatureAction = "Combo1"
    ActionCombo2     CreatureAction = "Combo2"
    ActionCombo3     CreatureAction = "Combo3"
    ActionTeamSkill1 CreatureAction = "TeamSkill1"
    ActionTeamSkill2 CreatureAction = "TeamSkill2"
    ActionTeamSkill3 CreatureAction = "TeamSkill3"
    ActionDie        CreatureAction = "Die"
)

type AIState string

const (
	AIStateIdle        AIState = "Idle"
	AIStatePatrolling  AIState = "Patrolling"
	AIStateChasing     AIState = "Chasing"
	AIStateFleeing     AIState = "Fleeing"
	AIStateReturning   AIState = "ReturningHome"
	AIStateStaggered   AIState = "Staggered"
	AIStateAmbushing   AIState = "Ambushing"
	AIStateSubStealth  AIState = "SubStealth"
	AIStateSubStealth  AIState = "Stealth"
)

