package ai

import (
	"log"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

// Nova estrutura com TargetTag incluído
type BehaviorRule struct {
	ObserverType consts.CreatureType
	TargetType   consts.CreatureType
	TargetTag    consts.CreatureTag
	Reaction     func(observer, target *creature.Creature)
}

var behaviorRules []BehaviorRule

func InitBehaviorRules() {
	behaviorRules = []BehaviorRule{
		// Exemplo 1: Wolf caça qualquer criatura com tag Prey
		{
			ObserverType: consts.Wolf,
			TargetTag:    consts.TagPrey,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Lobo [%s] vê [%s] como presa.",
					observer.GetHandle().ID, target.GetHandle().ID)
				observer.TargetCreatureHandle = target.GetHandle()
				observer.ChangeAIState(constslib.AIStateAlert)
			},
		},
		// Exemplo 2: Soldier odeia Wolf
		{
			ObserverType: consts.Soldier,
			TargetType:   consts.Wolf,
			Reaction: func(observer, target *creature.Creature) {
				log.Printf("[BehaviorRule] Soldado [%s] hostiliza o Lobo [%s].",
					observer.GetHandle().ID, target.GetHandle().ID)
				observer.TargetCreatureHandle = target.GetHandle()
				observer.ChangeAIState(constslib.AIStateAlert)
			},
		},
		// No futuro: mais regras...
	}
}

func EvaluateBehavior(observer, target *creature.Creature) {
	for _, rule := range behaviorRules {
		// Verifica se Observer bate com o tipo esperado
		typeMatch := rule.ObserverType == consts.AnyType || observer.HasType(rule.ObserverType)

		// Verifica se Target bate com o tipo ou tag esperada
		targetTypeMatch := rule.TargetType == consts.AnyType || target.HasType(rule.TargetType)
		targetTagMatch := rule.TargetTag == "" || target.HasTag(rule.TargetTag)

		if typeMatch && (targetTypeMatch || targetTagMatch) {
			rule.Reaction(observer, target)
		}
	}
}
