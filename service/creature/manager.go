package creature

import (
    "log"
    "math/rand"
)

var allCreatures = make(map[string]*Creature)

func AddCreature(c *Creature) {
    allCreatures[c.ID] = c
}

func GetCreature(id string) *Creature {
    return allCreatures[id]
}

func TickAll() {
    for _, c := range allCreatures {
        TickCreature(c)
    }
}

func TickCreature(c *Creature) {
    // Lógica de AI super simples por enquanto: escolher ação aleatória
    if c.CurrentAction == ActionIdle {
        possibleActions := c.Actions
        newAction := possibleActions[rand.Intn(len(possibleActions))]
        c.CurrentAction = newAction
        log.Printf("[Creature %s] mudou para ação: %s", c.ID, newAction)
    }
}
