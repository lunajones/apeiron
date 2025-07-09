package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type MoveTowardsTargetWithNeedFallbackNode struct {
// 	StopDistance   float64
// 	PriorityNeeds  []consts.NeedType
// 	CheckOnlyThese []consts.NeedType
// }

// func (n *MoveTowardsTargetWithNeedFallbackNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[MOVE WITH NEED FALLBACK] [%s (%s)] contexto inválido", c.GetHandle().ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	var target *creature.Creature
// 	for _, other := range svcCtx.GetServiceCreatures(c.GetPosition(), c.DetectionRadius) {
// 		if other.GetHandle().Equals(c.TargetCreatureHandle) {
// 			target = other
// 			break
// 		}
// 	}

// 	if target == nil || !target.IsAlive {
// 		log.Printf("[MOVE WITH NEED FALLBACK] [%s (%s)] alvo inválido/morto, limpando alvo e avaliando necessidades",
// 			c.GetHandle().ID, c.PrimaryType)
// 		c.ClearTargetHandles()

// 		needsEvaluator := &EvaluateNeedsNode{
// 			PriorityOrder:  n.PriorityNeeds,
// 			CheckOnlyThese: n.CheckOnlyThese,
// 		}
// 		if needsEvaluator.Tick(c, ctx) == core.StatusSuccess {
// 			log.Printf("[MOVE WITH NEED FALLBACK] [%s (%s)] fallback de necessidade acionado", c.GetHandle().ID, c.PrimaryType)
// 			return core.StatusFailure
// 		}

// 		c.ChangeAIState(consts.AIStateIdle)
// 		return core.StatusFailure
// 	}

// 	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
// 	if dist <= n.StopDistance {
// 		log.Printf("[MOVE WITH NEED FALLBACK] [%s (%s)] já na distância do alvo (%.2f ≤ %.2f)",
// 			c.GetHandle().ID, c.PrimaryType, dist, n.StopDistance)
// 		return core.StatusSuccess
// 	}

// 	// Agora usamos intent ao invés de SetTarget + Update direto
// 	c.MoveCtrl.SetMoveIntent(target.GetPosition(), c.RunSpeed, n.StopDistance)
// 	c.SetAction(consts.ActionRun)

// 	log.Printf("[MOVE WITH NEED FALLBACK] [%s (%s)] intenção de mover em direção ao alvo (%.2f > %.2f)",
// 		c.GetHandle().ID, c.PrimaryType, dist, n.StopDistance)

// 	return core.StatusRunning
// }

// func (n *MoveTowardsTargetWithNeedFallbackNode) Reset() {
// 	// Nada a resetar
// }
