package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// 	playerconsts "github.com/lunajones/apeiron/service/player/consts"
// )

// type DetectPlayerNode struct{}

// func (n *DetectPlayerNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[AI] [%s] contexto inválido em DetectPlayerNode", c.Handle.ID)
// 		return core.StatusFailure
// 	}

// 	if !c.IsAlive || c.IsBlind || c.IsDeaf {
// 		log.Printf("[AI] [%s] incapaz de detectar jogadores (morto/cego/surdo)", c.Handle.ID)
// 		return core.StatusFailure
// 	}

// 	players := svcCtx.GetServicePlayers(c.Position, c.DetectionRadius)
// 	if len(players) == 0 {
// 		log.Printf("[AI] [%s] não detectou jogadores no raio %.2f", c.Handle.ID, c.DetectionRadius)
// 		return core.StatusFailure
// 	}

// 	for _, p := range players {
// 		if !p.CheckIsAlive() {
// 			continue
// 		}

// 		dist := position.CalculateDistance(c.Position, p.Position)

// 		if p.CurrentRole == playerconsts.RoleMerchant {
// 			log.Printf("[AI] [%s] ignorou [%s] (comerciante)", c.Handle.ID, p.Handle.ID)
// 			continue
// 		}

// 		hunger := c.GetNeedValue(consts.NeedHunger)
// 		if hunger > 80 && c.HasTag(consts.TagPredator) {
// 			log.Printf("[AI] [%s] faminto, detectou [%s] como presa (%.2f)", c.Handle.ID, p.Handle.ID, dist)
// 			c.TargetPlayerHandle = p.Handle
// 			c.ChangeAIState(consts.AIStateAlert)
// 			return core.StatusSuccess
// 		}

// 		log.Printf("[AI] [%s] detectou [%s] (%.2f), iniciando alerta", c.Handle.ID, p.Handle.ID, dist)
// 		c.TargetPlayerHandle = p.Handle
// 		c.ChangeAIState(consts.AIStateAlert)
// 		return core.StatusSuccess
// 	}

// 	log.Printf("[AI] [%s] nenhum jogador relevante detectado", c.Handle.ID)
// 	return core.StatusFailure
// }

// func (n *DetectPlayerNode) Reset() {
// 	// Nada a resetar neste node
// }
