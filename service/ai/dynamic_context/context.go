package dynamic_context

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
)

type AIServiceContext struct {
	NavMesh                *navmesh.NavMesh
	SpatialIndex           navmesh.SpatialIndex
	perHandleCachedTargets map[handle.EntityHandle][]model.Targetable
	claimedPositions       map[string]handle.EntityHandle
}

// Construtor do contexto
func NewAIServiceContext(navMesh *navmesh.NavMesh, index navmesh.SpatialIndex) *AIServiceContext {
	return &AIServiceContext{
		NavMesh:                navMesh,
		SpatialIndex:           index,
		perHandleCachedTargets: make(map[handle.EntityHandle][]model.Targetable),
		claimedPositions:       make(map[string]handle.EntityHandle),
	}
}

func (ctx *AIServiceContext) GetCachedTargets(h handle.EntityHandle) []model.Targetable {
	return ctx.perHandleCachedTargets[h]
}

func (ctx *AIServiceContext) CacheFor(handle handle.EntityHandle, center position.Position, detectionRadius float64) {
	ctx.perHandleCachedTargets[handle] = ctx.SpatialIndex.Query(center, detectionRadius)
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

func (ctx *AIServiceContext) ClaimPosition(pos position.Position, h handle.EntityHandle) bool {
	key := pos.Key() // tipo "12:7"
	if existing, ok := ctx.claimedPositions[key]; ok && existing != h {
		return false
	}
	ctx.claimedPositions[key] = h
	return true
}

func (ctx *AIServiceContext) ClearClaims(h handle.EntityHandle) {
	for k, v := range ctx.claimedPositions {
		if v == h {
			delete(ctx.claimedPositions, k)
		}
	}
}

func (ctx *AIServiceContext) IsClaimedByOther(pos position.Position, h handle.EntityHandle) bool {
	key := pos.Key()
	if existing, ok := ctx.claimedPositions[key]; ok && existing != h {
		return true
	}
	return false
}

func (ctx *AIServiceContext) GetClaimer(pos position.Position) (handle.EntityHandle, bool) {
	h, ok := ctx.claimedPositions[pos.Key()]
	return h, ok
}
