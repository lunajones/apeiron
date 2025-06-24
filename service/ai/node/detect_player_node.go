package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

type DetectPlayerNode struct{}

func (n *DetectPlayerNode) Tick(c *creature.Creature, ctx dynamic_context.AIServiceContext) core.BehaviorStatus {
	log.Printf("[AI] %s executando DetectPlayerNode", c.ID)

	for _, p := range ctx.GetServicePlayers() {
		// Se a criatura estiver com medo, evita até olhar
		if c.MentalState == creature.MentalStateAfraid {
			log.Printf("[AI] %s está com medo, ignorando detecção de players.", c.ID)
			continue
		}

		// Se for pacífica (ex: merchant), evita engajamento
		if c.CurrentRole == creature.RoleMerchant {
			log.Printf("[AI] %s é um Merchant, evitando engajamento com player %s.", c.ID, p.ID)
			continue
		}

		// Se tem fome extrema e é predadora, vê jogadores como presa
		hunger := c.GetNeedValue(creature.NeedHunger)
		if hunger > 80 && c.HasTag(creature.TagPredator) {
			log.Printf("[AI] %s está faminto e vê %s como presa.", c.ID, p.ID)
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		// Comportamento padrão: ver ou ouvir o jogador
		if creature.CanSeePlayer(c, []*player.Player{p}) || creature.CanHearPlayer(c, []*player.Player{p}) {
			log.Printf("[AI] %s detectou o player %s.", c.ID, p.ID)
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}
