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
