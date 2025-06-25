package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/lib/position"
)

type DetectPlayerNode struct{}

func (n *DetectPlayerNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx := ctx.(dynamic_context.AIServiceContext)

	log.Printf("[AI] %s executando DetectPlayerNode", c.ID)

	if !c.IsAlive || c.IsBlind || c.IsDeaf {
		log.Printf("[AI] %s está incapacitado para detectar jogadores.", c.ID)
		return core.StatusFailure
	}

	for _, p := range svcCtx.GetServicePlayers() {
		if !p.IsAlive {
			continue
		}
		dist := position.CalculateDistance(c.Position, p.Position)
		if dist <= c.DetectionRadius {
			if p.CurrentRole == player.RoleMerchant {
				log.Printf("[AI] %s ignorou %s por ser comerciante.", c.ID, p.ID)
				continue
			}

			hunger := c.GetNeedValue(creature.NeedHunger)
			if hunger > 80 && c.HasTag(creature.TagPredator) {
				log.Printf("[AI] %s está faminto e detectou %s como possível alvo.", c.ID, p.ID)
				c.TargetPlayerID = p.ID
				c.ChangeAIState(creature.AIStateAlert)
				return core.StatusSuccess
			}

			log.Printf("[AI] %s detectou %s, iniciando perseguição.", c.ID, p.ID)
			c.TargetPlayerID = p.ID
			c.ChangeAIState(creature.AIStateAlert)
			return core.StatusSuccess
		}
	}

	log.Printf("[AI] %s não detectou nenhum jogador próximo.", c.ID)
	return core.StatusFailure
}
