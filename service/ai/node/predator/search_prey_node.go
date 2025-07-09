package predator

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type SearchPreyNode struct {
	TargetTags []consts.CreatureTag
}

func (n *SearchPreyNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[SEARCH-PREY] [%s (%s)] contexto inválido", c.Handle.String(), c.PrimaryType)
		return core.StatusFailure
	}

	for _, t := range svcCtx.GetCachedTargets(c.Handle) {
		other, ok := t.(*creature.Creature)
		if !ok || !other.Alive || other.Handle.Equals(c.Handle) {
			continue
		}
		if !hasAnyTag(other, n.TargetTags) {
			continue
		}

		c.TargetCreatureHandle = other.Handle
		c.ChangeAIState(constslib.AIStateChasing)
		log.Printf("[SEARCH-PREY] [%s (%s)] prey encontrado: [%s], mudando para Chasing",
			c.Handle.String(), c.PrimaryType, other.Handle.String())
		return core.StatusSuccess
	}

	log.Printf("[SEARCH-PREY] [%s (%s)] nenhuma prey encontrada", c.Handle.String(), c.PrimaryType)
	return core.StatusFailure
}

func hasAnyTag(other *creature.Creature, tags []consts.CreatureTag) bool {
	for _, requiredTag := range tags {
		if other.HasTag(requiredTag) {
			return true
		}
	}
	return false
}

func (n *SearchPreyNode) Reset() {
	// Nada a resetar
}
