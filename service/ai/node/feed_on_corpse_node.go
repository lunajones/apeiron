package node

// import (
// 	"log"
// 	"time"

// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type FeedOnCorpseNode struct{}

// func (n *FeedOnCorpseNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[FEED ON CORPSE] [%s (%s)] contexto inválido em FeedOnCorpseNode", c.Handle.String(), c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	if !c.HasTag(consts.TagPredator) {
// 		log.Printf("[FEED ON CORPSE] [%s (%s)] não é predador, ignorando cadáveres", c.Handle.String(), c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	hunger := c.GetNeedByType(consts.NeedHunger)
// 	if hunger == nil || hunger.Value < hunger.Threshold {
// 		log.Printf("[FEED ON CORPSE] [%s (%s)] sem fome suficiente para se alimentar (%.2f < %.2f)", c.Handle.String(), c.PrimaryType, hunger.Value, hunger.Threshold)
// 		return core.StatusFailure
// 	}

// 	corpses := svcCtx.GetServiceCorpses(c.Position, c.DetectionRadius)

// 	for _, corpse := range corpses {
// 		if corpse.Handle.ID == c.Handle.ID || corpse.IsAlive || !corpse.IsCorpse {
// 			continue
// 		}

// 		dist := position.CalculateDistance(c.Position, corpse.Position)

// 		log.Printf("[FEED ON CORPSE] [%s (%s)] avaliando cadáver %s (distância %.2f, isCorpse: %t, isAlive: %t)",
// 			c.Handle.String(), c.PrimaryType, corpse.Handle.String(), dist, corpse.IsCorpse, corpse.IsAlive)

// 		if !creature.AreEnemies(c, corpse) {
// 			continue
// 		}

// 		if dist <= (c.HitboxRadius + corpse.HitboxRadius + 0.2) {
// 			log.Printf("[FEED ON CORPSE] [%s (%s)] encontrou cadáver inimigo próximo (%.2f unidades), vai se alimentar", c.Handle.String(), c.PrimaryType, dist)

// 			c.SetAction(consts.ActionSkill1)
// 			c.ChangeAIState(consts.AIStateFeeding)

// 			creature.ModifyNeed(c, consts.NeedHunger, -25)
// 			creature.ModifyNeed(c, consts.NeedSleep, 25)
// 			log.Printf("[FEED ON CORPSE] [%s (%s)] sono aumentado em 25 após se alimentar", c.Handle.String(), c.PrimaryType)

// 			corpse.ConsumeCorpse()

// 			c.Memory = append(c.Memory, creature.MemoryEvent{
// 				Description: "Alimentou-se de cadáver",
// 				Timestamp:   time.Now(),
// 			})

// 			return core.StatusSuccess
// 		}
// 	}

// 	log.Printf("[FEED ON CORPSE] [%s (%s)] não encontrou cadáveres adequados", c.Handle.String(), c.PrimaryType)
// 	return core.StatusFailure
// }

// func (n *FeedOnCorpseNode) Reset() {
// 	// Não há estado interno para resetar neste node
// }
