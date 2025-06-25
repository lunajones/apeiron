package node

import (
    "log"

    "github.com/lunajones/apeiron/service/ai/core"
    "github.com/lunajones/apeiron/service/ai/dynamic_context"
    "github.com/lunajones/apeiron/lib/combat"
    "github.com/lunajones/apeiron/lib/position"
    "github.com/lunajones/apeiron/service/creature"
)

type AttackTargetNode struct {
    SkillName string
}

func (n *AttackTargetNode) Tick(c *creature.Creature, ctx dynamic_context.AIServiceContext) core.BehaviorStatus {
    log.Printf("[AI] %s executando AttackTargetNode", c.ID)

    if c.TargetCreatureID == "" {
        log.Printf("[AI] %s não tem alvo para atacar.", c.ID)
        return core.StatusFailure
    }

    target := creature.FindServiceByID(ctx.GetServiceCreatures(), c.TargetCreatureID)
    if target == nil || !target.IsAlive {
        log.Printf("[AI] %s: Target inválido ou morto.", c.ID)
        return core.StatusFailure
    }

    // Regra mental
    if c.MentalState == creature.MentalStateAfraid {
        log.Printf("[AI] %s está com medo, recusando-se a atacar.", c.ID)
        return core.StatusFailure
    }

    // Fome extrema
    hunger := c.GetNeedValue(creature.NeedHunger)
    if hunger <= 90 {
        // Se for animal, só atacar se target for presa
        if c.HasTag(creature.TagAnimal) && !target.HasTag(creature.TagPrey) {
            log.Printf("[AI] %s é animal e alvo %s não é presa.", c.ID, target.ID)
            return core.StatusFailure
        }
    }

    distance := position.CalculateDistance(c.Position, target.Position)
    skill, exists := combat.SkillRegistry[n.SkillName]
    if !exists {
        log.Printf("[AI] Skill %s não encontrada para %s.", n.SkillName, c.ID)
        return core.StatusFailure
    }
    if distance > skill.Range {
        log.Printf("[AI] %s: alvo %s fora de alcance da skill %s.", c.ID, target.ID, n.SkillName)
        return core.StatusFailure
    }

    combat.UseSkill(c, target, target.Position, n.SkillName, ctx.GetServiceCreatures(), ctx.GetServicePlayers())
    log.Printf("[AI] %s atacou %s com %s.", c.ID, target.ID, n.SkillName)
    return core.StatusSuccess
}
