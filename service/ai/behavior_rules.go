package ai

import (
	"github.com/lunajones/apeiron/service/creature"
	"log"
)

type BehaviorRule struct {
	ObserverType creature.Type
	TargetType   creature.Type
	Reaction        func(observer, target *creature.Creature)
}

var behaviorRules []BehaviorRule

func InitBehaviorRules() {
	behaviorRules = []BehaviorRule{
		{
			ObserverType: creature.Acolyte,
			TargetType:   creature.Zombie,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Acolyte %s tenta dominar Zombie %s", observer.ID, target.ID)
				// TODO: Implementar lógica de dominação
			},
		},
		{
			ObserverType: creature.Wolf,
			TargetType:   creature.Rabbit,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Wolf %s ataca Rabbit %s", observer.ID, target.ID)
				// TODO: Implementar lógica de ataque
			},
		},
	}
}

func EvaluateBehavior(observer, target *creature.Creature) {
	for _, rule := range behaviorRules {
		if rule.ObserverType == observer.Type && rule.TargetType == target.Type {
			rule.Reaction(observer, target)
		}
	}
}
