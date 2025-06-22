package combat

import (
	"math"

	"github.com/lunajones/apeiron/service/creature"
)

// Calcula dano físico direto (espadas, machados, etc)
func CalculatePhysicalDamage(attacker *creature.Creature, target *creature.Creature, skillMultiplier float64) int {
	baseAttack := attacker.Strength
	dexBonus := float64(attacker.Dexterity) * 0.1 // Pequeno bônus por destreza
	rawDamage := (float64(baseAttack) + dexBonus) * skillMultiplier

	// Aplicar defesa física do alvo
	finalDamage := rawDamage * (1 - target.PhysicalDefense)

	if finalDamage < 1 {
		finalDamage = 1
	}

	return int(math.Round(finalDamage))
}

// Calcula dano mágico (feitiços, projéteis mágicos)
func CalculateMagicDamage(attacker *creature.Creature, target *creature.Creature, skillMultiplier float64) int {
	baseMagic := attacker.Intelligence
	focusBonus := float64(attacker.Focus) * 0.05 // Pequeno bônus por Focus
	rawDamage := (float64(baseMagic) + focusBonus) * skillMultiplier

	// Aplicar defesa mágica do alvo
	finalDamage := rawDamage * (1 - target.MagicDefense)

	if finalDamage < 1 {
		finalDamage = 1
	}

	return int(math.Round(finalDamage))
}

// Calcula dano de Poison (DOT)
func CalculatePoisonDamage(attacker *creature.Creature, target *creature.Creature) int {
	base := 5
	strBonus := float64(attacker.Strength) * 0.2
	intBonus := float64(attacker.Intelligence) * 0.1
	rawDOT := float64(base) + strBonus + intBonus

	// Aplicar StatusResistance do alvo (reduz DOT)
	finalDOT := rawDOT * (1 - target.StatusResistance)

	if finalDOT < 1 {
		finalDOT = 1
	}

	return int(math.Round(finalDOT))
}

// Calcula dano de Burn (DOT)
func CalculateBurnDamage(attacker *creature.Creature, target *creature.Creature) int {
	base := 7
	intBonus := float64(attacker.Intelligence) * 0.3
	rawDOT := float64(base) + intBonus

	finalDOT := rawDOT * (1 - target.StatusResistance)

	if finalDOT < 1 {
		finalDOT = 1
	}

	return int(math.Round(finalDOT))
}

// Calcula cura baseada em Focus e Intelligence (para habilidades de cura)
func CalculateHealing(attacker *creature.Creature, skillMultiplier float64) int {
	baseHeal := attacker.Focus*2 + attacker.Intelligence
	finalHeal := float64(baseHeal) * skillMultiplier

	if finalHeal < 1 {
		finalHeal = 1
	}

	return int(math.Round(finalHeal))
}

// Calcula duração efetiva de um CC, considerando ControlResistance
func CalculateEffectiveCCDuration(baseDuration float64, target *creature.Creature) float64 {
	reduction := baseDuration * target.ControlResistance
	finalDuration := baseDuration - reduction

	if finalDuration < 0.1 {
		finalDuration = 0.1 // Mínimo de 0.1s pra evitar CC zero
	}

	return finalDuration
}
