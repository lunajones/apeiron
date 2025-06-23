package node

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type FleeIfLowHPNode struct{}

func (n *FleeIfLowHPNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando FleeIfLowHPNode", c.ID)

	hpThreshold := 30

	// 1. Se HP ainda tá acima do limite, não precisa fugir
	if c.HP > hpThreshold {
		log.Printf("[AI] %s falhou ao executar FleeIfLowHPNode", c.ID)
		return core.StatusFailure
	}

	// 2. Estado mental pode influenciar: agressivos ou enraivecidos relutam em fugir
	if c.MentalState == creature.MentalStateAggressive || c.MentalState == creature.MentalStateEnraged {
		log.Printf("[AI] %s está agressivo/enraivecido, mesmo com HP baixo, vai continuar lutando.", c.ID)
		return core.StatusFailure
	}

	// 3. Se criatura estiver com fome extrema, e o target for prey, talvez lute até a morte
	hunger := c.GetNeedValue(creature.NeedHunger)
	if hunger > 90 && c.HasTag(creature.TagPredator) {
		log.Printf("[AI] %s está morrendo de fome e vai arriscar a vida por comida.", c.ID)
		return core.StatusFailure
	}

	// 4. Comportamento normal de fuga
	log.Printf("[AI] %s com HP baixo e estado mental %s, tentando fugir!", c.ID, c.MentalState)
	c.SetAction(creature.ActionRun)
	c.ChangeAIState(creature.AIStateAlert)

	// Opcional: grava na memória que já esteve em perigo
	c.Memory = append(c.Memory, creature.MemoryEvent{
		Description: "Fugiu com HP crítico",
		Timestamp: time.Now(), // Garanta que o AIContext tenha um campo Now com time.Now()
	})

	return core.StatusSuccess
}
