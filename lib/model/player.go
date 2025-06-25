package model

type Player struct {
	ID               string
	Name             string
	MaxHP            int  // SkillID → Progress (ex: pontos gastos)
}

func (p *Player) GetID() string {
	return p.ID
}