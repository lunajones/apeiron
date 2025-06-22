package ai

import (
	"github.com/lunajones/apeiron/service/creature"
	"log"
)

type BehaviorRule struct {
	ObserverSubtype creature.CreatureSubtype
	TargetSubtype   creature.CreatureSubtype
	Reaction        func(observer, target *creature.Creature)
}

var behaviorRules []BehaviorRule

func InitBehaviorRules() {
	behaviorRules = []BehaviorRule{
		{
			ObserverSubtype: creature.Acolyte,
			TargetSubtype:   creature.Zombie,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Acolyte %s tenta dominar Zombie %s", observer.ID, target.ID)
				// TODO: Implementar lógica de dominação
			},
		},
		{
			ObserverSubtype: creature.Wolf,
			TargetSubtype:   creature.Rabbit,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Wolf %s ataca Rabbit %s", observer.ID, target.ID)
				// TODO: Implementar lógica de ataque
			},
		},
	}
}

func EvaluateBehavior(observer, target *creature.Creature) {
	for _, rule := range behaviorRules {
		if rule.ObserverSubtype == observer.Subtype && rule.TargetSubtype == target.Subtype {
			rule.Reaction(observer, target)
		}
	}
}
