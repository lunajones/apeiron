package creature

import (
	"math/rand"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/creature/consts"
)

func (c *Creature) SetNeedValue(needType constslib.NeedType, value float64) {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			c.Needs[i].Value = value
			return
		}
	}
}

func (c *Creature) GetNeedValue(needType constslib.NeedType) float64 {
	for _, n := range c.Needs {
		if n.Type == needType {
			return n.Value
		}
	}
	return 0
}

func ModifyNeed(c *Creature, needType constslib.NeedType, amount float64) {
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

			if newVal != old {
				c.Needs[i].Value = newVal

				// Só loga se for Need de combate
				if needType == constslib.NeedAdvance ||
					needType == constslib.NeedGuard ||
					needType == constslib.NeedRetreat ||
					needType == constslib.NeedRage ||
					needType == constslib.NeedPlan ||
					needType == constslib.NeedFake {
					// log.Printf("[Creature] %s (%s) teve %s modificada de %.2f → %.2f",
					// 	c.Handle.String(), c.PrimaryType, needType, old, newVal)
				}
			}

			break
		}
	}
}

func (c *Creature) GetNeedThreshold(needType constslib.NeedType) float64 {
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

var GlobalNeedDefaults = map[constslib.NeedType]NeedDefaults{
	constslib.NeedHunger:  {Min: 10, Max: 40, Threshold: 30},
	constslib.NeedSleep:   {Min: 20, Max: 50, Threshold: 40},
	constslib.NeedThirst:  {Min: 5, Max: 30, Threshold: 25},
	constslib.NeedSocial:  {Min: 0, Max: 50, Threshold: 45},
	constslib.NeedFuck:    {Min: 10, Max: 40, Threshold: 30},
	constslib.NeedKill:    {Min: 5, Max: 30, Threshold: 25},
	constslib.NeedDrink:   {Min: 5, Max: 30, Threshold: 20},
	constslib.NeedProvoke: {Min: 0, Max: 10, Threshold: 8},
	constslib.NeedAdvance: {Min: 0, Max: 10, Threshold: 7},
	constslib.NeedGuard:   {Min: 0, Max: 10, Threshold: 7},
}

var CreatureNeedDefaults = map[consts.CreatureType]map[constslib.NeedType]NeedDefaults{
	consts.Rabbit: {
		constslib.NeedHunger: {Min: 30, Max: 60, Threshold: 50},
		constslib.NeedSleep:  {Min: 40, Max: 70, Threshold: 60},
	},
	consts.Wolf: {
		constslib.NeedHunger: {Min: 20, Max: 50, Threshold: 40},
		constslib.NeedSleep:  {Min: 10, Max: 35, Threshold: 25},
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

func (c *Creature) GetNeedByType(needType constslib.NeedType) *constslib.Need {
	for i := range c.Needs {
		if c.Needs[i].Type == needType {
			return &c.Needs[i]
		}
	}
	return nil
}
