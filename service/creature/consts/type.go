package consts

type CreatureLevel string

const (
	Normal     CreatureLevel = "Normal"
	Elite      CreatureLevel = "Elite"
	Boss       CreatureLevel = "Boss"
	RegionBoss CreatureLevel = "RegionBoss"
	WorldBoss  CreatureLevel = "WorldBoss"
)

const AnyType CreatureType = "Any"

type CreatureType string

const (
	Zombie    CreatureType = "Zombie"
	Spider    CreatureType = "Spider"
	Wolf      CreatureType = "Wolf"
	Acolyte   CreatureType = "Acolyte"
	Rabbit    CreatureType = "Rabbit"
	Human     CreatureType = "Human"
	Soldier   CreatureType = "Soldier"
	Victim    CreatureType = "Victim"
	Concubine CreatureType = "Concubine"
	Slave     CreatureType = "Slave"
)

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
	TagUndead   CreatureTag = "Undead"
	TagAnimal   CreatureTag = "Animal"
	TagPrey     CreatureTag = "Prey"
	TagCoward   CreatureTag = "Coward"
	TagPredator CreatureTag = "Predator"
	TagMerchant CreatureTag = "Merchant"
)

type MentalState string

const (
	MentalStateCalm       MentalState = "Calm"
	MentalStateAfraid     MentalState = "Afraid"
	MentalStateAggressive MentalState = "Aggressive"
	MentalStateEnraged    MentalState = "Enraged"
	MentalStateDesperate  MentalState = "Desperate"
)
