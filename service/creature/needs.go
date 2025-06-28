package creature

import (
	"log"
	"math/rand"

	"github.com/lunajones/apeiron/service/creature/consts"
)

func (c *Creature) SetNeedValue(needType consts.NeedType, value float64) {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			c.Needs[i].Value = value
			return
		}
	}
}

func (c *Creature) GetNeedValue(needType consts.NeedType) float64 {
	for _, n := range c.Needs {
		if n.Type == needType {
			return n.Value
		}
	}
	return 0
}

func ModifyNeed(c *Creature, needType consts.NeedType, amount float64) {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			old := c.Needs[i].Value
			newVal := old + amount

			if newVal < 0 {
				newVal = 0
			}
			if newVal > 100 {
				newVal = 100
			}

			// Só aplica e loga se houver mudança real
			if newVal != old {
				c.Needs[i].Value = newVal
				log.Printf("[Creature] %s (%s) teve %s modificada de %.2f → %.2f", c.Handle.String(), c.PrimaryType, needType, old, newVal)
			}

			break
		}
	}
}

func (c *Creature) GetNeedThreshold(needType consts.NeedType) float64 {
	for _, n := range c.Needs {
		if n.Type == needType {
			return n.Threshold
		}
	}
	return 100 // default conservador
}

type NeedDefaults struct {
	Min       float64
	Max       float64
	Threshold float64
}

var GlobalNeedDefaults = map[consts.NeedType]NeedDefaults{
	consts.NeedHunger: {Min: 10, Max: 40, Threshold: 30},
	consts.NeedSleep:  {Min: 20, Max: 50, Threshold: 40},
	consts.NeedThirst: {Min: 5, Max: 30, Threshold: 25},
	consts.NeedSocial: {Min: 0, Max: 50, Threshold: 45},
	consts.NeedFuck:   {Min: 10, Max: 40, Threshold: 30},
	consts.NeedKill:   {Min: 5, Max: 30, Threshold: 25},
	consts.NeedDrink:  {Min: 5, Max: 30, Threshold: 20},
}

var CreatureNeedDefaults = map[consts.CreatureType]map[consts.NeedType]NeedDefaults{
	consts.Rabbit: {
		consts.NeedHunger: {Min: 30, Max: 60, Threshold: 50},
		consts.NeedSleep:  {Min: 40, Max: 70, Threshold: 60},
	},
	consts.Wolf: {
		consts.NeedHunger: {Min: 20, Max: 50, Threshold: 40},
		consts.NeedSleep:  {Min: 10, Max: 35, Threshold: 25},
	},
}

func (c *Creature) ResetNeeds() {
	for i := range c.Needs {
		needType := c.Needs[i].Type
		if speciesDefaults, ok := CreatureNeedDefaults[c.PrimaryType]; ok {
			if def, ok := speciesDefaults[needType]; ok {
				c.Needs[i].Value = def.Min + rand.Float64()*(def.Max-def.Min)
				continue
			}
		}
		// fallback pro global
		if def, ok := GlobalNeedDefaults[needType]; ok {
			c.Needs[i].Value = def.Min + rand.Float64()*(def.Max-def.Min)
		} else {
			c.Needs[i].Value = rand.Float64() * 30
		}
	}
}

func (c *Creature) GetNeedByType(needType consts.NeedType) *consts.Need {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			return &c.Needs[i]
		}
	}
	return nil
}
