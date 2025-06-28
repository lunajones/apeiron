package physics

import (
	"log"

	"github.com/lunajones/apeiron/lib/position"
)

// CheckWorldCollision centraliza verificação de colisão no mundo.
// Retorna true se há colisão no ponto informado.
func CheckWorldCollision(pos position.Position, hitboxRadius float64) bool {
	// Checa colisão no próprio sistema físico do mundo
	if CheckCollision(pos, hitboxRadius) {
		log.Printf("[WORLD COLLISION] Colisão detectada em (%.2f, %.2f, %.2f)", pos.FastGlobalX(), pos.FastGlobalY(), pos.Z)
		return true
	}

	// Aqui podemos expandir no futuro para checar outras coisas (portas, triggers, etc.)
	return false
}
