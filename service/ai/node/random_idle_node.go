package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomIdleNode struct{}

func (n *RandomIdleNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	log.Printf("[AI] %s executando RandomIdleNode", c.ID)

	// Type assertion correto para o contexto dinâmico
	_, ok := ctx.(dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[AI] %s: contexto inválido em RandomIdleNode", c.ID)
		return core.StatusFailure
	}

	// Checa necessidades básicas
	hunger := c.GetNeedValue(creature.NeedHunger)
	thirst := c.GetNeedValue(creature.NeedThirst)
	sleep := c.GetNeedValue(creature.NeedSleep)

	if hunger > 80 {
		log.Printf("[AI] %s está com fome alta, saindo do idle pra buscar comida.", c.ID)
		c.ChangeAIState(creature.AIStateSearchFood)
		return core.StatusFailure
	}

	if thirst > 80 {
		log.Printf("[AI] %s com sede alta, saindo do idle pra buscar água.", c.ID)
		c.ChangeAIState(creature.AIStateSearchWater)
		return core.StatusFailure
	}

	if sleep > 80 {
		log.Printf("[AI] %s está sonolento e prefere descansar.", c.ID)
		c.SetAction(creature.ActionSleep)
		return core.StatusSuccess
	}

	if c.CurrentRole == creature.RoleMerchant {
		log.Printf("[AI] %s é um Merchant, executando idle de comerciante.", c.ID)
		c.SetAction(creature.ActionWalk)
		return core.StatusSuccess
	}

	rand.Seed(time.Now().UnixNano())
	roll := rand.Float64()

	if roll < 0.3 {
		c.SetAction(creature.ActionWalk)
		log.Printf("[AI] %s escolheu andar aleatoriamente.", c.ID)
	} else {
		c.SetAction(creature.ActionIdle)
		log.Printf("[AI] %s permanece parado, observando o ambiente.", c.ID)
	}

	return core.StatusSuccess
}
