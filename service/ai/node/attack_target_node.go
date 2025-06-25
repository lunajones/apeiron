package node

import (
	"log"

	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type AttackTargetNode struct {
	AttackSkill string
}

func (n *AttackTargetNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx := ctx.(dynamic_context.AIServiceContext)

	log.Printf("[AI] %s executando AttackTargetNode", c.ID)

	if c.TargetCreatureID == "" {
		log.Printf("[AI] %s não tem alvo.", c.ID)
		return core.StatusFailure
	}

	var target *creature.Creature
	for _, other := range svcCtx.GetServiceCreatures() {
		if other.ID == c.TargetCreatureID && other.IsAlive {
			target = other
			break
		}
	}

	if target == nil {
		log.Printf("[AI] %s não encontrou o alvo %s ou ele está morto.", c.ID, c.TargetCreatureID)
		c.TargetCreatureID = ""
		return core.StatusFailure
	}

	dist := position.CalculateDistance(c.Position, target.Position)
	if dist > c.AttackRange {
		log.Printf("[AI] %s está fora do alcance de ataque (%.2f > %.2f).", c.ID, dist, c.AttackRange)
		c.MoveTowards(target.Position, c.MoveSpeed)
		c.SetAction(creature.ActionRun)
		return core.StatusRunning
	}

	log.Printf("[AI] %s ataca %s com %s!", c.ID, target.ID, n.AttackSkill)

	// Corrigido: usa UseSkill do combat
	combat.UseSkill(c, target, target.Position, n.AttackSkill, svcCtx.GetServiceCreatures(), svcCtx.GetServicePlayers())

	c.SetAction(creature.ActionSkill1)
	c.ChangeAIState(creature.AIStateAttack)

	return core.StatusSuccess
}
