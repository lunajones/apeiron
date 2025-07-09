package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type MoveTowardsCorpseNode struct{}

// func (n *MoveTowardsCorpseNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] contexto inválido", c.GetHandle().ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	corpses := svcCtx.GetServiceCorpses(c.GetPosition(), c.DetectionRadius)

// 	var targetCorpse *creature.Creature
// 	for _, corpse := range corpses {
// 		if corpse.GetHandle().Equals(c.GetHandle()) || corpse.IsAlive || !corpse.IsCorpse {
// 			continue
// 		}
// 		if !creature.AreEnemies(c, corpse) {
// 			continue
// 		}
// 		targetCorpse = corpse
// 		break
// 	}

// 	if targetCorpse == nil {
// 		log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] nenhum cadáver válido encontrado", c.GetHandle().ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	stopAt := c.GetHitboxRadius() + targetCorpse.GetHitboxRadius() + c.GetDesiredBufferDistance()
// 	dist := position.CalculateDistance(c.GetPosition(), targetCorpse.GetPosition())

// 	if dist <= stopAt {
// 		log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] chegou no cadáver (%.2f ≤ %.2f)", c.GetHandle().ID, c.PrimaryType, dist, stopAt)
// 		return core.StatusSuccess
// 	}

// 	var speed float64
// 	if dist < 3.0 {
// 		speed = c.WalkSpeed * 0.5
// 		c.SetAction(consts.ActionWalk)
// 		log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] aproximando com cautela (speed: %.2f)", c.GetHandle().ID, c.PrimaryType, speed)
// 	} else {
// 		speed = c.RunSpeed
// 		c.SetAction(consts.ActionRun)
// 		log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] aproximando em corrida (speed: %.2f)", c.GetHandle().ID, c.PrimaryType, speed)
// 	}

// 	// Agora usamos intent ao invés de SetTarget + Update direto
// 	c.MoveCtrl.SetMoveIntent(targetCorpse.GetPosition(), speed, stopAt)

// 	log.Printf("[MOVE TOWARDS CORPSE] [%s (%s)] intenção de mover até cadáver (%.2f > %.2f)",
// 		c.GetHandle().ID, c.PrimaryType, dist, stopAt)

// 	return core.StatusRunning
// }

// func (n *MoveTowardsCorpseNode) Reset() {
// 	// Nada a resetar
// }
