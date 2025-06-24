package node

import (
	"log"

	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/service/creature"
)

type AttackIfVulnerableNode struct {
	SkillName string
}

func (n *AttackIfVulnerableNode) Tick(c *creature.Creature, ctx dynamic_context.AIServiceContext) core.BehaviorStatus {
	if c.TargetCreatureID == "" {
		log.Printf("[AI] %s não tem target para avaliar vulnerabilidade.", c.ID)
		return core.StatusFailure
	}

	target := creature.FindServiceByID(ctx.GetServiceCreatures(), c.TargetCreatureID)
	if target == nil || !target.IsAlive {
		log.Printf("[AI] %s: Target %s inválido ou morto.", c.ID, c.TargetCreatureID)
		return core.StatusFailure
	}

	// Calcular percentual de HP do alvo
	hpPercent := (float64(target.HP) / 100.0) * 100
	if hpPercent > 30 {
		log.Printf("[AI] %s: Target %s não está vulnerável (HP: %.2f%%).", c.ID, target.ID, hpPercent)
		return core.StatusFailure
	}

	// Regra: Se criatura está com medo e não enraivecida, recusa atacar
	if c.MentalState == creature.MentalStateAfraid && c.MentalState != creature.MentalStateEnraged {
		log.Printf("[AI] %s está com medo, ignorando target vulnerável.", c.ID)
		return core.StatusFailure
	}

	// Regra: Se a criatura tem fome e é predador, ataca alvos vulneráveis com mais frequência
	hunger := c.GetNeedValue(creature.NeedHunger)
	if hunger > 80 && c.HasTag(creature.TagPredator) {
		log.Printf("[AI] %s faminto, atacando alvo vulnerável %s.", c.ID, target.ID)
		combat.UseSkill(c, target, target.Position, n.SkillName, ctx.GetServiceCreatures(), ctx.GetServicePlayers())
		return core.StatusSuccess
	}

	// Regra final: Se está agressivo ou enraivecido, ataca mesmo sem estar com fome
	if c.MentalState == creature.MentalStateAggressive || c.MentalState == creature.MentalStateEnraged {
		log.Printf("[AI] %s agressivo/enraivecido, atacando %s.", c.ID, target.ID)
		combat.UseSkill(c, target, target.Position, n.SkillName, ctx.GetServiceCreatures(), ctx.GetServicePlayers())
		return core.StatusSuccess
	}

	log.Printf("[AI] %s decidiu não atacar o alvo vulnerável %s.", c.ID, target.ID)
	return core.StatusFailure
}
