// effect_helper.go
package helper

import (
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/creature"
)

func HasDOTEffectOfType(target *creature.Creature, effectType constslib.EffectType) bool {
	for _, eff := range target.ActiveEffects {
		if eff.IsDOT && eff.Type == effectType {
			return true
		}
	}
	return false
}
