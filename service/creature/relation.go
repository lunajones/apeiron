package creature

import (
	"github.com/lunajones/apeiron/service/creature/consts"
)

func AreEnemies(a, b *Creature) bool {
	// Mesmo grupo ou não-hostis
	if a.Faction == b.Faction || (!a.IsHostile && !b.IsHostile) {
		return false
	}

	// Se A é predador e B é presa
	if a.HasTag(consts.TagPredator) && b.HasTag(consts.TagPrey) {
		return true
	}

	// Se A é presa e B é predador — coelho detectando lobo!
	if a.HasTag(consts.TagPrey) && b.HasTag(consts.TagPredator) {
		return true
	}

	// Predadores atacam humanoides se estiverem com fome
	if a.HasTag(consts.TagPredator) && b.HasTag(consts.TagHumanoid) {
		hunger := a.GetNeedValue(consts.NeedHunger)
		if hunger > 70 {
			return true
		}
	}

	// Humanoides se defendem se veem predador
	if a.HasTag(consts.TagHumanoid) && b.HasTag(consts.TagPredator) {
		return true
	}

	// Fallback: criaturas de facções diferentes e com hostilidade ativada
	return true
}
