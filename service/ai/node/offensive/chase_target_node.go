package offensive

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type ChaseTargetNode struct{}

func (n *ChaseTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.IsDodging() {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if drive.Caution < 0.4 && drive.Caution > 0.8 {
		log.Printf("[APPROACH-IN-RANGE] [%s] Caution baixo (%.2f), ignorando aproximação tática", c.Handle.String(), drive.Caution)
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[CHASE-TARGET] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if c.TargetCreatureHandle.IsEmpty() {
		log.Printf("[CHASE-TARGET] [%s (%s)] nenhum target lockado", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	var target *creature.Creature
	for _, t := range svcCtx.GetCachedTargets(c.Handle) {
		if other, ok := t.(*creature.Creature); ok {
			if other.GetHandle().Equals(c.TargetCreatureHandle) && other.Alive {
				target = other
				break
			}
		}
	}

	if target == nil {
		log.Printf("[CHASE-TARGET] [%s (%s)] target não encontrado ou morto",
			c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	stopAt := c.GetHitboxRadius() + target.GetHitboxRadius() + c.GetDesiredBufferDistance()
	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())

	if dist <= stopAt {
		c.SetAnimationState(constslib.AnimationIdle)
		c.MoveCtrl.IsMoving = false
		return core.StatusSuccess
	}

	if (!c.MoveCtrl.IsMoving || len(c.MoveCtrl.CurrentPath) == 0) && dist > stopAt {
		c.MoveCtrl.SetMoveTarget(target.GetPosition(), c.RunSpeed, stopAt)
		log.Printf("[CHASE-TARGET] [%s] Novo intent direto ao alvo criado. Dist=%.2f stopAt=%.2f",
			c.Handle.String(), dist, stopAt)
	}

	// Atualiza drive por estar perseguindo algo (fome, ação, estímulo)
	now := time.Now()
	drive.Rage += 0.02
	drive.Caution -= 0.01
	drive.Termination = 0
	drive.LastUpdated = now
	drive.Value = creature.RecalculateCombatDrive(drive)

	c.SetAnimationState(constslib.AnimationRun)
	log.Printf("[CHASE-TARGET] [%s] Perseguindo alvo. Dist=%.2f stopAt=%.2f",
		c.Handle.String(), dist, stopAt)
	return core.StatusRunning
}

func (n *ChaseTargetNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
}
