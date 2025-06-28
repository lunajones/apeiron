package player

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
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
	IsPvPEnabled          bool // usado s
	HitboxRadius          float64
	DesiredBufferDistance float64
	Position              position.Position
	LastPosition          position.Position
	EquippedWeapon        string         // Exemplo: "Sword", "Bow", "Staff"
	LearnedSkills         map[string]int // SkillID → Skill Level
	EquippedSkills        []string       // IDs das skills equipadas (máximo 6)
	SkillPoints           int
	SkillTreeProgress     map[string]int
	IsAlive               bool
	CurrentRole           consts.PlayerRole
	HP                    int
	ActiveEffects         []creatureconsts.ActiveEffect

	FacingDirection position.Vector2D
}

func (p *Player) GetHandle() handle.EntityHandle {
	return p.Handle
}

// Interface Targetable
func (p *Player) GetPosition() position.Position {
	return p.Position
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) CheckIsAlive() bool {
	return p.IsAlive
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
	if !p.IsAlive {
		return
	}
	p.HP -= amount
	if p.HP <= 0 {
		p.IsAlive = false
		// Qualquer lógica adicional de morte do player
	}
}

func (p *Player) ApplyEffect(effect creatureconsts.ActiveEffect) {
	p.ActiveEffects = append(p.ActiveEffects, effect)
}

func (p *Player) GetFacingDirection() position.Vector2D {
	return p.FacingDirection
}

func (c *Player) IsCreature() bool { return true }
