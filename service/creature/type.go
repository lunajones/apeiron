package creature

type CreatureLevel string

const (
    Normal           CreatureLevel = "Normal"
    Elite            CreatureLevel = "Elite"
    Boss             CreatureLevel = "Boss"
    RegionBoss       CreatureLevel = "RegionBoss"
    WorldBoss        CreatureLevel = "WorldBoss"
)

type CreatureType string

const (
    Zombie           CreatureType = "Zombie"
    Spider           CreatureType = "Spider"
    Wolf             CreatureType = "Wolf"
    Acolyte          CreatureType = "Acolyte"
    Rabbit           CreatureType = "Rabbit"
    Human            CreatureType = "Human"
    Soldier          CreatureType = "Soldier"
	Victim           CreatureType = "Victim"
	Concubine        CreatureType = "Concubine"
	Slave            CreatureType = "Slave"
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
	AIStateStealth     AIState = "Stealth"
)

