package offensive

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type CheckCombatStateNode struct {
	Expected []constslib.CombatState
}

func (n *CheckCombatStateNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	for _, state := range n.Expected {
		if c.CombatState&state != 0 || (c.CombatState == 0 && state == 0) {
			log.Printf("[CHECK-COMBAT-STATE] [%s (%s)] estados atuais = %b, bateu com esperado = %b (OK)",
				c.Handle.String(), c.PrimaryType, c.CombatState, state)
			return core.StatusSuccess
		}
	}

	log.Printf("[CHECK-COMBAT-STATE] [%s (%s)] estados atuais = %b, nenhum esperado bateu (FAIL)",
		c.Handle.String(), c.PrimaryType, c.CombatState)
	return core.StatusFailure
}

func (n *CheckCombatStateNode) Reset() {
	// Nada a resetar
}
