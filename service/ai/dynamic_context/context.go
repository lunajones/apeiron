package dynamic_context

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
)

type AIServiceContext struct {
	NavMesh                *navmesh.NavMesh
	SpatialIndex           navmesh.SpatialIndex
	perHandleCachedTargets map[handle.EntityHandle][]model.Targetable
	CombatBehaviors        []CombatBehaviorEvent
}

// Construtor do contexto
func NewAIServiceContext(navMesh *navmesh.NavMesh, index navmesh.SpatialIndex) *AIServiceContext {
	return &AIServiceContext{
		NavMesh:                navMesh,
		SpatialIndex:           index,
		perHandleCachedTargets: make(map[handle.EntityHandle][]model.Targetable),
	}
}

func (ctx *AIServiceContext) GetCachedTargets(h handle.EntityHandle) []model.Targetable {
	return ctx.perHandleCachedTargets[h]
}

func (ctx *AIServiceContext) CacheFor(handle handle.EntityHandle, center position.Position, detectionRadius float64) {
	ctx.perHandleCachedTargets[handle] = ctx.SpatialIndex.Query(center, detectionRadius)
}

func (ctx *AIServiceContext) RegisterCombatBehavior(event CombatBehaviorEvent) {
	ctx.CombatBehaviors = append(ctx.CombatBehaviors, event)
}

func (ctx *AIServiceContext) GetRecentCombatBehaviors(target handle.EntityHandle, since time.Time) []CombatBehaviorEvent {
	var result []CombatBehaviorEvent
	for _, e := range ctx.CombatBehaviors {
		if e.SourceHandle.Equals(target) && e.Timestamp.After(since) {
			result = append(result, e)
		}
	}
	return result
}

func (ctx *AIServiceContext) GetRecentCombatBehaviorsAsTarget(target handle.EntityHandle, since time.Time) []CombatBehaviorEvent {
	var result []CombatBehaviorEvent
	for _, e := range ctx.CombatBehaviors {
		if e.TargetHandle.Equals(target) && e.Timestamp.After(since) {
			result = append(result, e)
		}
	}
	return result
}
func (ctx *AIServiceContext) GetRecentAggressorsAgainst(target handle.EntityHandle, since time.Time) []CombatBehaviorEvent {
	var events []CombatBehaviorEvent

	for _, e := range ctx.CombatBehaviors {
		log.Printf(
			"[FILTER] analisando evento: Target=%s | Type=%s | Timestamp=%v | since=%v | MatchTarget=%v | MatchType=%v | PassaBefore=%v",
			e.TargetHandle.String(),
			e.BehaviorType,
			e.Timestamp.Format("15:04:05.000"),
			since.Format("15:04:05.000"),
			e.TargetHandle.Equals(target),
			e.BehaviorType == "AggressiveIntention",
			!e.Timestamp.Before(since),
		)

		if e.TargetHandle.Equals(target) &&
			e.BehaviorType == "AggressiveIntention" &&
			!e.Timestamp.Before(since) {
			events = append(events, e)
		}
	}

	return events
}

// Busca alvo específico pelo handle no índice
func (ctx *AIServiceContext) FindByHandle(h handle.EntityHandle) model.Targetable {
	candidates := ctx.SpatialIndex.Query(position.Position{}, 1e6) // Busca todos (ou ajuste se tiver melhor API)
	for _, t := range candidates {
		if t.GetHandle().Equals(h) {
			return t
		}
	}
	return nil
}
