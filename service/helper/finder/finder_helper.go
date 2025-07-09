package finder

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
)

func FindTargetByHandles(
	searcherHandle, creatureHandle, playerHandle handle.EntityHandle,
	svcCtx *dynamic_context.AIServiceContext,
) model.Targetable {
	targets := svcCtx.GetCachedTargets(searcherHandle)

	for _, t := range targets {
		if t.GetHandle().Equals(creatureHandle) && t.IsAlive() {
			return t
		}
		if t.GetHandle().Equals(playerHandle) && t.IsAlive() {
			return t
		}
	}
	return nil
}
