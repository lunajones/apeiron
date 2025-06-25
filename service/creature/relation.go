package creature

func AreEnemies(a, b *Creature) bool {
	// Se forem da mesma facção, não são inimigos
	if a.Faction == b.Faction {
		return false
	}

	// Se qualquer um dos dois for não-hostil, assume-se que não são inimigos
	if !a.IsHostile || !b.IsHostile {
		return false
	}

	return true
}
