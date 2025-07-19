package offensive

import (
	"log"
	"math/rand"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CheckSkillRangeNode struct{}

func (n *CheckSkillRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.NextSkillToUse == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] Sem skill planejada, permite próximo node",
			c.Handle.String(), c.PrimaryType)
		return core.StatusSuccess
	}

	if c.MoveCtrl.IsMoving {
		return core.StatusRunning
	}

	if c.IsInTacticalMovement() {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] em movimento tático (ex: Circle), impedido de castar agora",
			c.Handle.String(), c.PrimaryType)
		return core.StatusRunning
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] nenhum alvo válido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		drive := c.GetCombatDrive()
		stamina := c.Stamina
		blocking := target.IsBlocking()
		recentMiss := c.HasRecentlyMissedSkill()
		hpRatio := float64(c.HP) / 100.0

		// Hesita por miss recente + cautela
		if recentMiss && (drive.Caution > drive.Rage || drive.Counter < 0.3) {
			c.CombatState = constslib.CombatStateMoving

			log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] hesitando cast por miss recente + cautela/counter (%.2f > %.2f | counter %.2f)",
				c.Handle.String(), c.PrimaryType, drive.Caution, drive.Rage, drive.Counter)
			return core.StatusRunning
		}

		// Hesita por alvo bloqueando + cautela
		if blocking && stamina > 10.0 && drive.Caution > drive.Rage {
			c.CombatState = constslib.CombatStateMoving

			log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] hesitando cast por bloqueio do alvo + cautela alta (stamina %.1f)",
				c.Handle.String(), c.PrimaryType, stamina)
			return core.StatusRunning
		}

		// HP alto + Counter baixo → chance de hesitar
		if hpRatio > 0.7 && drive.Counter < 0.2 {
			c.CombatState = constslib.CombatStateMoving

			rand.Seed(time.Now().UnixNano())
			if rand.Float64() < 0.4 {
				log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] hesitando cast por HP alto (%.2f) + counter baixo (%.2f) + chance",
					c.Handle.String(), c.PrimaryType, hpRatio, drive.Counter)
				return core.StatusRunning
			}
		}

		// Tudo certo, pode castar
		c.ClearMovementIntent()
		c.CombatState = constslib.CombatStateCasting
		log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] dentro do range (%.2f <= %.2f), castando skill",
			c.Handle.String(), c.PrimaryType, dist, rangeNeeded)
		return core.StatusSuccess
	}

	log.Printf("[CHECK-SKILL-RANGE] [%s (%s)] fora do range (%.2f > %.2f), mantendo movimento",
		c.Handle.String(), c.PrimaryType, dist, rangeNeeded)

	c.CombatState = constslib.CombatStateMoving
	return core.StatusRunning
}

func (n *CheckSkillRangeNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
}
