package node

import (
	"log"
	"sort"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type EvaluateNeedsNode struct {
	// Ordem de prioridade: NeedType mais importante primeiro
	PriorityOrder []consts.NeedType
	// Se definido, ignora todas as necessidades fora dessa lista
	CheckOnlyThese []consts.NeedType
}

func (n *EvaluateNeedsNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	filteredNeeds := []consts.Need{}
	for _, need := range c.Needs {
		if n.shouldCheckNeed(need.Type) {
			filteredNeeds = append(filteredNeeds, need)
		}
	}

	sort.SliceStable(filteredNeeds, func(i, j int) bool {
		return n.indexOf(filteredNeeds[i].Type) < n.indexOf(filteredNeeds[j].Type)
	})

	for _, need := range filteredNeeds {
		if need.Value < need.Threshold {
			continue
		}
		switch need.Type {
		case consts.NeedHunger:
			log.Printf("[AI] [%s (%s)] com fome (%.2f ≥ %.2f), AIState: SearchFood", c.Handle.String(), c.PrimaryType, need.Value, need.Threshold)
			c.ChangeAIState(consts.AIStateSearchFood)
			return core.StatusSuccess

		case consts.NeedSleep:
			if !c.LastThreatSeen.IsZero() && time.Since(c.LastThreatSeen) < 15*time.Second {
				log.Printf("[AI] [%s (%s)] sono adiado: ameaça vista há %.1fs", c.Handle.String(), c.PrimaryType, time.Since(c.LastThreatSeen).Seconds())
				return core.StatusFailure
			}
			log.Printf("[AI] [%s (%s)] com sono (%.2f ≥ %.2f), AIState: Drowsy", c.Handle.String(), c.PrimaryType, need.Value, need.Threshold)
			c.ChangeAIState(consts.AIStateDrowsy)
			return core.StatusSuccess

		case consts.NeedThirst:
			log.Printf("[AI] [%s (%s)] com sede (%.2f ≥ %.2f), AIState: SearchWater", c.Handle.String(), c.PrimaryType, need.Value, need.Threshold)
			c.ChangeAIState(consts.AIStateSearchWater)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}

func (n *EvaluateNeedsNode) shouldCheckNeed(t consts.NeedType) bool {
	if len(n.CheckOnlyThese) == 0 {
		return true
	}
	for _, allowed := range n.CheckOnlyThese {
		if allowed == t {
			return true
		}
	}
	return false
}

func (n *EvaluateNeedsNode) indexOf(t consts.NeedType) int {
	for i, needType := range n.PriorityOrder {
		if needType == t {
			return i
		}
	}
	return len(n.PriorityOrder)
}

func (n *EvaluateNeedsNode) Reset() {
	// Não há estado interno para resetar neste node
}
