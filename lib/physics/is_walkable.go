package physics

import (
	"github.com/lunajones/apeiron/lib/position"
)

// IsWalkable verifica se o ponto pode ser pisado por uma criatura (para IA, pathfinding etc).
func IsWalkable(pos position.Position, hitboxRadius float64) bool {
	// Primeiro, colisão com o mundo
	if CheckWorldCollision(pos, hitboxRadius) {
		return false
	}

	// Aqui podemos expandir no futuro:
	// - verificar altura do terreno
	// - verificar se é água, lava etc.
	// - checar se área é "fora do mapa"

	return true
}
