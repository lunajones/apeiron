package node

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type InterruptIfAttackedRecentlyNode struct {
	InterruptAIState   constslib.AIState
	InterruptAnimation constslib.AnimationState
}

func (n *InterruptIfAttackedRecentlyNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if !c.ReceivedDamageRecently {
		return core.StatusFailure
	}

	log.Printf("[AI-INTERRUPT-DAMAGE] [%s (%s)] foi atacado recentemente. Interrompendo.",
		c.Handle.String(), c.PrimaryType)

	c.LastThreatSeen = time.Now()

	if n.InterruptAnimation != "" {
		c.SetAnimationState(n.InterruptAnimation)
	}
	if n.InterruptAIState != "" {
		c.ChangeAIState(n.InterruptAIState)
	}
	c.ReceivedDamageRecently = false // Consome o evento

	return core.StatusSuccess
}

func (n *InterruptIfAttackedRecentlyNode) Reset() {}
