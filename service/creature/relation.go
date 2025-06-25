package creature

func AreEnemies(a, b *Creature) bool {
	// Mesmo grupo ou não-hostis
	if a.Faction == b.Faction || !a.IsHostile || !b.IsHostile {
		return false
	}

	// Predadores caçam presas
	if a.HasTag(TagPredator) && b.HasTag(TagPrey) {
		return true
	}

	// Predadores atacam humanoides se estiverem com fome
	if a.HasTag(TagPredator) && b.HasTag(TagHumanoid) {
		hunger := a.GetNeedValue(NeedHunger)
		if hunger > 70 {
			return true
		}
	}

	// Presas não atacam nada — não são inimigas ativamente
	if a.HasTag(TagPrey) {
		return false
	}

	// Humanoides se defendem
	if a.HasTag(TagHumanoid) && b.HasTag(TagPredator) {
		return true
	}

	// Fallback padrão: inimigos se forem hostis e diferentes
	return true
}
