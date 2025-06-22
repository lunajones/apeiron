package combat

import "log"

func Init() {
    log.Println("Combat service initialized")
}

func TickAll() {
    // atualizações como cooldowns e efeitos
}

func ResolveAttack(attacker, defender, ability string) {
    log.Printf("Resolving attack: %s -> %s with %s", attacker, defender, ability)
}
