package helper

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
)

type ValidateCombatStateNode struct{}

func (n *ValidateCombatStateNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.GetCombatState() == constslib.CombatStateCasting {
		state := c.CurrentSkillState()
		if state == nil || state.WasInterrupted {
			log.Printf("[VALIDATOR] [%s] travado em casting com estado inválido ou interrompido. Forçando reset", c.Handle.String())

			if c.NextSkillToUse != nil {
				skillState := c.SkillStates[c.NextSkillToUse.Action]
				if skillState != nil {
					// skillState.InUse = false
					skillState.WindUpFired = false
					skillState.CastFired = false
					skillState.RecoveryFired = false
					skillState.WasInterrupted = false
				}
				c.NextSkillToUse = nil
			}
			c.SetCombatState(constslib.CombatStateMoving)
		}

	}
	return core.StatusSuccess
}

func (n *ValidateCombatStateNode) Reset(c *creature.Creature) {
}
