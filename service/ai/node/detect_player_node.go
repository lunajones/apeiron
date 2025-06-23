package node

import (
	"log"
	"github.com/lunajones/apeiron/service/player"


	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type DetectPlayerNode struct {
	Players []*player.Player
}

func (n *DetectPlayerNode) Tick(c *creature.Creature, ctx core.AIContext) core.BehaviorStatus {
	for _, p := range ctx.Players {
		// Se a criatura estiver com medo, ela evita até olhar
		if c.MentalState == creature.MentalStateAfraid {
			log.Printf("[AI] %s está com medo, ignorando detecção de players.", c.ID)
			continue
		}

		// Se for pacífica (ex: merchant), talvez não queira reagir
		if c.CurrentRole == creature.RoleMerchant {
			log.Printf("[AI] %s é um Merchant, evitando engajamento com player %s.", c.ID, p.ID)
			continue
		}

		// Se a criatura tem fome extrema e é predadora, ela passa a detectar jogadores como presa
		hunger := c.GetNeedValue(creature.NeedHunger)
		if hunger > 80 && c.HasTag(creature.TagPredator) {
			log.Printf("[AI] %s está faminto e vê %s como presa.", c.ID, p.ID)
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}

		// Comportamento padrão: detectar se consegue ver ou ouvir
		if creature.CanSeePlayer(c, []*player.Player{p}) || creature.CanHearPlayer(c, []*player.Player{p}) {
			log.Printf("[AI] %s detectou o player %s.", c.ID, p.ID)
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	return core.StatusFailure
}
