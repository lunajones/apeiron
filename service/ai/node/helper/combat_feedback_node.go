package helper

import (
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type CombatFeedbackNode struct{}

func (n *CombatFeedbackNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[COMBAT-FEEDBACK] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	eventsAsTarget := svcCtx.GetRecentCombatBehaviorsAsTarget(c.Handle, time.Now().Add(-2*time.Second))

	for _, evt := range eventsAsTarget {
		if evt.BehaviorType == "AggressiveIntention" {
			boost := 5 + rand.Float64()*5
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost)
			creature.ModifyNeed(c, constslib.NeedGuard, boost)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.6)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.2)

			log.Printf("[FEEDBACK] [%s] recebeu intenção hostil → +%.1f Guard, +%.1f Plan, +%.1f Retreat",
				c.Handle.String(), boost*0.2, boost*0.1, boost*0.1)
		}
	}

	// 📝 Processa eventos recentes
	events := svcCtx.GetRecentCombatBehaviors(c.Handle, time.Now().Add(-2*time.Second))
	for _, e := range events {
		switch e.BehaviorType {

		case "AggressiveIntention":
			boost := 3 + rand.Float64()*5 // Total: ~3-8
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*0.02)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*0.02)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.2)
			log.Printf("%s", color.New(color.FgHiMagenta).Sprintf(
				"[FEEDBACK] [%s] skill ofensiva planejada → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.5, boost*0.2, boost*0.1, boost*0.2, boost*0.2, boost*0.2))

		case "OffensiveSkillPlanned":
			boost := 3 + rand.Float64()*5 // Total: ~3-8
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*0.02)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*0.02)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.2)
			log.Printf("%s", color.New(color.FgHiMagenta).Sprintf(
				"[FEEDBACK] [%s] skill ofensiva planejada → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.5, boost*0.2, boost*0.1, boost*0.2, boost*0.2, boost*0.2))

		case "SkillExecuted":
			boost := 3 + rand.Float64()*5 // Total: ~3-8
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost*0.02)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*0.02)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.2)
			log.Printf("%s", color.New(color.FgHiBlue).Sprintf(
				"[FEEDBACK] [%s] skill executada → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.5, boost*0.2, boost*0.1, boost*0.4, boost*0.2, boost*0.2))

		case "SkillCycleCompleted":
			boost := 2 + rand.Float64()*4 // Total: ~2-6
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost*0.02)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*0.02)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.3)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.1)
			log.Printf("%s", color.New(color.FgHiMagenta).Sprintf(
				"[FEEDBACK] [%s] ciclo da skill concluído → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.5, boost*0.2, boost*0.1, boost*0.3, boost*0.2, boost*0.1))

		case "ChasePerformed":
			boost := 1 + rand.Float64()*2 // Total: ~2-6
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*0.0005)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*0.0005)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.0005)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.0005)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.0005)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.0005)
			log.Printf("%s", color.New(color.FgHiCyan).Sprintf(
				"[FEEDBACK] [%s] perseguição realizada → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.5, boost*0.2, boost*0.1, boost*0.3, boost*0.2, boost*0.1))

		case "DamageApplied":
			boost := 3 + rand.Float64()*5 // Total: 10-20 (média ~15)
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.1)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.1)
			log.Printf("%s", color.New(color.FgHiRed).Sprintf(
				"[FEEDBACK] [%s] aplicou dano → +%.1f Guard, -%.1f Advance, +%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.3, boost*0.3, boost*0.1, boost*0.2, boost*0.2, boost*0.1))

		case "DamageAvoided":
			boost := 5 + rand.Float64()*8 // Total: 10-20 (média ~15)
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*1)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*1)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, -boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.2)
			log.Printf("%s", color.New(color.FgHiGreen).Sprintf(
				"[FEEDBACK] [%s] evitou dano → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*2, boost*2, boost*0.5, boost*0.3, boost*0.3, boost*0.2))

		case "FlankExecuted":
			boost := 10 + rand.Float64()*20
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*1)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*1)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, -boost)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.5)
			log.Printf("%s", color.New(color.FgHiBlue).Sprintf(
				"[FEEDBACK] [%s] flanco → +%.1f Plan, +%.1f Fake, +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage",
				c.Handle.String(), boost, boost*0.5, boost*0.3, boost*0.2, boost*0.2, boost*0.4))

		case "DodgePerformed":
			boost := 15 + rand.Float64()*25
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*1)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*1)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, -boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.2)
			log.Printf("%s", color.New(color.FgHiYellow).Sprintf(
				"[FEEDBACK] [%s] esquiva realizada → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost, boost*0.5, boost*0.3, boost*0.2, boost*0.2, boost*0.2))

		case "FakeAdvancePerformed":
			boost := 10 + rand.Float64()*15
			creature.ModifyNeed(c, constslib.NeedAdvance, boost*2)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost*2)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.5)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.3)
			creature.ModifyNeed(c, constslib.NeedFake, boost)
			log.Printf("%s", color.New(color.FgHiWhite).Sprintf(
				"[FEEDBACK] [%s] fake advance → +%.1f Fake, +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan",
				c.Handle.String(), boost, boost*0.2, boost*0.1, boost*0.1, boost*0.2, boost*0.3))

		case "ReactedToFakeAdvance":
			boost := 5 + rand.Float64()*10
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost*2)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*2)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedRage, boost*5)
			creature.ModifyNeed(c, constslib.NeedPlan, -boost*0.3)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.2)
			log.Printf("%s", color.New(color.FgHiWhite).Sprintf(
				"[FEEDBACK] [%s] reagiu a FakeAdvance → +%.1f Guard, +%.1f Retreat, -%.1f Advance, +%.1f Rage, +%.1f Plan, -%.1f Fake",
				c.Handle.String(), boost, boost*0.5, boost*0.5, boost*0.1, boost*0.3, boost*0.1))

		case "TargetRetreatingDetected":
			boost := 6 + rand.Float64()*4 // Total: 6 a 10 → média ~8
			creature.ModifyNeed(c, constslib.NeedAdvance, boost)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, -boost*0.3)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.1)
			log.Printf("%s", color.New(color.FgHiGreen).Sprintf(
				"[FEEDBACK] [%s] detectou alvo recuando → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost, boost*0.2, boost*0.2, boost*0.4, boost*0.3, boost*0.1))

		case "TargetDefendingDetected":
			boost := 6 + rand.Float64()*4 // Total: 6 a 10 → média ~8
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost*0.5)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedPlan, boost)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.2)
			log.Printf("%s", color.New(color.FgHiYellow).Sprintf(
				"[FEEDBACK] [%s] detectou alvo defendendo → +%.1f Plan, -%.1f Advance, +%.1f Guard, +%.1f Rage, +%.1f Retreat, +%.1f Fake",
				c.Handle.String(), boost, boost*0.5, boost*0.5, boost*0.4, boost*0.2, boost*0.2))

		case "TargetVulnerableDetected":
			boost := 20 + rand.Float64()*10 // Total: 20 a 30 → média ~25
			creature.ModifyNeed(c, constslib.NeedAdvance, boost)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedPlan, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.1)
			log.Printf("%s", color.New(color.FgHiGreen).Sprintf(
				"[FEEDBACK] [%s] detectou alvo vulnerável → +%.1f Advance, -%.1f Guard, -%.1f Retreat, +%.1f Rage, -%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost, boost*0.2, boost*0.2, boost*0.5, boost*0.2, boost*0.1))

		case "DefendPerformed":
			boost := 6 + rand.Float64()*4 // Total: 6 a 10 → média ~8
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost)
			creature.ModifyNeed(c, constslib.NeedGuard, boost)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.3)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.1)
			log.Printf("%s", color.New(color.FgHiBlue).Sprintf(
				"[FEEDBACK] [%s] bloqueio realizado → +%.1f Guard, -%.1f Advance, +%.1f Retreat, +%.1f Rage, +%.1f Plan, -%.1f Fake",
				c.Handle.String(), boost, boost*0.5, boost*0.3, boost*0.2, boost*0.2, boost*0.1))

		case "ParryPerformed":
			boost := 20 + rand.Float64()*10 // Total: 20 a 30 → média ~25
			creature.ModifyNeed(c, constslib.NeedAdvance, boost)
			creature.ModifyNeed(c, constslib.NeedGuard, -boost)
			creature.ModifyNeed(c, constslib.NeedRetreat, -boost*0.2)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.3)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.4)
			creature.ModifyNeed(c, constslib.NeedFake, -boost*0.1)
			log.Printf("%s", color.New(color.FgHiCyan).Sprintf(
				"[FEEDBACK] [%s] parry realizado → +%.1f Guard, +%.1f Advance, -%.1f Retreat, +%.1f Rage, +%.1f Plan, -%.1f Fake",
				c.Handle.String(), boost, boost*0.2, boost*0.2, boost*0.3, boost*0.4, boost*0.1))

		case "DamageReceived":
			boost := 7 + rand.Float64()*6                            // Total: ~7-13, média ~10
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost)    // Menor avanço ao levar dano
			creature.ModifyNeed(c, constslib.NeedGuard, boost)       // Guard no mesmo nível do DamageApplied
			creature.ModifyNeed(c, constslib.NeedRetreat, boost*0.4) // Recuo mais forte que no DamageApplied
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.4)    // Um pouco mais de rage do que DamageApplied
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.2)    // Menor planejamento sob pressão
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.1)    // Menos chance de blefe quando apanha
			log.Printf("%s", color.New(color.FgHiRed).Sprintf(
				"[FEEDBACK] [%s] recebeu dano → +%.1f Guard, +%.1f Retreat, -%.1f Advance, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost*0.4, boost*0.4, boost*0.4, boost*0.4, boost*0.2, boost*0.1))

		case "RetreatPerformed":
			boost := 5 + rand.Float64()*6 // Total: ~5-11, média ~8
			creature.ModifyNeed(c, constslib.NeedAdvance, -boost*0.3)
			creature.ModifyNeed(c, constslib.NeedGuard, boost*0.5)
			creature.ModifyNeed(c, constslib.NeedRetreat, boost)
			creature.ModifyNeed(c, constslib.NeedRage, boost*0.2)
			creature.ModifyNeed(c, constslib.NeedPlan, boost*0.3)
			creature.ModifyNeed(c, constslib.NeedFake, boost*0.1)
			log.Printf("%s", color.New(color.FgHiYellow).Sprintf(
				"[FEEDBACK] [%s] recuo realizado → +%.1f Retreat, +%.1f Guard, -%.1f Advance, +%.1f Rage, +%.1f Plan, +%.1f Fake",
				c.Handle.String(), boost, boost*0.5, boost*0.3, boost*0.2, boost*0.3, boost*0.1))

		}

	}

	// 📝 Analisa CombatState atual (reforço de need)
	// switch c.CombatState {
	// case constslib.CombatStateDodging:
	// 	creature.ModifyNeed(c, constslib.NeedAdvance, 10)
	// 	creature.ModifyNeed(c, constslib.NeedGuard, -5)

	// case constslib.CombatStateParrying, constslib.CombatStateBlocking:
	// 	creature.ModifyNeed(c, constslib.NeedGuard, 15)
	// 	creature.ModifyNeed(c, constslib.NeedAdvance, 5)

	// case constslib.CombatStateExecutingSkill:
	// 	creature.ModifyNeed(c, constslib.NeedAdvance, 20)
	// 	creature.ModifyNeed(c, constslib.NeedRage, 10)

	// case constslib.CombatStateStaggered, constslib.CombatStatePostureBroken:
	// 	creature.ModifyNeed(c, constslib.NeedRetreat, 25)
	// 	creature.ModifyNeed(c, constslib.NeedGuard, 15)
	// 	creature.ModifyNeed(c, constslib.NeedAdvance, -15)

	// case constslib.CombatStateRaging:
	// 	creature.ModifyNeed(c, constslib.NeedRage, -c.GetNeedValue(constslib.NeedRage))
	// 	log.Printf("%s", color.New(color.FgHiRed).Sprintf(
	// 		"[FEEDBACK] [%s] entrou em Raging → NeedRage zerado",
	// 		c.Handle.String()))
	// }

	// 🧠 Decide CombatState macro baseado nos Needs
	var newState constslib.CombatState
	var stateName string

	if c.GetNeedValue(constslib.NeedRetreat) > c.GetNeedThreshold(constslib.NeedRetreat) {
		newState = constslib.CombatStateFleeing
		stateName = "Fleeing"
	} else if c.GetNeedValue(constslib.NeedAdvance) > c.GetNeedThreshold(constslib.NeedAdvance) {
		newState = constslib.CombatStateAggressive
		stateName = "Aggressive"
	} else if c.GetNeedValue(constslib.NeedGuard) > c.GetNeedThreshold(constslib.NeedGuard) {
		newState = constslib.CombatStateDefensive
		stateName = "Defensive"
	} else if c.GetNeedValue(constslib.NeedPlan) > c.GetNeedThreshold(constslib.NeedPlan) {
		newState = constslib.CombatStateStrategic
		stateName = "Strategic"
	} else if c.GetNeedValue(constslib.NeedRage) > c.GetNeedThreshold(constslib.NeedRage) {
		newState = constslib.CombatStateRaging
		stateName = "Raging"
	}
	if newState != 0 && c.CombatState != newState {
		c.CombatState = newState
		log.Printf("%s", color.New(color.FgHiWhite).Sprintf(
			"[FEEDBACK] [%s] CombatState mudou para %s", c.Handle.String(), stateName))

		svcCtx.RegisterCombatBehavior(dynamic_context.CombatBehaviorEvent{
			SourceHandle: c.Handle,
			BehaviorType: "StateChangedTo" + stateName,
			Timestamp:    time.Now(),
		})
	}

	return core.StatusSuccess
}

func (n *CombatFeedbackNode) Reset() {}
