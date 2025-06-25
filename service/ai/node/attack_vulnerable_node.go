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

func (n *AttackIfVulnerableNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
    svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
    if !ok {
        log.Printf("[AI] %s: contexto inválido para AttackIfVulnerableNode", c.ID)
        return core.StatusFailure
    }

    if c.TargetCreatureID == "" {
        return core.StatusFailure
    }
    target := creature.FindServiceByID(svcCtx.GetServiceCreatures(), c.TargetCreatureID)
    if target == nil || !target.IsAlive {
        return core.StatusFailure
    }

    hpPercent := float64(target.HP) / float64(target.MaxHP) * 100
    if hpPercent > 30 {
        return core.StatusFailure
    }

    if c.MentalState == creature.MentalStateAfraid && c.MentalState != creature.MentalStateEnraged {
        return core.StatusFailure
    }

    hunger := c.GetNeedValue(creature.NeedHunger)
    if hunger > 80 && c.HasTag(creature.TagPredator) {
        combat.UseSkill(c, target, target.Position, n.SkillName, svcCtx.GetServiceCreatures(), svcCtx.GetServicePlayers())
        return core.StatusSuccess
    }

    if c.MentalState == creature.MentalStateAggressive || c.MentalState == creature.MentalStateEnraged {
        combat.UseSkill(c, target, target.Position, n.SkillName, svcCtx.GetServiceCreatures(), svcCtx.GetServicePlayers())
        return core.StatusSuccess
    }

    return core.StatusFailure
}
