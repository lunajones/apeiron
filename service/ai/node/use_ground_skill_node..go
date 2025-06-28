package node

// import (
// 	"log"

// 	"github.com/lunajones/apeiron/lib/combat"
// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// type UseGroundSkillNode struct {
// 	SkillName string
// }

// func (n *UseGroundSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[AI] [%s (%s)] contexto inválido em UseGroundSkillNode", c.ID, c.PrimaryType)
// 		return core.StatusFailure
// 	}

// 	log.Printf("[AI] [%s (%s)] executando UseGroundSkillNode", c.ID, c.PrimaryType)

// 	skill, exists := combat.SkillRegistry[n.SkillName]
// 	if !exists {
// 		log.Printf("[AI] [%s (%s)] skill %s não encontrada", c.ID, c.PrimaryType, n.SkillName)
// 		return core.StatusFailure
// 	}

// 	if c.MentalState == consts.MentalStateAfraid {
// 		log.Printf("[AI] [%s (%s)] está com medo, recusando-se a usar %s", c.ID, c.PrimaryType, n.SkillName)
// 		return core.StatusFailure
// 	}

// 	var bestTargetPos position.Position
// 	var targetsInRange int

// 	for _, p := range svcCtx.GetServicePlayers(c.Position, skill.Range) {
// 		dist := position.CalculateDistance(c.Position, p.Position)
// 		if dist <= skill.Range {
// 			targetsInRange++
// 			bestTargetPos = p.Position
// 		}
// 	}

// 	for _, other := range svcCtx.GetServiceCreatures(c.Position, skill.Range) {
// 		if other.ID == c.ID || !other.IsAlive {
// 			continue
// 		}
// 		dist := position.CalculateDistance(c.Position, other.Position)
// 		if dist <= skill.Range {
// 			targetsInRange++
// 			bestTargetPos = other.Position
// 		}
// 	}

// 	if targetsInRange == 0 {
// 		log.Printf("[AI] [%s (%s)] não encontrou alvos próximos para usar %s", c.ID, c.PrimaryType, n.SkillName)
// 		return core.StatusFailure
// 	}

// 	hunger := c.GetNeedValue(consts.NeedHunger)
// 	if targetsInRange == 1 && !(hunger > 80 && c.HasTag(consts.TagPredator)) {
// 		log.Printf("[AI] [%s (%s)] preferiu guardar a skill %s para mais inimigos", c.ID, c.PrimaryType, n.SkillName)
// 		return core.StatusFailure
// 	}

// 	log.Printf("[AI] [%s (%s)] usando %s em posição (%.2f, %.2f, %.2f)", c.ID, c.PrimaryType, n.SkillName,
// 		bestTargetPos.X, bestTargetPos.Y, bestTargetPos.Z)

// 	combat.UseSkill(
// 		c,
// 		nil,
// 		bestTargetPos,
// 		n.SkillName,
// 		svcCtx.GetServiceCreatures(c.Position, skill.Range),
// 		svcCtx.GetServicePlayers(c.Position, skill.Range),
// 	)

// 	return core.StatusSuccess
// }
