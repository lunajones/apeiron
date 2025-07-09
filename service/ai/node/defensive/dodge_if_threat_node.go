package defensive

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type DodgeIfThreatNode struct{}

func (n *DodgeIfThreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[DODGE-THREAT] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[DODGE-THREAT] [%s] sem alvo para avaliar dodge", c.Handle.String())
		return core.StatusFailure
	}

	if target.IsAlive() && target.IsHostile() {
		distance := position.CalculateDistance2D(c.Position, target.GetPosition())
		if distance < c.DesiredBufferDistance+1.0 {
			if ct, ok := target.(*creature.Creature); ok {
				if ct.CombatState == constslib.CombatStateAttacking {
					c.PerformDodge(svcCtx)

					// Seta estado e registra evento
					c.CombatState = constslib.CombatStateDodging
					svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
						SourceHandle: c.Handle,
						BehaviorType: "DodgePerformed",
						Timestamp:    time.Now(),
					})

					log.Printf("[DODGE-THREAT] [%s] dodge acionado contra [%s]", c.Handle.String(), target.GetHandle().String())
					return core.StatusSuccess
				}
			}
		}
	}

	log.Printf("[DODGE-THREAT] [%s] sem ameaça que justifique dodge", c.Handle.String())
	return core.StatusFailure
}

func (n *DodgeIfThreatNode) Reset() {}
