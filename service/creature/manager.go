package creature

import (
	"log"
	"math/rand"
	"time"
)

type Manager struct {
	creatures map[string]*Creature
}

func NewManager() *Manager {
	return &Manager{
		creatures: make(map[string]*Creature),
	}
}

func (m *Manager) AddCreature(c *Creature) {
	m.creatures[c.ID] = c
}

func (m *Manager) GetCreature(id string) *Creature {
	return m.creatures[id]
}

func (m *Manager) RemoveCreature(id string) {
	delete(m.creatures, id)
}

func (m *Manager) TickAll() {
	now := time.Now()

	for _, c := range m.creatures {
		if !c.IsAlive {
			if c.RespawnTimeSec > 0 && !c.TimeOfDeath.IsZero() && now.Sub(c.TimeOfDeath).Seconds() >= float64(c.RespawnTimeSec) {
				m.RespawnCreature(c)
			}
			continue
		}

		m.TickCreature(c)
	}
}

func (m *Manager) TickCreature(c *Creature) {
	// AI básica de exemplo: escolher ação aleatória se idle
	if c.CurrentAction == ActionIdle {
		if len(c.Actions) == 0 {
			return
		}
		newAction := c.Actions[rand.Intn(len(c.Actions))]
		c.CurrentAction = newAction
		log.Printf("[Creature %s] mudou para ação: %s", c.ID, newAction)
	}
}

func (m *Manager) KillCreature(c *Creature) {
	c.IsAlive = false
	c.TimeOfDeath = time.Now()
	c.CurrentAction = ActionDie
	log.Printf("[Creature %s] morreu. Respawn em %d segundos", c.ID, c.RespawnTimeSec)
}

func (m *Manager) RespawnCreature(c *Creature) {
	c.HP = c.MaxHP
	c.IsAlive = true
	c.TimeOfDeath = time.Time{}
	c.CurrentAction = ActionIdle
	log.Printf("[Creature %s] Respawned at position: %+v", c.ID, c.Position)
}
