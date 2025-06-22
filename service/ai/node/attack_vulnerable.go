package node

import "github.com/lunajones/apeiron/service/creature"

type AttackIfEnemyVulnerableNode struct{}

func (a *AttackIfEnemyVulnerableNode) Tick(c *creature.Creature) BehaviorStatus {
	// TODO: Substituir por lógica real de detecção de vulnerabilidade
	if c.TargetPlayerID != "" {
		c.SetAction(creature.ActionAttack)
		c.ChangeAIState(creature.AIStateIdle)
		return StatusSuccess
	}
	return StatusFailure
}
