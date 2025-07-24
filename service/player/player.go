package player

import (
	"log"
	"math/rand/v2"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	creatureconsts "github.com/lunajones/apeiron/service/creature/consts"
	"github.com/lunajones/apeiron/service/npc"
	"github.com/lunajones/apeiron/service/player/consts"
)

// Wrapper Player (composição)
type Player struct {
	Handle handle.EntityHandle
	model.Player
	CityAffiliation       string
	Relatives             []*npc.NPC
	HitboxRadius          float64
	DesiredBufferDistance float64
	Position              position.Position
	LastPosition          position.Position
	EquippedWeapon        string         // Exemplo: "Sword", "Bow", "Staff"
	LearnedSkills         map[string]int // SkillID → Skill Level
	EquippedSkills        []string       // IDs das skills equipadas (máximo 6)
	SkillPoints           int
	SkillTreeProgress     map[string]int
	Alive                 bool
	CurrentRole           consts.PlayerRole
	HP                    int
	ActiveEffects         []constslib.ActiveEffect

	CombatState     constslib.CombatState
	FacingDirection position.Vector2D
	Tags            []creatureconsts.CreatureTag

	Needs []constslib.Need

	PvPEnabled bool

	blocking           bool
	Stamina            float64
	MaxStamina         float64
	StaminaRegenPerSec float64
	DodgeDistance      float64
	DodgeStaminaCost   float64
	DodgeDisabledUntil time.Time
	invulnerableUntil  time.Time

	MaxPosture              float64
	Posture                 float64
	PostureRegenRate        float64
	PostureBroken           bool
	TimePostureBroken       int64
	PostureBreakDurationSec int

	combatDrive model.CombatDrive
	casting     bool

	combatEvents []model.CombatEvent
	MoveCtrl     *movement.MovementController
}

func (p *Player) InitNeeds() {
	p.Needs = []constslib.Need{
		{
			Type:      constslib.NeedAdvance,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
		{
			Type:      constslib.NeedGuard,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
		{
			Type:      constslib.NeedRetreat,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
		{
			Type:      constslib.NeedProvoke,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
		{
			Type:      constslib.NeedRecover,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
		{
			Type:      constslib.NeedPlan,
			Value:     rand.Float64() * 30,
			Threshold: 50,
		},
	}
}

// func SpawnPlayer(session *Session, world *World) *Player {
// 	p := &Player{
// 		Handle: handle.NewEntityHandle(lib.NewUUID(), 1),
// 		// ... outros campos obrigatórios para o Player
// 	}

// 	p.Needs = []constslib.Need{
// 		{Type: constslib.NeedAdvance, Value: 0, Threshold: 50},
// 		{Type: constslib.NeedGuard, Value: 0, Threshold: 50},
// 		{Type: constslib.NeedRetreat, Value: 0, Threshold: 50},
// 		{Type: constslib.NeedProvoke, Value: 0, Threshold: 50},
// 		{Type: constslib.NeedRecover, Value: 0, Threshold: 50},
// 		{Type: constslib.NeedCircle, Value: 0, Threshold: 50},
// 	}

// 	p.InitNeeds()

// 	world.AddPlayer(p)

// 	return p
// }

func (p *Player) HasTag(tag creatureconsts.CreatureTag) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (p *Player) GetHandle() handle.EntityHandle {
	return p.Handle
}

// Interface Targetable
func (p *Player) GetPosition() position.Position {
	return p.Position
}

func (p *Player) IsAlive() bool {
	return p.Alive
}

func (p *Player) GetLastPosition() position.Position {
	return p.LastPosition
}

func (p *Player) SetPosition(newPos position.Position) {
	p.LastPosition = p.Position
	p.Position = newPos
}

func (p *Player) GetHitboxRadius() float64 {
	return p.HitboxRadius
}

func (p *Player) GetDesiredBufferDistance() float64 {
	return p.DesiredBufferDistance
}

func (p *Player) TakeDamage(amount int) {
	if !p.Alive {
		return
	}
	p.HP -= amount
	if p.HP <= 0 {
		p.Alive = false
		// Qualquer lógica adicional de morte do player
	}
}

func (p *Player) ApplyEffect(effect constslib.ActiveEffect) {
	p.ActiveEffects = append(p.ActiveEffects, effect)
}

func (p *Player) GetFacingDirection() position.Vector2D {
	return p.FacingDirection
}

func (c *Player) IsCreature() bool { return true }

func (c *Player) IsStaticObstacle() bool {
	return false // criaturas nunca são obstáculo absoluto
}

func (p *Player) IsHungry() bool {
	for _, n := range p.Needs {
		if n.Type == constslib.NeedHunger && n.Value >= n.Threshold {
			return true
		}
	}
	return false
}

func (p *Player) IsPvPEnabled() bool {
	return p.PvPEnabled
}

func (p *Player) IsBlocking() bool {
	return p.blocking
}

func (p *Player) SetBlocking(blocking bool) {
	p.blocking = blocking
}

func (p *Player) IsInvulnerableNow() bool {
	return time.Now().Before(p.invulnerableUntil)
}

func (p *Player) ApplyPostureDamage(amount float64) {
	if p.PostureBroken || !p.Alive {
		return
	}

	p.Posture -= amount
	if p.Posture <= 0 {
		p.Posture = 0
		p.PostureBroken = true
		p.TimePostureBroken = time.Now().Unix()
		p.CombatState = constslib.CombatStateStaggered
		// Exemplo de possível animação: c.SetAnimationState(constslib.AnimationIdle) ou custom para stagger
		// physics.StartStagger(&c.Stagger, c.TimePostureBroken, float64(c.PostureBreakDurationSec))
		log.Printf("[Player %s] Posture quebrada! Entrando em stagger.", p.Handle.ID)
	}
}

func (p *Player) IsInParryWindow() bool {
	return true
}

func (p *Player) GetFaction() string {
	return p.Faction
}

func (p *Player) IsCasting() bool {
	return p.casting
}

func (p *Player) SetCasting(casting bool) {
	p.casting = casting
}

func (p *Player) GetCombatDrive() *model.CombatDrive {
	return &p.combatDrive
}

func (p *Player) GetCombatEvents() []model.CombatEvent {
	return p.combatEvents
}

func (p *Player) RegisterCombatEvent(evt model.CombatEvent) {
	p.combatEvents = append(p.combatEvents, evt)
}

func (p *Player) ApplyImpulseFrom(from position.Position, duration time.Duration) {
	dir := position.CalculateDirection2D(from, p.Position)
	if dir.Length() == 0 {
		return
	}
	dir = dir.Normalize()
	dist := position.CalculateDistance2D(from, p.Position)
	dest := from.AddOffset(dir.X*dist, dir.Z*dist)

	p.MoveCtrl.SetImpulseMovement(from, dest, duration)
}
