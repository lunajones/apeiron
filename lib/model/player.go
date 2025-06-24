package model

import "github.com/lunajones/apeiron/lib/position"

type Player struct {
	ID               string
	Name             string
	Position         position.Position
	HP               int
	MaxHP            int
	EquippedWeapon   string            // Exemplo: "Sword", "Bow", "Staff"
	LearnedSkills    map[string]int    // SkillID → Skill Level
	EquippedSkills   []string          // IDs das skills equipadas (máximo 6)
	SkillPoints      int
	SkillTreeProgress map[string]int   // SkillID → Progress (ex: pontos gastos)
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) GetPosition() position.Position {
	return p.Position
}