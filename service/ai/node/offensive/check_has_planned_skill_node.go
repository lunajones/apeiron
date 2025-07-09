package offensive

// NAO NECESSARIO PQ O CORNO DO CHATGPT ME FEZ CRIAR, DIZENDO QUE ERA USADO NO DEATHMARTH, MAS DPS DISSE QUE NAO ERA.
// import (
// 	"log"

// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// )

// type CheckHasPlannedSkillNode struct{}

// func (n *CheckHasPlannedSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	_, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[CheckHasPlannedSkillNode] contexto inválido para criatura %s", c.Name)
// 		return core.StatusFailure
// 	}

// 	if c.NextComboSkillToUse == nil {
// 		log.Printf("[CHECK-PLANNED-SKILL] [%s] Nenhuma skill de combo planejada", c.Handle)
// 		return core.StatusFailure
// 	}

// 	log.Printf("[CHECK-PLANNED-SKILL] [%s] Skill de combo planejada: %s", c.Handle, c.NextComboSkillToUse.Name)
// 	return core.StatusSuccess
// }

// func (n *CheckHasPlannedSkillNode) Reset() {
// 	// Nada a resetar
// }
