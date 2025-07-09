package model

type Player struct {
	ID           string
	Name         string
	MaxHP        int // SkillID → Progress (ex: pontos gastos)
	Strength     int
	Dexterity    int
	Intelligence int
	Focus        int
	Hostile      bool
	Alive        bool
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) GetStrength() int {
	return p.Strength
}

func (p *Player) GetDexterity() int {
	return p.Dexterity
}

func (p *Player) GetIntelligence() int {
	return p.Intelligence
}

func (p *Player) GetFocus() int {
	return p.Focus
}

func (p *Player) IsAlive() bool {
	return p.Alive
}

func (p *Player) IsHostile() bool {
	return p.Hostile
}
