package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type SleepNode struct{}

func (n *SleepNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	log.Printf("[AI] [%s (%s)] está dormindo", c.Handle.String(), c.PrimaryType)
	c.ChangeAIState(consts.AIStateSleeping)
	creature.ModifyNeed(c, consts.NeedSleep, -2.0) // regeneração acelerada
	c.CurrentAction = consts.ActionIdle            // talvez no futuro: ActionSleep
	return core.StatusRunning
}

// Reset atende ao contrato da interface BehaviorNode
func (n *SleepNode) Reset() {
	// Nenhum estado interno a resetar no momento
}
