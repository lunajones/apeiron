package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type MaintainMediumDistanceNode struct{}

// func (n *MaintainMediumDistanceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[AI] [%s (%s)] contexto inválido em MaintainMediumDistanceNode", c.Handle.String(), c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	log.Printf("[AI] [%s (%s)] executando MaintainMediumDistanceNode", c.Handle.String(), c.PrimaryType)

// 	players := svcCtx.GetServicePlayers(c.GetPosition(), c.DetectionRadius)
// 	if len(players) == 0 {
// 		log.Printf("[AI] [%s (%s)] nenhum jogador próximo detectado", c.Handle.String(), c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	target := players[0]
// 	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
// 	stopAt := c.GetHitboxRadius() + c.GetDesiredBufferDistance()

// 	if dist < 4.0 {
// 		log.Printf("[AI] [%s (%s)] muito próximo de %s (%.2f u), recuando", c.Handle.String(), c.PrimaryType, target.Handle.String(), dist)

// 		dir := position.NewVector3DFromTo(target.GetPosition(), c.GetPosition())
// 		if dir.Magnitude() == 0 {
// 			dir = position.Vector3D{X: 1, Y: 0, Z: 0}
// 		}
// 		dir = dir.Normalize()

// 		newPos := position.FromGlobal(
// 			c.GetPosition().FastGlobalX()+dir.X*2.0,
// 			c.GetPosition().FastGlobalZ()+dir.Z*2.0,
// 			c.GetPosition().Y+dir.Y*2.0,
// 		)

// 		c.MoveCtrl.SetMoveIntent(newPos, c.RunSpeed, stopAt)
// 		c.SetAction(consts.ActionRun)

// 		return core.StatusRunning

// 	} else if dist > 8.0 {
// 		log.Printf("[AI] [%s (%s)] muito longe de %s (%.2f u), aproximando", c.Handle.String(), c.PrimaryType, target.Handle.String(), dist)

// 		c.MoveCtrl.SetMoveIntent(target.GetPosition(), c.RunSpeed, stopAt)
// 		c.SetAction(consts.ActionRun)

// 		return core.StatusRunning
// 	}

// 	log.Printf("[AI] [%s (%s)] distância ideal com %s (%.2f u), atacando", c.Handle.String(), c.PrimaryType, target.Handle.String(), dist)
// 	c.SetAction(consts.ActionSkill2)
// 	c.ChangeAIState(consts.AIStateAttack)

// 	return core.StatusSuccess
// }

// func (n *MaintainMediumDistanceNode) Reset() {
// 	// Este node não mantém estado interno
// }
