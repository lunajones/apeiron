package entityresolver

import (
	"github.com/lunajones/apeiron/lib/handle"
	"github.com/lunajones/apeiron/service/zone"
)

// Evita importar creature ou player. Resolve direto via zona global
func ResolveEntityByHandle(h handle.EntityHandle) any {
	for _, z := range zone.Zones {
		for _, c := range z.GetCreatures() {
			if c.GetHandle().Equals(h) {
				return c
			}
		}
		for _, p := range z.GetPlayers() {
			if p.GetHandle().Equals(h) {
				return p
			}
		}
	}
	return nil
}
