package creature

type CreatureLevel string

const (
    Normal           CreatureLevel = "Normal"
    Elite            CreatureLevel = "Elite"
    Boss             CreatureLevel = "Boss"
    RegionBoss       CreatureLevel = "RegionBoss"
    WorldBoss        CreatureLevel = "WorldBoss"
)

const AnyType        CreatureType = "Any"

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
    ActionSkill4     CreatureAction = "Skill4"
    ActionSkill5     CreatureAction = "Skill5"
    ActionCombo1     CreatureAction = "Combo1"
    ActionCombo2     CreatureAction = "Combo2"
    ActionCombo3     CreatureAction = "Combo3"
    ActionTeamSkill1 CreatureAction = "TeamSkill1"
    ActionTeamSkill2 CreatureAction = "TeamSkill2"
    ActionTeamSkill3 CreatureAction = "TeamSkill3"
    ActionDie        CreatureAction = "Die"
    ActionAttack     CreatureAction = "Attack"
    ActionSleep      CreatureAction = "Sleep"
)

type AIState string

const (
	AIStateIdle          AIState = "Idle"
	AIStatePatrolling    AIState = "Patrolling"
	AIStateChasing       AIState = "Chasing"
	AIStateFleeing       AIState = "Fleeing"
	AIStateDefending     AIState = "Defending"
	AIStateReturning     AIState = "ReturningHome"
	AIStateStaggered     AIState = "Staggered"
	AIStateAmbushing     AIState = "Ambushing"
	AIStateSubStealth    AIState = "SubStealth"
	AIStateStealth       AIState = "Stealth"
    AIStateAlert         AIState = "Alert"
	AIStateAttack        AIState = "Attack"
	AIStateCombat        AIState = "Combat"
	AIStateDead          AIState = "Dead"
	AIStatePostureBroken AIState = "PostureBroken"
    AIStateSearchFood    AIState = "SearchFood"
	AIStateSearchWater   AIState = "SearchWater"
	AIStateFeeding       AIState = "Feeding" 
)


type NeedType int

const (
	NeedHunger NeedType = iota
	NeedThirst
	NeedSleep
	NeedSocial
)

type Need struct {
	Type      NeedType
	Value     float64  // Exemplo: 0 a 100
	Threshold float64  // Valor em que a necessidade vira urgente
}

type Role string

const (
	RoleMerchant Role = "Merchant"
	RoleHunter   Role = "Hunter"
	RoleGuard    Role = "Guard"
	RoleNone     Role = "None"
)

type CreatureTag string

const (
	TagHumanoid CreatureTag = "Humanoid"
	TagAnimal   CreatureTag = "Animal"
	TagPrey     CreatureTag = "Prey"
	TagPredator CreatureTag = "Predator"
	TagMerchant CreatureTag = "Merchant"
)

type MentalState string

const (
	MentalStateCalm     MentalState = "Calm"
	MentalStateAfraid   MentalState = "Afraid"
	MentalStateAggressive MentalState = "Aggressive"
	MentalStateEnraged  MentalState = "Enraged"
	MentalStateDesperate MentalState = "Desperate"
)