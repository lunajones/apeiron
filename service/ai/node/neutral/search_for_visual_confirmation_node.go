package neutral

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/ai/sensor"
	"github.com/lunajones/apeiron/service/creature"
)

type SearchForVisualConfirmationNode struct{}

func (n *SearchForVisualConfirmationNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)

	// if c.IsMovementLocked() {
	// 	log.Printf("[CHASE-IN-RANGE] [%s] movimento travado até %v", c.Handle.String(), c.GetMovementLockUntil())
	// 	return core.StatusRunning
	// }

	if !ok {
		log.Printf("[VISUAL-CONFIRM] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	target := svcCtx.FindByHandle(c.TargetCreatureHandle)
	if target == nil {
		return core.StatusFailure
	}

	creatureTarget, ok := target.(*creature.Creature)
	if !ok {
		log.Printf("[VISUAL-CONFIRM] [%s (%s)] alvo não é uma criatura válida", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	if sensor.CanSee(c, creatureTarget) {
		hunger := c.GetNeedValue(consts.NeedHunger)
		thirst := c.GetNeedValue(consts.NeedThirst)
		sleep := c.GetNeedValue(consts.NeedSleep)

		// Aplicar fórmula que favorece valores baixos
		// ex: (100 - need) / 100 → vai de 1 (need=0) até 0 (need=100)
		// pode multiplicar esse fator para ajustar a influência
		hungerBoost := (100 - hunger) / 100.0 * 0.015 // aumenta impacto para necessidades urgentes
		thirstBoost := (100 - thirst) / 100.0 * 0.010
		sleepBoost := (100 - sleep) / 100.0 * 0.007

		drive := c.GetCombatDrive()
		drive.Rage += hungerBoost
		drive.Caution += thirstBoost + sleepBoost
		drive.Termination = 0
		drive.LastUpdated = time.Now()
		drive.Value = creature.RecalculateCombatDrive(drive)

		log.Printf("%s", color.New(color.FgHiYellow).Sprintf(
			"[VISUAL-CONFIRM] [%s (%s)] visão confirmada → ajustando Drive (Rage=%.2f, Caution=%.2f, Value=%.2f)",
			c.Handle.String(), c.PrimaryType, drive.Rage, drive.Caution, drive.Value,
		))

		return core.StatusRunning
	}

	if time.Since(c.LastThreatSeen) > 5*time.Second {
		c.CombatState = consts.CombatStateStrategic
		log.Printf("%s", color.New(color.FgHiCyan).Sprintf(
			"[VISUAL-CONFIRM] [%s (%s)] ameaça não confirmada, retornando para Strategic",
			c.Handle.String(), c.PrimaryType))
		return core.StatusSuccess
	}

	return core.StatusRunning
}

func (n *SearchForVisualConfirmationNode) Reset(c *creature.Creature) {}
