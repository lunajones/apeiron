package ai

import (
	"log"

	"github.com/lunajones/apeiron/service/creature"
)

// Nova estrutura com TargetTag incluído
type BehaviorRule struct {
	ObserverType creature.CreatureType
	TargetType   creature.CreatureType
	TargetTag    creature.CreatureTag
	Reaction     func(observer, target *creature.Creature)
}

var behaviorRules []BehaviorRule

func InitBehaviorRules() {
	behaviorRules = []BehaviorRule{
		// Exemplo 1: Wolf caça qualquer criatura com tag Prey
		{
			ObserverType: creature.Wolf,
			TargetTag:    creature.TagPrey,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Lobo %s vê %s como presa.", observer.ID, target.ID)
				observer.TargetCreatureID = target.ID
				observer.ChangeAIState(creature.AIStateAlert)
			},
		},
		// Exemplo 2: Soldier odeia Wolf
		{
			ObserverType: creature.Soldier,
			TargetType:   creature.Wolf,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Soldado %s hostiliza o Lobo %s.", observer.ID, target.ID)
				observer.TargetCreatureID = target.ID
				observer.ChangeAIState(creature.AIStateAlert)
			},
		},
		// No futuro: mais regras...
	}
}

func EvaluateBehavior(observer, target *creature.Creature) {
	for _, rule := range behaviorRules {
		// Verifica se Observer bate com o tipo esperado
		typeMatch := rule.ObserverType == creature.AnyType || observer.HasType(rule.ObserverType)

		// Verifica se Target bate com o tipo ou tag esperada
		targetTypeMatch := rule.TargetType == creature.AnyType || target.HasType(rule.TargetType)
		targetTagMatch := rule.TargetTag == "" || target.HasTag(rule.TargetTag)

		if typeMatch && (targetTypeMatch || targetTagMatch) {
			rule.Reaction(observer, target)
		}
	}
}
