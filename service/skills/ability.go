package ability

import "log"

func UseAbility(player, target, name string) {
    log.Printf("%s uses %s on %s", player, name, target)
}