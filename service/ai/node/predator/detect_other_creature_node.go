package predator

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/handle"
// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type DetectOtherCreatureNode struct{}

// func (n *DetectOtherCreatureNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[AI] [%s (%s)] contexto inválido para DetectOtherCreatureNode (predator)", c.Handle.ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	if c.IsBlind {
// 		log.Printf("[AI] [%s (%s)] está cego, não pode detectar criaturas", c.Handle.ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	var bestTarget *creature.Creature
// 	var bestScore float64 = -1
// 	var bestDistance float64

// 	hunger := c.GetNeedValue(consts.NeedHunger)

// 	for _, other := range svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius) {
// 		if other.GetHandle().Equals(c.Handle) || !other.IsAlive {
// 			continue
// 		}

// 		log.Printf("[AI] [%s (%s)] avaliando %s (%s) para inimigo/prey",
// 			c.Handle.ID, c.PrimaryType, other.Handle.ID, other.PrimaryType)

// 		if !creature.AreEnemies(c, other) {
// 			log.Printf("[AI] [%s (%s)] não considera %s (%s) inimigo",
// 				c.Handle.ID, c.PrimaryType, other.Handle.ID, other.PrimaryType)
// 			continue
// 		}

// 		if other.HasTag(consts.TagPrey) && hunger < 40 {
// 			log.Printf("[AI] [%s (%s)] ignorou %s (%s) porque está saciado (fome: %d)",
// 				c.Handle.ID, c.PrimaryType, other.Handle.ID, other.PrimaryType, hunger)
// 			continue
// 		}

// 		dist := position.CalculateDistance(c.Position, other.Position)

// 		score := 0.0
// 		if other.HasTag(consts.TagPrey) {
// 			score += 2.0
// 		}
// 		if other.HP < 15 {
// 			score += 1.5
// 		}
// 		score += float64(hunger) / 100.0
// 		score += 1.0 - (dist / c.DetectionRadius)

// 		if score > bestScore {
// 			bestScore = score
// 			bestTarget = other
// 			bestDistance = dist
// 		}
// 	}

// 	if bestTarget == nil {
// 		c.TargetCreatureHandle = handle.EntityHandle{}
// 		return core.StatusFailure
// 	}

// 	if !c.TargetCreatureHandle.Equals(handle.EntityHandle{}) {
// 		current := svcCtx.FindCreatureByHandle(c.TargetCreatureHandle)
// 		if current != nil && current.IsAlive {
// 			currentDist := position.CalculateDistance(c.Position, current.Position)
// 			if !bestTarget.GetHandle().Equals(current.GetHandle()) {
// 				if bestTarget.HP < current.HP && bestDistance < currentDist*0.5 {
// 					log.Printf("[AI] [%s (%s)] trocando alvo: %s → %s (melhor oportunidade)",
// 						c.Handle.ID, c.PrimaryType, current.Handle.ID, bestTarget.Handle.ID)
// 					c.TargetCreatureHandle = bestTarget.GetHandle()
// 				} else {
// 					log.Printf("[AI] [%s (%s)] mantém alvo atual: %s",
// 						c.Handle.ID, c.PrimaryType, current.Handle.ID)
// 				}
// 			}
// 		} else {
// 			c.TargetCreatureHandle = bestTarget.GetHandle()
// 		}
// 	} else {
// 		c.TargetCreatureHandle = bestTarget.GetHandle()
// 	}

// 	if bestTarget.HasTag(consts.TagPrey) && hunger > 40 {
// 		log.Printf("[AI] [%s (%s)] detectou [%s (%s)] como presa (fome: %d, distância: %.2f) → AIState: Chasing",
// 			c.Handle.ID, c.PrimaryType, bestTarget.Handle.ID, bestTarget.PrimaryType, hunger, bestDistance)
// 		c.TargetCreatureHandle = bestTarget.GetHandle()
// 		c.ChangeAIState(consts.AIStateChasing)
// 	} else {
// 		log.Printf("[AI] [%s (%s)] detectou [%s (%s)] como inimigo (distância: %.2f) → AIState: Combat",
// 			c.Handle.ID, c.PrimaryType, bestTarget.Handle.ID, bestTarget.PrimaryType, bestDistance)
// 		c.TargetCreatureHandle = bestTarget.GetHandle()
// 		c.ChangeAIState(consts.AIStateCombat)
// 	}

// 	return core.StatusSuccess
// }

// func (n *DetectOtherCreatureNode) Reset() {}
