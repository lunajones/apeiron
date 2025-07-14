package defensive

// import (
// 	"log"
// 	"time"

// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/helper/finder"
// )

// type IsTargetVulnerableNode struct{}

// func (n *IsTargetVulnerableNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[DEF-VULN-CHECK] [%s] contexto inválido", c.Handle.String())
// 		return core.StatusFailure
// 	}

// 	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

// 	if target == nil {
// 		log.Printf("[DEF-VULN-CHECK] [%s] sem alvo para checar vulnerabilidade", c.Handle.String())
// 		return core.StatusFailure
// 	}

// 	isBlocking := target.IsBlocking()
// 	facing := target.GetFacingDirection()
// 	toTarget := c.Position.ToVector2D().Sub(target.GetPosition().ToVector2D()).Normalize()
// 	dot := facing.Dot(toTarget)

// 	if !isBlocking && dot < -0.7 {
// 		log.Printf("[DEF-VULN-CHECK] [%s] alvo [%s] vulnerável (não bloqueando e de costas)", c.Handle.String(), target.GetHandle().String())

// 		// Registra o evento
// 		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
// 			SourceHandle: c.Handle,
// 			BehaviorType: "TargetVulnerableDetected",
// 			Timestamp:    time.Now(),
// 		})

// 		return core.StatusSuccess
// 	}

// 	log.Printf("[DEF-VULN-CHECK] [%s] alvo [%s] não vulnerável", c.Handle.String(), target.GetHandle().String())
// 	return core.StatusFailure
// }

// func (n *IsTargetVulnerableNode) Reset(c *creature.Creature) {
// 	// Nada a resetar
// }
