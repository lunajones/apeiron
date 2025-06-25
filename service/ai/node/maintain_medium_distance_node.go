package node

import (
    "math"
    "log"

    "github.com/lunajones/apeiron/service/ai/core"
    "github.com/lunajones/apeiron/service/ai/dynamic_context"
    "github.com/lunajones/apeiron/service/creature"
)

type MaintainMediumDistanceNode struct{}

func (n *MaintainMediumDistanceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
    svcCtx := ctx.(dynamic_context.AIServiceContext)
    log.Printf("[AI] %s executando MaintainMediumDistanceNode", c.ID)

    for _, p := range svcCtx.GetServicePlayers() {
        dx := p.Position.X - c.Position.X
        dz := p.Position.Z - c.Position.Z
        distance := math.Sqrt(dx*dx + dz*dz)

        if distance < 4.0 || distance > 8.0 {
            c.SetAction(creature.ActionRun)
        } else {
            c.SetAction(creature.ActionSkill2)
            c.ChangeAIState(creature.AIStateAttack)
        }

        return core.StatusSuccess
    }

    return core.StatusFailure
}
