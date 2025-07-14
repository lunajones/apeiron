package node

import (
	"log"

	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type RegenerateNeedNode struct {
	NeedType            consts.NeedType
	CompletionThreshold float64               // ponto para "parar de regenerar" (ou mudar de estado, opcional)
	RegenAmount         float64               // valor por tick
	OnCompleteAI        consts.AIState        // opcional: AIState ao finalizar, ex: Idle
	OnCompleteAnim      consts.AnimationState // opcional: animação ao finalizar, ex: Wake
	RunningAnim         consts.AnimationState // opcional: animação enquanto regenera, ex: Sleep
}

func (n *RegenerateNeedNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	need := c.GetNeedByType(n.NeedType)

	if need == nil {
		log.Printf("[AI] [%s (%s)] erro: necessidade %s não encontrada", c.Handle.String(), c.PrimaryType, n.NeedType)
		return core.StatusFailure
	}

	if need.Value <= n.CompletionThreshold {
		log.Printf("[AI] [%s (%s)] %s recuperado (%.2f ≤ %.2f), AIState: %s",
			c.Handle.String(), c.PrimaryType, n.NeedType, need.Value, n.CompletionThreshold, n.OnCompleteAI)

		if n.OnCompleteAnim != "" {
			c.SetAnimationState(n.OnCompleteAnim)
		}
		if n.OnCompleteAI != "" {
			c.ChangeAIState(n.OnCompleteAI)
		}
		return core.StatusSuccess
	}

	creature.ModifyNeed(c, n.NeedType, n.RegenAmount)
	// log.Printf("[AI] [%s (%s)] regenerando %s. Novo valor: %.2f",
	// 	c.Handle.String(), c.PrimaryType, n.NeedType, c.GetNeedValue(n.NeedType))

	if n.RunningAnim != "" {
		c.SetAnimationState(n.RunningAnim)
	}

	return core.StatusRunning
}

func (n *RegenerateNeedNode) Reset(c *creature.Creature) {}
