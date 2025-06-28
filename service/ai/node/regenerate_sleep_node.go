package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type RegenerateSleepNode struct{}

func (n *RegenerateSleepNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	sleepNeed := c.GetNeedByType(consts.NeedSleep)

	if sleepNeed == nil {
		log.Printf("[AI] [%s (%s)] erro: necessidade de sono não encontrada", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if sleepNeed.Value < sleepNeed.Threshold {
		log.Printf("[AI] [%s (%s)] descansou o suficiente (%.2f < %.2f), AIState: Idle",
			c.Handle.String(), c.PrimaryType, sleepNeed.Value, sleepNeed.Threshold)

		c.ChangeAIState(consts.AIStateIdle)
		return core.StatusSuccess
	}

	creature.ModifyNeed(c, consts.NeedSleep, -2.0)
	log.Printf("[AI] [%s (%s)] recuperando sono. Novo valor: %.2f", c.Handle.String(), c.PrimaryType, c.GetNeedValue(consts.NeedSleep))

	c.CurrentAction = consts.ActionIdle
	return core.StatusRunning
}

// Reset atende ao contrato da interface BehaviorNode
func (n *RegenerateSleepNode) Reset() {
	// Nenhum estado interno a resetar no momento
}
