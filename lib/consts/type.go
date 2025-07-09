package consts

// AIState representa o macro comportamento da criatura (decisão da AI)
type AIState string

const (
	AIStateIdle             AIState = "Idle"
	AIStatePatrolling       AIState = "Patrolling"
	AIStateChasing          AIState = "Chasing"
	AIStateFleeing          AIState = "Fleeing"
	AIStateDefending        AIState = "Defending"
	AIStateReturning        AIState = "ReturningHome"
	AIStateStaggered        AIState = "Staggered"
	AIStateAmbushing        AIState = "Ambushing"
	AIStateSubStealth       AIState = "SubStealth"
	AIStateStealth          AIState = "Stealth"
	AIStateAlert            AIState = "Alert"
	AIStateAttack           AIState = "Attack"
	AIStateCombat           AIState = "Combat"
	AIStateDead             AIState = "Dead"
	AIStatePostureBroken    AIState = "PostureBroken"
	AIStateSearchFood       AIState = "SearchFood"
	AIStateSearchWater      AIState = "SearchWater"
	AIStateFeeding          AIState = "Feeding"
	AIStateDrowsy           AIState = "Drowsy"
	AIStateSleeping         AIState = "Sleeping"
	AIStateSeekingSafePlace AIState = "SeekingSafePlace"
)

// CombatState representa o estado atual do combate
type CombatState uint32

const (
	CombatStateIdle           CombatState = 0
	CombatStateAttacking      CombatState = 1 << 0
	CombatStateParrying       CombatState = 1 << 1
	CombatStateBlocking       CombatState = 1 << 2
	CombatStatePostureBroken  CombatState = 1 << 3
	CombatStateStaggered      CombatState = 1 << 4
	CombatStateExecutingSkill CombatState = 1 << 5
	CombatStateRecovering     CombatState = 1 << 6
	CombatStateCombo          CombatState = 1 << 7
	CombatStateTeamSkill      CombatState = 1 << 8
	CombatStateDead           CombatState = 1 << 9
	CombatStateDodging        CombatState = 1 << 10
	CombatStateFleeing        CombatState = 1 << 11
	CombatStateAggressive     CombatState = 1 << 12
	CombatStateDefensive      CombatState = 1 << 13
	CombatStateStrategic      CombatState = 1 << 14
	CombatStateRaging         CombatState = 1 << 15
	CombatStateCautious       CombatState = 1 << 16
)

// AnimationState representa o estado visual da criatura (animação, som, FX)
type AnimationState string

const (
	AnimationIdle        AnimationState = "Idle"
	AnimationWalk        AnimationState = "Walk"
	AnimationRun         AnimationState = "Run"
	AnimationCrouchWalk  AnimationState = "CrouchWalk"
	AnimationCombatReady AnimationState = "CombatReady"
	AnimationSniff       AnimationState = "Sniff"
	AnimationParry       AnimationState = "Parry"
	AnimationBlock       AnimationState = "Block"
	AnimationJump        AnimationState = "Jump"
	AnimationAttack      AnimationState = "Attack"
	AnimationSleep       AnimationState = "Sleep"
	AnimationDie         AnimationState = "Die"
	AnimationVocalize    AnimationState = "Vocalize"
	AnimationPlay        AnimationState = "Play"
	AnimationThreat      AnimationState = "Threat"
	AnimationCurious     AnimationState = "Curious"
	AnimationLookAround  AnimationState = "LookAround"
	AnimationScratch     AnimationState = "Scratch"
	AnimationWake        AnimationState = "Wake"
	AnimationRecovery    AnimationState = "Recovery"
	AnimationWindup      AnimationState = "Windup"
	AnimationCast        AnimationState = "Cast"
)

// SkillAction representa a skill ou combo em execução
type SkillAction string

const (
	Basic      SkillAction = "Basic"
	Skill1     SkillAction = "Skill1"
	Skill2     SkillAction = "Skill2"
	Skill3     SkillAction = "Skill3"
	Skill4     SkillAction = "Skill4"
	Skill5     SkillAction = "Skill5"
	Combo1     SkillAction = "Combo1"
	Combo2     SkillAction = "Combo2"
	Combo3     SkillAction = "Combo3"
	TeamSkill1 SkillAction = "TeamSkill1"
	TeamSkill2 SkillAction = "TeamSkill2"
	TeamSkill3 SkillAction = "TeamSkill3"
)

// DamageType representa o tipo de dano
type DamageType string

const (
	DamageTypePhysical DamageType = "Physical"
	DamageTypeMagic    DamageType = "Magic"
	DamageTypeFire     DamageType = "Fire"
	DamageTypeIce      DamageType = "Ice"
	DamageTypePoison   DamageType = "Poison"
)

// StanceState representa a postura tática da criatura
type StanceState string

const (
	StanceAggressive StanceState = "Aggressive"
	StanceDefensive  StanceState = "Defensive"
	StanceNeutral    StanceState = "Neutral"
)

// EmoteState representa expressões faciais ou gestos secundários
type EmoteState string

const (
	EmoteGrowl EmoteState = "Growl"
	EmoteSnarl EmoteState = "Snarl"
	EmoteNone  EmoteState = "None"
)

type NeedType string

const (
	NeedHunger  NeedType = "Hunger"
	NeedThirst  NeedType = "Thirst"
	NeedSleep   NeedType = "Sleep"
	NeedSocial  NeedType = "Social"
	NeedFuck    NeedType = "Fuck"
	NeedKill    NeedType = "Kill"
	NeedDrink   NeedType = "Drink"
	NeedAdvance NeedType = "Advance"
	NeedGuard   NeedType = "Guard"
	NeedRetreat NeedType = "Retreat"
	NeedProvoke NeedType = "Provoke"
	NeedRecover NeedType = "Recover"
	NeedPlan    NeedType = "Plan"
	NeedFake    NeedType = "Fake"
	NeedRage    NeedType = "Rage"
)

type Need struct {
	Type         NeedType
	Value        float64 // Exemplo: 0 a 100
	LowThreshold float64
	Threshold    float64 // Valor em que a necessidade vira urgente
}

type SkillPushType string

const (
	PushNone     SkillPushType = "None"
	PushOnImpact SkillPushType = "OnImpact"
	MoveToImpact SkillPushType = "MoveToImpact"
	PushOnEnd    SkillPushType = "OnEnd"
)

type SkillType string

const (
	SkillTypePhysical SkillType = "Physical" // Habilidades corpo-a-corpo ou armas físicas
	SkillTypeMagic    SkillType = "Magic"    // Habilidades de dano ou efeito mágico
	SkillTypeUtility  SkillType = "Utility"  // Suporte, controle, buffs, debuffs
	SkillTypeEffect   SkillType = "Effect"   // Efeitos de cenário, armadilhas, armadilhas, invocações visuais
)

type MovementStyle string

const (
	MoveToFront MovementStyle = "MoveToFront" // Para na frente do alvo
	MoveThrough MovementStyle = "MoveThrough" // Atravessa o alvo (ex: Leap)
	MoveToBack  MovementStyle = "MoveToBack"  // Para atrás do alvo
)

func (cs CombatState) String() string {
	switch cs {
	case CombatStateIdle:
		return "Idle"
	case CombatStateAggressive:
		return "Aggressive"
	case CombatStateDefensive:
		return "Defensive"
	case CombatStateStrategic:
		return "Strategic"
	case CombatStateFleeing:
		return "Fleeing"
	case CombatStateRaging:
		return "Raging"
	default:
		return "Unknown"
	}
}
