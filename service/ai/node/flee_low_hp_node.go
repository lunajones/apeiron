package node

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type FleeIfLowHPNode struct{}

func (n *FleeIfLowHPNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	log.Printf("[FLEE LOW HP] [%s (%s)] executando FleeIfLowHPNode", c.Handle.String(), c.PrimaryType)

	hpThreshold := 30

	if c.HP > hpThreshold {
		log.Printf("[FLEE LOW HP] [%s (%s)] HP acima do limiar (%d > %d), não vai fugir", c.Handle.String(), c.PrimaryType, c.HP, hpThreshold)
		return core.StatusFailure
	}

	if c.MentalState == consts.MentalStateAggressive || c.MentalState == consts.MentalStateEnraged {
		log.Printf("[FLEE LOW HP] [%s (%s)] está agressivo/enraivecido, não vai fugir mesmo com HP baixo", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	hunger := c.GetNeedValue(consts.NeedHunger)
	if hunger > 90 && c.HasTag(consts.TagPredator) {
		log.Printf("[FLEE LOW HP] [%s (%s)] com fome extrema (%.2f), arrisca a vida por comida", c.Handle.String(), c.PrimaryType, hunger)
		return core.StatusFailure
	}

	log.Printf("[FLEE LOW HP] [%s (%s)] HP baixo e mentalState %s, iniciando fuga!", c.Handle.String(), c.PrimaryType, c.MentalState)
	c.SetAction(consts.ActionRun)
	c.ChangeAIState(consts.AIStateFleeing)

	c.Memory = append(c.Memory, creature.MemoryEvent{
		Description: "Fugiu com HP crítico",
		Timestamp:   time.Now(),
	})

	return core.StatusSuccess
}

func (n *FleeIfLowHPNode) Reset() {
	// Este node não mantém estado interno
}
