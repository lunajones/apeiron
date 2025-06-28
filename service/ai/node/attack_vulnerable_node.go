package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/combat"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type AttackIfVulnerableNode struct {
// 	SkillName string
// }

// func (n *AttackIfVulnerableNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[AI] [%s (%s)] contexto inválido para AttackIfVulnerableNode", c.ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	if c.TargetCreatureID == "" {
// 		return core.StatusFailure
// 	}

// 	target := creature.FindServiceByID(svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius), c.TargetCreatureID)
// 	if target == nil || !target.IsAlive {
// 		return core.StatusFailure
// 	}

// 	hpPercent := float64(target.HP) / float64(target.MaxHP) * 100
// 	if hpPercent > 30 {
// 		log.Printf("[AI] [%s (%s)] considera %s (%s) não vulnerável (%.2f%% HP)", c.ID, c.PrimaryType, target.ID, target.PrimaryType, hpPercent)
// 		return core.StatusFailure
// 	}

// 	if c.MentalState == consts.MentalStateAfraid && c.MentalState != consts.MentalStateEnraged {
// 		log.Printf("[AI] [%s (%s)] não ataca devido ao estado mental: %s", c.ID, c.PrimaryType, c.MentalState)
// 		return core.StatusFailure
// 	}

// 	hunger := c.GetNeedValue(consts.NeedHunger)

// 	if hunger > 80 && c.HasTag(consts.TagPredator) {
// 		log.Printf("[AI] [%s (%s)] ataca [%s (%s)] com %s (fome: %d%%, alvo vulnerável: %.2f%% HP)",
// 			c.ID, c.PrimaryType, target.ID, target.PrimaryType, n.SkillName, hunger, hpPercent)

// 		combat.UseSkill(
// 			c,
// 			target,
// 			target.Position,
// 			n.SkillName,
// 			svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius),
// 			svcCtx.GetServicePlayers(c.Position, c.DetectionRadius),
// 		)

// 		return core.StatusSuccess
// 	}

// 	if c.MentalState == consts.MentalStateAggressive || c.MentalState == consts.MentalStateEnraged {
// 		log.Printf("[AI] [%s (%s)] ataca [%s (%s)] com %s (estado mental: %s, alvo vulnerável: %.2f%% HP)",
// 			c.ID, c.PrimaryType, target.ID, target.PrimaryType, n.SkillName, c.MentalState, hpPercent)

// 		combat.UseSkill(
// 			c,
// 			target,
// 			target.Position,
// 			n.SkillName,
// 			svcCtx.GetServiceCreatures(c.Position, c.DetectionRadius),
// 			svcCtx.GetServicePlayers(c.Position, c.DetectionRadius),
// 		)

// 		return core.StatusSuccess
// 	}

// 	return core.StatusFailure
// }
