package player

import (
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
)

// Wrapper Player (composição)
type Player struct {
	model.Player
	HP               int
	Position         position.Position
	EquippedWeapon   string            // Exemplo: "Sword", "Bow", "Staff"
	LearnedSkills    map[string]int    // SkillID → Skill Level
	EquippedSkills   []string          // IDs das skills equipadas (máximo 6)
	SkillPoints      int
	SkillTreeProgress map[string]int 
	IsAlive          bool
	CurrentRole      PlayerRole
}

// Interface Targetable
func (p *Player) GetPosition() position.Position {
	return p.Position
}

func (p *Player) GetID() string {
	return p.ID
}