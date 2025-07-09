package offensive

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

type ChaseUntilInRangeNode struct{}

func (n *ChaseUntilInRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.NextSkillToUse == nil {
		log.Printf("[CHASE-IN-RANGE] [%s (%s)] Sem skill planejada", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[CHASE-IN-RANGE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[CHASE-IN-RANGE] [%s (%s)] alvo não encontrado", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		c.ClearMovementIntent()
		c.SetAnimationState(constslib.AnimationIdle)
		log.Printf("[CHASE-IN-RANGE] [%s] no range para %s (dist=%.2f)", c.Handle.String(), c.NextSkillToUse.Name, dist)
		return core.StatusSuccess
	}

	// Se não tem intent ativo, ou está parado, ou o alvo se mexeu, inicia perseguição
	if !c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0 {
		stopAt := rangeNeeded
		c.MoveCtrl.SetMoveIntent(target.GetPosition(), c.RunSpeed, stopAt)
		log.Printf("[CHASE-IN-RANGE] [%s] novo intent: alvo fora do range (dist=%.2f, range=%.2f)", c.Handle.String(), dist, rangeNeeded)

		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "ChasePerformed",
			Timestamp:    time.Now(),
		})
	}

	c.SetAnimationState(constslib.AnimationRun)
	log.Printf("[CHASE-IN-RANGE] [%s] perseguindo. Dist=%.2f target=%s", c.Handle.String(), dist, target.GetHandle().String())

	return core.StatusRunning
}

func (n *ChaseUntilInRangeNode) Reset() {
	// Nada a resetar
}
